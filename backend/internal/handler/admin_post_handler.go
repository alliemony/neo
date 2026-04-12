package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/service"
)

// AdminPostHandler handles admin HTTP requests for post management.
type AdminPostHandler struct {
	service *service.PostService
}

// NewAdminPostHandler creates a new AdminPostHandler.
func NewAdminPostHandler(svc *service.PostService) *AdminPostHandler {
	return &AdminPostHandler{service: svc}
}

// Create handles POST /api/v1/admin/posts.
func (h *AdminPostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.CreatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.service.Create(input)
	if err != nil {
		switch err {
		case model.ErrTitleRequired:
			writeError(w, http.StatusBadRequest, err.Error())
		case model.ErrSlugExists:
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, post)
}

// Update handles PUT /api/v1/admin/posts/{slug}.
func (h *AdminPostHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var input model.UpdatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.service.Update(slug, input)
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

// Delete handles DELETE /api/v1/admin/posts/{slug}.
func (h *AdminPostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	err := h.service.Delete(slug)
	if err == model.ErrNotFound {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListAll handles GET /api/v1/admin/posts (includes drafts).
func (h *AdminPostHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	posts, total, err := h.service.ListAll(model.ListOptions{Page: 1, PerPage: 100})
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
