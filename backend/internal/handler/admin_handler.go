package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/service"
)

type AdminHandler struct {
	postService *service.PostService
	pageService *service.PageService
}

func NewAdminHandler(postSvc *service.PostService, pageSvc *service.PageService) *AdminHandler {
	return &AdminHandler{postService: postSvc, pageService: pageSvc}
}

func (h *AdminHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var input model.CreatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.postService.Create(input)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrTitleRequired):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, model.ErrSlugExists):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, post)
}

func (h *AdminHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	posts, total, err := h.postService.ListAll(model.ListOptions{Page: page, PerPage: perPage})
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

func (h *AdminHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var input model.UpdatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.postService.Update(slug, input)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func (h *AdminHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	err := h.postService.Delete(slug)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
	var input model.CreatePageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	page, err := h.pageService.Create(input)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrTitleRequired):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, model.ErrSlugExists):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, page)
}

func (h *AdminHandler) ListPages(w http.ResponseWriter, r *http.Request) {
	pages, err := h.pageService.ListAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if pages == nil {
		pages = []model.Page{}
	}

	writeJSON(w, http.StatusOK, pages)
}

func (h *AdminHandler) UpdatePage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var input model.UpdatePageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	page, err := h.pageService.Update(slug, input)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, page)
}

func (h *AdminHandler) DeletePage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	err := h.pageService.Delete(slug)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
