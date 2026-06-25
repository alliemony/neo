package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/service"
)

// RecsHandler handles public HTTP requests for recommendations.
type RecsHandler struct {
	service *service.RecsService
}

// NewRecsHandler creates a new RecsHandler.
func NewRecsHandler(svc *service.RecsService) *RecsHandler {
	return &RecsHandler{service: svc}
}

// ListPublic handles GET /api/v1/recs — returns published recs grouped by section.
func (h *RecsHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	recs, err := h.service.ListPublished()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	sections := groupBySection(recs)
	writeJSON(w, http.StatusOK, map[string]interface{}{"sections": sections})
}

// groupBySection groups a flat slice of recs into RecSection slices,
// preserving the first-appearance order of section names.
func groupBySection(recs []model.Rec) []model.RecSection {
	var order []string
	groups := make(map[string][]model.Rec)
	for _, rec := range recs {
		if _, seen := groups[rec.Section]; !seen {
			order = append(order, rec.Section)
		}
		groups[rec.Section] = append(groups[rec.Section], rec)
	}

	sections := make([]model.RecSection, 0, len(order))
	for _, title := range order {
		sections = append(sections, model.RecSection{Title: title, Items: groups[title]})
	}
	return sections
}

// AdminRecsHandler handles admin HTTP requests for recommendations.
type AdminRecsHandler struct {
	service *service.RecsService
}

// NewAdminRecsHandler creates a new AdminRecsHandler.
func NewAdminRecsHandler(svc *service.RecsService) *AdminRecsHandler {
	return &AdminRecsHandler{service: svc}
}

// ListAll handles GET /api/v1/admin/recs.
func (h *AdminRecsHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	recs, err := h.service.ListAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if recs == nil {
		recs = []model.Rec{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"recs": recs})
}

// Create handles POST /api/v1/admin/recs.
func (h *AdminRecsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.CreateRecInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	rec, err := h.service.Create(input)
	if err == model.ErrRecValidation {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, rec)
}

// Update handles PUT /api/v1/admin/recs/{id}.
func (h *AdminRecsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input model.UpdateRecInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	rec, err := h.service.Update(id, input)
	if err == model.ErrNotFound {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err == model.ErrRecValidation {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, rec)
}

// Delete handles DELETE /api/v1/admin/recs/{id}.
func (h *AdminRecsHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
