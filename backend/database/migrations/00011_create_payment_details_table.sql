-- +goose Up
CREATE TABLE IF NOT EXISTS payment_details (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    payment_id INT UNSIGNED NOT NULL,
    student_bill_id INT UNSIGNED NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pd_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE,
    CONSTRAINT fk_pd_bill FOREIGN KEY (student_bill_id) REFERENCES student_bills(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS payment_details;
