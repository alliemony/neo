package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alliemony/neo/backend/internal/model"
)

// RecsRepo manages rec data access.
type RecsRepo struct {
	db *sql.DB
}

// NewRecsRepo creates a new RecsRepo.
func NewRecsRepo(db *sql.DB) *RecsRepo {
	return &RecsRepo{db: db}
}

// ListAll returns all recs ordered by sort_order ASC, created_at ASC.
func (r *RecsRepo) ListAll() ([]model.Rec, error) {
	rows, err := r.db.Query(
		`SELECT id, name, href, section, tag, tag_bg, tag_color, description, published, sort_order, created_at, updated_at
		 FROM recs ORDER BY sort_order ASC, created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list recs: %w", err)
	}
	defer rows.Close()
	return scanRecs(rows)
}

// ListPublished returns only published recs ordered by sort_order ASC, created_at ASC.
func (r *RecsRepo) ListPublished() ([]model.Rec, error) {
	rows, err := r.db.Query(
		`SELECT id, name, href, section, tag, tag_bg, tag_color, description, published, sort_order, created_at, updated_at
		 FROM recs WHERE published = 1 ORDER BY sort_order ASC, created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list published recs: %w", err)
	}
	defer rows.Close()
	return scanRecs(rows)
}

// GetByID returns a single rec by id.
func (r *RecsRepo) GetByID(id int64) (*model.Rec, error) {
	row := r.db.QueryRow(
		`SELECT id, name, href, section, tag, tag_bg, tag_color, description, published, sort_order, created_at, updated_at
		 FROM recs WHERE id = ?`, id,
	)
	return scanRec(row)
}

// Create inserts a new rec and sets its ID and timestamps.
func (r *RecsRepo) Create(rec *model.Rec) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`INSERT INTO recs (name, href, section, tag, tag_bg, tag_color, description, published, sort_order, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.Name, rec.Href, rec.Section, rec.Tag, rec.TagBg, rec.TagColor, rec.Description,
		rec.Published, rec.SortOrder, now, now,
	)
	if err != nil {
		return fmt.Errorf("insert rec: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	rec.ID = id
	rec.CreatedAt = now
	rec.UpdatedAt = now
	return nil
}

// Update saves changes to an existing rec by id.
func (r *RecsRepo) Update(rec *model.Rec) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`UPDATE recs SET name=?, href=?, section=?, tag=?, tag_bg=?, tag_color=?, description=?, published=?, sort_order=?, updated_at=?
		 WHERE id=?`,
		rec.Name, rec.Href, rec.Section, rec.Tag, rec.TagBg, rec.TagColor, rec.Description,
		rec.Published, rec.SortOrder, now, rec.ID,
	)
	if err != nil {
		return fmt.Errorf("update rec: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	rec.UpdatedAt = now
	return nil
}

// Delete removes a rec by id.
func (r *RecsRepo) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM recs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete rec: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

func scanRec(row *sql.Row) (*model.Rec, error) {
	var rec model.Rec
	var published int
	err := row.Scan(
		&rec.ID, &rec.Name, &rec.Href, &rec.Section, &rec.Tag, &rec.TagBg, &rec.TagColor,
		&rec.Description, &published, &rec.SortOrder, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan rec: %w", err)
	}
	rec.Published = published != 0
	return &rec, nil
}

func scanRecs(rows *sql.Rows) ([]model.Rec, error) {
	var recs []model.Rec
	for rows.Next() {
		var rec model.Rec
		var published int
		if err := rows.Scan(
			&rec.ID, &rec.Name, &rec.Href, &rec.Section, &rec.Tag, &rec.TagBg, &rec.TagColor,
			&rec.Description, &published, &rec.SortOrder, &rec.CreatedAt, &rec.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan rec row: %w", err)
		}
		rec.Published = published != 0
		recs = append(recs, rec)
	}
	return recs, rows.Err()
}
