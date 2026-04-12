package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alliemony/neo/backend/internal/model"
)

// PostRepo implements post data access against SQLite.
type PostRepo struct {
	db *sql.DB
}

// NewPostRepo creates a new PostRepo.
func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

// Create inserts a new post and populates ID and timestamps.
func (r *PostRepo) Create(post *model.Post) error {
	tagsJSON, err := json.Marshal(post.Tags)
	if err != nil {
		return fmt.Errorf("marshal tags: %w", err)
	}

	now := time.Now().UTC()
	result, err := r.db.Exec(
		`INSERT INTO posts (slug, title, content, content_type, tags, published, like_count, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		post.Slug, post.Title, post.Content, post.ContentType, string(tagsJSON), post.Published, post.LikeCount, now, now,
	)
	if err != nil {
		return fmt.Errorf("insert post: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	post.ID = id
	post.CreatedAt = now
	post.UpdatedAt = now
	return nil
}

// GetBySlug returns a single post by its slug.
func (r *PostRepo) GetBySlug(slug string) (*model.Post, error) {
	row := r.db.QueryRow(
		`SELECT id, slug, title, content, content_type, tags, published, like_count, created_at, updated_at
		 FROM posts WHERE slug = ?`, slug,
	)
	return scanPost(row)
}

// List returns paginated posts, optionally filtered to published only.
func (r *PostRepo) List(opts model.ListOptions, publishedOnly bool) ([]model.Post, int, error) {
	opts = normalizeOpts(opts)

	var countQuery, listQuery string
	var args []any

	if publishedOnly {
		countQuery = "SELECT COUNT(*) FROM posts WHERE published = 1"
		listQuery = `SELECT id, slug, title, content, content_type, tags, published, like_count, created_at, updated_at
		             FROM posts WHERE published = 1 ORDER BY created_at DESC LIMIT ? OFFSET ?`
	} else {
		countQuery = "SELECT COUNT(*) FROM posts"
		listQuery = `SELECT id, slug, title, content, content_type, tags, published, like_count, created_at, updated_at
		             FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?`
	}

	var total int
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count posts: %w", err)
	}

	offset := (opts.Page - 1) * opts.PerPage
	args = append(args, opts.PerPage, offset)

	rows, err := r.db.Query(listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	return scanPosts(rows, total)
}

// ListByTag returns paginated published posts that contain the given tag.
func (r *PostRepo) ListByTag(tag string, opts model.ListOptions) ([]model.Post, int, error) {
	opts = normalizeOpts(opts)

	// SQLite JSON: tags is stored as a JSON array like '["go","python"]'.
	// We use json_each to check membership.
	countQuery := `SELECT COUNT(*) FROM posts, json_each(posts.tags) WHERE json_each.value = ? AND published = 1`

	var total int
	if err := r.db.QueryRow(countQuery, tag).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count posts by tag: %w", err)
	}

	offset := (opts.Page - 1) * opts.PerPage
	listQuery := `SELECT posts.id, posts.slug, posts.title, posts.content, posts.content_type,
	              posts.tags, posts.published, posts.like_count, posts.created_at, posts.updated_at
	              FROM posts, json_each(posts.tags)
	              WHERE json_each.value = ? AND posts.published = 1
	              ORDER BY posts.created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(listQuery, tag, opts.PerPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list posts by tag: %w", err)
	}
	defer rows.Close()

	return scanPosts(rows, total)
}

// AllTags returns all unique tags with their post counts (published posts only).
func (r *PostRepo) AllTags() ([]model.TagCount, error) {
	rows, err := r.db.Query(
		`SELECT json_each.value AS tag, COUNT(*) AS cnt
		 FROM posts, json_each(posts.tags)
		 WHERE posts.published = 1
		 GROUP BY json_each.value
		 ORDER BY cnt DESC, tag ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("query all tags: %w", err)
	}
	defer rows.Close()

	var tags []model.TagCount
	for rows.Next() {
		var tc model.TagCount
		if err := rows.Scan(&tc.Name, &tc.Count); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, tc)
	}
	return tags, rows.Err()
}

// Update updates an existing post by slug.
func (r *PostRepo) Update(post *model.Post) error {
	tagsJSON, err := json.Marshal(post.Tags)
	if err != nil {
		return fmt.Errorf("marshal tags: %w", err)
	}

	now := time.Now().UTC()
	_, err = r.db.Exec(
		`UPDATE posts SET title = ?, content = ?, content_type = ?, tags = ?, published = ?, updated_at = ?
		 WHERE slug = ?`,
		post.Title, post.Content, post.ContentType, string(tagsJSON), post.Published, now, post.Slug,
	)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}
	post.UpdatedAt = now
	return nil
}

// Delete removes a post by slug.
func (r *PostRepo) Delete(slug string) error {
	result, err := r.db.Exec("DELETE FROM posts WHERE slug = ?", slug)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

// IncrementLikeCount atomically increments a post's like_count and returns the new value.
func (r *PostRepo) IncrementLikeCount(slug string) (int, error) {
	result, err := r.db.Exec(
		`UPDATE posts SET like_count = like_count + 1 WHERE slug = ?`, slug,
	)
	if err != nil {
		return 0, fmt.Errorf("increment like count: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return 0, model.ErrNotFound
	}

	var count int
	err = r.db.QueryRow(`SELECT like_count FROM posts WHERE slug = ?`, slug).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("read like count: %w", err)
	}
	return count, nil
}

// scanPost scans a single row into a Post.
func scanPost(row *sql.Row) (*model.Post, error) {
	var p model.Post
	var tagsJSON string
	var published int

	err := row.Scan(&p.ID, &p.Slug, &p.Title, &p.Content, &p.ContentType, &tagsJSON, &published, &p.LikeCount, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan post: %w", err)
	}

	p.Published = published != 0
	if err := json.Unmarshal([]byte(tagsJSON), &p.Tags); err != nil {
		return nil, fmt.Errorf("unmarshal tags: %w", err)
	}
	return &p, nil
}

// scanPosts scans multiple rows into a slice of Post.
func scanPosts(rows *sql.Rows, total int) ([]model.Post, int, error) {
	var posts []model.Post
	for rows.Next() {
		var p model.Post
		var tagsJSON string
		var published int
		if err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Content, &p.ContentType, &tagsJSON, &published, &p.LikeCount, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan post row: %w", err)
		}
		p.Published = published != 0
		if err := json.Unmarshal([]byte(tagsJSON), &p.Tags); err != nil {
			return nil, 0, fmt.Errorf("unmarshal tags: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, total, rows.Err()
}

func normalizeOpts(opts model.ListOptions) model.ListOptions {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PerPage < 1 {
		opts.PerPage = 10
	}
	return opts
}
