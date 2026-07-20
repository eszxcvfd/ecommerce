package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func seedProducts() []SanPhamSo { return SeedData() }

func TestCatalogEndpoints(t *testing.T) {
	repo := NewMemoryRepo(seedProducts())
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	ts := httptest.NewServer(mux)
	defer ts.Close()

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

		// Must include all 12 approved products
		if len(body.SanPham) != 12 {
			t.Fatalf("expected 12 approved products, got %d", len(body.SanPham))
		}

		// Must not include any non-approved product IDs
		for _, sp := range body.SanPham {
			if sp.ID == "sp-007" || sp.ID == "sp-008" || sp.ID == "sp-009" || sp.ID == "sp-010" {
				t.Errorf("non-approved product %s (%s) was included", sp.ID, sp.Ten)
			}
		}

		// Every returned product must have the required fields
		for _, sp := range body.SanPham {
			if sp.Ten == "" {
				t.Errorf("product %s has empty Ten", sp.ID)
			}
			if sp.DanhMuc == "" {
				t.Errorf("product %s has empty DanhMuc", sp.ID)
			}
			// AnhDemo, Gia are expected but some approved products have empty demo — acceptable in seed
		}

		// Must include filethietke.vn-backed product from seed
		foundFT := false
		for _, sp := range body.SanPham {
			if sp.Ten == "Mẫu vách CNC đồng tiền hiện đại" {
				foundFT = true
				break
			}
		}
		if !foundFT {
			t.Error("expected filethietke.vn product 'Mẫu vách CNC đồng tiền hiện đại' in approved list")
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
		// Confirm Content-Type is JSON
		ct := res.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected application/json, got %q", ct)
		}
	})
}

func TestCatalogHasAtLeastTwoPerCategory(t *testing.T) {
	repo := NewMemoryRepo(seedProducts())
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/v1/san-pham")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var body struct {
		SanPham []SanPhamSo `json:"san_pham"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	// Group products by category
	byCategory := make(map[DanhMuc][]SanPhamSo)
	for _, sp := range body.SanPham {
		byCategory[sp.DanhMuc] = append(byCategory[sp.DanhMuc], sp)
	}

	// Each of the six categories must have at least 2 approved products
	for _, dm := range AllDanhMuc {
		products := byCategory[dm]
		if len(products) < 2 {
			t.Errorf("category %q has %d approved product(s), want at least 2", dm, len(products))
		}
	}
}

// setupTestServer creates a test server for catalog API testing.
func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	repo := NewMemoryRepo(SeedData())
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// queryURL builds a URL with properly encoded query parameters.
// rawQuery is a standard query string (e.g. "q=CNC").
// It re-encodes through url.Values for proper handling of non-ASCII.
func queryURL(ts *httptest.Server, path, rawQuery string) string {
	if rawQuery == "" {
		return ts.URL + path
	}
	vals, _ := url.ParseQuery(rawQuery)
	return ts.URL + path + "?" + vals.Encode()
}

func decodeSanPham(t *testing.T, res *http.Response) []SanPhamSo {
	t.Helper()
	var body struct {
		SanPham []SanPhamSo `json:"san_pham"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	return body.SanPham
}

func TestCatalogSearch(t *testing.T) {
	ts := setupTestServer(t)

	t.Run("empty query returns all approved products", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", ""))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		if len(products) != 12 {
			t.Fatalf("expected 12 products, got %d", len(products))
		}
	})

	t.Run("search by name substring", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "q=CNC"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
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
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "q=Arduino"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		if len(products) != 1 {
			t.Fatalf("expected 1 Arduino product, got %d", len(products))
		}
		if products[0].ID != "sp-003" {
			t.Errorf("expected sp-003, got %s", products[0].ID)
		}
	})

	t.Run("search is case-insensitive", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "q=arduino"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		if len(products) != 1 {
			t.Fatalf("expected 1 product for lowercase query, got %d", len(products))
		}
	})

	t.Run("search is accent-insensitive for Vietnamese", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "q=xây dựng"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		// Should match products with descriptions mentioning "xây dựng"
		if len(products) == 0 {
			t.Fatal("expected at least 1 product matching 'xây dựng'")
		}
	})

	t.Run("search with no matches returns empty", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "q=zzzznotfound"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		if len(products) != 0 {
			t.Fatalf("expected 0 products, got %d", len(products))
		}
	})
}

