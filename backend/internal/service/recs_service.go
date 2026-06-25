package service

import (
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

// RecsService contains business logic for recommendations.
type RecsService struct {
	repo *repository.RecsRepo
}

// NewRecsService creates a new RecsService.
func NewRecsService(repo *repository.RecsRepo) *RecsService {
	return &RecsService{repo: repo}
}

// ListPublished returns all published recs.
func (s *RecsService) ListPublished() ([]model.Rec, error) {
	return s.repo.ListPublished()
}

// ListAll returns all recs (for admin).
func (s *RecsService) ListAll() ([]model.Rec, error) {
	return s.repo.ListAll()
}

// GetByID returns a rec by its id.
func (s *RecsService) GetByID(id int64) (*model.Rec, error) {
	return s.repo.GetByID(id)
}

// Create validates input and creates a new rec.
func (s *RecsService) Create(input model.CreateRecInput) (*model.Rec, error) {
	if err := validateRecFields(input.Name, input.Section, input.Tag, input.TagBg, input.TagColor, input.Description); err != nil {
		return nil, err
	}

	rec := &model.Rec{
		Name:        input.Name,
		Href:        input.Href,
		Section:     input.Section,
		Tag:         input.Tag,
		TagBg:       input.TagBg,
		TagColor:    input.TagColor,
		Description: input.Description,
		Published:   input.Published,
		SortOrder:   input.SortOrder,
	}

	if err := s.repo.Create(rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// Update applies a partial update to an existing rec.
func (s *RecsService) Update(id int64, input model.UpdateRecInput) (*model.Rec, error) {
	rec, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		rec.Name = *input.Name
	}
	if input.Href != nil {
		if *input.Href == "" {
			rec.Href = nil
		} else {
			rec.Href = input.Href
		}
	}
	if input.Section != nil {
		rec.Section = *input.Section
	}
	if input.Tag != nil {
		rec.Tag = *input.Tag
	}
	if input.TagBg != nil {
		rec.TagBg = *input.TagBg
	}
	if input.TagColor != nil {
		rec.TagColor = *input.TagColor
	}
	if input.Description != nil {
		rec.Description = *input.Description
	}
	if input.Published != nil {
		rec.Published = *input.Published
	}
	if input.SortOrder != nil {
		rec.SortOrder = *input.SortOrder
	}

	if err := validateRecFields(rec.Name, rec.Section, rec.Tag, rec.TagBg, rec.TagColor, rec.Description); err != nil {
		return nil, err
	}

	if err := s.repo.Update(rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// Delete removes a rec by id.
func (s *RecsService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func validateRecFields(name, section, tag, tagBg, tagColor, description string) error {
	if name == "" || section == "" || tag == "" || tagBg == "" || tagColor == "" || description == "" {
		return model.ErrRecValidation
	}
	return nil
}
