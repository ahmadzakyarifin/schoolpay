SHELL := /bin/sh

APP_NAME ?= schoolpay
BINARY ?= main

# Cari .env dan compose file dari folder sekarang naik sampai root.
# Override contoh: make env-info ENV_FILE=/path/.env COMPOSE_FILE=/path/docker-compose.yml
ENV_NAMES ?= .env env .env.local
COMPOSE_NAMES ?= docker-compose.yml docker-compose.yaml compose.yml compose.yaml

ENV_FILE ?= $(shell pwd="$(CURDIR)"; while [ "$$pwd" != "/" ]; do \
	for file in $(ENV_NAMES); do \
		if [ -f "$$pwd/$$file" ]; then printf '%s\n' "$$pwd/$$file"; exit 0; fi; \
	done; \
	pwd=$$(dirname "$$pwd"); \
done)

COMPOSE_FILE ?= $(shell pwd="$(CURDIR)"; while [ "$$pwd" != "/" ]; do \
	for file in $(COMPOSE_NAMES); do \
		if [ -f "$$pwd/$$file" ]; then printf '%s\n' "$$pwd/$$file"; exit 0; fi; \
	done; \
	pwd=$$(dirname "$$pwd"); \
done)

APP_ENTRY_NAMES ?= ./cmd/api/main.go ./cmd/main.go ./main.go .
APP_ENTRY ?= $(shell for path in $(APP_ENTRY_NAMES); do \
	if [ "$$path" = "." ] || [ -f "$$path" ]; then printf '%s\n' "$$path"; exit 0; fi; \
done)

ifneq ($(strip $(ENV_FILE)),)
	include $(ENV_FILE)
endif

.EXPORT_ALL_VARIABLES:

DB_USER ?= root
DB_PASS ?= password
DB_HOST ?= 127.0.0.1
DB_PORT ?= 3306
DB_NAME ?= $(APP_NAME)
DB_DOCKER_HOST ?= db
DB_DOCKER_PORT ?= 3306
PORT ?= 8080
GOCACHE ?= /tmp/$(APP_NAME)-go-build-cache

MIGRATIONS_DIR ?= database/migrations
DOCKER_MIGRATIONS_DIR ?= /migrations
GOOSE_BIN ?= goose
COMPOSE_BIN ?= docker compose
COMPOSE_RUN_FLAGS ?= --rm
VERSION ?= $(v)
NAME ?= $(name)
TYPE ?= sql
INTEGRATION_DSN_ENV ?= SCHOOLPAY_TEST_DSN
INTEGRATION_DSN ?= schoolpay_test:schoolpay_test@tcp(127.0.0.1:3307)/schoolpay_test?parseTime=true

