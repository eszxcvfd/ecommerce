// Package catalog provides the import service for the catalog domain.
// ImportFromFile is the reusable non-CLI entrypoint for importing a
// versioned catalog JSON file into the configured SQLite database.
package catalog

import (
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// ImportFromFile reads a versioned catalog JSON file, resolves the database
// path from APP_ENV / SQLITE_DB_PATH (same convention as main.go), opens
// the database in production mode, and imports the catalog data.
//
// This is the reusable import service used by cmd/importcatalog. It handles
// file reading, environment configuration, database opening, and the import
// transaction — keeping the CLI thin.
func ImportFromFile(path string, allowDuplicates bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	env, err := ParseAppEnv(os.Getenv("APP_ENV"))
	if err != nil {
		return err
	}

	dbPath, err := ResolveDBPath(env, os.Getenv("SQLITE_DB_PATH"))
	if err != nil {
		return err
	}

	db, err := OpenSQLiteProd(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Printf("Importing from %s into %s", path, dbPath)

	if err := ImportCatalogJSON(db, data, allowDuplicates); err != nil {
		return err
	}

	log.Println("Import complete")
	return nil
}
