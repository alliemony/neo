package handler

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
	"github.com/alliemony/neo/backend/internal/service"
)

func setupRSSHandler(t *testing.T) (*RSSHandler, *service.PostService) {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("setup db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	repo := repository.NewPostRepo(db)
	svc := service.NewPostService(repo)
	h := NewRSSHandler(svc, "https://example.com")
	return h, svc
}

func rssRouter(h *RSSHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/v1/feed.xml", h.Feed)
	return r
}

func TestRSSHandler_Feed_ReturnsValidXML(t *testing.T) {
	h, svc := setupRSSHandler(t)
	svc.Create(model.CreatePostInput{Title: "First Post", Content: "Hello world", Tags: []string{"go"}, Published: true})
	svc.Create(model.CreatePostInput{Title: "Second Post", Content: "Another post", Tags: []string{"react"}, Published: true})

	r := rssRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/feed.xml", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/rss+xml") {
		t.Errorf("content-type = %q, want application/rss+xml", ct)
	}

	var feed rssFeed
	if err := xml.Unmarshal(w.Body.Bytes(), &feed); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}

	if feed.Version != "2.0" {
		t.Errorf("rss version = %q, want 2.0", feed.Version)
	}
	if feed.Channel.Title != "neo" {
		t.Errorf("channel title = %q, want neo", feed.Channel.Title)
	}
	if len(feed.Channel.Items) != 2 {
		t.Errorf("items = %d, want 2", len(feed.Channel.Items))
	}
}

func TestRSSHandler_Feed_ContainsPostDetails(t *testing.T) {
	h, svc := setupRSSHandler(t)
	svc.Create(model.CreatePostInput{Title: "My Post", Content: "Some content here", Tags: []string{"go"}, Published: true})

	r := rssRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/feed.xml", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var feed rssFeed
	if err := xml.Unmarshal(w.Body.Bytes(), &feed); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}

	if len(feed.Channel.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(feed.Channel.Items))
	}

	item := feed.Channel.Items[0]
	if item.Title != "My Post" {
		t.Errorf("item title = %q, want My Post", item.Title)
	}
	if item.Link != "https://example.com/blog/my-post" {
		t.Errorf("item link = %q, want https://example.com/blog/my-post", item.Link)
	}
	if item.Description != "Some content here" {
		t.Errorf("item description = %q, want Some content here", item.Description)
	}
	if item.PubDate == "" {
		t.Error("item pubDate is empty")
	}
	if item.GUID == "" {
		t.Error("item GUID is empty")
	}
}

func TestRSSHandler_Feed_ExcludesDrafts(t *testing.T) {
	h, svc := setupRSSHandler(t)
	svc.Create(model.CreatePostInput{Title: "Published", Content: "visible", Tags: []string{"go"}, Published: true})
	svc.Create(model.CreatePostInput{Title: "Draft", Content: "hidden", Tags: []string{"go"}, Published: false})

	r := rssRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/feed.xml", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var feed rssFeed
	if err := xml.Unmarshal(w.Body.Bytes(), &feed); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}

	if len(feed.Channel.Items) != 1 {
		t.Errorf("items = %d, want 1 (drafts should be excluded)", len(feed.Channel.Items))
	}
	if feed.Channel.Items[0].Title != "Published" {
		t.Errorf("item title = %q, want Published", feed.Channel.Items[0].Title)
	}
}

func TestRSSHandler_Feed_EmptyWhenNoPosts(t *testing.T) {
	h, _ := setupRSSHandler(t)

	r := rssRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/feed.xml", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var feed rssFeed
	if err := xml.Unmarshal(w.Body.Bytes(), &feed); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}

	if len(feed.Channel.Items) != 0 {
		t.Errorf("items = %d, want 0", len(feed.Channel.Items))
	}
}

func TestRSSHandler_Feed_TruncatesLongContent(t *testing.T) {
	h, svc := setupRSSHandler(t)
	longContent := strings.Repeat("a", 500)
	svc.Create(model.CreatePostInput{Title: "Long Post", Content: longContent, Tags: []string{"go"}, Published: true})

	r := rssRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/feed.xml", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var feed rssFeed
	if err := xml.Unmarshal(w.Body.Bytes(), &feed); err != nil {
		t.Fatalf("invalid XML: %v", err)
	}

	desc := feed.Channel.Items[0].Description
	// excerpt truncates to 300 chars + "…"
	if len(desc) > 305 {
		t.Errorf("description length = %d, want <= 305 (300 + ellipsis)", len(desc))
	}
}
