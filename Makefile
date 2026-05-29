.PHONY: test-db-up test-db-down test-backend test-frontend test-all mockery

test-db-up:
	docker compose -f docker-compose.test.yml up -d test-db test-redis

test-db-down:
	docker compose -f docker-compose.test.yml down -v

test-backend:
	$(MAKE) -C backend test-unit

test-frontend:
	cd frontend && npm test

test-all: test-backend test-frontend

mockery:
	$(MAKE) -C backend mockery
