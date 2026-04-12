package model

import (
	"errors"
	"time"
)

// Comment represents a comment on a blog post.
type Comment struct {
	ID         int64     `json:"id"`
	PostID     int64     `json:"post_id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateCommentInput is the input for creating a new comment.
type CreateCommentInput struct {
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
}

var (
	ErrAuthorRequired = errors.New("author name is required")
	ErrContentRequired = errors.New("content is required")
	ErrContentTooLong  = errors.New("content must be 2000 characters or less")
)
