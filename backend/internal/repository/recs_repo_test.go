package repository

import (
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
)

func newTestRec() *model.Rec {
	href := "https://example.com"
	return &model.Rec{
		Name:        "Test Tool",
		Href:        &href,
		Section:     "Tools",
		Tag:         "CLI",
		TagBg:       "#edf5f3",
		TagColor:    "#1a7060",
		Description: "A great tool.",
		Published:   true,
		SortOrder:   0,
	}
}

func TestRecsRepo_Create(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewRecsRepo(db)
	rec := newTestRec()

	if err := repo.Create(rec); err != nil {
		t.Fatalf("create rec: %v", err)
	}
	if rec.ID == 0 {
		t.Fatal("expected non-zero ID after create")
	}
	if rec.CreatedAt.IsZero() {
		t.Fatal("expected non-zero created_at")
	}
}

func TestRecsRepo_GetByID(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewRecsRepo(db)
	rec := newTestRec()
	repo.Create(rec)

	got, err := repo.GetByID(rec.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Name != rec.Name {
		t.Fatalf("expected name %q, got %q", rec.Name, got.Name)
	}
}

func TestRecsRepo_GetByID_NotFound(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewRecsRepo(db)
	_, err = repo.GetByID(9999)
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRecsRepo_ListPublished_FiltersUnpublished(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewRecsRepo(db)

	published := newTestRec()
	published.Published = true
	repo.Create(published)

	draft := newTestRec()
	draft.Name = "Draft Rec"
	draft.Published = false
	repo.Create(draft)

	recs, err := repo.ListPublished()
	if err != nil {
		t.Fatalf("list published: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 published rec, got %d", len(recs))
	}
	if recs[0].Name != published.Name {
		t.Fatalf("expected %q, got %q", published.Name, recs[0].Name)
	}
}

func TestRecsRepo_ListAll_IncludesDrafts(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewRecsRepo(db)

	r1 := newTestRec()
	r1.Published = true
	repo.Create(r1)

	r2 := newTestRec()
	r2.Name = "Draft"
	r2.Published = false
	repo.Create(r2)

	all, err := repo.ListAll()
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 recs, got %d", len(all))
	}
}

func TestRecsRepo_Update(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewRecsRepo(db)
	rec := newTestRec()
	repo.Create(rec)

	rec.Name = "Updated Name"
	rec.Published = false
	if err := repo.Update(rec); err != nil {
		t.Fatalf("update rec: %v", err)
	}

	got, _ := repo.GetByID(rec.ID)
	if got.Name != "Updated Name" {
		t.Fatalf("expected updated name, got %q", got.Name)
	}
	if got.Published {
		t.Fatal("expected published=false after update")
	}
}

func TestRecsRepo_Delete(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewRecsRepo(db)
	rec := newTestRec()
	repo.Create(rec)

	if err := repo.Delete(rec.ID); err != nil {
		t.Fatalf("delete rec: %v", err)
	}
	_, err = repo.GetByID(rec.ID)
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestRecsRepo_Delete_NotFound(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewRecsRepo(db)
	err = repo.Delete(9999)
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
