package service

import (
	"strings"
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

func setupCommentService(t *testing.T) (*CommentService, *PostService) {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("setup test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	postRepo := repository.NewPostRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	postSvc := NewPostService(postRepo)
	commentSvc := NewCommentService(commentRepo, postRepo)
	return commentSvc, postSvc
}

func TestCommentService_Create(t *testing.T) {
	commentSvc, postSvc := setupCommentService(t)
	postSvc.Create(model.CreatePostInput{Title: "Test Post", Content: "c", Published: true})

	comment, err := commentSvc.Create("test-post", model.CreateCommentInput{
		AuthorName: "alice",
		Content:    "Great post!",
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if comment.AuthorName != "alice" {
		t.Errorf("author = %q, want alice", comment.AuthorName)
	}
}

func TestCommentService_Create_EmptyAuthorRejects(t *testing.T) {
	commentSvc, postSvc := setupCommentService(t)
	postSvc.Create(model.CreatePostInput{Title: "Test Post", Content: "c", Published: true})

	_, err := commentSvc.Create("test-post", model.CreateCommentInput{
		AuthorName: "",
		Content:    "content",
	})
	if err != model.ErrAuthorRequired {
		t.Errorf("expected ErrAuthorRequired, got %v", err)
	}
}

func TestCommentService_Create_EmptyContentRejects(t *testing.T) {
	commentSvc, postSvc := setupCommentService(t)
	postSvc.Create(model.CreatePostInput{Title: "Test Post", Content: "c", Published: true})

	_, err := commentSvc.Create("test-post", model.CreateCommentInput{
		AuthorName: "alice",
		Content:    "",
	})
	if err != model.ErrContentRequired {
		t.Errorf("expected ErrContentRequired, got %v", err)
	}
}

func TestCommentService_Create_ContentTooLongRejects(t *testing.T) {
	commentSvc, postSvc := setupCommentService(t)
	postSvc.Create(model.CreatePostInput{Title: "Test Post", Content: "c", Published: true})

	longContent := strings.Repeat("a", 2001)
	_, err := commentSvc.Create("test-post", model.CreateCommentInput{
		AuthorName: "alice",
		Content:    longContent,
	})
	if err != model.ErrContentTooLong {
		t.Errorf("expected ErrContentTooLong, got %v", err)
	}
}

func TestCommentService_Create_PostNotFoundRejects(t *testing.T) {
	commentSvc, _ := setupCommentService(t)

	_, err := commentSvc.Create("nonexistent", model.CreateCommentInput{
		AuthorName: "alice",
		Content:    "hello",
	})
	if err != model.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCommentService_ListByPostSlug(t *testing.T) {
	commentSvc, postSvc := setupCommentService(t)
	postSvc.Create(model.CreatePostInput{Title: "Test Post", Content: "c", Published: true})

	commentSvc.Create("test-post", model.CreateCommentInput{AuthorName: "alice", Content: "one"})
	commentSvc.Create("test-post", model.CreateCommentInput{AuthorName: "bob", Content: "two"})

	comments, err := commentSvc.ListByPostSlug("test-post")
	if err != nil {
		t.Fatalf("ListByPostSlug error: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("len(comments) = %d, want 2", len(comments))
	}
}
