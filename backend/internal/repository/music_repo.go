package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alliemony/neo/backend/internal/model"
)

// MusicRepo manages music_recs data access.
type MusicRepo struct {
	db *sql.DB
}

// NewMusicRepo creates a new MusicRepo.
func NewMusicRepo(db *sql.DB) *MusicRepo {
	return &MusicRepo{db: db}
}

// ListAll returns all music recs ordered by sort_order ASC, created_at DESC.
func (r *MusicRepo) ListAll() ([]model.MusicRec, error) {
	rows, err := r.db.Query(
		`SELECT id, title, artist, album, cover_url, spotify_url, apple_url, note, published, sort_order, created_at, updated_at
		 FROM music_recs ORDER BY sort_order ASC, created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list music recs: %w", err)
	}
	defer rows.Close()
	return scanMusicRecs(rows)
}

// ListPublished returns only published music recs ordered by sort_order ASC, created_at DESC.
func (r *MusicRepo) ListPublished() ([]model.MusicRec, error) {
	rows, err := r.db.Query(
		`SELECT id, title, artist, album, cover_url, spotify_url, apple_url, note, published, sort_order, created_at, updated_at
		 FROM music_recs WHERE published = 1 ORDER BY sort_order ASC, created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list published music recs: %w", err)
	}
	defer rows.Close()
	return scanMusicRecs(rows)
}

// GetByID returns a single music rec by id.
func (r *MusicRepo) GetByID(id int64) (*model.MusicRec, error) {
	row := r.db.QueryRow(
		`SELECT id, title, artist, album, cover_url, spotify_url, apple_url, note, published, sort_order, created_at, updated_at
		 FROM music_recs WHERE id = ?`, id,
	)
	return scanMusicRec(row)
}

// Create inserts a new music rec and sets its ID and timestamps.
func (r *MusicRepo) Create(rec *model.MusicRec) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`INSERT INTO music_recs (title, artist, album, cover_url, spotify_url, apple_url, note, published, sort_order, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.Title, rec.Artist, rec.Album, rec.CoverURL, rec.SpotifyURL, rec.AppleURL, rec.Note,
		rec.Published, rec.SortOrder, now, now,
	)
	if err != nil {
		return fmt.Errorf("insert music rec: %w", err)
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

// Update saves changes to an existing music rec by id.
func (r *MusicRepo) Update(rec *model.MusicRec) error {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`UPDATE music_recs SET title=?, artist=?, album=?, cover_url=?, spotify_url=?, apple_url=?, note=?, published=?, sort_order=?, updated_at=?
		 WHERE id=?`,
		rec.Title, rec.Artist, rec.Album, rec.CoverURL, rec.SpotifyURL, rec.AppleURL, rec.Note,
		rec.Published, rec.SortOrder, now, rec.ID,
	)
	if err != nil {
		return fmt.Errorf("update music rec: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	rec.UpdatedAt = now
	return nil
}

// Delete removes a music rec by id.
func (r *MusicRepo) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM music_recs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete music rec: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

func scanMusicRec(row *sql.Row) (*model.MusicRec, error) {
	var rec model.MusicRec
	var published int
	err := row.Scan(
		&rec.ID, &rec.Title, &rec.Artist, &rec.Album, &rec.CoverURL,
		&rec.SpotifyURL, &rec.AppleURL, &rec.Note,
		&published, &rec.SortOrder, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan music rec: %w", err)
	}
	rec.Published = published != 0
	return &rec, nil
}

func scanMusicRecs(rows *sql.Rows) ([]model.MusicRec, error) {
	var recs []model.MusicRec
	for rows.Next() {
		var rec model.MusicRec
		var published int
		if err := rows.Scan(
			&rec.ID, &rec.Title, &rec.Artist, &rec.Album, &rec.CoverURL,
			&rec.SpotifyURL, &rec.AppleURL, &rec.Note,
			&published, &rec.SortOrder, &rec.CreatedAt, &rec.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan music rec row: %w", err)
		}
		rec.Published = published != 0
		recs = append(recs, rec)
	}
	return recs, rows.Err()
}
