package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alliemony/neo/backend/internal/config"
	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
	"github.com/alliemony/neo/backend/internal/service"
)

func TestNewRouter_PostsEndpointReturnsJSON(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("database.New: %v", err)
	}
	defer db.Close()

	postService := service.NewPostService(repository.NewPostRepo(db))
	created, err := postService.Create(model.CreatePostInput{
		Title:       "Integration Post",
		Content:     "Server-level test content",
		ContentType: "markdown",
		Tags:        []string{"integration"},
		Published:   true,
	})
	if err != nil {
		t.Fatalf("create test post: %v", err)
	}

	router := newRouter(config.Config{
		CORSOrigins:   "http://localhost:5173",
		AuthMode:      "basic",
		AdminUsername: "admin",
		AdminPassword: "changeme",
		BaseURL:       "http://localhost:8080",
		SessionSecret: "test-session-secret",
	}, db)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/posts")
	if err != nil {
		t.Fatalf("GET /api/v1/posts: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type = %q, want %q", got, "application/json")
	}

	var payload struct {
		Posts []model.Post `json:"posts"`
		Total int          `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Total < 1 {
		t.Fatalf("total = %d, want at least 1", payload.Total)
	}

	found := false
	for _, post := range payload.Posts {
		if post.Slug == created.Slug {
			found = true
			if post.Title != created.Title {
				t.Fatalf("title = %q, want %q", post.Title, created.Title)
			}
		}
	}

	if !found {
		t.Fatalf("response posts did not include %q", created.Slug)
	}
}