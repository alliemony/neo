package service

import (
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

// PageService contains business logic for static pages.
type PageService struct {
	repo *repository.PageRepo
}

// NewPageService creates a new PageService.
func NewPageService(repo *repository.PageRepo) *PageService {
	return &PageService{repo: repo}
}

// Create validates input, generates a slug, and creates a page.
func (s *PageService) Create(input model.CreatePageInput) (*model.Page, error) {
	if input.Title == "" {
		return nil, model.ErrTitleRequired
	}

	slug := slugify(input.Title)

	existing, _ := s.repo.GetBySlug(slug)
	if existing != nil {
		return nil, model.ErrSlugExists
	}

	page := &model.Page{
		Slug:        slug,
		Title:       input.Title,
		Content:     input.Content,
		ContentType: "markdown",
		Published:   input.Published,
		SortOrder:   input.SortOrder,
	}

	if err := s.repo.Create(page); err != nil {
		return nil, err
	}
	return page, nil
}

// GetBySlug returns a page by its slug.
func (s *PageService) GetBySlug(slug string) (*model.Page, error) {
	return s.repo.GetBySlug(slug)
}

// ListPublished returns all published pages.
func (s *PageService) ListPublished() ([]model.Page, error) {
	return s.repo.ListPublished()
}

// ListAll returns all pages (for admin).
func (s *PageService) ListAll() ([]model.Page, error) {
	return s.repo.ListAll()
}

// Update updates an existing page.
func (s *PageService) Update(slug string, input model.UpdatePageInput) (*model.Page, error) {
	page, err := s.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		page.Title = *input.Title
	}
	if input.Content != nil {
		page.Content = *input.Content
	}
	if input.Published != nil {
		page.Published = *input.Published
	}
	if input.SortOrder != nil {
		page.SortOrder = *input.SortOrder
	}

	if err := s.repo.Update(page); err != nil {
		return nil, err
	}
	return page, nil
}

// Delete deletes a page by slug.
func (s *PageService) Delete(slug string) error {
	return s.repo.Delete(slug)
}
