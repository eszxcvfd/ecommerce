// Command testserver starts the catalog API with seeded data for e2e tests.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Test API server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down gracefully…")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("Server stopped")
}

func timestamp() string {
	return time.Now().Format("150405.000000000")
}
