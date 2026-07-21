// Command test starts the API with seeded data for E2E tests.
// It uses a temporary SQLite database, runs migrations and seed data,
// then starts an HTTP server on the configured port.
//
// This is intentionally NOT production-grade: no signal handling, graceful
// shutdown, or backup/import policy. It blocks serving until interrupted.
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ecommerce/api/account"
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

	// Seed the catalog database
	if err := catalog.SeedSQLite(db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	// Run account migrations
	if err := account.RunMigrations(db); err != nil {
		log.Fatalf("account migrations: %v", err)
	}

	catalogRepo := catalog.NewSQLiteRepo(db)
	accountRepo := account.NewSQLiteRepo(db)

	// Build mux with all module routes
	mux := http.NewServeMux()
	catalog.RegisterRoutes(mux, catalogRepo)
	checkReady := func() error { return catalog.VerifySchema(db) }
	catalog.RegisterHealthRoutes(mux, checkReady)
	account.RegisterRoutes(mux, accountRepo)
	catalog.RegisterSellerRoutes(mux, catalogRepo, accountRepo)

	srv := catalog.NewServerWithHandler(":"+port, mux, db)

	log.Printf("Test API server listening on :%s (db=%s)", port, dbPath)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("server stopped: %v", err)
	}
}

func timestamp() string {
	return time.Now().Format("150405.000000000")
}
