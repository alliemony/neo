package database

import (
	"testing"
)

func TestNew_CreatesDatabase(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("New(:memory:) error: %v", err)
	}
	defer db.Close()

	// Verify posts table exists by querying it.
	_, err = db.Exec("SELECT id, slug, title, content, content_type, tags, published, created_at, updated_at FROM posts LIMIT 0")
	if err != nil {
		t.Fatalf("posts table not created: %v", err)
	}
}

func TestNew_MigrationsAreIdempotent(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("first New error: %v", err)
	}

	// Run migrate again — should not fail.
	if err := migrate(db, DialectSQLite); err != nil {
		t.Fatalf("second migrate error: %v", err)
	}
	db.Close()
}

func TestNew_RecordsMigrationVersions(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if count != len(sqliteMigrations) {
		t.Errorf("expected %d migration(s) recorded, got %d", len(sqliteMigrations), count)
	}
}

func TestParseDriverDSN_SQLite(t *testing.T) {
	tests := []struct {
		input      string
		wantDriver string
		wantDSN    string
	}{
		{"sqlite://neo.db", "sqlite", "neo.db"},
		{":memory:", "sqlite", ":memory:"},
		{"neo.db", "sqlite", "neo.db"},
		{"sqlite:///data/neo.db", "sqlite", "/data/neo.db"},
	}
	for _, tt := range tests {
		driver, dsn := parseDriverDSN(tt.input)
		if driver != tt.wantDriver {
			t.Errorf("parseDriverDSN(%q) driver = %q, want %q", tt.input, driver, tt.wantDriver)
		}
		if dsn != tt.wantDSN {
			t.Errorf("parseDriverDSN(%q) dsn = %q, want %q", tt.input, dsn, tt.wantDSN)
		}
	}
}

func TestParseDriverDSN_Postgres(t *testing.T) {
	tests := []struct {
		input      string
		wantDriver string
		wantDSN    string
	}{
		{"postgres://user:pass@host:5432/neo", "pgx", "postgres://user:pass@host:5432/neo"},
		{"postgresql://user:pass@host/neo", "pgx", "postgresql://user:pass@host/neo"},
	}
	for _, tt := range tests {
		driver, dsn := parseDriverDSN(tt.input)
		if driver != tt.wantDriver {
			t.Errorf("parseDriverDSN(%q) driver = %q, want %q", tt.input, driver, tt.wantDriver)
		}
		if dsn != tt.wantDSN {
			t.Errorf("parseDriverDSN(%q) dsn = %q, want %q", tt.input, dsn, tt.wantDSN)
		}
	}
}

func TestNew_SetsCurrentDialect(t *testing.T) {
	_, err := New(":memory:")
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	if CurrentDialect != DialectSQLite {
		t.Errorf("CurrentDialect = %q, want %q", CurrentDialect, DialectSQLite)
	}
}
