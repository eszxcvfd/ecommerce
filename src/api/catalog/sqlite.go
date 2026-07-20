package catalog

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// defaultBusyTimeout is the busy timeout in milliseconds for SQLite connections.
const defaultBusyTimeout = 5000

// OpenSQLite opens a SQLite database at the given path, applies runtime PRAGMAs,
// runs embedded versioned migrations, and returns the *sql.DB handle.
// The returned DB is configured with MaxOpenConns=1 (single writer).
func OpenSQLite(path string) (*sql.DB, error) {
	// Ensure parent directory exists before opening the database file.
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return nil, fmt.Errorf("create database directory %s: %w", parent, err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// Apply runtime PRAGMAs
	pragmas := []struct {
		stmt string
		name string
	}{
		{"PRAGMA foreign_keys = ON", "foreign_keys"},
		{"PRAGMA journal_mode = WAL", "journal_mode"},
		{fmt.Sprintf("PRAGMA busy_timeout = %d", defaultBusyTimeout), "busy_timeout"},
		{"PRAGMA synchronous = NORMAL", "synchronous"},
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p.stmt); err != nil {
			db.Close()
			return nil, fmt.Errorf("set pragma %s: %w", p.name, err)
		}
	}

	// Limit to one writer connection (matches production topology)
	db.SetMaxOpenConns(1)

	// Run embedded versioned migrations via goose
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		db.Close()
		return nil, fmt.Errorf("goose dialect: %w", err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		db.Close()
		return nil, fmt.Errorf("goose up: %w", err)
	}

	return db, nil
}

// VerifySchema checks that all embedded migrations have been applied.
// Returns nil if the schema is up to date, or an error listing pending migrations.
// This is the production counterpart of OpenSQLite's auto-migration: it verifies
// but does NOT run migrations.
func VerifySchema(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("goose dialect: %w", err)
	}

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var maxVersion int64
	for _, entry := range entries {
		v, err := goose.NumericComponent(entry.Name())
		if err == nil && v > maxVersion {
			maxVersion = v
		}
	}

	current, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("get db version: %w", err)
	}

	if current < maxVersion {
		return fmt.Errorf("schema has pending migrations (current: %d, latest: %d)", current, maxVersion)
	}
	return nil
}

// ---------------------------------------------------------------------------
// SQLite repository adapter
// ---------------------------------------------------------------------------

// sqliteRepo implements CatalogRepository backed by a SQLite database.
type sqliteRepo struct {
	db *sql.DB
}

// NewSQLiteRepo creates a CatalogRepository backed by the given SQLite DB.
func NewSQLiteRepo(db *sql.DB) CatalogRepository {
	return &sqliteRepo{db: db}
}

// Products returns only approved (public) products.
func (r *sqliteRepo) Products() ([]SanPhamSo, error) {
	return queryApproved(r.db, CatalogQuery{})
}

// Search returns approved products filtered and sorted by the given query.
func (r *sqliteRepo) Search(query CatalogQuery) ([]SanPhamSo, error) {
	return queryApproved(r.db, query)
}

// queryApproved returns approved products optionally filtered/sorted.
func queryApproved(db *sql.DB, q CatalogQuery) ([]SanPhamSo, error) {
	var conditions []string
	var args []any

	// Always filter by approved status
	conditions = append(conditions, "s.trang_thai = 'approved'")

	// Text search
	if q.Q != "" {
		normalized := normalizeSearch(q.Q)
		// ADR-0001 requires LIKE over normalized columns for accent-insensitive search.
		// Both the stored search columns and the query are lowercased, stripped of
		// combining marks, and have đ→d mapped, so LIKE '%q%' on the concatenated
		// normalized text provides case/accent-insensitive substring matching.
		conditions = append(conditions, "(s.ten_search || ' ' || s.mo_ta_search) LIKE '%' || ? || '%'")
		args = append(args, normalized)
	}

	// danh_muc filter
	if q.DanhMuc != "" {
		conditions = append(conditions, "s.danh_muc = ?")
		args = append(args, q.DanhMuc)
	}

	// dinh_dang filter
	if q.DinhDang != "" {
		conditions = append(conditions, "EXISTS (SELECT 1 FROM san_pham_dinh_dang f WHERE f.san_pham_id = s.id AND f.dinh_dang = ?)")
		args = append(args, q.DinhDang)
	}

	// Price range (inclusive). Free products have so_xu = 0.
	if q.MinXu != nil || q.MaxXu != nil {
		minVal := int64(0)
		maxVal := int64(1<<63 - 1)
		if q.MinXu != nil {
			minVal = *q.MinXu
		}
		if q.MaxXu != nil {
			maxVal = *q.MaxXu
		}
		conditions = append(conditions, "CASE WHEN s.mien_phi = 1 THEN 0 ELSE s.so_xu END BETWEEN ? AND ?")
		args = append(args, minVal, maxVal)
	}

	// Build WHERE clause
	whereClause := strings.Join(conditions, " AND ")

	// Determine ORDER BY
	orderClause := buildOrderBy(q.Sort)

	// Build and execute main product query
	query := fmt.Sprintf(`
		SELECT s.id, s.ten, s.mo_ta, s.anh_demo,
		       s.mien_phi, s.so_xu, s.danh_muc,
		       s.diem_danh_gia, s.so_luong_danh_gia, s.ngay_tao,
		       s.so_luot_tai, s.trang_thai
		FROM san_pham_so s
		WHERE %s
		%s
	`, whereClause, orderClause)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	// Scan all products first (closes rows before loading formats)
	products, ids, err := scanProductIDs(rows)
	if err != nil {
		return nil, fmt.Errorf("scan products: %w", err)
	}
	if len(products) == 0 {
		return products, nil
	}

	// Batch-load all formats for matched product IDs
	formatMap, err := batchLoadDinhDang(db, ids)
	if err != nil {
		return nil, fmt.Errorf("load formats: %w", err)
	}
	for i := range products {
		products[i].DinhDang = formatMap[products[i].ID]
	}

	return products, nil
}

