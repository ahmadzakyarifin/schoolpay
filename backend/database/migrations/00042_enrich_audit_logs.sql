-- +goose Up
ALTER TABLE audit_logs 
    MODIFY user_id INT UNSIGNED NULL,
    MODIFY user_name VARCHAR(255) NULL,
    MODIFY role VARCHAR(100) NULL,
    MODIFY entity_type VARCHAR(100) NULL,
    MODIFY entity_id INT UNSIGNED NULL,
    ADD COLUMN method VARCHAR(10) NULL,
    ADD COLUMN path VARCHAR(255) NULL,
    ADD COLUMN device_id VARCHAR(120) NULL,
    ADD COLUMN app_platform VARCHAR(50) NULL,
    ADD COLUMN app_version VARCHAR(50) NULL,
    ADD COLUMN description TEXT NULL;

-- +goose Down
ALTER TABLE audit_logs 
    DROP COLUMN method,
    DROP COLUMN path,
    DROP COLUMN device_id,
    DROP COLUMN app_platform,
    DROP COLUMN app_version,
    DROP COLUMN description,
    MODIFY user_id INT UNSIGNED NOT NULL,
    MODIFY user_name VARCHAR(255) NOT NULL,
    MODIFY role VARCHAR(100) NOT NULL,
    MODIFY entity_type VARCHAR(100) NOT NULL,
    MODIFY entity_id INT UNSIGNED NOT NULL;
