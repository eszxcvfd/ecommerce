// Command test starts the catalog API with seeded data for e2e tests.
// This test verifies the test binary starts correctly, uses a unique
// temporary SQLite database (not data/dev.sqlite3), and serves the expected
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

func TestTest_UsesTempSQLiteAndServesCatalog(t *testing.T) {
	// Build the test binary.
	bin := filepath.Join(t.TempDir(), "e2e-test")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Skipf("build failed: %v (integration test requires build)", err)
	}

	// Start on a free-ish port.
	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), "API_PORT=19998")
	// Intentionally do NOT set APP_ENV or SQLITE_DB_PATH — the test command
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
		t.Fatalf("test did not respond: %v", err)
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

func TestTest_DBPathIsInTempDir_NotDataDir(t *testing.T) {
	// Build the test binary.
	bin := filepath.Join(t.TempDir(), "e2e-test-dbpath")
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

	// 1. API works.
	resp, err := http.Get("http://localhost:19997/api/v1/san-pham")
	if err != nil {
		t.Fatalf("test did not respond: %v", err)
	}
	resp.Body.Close()

	// 2. No data/dev.sqlite3 was created by the test process.
	badPaths := []string{
		"data/dev.sqlite3",
		filepath.Join("..", "data", "dev.sqlite3"),
	}
	for _, p := range badPaths {
		abs, _ := filepath.Abs(p)
		if _, err := os.Stat(abs); err == nil {
			info, _ := os.Stat(abs)
			if time.Since(info.ModTime()) < 30*time.Second {
				data, readErr := os.ReadFile(abs)
				if readErr == nil && strings.Contains(string(data), "sp-001") {
					t.Errorf("test created/used database at %s instead of temp dir", abs)
				}
			}
		}
	}
}
