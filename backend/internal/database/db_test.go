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
	if err := migrate(db); err != nil {
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
	if count != len(migrations) {
		t.Errorf("expected %d migration(s) recorded, got %d", len(migrations), count)
	}
}
