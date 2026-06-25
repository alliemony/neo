package repository

import (
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
)

func newTestMusicRec() *model.MusicRec {
	return &model.MusicRec{
		Title:     "Kind of Blue",
		Artist:    "Miles Davis",
		Published: true,
	}
}

func TestMusicRepo_Create(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewMusicRepo(db)
	rec := newTestMusicRec()
	if err := repo.Create(rec); err != nil {
		t.Fatalf("create: %v", err)
	}
	if rec.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
	if rec.CreatedAt.IsZero() {
		t.Fatal("expected non-zero created_at")
	}
}

func TestMusicRepo_GetByID(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewMusicRepo(db)
	rec := newTestMusicRec()
	repo.Create(rec)

	got, err := repo.GetByID(rec.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Title != rec.Title {
		t.Fatalf("expected %q, got %q", rec.Title, got.Title)
	}
}

func TestMusicRepo_GetByID_NotFound(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewMusicRepo(db)
	_, err = repo.GetByID(9999)
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMusicRepo_ListPublished_FiltersUnpublished(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewMusicRepo(db)
	pub := newTestMusicRec()
	pub.Published = true
	repo.Create(pub)

	draft := newTestMusicRec()
	draft.Title = "Draft Album"
	draft.Published = false
	repo.Create(draft)

	recs, err := repo.ListPublished()
	if err != nil {
		t.Fatalf("list published: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 published, got %d", len(recs))
	}
	if recs[0].Title != pub.Title {
		t.Fatalf("expected %q, got %q", pub.Title, recs[0].Title)
	}
}

func TestMusicRepo_ListAll_IncludesDrafts(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewMusicRepo(db)
	r1 := newTestMusicRec()
	r1.Published = true
	repo.Create(r1)

	r2 := newTestMusicRec()
	r2.Title = "Draft"
	r2.Published = false
	repo.Create(r2)

	all, err := repo.ListAll()
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2, got %d", len(all))
	}
}

func TestMusicRepo_Update(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewMusicRepo(db)
	rec := newTestMusicRec()
	repo.Create(rec)

	rec.Title = "Bitches Brew"
	rec.Published = false
	if err := repo.Update(rec); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, _ := repo.GetByID(rec.ID)
	if got.Title != "Bitches Brew" {
		t.Fatalf("expected updated title, got %q", got.Title)
	}
	if got.Published {
		t.Fatal("expected published=false after update")
	}
}

func TestMusicRepo_Delete(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewMusicRepo(db)
	rec := newTestMusicRec()
	repo.Create(rec)

	if err := repo.Delete(rec.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err = repo.GetByID(rec.ID)
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestMusicRepo_Delete_NotFound(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	repo := NewMusicRepo(db)
	err = repo.Delete(9999)
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
