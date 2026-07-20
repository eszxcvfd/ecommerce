package catalog

import (
	"fmt"
	"path/filepath"
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
// For development, if path is empty, it returns the default "../var/dev.sqlite3".
// For production, path must be non-empty and absolute.
func ResolveDBPath(env AppEnv, path string) (string, error) {
	if path != "" {
		if env == AppEnvProduction && !filepath.IsAbs(path) {
			return "", fmt.Errorf("SQLITE_DB_PATH must be an absolute path in production")
		}
		return path, nil
	}
	if env == AppEnvDevelopment {
		return "../var/dev.sqlite3", nil
	}
	return "", fmt.Errorf("SQLITE_DB_PATH is required in production")
}
