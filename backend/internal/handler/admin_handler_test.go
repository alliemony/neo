package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/middleware"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
	"github.com/alliemony/neo/backend/internal/service"
)

func setupAdminRouter(t *testing.T) http.Handler {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	postRepo := repository.NewPostRepo(db)
	pageRepo := repository.NewPageRepo(db)
	postSvc := service.NewPostService(postRepo)
	pageSvc := service.NewPageService(pageRepo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	auth := middleware.NewBasicAuth("admin", string(hash))

	adminHandler := NewAdminHandler(postSvc, pageSvc)

	r := chi.NewRouter()
	r.Route("/api/v1/admin", func(r chi.Router) {
		r.Use(auth.Middleware())
		r.Get("/posts", adminHandler.ListPosts)
		r.Post("/posts", adminHandler.CreatePost)
		r.Put("/posts/{slug}", adminHandler.UpdatePost)
		r.Delete("/posts/{slug}", adminHandler.DeletePost)
		r.Get("/pages", adminHandler.ListPages)
		r.Post("/pages", adminHandler.CreatePage)
		r.Put("/pages/{slug}", adminHandler.UpdatePage)
		r.Delete("/pages/{slug}", adminHandler.DeletePage)
	})
	return r
}

func adminRequest(method, path string, body any) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "secret")
	return req
}

func TestAdmin_Unauthorized(t *testing.T) {
	router := setupAdminRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/admin/posts", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAdmin_CreatePost(t *testing.T) {
	router := setupAdminRouter(t)

	req := adminRequest("POST", "/api/v1/admin/posts", model.CreatePostInput{
		Title:   "Test Post",
		Content: "Hello world",
		Tags:    []string{"go"},
	})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	var post model.Post
	json.NewDecoder(rr.Body).Decode(&post)
	if post.Slug != "test-post" {
		t.Fatalf("expected slug 'test-post', got %q", post.Slug)
	}
}

func TestAdmin_ListPosts(t *testing.T) {
	router := setupAdminRouter(t)

	// Create a draft post
	req := adminRequest("POST", "/api/v1/admin/posts", model.CreatePostInput{
		Title:     "Draft",
		Content:   "hidden",
		Published: false,
	})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Create a published post
	req = adminRequest("POST", "/api/v1/admin/posts", model.CreatePostInput{
		Title:     "Public",
		Content:   "visible",
		Published: true,
	})
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Admin list should include both
	req = adminRequest("GET", "/api/v1/admin/posts", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var result struct {
		Posts []model.Post `json:"posts"`
		Total int          `json:"total"`
	}
	json.NewDecoder(rr.Body).Decode(&result)
	if result.Total != 2 {
		t.Fatalf("expected 2 posts, got %d", result.Total)
	}
}

func TestAdmin_UpdatePost(t *testing.T) {
	router := setupAdminRouter(t)

	req := adminRequest("POST", "/api/v1/admin/posts", model.CreatePostInput{
		Title:   "Original",
		Content: "old content",
	})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	newTitle := "Updated"
	req = adminRequest("PUT", "/api/v1/admin/posts/original", model.UpdatePostInput{
		Title: &newTitle,
	})
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var post model.Post
	json.NewDecoder(rr.Body).Decode(&post)
	if post.Title != "Updated" {
		t.Fatalf("expected title 'Updated', got %q", post.Title)
	}
}

func TestAdmin_UpdatePost_NotFound(t *testing.T) {
	router := setupAdminRouter(t)

	title := "nope"
	req := adminRequest("PUT", "/api/v1/admin/posts/nonexistent", model.UpdatePostInput{Title: &title})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestAdmin_DeletePost(t *testing.T) {
	router := setupAdminRouter(t)

	req := adminRequest("POST", "/api/v1/admin/posts", model.CreatePostInput{
		Title:   "To Delete",
		Content: "bye",
	})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	req = adminRequest("DELETE", "/api/v1/admin/posts/to-delete", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestAdmin_DeletePost_NotFound(t *testing.T) {
	router := setupAdminRouter(t)

	req := adminRequest("DELETE", "/api/v1/admin/posts/nonexistent", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestAdmin_CreatePage(t *testing.T) {
	router := setupAdminRouter(t)

	req := adminRequest("POST", "/api/v1/admin/pages", model.CreatePageInput{
		Title:   "About",
		Content: "About us",
	})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	var page model.Page
	json.NewDecoder(rr.Body).Decode(&page)
	if page.Slug != "about" {
		t.Fatalf("expected slug 'about', got %q", page.Slug)
	}
}

func TestAdmin_ListPages(t *testing.T) {
	router := setupAdminRouter(t)

	req := adminRequest("POST", "/api/v1/admin/pages", model.CreatePageInput{Title: "About", Content: "a"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	req = adminRequest("GET", "/api/v1/admin/pages", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var pages []model.Page
	json.NewDecoder(rr.Body).Decode(&pages)
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
}

func TestAdmin_UpdatePage(t *testing.T) {
	router := setupAdminRouter(t)

	req := adminRequest("POST", "/api/v1/admin/pages", model.CreatePageInput{Title: "About", Content: "old"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	newContent := "new content"
	req = adminRequest("PUT", "/api/v1/admin/pages/about", model.UpdatePageInput{Content: &newContent})
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var page model.Page
	json.NewDecoder(rr.Body).Decode(&page)
	if page.Content != "new content" {
		t.Fatalf("expected 'new content', got %q", page.Content)
	}
}

func TestAdmin_DeletePage(t *testing.T) {
	router := setupAdminRouter(t)

	req := adminRequest("POST", "/api/v1/admin/pages", model.CreatePageInput{Title: "About", Content: "a"})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	req = adminRequest("DELETE", "/api/v1/admin/pages/about", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}
