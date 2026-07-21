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

// ProductByID returns one approved product by ID, or nil if not found or not public.
func (r *sqliteRepo) ProductByID(id string) (*SanPhamSo, error) {
	return queryProductApproved(r.db, id)
}

// ProductsByCategory returns approved products in the given category, excluding one ID.
// Ordered by publish date descending (newest first) with ID tie-break, limited to max.
func (r *sqliteRepo) ProductsByCategory(category DanhMuc, excludeID string, max int) ([]SanPhamSo, error) {
	return queryRecommendations(r.db, category, excludeID, max)
}

func (r *sqliteRepo) DraftsBySeller(sellerID string) ([]SanPhamSo, error) {
	return queryDraftsBySeller(r.db, sellerID)
}

func (r *sqliteRepo) DraftByID(id, sellerID string) (*SanPhamSo, error) {
	return queryDraftByID(r.db, id, sellerID)
}

func (r *sqliteRepo) CreateDraft(input DraftInput, sellerID string) (*SanPhamSo, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	id := newProductID()
	now := time.Now().UTC().Format(time.RFC3339)

	_, err = tx.Exec(`
		INSERT INTO san_pham_so (id, ten, mo_ta, mo_ta_chi_tiet, anh_demo, mien_phi, so_xu, danh_muc,
		                         ngay_tao, trang_thai, giay_phep, nguoi_ban_hien_thi, nguoi_ban_id,
		                         ten_search, mo_ta_search, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'draft', ?, '', ?, '', ?, datetime('now'), datetime('now'))
	`, id, input.Ten, input.MoTa, input.MoTaChiTiet, input.AnhDemo,
		boolToInt(input.MienPhi), input.SoXu, string(input.DanhMuc),
		now, input.GiayPhep, sellerID,
		normalizeSearch(input.Ten), normalizeSearch(input.MoTa))
	if err != nil {
		return nil, fmt.Errorf("insert draft: %w", err)
	}

	// Insert file entries
	formatSet := make(map[string]bool)
	for _, tep := range input.Tep {
		_, err = tx.Exec(`
			INSERT INTO san_pham_tep (san_pham_id, ten_tep, dinh_dang, dung_luong_bytes)
			VALUES (?, ?, ?, ?)
		`, id, tep.TenTep, tep.DinhDang, tep.DungLuongBytes)
		if err != nil {
			return nil, fmt.Errorf("insert file %s: %w", tep.TenTep, err)
		}
		if !formatSet[tep.DinhDang] {
			formatSet[tep.DinhDang] = true
			_, err = tx.Exec(`
				INSERT INTO san_pham_dinh_dang (san_pham_id, dinh_dang)
				VALUES (?, ?)
			`, id, tep.DinhDang)
			if err != nil {
				return nil, fmt.Errorf("insert format %s: %w", tep.DinhDang, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit draft: %w", err)
	}

	return queryDraftByID(r.db, id, sellerID)
}

func (r *sqliteRepo) UpdateDraft(id, sellerID string, input DraftUpdateInput) (*SanPhamSo, error) {
	// First verify ownership
	existing, err := queryDraftByID(r.db, id, sellerID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var sets []string
	var args []any

	if input.Ten != nil {
		sets = append(sets, "ten = ?", "ten_search = ?")
		args = append(args, *input.Ten, normalizeSearch(*input.Ten))
	}
	if input.MoTa != nil {
		sets = append(sets, "mo_ta = ?", "mo_ta_search = ?")
		args = append(args, *input.MoTa, normalizeSearch(*input.MoTa))
	}
	if input.MoTaChiTiet != nil {
		sets = append(sets, "mo_ta_chi_tiet = ?")
		args = append(args, *input.MoTaChiTiet)
	}
	if input.AnhDemo != nil {
		sets = append(sets, "anh_demo = ?")
		args = append(args, *input.AnhDemo)
	}
	if input.MienPhi != nil {
		sets = append(sets, "mien_phi = ?")
		args = append(args, boolToInt(*input.MienPhi))
	}
	if input.SoXu != nil {
		sets = append(sets, "so_xu = ?")
		args = append(args, *input.SoXu)
	}
	if input.DanhMuc != nil {
		sets = append(sets, "danh_muc = ?")
		args = append(args, string(*input.DanhMuc))
	}
	if input.GiayPhep != nil {
		sets = append(sets, "giay_phep = ?")
		args = append(args, *input.GiayPhep)
	}

	if len(sets) > 0 {
		sets = append(sets, "updated_at = datetime('now')")
		q := fmt.Sprintf("UPDATE san_pham_so SET %s WHERE id = ? AND nguoi_ban_id = ? AND trang_thai = 'draft'",
			strings.Join(sets, ", "))
		args = append(args, id, sellerID)
		_, err = tx.Exec(q, args...)
		if err != nil {
			return nil, fmt.Errorf("update draft: %w", err)
		}
	}

	// Replace files if provided
	if input.Tep != nil {
		_, err = tx.Exec("DELETE FROM san_pham_tep WHERE san_pham_id = ?", id)
		if err != nil {
			return nil, fmt.Errorf("delete old files: %w", err)
		}
		_, err = tx.Exec("DELETE FROM san_pham_dinh_dang WHERE san_pham_id = ?", id)
		if err != nil {
			return nil, fmt.Errorf("delete old formats: %w", err)
		}

		formatSet := make(map[string]bool)
		for _, tep := range input.Tep {
			_, err = tx.Exec(`
				INSERT INTO san_pham_tep (san_pham_id, ten_tep, dinh_dang, dung_luong_bytes)
				VALUES (?, ?, ?, ?)
			`, id, tep.TenTep, tep.DinhDang, tep.DungLuongBytes)
			if err != nil {
				return nil, fmt.Errorf("insert file: %w", err)
			}
			if !formatSet[tep.DinhDang] {
				formatSet[tep.DinhDang] = true
				_, err = tx.Exec(`
					INSERT INTO san_pham_dinh_dang (san_pham_id, dinh_dang)
					VALUES (?, ?)
				`, id, tep.DinhDang)
				if err != nil {
					return nil, fmt.Errorf("insert format: %w", err)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit update: %w", err)
	}

	return queryDraftByID(r.db, id, sellerID)
}

func (r *sqliteRepo) DeleteDraft(id, sellerID string) error {
	result, err := r.db.Exec(
		"DELETE FROM san_pham_so WHERE id = ? AND nguoi_ban_id = ? AND trang_thai = 'draft'",
		id, sellerID,
	)
	if err != nil {
		return fmt.Errorf("delete draft: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("không tìm thấy bản nháp")
	}
	return nil
}

// queryDraftsBySeller returns all drafts owned by the given seller.
func queryDraftsBySeller(db *sql.DB, sellerID string) ([]SanPhamSo, error) {
	rows, err := db.Query(`
		SELECT s.id, s.ten, s.mo_ta, s.mo_ta_chi_tiet, s.anh_demo,
		       s.mien_phi, s.so_xu, s.danh_muc,
		       s.diem_danh_gia, s.so_luong_danh_gia, s.ngay_tao,
		       s.so_luot_tai, s.trang_thai,
		       s.giay_phep, s.nguoi_ban_hien_thi, s.ngay_dang,
		       s.nguoi_ban_id
		FROM san_pham_so s
		WHERE s.nguoi_ban_id = ? AND s.trang_thai = 'draft'
		ORDER BY s.ngay_tao DESC
	`, sellerID)
	if err != nil {
		return nil, fmt.Errorf("query drafts: %w", err)
	}
	defer rows.Close()

	return scanDrafts(rows, db)
}

// queryDraftByID returns one draft by ID, scoped to the seller.
func queryDraftByID(db *sql.DB, id, sellerID string) (*SanPhamSo, error) {
	var sp SanPhamSo
	var ngayTao string
	var ngayDang string
	err := db.QueryRow(`
		SELECT s.id, s.ten, s.mo_ta, s.mo_ta_chi_tiet, s.anh_demo,
		       s.mien_phi, s.so_xu, s.danh_muc,
		       s.diem_danh_gia, s.so_luong_danh_gia, s.ngay_tao,
		       s.so_luot_tai, s.trang_thai,
		       s.giay_phep, s.nguoi_ban_hien_thi, s.ngay_dang,
		       s.nguoi_ban_id
		FROM san_pham_so s
		WHERE s.id = ? AND s.nguoi_ban_id = ? AND s.trang_thai = 'draft'
	`, id, sellerID).Scan(
		&sp.ID, &sp.Ten, &sp.MoTa, &sp.MoTaChiTiet, &sp.AnhDemo,
		&sp.Gia.MienPhi, &sp.Gia.SoXu,
		&sp.DanhMuc, &sp.DiemDanhGia, &sp.SoLuongDanhGia,
		&ngayTao, &sp.SoLuotTai, &sp.TrangThai,
		&sp.GiayPhep, &sp.NguoiBanHienThi, &ngayDang,
		&sp.NguoiBanID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query draft %s: %w", id, err)
	}
	t, _ := parseTime(ngayTao)
	sp.NgayTao = t
	if ngayDang != "" {
		sp.NgayDang, _ = parseTime(ngayDang)
	}

	// Load formats
	formatMap, _ := batchLoadDinhDang(db, []string{id})
	sp.DinhDang = formatMap[id]

	// Load files
	files, _ := loadTep(db, id)
	sp.Tep = files

	return &sp, nil
}

// scanDrafts scans draft rows and loads format/file data.
func scanDrafts(rows *sql.Rows, db *sql.DB) ([]SanPhamSo, error) {
	var products []SanPhamSo
	var ids []string
	for rows.Next() {
		var sp SanPhamSo
		var ngayTao string
		var ngayDang string
		if err := rows.Scan(
			&sp.ID, &sp.Ten, &sp.MoTa, &sp.MoTaChiTiet, &sp.AnhDemo,
			&sp.Gia.MienPhi, &sp.Gia.SoXu,
			&sp.DanhMuc, &sp.DiemDanhGia, &sp.SoLuongDanhGia,
			&ngayTao, &sp.SoLuotTai, &sp.TrangThai,
			&sp.GiayPhep, &sp.NguoiBanHienThi, &ngayDang,
			&sp.NguoiBanID,
		); err != nil {
			return nil, fmt.Errorf("scan draft: %w", err)
		}
		t, _ := parseTime(ngayTao)
		sp.NgayTao = t
		if ngayDang != "" {
			sp.NgayDang, _ = parseTime(ngayDang)
		}
		products = append(products, sp)
		ids = append(ids, sp.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Batch-load formats
	if len(ids) > 0 {
		formatMap, err := batchLoadDinhDang(db, ids)
		if err == nil {
			for i := range products {
				products[i].DinhDang = formatMap[products[i].ID]
			}
		}
		// Load files per product
		for i := range products {
			files, _ := loadTep(db, products[i].ID)
			products[i].Tep = files
		}
	}

	return products, nil
}

// newProductID generates a unique product ID for drafts.
func newProductID() string {
	return fmt.Sprintf("sp_draft_%d", time.Now().UnixNano())
}

// queryProductApproved fetches a single approved product by ID.
// Returns nil, nil if not found or not approved.
func queryProductApproved(db *sql.DB, id string) (*SanPhamSo, error) {
	var sp SanPhamSo
	var ngayTao string
	var ngayDang string
	err := db.QueryRow(`
		SELECT s.id, s.ten, s.mo_ta, s.mo_ta_chi_tiet, s.anh_demo,
		       s.mien_phi, s.so_xu, s.danh_muc,
		       s.diem_danh_gia, s.so_luong_danh_gia, s.ngay_tao,
		       s.so_luot_tai, s.trang_thai,
		       s.giay_phep, s.nguoi_ban_hien_thi, s.ngay_dang
		FROM san_pham_so s
		WHERE s.id = ? AND s.trang_thai = 'approved'
	`, id).Scan(
		&sp.ID, &sp.Ten, &sp.MoTa, &sp.MoTaChiTiet, &sp.AnhDemo,
		&sp.Gia.MienPhi, &sp.Gia.SoXu,
		&sp.DanhMuc, &sp.DiemDanhGia, &sp.SoLuongDanhGia,
		&ngayTao, &sp.SoLuotTai, &sp.TrangThai,
		&sp.GiayPhep, &sp.NguoiBanHienThi, &ngayDang,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query product %s: %w", id, err)
	}
	t, err := parseTime(ngayTao)
	if err != nil {
		t = time.Time{}
	}
	sp.NgayTao = t

	if ngayDang != "" {
		sp.NgayDang, _ = parseTime(ngayDang)
	}

	// Load formats
	formatMap, err := batchLoadDinhDang(db, []string{id})
	if err != nil {
		return nil, fmt.Errorf("load formats for %s: %w", id, err)
	}
	sp.DinhDang = formatMap[id]

	// Load files
	files, err := loadTep(db, id)
	if err != nil {
		return nil, fmt.Errorf("load files for %s: %w", id, err)
	}
	sp.Tep = files

	return &sp, nil
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
		SELECT s.id, s.ten, s.mo_ta, s.mo_ta_chi_tiet, s.anh_demo,
		       s.mien_phi, s.so_xu, s.danh_muc,
		       s.diem_danh_gia, s.so_luong_danh_gia, s.ngay_tao,
		       s.so_luot_tai, s.trang_thai,
		       s.giay_phep, s.nguoi_ban_hien_thi, s.ngay_dang
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
		var ngayDang string
		if err := rows.Scan(
			&sp.ID, &sp.Ten, &sp.MoTa, &sp.MoTaChiTiet, &sp.AnhDemo,
			&sp.Gia.MienPhi, &sp.Gia.SoXu,
			&sp.DanhMuc, &sp.DiemDanhGia, &sp.SoLuongDanhGia,
			&ngayTao, &sp.SoLuotTai, &sp.TrangThai,
			&sp.GiayPhep, &sp.NguoiBanHienThi, &ngayDang,
		); err != nil {
			return nil, nil, err
		}
		t, err := parseTime(ngayTao)
		if err != nil {
			t = time.Time{}
		}
		sp.NgayTao = t

		if ngayDang != "" {
			sp.NgayDang, _ = parseTime(ngayDang)
		}

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

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	rows, err := db.Query(
		"SELECT san_pham_id, dinh_dang FROM san_pham_dinh_dang WHERE san_pham_id IN ("+strings.Join(placeholders, ",")+") ORDER BY san_pham_id, dinh_dang",
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("query formats: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]string, len(ids))
	for rows.Next() {
		var pid, ext string
		if err := rows.Scan(&pid, &ext); err != nil {
			return nil, fmt.Errorf("scan format: %w", err)
		}
		result[pid] = append(result[pid], ext)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iter: %w", err)
	}

	return result, nil
}

// loadTep loads all file entries for a given product ID.
func loadTep(db *sql.DB, productID string) ([]Tep, error) {
	rows, err := db.Query(
		`SELECT ten_tep, dinh_dang, dung_luong_bytes FROM san_pham_tep WHERE san_pham_id = ? ORDER BY ten_tep`,
		productID,
	)
	if err != nil {
		return nil, fmt.Errorf("query files for %s: %w", productID, err)
	}
	defer rows.Close()

	var result []Tep
	for rows.Next() {
		var f Tep
		if err := rows.Scan(&f.TenTep, &f.DinhDang, &f.DungLuongBytes); err != nil {
			return nil, fmt.Errorf("scan file: %w", err)
		}
		result = append(result, f)
	}
	return result, rows.Err()
}

// queryRecommendations returns approved products in the given category, excluding one ID,
// ordered by publish date descending with ID tie-break, limited to max.
func queryRecommendations(db *sql.DB, category DanhMuc, excludeID string, max int) ([]SanPhamSo, error) {
	query := fmt.Sprintf(`
		SELECT s.id, s.ten, s.mo_ta, s.mo_ta_chi_tiet, s.anh_demo,
		       s.mien_phi, s.so_xu, s.danh_muc,
		       s.diem_danh_gia, s.so_luong_danh_gia, s.ngay_tao,
		       s.so_luot_tai, s.trang_thai,
		       s.giay_phep, s.nguoi_ban_hien_thi, s.ngay_dang
		FROM san_pham_so s
		WHERE s.trang_thai = 'approved'
		  AND s.danh_muc = ?
		  AND s.id != ?
		ORDER BY s.ngay_dang DESC, s.id ASC
		LIMIT ?
	`)

	rows, err := db.Query(query, string(category), excludeID, max)
	if err != nil {
		return nil, fmt.Errorf("query recommendations: %w", err)
	}
	defer rows.Close()

	products, ids, err := scanProductIDs(rows)
	if err != nil {
		return nil, fmt.Errorf("scan recommendations: %w", err)
	}
	if len(products) == 0 {
		return products, nil
	}

	// Batch-load formats
	formatMap, err := batchLoadDinhDang(db, ids)
	if err != nil {
		return nil, fmt.Errorf("load recommendation formats: %w", err)
	}
	for i := range products {
		products[i].DinhDang = formatMap[products[i].ID]
	}

	// Load files for each product
	for i := range products {
		files, err := loadTep(db, products[i].ID)
		if err != nil {
			return nil, fmt.Errorf("load files for %s: %w", products[i].ID, err)
		}
		products[i].Tep = files
	}

	return products, nil
}

// parseTime tries to parse a timestamp string, attempting RFC3339 first.
func parseTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02T15:04:05Z07:00", s)
}

// OpenSQLiteProd opens a SQLite database with production settings.
// It applies runtime PRAGMAs, restricts to a single writer, and verifies
// that all embedded migrations have been applied (without running them).
// This is the production counterpart of OpenSQLite (which auto-migrates).
func OpenSQLiteProd(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite %s: %w", path, err)
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

	// Verify all migrations have been applied (verify only, no auto-migration)
	if err := VerifySchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("verify schema: %w", err)
	}
	return db, nil
}
