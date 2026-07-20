package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
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

	env, err := catalog.ParseAppEnv(os.Getenv("APP_ENV"))
	if err != nil {
		log.Fatalf("APP_ENV: %v", err)
	}

	dbPath, err := catalog.ResolveDBPath(env, os.Getenv("SQLITE_DB_PATH"))
	if err != nil {
		log.Fatalf("SQLITE_DB_PATH: %v", err)
	}

	var db *sql.DB
	var checkReady func() error

	if env == catalog.AppEnvDevelopment {
		db, err = catalog.OpenSQLite(dbPath)
		if err != nil {
			log.Fatalf("open SQLite: %v", err)
		}
		checkReady = func() error { return catalog.VerifySchema(db) }
		log.Printf("Development mode — database at %s", dbPath)
	} else {
		db, err = openProdDB(dbPath)
		if err != nil {
			log.Fatalf("open SQLite: %v", err)
		}
		checkReady = func() error { return catalog.VerifySchema(db) }
		log.Printf("Production mode — database at %s", dbPath)
	}

	repo := catalog.NewSQLiteRepo(db)
	srv := catalog.NewServer(":"+port, repo, db, checkReady)

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("API server listening on :%s", port)
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

// openProdDB opens a SQLite database for production: applies runtime PRAGMAs
// and verifies schema, but does NOT auto-migrate.
func openProdDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	pragmas := []struct {
		stmt string
		name string
	}{
		{"PRAGMA foreign_keys = ON", "foreign_keys"},
		{"PRAGMA journal_mode = WAL", "journal_mode"},
		{fmt.Sprintf("PRAGMA busy_timeout = %d", 5000), "busy_timeout"},
		{"PRAGMA synchronous = NORMAL", "synchronous"},
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p.stmt); err != nil {
			db.Close()
			return nil, fmt.Errorf("set pragma %s: %w", p.name, err)
		}
	}
	db.SetMaxOpenConns(1)

	if err := catalog.VerifySchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("verify schema: %w", err)
	}

	return db, nil
}
