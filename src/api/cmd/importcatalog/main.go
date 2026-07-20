// Command importcatalog validates and imports a versioned catalog JSON file
// into the configured SQLite database. It delegates all orchestration to
// catalog.ImportFromFile, keeping this CLI thin.
//
// Usage:
//
//	APP_ENV=production SQLITE_DB_PATH=/data/prod.sqlite3 \
//	Allow duplicate IDs (INSERT OR IGNORE):
//	  go run ./cmd/importcatalog -path data.json -allow-duplicates
package main

import (
	"flag"
	"log"

	"ecommerce/api/catalog"
)

func main() {
	pathFlag := flag.String("path", "", "path to versioned catalog JSON file")
	allowDups := flag.Bool("allow-duplicates", false, "allow existing IDs to be skipped instead of rejected")
	flag.Parse()

	if *pathFlag == "" {
		log.Fatal("-path is required")
	}

	log.Printf("Importing catalog from %s", *pathFlag)

	if err := catalog.ImportFromFile(*pathFlag, *allowDups); err != nil {
		log.Fatalf("import: %v", err)
	}
}
