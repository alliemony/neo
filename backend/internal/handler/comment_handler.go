package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
	"github.com/alliemony/neo/backend/internal/service"
)

// CommentHandler handles HTTP requests for comments and post likes.
type CommentHandler struct {
	commentService *service.CommentService
	postRepo       *repository.PostRepo
}

// NewCommentHandler creates a new CommentHandler.
func NewCommentHandler(svc *service.CommentService, postRepo *repository.PostRepo) *CommentHandler {
	return &CommentHandler{commentService: svc, postRepo: postRepo}
}

// ListComments handles GET /api/v1/posts/{slug}/comments.
func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	comments, err := h.commentService.ListByPostSlug(slug)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if comments == nil {
		comments = []model.Comment{}
	}

	writeJSON(w, http.StatusOK, comments)
}

// List is an alias for ListComments for backward compatibility.
func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	h.ListComments(w, r)
}

// CreateComment handles POST /api/v1/posts/{slug}/comments.
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var input model.CreateCommentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	comment, err := h.commentService.Create(slug, input)
	if err != nil {
		switch err {
		case model.ErrAuthorRequired, model.ErrContentRequired, model.ErrContentTooLong:
			writeError(w, http.StatusBadRequest, err.Error())
		case model.ErrNotFound:
			writeError(w, http.StatusNotFound, "post not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

// Create is an alias for CreateComment for backward compatibility.
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.CreateComment(w, r)
}

// Like handles POST /api/v1/posts/{slug}/like.
func (h *CommentHandler) Like(w http.ResponseWriter, r *http.Request) {
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

