package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/service"
)

// MusicHandler handles public HTTP requests for music recommendations.
type MusicHandler struct {
	service *service.MusicService
}

// NewMusicHandler creates a new MusicHandler.
func NewMusicHandler(svc *service.MusicService) *MusicHandler {
	return &MusicHandler{service: svc}
}

// ListPublic handles GET /api/v1/music — returns published music recs.
func (h *MusicHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	recs, err := h.service.ListPublished()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if recs == nil {
		recs = []model.MusicRec{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"recs": recs})
}

// AdminMusicHandler handles admin HTTP requests for music recommendations.
type AdminMusicHandler struct {
	service *service.MusicService
}

// NewAdminMusicHandler creates a new AdminMusicHandler.
func NewAdminMusicHandler(svc *service.MusicService) *AdminMusicHandler {
	return &AdminMusicHandler{service: svc}
}

// ListAll handles GET /api/v1/admin/music.
func (h *AdminMusicHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	recs, err := h.service.ListAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if recs == nil {
		recs = []model.MusicRec{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"recs": recs})
}

// Create handles POST /api/v1/admin/music.
func (h *AdminMusicHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.CreateMusicRecInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	rec, err := h.service.Create(input)
	if err == model.ErrMusicRecValidation {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, rec)
}

// Update handles PUT /api/v1/admin/music/{id}.
func (h *AdminMusicHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var input model.UpdateMusicRecInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	rec, err := h.service.Update(id, input)
	if err == model.ErrNotFound {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err == model.ErrMusicRecValidation {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

// Delete handles DELETE /api/v1/admin/music/{id}.
func (h *AdminMusicHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.service.Delete(id); err == model.ErrNotFound {
		writeError(w, http.StatusNotFound, "not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
