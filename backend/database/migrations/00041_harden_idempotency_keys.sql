-- +goose Up
ALTER TABLE idempotency_keys
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'PROCESSING' AFTER `key`,
    ADD COLUMN IF NOT EXISTS request_hash CHAR(64) NOT NULL DEFAULT '' AFTER status,
    MODIFY COLUMN response_payload LONGTEXT NULL,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER created_at;

UPDATE idempotency_keys
SET status = 'COMPLETED'
WHERE response_payload IS NOT NULL
  AND response_payload != '';

CREATE INDEX IF NOT EXISTS idx_idempotency_keys_cleanup
    ON idempotency_keys (status, updated_at);

-- `key` is already unique because it is the PRIMARY KEY in 00039.

-- +goose Down
DROP INDEX IF EXISTS idx_idempotency_keys_cleanup ON idempotency_keys;

ALTER TABLE idempotency_keys
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS request_hash,
    DROP COLUMN IF EXISTS status,
    MODIFY COLUMN response_payload LONGTEXT NOT NULL;
