// Package catalog defines domain types for the public product catalog.
package catalog

import "time"

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
