# Load environment variables from backend/.env if present
-include backend/.env
export

# Database defaults
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= inventory
DB_SSLMODE ?= disable
DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Migrations configuration
MIGRATIONS_DIR := backend/internal/database/migrations
MIGRATE := migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)"

.PHONY: help
help: ## Display this help message
	@echo "Available commands:"
	@echo ""
	@echo "  Database Migrations:"
	@echo "    make migrate-up          Apply all pending migrations"
	@echo "    make migrate-down        Roll back the most recent migration (1 step)"
	@echo "    make migrate-down-all    Roll back all migrations"
	@echo "    make migrate-step N=1    Apply (+N) or rollback (-N) steps"
	@echo "    make migrate-status      Show the current migration version"
	@echo "    make migrate-force V=1   Force migration version (resolves dirty state)"
	@echo "    make migrate-create NAME=add_foo  Create a new numbered migration pair"
	@echo ""
	@echo "  Backend Development:"
	@echo "    make build               Build backend binary"
	@echo "    make run                 Run backend service"
	@echo "    make test                Run backend unit and integration tests"
	@echo "    make test-coverage       Run tests and generate HTML coverage report"
	@echo "    make tidy                Clean and verify Go module dependencies"
	@echo "    make swagger             Regenerate Swagger API documentation"
	@echo "    make clean               Remove build artifacts"
	@echo ""
	@echo "  Frontend Development:"
	@echo "    make frontend-install    Install frontend dependencies"
	@echo "    make frontend-dev        Start frontend Vite dev server"
	@echo "    make frontend-build      Build frontend production bundle"
	@echo "    make frontend-lint       Run frontend ESLint"
	@echo ""

# Database Migrations (golang-migrate)

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	@echo "Applying migrations to $(DB_NAME)..."
	$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## Roll back the most recent migration (1 step)
	@echo "Rolling back 1 migration..."
	$(MIGRATE) down 1

.PHONY: migrate-down-all
migrate-down-all: ## Roll back all migrations
	@echo "Rolling back all migrations..."
	$(MIGRATE) down -all

.PHONY: migrate-step
migrate-step: ## Apply (+N) or rollback (-N) steps: make migrate-step N=2
	@if [ -z "$(N)" ]; then echo "Error: Specify step count, e.g. make migrate-step N=1"; exit 1; fi
	$(MIGRATE) step $(N)

.PHONY: migrate-status
migrate-status: ## Show the current migration version
	$(MIGRATE) version

.PHONY: migrate-force
migrate-force: ## Force set version to recover from dirty state: make migrate-force V=1
	@if [ -z "$(V)" ]; then echo "Error: Specify version, e.g. make migrate-force V=1"; exit 1; fi
	$(MIGRATE) force $(V)

.PHONY: migrate-create
migrate-create: ## Create new migration: make migrate-create NAME=create_foo
	@if [ -z "$(NAME)" ]; then echo "Error: Specify NAME, e.g. make migrate-create NAME=create_foo"; exit 1; fi
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

# Backend Commands

.PHONY: build
build: ## Build backend binary
	@echo "Building backend..."
	cd backend && go build -o bin/api cmd/i_m_s/main.go

.PHONY: run
run: ## Run backend server
	@echo "Starting backend..."
	cd backend && go run cmd/i_m_s/main.go

.PHONY: test
test: ## Run backend tests
	@echo "Running tests..."
	cd backend && go test -v -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests and produce coverage report
	@echo "Running tests with coverage..."
	cd backend && go test -v -race -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to backend/coverage.html"

.PHONY: tidy
tidy: ## Go mod tidy
	cd backend && go mod tidy

.PHONY: swagger
swagger: ## Generate Swagger API documentation
	@echo "Generating Swagger documentation..."
	cd backend && swag init -g cmd/i_m_s/main.go -o docs

.PHONY: clean
clean: ## Clean binaries and coverage reports
	rm -rf backend/bin backend/coverage.out backend/coverage.html
	rm -rf frontend/dist

# Frontend Commands

.PHONY: frontend-install
frontend-install:
	cd frontend && npm install

.PHONY: frontend-dev
frontend-dev:
	cd frontend && npm run dev

.PHONY: frontend-build
frontend-build:
	cd frontend && npm run build

.PHONY: frontend-lint
frontend-lint:
	cd frontend && npm run lint
