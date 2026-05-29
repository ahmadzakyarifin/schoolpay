-- +goose Up
CREATE TABLE IF NOT EXISTS academic_year_classes (
    academic_year_id INT UNSIGNED NOT NULL,
    class_id INT UNSIGNED NOT NULL,
    PRIMARY KEY (academic_year_id, class_id),
    CONSTRAINT fk_ayc_ay FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE CASCADE,
    CONSTRAINT fk_ayc_class FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE
);

-- Allow classes to exist without a year initially (templates)
ALTER TABLE classes MODIFY academic_year_id INT UNSIGNED NULL;

-- +goose Down
ALTER TABLE classes MODIFY academic_year_id INT UNSIGNED NOT NULL;
DROP TABLE IF EXISTS academic_year_classes;
