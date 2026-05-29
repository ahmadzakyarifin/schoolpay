-- +goose Up
CREATE TABLE IF NOT EXISTS student_bills (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    student_id INT UNSIGNED NOT NULL,
    bill_type_id INT UNSIGNED NOT NULL,
    billing_rule_id INT UNSIGNED,
    academic_year VARCHAR(20),
    period VARCHAR(20),
    amount DECIMAL(15, 2) NOT NULL,
    total_paid DECIMAL(15, 2) DEFAULT 0,
    status ENUM('unpaid', 'partial', 'paid', 'overdue') DEFAULT 'unpaid',
    due_date DATE NOT NULL,
    end_date DATE,
    last_notified_at TIMESTAMP NULL,
    next_notified_at TIMESTAMP NULL,
    uq_key_period VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_bill_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_bill_type FOREIGN KEY (bill_type_id) REFERENCES bill_types(id) ON DELETE CASCADE,
    CONSTRAINT fk_bill_rule FOREIGN KEY (billing_rule_id) REFERENCES billing_rules(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS student_bills;
