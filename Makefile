## Project metadata
APP_NAME := everato
CMD_PATH := ./cmd
BIN_DIR := ./bin
LOGS_DIR := ./logs
BIN_FILE := $(BIN_DIR)/$(APP_NAME)
TEMPL_FILES := ./templates/views/**/*.templ
GENERATED_TEMPL_FILES := ./templates/views/**/*_templ.go

## Tools
GO := go
TEMPL := templ
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

## Default target
all: build
.PHONY: all

bootstrap: sqlc golangci air golang-migrate
	sudo chmod -R +x ./scripts
.PHONY: bootstrap

sqlc:
	@echo ">> Downloading sqlc..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
.PHONY: sqlc

golangci:
	@echo ">> Downloading golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6
.PHONY: golangci

air:
	@echo ">> Downloading air..."
	go install github.com/air-verse/air@latest
.PHONY: air

golang-migrate:
	@echo ">> Downloading golang-migrate..."
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
.PHONY: golang-migrate

## Build the Go project
build:
	@echo ">> Building binary..."
	mkdir -p $(BIN_DIR)
	$(TEMPL) generate
	$(GO) build -o $(BIN_FILE) $(CMD_PATH)
.PHONY: build

## Run the server
run: build
	@echo ">> Running $(APP_NAME)..."
	$(BIN_FILE)
.PHONY: run

dev:
	@air
.PHONY: dev

db:
	@echo ">> Running the database ..."
	@sudo docker compose -f docker/docker-compose.yaml up postgres
.PHONY: db

logs:
	@echo ">> Running the logs ..."
	@sudo docker compose -f docker/docker-compose.yaml up loki grafana promtail prometheus
.PHONY: logs

## Format all Go code
fmt:
	@echo ">> Formatting..."
	$(TEMPL) fmt $(TEMPL_FILES)
	$(GOFMT) -w $(GO_FILES)
	$(GO) fmt ./...
.PHONY: fmt

## Run linter (requires golangci-lint)
lint:
	@echo ">> Linting..."
	- $(GOLANGCI_LINT) run ./... || true
.PHONY: lint

## Run tests
test:
	@echo ">> Running tests..."
	$(GO) test ./... -v -race -cover
.PHONY: test

## Seed database with mockup data
seed:
	@echo ">> Seeding database with mockup data..."
	go run internal/db/seed/main.go
.PHONY: seed

## Generate Go code from SQL (via sqlc)
sqlc-gen:
	@echo ">> Generating code from SQL using sqlc..."
	$(SQLC) generate
.PHONY: sqlc-gen

## Apply all up migrations
migrate-up:
	@echo ">> Running migrations up..."
	@./scripts/migrate-up.sh
.PHONY: migrate-up

## Rollback last migration
migrate-down:
	@echo ">> Rolling back last migration..."
	@./scripts/migrate-down.sh
.PHONY: migrate-down

## Force set migration version (useful for fixing state)
migrate-force:
	@echo ">> Forcing migration version..."
	@./scripts/migrate-force.sh
.PHONY: migrate-force

## Drop the state
migrate-drop:
	@echo ">> Dropping migration state..."
	@./scripts/migrate-drop.sh
.PHONY: migrate-drop

## Create a new migration file
migrate-new:
	@read -p "Enter migration name: " name; \
		echo ">> Created:"; \
		migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq "$$name";
.PHONY: migrate-new

## Clean the binary
clean:
	@echo ">> Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	@echo ">> Removing log files..."
	rm -rf $(LOGS_DIR)
	@echo ">> Removing auto-generated templ files..."
	rm -rf
.PHONY: clean