// buildOrderBy returns the ORDER BY clause for the given sort order.
// Default sort is newest first (ngay_tao DESC, id ASC).
func buildOrderBy(sortStr string) string {
	switch SortOrder(sortStr) {
	case SortPopular:
		return "ORDER BY s.so_luot_tai DESC, s.id ASC"
	case SortPriceAsc:
		return "ORDER BY CASE WHEN s.mien_phi = 1 THEN 0 ELSE s.so_xu END ASC, s.id ASC"
	case SortPriceDesc:
		return "ORDER BY CASE WHEN s.mien_phi = 1 THEN 0 ELSE s.so_xu END DESC, s.id ASC"
	case SortRating:
		return "ORDER BY CASE WHEN s.so_luong_danh_gia > 0 THEN 0 ELSE 1 END ASC, s.diem_danh_gia DESC, s.id ASC"
	case SortNewest:
		return "ORDER BY s.ngay_tao DESC, s.id ASC"
	default:
		return "ORDER BY s.ngay_tao DESC, s.id ASC"
	}
}

// scanProductIDs scans all rows and returns products along with their IDs.
// This function does NOT execute any DB queries, avoiding MaxOpenConns deadlocks.
func scanProductIDs(rows *sql.Rows) ([]SanPhamSo, []string, error) {
	var result []SanPhamSo
	var ids []string
	for rows.Next() {
		var sp SanPhamSo
		var ngayTao string
		if err := rows.Scan(
			&sp.ID, &sp.Ten, &sp.MoTa, &sp.AnhDemo,
			&sp.Gia.MienPhi, &sp.Gia.SoXu,
			&sp.DanhMuc, &sp.DiemDanhGia, &sp.SoLuongDanhGia,
			&ngayTao, &sp.SoLuotTai, &sp.TrangThai,
		); err != nil {
			return nil, nil, err
		}
		t, err := time.Parse(time.RFC3339, ngayTao)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05Z07:00", ngayTao)
			if err != nil {
				t = time.Time{}
			}
		}
		sp.NgayTao = t
		result = append(result, sp)
		ids = append(ids, sp.ID)
	}
	return result, ids, rows.Err()
}

// batchLoadDinhDang loads all product formats for the given product IDs in one query.
func batchLoadDinhDang(db *sql.DB, ids []string) (map[string][]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	// Build placeholders for IN clause
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(
		"SELECT san_pham_id, dinh_dang FROM san_pham_dinh_dang WHERE san_pham_id IN (%s) ORDER BY san_pham_id, dinh_dang",
		strings.Join(placeholders, ","),
	)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]string)
	for rows.Next() {
		var pid, ext string
		if err := rows.Scan(&pid, &ext); err != nil {
			return nil, err
		}
		result[pid] = append(result[pid], ext)
	}
	return result, rows.Err()
}

// OpenSQLiteProd opens a SQLite database with production settings.
// It applies runtime PRAGMAs, restricts to a single writer, and verifies
// that all embedded migrations have been applied (without running them).
// This is the production counterpart of OpenSQLite (which auto-migrates).
func OpenSQLiteProd(path string) (*sql.DB, error) {
	// Ensure parent directory exists before opening the database file.
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return nil, fmt.Errorf("create database directory %s: %w", parent, err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	pragmas := []struct {
		stmt string
		name string
	}{
		{"PRAGMA foreign_keys = ON", "foreign_keys"},
		{"PRAGMA journal_mode = WAL", "journal_mode"},
		{fmt.Sprintf("PRAGMA busy_timeout = %d", defaultBusyTimeout), "busy_timeout"},
		{"PRAGMA synchronous = NORMAL", "synchronous"},
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p.stmt); err != nil {
			db.Close()
			return nil, fmt.Errorf("set pragma %s: %w", p.name, err)
		}
	}
	db.SetMaxOpenConns(1)

	if err := VerifySchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("verify schema: %w", err)
	}

	return db, nil
}
