-- +goose Up
ALTER TABLE payments ADD COLUMN IF NOT EXISTS is_bypass_rule BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE payments DROP COLUMN IF EXISTS is_bypass_rule;
