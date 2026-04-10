package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/service"
)

// PostHandler handles HTTP requests for blog posts.
type PostHandler struct {
	service *service.PostService
}

// NewPostHandler creates a new PostHandler.
func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{service: svc}
}

// List handles GET /api/v1/posts with optional ?tag= and ?page= query params.
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	opts := model.ListOptions{Page: page, PerPage: perPage}

	var posts []model.Post
	var total int
	var err error

	if tag != "" {
		posts, total, err = h.service.ListByTag(tag, opts)
	} else {
		posts, total, err = h.service.ListPublished(opts)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if posts == nil {
		posts = []model.Post{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"posts": posts,
		"total": total,
	})
}

// GetBySlug handles GET /api/v1/posts/{slug}.
func (h *PostHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	post, err := h.service.GetBySlug(slug)
	if err == model.ErrNotFound {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, post)
}

// ListTags handles GET /api/v1/tags.
func (h *PostHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.service.AllTags()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if tags == nil {
		tags = []model.TagCount{}
	}

	writeJSON(w, http.StatusOK, tags)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
