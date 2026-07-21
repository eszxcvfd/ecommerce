package account

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const (
	// ContextKeyTaiKhoanID holds the authenticated account's ID in request context.
	ContextKeyTaiKhoanID contextKey = "tai_khoan_id"
	// ContextKeyVaiTro holds the authenticated account's role in request context.
	ContextKeyVaiTro contextKey = "vai_tro"

	// sessionCookieName is the name of the session cookie.
	sessionCookieName = "session_token"
)

// RegisterRoutes mounts account HTTP endpoints on the given mux.
func RegisterRoutes(mux *http.ServeMux, repo AccountRepository) {
	h := &handler{repo: repo}
	// Public endpoints
	mux.HandleFunc("POST /api/v1/dang-ky", h.handleDangKy)
	mux.HandleFunc("POST /api/v1/dang-nhap", h.handleDangNhap)

	// Authenticated endpoints (use middleware wrapper)
	mux.HandleFunc("POST /api/v1/dang-xuat", h.requireAuth(h.handleDangXuat))
	mux.HandleFunc("GET /api/v1/tai-khoan/me", h.requireAuth(h.handleTaiKhoanMe))
	mux.HandleFunc("POST /api/v1/ho-so-nguoi-ban", h.requireAuth(h.handleTaoHOSoBan))
	mux.HandleFunc("GET /api/v1/ho-so-nguoi-ban", h.requireAuth(h.handleHOSoBan))
}

type handler struct {
	repo AccountRepository
}

// requireAuth wraps a handler to require a valid session token.
// The token is extracted from the Authorization header (Bearer) or session cookie.
func (h *handler) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			writeAuthError(w, http.StatusUnauthorized, "thieu_token", "Thiếu token xác thực")
			return
		}

		session, err := h.repo.PhienDangNhapByToken(r.Context(), token)
		if err != nil {
			if errors.Is(err, ErrInvalidToken) {
				writeAuthError(w, http.StatusUnauthorized, "token_khong_hop_le", "Phiên đăng nhập không hợp lệ hoặc đã hết hạn")
				return
			}
			writeAuthError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi xác thực")
			return
		}

		// Fetch account to verify it exists and get role
		acc, err := h.repo.TaiKhoanByID(r.Context(), session.TaiKhoanID)
		if err != nil {
			writeAuthError(w, http.StatusUnauthorized, "tai_khoan_khong_ton_tai", "Tài khoản không tồn tại")
			return
		}
		if acc == nil {
			writeAuthError(w, http.StatusUnauthorized, "tai_khoan_khong_ton_tai", "Tài khoản không tồn tại")
			return
		}

		// Store account info in context
		ctx := context.WithValue(r.Context(), ContextKeyTaiKhoanID, session.TaiKhoanID)
		ctx = context.WithValue(ctx, ContextKeyVaiTro, string(acc.VaiTro))
		next(w, r.WithContext(ctx))
	}
}

// RequireAuth wraps a handler to require a valid session token.
// It is exported so other modules (e.g. catalog) can protect their routes.
func RequireAuth(repo AccountRepository) func(http.HandlerFunc) http.HandlerFunc {
	h := &handler{repo: repo}
	return h.requireAuth
}

// TaiKhoanIDFromContext extracts the authenticated account's ID from context.
func TaiKhoanIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ContextKeyTaiKhoanID).(string)
	return v
}

// VaiTroFromContext extracts the authenticated account's role from context.
func VaiTroFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ContextKeyVaiTro).(string)
	return v
}

// ExtractToken extracts the bearer token from Authorization header or session cookie.
func ExtractToken(r *http.Request) string {
	return extractToken(r)
}

// WriteAuthError writes an unauthorized JSON error response with WWW-Authenticate header.
func WriteAuthError(w http.ResponseWriter, status int, code, message string) {
	writeAuthError(w, status, code, message)
}

// WriteError writes a JSON error response.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	writeError(w, status, code, message)
}

// extractToken extracts the bearer token from Authorization header or session cookie.
func extractToken(r *http.Request) string {
	// Check Authorization header first
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	// Fall back to cookie
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

// setSessionCookie sets the session cookie on the response.
func setSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
}

