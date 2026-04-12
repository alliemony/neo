package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alliemony/neo/backend/internal/model"
)

// CommentRepo implements comment data access against SQLite.
type CommentRepo struct {
	db *sql.DB
}

// NewCommentRepo creates a new CommentRepo.
func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

// Create inserts a new comment.
func (r *CommentRepo) Create(comment *model.Comment) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`INSERT INTO comments (post_id, author_name, content, created_at)
		 VALUES (?, ?, ?, ?)`,
		comment.PostID, comment.AuthorName, comment.Content, now,
	)
	if err != nil {
		return fmt.Errorf("insert comment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	comment.ID = id
	comment.CreatedAt = now
	return nil
}

// ListByPostSlug returns all comments for a post ordered by created_at ascending.
func (r *CommentRepo) ListByPostSlug(slug string) ([]model.Comment, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.post_id, c.author_name, c.content, c.created_at
		 FROM comments c
		 JOIN posts p ON p.id = c.post_id
		 WHERE p.slug = ?
		 ORDER BY c.created_at ASC`, slug,
	)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorName, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
