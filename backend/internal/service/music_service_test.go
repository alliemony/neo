package service

import (
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

func newMusicService(t *testing.T) *MusicService {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewMusicService(repository.NewMusicRepo(db))
}

func validMusicInput() model.CreateMusicRecInput {
	return model.CreateMusicRecInput{
		Title:     "Kind of Blue",
		Artist:    "Miles Davis",
		Published: true,
	}
}

func TestMusicService_Create_Valid(t *testing.T) {
	svc := newMusicService(t)
	rec, err := svc.Create(validMusicInput())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
}

func TestMusicService_Create_MissingTitle(t *testing.T) {
	svc := newMusicService(t)
	input := validMusicInput()
	input.Title = ""
	_, err := svc.Create(input)
	if err != model.ErrMusicRecValidation {
		t.Fatalf("expected ErrMusicRecValidation, got %v", err)
	}
}

func TestMusicService_Create_MissingArtist(t *testing.T) {
	svc := newMusicService(t)
	input := validMusicInput()
	input.Artist = ""
	_, err := svc.Create(input)
	if err != model.ErrMusicRecValidation {
		t.Fatalf("expected ErrMusicRecValidation, got %v", err)
	}
}

func TestMusicService_Update_PartialFields(t *testing.T) {
	svc := newMusicService(t)
	rec, _ := svc.Create(validMusicInput())

	newTitle := "Bitches Brew"
	updated, err := svc.Update(rec.ID, model.UpdateMusicRecInput{Title: &newTitle})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Title != "Bitches Brew" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}
	if updated.Artist != rec.Artist {
		t.Fatal("artist should not have changed")
	}
}

func TestMusicService_Update_NotFound(t *testing.T) {
	svc := newMusicService(t)
	newTitle := "x"
	_, err := svc.Update(9999, model.UpdateMusicRecInput{Title: &newTitle})
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMusicService_Delete(t *testing.T) {
	svc := newMusicService(t)
	rec, _ := svc.Create(validMusicInput())
	if err := svc.Delete(rec.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	_, err := svc.GetByID(rec.ID)
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestMusicService_ListPublished(t *testing.T) {
	svc := newMusicService(t)
	pub := validMusicInput()
	pub.Published = true
	svc.Create(pub)

	draft := validMusicInput()
	draft.Title = "Draft"
	draft.Published = false
	svc.Create(draft)

	recs, err := svc.ListPublished()
	if err != nil {
		t.Fatalf("list published: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 published, got %d", len(recs))
	}
}
