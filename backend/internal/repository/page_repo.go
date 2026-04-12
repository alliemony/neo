package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alliemony/neo/backend/internal/model"
)

// PageRepo implements page data access against SQLite.
type PageRepo struct {
	db *sql.DB
}

// NewPageRepo creates a new PageRepo.
func NewPageRepo(db *sql.DB) *PageRepo {
	return &PageRepo{db: db}
}

// Create inserts a new page.
func (r *PageRepo) Create(page *model.Page) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`INSERT INTO pages (slug, title, content, published, sort_order, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		page.Slug, page.Title, page.Content, page.Published, page.SortOrder, now, now,
	)
	if err != nil {
		return fmt.Errorf("insert page: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}

	page.ID = id
	page.CreatedAt = now
	page.UpdatedAt = now
	return nil
}

// GetBySlug returns a single page by its slug.
func (r *PageRepo) GetBySlug(slug string) (*model.Page, error) {
	row := r.db.QueryRow(
		`SELECT id, slug, title, content, published, sort_order, created_at, updated_at
		 FROM pages WHERE slug = ?`, slug,
	)
	return scanPage(row)
}

// ListPublished returns all published pages ordered by sort_order.
func (r *PageRepo) ListPublished() ([]model.Page, error) {
	rows, err := r.db.Query(
		`SELECT id, slug, title, content, published, sort_order, created_at, updated_at
		 FROM pages WHERE published = 1 ORDER BY sort_order ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list pages: %w", err)
	}
	defer rows.Close()

	return scanPages(rows)
}

// ListAll returns all pages ordered by sort_order (for admin).
func (r *PageRepo) ListAll() ([]model.Page, error) {
	rows, err := r.db.Query(
		`SELECT id, slug, title, content, published, sort_order, created_at, updated_at
		 FROM pages ORDER BY sort_order ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list all pages: %w", err)
	}
	defer rows.Close()

	return scanPages(rows)
}

// Update updates an existing page by slug.
func (r *PageRepo) Update(page *model.Page) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(
		`UPDATE pages SET title = ?, content = ?, published = ?, sort_order = ?, updated_at = ?
		 WHERE slug = ?`,
		page.Title, page.Content, page.Published, page.SortOrder, now, page.Slug,
	)
	if err != nil {
		return fmt.Errorf("update page: %w", err)
	}
	page.UpdatedAt = now
	return nil
}

// Delete removes a page by slug.
func (r *PageRepo) Delete(slug string) error {
	result, err := r.db.Exec("DELETE FROM pages WHERE slug = ?", slug)
	if err != nil {
		return fmt.Errorf("delete page: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

func scanPage(row *sql.Row) (*model.Page, error) {
	var p model.Page
	var published int

	err := row.Scan(&p.ID, &p.Slug, &p.Title, &p.Content, &published, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan page: %w", err)
	}

	p.Published = published != 0
	return &p, nil
}

func scanPages(rows *sql.Rows) ([]model.Page, error) {
	var pages []model.Page
	for rows.Next() {
		var p model.Page
		var published int
		if err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Content, &published, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan page row: %w", err)
		}
		p.Published = published != 0
		pages = append(pages, p)
	}
	return pages, rows.Err()
}
