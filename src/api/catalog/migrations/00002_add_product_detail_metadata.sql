-- +goose Up
-- +goose StatementBegin

ALTER TABLE san_pham_so ADD COLUMN mo_ta_chi_tiet TEXT NOT NULL DEFAULT '';
ALTER TABLE san_pham_so ADD COLUMN giay_phep TEXT NOT NULL DEFAULT '';
ALTER TABLE san_pham_so ADD COLUMN nguoi_ban_hien_thi TEXT NOT NULL DEFAULT '';
ALTER TABLE san_pham_so ADD COLUMN ngay_dang TEXT NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS san_pham_tep (
    san_pham_id     TEXT NOT NULL REFERENCES san_pham_so(id) ON DELETE CASCADE,
    ten_tep         TEXT NOT NULL,
    dinh_dang       TEXT NOT NULL,
    dung_luong_bytes INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (san_pham_id, ten_tep, dinh_dang)
);

CREATE INDEX IF NOT EXISTS idx_san_pham_tep_id ON san_pham_tep(san_pham_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS san_pham_tep;

ALTER TABLE san_pham_so DROP COLUMN IF EXISTS ngay_dang;
ALTER TABLE san_pham_so DROP COLUMN IF EXISTS nguoi_ban_hien_thi;
ALTER TABLE san_pham_so DROP COLUMN IF EXISTS giay_phep;
ALTER TABLE san_pham_so DROP COLUMN IF EXISTS mo_ta_chi_tiet;

-- +goose StatementEnd
