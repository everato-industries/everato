## Project metadata
APP_NAME := everato
CMD_PATH := ./cmd/app
BIN_DIR := ./bin
LOGS_DIR := ./logs
BIN_FILE := $(BIN_DIR)/$(APP_NAME)

## Tools
GO := go
SQLC := sqlc
GOLANGCI_LINT := golangci-lint
MIGRATE := migrate

## DB config - use environment variables from .env file
DB_URL ?= postgres://piush:root_access@localhost:5432/everato?sslmode=disable
MIGRATIONS_DIR ?= internal/db/migrations

## Flags
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOFMT := gofmt -s
GOFLAGS := -mod=readonly

## Build target
.PHONY: all build run fmt lint test migrate-up migrate-down clean help bootstrap sqlc golangci air db golang-migrate

all: build

bootstrap: sqlc golangci air golang-migrate
	sudo chmod -R +x ./scripts

sqlc:
	@echo ">> Downloading sqlc..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

golangci:
	@echo ">> Downloading golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6

air:
	@echo ">> Downloading air..."
	go install github.com/air-verse/air@latest

golang-migrate:
	@echo ">> Downloading golang-migrate..."
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## Build the Go project
build:
	@echo ">> Building binary..."
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_FILE) $(CMD_PATH)

## Run the server
run: build
	@echo ">> Running $(APP_NAME)..."
	@$(BIN_FILE)

dev:
	@air

db:
	@echo ">> Running the database ..."
	@sudo docker compose -f docker/docker-compose.yaml up postgres

## Format all Go code
fmt:
	@echo ">> Formatting..."
	@$(GOFMT) -w $(GO_FILES)
	@go fmt ./...

## Run linter (requires golangci-lint)
lint:
	@echo ">> Linting..."
	$(GOLANGCI_LINT) run ./...

## Run tests
test:
	@echo ">> Running tests..."
	@$(GO) test ./... -v -race -cover

## Generate Go code from SQL (via sqlc)
sqlc-gen:
	@echo ">> Generating code from SQL using sqlc..."
	$(SQLC) generate

## Apply all up migrations
migrate-up:
	@echo ">> Running migrations up..."
	@./scripts/migrate-up.sh

## Rollback last migration
migrate-down:
	@echo ">> Rolling back last migration..."
	@./scripts/migrate-down.sh

## Force set migration version (useful for fixing state)
migrate-force:
	@echo ">> Forcing migration version..."
	@./scripts/migrate-force.sh

## Create a new migration file
migrate-new:
	@read -p "Enter migration name: " name; \
		migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq "$$name"; \
    	echo ">> Created: {SEQ_NUMBER}_$${name}.up.sql and .down.sql"

## Clean the binary
clean:
	@echo ">> Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	@echo ">> Removing log files..."
	rm -rf $(LOGS_DIR)
