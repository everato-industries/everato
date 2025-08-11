## Project metadata
APP_NAME := everato
CMD_PATH := .
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
TAILWIND := ./bin/tailwind

## DB config - use environment variables from .env file
DB_URL ?= postgres://piush:root_access@localhost:5432/everato?sslmode=disable
MIGRATIONS_DIR ?= internal/db/migrations

## Flags
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOFMT := gofmt -s
GOFLAGS := -mod=readonly
GOTAGS := -tags=prod

## Default target
all: build
.PHONY: all

install: sqlc golangci air golang-migrate tailwind templ browser-sync
	sudo chmod -R +x ./scripts
.PHONY: install

sqlc:
	@echo ">> Downloading sqlc..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
.PHONY: sqlc

golangci:
	@echo ">> Downloading golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6
.PHONY: golangci

tailwind:
	@echo ">> Downloading tailwindcss-cli into $(TAILWIND) ..."
	@curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
		chmod +x tailwindcss-linux-x64 && \
		mkdir -p ./bin && \
		mv tailwindcss-linux-x64 $(TAILWIND) && \
		$(TAILWIND) --version && \
		echo "Tailwind CSS downloaded and made executable."
.PHONY: tailwind

air:
	@echo ">> Downloading air..."
	@go install github.com/air-verse/air@latest
.PHONY: air

browser-sync:
	@echo ">> Downloading browser-sync..."
	@pnpm install -g browser-sync || \
		npm install -g browser-sync || \
		yarn global add browser-sync || \
		echo "Browser-sync installation failed. Please install it manually."
.PHONY: browser-sync

templ:
	@echo ">> Downloading templ..."
	@go install github.com/a-h/templ/cmd/templ@latest
.PHONY: templ

golang-migrate:
	@echo ">> Downloading golang-migrate..."
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
.PHONY: golang-migrate

## Build the Go project
build: clean
	@echo ">> Building the UI..."
	cd www && \
		pnpm install && \
		pnpm run build && \
	cd ..
	@echo ">> Building binary..."
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_FILE) $(GOTAGS) $(CMD_PATH)
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

## Generate the templ files
templ-watch:
	templ generate --watch
.PHONY: templ-watch

## Generate the css from tailwind classed
tailwind-watch:
	$(TAILWIND) -i ./styles/root.css -o ./public/css/styles.css --watch
.PHONY: tailwind-watch

serve:
	@echo ">> Starting the server..."
	browser-sync \
	    start \
		--server \
		--ss "./public/**/*" \
		--files "templates/pages/**.html" \
		--port 3000 \
		--no-open \
		--no-ui \
		--no-notify \
		--reload-delay 1000
		# --proxy "http://localhost:8080"
.PHONY: serve

watch:
	@echo ">> Watching for changes in templ files and css classes..."
	$(MAKE) --no-print-directory -j3 tailwind-watch serve dev
.PHONY: watch

## Clean the binary
clean:
	@echo ">> Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	@echo ">> Removing log files..."
	rm -rf $(LOGS_DIR)
	@echo ">> Removing auto-generated templ files..."
	rm -rf
.PHONY: clean
