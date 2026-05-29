-- +goose Up
CREATE TABLE IF NOT EXISTS students (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    nik VARCHAR(20) NOT NULL UNIQUE,
    nis VARCHAR(50),
    nisn VARCHAR(50) NOT NULL UNIQUE,
    class_id INT UNSIGNED NOT NULL,
    jurusan_id INT UNSIGNED NOT NULL,
    parent_id INT UNSIGNED,
    name VARCHAR(255) NOT NULL,
    image_path TEXT,
    gender ENUM('L', 'P') NOT NULL,
    birth_place VARCHAR(100) NOT NULL,
    birth_date DATE NOT NULL,
    religion VARCHAR(50) NOT NULL,
    rt VARCHAR(10) NOT NULL,
    rw VARCHAR(10) NOT NULL,
    village VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    entry_year INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_student_class FOREIGN KEY (class_id) REFERENCES classes(id),
    CONSTRAINT fk_student_jurusan FOREIGN KEY (jurusan_id) REFERENCES majors(id),
    CONSTRAINT fk_student_parent FOREIGN KEY (parent_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS students;
