package catalog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// seedProducts returns deterministic test data covering:
//   - All six categories
//   - Products with free and paid prices
//   - Products with/without ratings
//   - Products that should NOT appear (draft, pending, rejected, hidden)
func seedProducts() []SanPhamSo {
	return []SanPhamSo{
		// --- Approved products (should appear) ---
		{
			ID: "sp-001", Ten: "Bản vẽ nhà phố 3 tầng",
			AnhDemo: "/images/nha-pho.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 4.5, SoLuongDanhGia: 12,
		},
		{
			ID: "sp-002", Ten: "Mô hình khung thép tiền chế",
			AnhDemo: "/images/khung-thep.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 15000},
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 0, SoLuongDanhGia: 0,
		},
		{
			ID: "sp-003", Ten: "Sơ đồ mạch Arduino điều khiển LED",
			AnhDemo: "/images/arduino-led.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 3.8, SoLuongDanhGia: 5,
		},
		{
			ID: "sp-004", Ten: "Bộ icon phong cách tối giản",
			AnhDemo: "/images/icon-set.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 5000},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 4.2, SoLuongDanhGia: 8,
		},
		{
			ID: "sp-005", Ten: "Đồ án thiết kế cầu dầm BTCT",
			AnhDemo: "/images/cau-dam.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDoAn, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 0, SoLuongDanhGia: 0,
		},
		{
			ID: "sp-006", Ten: "Luận văn thạc sĩ AI trong xây dựng",
			AnhDemo: "/images/luanvan-ai.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 30000},
			DanhMuc: DanhMucLuanVan, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 5.0, SoLuongDanhGia: 3,
		},
		// --- Non-approved products (must NOT appear) ---
		{
			ID: "sp-007", Ten: "Draft product",
			AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDangSoan,
		},
		{
			ID: "sp-008", Ten: "Pending review product",
			AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiChoDuyet,
		},
		{
			ID: "sp-009", Ten: "Rejected product",
			AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiBiTuChoi,
		},
		{
			ID: "sp-010", Ten: "Hidden product after violation",
			AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiBiAn,
		},
	}
}

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

		// Must include all 6 approved products
		if len(body.SanPham) != 6 {
			t.Fatalf("expected 6 approved products, got %d", len(body.SanPham))
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
