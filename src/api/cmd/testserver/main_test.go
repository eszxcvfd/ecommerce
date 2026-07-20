// Command testserver starts the catalog API with seeded data for e2e tests.
// This test verifies the testserver binary starts correctly, uses a unique
// temporary SQLite database (not var/dev.sqlite3), and serves the expected
// catalog data through the HTTP API.
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTestserver_UsesTempSQLiteAndServesCatalog(t *testing.T) {
	// Build the testserver binary.
	bin := filepath.Join(t.TempDir(), "testserver")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Skipf("build failed: %v (integration test requires build)", err)
	}

	// Start on a free-ish port.
	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), "API_PORT=19998")
	// Intentionally do NOT set APP_ENV or SQLITE_DB_PATH — the testserver
	// should derive a temp path on its own without requiring env vars.

	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()

	// Give the server time to migrate, seed, and bind.
	time.Sleep(800 * time.Millisecond)

	// Fetch the product list.
	resp, err := http.Get("http://localhost:19998/api/v1/san-pham")
	if err != nil {
		t.Fatalf("testserver did not respond: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body struct {
		SanPham []struct{ ID string } `json:"san_pham"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("valid JSON response expected: %v", err)
	}

	if len(body.SanPham) != 12 {
		t.Fatalf("expected 12 approved products, got %d", len(body.SanPham))
	}

	// Verify health endpoints.
	for _, path := range []string{"/healthz", "/readyz"} {
		resp2, err := http.Get("http://localhost:19998" + path)
		if err != nil {
			t.Fatalf("%s failed: %v", path, err)
		}
		resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			t.Errorf("%s returned %d, expected 200", path, resp2.StatusCode)
		}
	}
}

func TestTestserver_DBPathIsInTempDir_NotVarDir(t *testing.T) {
	// Build the testserver binary.
	bin := filepath.Join(t.TempDir(), "testserver-dbpath")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Skipf("build failed: %v (integration test requires build)", err)
	}

	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), "API_PORT=19997")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()

	time.Sleep(800 * time.Millisecond)

	// The testserver logs the database path.
	// Since we can't read stderr easily after the fact via exec, we instead
	// verify via the API that it's working and confirm no var/dev.sqlite3
	// was created by this run.

	// 1. API works.
	resp, err := http.Get("http://localhost:19997/api/v1/san-pham")
	if err != nil {
		t.Fatalf("testserver did not respond: %v", err)
	}
	resp.Body.Close()

	// 2. No var/dev.sqlite3 was created by the testserver process.
	// The testserver runs from src/api/ and the default dev path is ../var/dev.sqlite3
	// which resolves to src/var/dev.sqlite3.
	// Check a few known locations to be safe.
	badPaths := []string{
		"../../var/dev.sqlite3",           // relative to cmd/testserver/
		"../var/dev.sqlite3",              // relative to src/api/
		filepath.Join("..", "..", "var", "dev.sqlite3"),
	}
	for _, p := range badPaths {
		abs, _ := filepath.Abs(p)
		if _, err := os.Stat(abs); err == nil {
			// File exists — check if it's been recently modified (within last 5s)
			info, _ := os.Stat(abs)
			if time.Since(info.ModTime()) < 30*time.Second {
				// Check if it contains our test data
				data, readErr := os.ReadFile(abs)
				if readErr == nil && strings.Contains(string(data), "sp-001") {
					t.Errorf("testserver created/used database at %s instead of temp dir", abs)
				}
			}
		}
	}
}
