// Package catalog defines domain types for the public product catalog.
package catalog

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
	ID             string           `json:"id"`
	Ten            string           `json:"ten"`
	AnhDemo        string           `json:"anh_demo"`
	Gia            Gia              `json:"gia"`
	DanhMuc        DanhMuc          `json:"danh_muc"`
	DiemDanhGia    float64          `json:"diem_danh_gia"`
	SoLuongDanhGia int              `json:"so_luong_danh_gia"`
	TrangThai      TrangThaiSanPham `json:"-"`
}