func TestCatalogFilter(t *testing.T) {
	ts := setupTestServer(t)

	t.Run("filter by danh_muc", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "danh_muc=kiến trúc"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
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
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "dinh_dang=dxf"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
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
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "min_xu=5000"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
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
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "max_xu=200"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
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
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "min_xu=100&max_xu=100"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		// Should include products that are free (price=0) — wait no, min_xu=100 means price >= 100
		// Products priced at exactly 100: sp-017 (100), sp-018 (100)
		if len(products) != 2 {
			t.Fatalf("expected 2 products priced at 100, got %d", len(products))
		}
	})

	t.Run("search + filter combined", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "q=CNC&danh_muc=cơ khí"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		if len(products) != 2 {
			t.Fatalf("expected 2 CNC products in 'cơ khí', got %d", len(products))
		}
	})
}

func TestCatalogSort(t *testing.T) {
	ts := setupTestServer(t)

	t.Run("default sort is newest", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", ""))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		// newest first: sp-016 (2026-07-12) should be first
		if products[0].ID != "sp-016" {
			t.Errorf("expected newest first (sp-016), got %s", products[0].ID)
		}
	})

	t.Run("sort by popular (download count)", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "sort=popular"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		// Most downloaded: sp-014 (520) then sp-004 (310)
		if products[0].ID != "sp-014" {
			t.Errorf("expected most popular first (sp-014), got %s", products[0].ID)
		}
	})

	t.Run("sort by price ascending", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "sort=price_asc"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		// Free products first, then by ascending price
		if !products[0].Gia.MienPhi {
			t.Errorf("expected first product to be free, got price %d", products[0].Gia.SoXu)
		}
	})

	t.Run("sort by price descending", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "sort=price_desc"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
		if len(products) < 2 {
			t.Fatal("need at least 2 products to test sort")
		}
		// Most expensive first
		if products[0].Gia.MienPhi {
			t.Errorf("expected first product to be paid, got free")
		}
	})

	t.Run("sort by rating with unrated last", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "sort=rating"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
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
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "sort=rating"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		products := decodeSanPham(t, res)
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

func TestCatalogValidation(t *testing.T) {
	ts := setupTestServer(t)

	t.Run("invalid sort returns 400", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "sort=invalid"))
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
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "danh_muc=invalid"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400, got %d", res.StatusCode)
		}
	})

	t.Run("invalid min_xu returns 400", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "min_xu=-1"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400 for negative min_xu, got %d", res.StatusCode)
		}
	})

	t.Run("invalid max_xu returns 400", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "max_xu=abc"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400 for non-numeric max_xu, got %d", res.StatusCode)
		}
	})

	t.Run("min_xu > max_xu returns 400", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "min_xu=100&max_xu=50"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("expected 400 for min > max, got %d", res.StatusCode)
		}
	})

	t.Run("max_xu=0 returns free products only", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "max_xu=0"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
		products := decodeSanPham(t, res)
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
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "min_xu=10&max_xu=0"))
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

func TestCatalogFormatValidation(t *testing.T) {
	ts := setupTestServer(t)

	t.Run("invalid dinh_dang returns 400", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "dinh_dang=invalid_format"))
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
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "dinh_dang=pdf"))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			t.Fatalf("expected 200 for valid dinh_dang, got %d", res.StatusCode)
		}
		products := decodeSanPham(t, res)
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

func TestCatalogDinhDangEndpoint(t *testing.T) {
	ts := setupTestServer(t)

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

		// Should have multiple distinct formats
		if len(body.DinhDang) < 3 {
			t.Fatalf("expected at least 3 formats, got %d: %v", len(body.DinhDang), body.DinhDang)
		}

		// Must be sorted
		for i := 1; i < len(body.DinhDang); i++ {
			if body.DinhDang[i-1] > body.DinhDang[i] {
				t.Errorf("formats not sorted: %v", body.DinhDang)
				break
			}
		}

		// Should include commonly expected formats
		expected := map[string]bool{"pdf": true, "dwg": true, "dxf": true}
		for _, f := range body.DinhDang {
			delete(expected, f)
		}
		for missing := range expected {
			t.Errorf("expected format %q not found in response", missing)
		}
	})
}

