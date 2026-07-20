package catalog

import (
	"database/sql"
	"fmt"
	"time"
)

// ImportCatalogJSON validates a versioned catalog JSON file and imports all
// products into the given database in a single transaction (all-or-nothing).
//
// When allowDuplicates is false, it rejects:
//   - Duplicate IDs inside the input JSON.
//   - IDs that already exist in the database.
//
// When allowDuplicates is true, conflicting rows are silently skipped (INSERT OR IGNORE).
func ImportCatalogJSON(db *sql.DB, data []byte, allowDuplicates bool) error {
	cf, err := ValidateCatalogJSON(data)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if !allowDuplicates {
		// Check for duplicate IDs within the input.
		seen := make(map[string]int, len(cf.Products))
		for _, p := range cf.Products {
			if _, dup := seen[p.ID]; dup {
				return fmt.Errorf("duplicate product ID %q in input", p.ID)
			}
			seen[p.ID] = 1
		}

		// Check for conflicts with existing rows.
		for _, p := range cf.Products {
			var exists int
			if err := db.QueryRow("SELECT COUNT(*) FROM san_pham_so WHERE id = ?", p.ID).Scan(&exists); err != nil {
				return fmt.Errorf("check id %s: %w", p.ID, err)
			}
			if exists > 0 {
				return fmt.Errorf("product ID %q already exists in database", p.ID)
			}
		}
	}

	// Convert and insert in a single transaction.
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	for _, cp := range cf.Products {
		sp, err := cp.ToSanPhamSo()
		if err != nil {
			return fmt.Errorf("convert %s: %w", cp.ID, err)
		}

		if allowDuplicates {
			_, err = tx.Exec(`
				INSERT OR IGNORE INTO san_pham_so (id, ten, mo_ta, anh_demo, mien_phi, so_xu, danh_muc,
				                         diem_danh_gia, so_luong_danh_gia, ngay_tao, so_luot_tai, trang_thai,
				                         ten_search, mo_ta_search)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, sp.ID, sp.Ten, sp.MoTa, sp.AnhDemo, boolToInt(sp.Gia.MienPhi), sp.Gia.SoXu,
				string(sp.DanhMuc), sp.DiemDanhGia, sp.SoLuongDanhGia,
				sp.NgayTao.Format(time.RFC3339), sp.SoLuotTai, string(sp.TrangThai),
				normalizeSearch(sp.Ten), normalizeSearch(sp.MoTa),
			)
		} else {
			_, err = tx.Exec(`
				INSERT INTO san_pham_so (id, ten, mo_ta, anh_demo, mien_phi, so_xu, danh_muc,
				                         diem_danh_gia, so_luong_danh_gia, ngay_tao, so_luot_tai, trang_thai,
				                         ten_search, mo_ta_search)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, sp.ID, sp.Ten, sp.MoTa, sp.AnhDemo, boolToInt(sp.Gia.MienPhi), sp.Gia.SoXu,
				string(sp.DanhMuc), sp.DiemDanhGia, sp.SoLuongDanhGia,
				sp.NgayTao.Format(time.RFC3339), sp.SoLuotTai, string(sp.TrangThai),
				normalizeSearch(sp.Ten), normalizeSearch(sp.MoTa),
			)
		}
		if err != nil {
			return fmt.Errorf("import %s: %w", sp.ID, err)
		}

		for _, ext := range sp.DinhDang {
			_, err := tx.Exec(
				"INSERT OR IGNORE INTO san_pham_dinh_dang (san_pham_id, dinh_dang) VALUES (?, ?)",
				sp.ID, ext,
			)
			if err != nil {
				return fmt.Errorf("import format %s for %s: %w", ext, sp.ID, err)
			}
		}
	}

	return tx.Commit()
}
