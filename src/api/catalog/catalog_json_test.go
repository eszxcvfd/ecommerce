package catalog

import (
	"testing"
	"time"
)

func TestValidateCatalogJSON_ValidMinimal(t *testing.T) {
	data := `{
		"version": 1,
		"products": [
			{
				"id": "sp-001",
				"ten": "Test Product",
				"mo_ta": "A description",
				"anh_demo": "/images/test.jpg",
				"gia": {"mien_phi": true, "so_xu": 0},
				"danh_muc": "kiến trúc",
				"dinh_dang": ["pdf"],
				"diem_danh_gia": 4.5,
				"so_luong_danh_gia": 10,
				"ngay_tao": "2026-07-01T00:00:00Z",
				"so_luot_tai": 100,
				"trang_thai": "approved"
			}
		]
	}`
	cf, err := ValidateCatalogJSON([]byte(data))
	if err != nil {
		t.Fatalf("expected valid catalog: %v", err)
	}
	if len(cf.Products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(cf.Products))
	}
	if cf.Products[0].ID != "sp-001" {
		t.Errorf("expected ID sp-001, got %s", cf.Products[0].ID)
	}
	if cf.Products[0].Ten != "Test Product" {
		t.Errorf("expected Ten 'Test Product', got %s", cf.Products[0].Ten)
	}
	if cf.Version != 1 {
		t.Errorf("expected version 1, got %d", cf.Version)
	}
}

func TestValidateCatalogJSON_RejectsUnsupportedVersion(t *testing.T) {
	data := `{"version": 42, "products": [
		{"id":"x","ten":"x","danh_muc":"kiến trúc","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"}
	]}`
	_, err := ValidateCatalogJSON([]byte(data))
	if err == nil {
		t.Fatal("expected error for unsupported version 42")
	}
}

func TestValidateCatalogJSON_RejectsEmptyProducts(t *testing.T) {
	data := `{"version": 1, "products": []}`
	_, err := ValidateCatalogJSON([]byte(data))
	if err == nil {
		t.Fatal("expected error for empty products")
	}
}

