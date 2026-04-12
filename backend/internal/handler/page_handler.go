package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/service"
)

// PageHandler handles public HTTP requests for static pages.
type PageHandler struct {
	service *service.PageService
}

// NewPageHandler creates a new PageHandler.
func NewPageHandler(svc *service.PageService) *PageHandler {
	return &PageHandler{service: svc}
}

// ListPublished handles GET /api/v1/pages.
func (h *PageHandler) ListPublished(w http.ResponseWriter, r *http.Request) {
	pages, err := h.service.ListPublished()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if pages == nil {
		pages = []model.Page{}
	}

	writeJSON(w, http.StatusOK, pages)
}

// GetBySlug handles GET /api/v1/pages/{slug}.
func (h *PageHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	page, err := h.service.GetBySlug(slug)
	if err == model.ErrNotFound {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if !page.Published {
		writeError(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, page)
}

// AdminPageHandler handles admin HTTP requests for page management.
type AdminPageHandler struct {
	service *service.PageService
}

// NewAdminPageHandler creates a new AdminPageHandler.
func NewAdminPageHandler(svc *service.PageService) *AdminPageHandler {
	return &AdminPageHandler{service: svc}
}

// Create handles POST /api/v1/admin/pages.
func (h *AdminPageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.CreatePageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	page, err := h.service.Create(input)
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

	writeJSON(w, http.StatusCreated, page)
}

// Update handles PUT /api/v1/admin/pages/{slug}.
func (h *AdminPageHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var input model.UpdatePageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	page, err := h.service.Update(slug, input)
	if err == model.ErrNotFound {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, page)
}

// Delete handles DELETE /api/v1/admin/pages/{slug}.
func (h *AdminPageHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

// ListAll handles GET /api/v1/admin/pages (includes unpublished).
func (h *AdminPageHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	pages, err := h.service.ListAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if pages == nil {
		pages = []model.Page{}
	}

	writeJSON(w, http.StatusOK, pages)
}
