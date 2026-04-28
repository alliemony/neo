package repository

import (
	"testing"

	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/model"
)

func setupCommentTestDB(t *testing.T) (*CommentRepo, *PostRepo) {
	t.Helper()
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("setup test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewCommentRepo(db), NewPostRepo(db)
}

func createTestPost(t *testing.T, postRepo *PostRepo, slug string) *model.Post {
	t.Helper()
	post := &model.Post{
		Slug: slug, Title: slug, Content: "content",
		ContentType: "markdown", Tags: []string{}, Published: true,
	}
	if err := postRepo.Create(post); err != nil {
		t.Fatalf("create test post: %v", err)
	}
	return post
}

func TestCommentRepo_Create(t *testing.T) {
	commentRepo, postRepo := setupCommentTestDB(t)
	post := createTestPost(t, postRepo, "hello-world")

	comment := &model.Comment{
		PostID:     post.ID,
		AuthorName: "alice",
		Content:    "Great post!",
	}
	err := commentRepo.Create(comment)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if comment.ID == 0 {
		t.Error("expected non-zero ID after create")
	}
	if comment.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestCommentRepo_ListByPostSlug(t *testing.T) {
	commentRepo, postRepo := setupCommentTestDB(t)
	post := createTestPost(t, postRepo, "hello-world")

	for _, name := range []string{"alice", "bob", "charlie"} {
		commentRepo.Create(&model.Comment{
			PostID: post.ID, AuthorName: name, Content: "comment by " + name,
		})
	}

	comments, err := commentRepo.ListByPostSlug("hello-world")
	if err != nil {
		t.Fatalf("ListByPostSlug error: %v", err)
	}
	if len(comments) != 3 {
		t.Errorf("len(comments) = %d, want 3", len(comments))
	}
	// Should be ordered oldest first.
	if comments[0].AuthorName != "alice" {
		t.Errorf("first comment author = %q, want alice", comments[0].AuthorName)
	}
	if comments[2].AuthorName != "charlie" {
		t.Errorf("last comment author = %q, want charlie", comments[2].AuthorName)
	}
}

func TestCommentRepo_ListByPostSlug_EmptyForNewPost(t *testing.T) {
	commentRepo, postRepo := setupCommentTestDB(t)
	createTestPost(t, postRepo, "no-comments")

	comments, err := commentRepo.ListByPostSlug("no-comments")
	if err != nil {
		t.Fatalf("ListByPostSlug error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}

func TestCommentRepo_CascadeDeleteOnPostDelete(t *testing.T) {
	commentRepo, postRepo := setupCommentTestDB(t)
	post := createTestPost(t, postRepo, "will-delete")

	commentRepo.Create(&model.Comment{
		PostID: post.ID, AuthorName: "alice", Content: "gone soon",
	})

	postRepo.Delete("will-delete")

	comments, err := commentRepo.ListByPostSlug("will-delete")
	if err != nil {
		t.Fatalf("ListByPostSlug error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments after post delete, got %d", len(comments))
	}
}
