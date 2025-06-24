# Project metadata
APP_NAME := everato
CMD_PATH := ./cmd
BIN_DIR := ./bin
BIN_FILE := $(BIN_DIR)/$(APP_NAME)

# Tools
GO := go
SQLC := sqlc
GOLANGCI_LINT := golangci-lint
MIGRATE := migrate

# DB config - use environment variables from .env file
DB_URL ?= postgres://piush:root_access@localhost:5432/everato?sslmode=disable
MIGRATIONS_DIR ?= internal/db/migrations

# Flags
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOFMT := gofmt -s
GOFLAGS := -mod=readonly
ENV_FILE := .env

# Build target
.PHONY: all build run fmt lint test migrate-up migrate-down clean help bootstrap sqlc golangci

all: build

bootstrap: sqlc golangci

sqlc:
	@echo ">> Downloading sqlc..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

golangci:
	@echo ">> Downloading golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6

## Build the Go project
build:
	@echo ">> Building binary..."
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_FILE) $(CMD_PATH)

## Run the server
run: build
	@echo ">> Running $(APP_NAME)..."
	@$(BIN_FILE)

## Run the server with .env variables loaded
run-env: build
	@echo ">> Running with .env"
	@set -a && source $(ENV_FILE) && set +a && $(BIN_FILE)

## Format all Go code
fmt:
	@echo ">> Formatting..."
	@$(GOFMT) -w $(GO_FILES)

## Run linter (requires golangci-lint)
lint:
	@echo ">> Linting..."
	@$(GOLANGCI_LINT) run ./...

## Run tests
test:
	@echo ">> Running tests..."
	@$(GO) test ./... -v -race -cover

## Generate Go code from SQL (via sqlc)
sqlc-gen:
	@echo ">> Generating code from SQL using sqlc..."
	@$(SQLC) generate

## Apply all up migrations
migrate-up:
	@echo ">> Running migrations up..."
	@set -a && source $(ENV_FILE) && set +a && $(MIGRATE) -database "$$DB_URL" -path "$$MIGRATIONS_DIR" up

## Rollback last migration
migrate-down:
	@echo ">> Rolling back last migration..."
	@set -a && source $(ENV_FILE) && set +a && $(MIGRATE) -database "$$DB_URL" -path "$$MIGRATIONS_DIR" down 1

## Force set migration version (useful for fixing state)
migrate-force:
	@set -a && source $(ENV_FILE) && set +a && \
	read -p "Enter version to force to: " version; \
	$(MIGRATE) -database "$$DB_URL" -path "$$MIGRATIONS_DIR" force $$version

## Create a new migration file
migrate-new:
	@set -a && source $(ENV_FILE) && set +a && \
	read -p "Enter migration name: " name; \
	timestamp=$$(date +%s); \
	mkdir -p "$$MIGRATIONS_DIR"; \
	touch "$$MIGRATIONS_DIR/$${timestamp}_$${name}.up.sql"; \
	touch "$$MIGRATIONS_DIR/$${timestamp}_$${name}.down.sql"; \
	echo ">> Created: $${timestamp}_$${name}.up.sql and .down.sql"

## Clean the binary
clean:
	@echo ">> Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
