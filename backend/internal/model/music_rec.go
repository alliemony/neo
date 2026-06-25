package model

import (
	"errors"
	"time"
)

// MusicRec represents a curated music recommendation.
type MusicRec struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Artist     string    `json:"artist"`
	Album      *string   `json:"album"`
	CoverURL   *string   `json:"cover_url"`
	SpotifyURL *string   `json:"spotify_url"`
	AppleURL   *string   `json:"apple_url"`
	Note       *string   `json:"note"`
	Published  bool      `json:"published"`
	SortOrder  int       `json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateMusicRecInput is the input for creating a new music rec.
type CreateMusicRecInput struct {
	Title      string  `json:"title"`
	Artist     string  `json:"artist"`
	Album      *string `json:"album"`
	CoverURL   *string `json:"cover_url"`
	SpotifyURL *string `json:"spotify_url"`
	AppleURL   *string `json:"apple_url"`
	Note       *string `json:"note"`
	Published  bool    `json:"published"`
	SortOrder  int     `json:"sort_order"`
}

// UpdateMusicRecInput is the input for updating a music rec (all fields optional).
type UpdateMusicRecInput struct {
	Title      *string `json:"title"`
	Artist     *string `json:"artist"`
	Album      *string `json:"album"`
	CoverURL   *string `json:"cover_url"`
	SpotifyURL *string `json:"spotify_url"`
	AppleURL   *string `json:"apple_url"`
	Note       *string `json:"note"`
	Published  *bool   `json:"published"`
	SortOrder  *int    `json:"sort_order"`
}

var ErrMusicRecValidation = errors.New("title and artist are required")
