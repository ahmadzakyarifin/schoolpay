-- +goose Up
ALTER TABLE students ADD COLUMN deposit_balance DOUBLE PRECISION DEFAULT 0;

-- +goose Down
ALTER TABLE students DROP COLUMN deposit_balance;
