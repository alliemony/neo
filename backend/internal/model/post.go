package model

import (
	"errors"
	"time"
)

// Post represents a blog post.
type Post struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	ContentType string    `json:"content_type"`
	Tags        []string  `json:"tags"`
	Published   bool      `json:"published"`
	LikeCount   int       `json:"like_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TagCount represents a tag name with its associated post count.
type TagCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// ListOptions configures pagination for list queries.
type ListOptions struct {
	Page    int
	PerPage int
}

// CreatePostInput is the input for creating a new post.
type CreatePostInput struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	ContentType string   `json:"content_type"`
	Tags        []string `json:"tags"`
	Published   bool     `json:"published"`
}

// UpdatePostInput is the input for updating an existing post.
type UpdatePostInput struct {
	Title       *string  `json:"title"`
	Content     *string  `json:"content"`
	ContentType *string  `json:"content_type"`
	Tags        []string `json:"tags"`
	Published   *bool    `json:"published"`
}

var (
	ErrNotFound      = errors.New("not found")
	ErrTitleRequired = errors.New("title is required")
	ErrSlugExists    = errors.New("slug already exists")
)
