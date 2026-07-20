// Command migrate runs embedded goose migrations against the configured SQLite database.
// It applies runtime PRAGMAs and migrates to the latest version.
//
// Usage:
//
//	APP_ENV=development [SQLITE_DB_PATH=../data/dev.sqlite3] go run ./cmd/migrate
//	APP_ENV=production SQLITE_DB_PATH=/data/prod.sqlite3 go run ./cmd/migrate
package main

import (
	"log"
	"os"

	"ecommerce/api/catalog"

	_ "modernc.org/sqlite"
)

func main() {
	env, err := catalog.ParseAppEnv(os.Getenv("APP_ENV"))
	if err != nil {
		log.Fatalf("APP_ENV: %v", err)
	}

	dbPath, err := catalog.ResolveDBPath(env, os.Getenv("SQLITE_DB_PATH"))
	if err != nil {
		log.Fatalf("SQLITE_DB_PATH: %v", err)
	}

	log.Printf("Migrating database at %s (APP_ENV=%s)", dbPath, env)

	db, err := catalog.OpenSQLite(dbPath)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}
	db.Close()

	log.Println("Migration complete")
}
