package database

import (
	"database/sql"
	"fmt"
	"sort"

	_ "modernc.org/sqlite"
)

// New opens a SQLite database and runs all pending migrations.
func New(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Enable WAL mode for better concurrent read performance.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set WAL mode: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	// Create migrations tracking table.
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Get already-applied versions.
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("query migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return fmt.Errorf("scan migration version: %w", err)
		}
		applied[v] = true
	}

	// Run pending migrations in order.
	versions := make([]string, 0, len(migrations))
	for v := range migrations {
		versions = append(versions, v)
	}
	sort.Strings(versions)

	for _, v := range versions {
		if applied[v] {
			continue
		}
		if _, err := db.Exec(migrations[v]); err != nil {
			return fmt.Errorf("apply migration %s: %w", v, err)
		}
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", v); err != nil {
			return fmt.Errorf("record migration %s: %w", v, err)
		}
	}

	return nil
}

// migrations maps version strings to SQL statements.
var migrations = map[string]string{
	"001_create_posts": `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slug TEXT UNIQUE NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			content_type TEXT NOT NULL DEFAULT 'markdown',
			tags TEXT NOT NULL DEFAULT '[]',
			published INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
		CREATE INDEX IF NOT EXISTS idx_posts_published ON posts(published);
		CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
	`,
}
