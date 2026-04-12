package service

import (
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

func setupTestService(t *testing.T) *PostService {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("setup test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	repo := repository.NewPostRepo(db)
	return NewPostService(repo)
}

func TestPostService_Create_GeneratesSlug(t *testing.T) {
	svc := setupTestService(t)

	post, err := svc.Create(model.CreatePostInput{
		Title:   "My Great Post!",
		Content: "Some content",
		Tags:    []string{"go"},
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if post.Slug != "my-great-post" {
		t.Errorf("slug = %q, want %q", post.Slug, "my-great-post")
	}
	if post.ContentType != "markdown" {
		t.Errorf("content_type = %q, want %q", post.ContentType, "markdown")
	}
}

func TestPostService_Create_EmptyTitleRejects(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.Create(model.CreatePostInput{Title: "", Content: "c"})
	if err != model.ErrTitleRequired {
		t.Errorf("expected ErrTitleRequired, got %v", err)
	}
}

func TestPostService_Create_DuplicateSlugRejects(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.Create(model.CreatePostInput{Title: "Hello World", Content: "first"})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err = svc.Create(model.CreatePostInput{Title: "Hello World", Content: "second"})
	if err != model.ErrSlugExists {
		t.Errorf("expected ErrSlugExists, got %v", err)
	}
}

func TestPostService_GetBySlug(t *testing.T) {
	svc := setupTestService(t)

	svc.Create(model.CreatePostInput{Title: "Test Post", Content: "content", Published: true})

	post, err := svc.GetBySlug("test-post")
	if err != nil {
		t.Fatalf("GetBySlug error: %v", err)
	}
	if post.Title != "Test Post" {
		t.Errorf("title = %q, want %q", post.Title, "Test Post")
	}
}

func TestPostService_GetBySlug_NotFound(t *testing.T) {
	svc := setupTestService(t)

	_, err := svc.GetBySlug("nonexistent")
	if err != model.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPostService_ListPublished(t *testing.T) {
	svc := setupTestService(t)

	svc.Create(model.CreatePostInput{Title: "Published", Content: "c", Published: true})
	svc.Create(model.CreatePostInput{Title: "Draft", Content: "c", Published: false})

	posts, total, err := svc.ListPublished(model.ListOptions{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("ListPublished error: %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if len(posts) != 1 {
		t.Errorf("len(posts) = %d, want 1", len(posts))
	}
}

func TestPostService_ListByTag(t *testing.T) {
	svc := setupTestService(t)

	svc.Create(model.CreatePostInput{Title: "Go Post", Content: "c", Tags: []string{"go"}, Published: true})
	svc.Create(model.CreatePostInput{Title: "Python Post", Content: "c", Tags: []string{"python"}, Published: true})

	posts, total, err := svc.ListByTag("go", model.ListOptions{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("ListByTag error: %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if len(posts) != 1 || posts[0].Title != "Go Post" {
		t.Errorf("unexpected posts: %v", posts)
	}
}

func TestPostService_AllTags(t *testing.T) {
	svc := setupTestService(t)

	svc.Create(model.CreatePostInput{Title: "P1", Content: "c", Tags: []string{"go", "tutorial"}, Published: true})
	svc.Create(model.CreatePostInput{Title: "P2", Content: "c", Tags: []string{"go"}, Published: true})

	tags, err := svc.AllTags()
	if err != nil {
		t.Fatalf("AllTags error: %v", err)
	}

	tagMap := make(map[string]int)
	for _, tc := range tags {
		tagMap[tc.Name] = tc.Count
	}
	if tagMap["go"] != 2 {
		t.Errorf("go count = %d, want 2", tagMap["go"])
	}
	if tagMap["tutorial"] != 1 {
		t.Errorf("tutorial count = %d, want 1", tagMap["tutorial"])
	}
}
