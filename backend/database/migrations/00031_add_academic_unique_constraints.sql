-- +goose Up
-- Add unique constraint for major name
CREATE UNIQUE INDEX idx_majors_name ON majors (name);
-- Add unique constraint for academic year
CREATE UNIQUE INDEX idx_academic_years_year ON academic_years (year);

-- +goose Down
DROP INDEX idx_majors_name ON majors;
DROP INDEX idx_academic_years_year ON academic_years;
