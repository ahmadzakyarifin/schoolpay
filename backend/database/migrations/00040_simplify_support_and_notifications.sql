-- +goose Up
DROP TABLE IF EXISTS support_messages;

ALTER TABLE support_conversations
  DROP FOREIGN KEY fk_support_conversation_admin;

ALTER TABLE support_conversations
  DROP COLUMN IF EXISTS assigned_admin_id,
  DROP COLUMN IF EXISTS subject,
  DROP COLUMN IF EXISTS unread_count,
  DROP COLUMN IF EXISTS last_message,
  DROP COLUMN IF EXISTS last_message_at,
  DROP COLUMN IF EXISTS closed_at;

-- +goose Down
-- Tidak didukung untuk rollback penuh karena tabel support_messages sudah diubah secara destruktif demi efisiensi.
