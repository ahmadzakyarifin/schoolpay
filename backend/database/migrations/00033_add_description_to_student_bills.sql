-- +goose Up
ALTER TABLE student_bills
  ADD COLUMN IF NOT EXISTS description TEXT NULL AFTER end_date;

-- +goose Down
ALTER TABLE student_bills
  DROP COLUMN IF EXISTS description;
