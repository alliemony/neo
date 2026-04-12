package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/service"
)

// CommentHandler handles HTTP requests for comments.
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler creates a new CommentHandler.
func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: svc}
}

// List handles GET /api/v1/posts/{slug}/comments.
func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
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

// Create handles POST /api/v1/posts/{slug}/comments.
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
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
