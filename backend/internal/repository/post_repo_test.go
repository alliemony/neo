package repository

import (
	"fmt"
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
)

func setupTestDB(t *testing.T) *PostRepo {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("setup test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewPostRepo(db)
}

func TestPostRepo_Create(t *testing.T) {
	repo := setupTestDB(t)

	post := &model.Post{
		Slug:        "hello-world",
		Title:       "Hello World",
		Content:     "First post",
		ContentType: "markdown",
		Tags:        []string{"go", "intro"},
		Published:   true,
	}
	err := repo.Create(post)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if post.ID == 0 {
		t.Error("expected non-zero ID after create")
	}
	if post.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if post.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}
}

func TestPostRepo_GetBySlug(t *testing.T) {
	repo := setupTestDB(t)

	original := &model.Post{
		Slug:        "hello-world",
		Title:       "Hello World",
		Content:     "First post",
		ContentType: "markdown",
		Tags:        []string{"go"},
		Published:   true,
	}
	if err := repo.Create(original); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	got, err := repo.GetBySlug("hello-world")
	if err != nil {
		t.Fatalf("GetBySlug error: %v", err)
	}
	if got.Title != "Hello World" {
		t.Errorf("title = %q, want %q", got.Title, "Hello World")
	}
	if len(got.Tags) != 1 || got.Tags[0] != "go" {
		t.Errorf("tags = %v, want [go]", got.Tags)
	}
}

func TestPostRepo_GetBySlug_NotFound(t *testing.T) {
	repo := setupTestDB(t)

	_, err := repo.GetBySlug("nonexistent")
	if err != model.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPostRepo_List_PublishedOnly(t *testing.T) {
	repo := setupTestDB(t)

	for i := range 15 {
		post := &model.Post{
			Slug:        fmt.Sprintf("published-%d", i),
			Title:       fmt.Sprintf("Published %d", i),
			Content:     "content",
			ContentType: "markdown",
			Tags:        []string{},
			Published:   true,
		}
		if err := repo.Create(post); err != nil {
			t.Fatalf("create published post: %v", err)
		}
	}
	for i := range 3 {
		post := &model.Post{
			Slug:        fmt.Sprintf("draft-%d", i),
			Title:       fmt.Sprintf("Draft %d", i),
			Content:     "content",
			ContentType: "markdown",
			Tags:        []string{},
			Published:   false,
		}
		if err := repo.Create(post); err != nil {
			t.Fatalf("create draft post: %v", err)
		}
	}

	posts, total, err := repo.List(model.ListOptions{Page: 1, PerPage: 10}, true)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if total != 15 {
		t.Errorf("total = %d, want 15", total)
	}
	if len(posts) != 10 {
		t.Errorf("len(posts) = %d, want 10", len(posts))
	}
	// Should be ordered by created_at descending (newest first).
	for _, p := range posts {
		if !p.Published {
			t.Error("got unpublished post in published-only list")
		}
	}
}

func TestPostRepo_List_IncludeDrafts(t *testing.T) {
	repo := setupTestDB(t)

	repo.Create(&model.Post{Slug: "pub", Title: "Pub", Content: "c", ContentType: "markdown", Tags: []string{}, Published: true})
	repo.Create(&model.Post{Slug: "draft", Title: "Draft", Content: "c", ContentType: "markdown", Tags: []string{}, Published: false})

	posts, total, err := repo.List(model.ListOptions{Page: 1, PerPage: 10}, false)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(posts) != 2 {
		t.Errorf("len(posts) = %d, want 2", len(posts))
	}
}

func TestPostRepo_ListByTag(t *testing.T) {
	repo := setupTestDB(t)

	for i := range 3 {
		repo.Create(&model.Post{
			Slug: fmt.Sprintf("python-%d", i), Title: fmt.Sprintf("Python %d", i),
			Content: "c", ContentType: "markdown", Tags: []string{"python"}, Published: true,
		})
	}
	for i := range 2 {
		repo.Create(&model.Post{
			Slug: fmt.Sprintf("go-%d", i), Title: fmt.Sprintf("Go %d", i),
			Content: "c", ContentType: "markdown", Tags: []string{"go"}, Published: true,
		})
	}

	posts, total, err := repo.ListByTag("python", model.ListOptions{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("ListByTag error: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(posts) != 3 {
		t.Errorf("len(posts) = %d, want 3", len(posts))
	}
}

func TestPostRepo_AllTags(t *testing.T) {
	repo := setupTestDB(t)

	repo.Create(&model.Post{Slug: "p1", Title: "P1", Content: "c", ContentType: "markdown", Tags: []string{"python", "ml"}, Published: true})
	repo.Create(&model.Post{Slug: "p2", Title: "P2", Content: "c", ContentType: "markdown", Tags: []string{"python", "tutorial"}, Published: true})

	tags, err := repo.AllTags()
	if err != nil {
		t.Fatalf("AllTags error: %v", err)
	}

	tagMap := make(map[string]int)
	for _, tc := range tags {
		tagMap[tc.Name] = tc.Count
	}

	if tagMap["python"] != 2 {
		t.Errorf("python count = %d, want 2", tagMap["python"])
	}
	if tagMap["ml"] != 1 {
		t.Errorf("ml count = %d, want 1", tagMap["ml"])
	}
	if tagMap["tutorial"] != 1 {
		t.Errorf("tutorial count = %d, want 1", tagMap["tutorial"])
	}
}

func TestPostRepo_Update(t *testing.T) {
	repo := setupTestDB(t)

	post := &model.Post{
		Slug: "test", Title: "Test", Content: "old",
		ContentType: "markdown", Tags: []string{}, Published: false,
	}
	repo.Create(post)

	post.Content = "new content"
	post.Published = true
	err := repo.Update(post)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}

	got, _ := repo.GetBySlug("test")
	if got.Content != "new content" {
		t.Errorf("content = %q, want %q", got.Content, "new content")
	}
	if !got.Published {
		t.Error("expected published = true after update")
	}
}

func TestPostRepo_Delete(t *testing.T) {
	repo := setupTestDB(t)

	repo.Create(&model.Post{Slug: "del-me", Title: "Delete Me", Content: "c", ContentType: "markdown", Tags: []string{}, Published: true})

	err := repo.Delete("del-me")
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	_, err = repo.GetBySlug("del-me")
	if err != model.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}
