package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// Run account module migrations
	if err := account.RunMigrations(db); err != nil {
		log.Fatalf("account migrations: %v", err)
	}

	catalogRepo := catalog.NewSQLiteRepo(db)
	accountRepo := account.NewSQLiteRepo(db)

	// Build mux with all module routes
	mux := http.NewServeMux()
	catalog.RegisterRoutes(mux, catalogRepo)
	catalog.RegisterHealthRoutes(mux, checkReady)
	catalog.RegisterSellerRoutes(mux, catalogRepo, accountRepo)

	// Seed admin account for development
	if env == catalog.AppEnvDevelopment {
		if err := account.SeedAdmin(context.Background(), accountRepo); err != nil {
			log.Printf("seed admin (non-fatal): %v", err)
		}
	}

	// Create server with the pre-built mux
	srv := catalog.NewServerWithHandler(":"+port, mux, db)

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
