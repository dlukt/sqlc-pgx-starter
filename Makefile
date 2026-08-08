# Common task runner for the sqlc + pgx starter.
# Run `make` (or `make help`) to see available targets.

BINARY   := bin/app
VERSION  := dev
LDFLAGS  := -X github.com/dlukt/sqlc-pgx-starter/internal/cli.Version=$(VERSION)

.DEFAULT_GOAL := help
.PHONY: help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
	/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# --- Go ---

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: build
build: ## Build the CLI into ./bin/app
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/app

.PHONY: run
run: ## Build and run the CLI (pass subcommands via ARGS, e.g. make run ARGS=serve)
	go run -ldflags "$(LDFLAGS)" ./cmd/app $(ARGS)

.PHONY: check
check: fmt vet build ## Run fmt, vet and build

.PHONY: test
test: ## Run all tests (integration tests require Docker)
	go test -race ./...

.PHONY: test-short
test-short: ## Run only unit tests (skip Docker integration tests)
	go test -race -short ./...

.PHONY: clean
clean: ## Remove built artifacts
	rm -rf bin/ dist/

# --- sqlc ---

.PHONY: generate
generate: ## Regenerate Go code from SQL with sqlc
	sqlc generate

.PHONY: sqlc-verify
sqlc-verify: generate ## Fail if generated code under internal/db is out of sync
	@git diff --exit-code -- internal/db || { \
		echo ""; \
		echo "error: internal/db is out of date. Run 'make generate' and commit the result."; \
		exit 1; }

# --- Migrations (goose, via the CLI) ---

.PHONY: migrate-up migrate-down migrate-redo migrate-status new-migration
migrate-up: ## Apply all pending migrations
	go run ./cmd/app migrate up

migrate-down: ## Roll back the most recent migration
	go run ./cmd/app migrate down

migrate-redo: ## Roll back and re-apply the most recent migration
	go run ./cmd/app migrate redo

migrate-status: ## Print migration status
	go run ./cmd/app migrate status

new-migration: ## Create a new migration: make new-migration NAME=add_index
	@test -n "$(NAME)" || { echo "Usage: make new-migration NAME=add_index"; exit 1; }
	go run ./cmd/app migrate create $(NAME)

# --- Local database (Docker) ---

.PHONY: db-up db-down db-reset db-logs
db-up: ## Start the local Postgres container
	docker compose up -d postgres

db-down: ## Stop the local Postgres container
	docker compose down

db-reset: ## Drop and recreate the local Postgres volume
	docker compose down -v
	docker compose up -d postgres

db-logs: ## Tail Postgres container logs
	docker compose logs -f postgres
