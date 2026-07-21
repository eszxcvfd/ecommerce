// Package account defines domain types for the account and authentication system.
package account

import "time"

// VaiTro represents an account's authorization role.
type VaiTro string

const (
	// VaiTroNguoiMua is a regular buyer account (Người mua cá nhân).
	VaiTroNguoiMua VaiTro = "nguoi_mua"
	// VaiTroAdmin is a marketplace administrator.
	VaiTroAdmin VaiTro = "admin"
)

// TrangThaiHOSoBan is the activation state of a seller profile.
type TrangThaiHOSoBan string

const (
	// TrangThaiHOSoBanKichHoat means the seller profile is active.
	TrangThaiHOSoBanKichHoat TrangThaiHOSoBan = "kich_hoat"
	// TrangThaiHOSoBanVoHieuHoa means the seller profile is deactivated.
	TrangThaiHOSoBanVoHieuHoa TrangThaiHOSoBan = "vo_hieu_hoa"
)

// TaiKhoan represents a user account that can act as a buyer and optionally
// activate a seller profile.
type TaiKhoan struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	MatKhauHash string   `json:"-"` // never serialized
	Ten        string    `json:"ten"`
	VaiTro     VaiTro    `json:"vai_tro"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TaiKhoanPublic is the public-safe representation of an account.
type TaiKhoanPublic struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Ten       string    `json:"ten"`
	VaiTro    VaiTro    `json:"vai_tro"`
	CreatedAt time.Time `json:"created_at"`
}

// ToPublic converts a TaiKhoan to its public-safe representation.
func (a *TaiKhoan) ToPublic() TaiKhoanPublic {
	return TaiKhoanPublic{
		ID:        a.ID,
		Email:     a.Email,
		Ten:       a.Ten,
		VaiTro:    a.VaiTro,
		CreatedAt: a.CreatedAt,
	}
}

// HOSoBan represents a seller profile activated on an account.
type HOSoBan struct {
	ID         string            `json:"id"`
	TaiKhoanID string            `json:"tai_khoan_id"`
	TrangThai  TrangThaiHOSoBan  `json:"trang_thai"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// PhienDangNhap represents an authentication session (login session).
type PhienDangNhap struct {
	ID         string    `json:"id"`
	TaiKhoanID string    `json:"tai_khoan_id"`
	Token      string    `json:"token"` // bearer token stored as SHA-256 hash
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// DangKyRequest is the registration request body.
type DangKyRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Ten      string `json:"ten"`
}

// DangNhapRequest is the login request body.
type DangNhapRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// DangNhapResponse is the login response body.
type DangNhapResponse struct {
	TaiKhoan TaiKhoanPublic `json:"tai_khoan"`
	Token    string         `json:"token"`
}

// TaoHOSoBanRequest is the request to activate a seller profile.
type TaoHOSoBanRequest struct {
	// Optional display name for the seller profile; defaults to account ten.
	TenHienThi string `json:"ten_hien_thi,omitempty"`
}
