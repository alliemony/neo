package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
	"github.com/alliemony/neo/backend/internal/service"
)

func setupPageRouter(t *testing.T) (*service.PageService, http.Handler) {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	pageRepo := repository.NewPageRepo(db)
	pageSvc := service.NewPageService(pageRepo)
	pageHandler := NewPageHandler(pageSvc)

	r := chi.NewRouter()
	r.Get("/api/v1/pages", pageHandler.ListPublished)
	r.Get("/api/v1/pages/{slug}", pageHandler.GetBySlug)
	return pageSvc, r
}

func TestPageHandler_ListPublished(t *testing.T) {
	pageSvc, router := setupPageRouter(t)

	pageSvc.Create(model.CreatePageInput{Title: "About", Content: "a", Published: true, SortOrder: 1})
	pageSvc.Create(model.CreatePageInput{Title: "Projects", Content: "b", Published: true, SortOrder: 2})
	pageSvc.Create(model.CreatePageInput{Title: "Draft", Content: "c", Published: false})

	req := httptest.NewRequest("GET", "/api/v1/pages", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var pages []model.Page
	json.NewDecoder(rr.Body).Decode(&pages)
	if len(pages) != 2 {
		t.Fatalf("expected 2 published, got %d", len(pages))
	}
	if pages[0].Title != "About" || pages[1].Title != "Projects" {
		t.Fatalf("unexpected order: %q, %q", pages[0].Title, pages[1].Title)
	}
}

func TestPageHandler_GetBySlug(t *testing.T) {
	pageSvc, router := setupPageRouter(t)
	pageSvc.Create(model.CreatePageInput{Title: "About", Content: "Hello world", Published: true})

	req := httptest.NewRequest("GET", "/api/v1/pages/about", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var page model.Page
	json.NewDecoder(rr.Body).Decode(&page)
	if page.Title != "About" {
		t.Fatalf("expected 'About', got %q", page.Title)
	}
}

func TestPageHandler_GetBySlug_NotFound(t *testing.T) {
	_, router := setupPageRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/pages/nonexistent", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestPageHandler_GetBySlug_DraftReturns404(t *testing.T) {
	pageSvc, router := setupPageRouter(t)
	pageSvc.Create(model.CreatePageInput{Title: "Secret", Content: "hidden", Published: false})

	req := httptest.NewRequest("GET", "/api/v1/pages/secret", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for draft page, got %d", rr.Code)
	}
}
