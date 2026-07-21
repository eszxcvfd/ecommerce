// Package catalog defines domain types for the public product catalog.
package catalog

import (
	"fmt"
	"time"
)

// Tep represents a single file within a SanPhamSo.
type Tep struct {
	TenTep         string `json:"ten_tep"`
	DinhDang       string `json:"dinh_dang"`
	DungLuongBytes int64  `json:"dung_luong_bytes"`
}

// DanhMuc represents a product category.
type DanhMuc string

const (
	DanhMucKienTruc DanhMuc = "kiến trúc"
	DanhMucCoKhi    DanhMuc = "cơ khí"
	DanhMucDienTu   DanhMuc = "điện tử"
	DanhMucDoHoa    DanhMuc = "đồ họa"
	DanhMucDoAn     DanhMuc = "đồ án"
	DanhMucLuanVan  DanhMuc = "luận văn"
)

// AllDanhMuc lists the six MVP categories in display order.
var AllDanhMuc = []DanhMuc{
	DanhMucKienTruc,
	DanhMucCoKhi,
	DanhMucDienTu,
	DanhMucDoHoa,
	DanhMucDoAn,
	DanhMucLuanVan,
}

// TrangThaiSanPham is the moderation state of a product.
type TrangThaiSanPham string

const (
	TrangThaiDangSoan TrangThaiSanPham = "draft"    // seller is still editing
	TrangThaiChoDuyet TrangThaiSanPham = "pending"  // submitted for review
	TrangThaiDaDuyet  TrangThaiSanPham = "approved" // approved and public
	TrangThaiBiTuChoi TrangThaiSanPham = "rejected" // rejected by admin
	TrangThaiBiAn     TrangThaiSanPham = "hidden"   // hidden after violation report
)

// Gia represents the price of a product.
type Gia struct {
	MienPhi bool  `json:"mien_phi"`
	SoXu    int64 `json:"so_xu,omitempty"` // only meaningful when MienPhi is false
}

// SanPhamSo is a digital product listed on the marketplace.
type SanPhamSo struct {
	ID              string           `json:"id"`
	Ten             string           `json:"ten"`
	MoTa            string           `json:"mo_ta"`
	MoTaChiTiet     string           `json:"mo_ta_chi_tiet"`
	AnhDemo         string           `json:"anh_demo"`
	Gia             Gia              `json:"gia"`
	DanhMuc         DanhMuc          `json:"danh_muc"`
	DinhDang        []string         `json:"dinh_dang"`
	DiemDanhGia     float64          `json:"diem_danh_gia"`
	SoLuongDanhGia  int              `json:"so_luong_danh_gia"`
	NgayTao         time.Time        `json:"ngay_tao"`
	NgayDang        time.Time        `json:"ngay_dang"`
	SoLuotTai       int64            `json:"so_luot_tai"`
	GiayPhep        string           `json:"giay_phep"`
	NguoiBanHienThi string           `json:"nguoi_ban_hien_thi"`
	Tep             []Tep            `json:"tep"`
	NguoiBanID      string           `json:"-"`
	TrangThai       TrangThaiSanPham `json:"-"`
}

// CatalogQuery carries all search/filter/sort parameters for the public catalog.
type CatalogQuery struct {
	Q        string
	DanhMuc  string
	DinhDang string
	MinXu    *int64
	MaxXu    *int64
	Sort     string
}

// SortOrder represents a recognized sort key.
type SortOrder string

const (
	SortNewest    SortOrder = "newest"
	SortPopular   SortOrder = "popular"
	SortPriceAsc  SortOrder = "price_asc"
	SortPriceDesc SortOrder = "price_desc"
	SortRating    SortOrder = "rating"
)

// ValidSortOrders lists all supported sort orders.
var ValidSortOrders = []SortOrder{SortNewest, SortPopular, SortPriceAsc, SortPriceDesc, SortRating}

