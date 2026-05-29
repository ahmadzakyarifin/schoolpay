-- +goose Up
UPDATE majors SET code = NULL WHERE code IS NOT NULL AND TRIM(code) = '';
CREATE INDEX idx_majors_code ON majors (code);

-- +goose Down
DROP INDEX idx_majors_code ON majors;
