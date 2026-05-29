-- +goose Up
ALTER TABLE students 
ADD COLUMN email VARCHAR(255) AFTER entry_year,
MODIFY COLUMN phone_number VARCHAR(20) NOT NULL,
ADD UNIQUE INDEX idx_student_email (email),
ADD UNIQUE INDEX idx_student_phone (phone_number);

-- Drop foreign key because we want to make parent_id NOT NULL
ALTER TABLE students DROP FOREIGN KEY fk_student_parent;

-- Modify columns: parent_id NOT NULL, address NULL (optional)
ALTER TABLE students 
MODIFY COLUMN parent_id INT UNSIGNED NOT NULL,
MODIFY COLUMN address TEXT NULL;

-- Re-add foreign key without SET NULL
ALTER TABLE students ADD CONSTRAINT fk_student_parent FOREIGN KEY (parent_id) REFERENCES users(id);

-- +goose Down
ALTER TABLE students DROP FOREIGN KEY fk_student_parent;
ALTER TABLE students MODIFY COLUMN parent_id INT UNSIGNED, MODIFY COLUMN address TEXT NOT NULL;
ALTER TABLE students ADD CONSTRAINT fk_student_parent FOREIGN KEY (parent_id) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE students 
DROP INDEX idx_student_email,
DROP INDEX idx_student_phone,
DROP COLUMN email;
