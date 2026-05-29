-- +goose Up
-- 1. Add academic_year_id to classes table
ALTER TABLE classes 
ADD COLUMN academic_year_id INT UNSIGNED AFTER major_id,
ADD CONSTRAINT fk_class_academic_year FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE SET NULL;

-- 2. Create join table for academic year majors
CREATE TABLE IF NOT EXISTS academic_year_majors (
    academic_year_id INT UNSIGNED NOT NULL,
    major_id INT UNSIGNED NOT NULL,
    PRIMARY KEY (academic_year_id, major_id),
    CONSTRAINT fk_aym_academic_year FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE CASCADE,
    CONSTRAINT fk_aym_major FOREIGN KEY (major_id) REFERENCES majors(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 3. Add unique index to prevent duplicate class names in the same year and major
ALTER TABLE classes ADD UNIQUE INDEX idx_unique_class_in_year (academic_year_id, name, major_id);

-- +goose Down
ALTER TABLE classes DROP INDEX idx_unique_class_in_year;
DROP TABLE IF EXISTS academic_year_majors;
ALTER TABLE classes 
DROP FOREIGN KEY fk_class_academic_year,
DROP COLUMN academic_year_id;
