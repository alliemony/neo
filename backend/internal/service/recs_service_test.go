package service

import (
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

func newRecsService(t *testing.T) *RecsService {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewRecsService(repository.NewRecsRepo(db))
}

func validCreateInput() model.CreateRecInput {
	return model.CreateRecInput{
		Name:        "Test Tool",
		Section:     "Tools",
		Tag:         "CLI",
		TagBg:       "#edf5f3",
		TagColor:    "#1a7060",
		Description: "A great tool.",
		Published:   true,
	}
}

func TestRecsService_Create_Valid(t *testing.T) {
	svc := newRecsService(t)
	rec, err := svc.Create(validCreateInput())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rec.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
}

func TestRecsService_Create_MissingName(t *testing.T) {
	svc := newRecsService(t)
	input := validCreateInput()
	input.Name = ""
	_, err := svc.Create(input)
	if err != model.ErrRecValidation {
		t.Fatalf("expected ErrRecValidation, got %v", err)
	}
}

func TestRecsService_Create_MissingSection(t *testing.T) {
	svc := newRecsService(t)
	input := validCreateInput()
	input.Section = ""
	_, err := svc.Create(input)
	if err != model.ErrRecValidation {
		t.Fatalf("expected ErrRecValidation, got %v", err)
	}
}

func TestRecsService_Update_PartialFields(t *testing.T) {
	svc := newRecsService(t)
	rec, _ := svc.Create(validCreateInput())

	newName := "Updated"
	updated, err := svc.Update(rec.ID, model.UpdateRecInput{Name: &newName})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Name != "Updated" {
		t.Fatalf("expected Updated, got %q", updated.Name)
	}
	if updated.Section != rec.Section {
		t.Fatal("section should not have changed")
	}
}

func TestRecsService_Update_NotFound(t *testing.T) {
	svc := newRecsService(t)
	newName := "x"
	_, err := svc.Update(9999, model.UpdateRecInput{Name: &newName})
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRecsService_Update_ClearsHref(t *testing.T) {
	svc := newRecsService(t)
	href := "https://example.com"
	input := validCreateInput()
	input.Href = &href
	rec, _ := svc.Create(input)

	empty := ""
	updated, err := svc.Update(rec.ID, model.UpdateRecInput{Href: &empty})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Href != nil {
		t.Fatalf("expected nil href after clearing, got %v", updated.Href)
	}
}

func TestRecsService_Delete(t *testing.T) {
	svc := newRecsService(t)
	rec, _ := svc.Create(validCreateInput())

	if err := svc.Delete(rec.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	_, err := svc.GetByID(rec.ID)
	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestRecsService_ListPublished(t *testing.T) {
	svc := newRecsService(t)

	pub := validCreateInput()
	pub.Published = true
	svc.Create(pub)

	draft := validCreateInput()
	draft.Name = "Draft"
	draft.Published = false
	svc.Create(draft)

	recs, err := svc.ListPublished()
	if err != nil {
		t.Fatalf("list published: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 published rec, got %d", len(recs))
	}
}
