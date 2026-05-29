-- +goose Up
ALTER TABLE students ADD COLUMN IF NOT EXISTS description TEXT;

-- +goose Down
ALTER TABLE students DROP COLUMN IF EXISTS description;
