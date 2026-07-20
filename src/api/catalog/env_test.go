package catalog

import (
	"testing"
)

func TestParseAppEnv(t *testing.T) {
	t.Run("development returns AppEnvDevelopment", func(t *testing.T) {
		env, err := ParseAppEnv("development")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if env != AppEnvDevelopment {
			t.Errorf("expected development, got %q", env)
		}
	})

	t.Run("production returns AppEnvProduction", func(t *testing.T) {
		env, err := ParseAppEnv("production")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if env != AppEnvProduction {
			t.Errorf("expected production, got %q", env)
		}
	})

	t.Run("empty string returns error", func(t *testing.T) {
		_, err := ParseAppEnv("")
		if err == nil {
			t.Fatal("expected error for empty APP_ENV")
		}
	})

	t.Run("invalid value returns error", func(t *testing.T) {
		_, err := ParseAppEnv("staging")
		if err == nil {
			t.Fatal("expected error for invalid APP_ENV")
		}
	})
}

func TestResolveDBPath(t *testing.T) {
	t.Run("development with explicit path returns that path", func(t *testing.T) {
		path, err := ResolveDBPath(AppEnvDevelopment, "/tmp/mydb.sqlite3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path != "/tmp/mydb.sqlite3" {
			t.Errorf("expected /tmp/mydb.sqlite3, got %q", path)
		}
	})

	t.Run("development with empty path returns ../data/dev.sqlite3", func(t *testing.T) {
		path, err := ResolveDBPath(AppEnvDevelopment, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path != "../data/dev.sqlite3" {
			t.Errorf("expected default '../data/dev.sqlite3', got %q", path)
		}
	})

	t.Run("production with explicit absolute path returns that path", func(t *testing.T) {
		path, err := ResolveDBPath(AppEnvProduction, "/var/lib/ecommerce/prod.sqlite3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path != "/var/lib/ecommerce/prod.sqlite3" {
			t.Errorf("expected /var/lib/ecommerce/prod.sqlite3, got %q", path)
		}
	})

	t.Run("production with empty path returns ../data/production.sqlite3", func(t *testing.T) {
		path, err := ResolveDBPath(AppEnvProduction, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path != "../data/production.sqlite3" {
			t.Errorf("expected default '../data/production.sqlite3', got %q", path)
		}
	})

	t.Run("production with relative path is accepted as override", func(t *testing.T) {
		path, err := ResolveDBPath(AppEnvProduction, "my/path.sqlite3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path != "my/path.sqlite3" {
			t.Errorf("expected 'my/path.sqlite3', got %q", path)
		}
	})
}
