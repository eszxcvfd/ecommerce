package catalog

import (
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestImportCatalogJSON_Success(t *testing.T) {
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "import_ok.db"))
	if err != nil {
		t.Fatalf("OpenSQLite: %v", err)
	}
	defer db.Close()

	// Use the embedded seed JSON (already valid).
	if err := ImportCatalogJSON(db, seedDataJSON, false); err != nil {
		t.Fatalf("ImportCatalogJSON failed: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM san_pham_so").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 16 {
		t.Fatalf("expected 16 products, got %d", count)
	}

	// Products must be queryable through the repository.
	repo := NewSQLiteRepo(db)
	products, err := repo.Products()
	if err != nil {
		t.Fatalf("Products() failed: %v", err)
	}
	if len(products) != 12 {
		t.Fatalf("expected 12 approved products, got %d", len(products))
	}
}

func TestImportCatalogJSON_RejectsDuplicateIDsInInput(t *testing.T) {
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "import_dup_input.db"))
	if err != nil {
		t.Fatalf("OpenSQLite: %v", err)
	}
	defer db.Close()

	data := []byte(`{
		"version": 1,
		"products": [
			{"id":"sp-001","ten":"A","danh_muc":"kiến trúc","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"},
			{"id":"sp-001","ten":"B","danh_muc":"cơ khí","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"}
		]
	}`)

	err = ImportCatalogJSON(db, data, false)
	if err == nil {
		t.Fatal("expected error for duplicate IDs in input")
	}

	// DB must be untouched (all-or-nothing).
	var count int
	db.QueryRow("SELECT COUNT(*) FROM san_pham_so").Scan(&count)
	if count != 0 {
		t.Fatalf("expected 0 products after failed import, got %d", count)
	}
}

func TestImportCatalogJSON_RejectsExistingIDs(t *testing.T) {
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "import_existing.db"))
	if err != nil {
		t.Fatalf("OpenSQLite: %v", err)
	}
	defer db.Close()

	// Pre-insert a product with ID sp-001.
	if _, err := db.Exec(`
		INSERT INTO san_pham_so (id, ten, mo_ta, anh_demo, mien_phi, so_xu, danh_muc,
		                         diem_danh_gia, so_luong_danh_gia, ngay_tao, so_luot_tai, trang_thai,
		                         ten_search, mo_ta_search)
		VALUES ('sp-001', 'Existing', '', '', 1, 0, 'kiến trúc', 0, 0, '2026-07-01T00:00:00Z', 0, 'approved', '', '')
	`); err != nil {
		t.Fatal(err)
	}

	data := []byte(`{
		"version": 1,
		"products": [
			{"id":"sp-001","ten":"New","danh_muc":"cơ khí","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"},
			{"id":"sp-002","ten":"New Two","danh_muc":"điện tử","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"}
		]
	}`)

	err = ImportCatalogJSON(db, data, false)
	if err == nil {
		t.Fatal("expected error for existing ID in DB")
	}

	// DB must be untouched — sp-002 should NOT have been inserted.
	var count int
	db.QueryRow("SELECT COUNT(*) FROM san_pham_so").Scan(&count)
	if count != 1 {
		t.Fatalf("expected only 1 pre-existing product, got %d", count)
	}
}

func TestImportCatalogJSON_AllowDuplicatesSkipsConflicts(t *testing.T) {
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "import_allow.db"))
	if err != nil {
		t.Fatalf("OpenSQLite: %v", err)
	}
	defer db.Close()

	// Pre-insert a product.
	if _, err := db.Exec(`
		INSERT INTO san_pham_so (id, ten, mo_ta, anh_demo, mien_phi, so_xu, danh_muc,
		                         diem_danh_gia, so_luong_danh_gia, ngay_tao, so_luot_tai, trang_thai,
		                         ten_search, mo_ta_search)
		VALUES ('sp-001', 'Original', '', '', 1, 0, 'kiến trúc', 0, 0, '2026-07-01T00:00:00Z', 0, 'approved', '', '')
	`); err != nil {
		t.Fatal(err)
	}

	data := []byte(`{
		"version": 1,
		"products": [
			{"id":"sp-001","ten":"Duplicate","danh_muc":"cơ khí","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"},
			{"id":"sp-002","ten":"New","danh_muc":"điện tử","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"}
		]
	}`)

	if err := ImportCatalogJSON(db, data, true); err != nil {
		t.Fatalf("ImportCatalogJSON with allowDuplicates failed: %v", err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM san_pham_so").Scan(&count)
	if count != 2 {
		t.Fatalf("expected 2 products (original + new), got %d", count)
	}

	// Original must not have been overwritten.
	var ten string
	db.QueryRow("SELECT ten FROM san_pham_so WHERE id = 'sp-001'").Scan(&ten)
	if ten != "Original" {
		t.Errorf("expected sp-001 to keep 'Original', got %q", ten)
	}
}

func TestImportCatalogJSON_RejectsInvalidJSON(t *testing.T) {
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "import_invalid.db"))
	if err != nil {
		t.Fatalf("OpenSQLite: %v", err)
	}
	defer db.Close()

	err = ImportCatalogJSON(db, []byte(`{"version": 42, "products": []}`), false)
	if err == nil {
		t.Fatal("expected error for invalid version")
	}
}
