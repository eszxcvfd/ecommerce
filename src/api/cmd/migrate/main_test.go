package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMigrate_Success(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "migrate")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Skipf("build failed: %v (integration test requires build)", err)
	}

	dbPath := filepath.Join(t.TempDir(), "migrated.db")

	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(),
		"APP_ENV=development",
		"SQLITE_DB_PATH="+dbPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("migrate failed: %v\noutput: %s", err, out)
	}

	if _, err := os.Stat(dbPath); err != nil {
		t.Errorf("database file not created: %v", err)
	}
}

func TestMigrate_FailsWithoutAppEnv(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "migrate_noenv")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Skipf("build failed: %v", err)
	}

	cmd := exec.Command(bin)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected migrate to fail without APP_ENV")
	}
	if len(out) == 0 {
		t.Fatal("expected error output for missing APP_ENV")
	}
}
