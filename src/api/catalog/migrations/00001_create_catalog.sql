-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS san_pham_so (
    id              TEXT PRIMARY KEY,
    ten             TEXT NOT NULL,
    mo_ta           TEXT NOT NULL DEFAULT '',
    anh_demo        TEXT NOT NULL DEFAULT '',
    mien_phi        INTEGER NOT NULL DEFAULT 0,  -- 0=false, 1=true (SQLite lacks BOOLEAN)
    so_xu           INTEGER NOT NULL DEFAULT 0,
    danh_muc        TEXT NOT NULL,
    diem_danh_gia   REAL NOT NULL DEFAULT 0.0,
    so_luong_danh_gia INTEGER NOT NULL DEFAULT 0,
    ngay_tao        TEXT NOT NULL,  -- ISO 8601 / RFC3339
    so_luot_tai     INTEGER NOT NULL DEFAULT 0,
    trang_thai      TEXT NOT NULL DEFAULT 'draft',
    ten_search      TEXT NOT NULL DEFAULT '',  -- normalized for accent-insensitive search
    mo_ta_search    TEXT NOT NULL DEFAULT '',  -- normalized for accent-insensitive search
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS san_pham_dinh_dang (
    san_pham_id     TEXT NOT NULL REFERENCES san_pham_so(id) ON DELETE CASCADE,
    dinh_dang       TEXT NOT NULL,
    PRIMARY KEY (san_pham_id, dinh_dang)
);

-- Targeted indexes
CREATE INDEX IF NOT EXISTS idx_san_pham_trang_thai ON san_pham_so(trang_thai);
CREATE INDEX IF NOT EXISTS idx_san_pham_danh_muc ON san_pham_so(danh_muc);
CREATE INDEX IF NOT EXISTS idx_san_pham_so_xu ON san_pham_so(so_xu);
CREATE INDEX IF NOT EXISTS idx_san_pham_ngay_tao ON san_pham_so(ngay_tao);
CREATE INDEX IF NOT EXISTS idx_san_pham_so_luot_tai ON san_pham_so(so_luot_tai);
CREATE INDEX IF NOT EXISTS idx_san_pham_diem_danh_gia ON san_pham_so(diem_danh_gia);
CREATE INDEX IF NOT EXISTS idx_san_pham_dinh_dang ON san_pham_dinh_dang(dinh_dang);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS san_pham_dinh_dang;
DROP TABLE IF EXISTS san_pham_so;

-- +goose StatementEnd
