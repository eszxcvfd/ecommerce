package catalog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// setupSQLiteTestServer creates a test server backed by a seeded SQLite database.
// This exercises the real SQLite adapter through the HTTP API.
func setupSQLiteTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	db := openSQLiteAndSeed(t)
	repo := NewSQLiteRepo(db)
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	RegisterHealthRoutes(mux, nil)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// queryURL builds a URL with properly encoded query parameters.
func sqliteQueryURL(ts *httptest.Server, path, rawQuery string) string {
	if rawQuery == "" {
		return ts.URL + path
	}
	vals, _ := url.ParseQuery(rawQuery)
	return ts.URL + path + "?" + vals.Encode()
}

// decodeSanPham decodes the san_pham response body.
func decodeSP(t *testing.T, res *http.Response) []SanPhamSo {
	t.Helper()
	var body struct {
		SanPham []SanPhamSo `json:"san_pham"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode san_pham: %v", err)
	}
	return body.SanPham
}

// ---------------------------------------------------------------------------
// Endpoints
// ---------------------------------------------------------------------------

func TestCatalogEndpoints_SQLite(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	t.Run("GET /api/v1/danh-muc returns all six categories", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/danh-muc")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
		var body struct {
			DanhMuc []DanhMuc `json:"danh_muc"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.DanhMuc) != 6 {
			t.Fatalf("expected 6 categories, got %d: %v", len(body.DanhMuc), body.DanhMuc)
		}
		expected := []DanhMuc{"kiến trúc", "cơ khí", "điện tử", "đồ họa", "đồ án", "luận văn"}
		for i, dm := range body.DanhMuc {
			if dm != expected[i] {
				t.Errorf("position %d: expected %q, got %q", i, expected[i], dm)
			}
		}
	})

	t.Run("GET /api/v1/san-pham returns only approved products", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
		var body struct {
			SanPham []SanPhamSo `json:"san_pham"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.SanPham) != 12 {
			t.Fatalf("expected 12 approved products, got %d", len(body.SanPham))
		}
		for _, sp := range body.SanPham {
			if sp.ID == "sp-007" || sp.ID == "sp-008" || sp.ID == "sp-009" || sp.ID == "sp-010" {
				t.Errorf("non-approved product %s (%s) was included", sp.ID, sp.Ten)
			}
			if sp.Ten == "" {
				t.Errorf("product %s has empty Ten", sp.ID)
			}
			if sp.DanhMuc == "" {
				t.Errorf("product %s has empty DanhMuc", sp.ID)
			}
		}
	})

	t.Run("GET /api/v1/san-pham serves valid JSON", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
		ct := res.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected application/json, got %q", ct)
		}
	})
}

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

func TestCatalogSearch_SQLite(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	t.Run("empty query returns all approved products", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", ""))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) != 12 {
			t.Fatalf("expected 12 products, got %d", len(products))
		}
	})

	t.Run("search by name substring", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "q=CNC"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) != 2 {
			t.Fatalf("expected 2 CNC products, got %d", len(products))
		}
		for _, p := range products {
			if p.Ten != "Mẫu vách CNC đồng tiền hiện đại" && p.Ten != "Mẫu vách cổng CNC cây nghệ thuật" {
				t.Errorf("unexpected product: %s", p.Ten)
			}
		}
	})

	t.Run("search by description content", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "q=Arduino"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) != 1 {
			t.Fatalf("expected 1 Arduino product, got %d", len(products))
		}
		if products[0].ID != "sp-003" {
			t.Errorf("expected sp-003, got %s", products[0].ID)
		}
	})

	t.Run("search is case-insensitive", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "q=arduino"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) != 1 {
			t.Fatalf("expected 1 product for lowercase query, got %d", len(products))
		}
	})

	t.Run("search is accent-insensitive for Vietnamese", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "q=xây dựng"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) == 0 {
			t.Fatal("expected at least 1 product matching 'xây dựng'")
		}
	})

	t.Run("search with no matches returns empty", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "q=zzzznotfound"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) != 0 {
			t.Fatalf("expected 0 products, got %d", len(products))
		}
	})
}

// ---------------------------------------------------------------------------
// Filter
// ---------------------------------------------------------------------------

func TestCatalogFilter_SQLite(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	t.Run("filter by danh_muc", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "danh_muc=kiến trúc"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
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
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "dinh_dang=dxf"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
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
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "min_xu=5000"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
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
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "max_xu=200"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
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
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "min_xu=100&max_xu=100"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) != 2 {
			t.Fatalf("expected 2 products priced at 100, got %d", len(products))
		}
	})

	t.Run("search + filter combined", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "q=CNC&danh_muc=cơ khí"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) != 2 {
			t.Fatalf("expected 2 CNC products in 'cơ khí', got %d", len(products))
		}
	})
}

// ---------------------------------------------------------------------------
// Sort
// ---------------------------------------------------------------------------

func TestCatalogSort_SQLite(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	t.Run("default sort is newest", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", ""))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		if products[0].ID != "sp-016" {
			t.Errorf("expected newest first (sp-016), got %s", products[0].ID)
		}
	})

	t.Run("sort by popular (download count)", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "sort=popular"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		if products[0].ID != "sp-014" {
			t.Errorf("expected most popular first (sp-014), got %s", products[0].ID)
		}
	})

	t.Run("sort by price ascending", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "sort=price_asc"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		if !products[0].Gia.MienPhi {
			t.Errorf("expected first product to be free, got price %d", products[0].Gia.SoXu)
		}
	})

	t.Run("sort by price descending", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "sort=price_desc"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		if products[0].Gia.MienPhi {
			t.Errorf("expected first product to be paid, got free")
		}
	})

	t.Run("sort by rating with unrated last", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "sort=rating"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		foundUnrated := false
		for _, p := range products {
			if p.SoLuongDanhGia == 0 {
				foundUnrated = true
			} else if foundUnrated {
				t.Errorf("rated product %s appears after unrated products", p.ID)
			}
		}
		if products[0].ID != "sp-006" {
			t.Errorf("expected highest rated first (sp-006), got %s", products[0].ID)
		}
	})

	t.Run("stable tie-break by id for same rating", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "sort=rating"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSP(t, res)
		unratedIDs := []string{}
		for _, p := range products {
			if p.SoLuongDanhGia == 0 {
				unratedIDs = append(unratedIDs, p.ID)
			}
		}
		if len(unratedIDs) >= 2 && unratedIDs[0] > unratedIDs[1] {
			t.Errorf("unrated products not sorted by ID: %v", unratedIDs)
		}
	})
}

// ---------------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------------

func TestCatalogValidation_SQLite(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	t.Run("invalid sort returns 400", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "sort=invalid"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400, got %d", res.StatusCode)
		}
		var body map[string]string
		json.NewDecoder(res.Body).Decode(&body)
		if body["error"] != "invalid_filter" {
			t.Errorf("expected error 'invalid_filter', got %q", body["error"])
		}
	})

	t.Run("invalid danh_muc returns 400", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "danh_muc=invalid"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400, got %d", res.StatusCode)
		}
	})

	t.Run("invalid min_xu returns 400", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "min_xu=-1"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400 for negative min_xu, got %d", res.StatusCode)
		}
	})

	t.Run("invalid max_xu returns 400", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "max_xu=abc"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400 for non-numeric max_xu, got %d", res.StatusCode)
		}
	})

	t.Run("min_xu > max_xu returns 400", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "min_xu=100&max_xu=50"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400 for min > max, got %d", res.StatusCode)
		}
	})

	t.Run("max_xu=0 returns free products only", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "max_xu=0"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
		products := decodeSP(t, res)
		if len(products) == 0 {
			t.Fatal("expected at least one free product")
		}
		for _, p := range products {
			if !p.Gia.MienPhi {
				t.Errorf("product %s is not free but returned for max_xu=0", p.ID)
			}
		}
	})

	t.Run("min_xu > max_xu with max_xu=0 returns 400", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "min_xu=10&max_xu=0"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400 for min_xu=10 > max_xu=0, got %d", res.StatusCode)
		}
		var body map[string]string
		json.NewDecoder(res.Body).Decode(&body)
		if body["error"] != "invalid_filter" {
			t.Errorf("expected error 'invalid_filter', got %q", body["error"])
		}
	})
}

// ---------------------------------------------------------------------------
// Format validation and endpoint
// ---------------------------------------------------------------------------

func TestCatalogFormatValidation_SQLite(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	t.Run("invalid dinh_dang returns 400", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "dinh_dang=invalid_format"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400 for invalid dinh_dang, got %d", res.StatusCode)
		}
		var body map[string]string
		json.NewDecoder(res.Body).Decode(&body)
		if body["error"] != "invalid_filter" {
			t.Errorf("expected error 'invalid_filter', got %q", body["error"])
		}
	})

	t.Run("valid dinh_dang passes through", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "dinh_dang=pdf"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			t.Fatalf("expected 200 for valid dinh_dang, got %d", res.StatusCode)
		}
		products := decodeSP(t, res)
		if len(products) == 0 {
			t.Fatal("expected at least one PDF product")
		}
		for _, p := range products {
			has := false
			for _, ext := range p.DinhDang {
				if ext == "pdf" {
					has = true
					break
				}
			}
			if !has {
				t.Errorf("product %s does not have pdf format", p.ID)
			}
		}
	})
}

func TestCatalogDinhDangEndpoint_SQLite(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	t.Run("GET /api/v1/dinh-dang returns distinct formats sorted", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/dinh-dang")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
		var body struct {
			DinhDang []string `json:"dinh_dang"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.DinhDang) < 3 {
			t.Fatalf("expected at least 3 formats, got %d: %v", len(body.DinhDang), body.DinhDang)
		}
		for i := 1; i < len(body.DinhDang); i++ {
			if body.DinhDang[i-1] > body.DinhDang[i] {
				t.Errorf("formats not sorted: %v", body.DinhDang)
				break
			}
		}
		expected := map[string]bool{"pdf": true, "dwg": true, "dxf": true}
		for _, f := range body.DinhDang {
			delete(expected, f)
		}
		for missing := range expected {
			t.Errorf("expected format %q not found in response", missing)
		}
	})
}

