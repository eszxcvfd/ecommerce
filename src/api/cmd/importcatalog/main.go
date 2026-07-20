// Command importcatalog validates and imports a versioned catalog JSON file
// into the configured SQLite database. It validates the entire input before
// beginning the write transaction and imports all-or-nothing.
//
// Usage:
//
//	APP_ENV=production SQLITE_DB_PATH=/data/prod.sqlite3 \
//	  go run ./cmd/importcatalog -path seed_data.json
//
//	Allow duplicate IDs (INSERT OR IGNORE):
//	  go run ./cmd/importcatalog -path data.json -allow-duplicates
package main

import (
	"flag"
	"log"
	"os"

	"ecommerce/api/catalog"

	_ "modernc.org/sqlite"
)

func main() {
	pathFlag := flag.String("path", "", "path to versioned catalog JSON file")
	allowDups := flag.Bool("allow-duplicates", false, "allow existing IDs to be skipped instead of rejected")
	flag.Parse()

	if *pathFlag == "" {
		log.Fatal("-path is required")
	}

	data, err := os.ReadFile(*pathFlag)
	if err != nil {
		log.Fatalf("read %s: %v", *pathFlag, err)
	}

	env, err := catalog.ParseAppEnv(os.Getenv("APP_ENV"))
	if err != nil {
		log.Fatalf("APP_ENV: %v", err)
	}

	dbPath, err := catalog.ResolveDBPath(env, os.Getenv("SQLITE_DB_PATH"))
	if err != nil {
		log.Fatalf("SQLITE_DB_PATH: %v", err)
	}

	db, err := catalog.OpenSQLiteProd(dbPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	log.Printf("Importing from %s into %s", *pathFlag, dbPath)

	if err := catalog.ImportCatalogJSON(db, data, *allowDups); err != nil {
		log.Fatalf("import: %v", err)
	}

	log.Println("Import complete")
}
