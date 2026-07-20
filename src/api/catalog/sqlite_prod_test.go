package catalog

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

// productionOpen opens a SQLite database with runtime PRAGMAs but without auto-migration.
// It verifies schema readiness instead. This simulates the production path.
func productionOpen(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	// Apply PRAGMAs (same as OpenSQLite but without migration)
	pragmas := []struct {
		stmt string
		name string
	}{
		{"PRAGMA foreign_keys = ON", "foreign_keys"},
		{"PRAGMA journal_mode = WAL", "journal_mode"},
		{"PRAGMA busy_timeout = 5000", "busy_timeout"},
		{"PRAGMA synchronous = NORMAL", "synchronous"},
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p.stmt); err != nil {
			db.Close()
			return nil, err
		}
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func TestVerifySchema_WithMigratedDB(t *testing.T) {
	// Open with dev path (auto-migrate), then verify schema
	devDB, err := OpenSQLite(filepath.Join(t.TempDir(), "migrated.db"))
	if err != nil {
		t.Fatalf("OpenSQLite failed: %v", err)
	}
	defer devDB.Close()

	if err := VerifySchema(devDB); err != nil {
		t.Errorf("VerifySchema on migrated DB should succeed: %v", err)
	}
}

func TestVerifySchema_WithoutMigrations_RawDB(t *testing.T) {
	// Open via productionOpen (no auto-migration, just PRAGMAs)
	db, err := productionOpen(filepath.Join(t.TempDir(), "raw.db"))
	if err != nil {
		t.Fatalf("productionOpen failed: %v", err)
	}
	defer db.Close()

	if err := VerifySchema(db); err == nil {
		t.Error("VerifySchema on unmigrated DB should fail")
	}
}

func TestVerifySchema_WithPartiallyAppliedMigrations(t *testing.T) {
	db, err := productionOpen(filepath.Join(t.TempDir(), "partial.db"))
	if err != nil {
		t.Fatalf("productionOpen failed: %v", err)
	}
	defer db.Close()

	// Create the goose tracking table but don't apply the migration
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS goose_db_version (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version_id INTEGER NOT NULL,
		is_applied INTEGER NOT NULL,
		tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		t.Fatal(err)
	}
	// Record version 0 as applied but NOT version 1 (the real migration)
	if _, err := db.Exec(
		"INSERT INTO goose_db_version (version_id, is_applied) VALUES (0, 1)",
	); err != nil {
		t.Fatal(err)
	}

	if err := VerifySchema(db); err == nil {
		t.Error("VerifySchema with pending migration should fail")
	}
}
