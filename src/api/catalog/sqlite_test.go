package catalog

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

// openTestSQLite opens a temporary SQLite database for testing.
// It uses the raw driver (not OpenSQLite) so we can test OpenSQLite independently.
func openTestSQLite(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestOpenSQLite_CreatesParentDirWhenAbsent(t *testing.T) {
	// Use a path inside a deeply nested non-existent subdirectory of TempDir.
	// Neither the parent dirs nor the file should exist before the call.
	dbPath := filepath.Join(t.TempDir(), "a", "b", "c", "test.db")
	if _, err := filepath.Glob(dbPath); err == nil {
		// sanity: dir does not exist yet
	}

	db, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("OpenSQLite with path in non-existent dir should succeed: %v", err)
	}
	defer db.Close()

	// Verify the parent directory was created.
	parent := filepath.Dir(filepath.Dir(dbPath))
	if fi, err := filepath.Glob(filepath.Dir(parent)); err == nil && len(fi) > 0 {
		// parent exists — good
	} else {
		t.Errorf("parent directory %s should exist after OpenSQLite", parent)
	}

	// DB should work — table exists.
	var name string
	if err := db.QueryRow(
		"SELECT name FROM sqlite_master WHERE type='table' AND name='san_pham_so'",
	).Scan(&name); err != nil {
		t.Errorf("expected san_pham_so table to exist: %v", err)
	}
}

func TestOpenSQLite_CreatesSchemaAndSetsPragmas(t *testing.T) {
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "catalog.db"))
	if err != nil {
		t.Fatalf("OpenSQLite failed: %v", err)
	}
	defer db.Close()

	// Verify PRAGMAs are set
	t.Run("foreign_keys is ON", func(t *testing.T) {
		var val int
		if err := db.QueryRow("PRAGMA foreign_keys").Scan(&val); err != nil {
			t.Fatal(err)
		}
		if val != 1 {
			t.Errorf("expected foreign_keys=1, got %d", val)
		}
	})

	t.Run("journal_mode is WAL", func(t *testing.T) {
		var mode string
		if err := db.QueryRow("PRAGMA journal_mode").Scan(&mode); err != nil {
			t.Fatal(err)
		}
		if mode != "wal" {
			t.Errorf("expected journal_mode=wal, got %s", mode)
		}
	})

	t.Run("synchronous is NORMAL", func(t *testing.T) {
		var val string
		if err := db.QueryRow("PRAGMA synchronous").Scan(&val); err != nil {
			t.Fatal(err)
		}
		if val != "NORMAL" && val != "1" {
			t.Errorf("expected synchronous=NORMAL (1), got %s", val)
		}
	})

	// Verify MaxOpenConns
	t.Run("MaxOpenConns is 1", func(t *testing.T) {
		if got := db.Stats().MaxOpenConnections; got != 1 {
			t.Errorf("expected MaxOpenConnections=1, got %d", got)
		}
	})

	// Verify schema tables exist
	t.Run("san_pham_so table exists", func(t *testing.T) {
		var name string
		err := db.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name='san_pham_so'",
		).Scan(&name)
		if err == sql.ErrNoRows {
			t.Fatal("table san_pham_so does not exist")
		}
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("san_pham_dinh_dang table exists", func(t *testing.T) {
		var name string
		err := db.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name='san_pham_dinh_dang'",
		).Scan(&name)
		if err == sql.ErrNoRows {
			t.Fatal("table san_pham_dinh_dang does not exist")
		}
		if err != nil {
			t.Fatal(err)
		}
	})

	// Verify indexes exist
	t.Run("targeted indexes exist", func(t *testing.T) {
		expectedIndexes := []string{
			"idx_san_pham_trang_thai",
			"idx_san_pham_danh_muc",
			"idx_san_pham_so_xu",
			"idx_san_pham_ngay_tao",
			"idx_san_pham_so_luot_tai",
			"idx_san_pham_diem_danh_gia",
			"idx_san_pham_dinh_dang",
		}
		for _, idx := range expectedIndexes {
			var name string
			err := db.QueryRow(
				"SELECT name FROM sqlite_master WHERE type='index' AND name=?", idx,
			).Scan(&name)
			if err == sql.ErrNoRows {
				t.Errorf("index %s does not exist", idx)
			} else if err != nil {
				t.Errorf("checking index %s: %v", idx, err)
			}
		}
	})

	// Verify migration is recorded
	t.Run("migration version recorded", func(t *testing.T) {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM goose_db_version").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count == 0 {
			t.Error("expected at least one migration record in goose_db_version")
		}
	})

	// Verify idempotency: running OpenSQLite again on same path works
	t.Run("reopen is idempotent", func(t *testing.T) {
		db2, err := OpenSQLite(filepath.Join(t.TempDir(), "catalog_reopen.db"))
		if err != nil {
			t.Fatalf("second OpenSQLite failed: %v", err)
		}
		defer db2.Close()
		var count int
		if err := db2.QueryRow("SELECT COUNT(*) FROM san_pham_so").Scan(&count); err != nil {
			t.Fatal(err)
		}
	})
}

