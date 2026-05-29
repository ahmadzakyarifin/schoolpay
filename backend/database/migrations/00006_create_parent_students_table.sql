-- +goose Up
CREATE TABLE IF NOT EXISTS parent_students (
    parent_id INT UNSIGNED NOT NULL,
    student_id INT UNSIGNED NOT NULL,
    relation VARCHAR(50),
    PRIMARY KEY (parent_id, student_id),
    CONSTRAINT fk_ps_parent FOREIGN KEY (parent_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_ps_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS parent_students;
