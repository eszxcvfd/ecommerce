package catalog

import (
	"encoding/json"
	"fmt"
	"time"
)

// CurrentCatalogVersion is the only supported catalog JSON version.
const CurrentCatalogVersion = 1

// CatalogFile is the versioned JSON contract shared by dev seed and production import.
type CatalogFile struct {
	Version  int              `json:"version"`
	Products []CatalogProduct `json:"products"`
}

// CatalogProduct is a single product entry in the versioned catalog JSON.
type CatalogProduct struct {
	ID             string   `json:"id"`
	Ten            string   `json:"ten"`
	MoTa           string   `json:"mo_ta"`
	AnhDemo        string   `json:"anh_demo"`
	Gia            GiaJSON  `json:"gia"`
	DanhMuc        string   `json:"danh_muc"`
	DinhDang       []string `json:"dinh_dang"`
	DiemDanhGia    float64  `json:"diem_danh_gia"`
	SoLuongDanhGia int      `json:"so_luong_danh_gia"`
	NgayTao        string   `json:"ngay_tao"`
	SoLuotTai      int64    `json:"so_luot_tai"`
	TrangThai      string   `json:"trang_thai"`
}

// GiaJSON represents price in the catalog JSON contract.
type GiaJSON struct {
	MienPhi bool  `json:"mien_phi"`
	SoXu    int64 `json:"so_xu,omitempty"`
}

// ValidateCatalogJSON parses and validates a versioned JSON catalog file.
// It checks the version, required fields, and known enum values.
// Returns the parsed CatalogFile on success.
func ValidateCatalogJSON(data []byte) (*CatalogFile, error) {
	var cf CatalogFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parse catalog JSON: %w", err)
	}
	if cf.Version != CurrentCatalogVersion {
		return nil, fmt.Errorf("unsupported catalog version %d (expected %d)", cf.Version, CurrentCatalogVersion)
	}
	if len(cf.Products) == 0 {
		return nil, fmt.Errorf("catalog must contain at least one product")
	}
	for i := range cf.Products {
		if err := validateCatalogProduct(&cf.Products[i]); err != nil {
			return nil, fmt.Errorf("product[%d]: %w", i, err)
		}
	}
	return &cf, nil
}

// validateCatalogProduct checks field-level constraints for one product entry.
func validateCatalogProduct(p *CatalogProduct) error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if p.Ten == "" {
		return fmt.Errorf("ten is required")
	}

	valid := false
	for _, dm := range AllDanhMuc {
		if string(dm) == p.DanhMuc {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid danh_muc %q", p.DanhMuc)
	}

	valid = false
	for _, st := range allTrangThai() {
		if st == p.TrangThai {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid trang_thai %q", p.TrangThai)
	}

	if _, err := time.Parse(time.RFC3339, p.NgayTao); err != nil {
		return fmt.Errorf("invalid ngay_tao %q: %w", p.NgayTao, err)
	}

	if !p.Gia.MienPhi && p.Gia.SoXu < 1 {
		return fmt.Errorf("paid product must have so_xu > 0")
	}
	return nil
}

// allTrangThai returns all valid TrangThaiSanPham values as strings.
func allTrangThai() []string {
	return []string{
		string(TrangThaiDangSoan),
		string(TrangThaiChoDuyet),
		string(TrangThaiDaDuyet),
		string(TrangThaiBiTuChoi),
		string(TrangThaiBiAn),
	}
}

// ToSanPhamSo converts a CatalogProduct to a domain SanPhamSo.
func (p *CatalogProduct) ToSanPhamSo() (SanPhamSo, error) {
	ngayTao, err := time.Parse(time.RFC3339, p.NgayTao)
	if err != nil {
		return SanPhamSo{}, fmt.Errorf("invalid ngay_tao %q: %w", p.NgayTao, err)
	}
	return SanPhamSo{
		ID:             p.ID,
		Ten:            p.Ten,
		MoTa:           p.MoTa,
		AnhDemo:        p.AnhDemo,
		Gia:            Gia{MienPhi: p.Gia.MienPhi, SoXu: p.Gia.SoXu},
		DanhMuc:        DanhMuc(p.DanhMuc),
		DinhDang:       p.DinhDang,
		DiemDanhGia:    p.DiemDanhGia,
		SoLuongDanhGia: p.SoLuongDanhGia,
		NgayTao:        ngayTao,
		SoLuotTai:      p.SoLuotTai,
		TrangThai:      TrangThaiSanPham(p.TrangThai),
	}, nil
}

// SanPhamSoToCatalogProduct converts a domain SanPhamSo to a JSON-serializable CatalogProduct.
func SanPhamSoToCatalogProduct(sp SanPhamSo) CatalogProduct {
	ngayTao := sp.NgayTao.Format(time.RFC3339)
	if sp.NgayTao.IsZero() {
		ngayTao = ""
	}
	return CatalogProduct{
		ID:             sp.ID,
		Ten:            sp.Ten,
		MoTa:           sp.MoTa,
		AnhDemo:        sp.AnhDemo,
		Gia:            GiaJSON{MienPhi: sp.Gia.MienPhi, SoXu: sp.Gia.SoXu},
		DanhMuc:        string(sp.DanhMuc),
		DinhDang:       sp.DinhDang,
		DiemDanhGia:    sp.DiemDanhGia,
		SoLuongDanhGia: sp.SoLuongDanhGia,
		NgayTao:        ngayTao,
		SoLuotTai:      sp.SoLuotTai,
		TrangThai:      string(sp.TrangThai),
	}
}
