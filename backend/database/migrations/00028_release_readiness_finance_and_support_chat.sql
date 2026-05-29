-- +goose Up
-- Add complete period metadata while keeping legacy `period` for compatibility.
ALTER TABLE student_bills
  ADD COLUMN IF NOT EXISTS name VARCHAR(255) NULL AFTER billing_rule_id,
  ADD COLUMN IF NOT EXISTS period_month TINYINT UNSIGNED NULL AFTER period,
  ADD COLUMN IF NOT EXISTS period_year SMALLINT UNSIGNED NULL AFTER period_month,
  ADD COLUMN IF NOT EXISTS period_start_date DATE NULL AFTER period_year,
  ADD COLUMN IF NOT EXISTS period_end_date DATE NULL AFTER period_start_date,
  ADD COLUMN IF NOT EXISTS voided_at TIMESTAMP NULL AFTER deleted_at,
  ADD COLUMN IF NOT EXISTS void_reason TEXT NULL AFTER voided_at;

UPDATE student_bills sb
JOIN bill_types bt ON bt.id = sb.bill_type_id
SET
  sb.period_month = CASE WHEN sb.period REGEXP '^[0-9]{4}-[0-9]{2}$' THEN CAST(SUBSTRING(sb.period, 6, 2) AS UNSIGNED) ELSE sb.period_month END,
  sb.period_year = CASE WHEN sb.period REGEXP '^[0-9]{4}-[0-9]{2}$' THEN CAST(SUBSTRING(sb.period, 1, 4) AS UNSIGNED) ELSE sb.period_year END,
  sb.period_start_date = CASE WHEN sb.period REGEXP '^[0-9]{4}-[0-9]{2}$' THEN STR_TO_DATE(CONCAT(sb.period, '-01'), '%Y-%m-%d') ELSE sb.period_start_date END,
  sb.period_end_date = CASE WHEN sb.period REGEXP '^[0-9]{4}-[0-9]{2}$' THEN LAST_DAY(STR_TO_DATE(CONCAT(sb.period, '-01'), '%Y-%m-%d')) ELSE sb.period_end_date END,
  sb.name = COALESCE(sb.name, CONCAT(bt.name, CASE WHEN sb.period IS NULL OR sb.period = '' THEN '' ELSE CONCAT(' ', sb.period) END));

-- Add payment-order / reconciliation metadata without renaming existing tables yet.
ALTER TABLE payments
  ADD COLUMN IF NOT EXISTS external_order_id VARCHAR(100) NULL AFTER transaction_ref,
  ADD COLUMN IF NOT EXISTS external_transaction_id VARCHAR(100) NULL AFTER external_order_id,
  ADD COLUMN IF NOT EXISTS idempotency_key VARCHAR(160) NULL AFTER external_transaction_id,
  ADD COLUMN IF NOT EXISTS intent_bill_ids JSON NULL AFTER idempotency_key,
  ADD COLUMN IF NOT EXISTS gateway_raw_response JSON NULL AFTER gateway_id,
  ADD COLUMN IF NOT EXISTS reconcile_attempts INT NOT NULL DEFAULT 0 AFTER note,
  ADD COLUMN IF NOT EXISTS last_reconcile_error TEXT NULL AFTER reconcile_attempts,
  ADD COLUMN IF NOT EXISTS last_checked_at TIMESTAMP NULL AFTER last_reconcile_error,
  ADD COLUMN IF NOT EXISTS reconciled_at TIMESTAMP NULL AFTER last_checked_at,
  ADD COLUMN IF NOT EXISTS voided_at TIMESTAMP NULL AFTER deleted_at,
  ADD COLUMN IF NOT EXISTS reversal_of_payment_id INT UNSIGNED NULL AFTER voided_at;

CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_external_order_id ON payments (external_order_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_external_transaction_id ON payments (external_transaction_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_idempotency_key ON payments (idempotency_key);
CREATE INDEX IF NOT EXISTS idx_payments_reconcile_status ON payments (status, gateway_provider, created_at);

-- Internal CS inbox for one official school WhatsApp number.
CREATE TABLE IF NOT EXISTS support_conversations (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  parent_id INT UNSIGNED NULL,
  phone_number VARCHAR(32) NOT NULL,
  parent_name VARCHAR(255) NULL,
  status ENUM('open', 'pending', 'closed') NOT NULL DEFAULT 'open',
  assigned_admin_id INT UNSIGNED NULL,
  subject VARCHAR(255) NULL,
  unread_count INT NOT NULL DEFAULT 0,
  last_message TEXT NULL,
  last_message_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  closed_at TIMESTAMP NULL,
  CONSTRAINT fk_support_conversation_parent FOREIGN KEY (parent_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_support_conversation_admin FOREIGN KEY (assigned_admin_id) REFERENCES users(id) ON DELETE SET NULL,
  INDEX idx_support_conversation_phone_status (phone_number, status),
  INDEX idx_support_conversation_status (status),
  INDEX idx_support_conversation_last_message_at (last_message_at)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS support_messages (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  conversation_id INT UNSIGNED NOT NULL,
  sender_type ENUM('parent', 'admin', 'system') NOT NULL,
  sender_id INT UNSIGNED NULL,
  message TEXT NOT NULL,
  whatsapp_message_id VARCHAR(160) NULL,
  delivery_status VARCHAR(50) NOT NULL DEFAULT 'received',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_support_message_conversation FOREIGN KEY (conversation_id) REFERENCES support_conversations(id) ON DELETE CASCADE,
  INDEX idx_support_message_conversation_created (conversation_id, created_at),
  INDEX idx_support_message_whatsapp_id (whatsapp_message_id)
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS support_messages;
DROP TABLE IF EXISTS support_conversations;

DROP INDEX IF EXISTS idx_payments_reconcile_status ON payments;
DROP INDEX IF EXISTS idx_payments_idempotency_key ON payments;
DROP INDEX IF EXISTS idx_payments_external_transaction_id ON payments;
DROP INDEX IF EXISTS idx_payments_external_order_id ON payments;
ALTER TABLE payments
  DROP COLUMN IF EXISTS reversal_of_payment_id,
  DROP COLUMN IF EXISTS voided_at,
  DROP COLUMN IF EXISTS reconciled_at,
  DROP COLUMN IF EXISTS last_checked_at,
  DROP COLUMN IF EXISTS last_reconcile_error,
  DROP COLUMN IF EXISTS reconcile_attempts,
  DROP COLUMN IF EXISTS gateway_raw_response,
  DROP COLUMN IF EXISTS intent_bill_ids,
  DROP COLUMN IF EXISTS idempotency_key,
  DROP COLUMN IF EXISTS external_transaction_id,
  DROP COLUMN IF EXISTS external_order_id;

ALTER TABLE student_bills
  DROP COLUMN IF EXISTS void_reason,
  DROP COLUMN IF EXISTS voided_at,
  DROP COLUMN IF EXISTS period_end_date,
  DROP COLUMN IF EXISTS period_start_date,
  DROP COLUMN IF EXISTS period_year,
  DROP COLUMN IF EXISTS period_month,
  DROP COLUMN IF EXISTS name;
