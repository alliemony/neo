package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

// LikeHandler handles HTTP requests for post likes.
type LikeHandler struct {
	postRepo *repository.PostRepo
}

// NewLikeHandler creates a new LikeHandler.
func NewLikeHandler(repo *repository.PostRepo) *LikeHandler {
	return &LikeHandler{postRepo: repo}
}

// Like handles POST /api/v1/posts/{slug}/like.
func (h *LikeHandler) Like(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	count, err := h.postRepo.IncrementLikeCount(slug)
	if err == model.ErrNotFound {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"like_count": count})
}
