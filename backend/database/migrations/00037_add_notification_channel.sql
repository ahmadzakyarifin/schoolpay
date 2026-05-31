-- +goose Up
ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS channel VARCHAR(20) NOT NULL DEFAULT 'system' AFTER type;

UPDATE notifications
SET channel = CASE
    WHEN type IN ('whatsapp', 'email', 'system') THEN type
    WHEN whatsapp_id IS NOT NULL AND whatsapp_id != '' THEN 'whatsapp'
    ELSE 'email'
END
WHERE channel = 'system' OR channel IS NULL OR channel = '';

-- +goose Down
ALTER TABLE notifications
DROP COLUMN IF EXISTS channel;