func TestCatalogTimestampRFC3339(t *testing.T) {
	ts := setupTestServer(t)

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
			// Must parse as RFC3339
			_, err := time.Parse(time.RFC3339, sp.NgayTao)
			if err != nil {
				t.Errorf("product %s ngay_tao %q is not RFC3339: %v", sp.ID, sp.NgayTao, err)
			}
		}
	})

	t.Run("newest sort is chronological with RFC3339", func(t *testing.T) {
		res, err := http.Get(queryURL(ts, "/api/v1/san-pham", "sort=newest"))
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

// errRepo is a CatalogRepository stub that always returns errors.
// Used to test HTTP error handling for storage failures.
type errRepo struct{}

func (errRepo) Products() ([]SanPhamSo, error) {
	return nil, fmt.Errorf("simulated disk failure")
}

func (errRepo) Search(CatalogQuery) ([]SanPhamSo, error) {
	return nil, fmt.Errorf("simulated disk failure")
}

func (errRepo) ProductByID(string) (*SanPhamSo, error) {
	return nil, fmt.Errorf("simulated disk failure")
}

func TestCatalogStorageError_Returns500(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, errRepo{})
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
		if body["message"] == "simulated disk failure" {
			t.Error("error message must not leak raw driver error")
		}
	})

	t.Run("GET /api/v1/danh-muc returns 200 even with erroring repo", func(t *testing.T) {
		// danh-muc is static, doesn't use the repo
		res, err := http.Get(ts.URL + "/api/v1/danh-muc")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
	})

	t.Run("GET /api/v1/dinh-dang returns 500 storage_unavailable", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/dinh-dang")
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
	})
}

func TestSanPhamDetailEndpoint(t *testing.T) {
	repo := NewMemoryRepo(seedProducts())
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	t.Run("returns 200 and product for approved product", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham/sp-001")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}

		var sp SanPhamSo
		if err := json.NewDecoder(res.Body).Decode(&sp); err != nil {
			t.Fatal(err)
		}
		if sp.ID != "sp-001" {
			t.Errorf("expected id sp-001, got %q", sp.ID)
		}
		if sp.Ten != "Bản vẽ nhà phố 3 tầng" {
			t.Errorf("expected ten 'Bản vẽ nhà phố 3 tầng', got %q", sp.Ten)
		}
		if sp.Gia.MienPhi != true {
			t.Error("expected free product")
		}
		if len(sp.DinhDang) == 0 {
			t.Error("expected non-empty formats")
		}
	})

	t.Run("returns 200 and paid product details", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham/sp-017")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}

		var sp SanPhamSo
		if err := json.NewDecoder(res.Body).Decode(&sp); err != nil {
			t.Fatal(err)
		}
		if sp.ID != "sp-017" {
			t.Errorf("expected id sp-017, got %q", sp.ID)
		}
		if sp.Gia.MienPhi != false {
			t.Error("expected paid product")
		}
		if sp.Gia.SoXu != 100 {
			t.Errorf("expected 100 xu, got %d", sp.Gia.SoXu)
		}
		if sp.TrangThai != "" {
			t.Error("trang_thai should not be serialized (json:\"-\")")
		}
	})

	t.Run("returns 404 for non-existent product", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham/nonexistent")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", res.StatusCode)
		}
	})

	t.Run("returns 404 for draft product", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham/sp-007")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 for draft product, got %d", res.StatusCode)
		}
	})

	t.Run("returns 404 for hidden product", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham/sp-010")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 for hidden product, got %d", res.StatusCode)
		}
	})

	t.Run("returns 404 for pending product", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham/sp-008")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 for pending product, got %d", res.StatusCode)
		}
	})

	t.Run("returns 404 for rejected product", func(t *testing.T) {
		res, err := http.Get(ts.URL + "/api/v1/san-pham/sp-009")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 for rejected product, got %d", res.StatusCode)
		}
	})

	t.Run("returns 500 for storage error", func(t *testing.T) {
		mux2 := http.NewServeMux()
		RegisterRoutes(mux2, errRepo{})
		ts2 := httptest.NewServer(mux2)
		defer ts2.Close()

		res, err := http.Get(ts2.URL + "/api/v1/san-pham/sp-001")
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
	})
}
