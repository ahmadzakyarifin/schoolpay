-- +goose Up
-- Clean up any inconsistent data first (if any)
UPDATE students SET status = 'active' WHERE status = 'aktif';
UPDATE students SET status = 'graduated' WHERE status LIKE 'lulus';
UPDATE students SET status = 'inactive' WHERE status IN ('keluar', 'non-aktif');

-- Change column to ENUM
ALTER TABLE students 
MODIFY COLUMN status ENUM('active', 'inactive', 'graduated') NOT NULL DEFAULT 'active';

-- +goose Down
ALTER TABLE students 
MODIFY COLUMN status VARCHAR(50) NOT NULL DEFAULT 'active';
