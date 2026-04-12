package service

import (
	"github.com/alliemony/neo/backend/internal/model"
	"github.com/alliemony/neo/backend/internal/repository"
)

// CommentService contains business logic for comments.
type CommentService struct {
	commentRepo *repository.CommentRepo
	postRepo    *repository.PostRepo
}

// NewCommentService creates a new CommentService.
func NewCommentService(commentRepo *repository.CommentRepo, postRepo *repository.PostRepo) *CommentService {
	return &CommentService{commentRepo: commentRepo, postRepo: postRepo}
}

// Create validates input and creates a comment on the given post.
func (s *CommentService) Create(slug string, input model.CreateCommentInput) (*model.Comment, error) {
	if input.AuthorName == "" {
		return nil, model.ErrAuthorRequired
	}
	if input.Content == "" {
		return nil, model.ErrContentRequired
	}
	if len(input.Content) > 2000 {
		return nil, model.ErrContentTooLong
	}

	post, err := s.postRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	comment := &model.Comment{
		PostID:     post.ID,
		AuthorName: input.AuthorName,
		Content:    input.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

// ListByPostSlug returns all comments for a post.
func (s *CommentService) ListByPostSlug(slug string) ([]model.Comment, error) {
	return s.commentRepo.ListByPostSlug(slug)
}
