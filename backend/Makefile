.PHONY: proto install dev test lint clean migrate

proto: ## Generate protobuf files
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

install: ## Install dependencies
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

dev: ## Run all services
	cd estimation-service && go run cmd/server/main.go &
	cd project-service && go run cmd/server/main.go &
	cd training-service && go run cmd/server/main.go &
	cd marketplace-service && go run cmd/server/main.go &
	cd auth-service && go run cmd/server/main.go

build: ## Build all services
	cd estimation-service && go build -o bin/estimation-service cmd/server/main.go
	cd project-service && go build -o bin/project-service cmd/server/main.go
	cd training-service && go build -o bin/training-service cmd/server/main.go
	cd marketplace-service && go build -o bin/marketplace-service cmd/server/main.go
	cd auth-service && go build -o bin/auth-service cmd/server/main.go

test: ## Run tests
	go test -v ./...

lint: ## Run linter
	golangci-lint run ./...

format: ## Format code
	gofmt -w .

clean: ## Clean build artifacts
	rm -rf */bin/

migrate: ## Run migrations
	cd estimation-service && go run migrations/migrate.go
	cd project-service && go run migrations/migrate.go
	cd training-service && go run migrations/migrate.go
	cd marketplace-service && go run migrations/migrate.go
	cd auth-service && go run migrations/migrate.go

seed: ## Seed database
	go run scripts/seed.go