LOCAL_DSN := $(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true
DOCKER_DSN := $(DB_USER):$(DB_PASS)@tcp($(DB_DOCKER_HOST):$(DB_DOCKER_PORT))/$(DB_NAME)?parseTime=true

ifneq ($(strip $(COMPOSE_FILE)),)
	COMPOSE := $(COMPOSE_BIN) -f "$(COMPOSE_FILE)"
endif
ifneq ($(strip $(ENV_FILE)),)
	COMPOSE_ENV := --env-file "$(ENV_FILE)"
endif

GOOSE_LOCAL := $(GOOSE_BIN) -dir "$(MIGRATIONS_DIR)" mysql "$(LOCAL_DSN)"
GOOSE_LOCAL_CREATE := $(GOOSE_BIN) -dir "$(MIGRATIONS_DIR)" create
GOOSE_DOCKER := $(COMPOSE) $(COMPOSE_ENV) run $(COMPOSE_RUN_FLAGS) migrate goose -dir "$(DOCKER_MIGRATIONS_DIR)" mysql "$(DOCKER_DSN)"
GOOSE_DOCKER_CREATE := $(COMPOSE) $(COMPOSE_ENV) run $(COMPOSE_RUN_FLAGS) --no-deps migrate goose -dir "$(DOCKER_MIGRATIONS_DIR)" create

.PHONY: help env-info run build test test-unit test-integration clean \
	check-goose install-goose \
	docker-build docker-up docker-down docker-restart docker-logs docker-ps docker-clean-orphans \
	migrate-status migrate-up migrate-up-one migrate-up-to migrate-down migrate-down-to \
	migrate-redo migrate-reset migrate-version migrate-create migrate-fix migrate-validate migrate-force \
	docker-migrate-status docker-migrate-up docker-migrate-up-one docker-migrate-up-to \
	docker-migrate-down docker-migrate-down-to docker-migrate-redo docker-migrate-reset \
	docker-migrate-version docker-migrate-create docker-migrate-fix docker-migrate-validate docker-migrate-force \
	app-up app-down app-restart app-logs app-ps app-clean \
	d-status d-up d-up-one d-up-to d-down d-down-to d-redo d-reset d-version d-create d-fix d-validate d-force \
	status up up-one up-to down down-to redo reset version create fix validate force \
	mockery

help:
	@printf '%s\n' \
		'App:' \
		'  make run                  Run app lokal kalau go.mod ada' \
		'  make test                 Run test kalau go.mod ada' \
		'  make app-up               Start Docker Compose kalau compose file ada' \
		'  make app-down             Stop Docker Compose kalau compose file ada' \
		'  make app-logs             Lihat log Docker Compose' \
		'' \
		'Migrasi lokal:' \
		'  make create name=create_users' \
		'  make up | down | status | version' \
		'  make reset | redo | fix | validate' \
		'  make up-to v=20260605010101' \
		'' \
		'Migrasi Docker:' \
		'  make d-create name=create_users' \
		'  make d-up | d-down | d-status | d-version' \
		'  make d-reset | d-redo | d-fix | d-validate' \
		'  make d-up-to v=20260605010101' \
		'' \
		'Auto discovery:' \
		'  ENV_FILE      auto cari .env/env/.env.local dari current dir ke root' \
		'  COMPOSE_FILE  auto cari docker-compose.yml/compose.yml dari current dir ke root'

env-info:
	@printf 'ENV_FILE=%s\n' "$(ENV_FILE)"
	@printf 'COMPOSE_FILE=%s\n' "$(COMPOSE_FILE)"
	@printf 'APP_NAME=%s\n' "$(APP_NAME)"
	@printf 'APP_ENTRY=%s\n' "$(APP_ENTRY)"
	@printf 'PORT=%s\n' "$(PORT)"
	@printf 'DB_USER=%s\n' "$(DB_USER)"
	@printf 'DB_HOST=%s\n' "$(DB_HOST)"
	@printf 'DB_PORT=%s\n' "$(DB_PORT)"
	@printf 'DB_NAME=%s\n' "$(DB_NAME)"
	@printf 'DB_DOCKER_HOST=%s\n' "$(DB_DOCKER_HOST)"
	@printf 'DB_DOCKER_PORT=%s\n' "$(DB_DOCKER_PORT)"

run:
	@if [ ! -f go.mod ]; then \
		printf 'go.mod tidak ada di %s, skip run.\n' "$(CURDIR)"; \
	else \
		go run "$(APP_ENTRY)"; \
	fi

build:
	@if [ ! -f go.mod ]; then \
		printf 'go.mod tidak ada di %s, skip build.\n' "$(CURDIR)"; \
	else \
		go build -o "$(BINARY)" "$(APP_ENTRY)"; \
	fi

test: test-unit

test-unit:
	@if [ ! -f go.mod ]; then \
		printf 'go.mod tidak ada di %s, skip test.\n' "$(CURDIR)"; \
	else \
		GOCACHE="$(GOCACHE)" go test ./...; \
	fi

test-integration:
	@if [ ! -f go.mod ]; then \
		printf 'go.mod tidak ada di %s, skip integration test.\n' "$(CURDIR)"; \
	else \
		GOCACHE="$(GOCACHE)" $(INTEGRATION_DSN_ENV)="$(INTEGRATION_DSN)" go test -tags=integration ./...; \
	fi

clean:
	@rm -f "$(BINARY)"

check-goose:
	@command -v "$(GOOSE_BIN)" >/dev/null 2>&1 || ( \
		printf '%s\n' \
			'Goose belum terinstall untuk migrasi lokal.' \
			'Pilih salah satu:' \
			'  make install-goose' \
			'  make d-create name=create_users_table' \
			'  make d-up'; \
		exit 127; \
	)

install-goose:
	@go install github.com/pressly/goose/v3/cmd/goose@latest

docker-build:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker-build.\n'; else $(COMPOSE) $(COMPOSE_ENV) build; fi

docker-up:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker-up.\n'; else $(COMPOSE) $(COMPOSE_ENV) up -d --build; fi

docker-down:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker-down.\n'; else $(COMPOSE) $(COMPOSE_ENV) down; fi

docker-restart: docker-down docker-up

docker-logs:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker-logs.\n'; else $(COMPOSE) $(COMPOSE_ENV) logs -f; fi

docker-ps:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker-ps.\n'; else $(COMPOSE) $(COMPOSE_ENV) ps; fi

docker-clean-orphans:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker-clean-orphans.\n'; else $(COMPOSE) $(COMPOSE_ENV) down --remove-orphans; fi

migrate-status:
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) status; fi

