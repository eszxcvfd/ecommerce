package account

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations runs account module migrations on the given database.
// It's safe to call multiple times (goose tracks applied migrations).
func RunMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("account migrations: %w", err)
	}
	return nil
}

// VerifyMigrations checks that all account migrations have been applied.
func VerifyMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}
	return goose.UpByOne(db, "migrations")
}

// ---------------------------------------------------------------------------
// SQLite repository adapter
// ---------------------------------------------------------------------------

type sqliteRepo struct {
	db *sql.DB
}

// NewSQLiteRepo creates an AccountRepository backed by the given SQLite DB.
func NewSQLiteRepo(db *sql.DB) AccountRepository {
	return &sqliteRepo{db: db}
}

func (r *sqliteRepo) CreateTaiKhoan(ctx context.Context, email, password, ten string) (*TaiKhoan, error) {
	var existing int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tai_khoan WHERE email = ?", email).Scan(&existing); err != nil {
		return nil, fmt.Errorf("check duplicate: %w", err)
	}
	if existing > 0 {
		return nil, ErrDuplicateEmail
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id := newID("TK")

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO tai_khoan (id, email, mat_khau_hash, ten, vai_tro, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, email, hash, ten, VaiTroNguoiMua, now, now)
	if err != nil {
		return nil, fmt.Errorf("insert account: %w", err)
	}

	createdAt, _ := time.Parse(time.RFC3339, now)
	return &TaiKhoan{
		ID:          id,
		Email:       email,
		MatKhauHash: hash,
		Ten:         ten,
		VaiTro:      VaiTroNguoiMua,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}, nil
}

func (r *sqliteRepo) TaiKhoanByEmail(ctx context.Context, email string) (*TaiKhoan, error) {
	return scanTaiKhoan(r.db.QueryRowContext(ctx,
		`SELECT id, email, mat_khau_hash, ten, vai_tro, created_at, updated_at FROM tai_khoan WHERE email = ?`, email))
}

func (r *sqliteRepo) TaiKhoanByID(ctx context.Context, id string) (*TaiKhoan, error) {
	return scanTaiKhoan(r.db.QueryRowContext(ctx,
		`SELECT id, email, mat_khau_hash, ten, vai_tro, created_at, updated_at FROM tai_khoan WHERE id = ?`, id))
}

func (r *sqliteRepo) TaoPhienDangNhap(ctx context.Context, taiKhoanID string) (*PhienDangNhap, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))

	now := time.Now().UTC()
	expiresAt := now.Add(24 * time.Hour)
	id := newID("PDN")

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO phien_dang_nhap (id, tai_khoan_id, token, created_at, expires_at) VALUES (?, ?, ?, ?, ?)`,
		id, taiKhoanID, tokenHash, now.Format(time.RFC3339), expiresAt.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("insert session: %w", err)
	}

	return &PhienDangNhap{
		ID:         id,
		TaiKhoanID: taiKhoanID,
		Token:      token,
		CreatedAt:  now,
		ExpiresAt:  expiresAt,
	}, nil
}

func (r *sqliteRepo) PhienDangNhapByToken(ctx context.Context, token string) (*PhienDangNhap, error) {
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))

	var createdAt, expiresAt string
	s := PhienDangNhap{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tai_khoan_id, token, created_at, expires_at FROM phien_dang_nhap WHERE token = ?`, tokenHash,
	).Scan(&s.ID, &s.TaiKhoanID, &s.Token, &createdAt, &expiresAt)

	if err == sql.ErrNoRows {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, fmt.Errorf("query session: %w", err)
	}

	s.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	s.ExpiresAt, _ = time.Parse(time.RFC3339, expiresAt)

	if time.Now().UTC().After(s.ExpiresAt) {
		_, _ = r.db.ExecContext(ctx, "DELETE FROM phien_dang_nhap WHERE id = ?", s.ID)
		return nil, ErrInvalidToken
	}

	return &s, nil
}

func (r *sqliteRepo) XoaPhienDangNhap(ctx context.Context, token string) error {
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
	_, err := r.db.ExecContext(ctx, "DELETE FROM phien_dang_nhap WHERE token = ?", tokenHash)
	return err
}

func (r *sqliteRepo) TaoHOSoBan(ctx context.Context, taiKhoanID string) (*HOSoBan, error) {
	var existing int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ho_so_ban WHERE tai_khoan_id = ?", taiKhoanID).Scan(&existing); err != nil {
		return nil, fmt.Errorf("check existing seller profile: %w", err)
	}
	if existing > 0 {
		return nil, ErrHOSoBanDaTonTai
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id := newID("HS")

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ho_so_ban (id, tai_khoan_id, trang_thai, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		id, taiKhoanID, TrangThaiHOSoBanKichHoat, now, now)
	if err != nil {
		return nil, fmt.Errorf("insert seller profile: %w", err)
	}

	createdAt, _ := time.Parse(time.RFC3339, now)
	return &HOSoBan{
		ID:         id,
		TaiKhoanID: taiKhoanID,
		TrangThai:  TrangThaiHOSoBanKichHoat,
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
	}, nil
}

func (r *sqliteRepo) HOSoBanByTaiKhoanID(ctx context.Context, taiKhoanID string) (*HOSoBan, error) {
	hs := HOSoBan{}
	var createdAt, updatedAt string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tai_khoan_id, trang_thai, created_at, updated_at FROM ho_so_ban WHERE tai_khoan_id = ?`, taiKhoanID,
	).Scan(&hs.ID, &hs.TaiKhoanID, &hs.TrangThai, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrHOSoBanNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query seller profile: %w", err)
	}

	hs.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	hs.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &hs, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// newID generates a unique ID with the given prefix.
func newID(prefix string) string {
	b := make([]byte, 12)
	rand.Read(b)
	return prefix + hex.EncodeToString(b)
}

// scanTaiKhoan scans a single row into a TaiKhoan.
func scanTaiKhoan(row *sql.Row) (*TaiKhoan, error) {
	acc := TaiKhoan{}
	var createdAt, updatedAt string
	err := row.Scan(&acc.ID, &acc.Email, &acc.MatKhauHash, &acc.Ten, &acc.VaiTro, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan account: %w", err)
	}
	acc.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	acc.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &acc, nil
}

// SeedAdmin creates a default admin account if none exists.
// Intended for development/testing bootstrap.
func SeedAdmin(ctx context.Context, repo AccountRepository) error {
	existing, err := repo.TaiKhoanByEmail(ctx, "admin@ecommerce.local")
	if err != nil {
		return nil // DB not available
	}
	if existing != nil {
		return nil // already exists
	}

	_, err = repo.CreateTaiKhoan(ctx, "admin@ecommerce.local", "admin123", "Quản trị viên")
	if err != nil {
		return fmt.Errorf("seed admin: %w", err)
	}

	// Promote to admin role via direct DB update if SQLite
	if sr, ok := repo.(*sqliteRepo); ok {
		_, err = sr.db.Exec("UPDATE tai_khoan SET vai_tro = ? WHERE email = ?", string(VaiTroAdmin), "admin@ecommerce.local")
		if err != nil {
			return fmt.Errorf("promote admin: %w", err)
		}
	}

	return nil
}
