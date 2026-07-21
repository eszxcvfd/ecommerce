package account

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrDuplicateEmail is returned when registering with an email that already exists.
var ErrDuplicateEmail = errors.New("email đã tồn tại")

// ErrInvalidCredentials is returned when login credentials are incorrect.
var ErrInvalidCredentials = errors.New("email hoặc mật khẩu không đúng")

// ErrAccountNotFound is returned when the account is not found.
var ErrAccountNotFound = errors.New("không tìm thấy tài khoản")

// ErrInvalidToken is returned when the session token is invalid or expired.
var ErrInvalidToken = errors.New("phiên đăng nhập không hợp lệ")

// ErrHOSoBanDaTonTai is returned when the seller profile already exists.
var ErrHOSoBanDaTonTai = errors.New("hồ sơ người bán đã tồn tại")

// ErrHOSoBanNotFound is returned when no seller profile exists.
var ErrHOSoBanNotFound = errors.New("không tìm thấy hồ sơ người bán")

// AccountRepository is the seam for account and authentication data access.
type AccountRepository interface {
	// CreateTaiKhoan creates a new account with the given details.
	// Returns ErrDuplicateEmail if the email is already in use.
	CreateTaiKhoan(ctx context.Context, email, password, ten string) (*TaiKhoan, error)

	// TaiKhoanByEmail looks up an account by email. Returns nil, nil if not found.
	TaiKhoanByEmail(ctx context.Context, email string) (*TaiKhoan, error)

	// TaiKhoanByID looks up an account by ID. Returns nil, nil if not found.
	TaiKhoanByID(ctx context.Context, id string) (*TaiKhoan, error)

	// TaoPhienDangNhap creates a new login session for the given account.
	// Returns the session with a plain-text token (also stored as hash).
	TaoPhienDangNhap(ctx context.Context, taiKhoanID string) (*PhienDangNhap, error)

	// PhienDangNhapByToken looks up a session by its plain-text token.
	// Returns ErrInvalidToken if not found or expired.
	PhienDangNhapByToken(ctx context.Context, token string) (*PhienDangNhap, error)

	// XoaPhienDangNhap deletes a session (logout).
	XoaPhienDangNhap(ctx context.Context, token string) error

	// TaoHOSoBan creates a seller profile for the given account.
	// Returns ErrHOSoBanDaTonTai if one already exists.
	TaoHOSoBan(ctx context.Context, taiKhoanID string) (*HOSoBan, error)

	// HOSoBanByTaiKhoanID returns the seller profile for an account.
	// Returns ErrHOSoBanNotFound if no profile exists.
	HOSoBanByTaiKhoanID(ctx context.Context, taiKhoanID string) (*HOSoBan, error)
}

// ---------------------------------------------------------------------------
// In-memory adapter (for tests)
// ---------------------------------------------------------------------------

type memoryAccountRepo struct {
	mu            sync.RWMutex
	taiKhoans     []TaiKhoan
	phienDangNhap map[string]PhienDangNhap // token hash -> session
	hoSoBans      []HOSoBan
	nextID        int64
}

// NewMemoryRepo creates an in-memory AccountRepository for testing.
func NewMemoryRepo() AccountRepository {
	return &memoryAccountRepo{
		taiKhoans:     []TaiKhoan{},
		phienDangNhap: map[string]PhienDangNhap{},
		hoSoBans:      []HOSoBan{},
		nextID:        1,
	}
}

func (r *memoryAccountRepo) CreateTaiKhoan(ctx context.Context, email, password, ten string) (*TaiKhoan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, a := range r.taiKhoans {
		if a.Email == email {
			return nil, ErrDuplicateEmail
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	id := fmt.Sprintf("TK%04d", r.nextID)
	r.nextID++

	acc := TaiKhoan{
		ID:          id,
		Email:       email,
		MatKhauHash: string(hash),
		Ten:         ten,
		VaiTro:      VaiTroNguoiMua,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.taiKhoans = append(r.taiKhoans, acc)
	return &acc, nil
}

func (r *memoryAccountRepo) TaiKhoanByEmail(ctx context.Context, email string) (*TaiKhoan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, a := range r.taiKhoans {
		if a.Email == email {
			return &a, nil
		}
	}
	return nil, nil
}

func (r *memoryAccountRepo) TaiKhoanByID(ctx context.Context, id string) (*TaiKhoan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, a := range r.taiKhoans {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, nil
}

func (r *memoryAccountRepo) TaoPhienDangNhap(ctx context.Context, taiKhoanID string) (*PhienDangNhap, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))

	now := time.Now().UTC()
	s := PhienDangNhap{
		ID:         fmt.Sprintf("PDN%04d", r.nextID),
		TaiKhoanID: taiKhoanID,
		Token:      tokenHash,
		CreatedAt:  now,
		ExpiresAt:  now.Add(24 * time.Hour),
	}
	r.nextID++

	r.mu.Lock()
	r.phienDangNhap[tokenHash] = s
	r.mu.Unlock()

	// Return with plain-text token for the caller
	s.Token = token
	return &s, nil
}

func (r *memoryAccountRepo) PhienDangNhapByToken(ctx context.Context, token string) (*PhienDangNhap, error) {
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))

	r.mu.RLock()
	s, ok := r.phienDangNhap[tokenHash]
	r.mu.RUnlock()

	if !ok {
		return nil, ErrInvalidToken
	}
	if time.Now().UTC().After(s.ExpiresAt) {
		r.mu.Lock()
		delete(r.phienDangNhap, tokenHash)
		r.mu.Unlock()
		return nil, ErrInvalidToken
	}
	return &s, nil
}

func (r *memoryAccountRepo) XoaPhienDangNhap(ctx context.Context, token string) error {
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
	r.mu.Lock()
	delete(r.phienDangNhap, tokenHash)
	r.mu.Unlock()
	return nil
}

func (r *memoryAccountRepo) TaoHOSoBan(ctx context.Context, taiKhoanID string) (*HOSoBan, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, hs := range r.hoSoBans {
		if hs.TaiKhoanID == taiKhoanID {
			return nil, ErrHOSoBanDaTonTai
		}
	}

	now := time.Now().UTC()
	hs := HOSoBan{
		ID:         fmt.Sprintf("HS%04d", r.nextID),
		TaiKhoanID: taiKhoanID,
		TrangThai:  TrangThaiHOSoBanKichHoat,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	r.nextID++
	r.hoSoBans = append(r.hoSoBans, hs)
	return &hs, nil
}

func (r *memoryAccountRepo) HOSoBanByTaiKhoanID(ctx context.Context, taiKhoanID string) (*HOSoBan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, hs := range r.hoSoBans {
		if hs.TaiKhoanID == taiKhoanID {
			return &hs, nil
		}
	}
	return nil, ErrHOSoBanNotFound
}

// WithAdminAccount seeds an admin account into a memory repo for testing.
func WithAdminAccount(repo AccountRepository, email, password, ten string) AccountRepository {
	mr := repo.(*memoryAccountRepo)
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mr.mu.Lock()
	defer mr.mu.Unlock()

	now := time.Now().UTC()
	id := fmt.Sprintf("TK%04d", mr.nextID)
	mr.nextID++

	mr.taiKhoans = append(mr.taiKhoans, TaiKhoan{
		ID:          id,
		Email:       email,
		MatKhauHash: string(hash),
		Ten:         ten,
		VaiTro:      VaiTroAdmin,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	return repo
}

// accounts returns all accounts for test inspection.
func (r *memoryAccountRepo) accounts() []TaiKhoan {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]TaiKhoan, len(r.taiKhoans))
	copy(result, r.taiKhoans)
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}
