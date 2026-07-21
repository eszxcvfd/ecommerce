-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS tai_khoan (
    id              TEXT PRIMARY KEY,
    email           TEXT NOT NULL UNIQUE,
    mat_khau_hash   TEXT NOT NULL,
    ten             TEXT NOT NULL DEFAULT '',
    vai_tro         TEXT NOT NULL DEFAULT 'nguoi_mua',
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tai_khoan_email ON tai_khoan(email);

CREATE TABLE IF NOT EXISTS phien_dang_nhap (
    id              TEXT PRIMARY KEY,
    tai_khoan_id    TEXT NOT NULL REFERENCES tai_khoan(id) ON DELETE CASCADE,
    token           TEXT NOT NULL,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    expires_at      TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_phien_dang_nhap_token ON phien_dang_nhap(token);
CREATE INDEX IF NOT EXISTS idx_phien_dang_nhap_tai_khoan ON phien_dang_nhap(tai_khoan_id);

CREATE TABLE IF NOT EXISTS ho_so_ban (
    id              TEXT PRIMARY KEY,
    tai_khoan_id    TEXT NOT NULL UNIQUE REFERENCES tai_khoan(id) ON DELETE CASCADE,
    trang_thai      TEXT NOT NULL DEFAULT 'kich_hoat',
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_ho_so_ban_tai_khoan ON ho_so_ban(tai_khoan_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS ho_so_ban;
DROP TABLE IF EXISTS phien_dang_nhap;
DROP TABLE IF EXISTS tai_khoan;

-- +goose StatementEnd
