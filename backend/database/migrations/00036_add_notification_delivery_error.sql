-- +goose Up
ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS delivery_error TEXT NULL AFTER delivery_status;

-- +goose Down
ALTER TABLE notifications
DROP COLUMN IF EXISTS delivery_error;
