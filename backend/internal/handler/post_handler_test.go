package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
	"github.com/alliemony/neo/backend/internal/service"
)

func setupPostHandler(t *testing.T) (*PostHandler, *service.PostService) {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("setup db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	repo := repository.NewPostRepo(db)
	svc := service.NewPostService(repo)
	h := NewPostHandler(svc)
	return h, svc
}

func newChiRouter(h *PostHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/v1/posts", h.List)
	r.Get("/api/v1/posts/{slug}", h.GetBySlug)
	r.Get("/api/v1/tags", h.ListTags)
	return r
}

func TestPostHandler_List_ReturnsPosts(t *testing.T) {
	h, svc := setupPostHandler(t)
	svc.Create(model.CreatePostInput{Title: "Post One", Content: "content one", Tags: []string{"go"}, Published: true})
	svc.Create(model.CreatePostInput{Title: "Post Two", Content: "content two", Tags: []string{"python"}, Published: true})

	r := newChiRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("content-type = %q, want application/json", ct)
	}

	var resp struct {
		Posts []model.Post `json:"posts"`
		Total int          `json:"total"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Total != 2 {
		t.Errorf("total = %d, want 2", resp.Total)
	}
	if len(resp.Posts) != 2 {
		t.Errorf("len(posts) = %d, want 2", len(resp.Posts))
	}
}

func TestPostHandler_List_FilterByTag(t *testing.T) {
	h, svc := setupPostHandler(t)
	svc.Create(model.CreatePostInput{Title: "Go Post", Content: "c", Tags: []string{"go"}, Published: true})
	svc.Create(model.CreatePostInput{Title: "Python Post", Content: "c", Tags: []string{"python"}, Published: true})

	r := newChiRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/posts?tag=go", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp struct {
		Posts []model.Post `json:"posts"`
		Total int          `json:"total"`
	}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Total != 1 {
		t.Errorf("total = %d, want 1", resp.Total)
	}
	if len(resp.Posts) != 1 || resp.Posts[0].Title != "Go Post" {
		t.Errorf("unexpected posts: %v", resp.Posts)
	}
}

func TestPostHandler_List_ExcludesDrafts(t *testing.T) {
	h, svc := setupPostHandler(t)
	svc.Create(model.CreatePostInput{Title: "Published", Content: "c", Published: true})
	svc.Create(model.CreatePostInput{Title: "Draft", Content: "c", Published: false})

	r := newChiRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp struct {
		Posts []model.Post `json:"posts"`
		Total int          `json:"total"`
	}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Total != 1 {
		t.Errorf("total = %d, want 1", resp.Total)
	}
}

func TestPostHandler_GetBySlug_Found(t *testing.T) {
	h, svc := setupPostHandler(t)
	svc.Create(model.CreatePostInput{Title: "Hello World", Content: "first post", Tags: []string{"intro"}, Published: true})

	r := newChiRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/posts/hello-world", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var post model.Post
	json.NewDecoder(w.Body).Decode(&post)
	if post.Title != "Hello World" {
		t.Errorf("title = %q, want %q", post.Title, "Hello World")
	}
}

func TestPostHandler_GetBySlug_NotFound(t *testing.T) {
	h, _ := setupPostHandler(t)

	r := newChiRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/posts/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] != "not found" {
		t.Errorf("error = %q, want %q", resp["error"], "not found")
	}
}

func TestPostHandler_ListTags(t *testing.T) {
	h, svc := setupPostHandler(t)
	svc.Create(model.CreatePostInput{Title: "P1", Content: "c", Tags: []string{"python", "ml"}, Published: true})
	svc.Create(model.CreatePostInput{Title: "P2", Content: "c", Tags: []string{"python", "tutorial"}, Published: true})

	r := newChiRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/tags", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var tags []model.TagCount
	json.NewDecoder(w.Body).Decode(&tags)

	tagMap := make(map[string]int)
	for _, tc := range tags {
		tagMap[tc.Name] = tc.Count
	}
	if tagMap["python"] != 2 {
		t.Errorf("python count = %d, want 2", tagMap["python"])
	}
	if tagMap["ml"] != 1 {
		t.Errorf("ml count = %d, want 1", tagMap["ml"])
	}
	if tagMap["tutorial"] != 1 {
		t.Errorf("tutorial count = %d, want 1", tagMap["tutorial"])
	}
}
