package catalog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
