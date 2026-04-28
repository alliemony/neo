package repository

import (
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
)

func TestPageRepo_Create(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewPageRepo(db)
	page := &model.Page{
		Slug:        "about",
		Title:       "About Me",
		Content:     "Hello world",
		ContentType: "markdown",
		Published:   true,
		SortOrder:   1,
	}

	if err := repo.Create(page); err != nil {
		t.Fatalf("create page: %v", err)
	}

	if page.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
	if page.CreatedAt.IsZero() {
		t.Fatal("expected non-zero created_at")
	}
}

func TestPageRepo_GetBySlug(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewPageRepo(db)
	page := &model.Page{Slug: "about", Title: "About", Content: "Hi", ContentType: "markdown", Published: true}
	repo.Create(page)

	got, err := repo.GetBySlug("about")
	if err != nil {
		t.Fatalf("get by slug: %v", err)
	}
	if got.Title != "About" {
		t.Fatalf("expected title 'About', got %q", got.Title)
	}
}

func TestPageRepo_GetBySlug_NotFound(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewPageRepo(db)
	_, err = repo.GetBySlug("nope")
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPageRepo_List(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewPageRepo(db)
	repo.Create(&model.Page{Slug: "about", Title: "About", Content: "a", ContentType: "markdown", Published: true, SortOrder: 1})
	repo.Create(&model.Page{Slug: "contact", Title: "Contact", Content: "b", ContentType: "markdown", Published: true, SortOrder: 2})
	repo.Create(&model.Page{Slug: "draft", Title: "Draft", Content: "c", ContentType: "markdown", Published: false, SortOrder: 3})

	pages, err := repo.List(false)
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if len(pages) != 3 {
		t.Fatalf("expected 3, got %d", len(pages))
	}

	published, err := repo.List(true)
	if err != nil {
		t.Fatalf("list published: %v", err)
	}
	if len(published) != 2 {
		t.Fatalf("expected 2 published, got %d", len(published))
	}
}

func TestPageRepo_Update(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewPageRepo(db)
	page := &model.Page{Slug: "about", Title: "About", Content: "old", ContentType: "markdown"}
	repo.Create(page)

	page.Title = "About Us"
	page.Content = "new content"
	if err := repo.Update(page); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, _ := repo.GetBySlug("about")
	if got.Title != "About Us" {
		t.Fatalf("expected 'About Us', got %q", got.Title)
	}
	if got.Content != "new content" {
		t.Fatalf("expected 'new content', got %q", got.Content)
	}
}

func TestPageRepo_Delete(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewPageRepo(db)
	repo.Create(&model.Page{Slug: "about", Title: "About", Content: "a", ContentType: "markdown"})

	if err := repo.Delete("about"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = repo.GetBySlug("about")
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestPageRepo_Delete_NotFound(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewPageRepo(db)
	if err := repo.Delete("nope"); err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