func TestValidateCatalogJSON_RejectsMissingID(t *testing.T) {
	data := `{"version": 1, "products": [
		{"ten":"x","danh_muc":"kiến trúc","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"}
	]}`
	_, err := ValidateCatalogJSON([]byte(data))
	if err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestValidateCatalogJSON_RejectsMissingTen(t *testing.T) {
	data := `{"version": 1, "products": [
		{"id":"sp-001","danh_muc":"kiến trúc","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"}
	]}`
	_, err := ValidateCatalogJSON([]byte(data))
	if err == nil {
		t.Fatal("expected error for missing ten")
	}
}

func TestValidateCatalogJSON_RejectsInvalidDanhMuc(t *testing.T) {
	data := `{"version": 1, "products": [
		{"id":"sp-001","ten":"x","danh_muc":"invalid","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"approved"}
	]}`
	_, err := ValidateCatalogJSON([]byte(data))
	if err == nil {
		t.Fatal("expected error for invalid danh_muc")
	}
}

func TestValidateCatalogJSON_RejectsInvalidTrangThai(t *testing.T) {
	data := `{"version": 1, "products": [
		{"id":"sp-001","ten":"x","danh_muc":"kiến trúc","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":true},"trang_thai":"unknown"}
	]}`
	_, err := ValidateCatalogJSON([]byte(data))
	if err == nil {
		t.Fatal("expected error for invalid trang_thai")
	}
}

func TestValidateCatalogJSON_RejectsBadTimestamp(t *testing.T) {
	data := `{"version": 1, "products": [
		{"id":"sp-001","ten":"x","danh_muc":"kiến trúc","ngay_tao":"not-a-timestamp","gia":{"mien_phi":true},"trang_thai":"approved"}
	]}`
	_, err := ValidateCatalogJSON([]byte(data))
	if err == nil {
		t.Fatal("expected error for bad timestamp")
	}
}

func TestValidateCatalogJSON_RejectsPaidProductWithoutSoXu(t *testing.T) {
	data := `{"version": 1, "products": [
		{"id":"sp-001","ten":"x","danh_muc":"kiến trúc","ngay_tao":"2026-07-01T00:00:00Z","gia":{"mien_phi":false,"so_xu":0},"trang_thai":"approved"}
	]}`
	_, err := ValidateCatalogJSON([]byte(data))
	if err == nil {
		t.Fatal("expected error for paid product with so_xu=0")
	}
}

func TestValidateCatalogJSON_RejectsMalformedJSON(t *testing.T) {
	data := `{not-json}`
	_, err := ValidateCatalogJSON([]byte(data))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestSanPhamSoToCatalogProduct_RoundTrip(t *testing.T) {
	original := SanPhamSo{
		ID: "sp-test", Ten: "Test", MoTa: "Desc",
		AnhDemo: "/img.jpg",
		Gia:     Gia{MienPhi: false, SoXu: 5000},
		DanhMuc: DanhMucDoHoa, DinhDang: []string{"svg", "png"},
		DiemDanhGia: 4.2, SoLuongDanhGia: 8,
		NgayTao:   testParseTime(t, "2026-07-04T00:00:00Z"),
		SoLuotTai: 310, TrangThai: TrangThaiDaDuyet,
	}

	cp := SanPhamSoToCatalogProduct(original)
	if cp.ID != original.ID {
		t.Errorf("ID: got %s, want %s", cp.ID, original.ID)
	}
	if cp.Ten != original.Ten {
		t.Errorf("Ten: got %s, want %s", cp.Ten, original.Ten)
	}
	if cp.DanhMuc != string(original.DanhMuc) {
		t.Errorf("DanhMuc: got %s, want %s", cp.DanhMuc, original.DanhMuc)
	}
	if cp.TrangThai != string(original.TrangThai) {
		t.Errorf("TrangThai: got %s, want %s", cp.TrangThai, original.TrangThai)
	}
	if cp.Gia.MienPhi {
		t.Errorf("expected paid product, got mien_phi=true")
	}
	if cp.Gia.SoXu != 5000 {
		t.Errorf("SoXu: got %d, want 5000", cp.Gia.SoXu)
	}
}

func TestToSanPhamSo_RoundTrip(t *testing.T) {
	cp := CatalogProduct{
		ID: "sp-rt", Ten: "Round Trip", MoTa: "Testing",
		AnhDemo: "/img.jpg",
		Gia:     GiaJSON{MienPhi: false, SoXu: 10000},
		DanhMuc: "điện tử", DinhDang: []string{"brd", "sch"},
		DiemDanhGia: 3.8, SoLuongDanhGia: 5,
		NgayTao: "2026-07-03T00:00:00Z", SoLuotTai: 230,
		TrangThai: "approved",
	}

	sp, err := cp.ToSanPhamSo()
	if err != nil {
		t.Fatalf("ToSanPhamSo failed: %v", err)
	}
	if sp.ID != cp.ID {
		t.Errorf("ID: got %s, want %s", sp.ID, cp.ID)
	}
	if sp.Ten != cp.Ten {
		t.Errorf("Ten: got %s, want %s", sp.Ten, cp.Ten)
	}
	if string(sp.DanhMuc) != cp.DanhMuc {
		t.Errorf("DanhMuc: got %s, want %s", sp.DanhMuc, cp.DanhMuc)
	}
	if string(sp.TrangThai) != cp.TrangThai {
		t.Errorf("TrangThai: got %s, want %s", sp.TrangThai, cp.TrangThai)
	}
	if sp.Gia.MienPhi {
		t.Errorf("expected paid product, got mien_phi=true")
	}
	if sp.Gia.SoXu != 10000 {
		t.Errorf("SoXu: got %d, want 10000", sp.Gia.SoXu)
	}
}

func TestToSanPhamSo_RejectsInvalidTimestamp(t *testing.T) {
	cp := CatalogProduct{
		ID: "sp-bad", Ten: "Bad", DanhMuc: "kiến trúc",
		NgayTao:   "invalid-date",
		Gia:       GiaJSON{MienPhi: true},
		TrangThai: "approved",
	}
	_, err := cp.ToSanPhamSo()
	if err == nil {
		t.Fatal("expected error for invalid timestamp")
	}
}

func testParseTime(t *testing.T, s string) time.Time {
	t.Helper()
	ts, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return ts
}