migrate-up:
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) up; fi

migrate-up-one:
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) up-by-one; fi

migrate-up-to:
	@test -n "$(VERSION)" || (printf 'Isi VERSION dulu. Contoh: make migrate-up-to VERSION=20260605010101\n' && exit 1)
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) up-to "$(VERSION)"; fi

migrate-down:
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) down; fi

migrate-down-to:
	@test -n "$(VERSION)" || (printf 'Isi VERSION dulu. Contoh: make migrate-down-to VERSION=0\n' && exit 1)
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) down-to "$(VERSION)"; fi

migrate-redo:
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) redo; fi

migrate-reset:
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) reset; fi

migrate-version:
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) version; fi

migrate-create: check-goose
	@mkdir -p "$(MIGRATIONS_DIR)"
	@if [ -n "$(NAME)" ]; then \
		$(GOOSE_LOCAL_CREATE) "$(NAME)" "$(TYPE)"; \
	else \
		read -p "Nama migrasi: " input_name; \
		$(GOOSE_LOCAL_CREATE) "$$input_name" "$(TYPE)"; \
	fi

migrate-fix:
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) fix; fi

migrate-validate:
	@if [ ! -d "$(MIGRATIONS_DIR)" ]; then printf 'Folder %s tidak ada, skip migrasi.\n' "$(MIGRATIONS_DIR)"; \
	elif ! command -v "$(GOOSE_BIN)" >/dev/null 2>&1; then $(MAKE) check-goose; \
	else $(GOOSE_LOCAL) validate; fi

migrate-force:
	@printf '%s\n' \
		'Goose tidak punya command force.' \
		'Pakai migrate-up-to, migrate-down-to, migrate-reset, atau migrate-fix.'
	@exit 1

docker-migrate-status:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) status; fi

docker-migrate-up:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) up; fi

docker-migrate-up-one:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) up-by-one; fi

docker-migrate-up-to:
	@test -n "$(VERSION)" || (printf 'Isi VERSION dulu. Contoh: make d-up-to v=20260605010101\n' && exit 1)
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) up-to "$(VERSION)"; fi

docker-migrate-down:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) down; fi

docker-migrate-down-to:
	@test -n "$(VERSION)" || (printf 'Isi VERSION dulu. Contoh: make d-down-to v=0\n' && exit 1)
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) down-to "$(VERSION)"; fi

docker-migrate-redo:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) redo; fi

docker-migrate-reset:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) reset; fi

docker-migrate-version:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) version; fi

docker-migrate-create:
	@mkdir -p "$(MIGRATIONS_DIR)"
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate create.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate create.\n'; \
	elif [ -n "$(NAME)" ]; then $(GOOSE_DOCKER_CREATE) "$(NAME)" "$(TYPE)"; \
	else read -p "Nama migrasi: " input_name; $(GOOSE_DOCKER_CREATE) "$$input_name" "$(TYPE)"; fi

docker-migrate-fix:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) fix; fi

docker-migrate-validate:
	@if [ -z "$(COMPOSE_FILE)" ]; then printf 'Compose file tidak ditemukan, skip docker migrate.\n'; \
	elif ! $(COMPOSE) $(COMPOSE_ENV) config --services | grep -qx migrate; then printf 'Service migrate tidak ada di compose, skip docker migrate.\n'; \
	else $(GOOSE_DOCKER) validate; fi

docker-migrate-force:
	@$(MAKE) migrate-force

app-up: docker-up
app-down: docker-down
app-restart: docker-restart
app-logs: docker-logs
app-ps: docker-ps
app-clean: docker-clean-orphans

d-status: docker-migrate-status
d-up: docker-migrate-up
d-up-one: docker-migrate-up-one
d-up-to: docker-migrate-up-to
d-down: docker-migrate-down
d-down-to: docker-migrate-down-to
d-redo: docker-migrate-redo
d-reset: docker-migrate-reset
d-version: docker-migrate-version
d-create: docker-migrate-create
d-fix: docker-migrate-fix
d-validate: docker-migrate-validate
d-force: docker-migrate-force

status: migrate-status
up: migrate-up
up-one: migrate-up-one
up-to: migrate-up-to
down: migrate-down
down-to: migrate-down-to
redo: migrate-redo
reset: migrate-reset
version: migrate-version
create: migrate-create
fix: migrate-fix
validate: migrate-validate
force: migrate-force

mockery:
	@if command -v mockery >/dev/null 2>&1; then mockery; else printf 'mockery tidak terinstall, skip.\n'; fi
