package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
	"github.com/alliemony/neo/backend/internal/service"
)

func setupRecsRouter(t *testing.T) (*service.RecsService, http.Handler) {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	svc := service.NewRecsService(repository.NewRecsRepo(db))
	pub := NewRecsHandler(svc)
	adm := NewAdminRecsHandler(svc)

	r := chi.NewRouter()
	r.Get("/api/v1/recs", pub.ListPublic)
	r.Get("/api/v1/admin/recs", adm.ListAll)
	r.Post("/api/v1/admin/recs", adm.Create)
	r.Put("/api/v1/admin/recs/{id}", adm.Update)
	r.Delete("/api/v1/admin/recs/{id}", adm.Delete)

	return svc, r
}

func validRecBody() map[string]interface{} {
	return map[string]interface{}{
		"name":        "Test Tool",
		"section":     "Tools",
		"tag":         "CLI",
		"tag_bg":      "#edf5f3",
		"tag_color":   "#1a7060",
		"description": "A great tool.",
		"published":   true,
	}
}

func TestRecsHandler_ListPublic_GroupsBySection(t *testing.T) {
	svc, router := setupRecsRouter(t)

	svc.Create(model.CreateRecInput{Name: "Tool A", Section: "Tools", Tag: "CLI", TagBg: "#fff", TagColor: "#000", Description: "d", Published: true})
	svc.Create(model.CreateRecInput{Name: "Book A", Section: "Reading", Tag: "Book", TagBg: "#fff", TagColor: "#000", Description: "d", Published: true})
	svc.Create(model.CreateRecInput{Name: "Draft", Section: "Tools", Tag: "CLI", TagBg: "#fff", TagColor: "#000", Description: "d", Published: false})

	req := httptest.NewRequest("GET", "/api/v1/recs", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var body struct {
		Sections []model.RecSection `json:"sections"`
	}
	json.NewDecoder(rr.Body).Decode(&body)
	if len(body.Sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(body.Sections))
	}
	if body.Sections[0].Title != "Tools" || len(body.Sections[0].Items) != 1 {
		t.Fatalf("unexpected Tools section: %+v", body.Sections[0])
	}
}

func TestRecsHandler_ListPublic_EmptySections(t *testing.T) {
	_, router := setupRecsRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/recs", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var body struct {
		Sections []model.RecSection `json:"sections"`
	}
	json.NewDecoder(rr.Body).Decode(&body)
	if body.Sections == nil {
		t.Fatal("expected non-nil sections slice")
	}
}

func TestAdminRecsHandler_Create_Valid(t *testing.T) {
	_, router := setupRecsRouter(t)

	body, _ := json.Marshal(validRecBody())
	req := httptest.NewRequest("POST", "/api/v1/admin/recs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
	var rec model.Rec
	json.NewDecoder(rr.Body).Decode(&rec)
	if rec.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
}

func TestAdminRecsHandler_Create_MissingField(t *testing.T) {
	_, router := setupRecsRouter(t)

	b := validRecBody()
	delete(b, "name")
	body, _ := json.Marshal(b)
	req := httptest.NewRequest("POST", "/api/v1/admin/recs", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestAdminRecsHandler_Update_NotFound(t *testing.T) {
	_, router := setupRecsRouter(t)

	body, _ := json.Marshal(map[string]string{"name": "New"})
	req := httptest.NewRequest("PUT", "/api/v1/admin/recs/9999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestAdminRecsHandler_Delete_Valid(t *testing.T) {
	svc, router := setupRecsRouter(t)
	rec, _ := svc.Create(model.CreateRecInput{Name: "x", Section: "S", Tag: "T", TagBg: "#fff", TagColor: "#000", Description: "d"})

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/admin/recs/%d", rec.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestAdminRecsHandler_Delete_NotFound(t *testing.T) {
	_, router := setupRecsRouter(t)

	req := httptest.NewRequest("DELETE", "/api/v1/admin/recs/9999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