// SanPhamChiTietResponse is the wrapped response for the product detail endpoint.
type SanPhamChiTietResponse struct {
	SanPham       SanPhamSo   `json:"san_pham"`
	SanPhamDeXuat []SanPhamSo `json:"san_pham_de_xuat"`
}

// TepInput represents a single file entry during draft creation/update.
type TepInput struct {
	TenTep         string `json:"ten_tep"`
	DinhDang       string `json:"dinh_dang"`
	DungLuongBytes int64  `json:"dung_luong_bytes"`
}

// DraftInput is the request body for creating a new draft product.
type DraftInput struct {
	Ten         string     `json:"ten"`
	MoTa        string     `json:"mo_ta"`
	MoTaChiTiet string     `json:"mo_ta_chi_tiet"`
	AnhDemo     string     `json:"anh_demo"`
	MienPhi     bool       `json:"mien_phi"`
	SoXu        int64      `json:"so_xu"`
	DanhMuc     DanhMuc    `json:"danh_muc"`
	GiayPhep    string     `json:"giay_phep"`
	Tep         []TepInput `json:"tep"`
}

// DraftUpdateInput is the request body for updating a draft product.
// All fields are optional; only non-nil fields are updated.
type DraftUpdateInput struct {
	Ten         *string     `json:"ten,omitempty"`
	MoTa        *string     `json:"mo_ta,omitempty"`
	MoTaChiTiet *string     `json:"mo_ta_chi_tiet,omitempty"`
	AnhDemo     *string     `json:"anh_demo,omitempty"`
	MienPhi     *bool       `json:"mien_phi,omitempty"`
	SoXu        *int64      `json:"so_xu,omitempty"`
	DanhMuc     *DanhMuc    `json:"danh_muc,omitempty"`
	GiayPhep    *string     `json:"giay_phep,omitempty"`
	Tep         []TepInput  `json:"tep,omitempty"`
}

// AllowedFormats is the system allowlist of supported file formats.
// The list is lowercased (no leading dot) for consistent matching.
var AllowedFormats = map[string]string{
	"dwg":  "AutoCAD",
	"dxf":  "AutoCAD",
	"skp":  "SketchUp",
	"rvt":  "Revit",
	"rfa":  "Revit",
	"max":  "3ds Max",
	"3ds":  "3ds Max",
	"psd":  "Photoshop",
	"ai":   "Illustrator",
	"eps":  "Illustrator",
}

// IsFormatAllowed checks if the given format (lowercase, no dot) is in the allowlist.
func IsFormatAllowed(format string) bool {
	_, ok := AllowedFormats[format]
	return ok
}

// ValidateDraftInput validates a draft creation request.
func ValidateDraftInput(input DraftInput) error {
	if input.Ten == "" {
		return fmt.Errorf("tên sản phẩm là bắt buộc")
	}
	if input.DanhMuc == "" {
		return fmt.Errorf("danh mục là bắt buộc")
	}
	valid := false
	for _, dm := range AllDanhMuc {
		if input.DanhMuc == dm {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("danh mục không hợp lệ: %s", input.DanhMuc)
	}
	if input.SoXu < 0 {
		return fmt.Errorf("số xu không được âm")
	}
	if len(input.Tep) == 0 {
		return fmt.Errorf("cần ít nhất một tệp")
	}
	for _, tep := range input.Tep {
		if tep.TenTep == "" {
			return fmt.Errorf("tên tệp là bắt buộc")
		}
		if tep.DinhDang == "" {
			return fmt.Errorf("định dạng tệp là bắt buộc")
		}
		if tep.DungLuongBytes <= 0 {
			return fmt.Errorf("dung lượng tệp phải lớn hơn 0")
		}
		if !IsFormatAllowed(tep.DinhDang) {
			return fmt.Errorf("định dạng không được hỗ trợ: %s", tep.DinhDang)
		}
	}
	return nil
}
