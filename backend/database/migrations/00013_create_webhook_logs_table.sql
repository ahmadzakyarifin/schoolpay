-- +goose Up
CREATE TABLE IF NOT EXISTS webhook_logs (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    provider VARCHAR(50),
    event_id VARCHAR(100) UNIQUE,
    payload JSON,
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS webhook_logs;
