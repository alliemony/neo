package handler

import (
	"bytes"
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

func setupCommentHandler(t *testing.T) (*CommentHandler, *service.PostService) {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("setup db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	postRepo := repository.NewPostRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	postSvc := service.NewPostService(postRepo)
	commentSvc := service.NewCommentService(commentRepo, postRepo)
	h := NewCommentHandler(commentSvc, postRepo)
	return h, postSvc
}

func newCommentRouter(h *CommentHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/v1/posts/{slug}/comments", h.ListComments)
	r.Post("/api/v1/posts/{slug}/comments", h.CreateComment)
	r.Post("/api/v1/posts/{slug}/like", h.Like)
	return r
}

func TestCommentHandler_ListComments(t *testing.T) {
	h, postSvc := setupCommentHandler(t)
	postSvc.Create(model.CreatePostInput{Title: "Test Post", Content: "c", Published: true})

	r := newCommentRouter(h)

	// Create a comment via POST first.
	body, _ := json.Marshal(map[string]string{"author_name": "alice", "content": "Nice!"})
	req := httptest.NewRequest("POST", "/api/v1/posts/test-post/comments", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create comment status = %d, want %d", w.Code, http.StatusCreated)
	}

	// Now list comments.
	req = httptest.NewRequest("GET", "/api/v1/posts/test-post/comments", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("list status = %d, want %d", w.Code, http.StatusOK)
	}

	var comments []model.Comment
	json.NewDecoder(w.Body).Decode(&comments)
	if len(comments) != 1 {
		t.Errorf("len(comments) = %d, want 1", len(comments))
	}
	if comments[0].AuthorName != "alice" {
		t.Errorf("author = %q, want alice", comments[0].AuthorName)
	}
}

func TestCommentHandler_CreateComment_Valid(t *testing.T) {
	h, postSvc := setupCommentHandler(t)
	postSvc.Create(model.CreatePostInput{Title: "Test Post", Content: "c", Published: true})

	r := newCommentRouter(h)
	body, _ := json.Marshal(map[string]string{"author_name": "bob", "content": "Hello!"})
	req := httptest.NewRequest("POST", "/api/v1/posts/test-post/comments", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var comment model.Comment
	json.NewDecoder(w.Body).Decode(&comment)
	if comment.AuthorName != "bob" {
		t.Errorf("author = %q, want bob", comment.AuthorName)
	}
}

func TestCommentHandler_CreateComment_PostNotFound(t *testing.T) {
	h, _ := setupCommentHandler(t)

	r := newCommentRouter(h)
	body, _ := json.Marshal(map[string]string{"author_name": "bob", "content": "Hello!"})
	req := httptest.NewRequest("POST", "/api/v1/posts/nonexistent/comments", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestCommentHandler_CreateComment_ValidationError(t *testing.T) {
	h, postSvc := setupCommentHandler(t)
	postSvc.Create(model.CreatePostInput{Title: "Test Post", Content: "c", Published: true})

	r := newCommentRouter(h)
	body, _ := json.Marshal(map[string]string{"author_name": "", "content": "Hello!"})
	req := httptest.NewRequest("POST", "/api/v1/posts/test-post/comments", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCommentHandler_Like(t *testing.T) {
	h, postSvc := setupCommentHandler(t)
	postSvc.Create(model.CreatePostInput{Title: "Likeable Post", Content: "c", Published: true})

	r := newCommentRouter(h)
	req := httptest.NewRequest("POST", "/api/v1/posts/likeable-post/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]int
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["like_count"] != 1 {
		t.Errorf("like_count = %d, want 1", resp["like_count"])
	}

	// Second like.
	req = httptest.NewRequest("POST", "/api/v1/posts/likeable-post/like", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.NewDecoder(w.Body).Decode(&resp)
	if resp["like_count"] != 2 {
		t.Errorf("like_count = %d, want 2", resp["like_count"])
	}
}

func TestCommentHandler_Like_PostNotFound(t *testing.T) {
	h, _ := setupCommentHandler(t)

	r := newCommentRouter(h)
	req := httptest.NewRequest("POST", "/api/v1/posts/nonexistent/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
