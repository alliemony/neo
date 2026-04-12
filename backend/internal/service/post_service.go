package service

import (
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

// PostService contains business logic for blog posts.
type PostService struct {
	repo *repository.PostRepo
}

// NewPostService creates a new PostService.
func NewPostService(repo *repository.PostRepo) *PostService {
	return &PostService{repo: repo}
}

// Create validates the input, generates a slug, and creates a post.
func (s *PostService) Create(input model.CreatePostInput) (*model.Post, error) {
	if input.Title == "" {
		return nil, model.ErrTitleRequired
	}

	slug := slugify(input.Title)

	existing, _ := s.repo.GetBySlug(slug)
	if existing != nil {
		return nil, model.ErrSlugExists
	}

	contentType := input.ContentType
	if contentType == "" {
		contentType = "markdown"
	}

	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}

	post := &model.Post{
		Slug:        slug,
		Title:       input.Title,
		Content:     input.Content,
		ContentType: contentType,
		Tags:        tags,
		Published:   input.Published,
	}

	if err := s.repo.Create(post); err != nil {
		return nil, err
	}
	return post, nil
}

// GetBySlug returns a post by its slug.
func (s *PostService) GetBySlug(slug string) (*model.Post, error) {
	return s.repo.GetBySlug(slug)
}

// ListPublished returns paginated published posts.
func (s *PostService) ListPublished(opts model.ListOptions) ([]model.Post, int, error) {
	return s.repo.List(opts, true)
}

// ListByTag returns paginated published posts filtered by tag.
func (s *PostService) ListByTag(tag string, opts model.ListOptions) ([]model.Post, int, error) {
	return s.repo.ListByTag(tag, opts)
}

// AllTags returns all unique tags with counts.
func (s *PostService) AllTags() ([]model.TagCount, error) {
	return s.repo.AllTags()
}

// ListAll returns all posts including drafts (for admin).
func (s *PostService) ListAll(opts model.ListOptions) ([]model.Post, int, error) {
	return s.repo.List(opts, false)
}

// Update updates an existing post by slug.
func (s *PostService) Update(slug string, input model.UpdatePostInput) (*model.Post, error) {
	post, err := s.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.ContentType != nil {
		post.ContentType = *input.ContentType
	}
	if input.Tags != nil {
		post.Tags = input.Tags
	}
	if input.Published != nil {
		post.Published = *input.Published
	}

	if err := s.repo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

// Delete removes a post by slug.
func (s *PostService) Delete(slug string) error {
	return s.repo.Delete(slug)
}
