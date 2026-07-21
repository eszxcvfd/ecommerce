-- +goose Up
-- +goose StatementBegin

ALTER TABLE san_pham_so ADD COLUMN nguoi_ban_id TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_san_pham_nguoi_ban ON san_pham_so(nguoi_ban_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_san_pham_nguoi_ban;
ALTER TABLE san_pham_so DROP COLUMN IF EXISTS nguoi_ban_id;

-- +goose StatementEnd
