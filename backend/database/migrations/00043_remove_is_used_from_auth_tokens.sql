-- +goose Up
ALTER TABLE auth_tokens DROP COLUMN is_used;

-- +goose Down
ALTER TABLE auth_tokens ADD COLUMN is_used BOOLEAN DEFAULT FALSE;
