-- +goose Up
ALTER TABLE deposit_movements
MODIFY COLUMN reason ENUM('OVERPAYMENT', 'BILL_VOIDED', 'BILL_AMOUNT_REDUCED', 'MANUAL_DEPOSIT', 'PAY_BILL', 'WITHDRAWAL', 'GATEWAY_ALLOCATION_FAILED') NOT NULL;

ALTER TABLE students
ADD CONSTRAINT chk_students_deposit_balance_non_negative CHECK (deposit_balance >= 0);

-- +goose Down
UPDATE deposit_movements SET reason = 'MANUAL_DEPOSIT' WHERE reason = 'GATEWAY_ALLOCATION_FAILED';
ALTER TABLE deposit_movements
MODIFY COLUMN reason ENUM('OVERPAYMENT', 'BILL_VOIDED', 'BILL_AMOUNT_REDUCED', 'MANUAL_DEPOSIT', 'PAY_BILL', 'WITHDRAWAL') NOT NULL;

ALTER TABLE students
DROP CONSTRAINT chk_students_deposit_balance_non_negative;
