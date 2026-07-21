package catalog

import (
	"encoding/json"
	"fmt"
	"time"
)

// CurrentCatalogVersion is the only supported catalog JSON version.
const CurrentCatalogVersion = 2

// CatalogFile is the versioned JSON contract shared by dev seed and production import.
type CatalogFile struct {
	Version  int              `json:"version"`
	Products []CatalogProduct `json:"products"`
}

// CatalogProduct is a single product entry in the versioned catalog JSON.
type CatalogProduct struct {
	ID              string             `json:"id"`
	Ten             string             `json:"ten"`
	MoTa            string             `json:"mo_ta"`
	MoTaChiTiet     string             `json:"mo_ta_chi_tiet,omitempty"`
	AnhDemo         string             `json:"anh_demo"`
	Gia             GiaJSON            `json:"gia"`
	DanhMuc         string             `json:"danh_muc"`
	DinhDang        []string           `json:"dinh_dang"`
	DiemDanhGia     float64            `json:"diem_danh_gia"`
	SoLuongDanhGia  int                `json:"so_luong_danh_gia"`
	NgayTao         string             `json:"ngay_tao"`
	NgayDang        string             `json:"ngay_dang,omitempty"`
	SoLuotTai       int64              `json:"so_luot_tai"`
	GiayPhep        string             `json:"giay_phep,omitempty"`
	NguoiBanHienThi string             `json:"nguoi_ban_hien_thi,omitempty"`
	Tep             []CatalogProductFile `json:"tep,omitempty"`
	TrangThai       string             `json:"trang_thai"`
}

// CatalogProductFile represents a single file entry in the catalog JSON.
type CatalogProductFile struct {
	TenTep         string `json:"ten_tep"`
	DinhDang       string `json:"dinh_dang"`
	DungLuongBytes int64  `json:"dung_luong_bytes"`
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
	var raw struct {
		Version  int               `json:"version"`
		Products []json.RawMessage `json:"products"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if raw.Version < 1 || raw.Version > CurrentCatalogVersion {
		return nil, fmt.Errorf("unsupported version %d (supported: 1–%d)", raw.Version, CurrentCatalogVersion)
	}
	if len(raw.Products) == 0 {
		return nil, fmt.Errorf("catalog file has no products")
	}

	cf := &CatalogFile{Version: raw.Version}

	for i, rawProd := range raw.Products {
		var cp CatalogProduct
		if err := json.Unmarshal(rawProd, &cp); err != nil {
			return nil, fmt.Errorf("product %d: %w", i, err)
		}
		if err := validateCatalogProduct(&cp); err != nil {
			return nil, fmt.Errorf("product %q: %w", cp.ID, err)
		}
		if raw.Version == 1 {
			backfillV1(&cp)
		}
		cf.Products = append(cf.Products, cp)
	}

	return cf, nil
}

// backfillV1 populates new v2 fields from legacy v1 data.
func backfillV1(cp *CatalogProduct) {
	// Derive file metadata from the existing format list: one Tep per format.
	if len(cp.Tep) == 0 && len(cp.DinhDang) > 0 {
		for _, ext := range cp.DinhDang {
			cp.Tep = append(cp.Tep, CatalogProductFile{
				TenTep:         fmt.Sprintf("file.%s", ext),
				DinhDang:       ext,
				DungLuongBytes: 0,
			})
		}
	}
	// Copy creation date as publish date for approved products.
	if cp.NgayDang == "" && cp.TrangThai == "approved" && cp.NgayTao != "" {
		cp.NgayDang = cp.NgayTao
	}
	// Leave MoTaChiTiet, GiayPhep, NguoiBanHienThi empty (zero value).
}

// validateCatalogProduct checks field-level constraints for one product entry.
func validateCatalogProduct(p *CatalogProduct) error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if p.Ten == "" {
		return fmt.Errorf("ten is required")
	}
	if p.DanhMuc == "" {
		return fmt.Errorf("missing danh_muc")
	}
	if p.TrangThai == "" {
		return fmt.Errorf("missing trang_thai")
	}
	if p.Gia.MienPhi && p.Gia.SoXu > 0 {
		return fmt.Errorf("free product must not have so_xu > 0")
	}
	if !p.Gia.MienPhi && p.Gia.SoXu < 1 {
		return fmt.Errorf("paid product must have so_xu > 0")
	}
	// Validate danh_muc
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
	// Validate trang_thai
	valid = false
	for _, tt := range allTrangThai() {
		if tt == p.TrangThai {
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

	var ngayDang time.Time
	if p.NgayDang != "" {
		ngayDang, err = time.Parse(time.RFC3339, p.NgayDang)
		if err != nil {
			return SanPhamSo{}, fmt.Errorf("invalid ngay_dang %q: %w", p.NgayDang, err)
		}
	}

	tepList := make([]Tep, len(p.Tep))
	for i, f := range p.Tep {
		tepList[i] = Tep{
			TenTep:         f.TenTep,
			DinhDang:       f.DinhDang,
			DungLuongBytes: f.DungLuongBytes,
		}
	}

	return SanPhamSo{
		ID:              p.ID,
		Ten:             p.Ten,
		MoTa:            p.MoTa,
		MoTaChiTiet:     p.MoTaChiTiet,
		AnhDemo:         p.AnhDemo,
		Gia:             Gia{MienPhi: p.Gia.MienPhi, SoXu: p.Gia.SoXu},
		DanhMuc:         DanhMuc(p.DanhMuc),
		DinhDang:        p.DinhDang,
		DiemDanhGia:     p.DiemDanhGia,
		SoLuongDanhGia:  p.SoLuongDanhGia,
		NgayTao:         ngayTao,
		NgayDang:        ngayDang,
		SoLuotTai:       p.SoLuotTai,
		GiayPhep:        p.GiayPhep,
		NguoiBanHienThi: p.NguoiBanHienThi,
		Tep:             tepList,
		TrangThai:       TrangThaiSanPham(p.TrangThai),
	}, nil
}

// SanPhamSoToCatalogProduct converts a domain SanPhamSo to a JSON-serializable CatalogProduct.
func SanPhamSoToCatalogProduct(sp SanPhamSo) CatalogProduct {
	tepList := make([]CatalogProductFile, len(sp.Tep))
	for i, f := range sp.Tep {
		tepList[i] = CatalogProductFile{
			TenTep:         f.TenTep,
			DinhDang:       f.DinhDang,
			DungLuongBytes: f.DungLuongBytes,
		}
	}

	var ngayDang string
	if !sp.NgayDang.IsZero() {
		ngayDang = sp.NgayDang.Format(time.RFC3339)
	}

	return CatalogProduct{
		ID:              sp.ID,
		Ten:             sp.Ten,
		MoTa:            sp.MoTa,
		MoTaChiTiet:     sp.MoTaChiTiet,
		AnhDemo:         sp.AnhDemo,
		Gia:             GiaJSON{MienPhi: sp.Gia.MienPhi, SoXu: sp.Gia.SoXu},
		DanhMuc:         string(sp.DanhMuc),
		DinhDang:        sp.DinhDang,
		DiemDanhGia:     sp.DiemDanhGia,
		SoLuongDanhGia:  sp.SoLuongDanhGia,
		NgayTao:         sp.NgayTao.Format(time.RFC3339),
		NgayDang:        ngayDang,
		SoLuotTai:       sp.SoLuotTai,
		GiayPhep:        sp.GiayPhep,
		NguoiBanHienThi: sp.NguoiBanHienThi,
		Tep:             tepList,
		TrangThai:       string(sp.TrangThai),
	}
}
