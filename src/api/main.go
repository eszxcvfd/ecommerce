package main

import (
	"context"
	"database/sql"
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
		if err := catalog.SeedSQLite(db); err != nil {
			log.Fatalf("seed: %v", err)
		}
		checkReady = func() error { return catalog.VerifySchema(db) }
		log.Printf("Development mode — database at %s", dbPath)
	} else {
		db, err = catalog.OpenSQLiteProd(dbPath)
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
