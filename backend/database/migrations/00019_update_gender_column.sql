-- +goose Up
-- Change gender column from ENUM to VARCHAR to allow full text "Laki-laki" and "Perempuan"
ALTER TABLE students MODIFY COLUMN gender VARCHAR(20) NOT NULL;

-- Update existing records to full text (optional, but good for consistency)
UPDATE students SET gender = 'Laki-laki' WHERE gender = 'L';
UPDATE students SET gender = 'Perempuan' WHERE gender = 'P';

-- +goose Down
-- Revert back to ENUM if needed (Warning: this will fail if values are not L or P)
UPDATE students SET gender = 'L' WHERE gender = 'Laki-laki';
UPDATE students SET gender = 'P' WHERE gender = 'Perempuan';
ALTER TABLE students MODIFY COLUMN gender ENUM('L', 'P') NOT NULL;
