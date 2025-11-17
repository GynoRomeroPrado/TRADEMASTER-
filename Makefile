.PHONY: help install dev build test clean migrate seed docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Installing mobile dependencies..."
	cd mobile && npm install
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing ML dependencies..."
	cd ml-services && pip install -r requirements.txt
	@echo "Done!"

dev: ## Start development environment
	docker-compose up -d postgres mongodb redis elasticsearch rabbitmq
	@echo "Starting backend services..."
	cd backend && make dev &
	@echo "Starting frontend..."
	cd frontend && npm run dev &
	@echo "Starting ML services..."
	cd ml-services && python main.py &
	@echo "Development environment ready!"

build: ## Build all services
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Building backend services..."
	cd backend && make build
	@echo "Building ML services..."
	cd ml-services && docker build -t trademaster-ml .
	@echo "Build complete!"

test: ## Run all tests
	@echo "Running frontend tests..."
	cd frontend && npm test
	@echo "Running backend tests..."
	cd backend && make test
	@echo "Running ML tests..."
	cd ml-services && pytest
	@echo "All tests passed!"

clean: ## Clean build artifacts
	rm -rf frontend/.next frontend/out
	rm -rf backend/*/bin
	rm -rf ml-services/__pycache__ ml-services/.pytest_cache
	@echo "Cleaned!"

migrate: ## Run database migrations
	cd backend && make migrate

seed: ## Seed database with sample data
	cd backend && make seed

docker-up: ## Start all Docker services
	docker-compose up -d

docker-down: ## Stop all Docker services
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

proto: ## Generate protobuf files
	cd backend && make proto

lint: ## Run linters
	cd frontend && npm run lint
	cd backend && make lint
	cd ml-services && flake8 . && black --check .

format: ## Format code
	cd frontend && npm run format
	cd backend && make format
	cd ml-services && black .

deploy-dev: ## Deploy to development
	./scripts/deploy-dev.sh

deploy-prod: ## Deploy to production
	./scripts/deploy-prod.sh
