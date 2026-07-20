package catalog

import (
	"fmt"
)

// AppEnv represents the application environment.
type AppEnv string

const (
	// AppEnvDevelopment is the local development environment.
	AppEnvDevelopment AppEnv = "development"
	// AppEnvProduction is the production environment.
	AppEnvProduction AppEnv = "production"
)

// ParseAppEnv validates and returns the AppEnv from a string.
// It returns an error if the value is not "development" or "production".
func ParseAppEnv(s string) (AppEnv, error) {
	switch s {
	case "development":
		return AppEnvDevelopment, nil
	case "production":
		return AppEnvProduction, nil
	default:
		return "", fmt.Errorf("invalid APP_ENV: %q (must be 'development' or 'production')", s)
	}
}

// ResolveDBPath returns the resolved SQLite database path for the given environment.
// The logical defaults are data/dev.sqlite3 and data/production.sqlite3 (documented
// in ADR-0001 and the runbook). Because the documented workflow runs from src/api/,
// the physical defaults use ../data/... to land in the repository-level data/ directory.
// A non-empty path is accepted as-is in both environments (override via SQLITE_DB_PATH).
func ResolveDBPath(env AppEnv, path string) (string, error) {
	if path != "" {
		return path, nil
	}
	switch env {
	case AppEnvDevelopment:
		return "../data/dev.sqlite3", nil
	case AppEnvProduction:
		return "../data/production.sqlite3", nil
	default:
		return "", fmt.Errorf("SQLITE_DB_PATH is required in %s", env)
	}
}
