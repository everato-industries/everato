## Project metadata
APP_NAME := everato
CMD_PATH := .
BIN_DIR := ./bin
LOGS_DIR := ./logs
BIN_FILE := $(BIN_DIR)/$(APP_NAME)

## Tools
GO := go
SQLC := sqlc
GOLANGCI_LINT := golangci-lint
MIGRATE := migrate
PNPM := pnpm

## DB config - default fallback or from .env file
DB_URL ?= postgres://piush:root_access@localhost:5432/everato?sslmode=disable
MIGRATIONS_DIR ?= internal/db/migrations

## Flags
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOFMT := gofmt -s
GOTAGS := -tags=prod

## Default target
all: build
.PHONY: all

## ─────────────────────────────────────────────────────────────────────────────
## Tool Installation & Setup
## ─────────────────────────────────────────────────────────────────────────────

install: sqlc golangci air golang-migrate
	@chmod +x ./scripts/*.sh 2>/dev/null || true
	@echo ">> Setup complete. Install frontend dependencies with: make frontend-install"
.PHONY: install

sqlc:
	@echo ">> Installing sqlc..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
.PHONY: sqlc

golangci:
	@echo ">> Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.6 2>/dev/null || \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
.PHONY: golangci

air:
	@echo ">> Installing air (hot reload)..."
	@go install github.com/air-verse/air@latest
.PHONY: air

golang-migrate:
	@echo ">> Installing golang-migrate..."
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
.PHONY: golang-migrate

## ─────────────────────────────────────────────────────────────────────────────
## Frontend Targets (React + TypeScript / Vite)
## ─────────────────────────────────────────────────────────────────────────────

frontend-install:
	@echo ">> Installing frontend dependencies..."
	@cd www && $(PNPM) install
.PHONY: frontend-install

frontend-dev:
	@echo ">> Starting frontend Vite dev server..."
	@cd www && $(PNPM) run dev
.PHONY: frontend-dev

frontend-build:
	@echo ">> Building frontend production bundle..."
	@cd www && $(PNPM) run build
.PHONY: frontend-build

frontend-lint:
	@echo ">> Linting frontend..."
	@cd www && $(PNPM) run lint
.PHONY: frontend-lint

frontend-check:
	@echo ">> Type-checking frontend with TypeScript..."
	@cd www && tsc -b
.PHONY: frontend-check

## ─────────────────────────────────────────────────────────────────────────────
## Backend & Full Build Targets
## ─────────────────────────────────────────────────────────────────────────────

## Build UI bundle and embed into single Go production binary
build: clean frontend-build
	@echo ">> Building Everato single binary with embedded assets..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_FILE) $(GOTAGS) $(CMD_PATH)
	@echo "✅ Build completed successfully: $(BIN_FILE)"
.PHONY: build

## Run the compiled production binary
run: build
	@echo ">> Running $(APP_NAME)..."
	$(BIN_FILE)
.PHONY: run

## Start backend with Air hot reload
dev:
	@echo ">> Starting Everato backend with Air hot reload..."
	@air
.PHONY: dev

## Start PostgreSQL database using Docker
db:
	@echo ">> Starting database container..."
	@docker compose -f docker/docker-compose.yaml up -d postgres 2>/dev/null || \
		sudo docker compose -f docker/docker-compose.yaml up -d postgres
.PHONY: db

## Start Loki, Grafana, Promtail monitoring stack
logs:
	@echo ">> Starting observability stack..."
	@docker compose -f docker/docker-compose.yaml up -d loki grafana promtail prometheus 2>/dev/null || \
		sudo docker compose -f docker/docker-compose.yaml up -d loki grafana promtail prometheus
.PHONY: logs

## Format all Go code
fmt:
	@echo ">> Formatting Go code..."
	@$(GOFMT) -w $(GO_FILES)
	@$(GO) fmt ./...
.PHONY: fmt

## Run linter
lint:
	@echo ">> Linting backend..."
	- $(GOLANGCI_LINT) run ./... || true
.PHONY: lint

## Run tests with race detection and coverage
test:
	@echo ">> Running Go test suite..."
	$(GO) test ./... -v -race -cover
.PHONY: test

## Run tests and output coverage report
test-coverage:
	@echo ">> Running tests with coverage profile..."
	$(GO) test ./... -v -race -cover -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo ">> Coverage report generated at coverage.html"
.PHONY: test-coverage

## Seed database with mockup data
seed:
	@echo ">> Seeding database with test data..."
	go run internal/db/seed/main.go
.PHONY: seed

## Generate Go repository code from SQL queries via sqlc
sqlc-gen:
	@echo ">> Generating code from SQL using sqlc..."
	$(SQLC) generate
.PHONY: sqlc-gen

## ─────────────────────────────────────────────────────────────────────────────
## Database Migration Targets
## ─────────────────────────────────────────────────────────────────────────────

## Apply all pending up migrations
migrate-up:
	@echo ">> Running migrations up..."
	@./scripts/migrate-up.sh
.PHONY: migrate-up

## Rollback the last migration
migrate-down:
	@echo ">> Rolling back last migration..."
	@./scripts/migrate-down.sh
.PHONY: migrate-down

## Force set migration version
migrate-force:
	@echo ">> Forcing migration version..."
	@./scripts/migrate-force.sh
.PHONY: migrate-force

## Drop all migration tables
migrate-drop:
	@echo ">> Dropping migration state..."
	@./scripts/migrate-drop.sh
.PHONY: migrate-drop

## Create a new migration file pair (.up.sql and .down.sql)
migrate-new:
	@read -p "Enter migration name: " name; \
		echo ">> Creating migration files..."; \
		migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq "$$name";
.PHONY: migrate-new

## ─────────────────────────────────────────────────────────────────────────────
## Clean
## ─────────────────────────────────────────────────────────────────────────────

clean:
	@echo ">> Cleaning build artifacts and logs..."
	@rm -rf $(BIN_DIR) $(LOGS_DIR) coverage.out coverage.html www/dist
.PHONY: clean
