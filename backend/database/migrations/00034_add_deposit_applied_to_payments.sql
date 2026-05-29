-- +goose Up
ALTER TABLE payments
  ADD COLUMN IF NOT EXISTS deposit_applied DECIMAL(15, 2) NOT NULL DEFAULT 0 AFTER amount;

-- +goose Down
ALTER TABLE payments
  DROP COLUMN IF EXISTS deposit_applied;
