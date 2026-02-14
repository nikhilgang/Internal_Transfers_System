.PHONY: help build run test lint clean local-db-create local-db-drop local-migrate-up local-setup db-up db-down migrate-up setup

# ── Variables ────────────────────────────────────────────────────────────────
APP_NAME   := internal-transfers
BIN_DIR    := bin
DB_DSN     := postgres://postgres:postgres@localhost:5432/transaction_manager
DB_USER    := postgres
DB_NAME    := transaction_manager

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

build: ## Build the server binary
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/server ./cmd/server

run: build ## Build and run the server
	./$(BIN_DIR)/server

test: ## Run all tests
	go test -v -race -count=1 ./...

lint: ## Run staticcheck linter
	staticcheck ./...

clean: ## Remove build artefacts
	rm -rf $(BIN_DIR)

# ── Local (no Docker) database targets ───────────────────────────────────────

local-db-create: ## Create the local PostgreSQL database (requires psql)
	@echo "Creating database '$(DB_NAME)'..."
	@psql -U $(DB_USER) -tc "SELECT 1 FROM pg_database WHERE datname = '$(DB_NAME)'" | grep -q 1 || psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME)"
	@echo "Database '$(DB_NAME)' is ready."

local-db-drop: ## Drop the local PostgreSQL database
	@echo "Dropping database '$(DB_NAME)'..."
	@psql -U $(DB_USER) -c "DROP DATABASE IF EXISTS $(DB_NAME)"
	@echo "Done."

local-migrate-up: ## Run migrations UP on local PostgreSQL (requires psql)
	@echo "Applying migrations..."
	@psql -U $(DB_USER) -d $(DB_NAME) -f migrations/000001_init_schema.up.sql
	@echo "Migrations applied."

local-setup: local-db-create local-migrate-up ## Full local setup without Docker: create DB + run migrations
	@echo "Local setup complete. Run 'make run' to start the server."

# ── Docker-based database targets (optional) ─────────────────────────────────

db-up: ## Start PostgreSQL via Docker Compose
	docker compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker compose exec postgres pg_isready -U postgres > /dev/null 2>&1; do sleep 1; done
	@echo "PostgreSQL is ready."

db-down: ## Stop PostgreSQL and remove volumes
	docker compose down -v

migrate-up: ## Run database migrations UP (Docker)
	@docker compose exec postgres psql -U postgres -d transaction_manager -f /dev/stdin < migrations/000001_init_schema.up.sql

setup: db-up ## Full Docker setup: start DB and run migrations
	@echo "Applying migrations..."
	@sleep 2
	@docker compose exec -T postgres psql -U postgres -d transaction_manager < migrations/000001_init_schema.up.sql
	@echo "Setup complete. Run 'make run' to start the server."
