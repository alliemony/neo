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

func setupMusicRouter(t *testing.T) (*service.MusicService, http.Handler) {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	svc := service.NewMusicService(repository.NewMusicRepo(db))
	pub := NewMusicHandler(svc)
	adm := NewAdminMusicHandler(svc)

	r := chi.NewRouter()
	r.Get("/api/v1/music", pub.ListPublic)
	r.Get("/api/v1/admin/music", adm.ListAll)
	r.Post("/api/v1/admin/music", adm.Create)
	r.Put("/api/v1/admin/music/{id}", adm.Update)
	r.Delete("/api/v1/admin/music/{id}", adm.Delete)

	return svc, r
}

func validMusicBody() map[string]interface{} {
	return map[string]interface{}{
		"title":     "Kind of Blue",
		"artist":    "Miles Davis",
		"published": true,
	}
}

func TestMusicHandler_ListPublic_ReturnsPublishedOnly(t *testing.T) {
	svc, router := setupMusicRouter(t)
	svc.Create(model.CreateMusicRecInput{Title: "Kind of Blue", Artist: "Miles Davis", Published: true})
	svc.Create(model.CreateMusicRecInput{Title: "Draft Album", Artist: "Unknown", Published: false})

	req := httptest.NewRequest("GET", "/api/v1/music", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var body struct {
		Recs []model.MusicRec `json:"recs"`
	}
	json.NewDecoder(rr.Body).Decode(&body)
	if len(body.Recs) != 1 {
		t.Fatalf("expected 1 published rec, got %d", len(body.Recs))
	}
}

func TestMusicHandler_ListPublic_EmptyList(t *testing.T) {
	_, router := setupMusicRouter(t)
	req := httptest.NewRequest("GET", "/api/v1/music", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var body struct {
		Recs []model.MusicRec `json:"recs"`
	}
	json.NewDecoder(rr.Body).Decode(&body)
	if body.Recs == nil {
		t.Fatal("expected non-nil recs slice")
	}
}

func TestAdminMusicHandler_Create_Valid(t *testing.T) {
	_, router := setupMusicRouter(t)
	body, _ := json.Marshal(validMusicBody())
	req := httptest.NewRequest("POST", "/api/v1/admin/music", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
	var rec model.MusicRec
	json.NewDecoder(rr.Body).Decode(&rec)
	if rec.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
}

func TestAdminMusicHandler_Create_MissingTitle(t *testing.T) {
	_, router := setupMusicRouter(t)
	b := validMusicBody()
	delete(b, "title")
	body, _ := json.Marshal(b)
	req := httptest.NewRequest("POST", "/api/v1/admin/music", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestAdminMusicHandler_Update_NotFound(t *testing.T) {
	_, router := setupMusicRouter(t)
	body, _ := json.Marshal(map[string]string{"title": "New"})
	req := httptest.NewRequest("PUT", "/api/v1/admin/music/9999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestAdminMusicHandler_Delete_Valid(t *testing.T) {
	svc, router := setupMusicRouter(t)
	rec, _ := svc.Create(model.CreateMusicRecInput{Title: "x", Artist: "y"})

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/admin/music/%d", rec.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestAdminMusicHandler_Delete_NotFound(t *testing.T) {
	_, router := setupMusicRouter(t)
	req := httptest.NewRequest("DELETE", "/api/v1/admin/music/9999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
