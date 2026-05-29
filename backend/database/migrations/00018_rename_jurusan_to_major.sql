-- +goose Up
-- Rename column in classes table
ALTER TABLE classes CHANGE jurusan_id major_id INT UNSIGNED;
ALTER TABLE classes DROP FOREIGN KEY fk_class_jurusan;
ALTER TABLE classes ADD CONSTRAINT fk_class_major FOREIGN KEY (major_id) REFERENCES majors(id) ON DELETE SET NULL;

-- Rename column in students table
ALTER TABLE students CHANGE jurusan_id major_id INT UNSIGNED;

-- +goose Down
-- Revert column in classes table
ALTER TABLE classes CHANGE major_id jurusan_id INT UNSIGNED;
ALTER TABLE classes DROP FOREIGN KEY fk_class_major;
ALTER TABLE classes ADD CONSTRAINT fk_class_jurusan FOREIGN KEY (jurusan_id) REFERENCES majors(id) ON DELETE SET NULL;

-- Revert column in students table
ALTER TABLE students CHANGE major_id jurusan_id INT UNSIGNED;
