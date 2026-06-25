package model

import (
	"errors"
	"time"
)

// Rec represents a curated recommendation.
type Rec struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Href        *string   `json:"href"`
	Section     string    `json:"section"`
	Tag         string    `json:"tag"`
	TagBg       string    `json:"tag_bg"`
	TagColor    string    `json:"tag_color"`
	Description string    `json:"description"`
	Published   bool      `json:"published"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RecSection is a group of recs sharing the same section label.
type RecSection struct {
	Title string `json:"title"`
	Items []Rec  `json:"items"`
}

// CreateRecInput is the input for creating a new rec.
type CreateRecInput struct {
	Name        string  `json:"name"`
	Href        *string `json:"href"`
	Section     string  `json:"section"`
	Tag         string  `json:"tag"`
	TagBg       string  `json:"tag_bg"`
	TagColor    string  `json:"tag_color"`
	Description string  `json:"description"`
	Published   bool    `json:"published"`
	SortOrder   int     `json:"sort_order"`
}

// UpdateRecInput is the input for updating an existing rec (all fields optional).
// Href: nil = don't change; pointer to empty string = clear the link.
type UpdateRecInput struct {
	Name        *string `json:"name"`
	Href        *string `json:"href"`
	Section     *string `json:"section"`
	Tag         *string `json:"tag"`
	TagBg       *string `json:"tag_bg"`
	TagColor    *string `json:"tag_color"`
	Description *string `json:"description"`
	Published   *bool   `json:"published"`
	SortOrder   *int    `json:"sort_order"`
}

var ErrRecValidation = errors.New("name, section, tag, tag_bg, tag_color, and description are required")