// ---------------------------------------------------------------------------
// Timestamp RFC3339
// ---------------------------------------------------------------------------

func TestCatalogTimestampRFC3339_SQLite(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	t.Run("ngay_tao is RFC3339 in JSON", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var body struct {
			SanPham []struct {
				ID      string `json:"id"`
				NgayTao string `json:"ngay_tao"`
			} `json:"san_pham"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.SanPham) == 0 {
			t.Fatal("expected at least one product")
		}
		for _, sp := range body.SanPham {
			_, err := time.Parse(time.RFC3339, sp.NgayTao)
			if err != nil {
				t.Errorf("product %s ngay_tao %q is not RFC3339: %v", sp.ID, sp.NgayTao, err)
			}
		}
	})

	t.Run("newest sort is chronological with RFC3339", func(t *testing.T) {
		res, err := http.Get(sqliteQueryURL(ts, "/api/v1/san-pham", "sort=newest"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var body struct {
			SanPham []struct {
				ID      string `json:"id"`
				NgayTao string `json:"ngay_tao"`
			} `json:"san_pham"`
		}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.SanPham) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		for i, sp := range body.SanPham {
			tm, err := time.Parse(time.RFC3339, sp.NgayTao)
			if err != nil {
				t.Fatalf("product %s at position %d has invalid ngay_tao: %v", sp.ID, i, err)
			}
			if i > 0 {
				prev, _ := time.Parse(time.RFC3339, body.SanPham[i-1].NgayTao)
				if tm.After(prev) {
					t.Errorf("position %d (product %s, %s) is newer than position %d (product %s, %s) — should be newest first",
						i, sp.ID, tm.Format(time.RFC3339), i-1, body.SanPham[i-1].ID, prev.Format(time.RFC3339))
				}
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Storage error: SQLite-specific
// ---------------------------------------------------------------------------

func TestCatalogSQLiteStorageError_Returns500(t *testing.T) {
	// Open a SQLite DB, close it immediately, then try to query it.
	// This tests that the SQLite adapter returns errors and the HTTP layer
	// maps them to 500 storage_unavailable without leaking details.
	db, err := OpenSQLite(t.TempDir() + "/gone.db")
	if err != nil {
		t.Fatalf("OpenSQLite: %v", err)
	}
	// Close the DB so queries fail
	db.Close()

	repo := NewSQLiteRepo(db)
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	t.Run("GET /api/v1/san-pham returns 500 storage_unavailable", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", res.StatusCode)
		}
		var body map[string]string
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["error"] != "storage_unavailable" {
			t.Errorf("expected error 'storage_unavailable', got %q", body["error"])
		}
		if body["message"] == "" {
			t.Error("expected non-empty error message")
		}
		// Must not leak SQL or paths
		if body["message"] == "database is closed" {
			t.Error("error message must not leak raw driver error")
		}
	})

	t.Run("GET /api/v1/danh-muc returns 200 even with closed db", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/danh-muc")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
	})
}
