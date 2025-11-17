# TRADEMASTER Backend Services

Go-based microservices architecture with gRPC communication.

## Services

1. **estimation-service**: Project estimation and blueprint analysis
2. **project-service**: Project and task management
3. **training-service**: Training courses and certifications
4. **marketplace-service**: Job marketplace and matching
5. **auth-service**: Authentication and authorization

## Tech Stack

- Go 1.21+
- gRPC + Protocol Buffers
- PostgreSQL (primary database)
- MongoDB (document storage)
- Redis (caching)
- RabbitMQ (message queue)

## Development

```bash
# Install dependencies
make install

# Generate protobuf files
make proto

# Run all services
make dev

# Run specific service
cd estimation-service && go run cmd/server/main.go

# Run tests
make test

# Run linter
make lint
```

## Project Structure

```
backend/
├── estimation-service/
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── api/             # gRPC handlers
│   │   ├── service/         # Business logic
│   │   ├── repository/      # Database access
│   │   └── ml/              # ML client
│   ├── pkg/                 # Public packages
│   ├── proto/               # Protobuf definitions
│   └── migrations/          # DB migrations
├── shared/                  # Shared code
│   ├── auth/                # Auth utilities
│   ├── db/                  # DB connections
│   └── utils/               # Common utils
└── Makefile
```
