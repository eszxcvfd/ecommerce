// Command test starts the catalog API with seeded data for E2E tests.
// It uses a temporary SQLite database, runs migrations and seed data,
// then starts an HTTP server on the configured port.
//
// This is intentionally NOT production-grade: no signal handling, graceful
// shutdown, or backup/import policy. It blocks serving until interrupted.
package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"ecommerce/api/catalog"

	_ "modernc.org/sqlite"
)

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	// Use a temporary SQLite database for each test run
	dbPath := filepath.Join(os.TempDir(), "ecommerce-test-"+timestamp()+".sqlite3")
	db, err := catalog.OpenSQLite(dbPath)
	if err != nil {
		log.Fatalf("open SQLite: %v", err)
	}

	// Seed the database
	if err := catalog.SeedSQLite(db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	repo := catalog.NewSQLiteRepo(db)
	checkReady := func() error { return catalog.VerifySchema(db) }
	srv := catalog.NewServer(":"+port, repo, db, checkReady)

	log.Printf("Test API server listening on :%s (db=%s)", port, dbPath)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("server stopped: %v", err)
	}
}

func timestamp() string {
	return time.Now().Format("150405.000000000")
}
