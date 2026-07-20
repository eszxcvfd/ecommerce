package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTestJSON(t *testing.T, path string, data interface{}) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
}

func TestImportCatalog_Success(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "importcatalog")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Skipf("build failed: %v (integration test requires build)", err)
	}

	// First create a database with migrations.
	dbPath := filepath.Join(t.TempDir(), "imported.db")
	migBin := filepath.Join(t.TempDir(), "migrate_prep")
	if err := exec.Command("go", "build", "-o", migBin, "../migrate").Run(); err != nil {
		t.Skipf("migrate build failed: %v", err)
	}
	migCmd := exec.Command(migBin)
	migCmd.Env = append(os.Environ(),
		"APP_ENV=development",
		"SQLITE_DB_PATH="+dbPath,
	)
	if out, err := migCmd.CombinedOutput(); err != nil {
		t.Fatalf("prepare migration failed: %v\n%s", err, out)
	}

	// Create a valid catalog JSON.
	jsonPath := filepath.Join(t.TempDir(), "products.json")
	writeTestJSON(t, jsonPath, map[string]interface{}{
		"version": 1,
		"products": []map[string]interface{}{
			{
				"id": "sp-import-1", "ten": "Imported One",
				"mo_ta":         "First imported product",
				"anh_demo":      "/img.jpg",
				"gia":           map[string]interface{}{"mien_phi": true},
				"danh_muc":      "điện tử",
				"dinh_dang":     []string{"pdf"},
				"diem_danh_gia": 4.0, "so_luong_danh_gia": 1,
				"ngay_tao": "2026-07-15T00:00:00Z", "so_luot_tai": 10,
				"trang_thai": "approved",
			},
			{
				"id": "sp-import-2", "ten": "Imported Two",
				"mo_ta":         "Second imported product",
				"anh_demo":      "/img2.jpg",
				"gia":           map[string]interface{}{"mien_phi": false, "so_xu": 15000},
				"danh_muc":      "cơ khí",
				"dinh_dang":     []string{"dwg"},
				"diem_danh_gia": 0, "so_luong_danh_gia": 0,
				"ngay_tao": "2026-07-16T00:00:00Z", "so_luot_tai": 5,
				"trang_thai": "approved",
			},
		},
	})

	cmd := exec.Command(bin, "-path", jsonPath)
	cmd.Env = append(os.Environ(),
		"APP_ENV=development",
		"SQLITE_DB_PATH="+dbPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("import failed: %v\noutput: %s", err, out)
	}
}

func TestImportCatalog_RejectsDuplicateIDs(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "importcatalog_dup")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Skipf("build failed: %v", err)
	}

	dbPath := filepath.Join(t.TempDir(), "import_dup.db")
	migBin := filepath.Join(t.TempDir(), "migrate_dup_prep")
	if err := exec.Command("go", "build", "-o", migBin, "../migrate").Run(); err != nil {
		t.Skipf("migrate build failed: %v", err)
	}
	migCmd := exec.Command(migBin)
	migCmd.Env = append(os.Environ(),
		"APP_ENV=development",
		"SQLITE_DB_PATH="+dbPath,
	)
	if out, err := migCmd.CombinedOutput(); err != nil {
		t.Fatalf("prepare migration failed: %v\n%s", err, out)
	}

	// JSON with duplicate IDs.
	jsonPath := filepath.Join(t.TempDir(), "dup.json")
	writeTestJSON(t, jsonPath, map[string]interface{}{
		"version": 1,
		"products": []map[string]interface{}{
			{
				"id": "sp-dup", "ten": "First",
				"danh_muc": "kiến trúc", "ngay_tao": "2026-07-01T00:00:00Z",
				"gia":        map[string]interface{}{"mien_phi": true},
				"trang_thai": "approved",
			},
			{
				"id": "sp-dup", "ten": "Second",
				"danh_muc": "cơ khí", "ngay_tao": "2026-07-01T00:00:00Z",
				"gia":        map[string]interface{}{"mien_phi": true},
				"trang_thai": "approved",
			},
		},
	})

	cmd := exec.Command(bin, "-path", jsonPath)
	cmd.Env = append(os.Environ(),
		"APP_ENV=development",
		"SQLITE_DB_PATH="+dbPath,
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected import to fail with duplicate IDs")
	}
	if len(out) == 0 {
		t.Fatal("expected error output")
	}
}

func TestImportCatalog_RejectsInvalidVersion(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "importcatalog_ver")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Skipf("build failed: %v", err)
	}

	jsonPath := filepath.Join(t.TempDir(), "badver.json")
	writeTestJSON(t, jsonPath, map[string]interface{}{
		"version":  99,
		"products": []map[string]interface{}{},
	})

	dbPath := filepath.Join(t.TempDir(), "import_ver.db")
	cmd := exec.Command(bin, "-path", jsonPath)
	cmd.Env = append(os.Environ(),
		"APP_ENV=development",
		"SQLITE_DB_PATH="+dbPath,
	)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected import to fail with invalid version")
	}
	if len(out) == 0 {
		t.Fatal("expected error output")
	}
}
