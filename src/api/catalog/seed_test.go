package catalog

import (
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSeedFromJSON_ProducesExpectedCount(t *testing.T) {
	products, err := SeedFromJSON()
	if err != nil {
		t.Fatalf("SeedFromJSON() failed: %v", err)
	}
	if len(products) != 16 {
		t.Fatalf("expected 16 products from seed JSON, got %d", len(products))
	}
}

func TestSeedFromJSON_ProductsContainFullData(t *testing.T) {
	products, err := SeedFromJSON()
	if err != nil {
		t.Fatalf("SeedFromJSON() failed: %v", err)
	}

	// Check sp-001 (free, approved, with formats)
	var sp001 *SanPhamSo
	for i := range products {
		if products[i].ID == "sp-001" {
			sp001 = &products[i]
			break
		}
	}
	if sp001 == nil {
		t.Fatal("expected sp-001 in seed")
	}
	if sp001.Ten != "Bản vẽ nhà phố 3 tầng" {
		t.Errorf("sp-001 Ten = %q", sp001.Ten)
	}
	if !sp001.Gia.MienPhi {
		t.Errorf("sp-001 should be free")
	}
	if len(sp001.DinhDang) != 2 {
		t.Errorf("sp-001 expected 2 dinh_dang, got %d", len(sp001.DinhDang))
	}
	if sp001.NgayTao.IsZero() {
		t.Errorf("sp-001 should have non-zero NgayTao")
	}

	// Check sp-007 (non-approved, from JSON)
	var sp007 *SanPhamSo
	for i := range products {
		if products[i].ID == "sp-007" {
			sp007 = &products[i]
			break
		}
	}
	if sp007 == nil {
		t.Fatal("expected sp-007 in seed")
	}
	if sp007.TrangThai != TrangThaiDangSoan {
		t.Errorf("sp-007 should be draft, got %s", sp007.TrangThai)
	}
	if sp007.NgayTao.IsZero() {
		t.Errorf("sp-007 should have non-zero NgayTao from JSON")
	}
}

func TestSeedSQLite_SeedsOnlyWhenEmpty(t *testing.T) {
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "seed_test.db"))
	if err != nil {
		t.Fatalf("OpenSQLite failed: %v", err)
	}
	defer db.Close()

	// First seed — DB is empty, should seed.
	if err := SeedSQLite(db); err != nil {
		t.Fatalf("first SeedSQLite failed: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM san_pham_so").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 16 {
		t.Fatalf("expected 16 products after first seed, got %d", count)
	}

	// Second seed — DB has products, should skip.
	if err := SeedSQLite(db); err != nil {
		t.Fatalf("second SeedSQLite failed: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM san_pham_so").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 16 {
		t.Fatalf("expected still 16 products after second seed (idempotent), got %d", count)
	}
}

func TestSeedSQLite_ProductsAreQueryable(t *testing.T) {
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "query_test.db"))
	if err != nil {
		t.Fatalf("OpenSQLite failed: %v", err)
	}
	defer db.Close()

	if err := SeedSQLite(db); err != nil {
		t.Fatalf("SeedSQLite failed: %v", err)
	}

	repo := NewSQLiteRepo(db)
	products, err := repo.Products()
	if err != nil {
		t.Fatalf("Products() failed: %v", err)
	}
	if len(products) != 12 {
		t.Fatalf("expected 12 approved products, got %d", len(products))
	}
}

func TestSeedSQLite_NoTransactionsOnSkip(t *testing.T) {
	// When DB has products, SeedSQLite must not start a write transaction.
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "noskip_test.db"))
	if err != nil {
		t.Fatalf("OpenSQLite failed: %v", err)
	}
	defer db.Close()

	// Pre-populate with one product
	if _, err := db.Exec(`
		INSERT INTO san_pham_so (id, ten, mo_ta, anh_demo, mien_phi, so_xu, danh_muc,
		                         diem_danh_gia, so_luong_danh_gia, ngay_tao, so_luot_tai, trang_thai,
		                         ten_search, mo_ta_search)
		VALUES ('sp-manual', 'Manual', '', '', 1, 0, 'kiến trúc', 0, 0, '2026-07-01T00:00:00Z', 0, 'approved', '', '')
	`); err != nil {
		t.Fatal(err)
	}

	// SeedSQLite should skip and not interfere
	if err := SeedSQLite(db); err != nil {
		t.Fatalf("SeedSQLite with existing data failed: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM san_pham_so").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 product after skip, got %d", count)
	}
}
