-- +goose Up
-- 1. Modify student_bills status enum to support 'voided' status
ALTER TABLE student_bills MODIFY COLUMN status ENUM('unpaid', 'partial', 'paid', 'overdue', 'voided') DEFAULT 'unpaid';

-- 2. Add audit fields to payments table
ALTER TABLE payments 
ADD COLUMN created_by VARCHAR(50) DEFAULT 'SYSTEM',
ADD COLUMN bypass_reason TEXT NULL,
ADD COLUMN proof_attachment VARCHAR(255) NULL,
ADD COLUMN note TEXT NULL;

-- 3. Create deposit_movements table
CREATE TABLE IF NOT EXISTS deposit_movements (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    student_id INT UNSIGNED NOT NULL,
    type ENUM('IN', 'OUT') NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    reason ENUM('OVERPAYMENT', 'BILL_VOIDED', 'MANUAL_DEPOSIT', 'PAY_BILL', 'WITHDRAWAL') NOT NULL,
    reference_id VARCHAR(50) NULL,
    created_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_dm_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS deposit_movements;

ALTER TABLE payments 
DROP COLUMN created_by,
DROP COLUMN bypass_reason,
DROP COLUMN proof_attachment,
DROP COLUMN note;

ALTER TABLE student_bills MODIFY COLUMN status ENUM('unpaid', 'partial', 'paid', 'overdue') DEFAULT 'unpaid';