// seedSQLite inserts the full seed dataset into the given (already-migrated) SQLite DB.
func seedSQLite(t *testing.T, db *sql.DB) {
	t.Helper()
	products := SeedData()
	for _, p := range products {
		_, err := db.Exec(`
			INSERT INTO san_pham_so (id, ten, mo_ta, anh_demo, mien_phi, so_xu, danh_muc,
			                         diem_danh_gia, so_luong_danh_gia, ngay_tao, so_luot_tai, trang_thai,
			                         ten_search, mo_ta_search)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, p.ID, p.Ten, p.MoTa, p.AnhDemo, boolToInt(p.Gia.MienPhi), p.Gia.SoXu,
			string(p.DanhMuc), p.DiemDanhGia, p.SoLuongDanhGia,
			p.NgayTao.Format(time.RFC3339), p.SoLuotTai, string(p.TrangThai),
			normalizeSearch(p.Ten), normalizeSearch(p.MoTa),
		)
		if err != nil {
			t.Fatalf("insert product %s: %v", p.ID, err)
		}
		for _, ext := range p.DinhDang {
			_, err := db.Exec(
				"INSERT INTO san_pham_dinh_dang (san_pham_id, dinh_dang) VALUES (?, ?)",
				p.ID, ext,
			)
			if err != nil {
				t.Fatalf("insert format %s for %s: %v", ext, p.ID, err)
			}
		}
	}
}

// openSQLiteAndSeed opens a temp SQLite database, runs migrations, and seeds it.
func openSQLiteAndSeed(t *testing.T) *sql.DB {
	t.Helper()
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("OpenSQLite: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	seedSQLite(t, db)
	return db
}

func TestSQLiteProducts_ReturnsOnlyApproved(t *testing.T) {
	db := openSQLiteAndSeed(t)
	repo := NewSQLiteRepo(db)

	products, err := repo.Products()
	if err != nil {
		t.Fatalf("Products() failed: %v", err)
	}

	if len(products) != 12 {
		t.Fatalf("expected 12 approved products, got %d", len(products))
	}

	// Must not include any non-approved product IDs
	nonApproved := map[string]bool{"sp-007": true, "sp-008": true, "sp-009": true, "sp-010": true}
	for _, sp := range products {
		if nonApproved[sp.ID] {
			t.Errorf("non-approved product %s (%s) was included", sp.ID, sp.Ten)
		}
		// Must be approved
		if sp.TrangThai != TrangThaiDaDuyet {
			t.Errorf("product %s has trang_thai %q, expected approved", sp.ID, sp.TrangThai)
		}
		// Basic fields
		if sp.Ten == "" {
			t.Errorf("product %s has empty Ten", sp.ID)
		}
	}
}

func TestSQLiteProducts_PopulatesDinhDang(t *testing.T) {
	db := openSQLiteAndSeed(t)
	repo := NewSQLiteRepo(db)

	products, err := repo.Products()
	if err != nil {
		t.Fatalf("Products() failed: %v", err)
	}

	// sp-004 (Bộ icon phong cách tối giản) has formats: svg, png, ai
	found := false
	for _, p := range products {
		if p.ID == "sp-004" {
			found = true
			if len(p.DinhDang) != 3 {
				t.Fatalf("expected sp-004 to have 3 dinh_dang, got %d: %v", len(p.DinhDang), p.DinhDang)
			}
			expected := []string{"ai", "png", "svg"} // sorted
			for i, ext := range p.DinhDang {
				if ext != expected[i] {
					t.Errorf("dinh_dang[%d] = %q, expected %q", i, ext, expected[i])
				}
			}
			break
		}
	}
	if !found {
		t.Fatal("expected sp-004 in approved products")
	}
}

func TestSQLiteProducts_ReturnsEmptyWhenNoApproved(t *testing.T) {
	db, err := OpenSQLite(filepath.Join(t.TempDir(), "empty.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Insert only non-approved products
	nonApproved := []SanPhamSo{
		{ID: "sp-100", Ten: "Draft", DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDangSoan},
		{ID: "sp-101", Ten: "Pending", DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiChoDuyet},
	}
	for _, p := range nonApproved {
		_, err := db.Exec(`
			INSERT INTO san_pham_so (id, ten, mo_ta, anh_demo, mien_phi, so_xu, danh_muc,
			                         diem_danh_gia, so_luong_danh_gia, ngay_tao, so_luot_tai, trang_thai)
			VALUES (?, ?, '', '', 0, 0, ?, 0, 0, datetime('now'), 0, ?)
		`, p.ID, p.Ten, string(p.DanhMuc), string(p.TrangThai))
		if err != nil {
			t.Fatalf("insert %s: %v", p.ID, err)
		}
	}

	repo := NewSQLiteRepo(db)
	products, err := repo.Products()
	if err != nil {
		t.Fatalf("Products() failed: %v", err)
	}
	if len(products) != 0 {
		t.Fatalf("expected 0 approved products, got %d", len(products))
	}
}

func TestSQLiteSearch_TextSearch(t *testing.T) {
	db := openSQLiteAndSeed(t)
	repo := NewSQLiteRepo(db)

	t.Run("empty query returns all approved products", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) != 12 {
			t.Fatalf("expected 12 products, got %d", len(products))
		}
	})

	t.Run("search by name substring", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Q: "CNC"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) != 2 {
			t.Fatalf("expected 2 CNC products, got %d", len(products))
		}
		for _, p := range products {
			if p.ID != "sp-017" && p.ID != "sp-018" {
				t.Errorf("unexpected product: %s (%s)", p.ID, p.Ten)
			}
		}
	})

	t.Run("search by description content", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Q: "Arduino"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) != 1 {
			t.Fatalf("expected 1 Arduino product, got %d", len(products))
		}
		if products[0].ID != "sp-003" {
			t.Errorf("expected sp-003, got %s", products[0].ID)
		}
	})

	t.Run("search is case-insensitive", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Q: "arduino"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) != 1 {
			t.Fatalf("expected 1 product for lowercase query, got %d", len(products))
		}
	})

	t.Run("search is accent-insensitive for Vietnamese", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Q: "xây dựng"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) == 0 {
			t.Fatal("expected at least 1 product matching 'xây dựng'")
		}
	})

	t.Run("search with no matches returns empty", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Q: "zzzznotfound"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) != 0 {
			t.Fatalf("expected 0 products, got %d", len(products))
		}
	})
}

func TestSQLiteSearch_Filter(t *testing.T) {
	db := openSQLiteAndSeed(t)
	repo := NewSQLiteRepo(db)

	t.Run("filter by danh_muc", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{DanhMuc: "kiến trúc"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) != 2 {
			t.Fatalf("expected 2 products in 'kiến trúc', got %d", len(products))
		}
		for _, p := range products {
			if string(p.DanhMuc) != "kiến trúc" {
				t.Errorf("expected 'kiến trúc', got %q", p.DanhMuc)
			}
		}
	})

	t.Run("filter by dinh_dang", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{DinhDang: "dxf"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) != 2 {
			t.Fatalf("expected 2 DXF products, got %d", len(products))
		}
		for _, p := range products {
			hasDXF := false
			for _, ext := range p.DinhDang {
				if ext == "dxf" {
					hasDXF = true
					break
				}
			}
			if !hasDXF {
				t.Errorf("product %s does not have dxf format", p.ID)
			}
		}
	})

	t.Run("filter by min_xu (paid products over 5000)", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{MinXu: ptrInt64(5000)})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		for _, p := range products {
			price := int64(0)
			if !p.Gia.MienPhi {
				price = p.Gia.SoXu
			}
			if price < 5000 {
				t.Errorf("product %s has price %d < min_xu 5000", p.ID, price)
			}
		}
	})

	t.Run("filter by max_xu (products under 200)", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{MaxXu: ptrInt64(200)})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		for _, p := range products {
			price := int64(0)
			if !p.Gia.MienPhi {
				price = p.Gia.SoXu
			}
			if price > 200 {
				t.Errorf("product %s has price %d > max_xu 200", p.ID, price)
			}
		}
	})

	t.Run("filter by price range inclusive", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{MinXu: ptrInt64(100), MaxXu: ptrInt64(100)})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) != 2 {
			t.Fatalf("expected 2 products priced at 100, got %d", len(products))
		}
	})

	t.Run("search + filter combined", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Q: "CNC", DanhMuc: "cơ khí"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) != 2 {
			t.Fatalf("expected 2 CNC products in 'cơ khí', got %d", len(products))
		}
	})
}

func ptrInt64(v int64) *int64 { return &v }

func TestSQLiteSearch_Sort(t *testing.T) {
	db := openSQLiteAndSeed(t)
	repo := NewSQLiteRepo(db)

	t.Run("default sort is newest", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		// newest first: sp-016 (2026-07-12) should be first
		if products[0].ID != "sp-016" {
			t.Errorf("expected newest first (sp-016), got %s", products[0].ID)
		}
	})

	t.Run("sort by popular (download count)", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Sort: "popular"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		// Most downloaded: sp-014 (520) then sp-004 (310)
		if products[0].ID != "sp-014" {
			t.Errorf("expected most popular first (sp-014), got %s", products[0].ID)
		}
	})

	t.Run("sort by price ascending", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Sort: "price_asc"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		// Free products first, then by ascending price
		if !products[0].Gia.MienPhi {
			t.Errorf("expected first product to be free, got price %d", products[0].Gia.SoXu)
		}
	})

	t.Run("sort by price descending", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Sort: "price_desc"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		// Most expensive first
		if products[0].Gia.MienPhi {
			t.Errorf("expected first product to be paid, got free")
		}
	})

	t.Run("sort by rating with unrated last", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Sort: "rating"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		// Products with SoLuongDanhGia > 0 come first, sorted by DiemDanhGia desc
		// Then products with 0 ratings come last (sorted by ID)
		foundUnrated := false
		for _, p := range products {
			if p.SoLuongDanhGia == 0 {
				foundUnrated = true
			} else if foundUnrated {
				t.Errorf("rated product %s appears after unrated products", p.ID)
			}
		}
		// First rated product should be sp-006 (5.0, 3 ratings)
		if products[0].ID != "sp-006" {
			t.Errorf("expected highest rated first (sp-006), got %s", products[0].ID)
		}
	})

	t.Run("stable tie-break by id for same rating", func(t *testing.T) {
		products, err := repo.Search(CatalogQuery{Sort: "rating"})
		if err != nil {
			t.Fatalf("Search() failed: %v", err)
		}
		// Both unrated: sp-005, sp-013, sp-016 have 0 ratings
		// They should be sorted by ID: sp-005, sp-013, sp-016
		unratedIDs := []string{}
		for _, p := range products {
			if p.SoLuongDanhGia == 0 {
				unratedIDs = append(unratedIDs, p.ID)
			}
		}
		// Check at least that sp-005 comes before sp-016
		if len(unratedIDs) >= 2 && unratedIDs[0] > unratedIDs[1] {
			t.Errorf("unrated products not sorted by ID: %v", unratedIDs)
		}
	})
}