// clearSessionCookie clears the session cookie.
func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func (h *handler) handleDangKy(w http.ResponseWriter, r *http.Request) {
	var req DangKyRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "thieu_thong_tin", "Email và mật khẩu là bắt buộc")
		return
	}

	acc, err := h.repo.CreateTaiKhoan(r.Context(), req.Email, req.Password, req.Ten)
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "email_da_ton_tai", "Email đã tồn tại trong hệ thống")
			return
		}
		log.Printf("create account error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi tạo tài khoản")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]TaiKhoanPublic{"tai_khoan": acc.ToPublic()})
}

func (h *handler) handleDangNhap(w http.ResponseWriter, r *http.Request) {
	var req DangNhapRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "thieu_thong_tin", "Email và mật khẩu là bắt buộc")
		return
	}

	acc, err := h.repo.TaiKhoanByEmail(r.Context(), req.Email)
	if err != nil {
		log.Printf("lookup account error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi hệ thống")
		return
	}

	if acc == nil {
		writeError(w, http.StatusUnauthorized, "sai_thong_tin", "Email hoặc mật khẩu không đúng")
		return
	}

	// Verify password using bcrypt
	if err := verifyPassword(acc.MatKhauHash, req.Password); err != nil {
		writeError(w, http.StatusUnauthorized, "sai_thong_tin", "Email hoặc mật khẩu không đúng")
		return
	}

	session, err := h.repo.TaoPhienDangNhap(r.Context(), acc.ID)
	if err != nil {
		log.Printf("create session error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi tạo phiên đăng nhập")
		return
	}

	// Set session cookie for browser clients
	setSessionCookie(w, session.Token, session.ExpiresAt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DangNhapResponse{
		TaiKhoan: acc.ToPublic(),
		Token:    session.Token,
	})
}

func (h *handler) handleDangXuat(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token != "" {
		if err := h.repo.XoaPhienDangNhap(r.Context(), token); err != nil {
			log.Printf("delete session error: %v", err)
		}
	}

	clearSessionCookie(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *handler) handleTaiKhoanMe(w http.ResponseWriter, r *http.Request) {
	taiKhoanID := r.Context().Value(ContextKeyTaiKhoanID).(string)

	acc, err := h.repo.TaiKhoanByID(r.Context(), taiKhoanID)
	if err != nil {
		log.Printf("get account error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi lấy thông tin tài khoản")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc.ToPublic())
}

func (h *handler) handleTaoHOSoBan(w http.ResponseWriter, r *http.Request) {
	taiKhoanID := r.Context().Value(ContextKeyTaiKhoanID).(string)

	hs, err := h.repo.TaoHOSoBan(r.Context(), taiKhoanID)
	if err != nil {
		if errors.Is(err, ErrHOSoBanDaTonTai) {
			writeError(w, http.StatusConflict, "ho_so_ban_da_ton_tai", "Hồ sơ người bán đã tồn tại")
			return
		}
		log.Printf("create seller profile error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi tạo hồ sơ người bán")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(hs)
}

func (h *handler) handleHOSoBan(w http.ResponseWriter, r *http.Request) {
	taiKhoanID := r.Context().Value(ContextKeyTaiKhoanID).(string)

	hs, err := h.repo.HOSoBanByTaiKhoanID(r.Context(), taiKhoanID)
	if err != nil {
		if errors.Is(err, ErrHOSoBanNotFound) {
			writeError(w, http.StatusNotFound, "khong_tim_thay", "Chưa có hồ sơ người bán")
			return
		}
		log.Printf("get seller profile error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi lấy hồ sơ người bán")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hs)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// decodeJSON decodes a JSON request body. Returns false and writes an error response on failure.
func decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	// Check Content-Type
	ct := r.Header.Get("Content-Type")
	if ct != "" && !strings.HasPrefix(ct, "application/json") {
		writeError(w, http.StatusUnsupportedMediaType, "sai_dinh_dang", "Yêu cầu Content-Type application/json")
		return false
	}

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "json_khong_hop_le", "Dữ liệu JSON không hợp lệ")
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": code, "message": message})
}

func writeAuthError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": code, "message": message})
}
