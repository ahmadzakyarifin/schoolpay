# Database Migration Variables
DB_URL="root:password@tcp(127.0.0.1:3306)/schoolpay?parseTime=true"
MIGRATIONS_DIR=database/migrations

.PHONY: migrate-status migrate-up migrate-down migrate-redo migrate-reset migrate-version migrate-create

# --- LOKAL (Jika Goose terinstall di Host) ---

migrate-status:
	@goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) status

migrate-up:
	@goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) up

migrate-down:
	@goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) down

migrate-redo:
	@goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) redo

migrate-reset:
	@goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) reset

migrate-version:
	@goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) version

migrate-create:
	@read -p "Masukkan nama migrasi: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql


# --- DOCKER (Menjalankan Goose di dalam Container) ---

docker-migrate-status:
	docker-compose exec app goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) status

docker-migrate-up:
	docker-compose exec app goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) up

docker-migrate-down:
	docker-compose exec app goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) down

docker-migrate-reset:
	docker-compose exec app goose -dir $(MIGRATIONS_DIR) mysql $(DB_URL) reset

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down


# --- RELEASE READINESS ---

test-unit:
	go test ./...

test-integration:
	SCHOOLPAY_TEST_DSN="schoolpay_test:schoolpay_test@tcp(127.0.0.1:3307)/schoolpay_test?parseTime=true" SCHOOLPAY_TEST_REDIS="127.0.0.1:6380" go test -tags=integration ./...

mockery:
	mockery

openapi-validate:
	@test -f docs/openapi.json
