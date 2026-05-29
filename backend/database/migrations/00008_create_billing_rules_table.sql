-- +goose Up
CREATE TABLE IF NOT EXISTS billing_rules (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    bill_type_id INT UNSIGNED NOT NULL,
    class_id INT UNSIGNED,
    target_type VARCHAR(50),
    target_id INT UNSIGNED,
    amount DECIMAL(15, 2) NOT NULL,
    period_type ENUM('bulanan', 'tahunan', 'sekali') DEFAULT 'bulanan',
    allow_installment TINYINT(1) DEFAULT 0,
    max_installment INT,
    due_day INT DEFAULT 10,
    start_date DATE,
    end_date DATE,
    is_active TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_rule_bill_type FOREIGN KEY (bill_type_id) REFERENCES bill_types(id) ON DELETE CASCADE,
    CONSTRAINT fk_rule_class FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS billing_rules;
