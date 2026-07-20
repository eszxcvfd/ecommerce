package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestProductionEntrypointServesDeterministicSeed(t *testing.T) {
	// Build the production binary — same as `go run .`
	bin := "/tmp/ecom-api-test"
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Skipf("build failed: %v (integration test requires build)", err)
	}
	defer os.Remove(bin)

	// Start on a free-ish port with required env vars
	cmd := exec.Command(bin)
	dbPath := t.TempDir() + "/test.sqlite3"
	cmd.Env = append(os.Environ(),
		"API_PORT=19999",
		"APP_ENV=development",
		"SQLITE_DB_PATH="+dbPath,
	)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()

	// Give the server time to bind
	time.Sleep(600 * time.Millisecond)

	resp, err := http.Get("http://localhost:19999/api/v1/san-pham")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var body struct {
		SanPham []struct{ ID string } `json:"san_pham"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("valid JSON response expected: %v", err)
	}

	if len(body.SanPham) != 12 {
		t.Fatalf("expected 12 products from development entrypoint (auto-seeded via SeedSQLite), got %d", len(body.SanPham))
	}
}
