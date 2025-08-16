# Newstrix Makefile
.PHONY: help build clean test run-api run-fetcher run-embedder docker-build docker-run migrate lint format

# Variables
BINARY_DIR=bin
API_BINARY=$(BINARY_DIR)/api
FETCHER_BINARY=$(BINARY_DIR)/fetcher
EMBEDDER_BINARY=$(BINARY_DIR)/embedder

# Default target
help: ## Show this help message
	@echo "Newstrix - AI-Powered News Aggregator"
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build targets
build: clean ## Build all binaries
	@echo "Building binaries..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(API_BINARY) ./cmd/api
	@go build -o $(FETCHER_BINARY) ./cmd/fetcher
	@go build -o $(EMBEDDER_BINARY) ./cmd/embedder
	@echo "Build completed!"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)
	@go clean

# Run targets
run-api: ## Run API service
	@echo "Starting API service..."
	@go run ./cmd/api

run-fetcher: ## Run fetcher service
	@echo "Starting fetcher service..."
	@go run ./cmd/fetcher

run-embedder: ## Run embedder service
	@echo "Starting embedder service..."
	@go run ./cmd/embedder

# Docker targets
docker-build: ## Build Docker images
	@echo "Building Docker images..."
	@docker build -f Dockerfile.api -t newstrix-api .
	@docker build -f Dockerfile.fetcher -t newstrix-fetcher .
	@echo "Docker images built!"

docker-run: ## Run with Docker Compose
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

docker-stop: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	@docker-compose down

# Database targets
goose-install: ## Install goose migration tool
	@echo "Installing goose migration tool..."
	@go install github.com/pressly/goose/v3/cmd/goose@latest

migrate: ## Run database migrations
	@echo "Running database migrations..."
	@goose -dir ./migrations postgres "$(shell grep POSTGRES_URL .env | cut -d '=' -f2)" up

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@goose -dir ./migrations postgres "$(shell grep POSTGRES_URL .env | cut -d '=' -f2)" down

migrate-status: ## Show migration status
	@echo "Migration status:"
	@goose -dir ./migrations postgres "$(shell grep POSTGRES_URL .env | cut -d '=' -f2)" status

migrate-create: ## Create new migration file
	@echo "Creating new migration file..."
	@goose -dir ./migrations postgres "$(shell grep POSTGRES_URL .env | cut -d '=' -f2)" create $(name) sql

# Development targets
test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@go vet ./...

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Development setup
dev-setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Development environment ready!"

# Quick start
quick-start: build migrate ## Quick start with local binaries and migrations
	@echo "Starting Newstrix..."
	@echo "API: http://localhost:8080"
	@echo "Fetcher: running in background"
	@$(FETCHER_BINARY) &
	@$(API_BINARY)

# Production build
prod-build: ## Build production binaries
	@echo "Building production binaries..."
	@mkdir -p $(BINARY_DIR)
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(API_BINARY) ./cmd/api
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(FETCHER_BINARY) ./cmd/fetcher
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(EMBEDDER_BINARY) ./cmd/embedder
	@echo "Production build completed!"

# Health check
health: ## Check service health
	@echo "Checking service health..."
	@curl -f http://localhost:8080/ || echo "API service is down"
	@echo "Health check completed!"

# Logs
logs: ## Show Docker Compose logs
	@docker-compose logs -f

# Clean everything
clean-all: clean ## Clean everything including Docker
	@echo "Cleaning everything..."
	@docker-compose down -v --remove-orphans
	@docker system prune -f
	@echo "Everything cleaned!"