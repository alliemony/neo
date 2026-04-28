package service

import (
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

func setupPageService(t *testing.T) *PageService {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewPageService(repository.NewPageRepo(db))
}

func TestPageService_Create(t *testing.T) {
	svc := setupPageService(t)
	page, err := svc.Create(model.CreatePageInput{
		Title:   "About",
		Content: "Hello",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if page.Slug != "about" {
		t.Fatalf("expected slug 'about', got %q", page.Slug)
	}
	if page.ContentType != "markdown" {
		t.Fatalf("expected content_type 'markdown', got %q", page.ContentType)
	}
}

func TestPageService_Create_TitleRequired(t *testing.T) {
	svc := setupPageService(t)
	_, err := svc.Create(model.CreatePageInput{Content: "no title"})
	if err != model.ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestPageService_Create_DuplicateSlug(t *testing.T) {
	svc := setupPageService(t)
	svc.Create(model.CreatePageInput{Title: "About", Content: "a"})
	_, err := svc.Create(model.CreatePageInput{Title: "About", Content: "b"})
	if err != model.ErrSlugExists {
		t.Fatalf("expected ErrSlugExists, got %v", err)
	}
}

func TestPageService_GetBySlug(t *testing.T) {
	svc := setupPageService(t)
	svc.Create(model.CreatePageInput{Title: "Contact", Content: "hi"})

	page, err := svc.GetBySlug("contact")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if page.Title != "Contact" {
		t.Fatalf("expected 'Contact', got %q", page.Title)
	}
}

func TestPageService_Update(t *testing.T) {
	svc := setupPageService(t)
	svc.Create(model.CreatePageInput{Title: "About", Content: "old"})

	newTitle := "About Us"
	newContent := "new content"
	updated, err := svc.Update("about", model.UpdatePageInput{
		Title:   &newTitle,
		Content: &newContent,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Title != "About Us" {
		t.Fatalf("expected 'About Us', got %q", updated.Title)
	}
}

func TestPageService_Update_NotFound(t *testing.T) {
	svc := setupPageService(t)
	title := "nope"
	_, err := svc.Update("nope", model.UpdatePageInput{Title: &title})
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPageService_Delete(t *testing.T) {
	svc := setupPageService(t)
	svc.Create(model.CreatePageInput{Title: "About", Content: "a"})

	if err := svc.Delete("about"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err := svc.GetBySlug("about")
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestPageService_ListPublished(t *testing.T) {
	svc := setupPageService(t)
	svc.Create(model.CreatePageInput{Title: "About", Content: "a", Published: true})
	svc.Create(model.CreatePageInput{Title: "Draft", Content: "b", Published: false})

	pages, err := svc.ListPublished()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("expected 1, got %d", len(pages))
	}
}

func TestPageService_ListAll(t *testing.T) {
	svc := setupPageService(t)
	svc.Create(model.CreatePageInput{Title: "About", Content: "a", Published: true})
	svc.Create(model.CreatePageInput{Title: "Draft", Content: "b", Published: false})

	pages, err := svc.ListAll()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(pages) != 2 {
		t.Fatalf("expected 2, got %d", len(pages))
	}
}
