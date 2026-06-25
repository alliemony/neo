package service

import (
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

// MusicService contains business logic for music recommendations.
type MusicService struct {
	repo *repository.MusicRepo
}

// NewMusicService creates a new MusicService.
func NewMusicService(repo *repository.MusicRepo) *MusicService {
	return &MusicService{repo: repo}
}

// ListPublished returns all published music recs.
func (s *MusicService) ListPublished() ([]model.MusicRec, error) {
	return s.repo.ListPublished()
}

// ListAll returns all music recs (for admin).
func (s *MusicService) ListAll() ([]model.MusicRec, error) {
	return s.repo.ListAll()
}

// GetByID returns a music rec by id.
func (s *MusicService) GetByID(id int64) (*model.MusicRec, error) {
	return s.repo.GetByID(id)
}

// Create validates input and creates a new music rec.
func (s *MusicService) Create(input model.CreateMusicRecInput) (*model.MusicRec, error) {
	if input.Title == "" || input.Artist == "" {
		return nil, model.ErrMusicRecValidation
	}
	rec := &model.MusicRec{
		Title:      input.Title,
		Artist:     input.Artist,
		Album:      input.Album,
		CoverURL:   input.CoverURL,
		SpotifyURL: input.SpotifyURL,
		AppleURL:   input.AppleURL,
		Note:       input.Note,
		Published:  input.Published,
		SortOrder:  input.SortOrder,
	}
	if err := s.repo.Create(rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// Update applies a partial update to an existing music rec.
func (s *MusicService) Update(id int64, input model.UpdateMusicRecInput) (*model.MusicRec, error) {
	rec, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Title != nil {
		rec.Title = *input.Title
	}
	if input.Artist != nil {
		rec.Artist = *input.Artist
	}
	if input.Album != nil {
		rec.Album = input.Album
	}
	if input.CoverURL != nil {
		rec.CoverURL = input.CoverURL
	}
	if input.SpotifyURL != nil {
		rec.SpotifyURL = input.SpotifyURL
	}
	if input.AppleURL != nil {
		rec.AppleURL = input.AppleURL
	}
	if input.Note != nil {
		rec.Note = input.Note
	}
	if input.Published != nil {
		rec.Published = *input.Published
	}
	if input.SortOrder != nil {
		rec.SortOrder = *input.SortOrder
	}
	if rec.Title == "" || rec.Artist == "" {
		return nil, model.ErrMusicRecValidation
	}
	if err := s.repo.Update(rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// Delete removes a music rec by id.
func (s *MusicService) Delete(id int64) error {
	return s.repo.Delete(id)
}
