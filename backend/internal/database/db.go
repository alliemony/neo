package database

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// Dialect represents the database type.
type Dialect string

const (
	DialectSQLite   Dialect = "sqlite"
	DialectPostgres Dialect = "postgres"
)

// CurrentDialect holds the active dialect after New() is called.
var CurrentDialect Dialect

// New opens a database connection and runs all pending migrations.
// It detects the driver from the DATABASE_URL scheme:
//   - "postgres://" or "postgresql://" → PostgreSQL
//   - anything else → SQLite
func New(dbPath string) (*sql.DB, error) {
	driver, dsn := parseDriverDSN(dbPath)
	CurrentDialect = Dialect(driver)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if driver == "sqlite" {
		// Enable WAL mode for better concurrent read performance.
		if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
			db.Close()
			return nil, fmt.Errorf("set WAL mode: %w", err)
		}
	}

	if err := migrate(db, Dialect(driver)); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

// parseDriverDSN returns the database/sql driver name and DSN for the given URL.
func parseDriverDSN(dbURL string) (driver, dsn string) {
	if strings.HasPrefix(dbURL, "postgres://") || strings.HasPrefix(dbURL, "postgresql://") {
		return "pgx", dbURL
	}
	// SQLite — strip sqlite:// prefix if present.
	return "sqlite", strings.TrimPrefix(dbURL, "sqlite://")
}

func migrate(db *sql.DB, dialect Dialect) error {
	// Create migrations tracking table.
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

	// Select migration set based on dialect.
	m := sqliteMigrations
	if dialect == DialectPostgres {
		m = postgresMigrations
	}

	// Run pending migrations in order.
	versions := make([]string, 0, len(m))
	for v := range m {
		versions = append(versions, v)
	}
	sort.Strings(versions)

	for _, v := range versions {
		if applied[v] {
			continue
		}
		if _, err := db.Exec(m[v]); err != nil {
			return fmt.Errorf("apply migration %s: %w", v, err)
		}
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", v); err != nil {
			return fmt.Errorf("record migration %s: %w", v, err)
		}
	}

	return nil
}

// sqliteMigrations are SQLite-specific DDL.
var sqliteMigrations = map[string]string{
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
	"002_add_like_count": `
		ALTER TABLE posts ADD COLUMN like_count INTEGER NOT NULL DEFAULT 0;
	`,
	"003_create_comments": `
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			author_name TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
	`,
	"004_create_pages": `
		CREATE TABLE IF NOT EXISTS pages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slug TEXT UNIQUE NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			published INTEGER NOT NULL DEFAULT 0,
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_pages_slug ON pages(slug);
		CREATE INDEX IF NOT EXISTS idx_pages_published ON pages(published);
	`,
	"005_add_content_type_to_pages": `
		ALTER TABLE pages ADD COLUMN content_type TEXT NOT NULL DEFAULT 'markdown';
	`,
}

// postgresMigrations are PostgreSQL-specific DDL.
var postgresMigrations = map[string]string{
	"001_create_posts": `
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			slug TEXT UNIQUE NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			content_type TEXT NOT NULL DEFAULT 'markdown',
			tags TEXT NOT NULL DEFAULT '[]',
			published BOOLEAN NOT NULL DEFAULT FALSE,
			like_count INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
		CREATE INDEX IF NOT EXISTS idx_posts_published ON posts(published);
		CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
	`,
	"003_create_comments": `
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			author_name TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
	`,
	"004_create_pages": `
		CREATE TABLE IF NOT EXISTS pages (
			id SERIAL PRIMARY KEY,
			slug TEXT UNIQUE NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			published BOOLEAN NOT NULL DEFAULT FALSE,
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_pages_slug ON pages(slug);
		CREATE INDEX IF NOT EXISTS idx_pages_published ON pages(published);
	`,
	"005_add_content_type_to_pages": `
		ALTER TABLE pages ADD COLUMN content_type TEXT NOT NULL DEFAULT 'markdown';
	`,
}
