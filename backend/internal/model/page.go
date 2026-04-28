package model

import "time"

// Page represents a static content page.
type Page struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	ContentType string    `json:"content_type"`
	Published   bool      `json:"published"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePageInput is the input for creating a new page.
type CreatePageInput struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
	SortOrder int    `json:"sort_order"`
}

// UpdatePageInput is the input for updating an existing page.
type UpdatePageInput struct {
	Title     *string `json:"title"`
	Content   *string `json:"content"`
	Published *bool   `json:"published"`
	SortOrder *int    `json:"sort_order"`
}
