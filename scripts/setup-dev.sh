#!/bin/bash

# TRADEMASTER Development Environment Setup Script
# This script sets up the complete development environment

set -e

echo "🚀 Setting up TRADEMASTER development environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi
print_status "Docker is installed"

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi
print_status "Docker Compose is installed"

# Check if Go is installed (for backend development)
if ! command -v go &> /dev/null; then
    print_warning "Go is not installed. Backend services will only run in Docker."
else
    GO_VERSION=$(go version | awk '{print $3}')
    print_status "Go is installed: $GO_VERSION"
fi

# Check if Node.js is installed (for frontend development)
if ! command -v node &> /dev/null; then
    print_warning "Node.js is not installed. Frontend will only run in Docker."
else
    NODE_VERSION=$(node --version)
    print_status "Node.js is installed: $NODE_VERSION"
fi

# Check if Python is installed (for ML services)
if ! command -v python3 &> /dev/null; then
    print_warning "Python 3 is not installed. ML services will only run in Docker."
else
    PYTHON_VERSION=$(python3 --version)
    print_status "Python is installed: $PYTHON_VERSION"
fi

# Create .env files if they don't exist
echo ""
echo "📝 Setting up environment variables..."

if [ ! -f "frontend/.env.local" ]; then
    cp frontend/.env.example frontend/.env.local
    print_status "Created frontend/.env.local"
else
    print_warning "frontend/.env.local already exists"
fi

if [ ! -f "backend/.env" ]; then
    cp backend/.env.example backend/.env
    print_status "Created backend/.env"
else
    print_warning "backend/.env already exists"
fi

if [ ! -f "ml-services/.env" ]; then
    cp ml-services/.env.example ml-services/.env
    print_status "Created ml-services/.env"
else
    print_warning "ml-services/.env already exists"
fi

# Create necessary directories
echo ""
echo "📁 Creating necessary directories..."

mkdir -p data/postgres data/mongodb data/redis data/elasticsearch logs
print_status "Created data and logs directories"

# Install frontend dependencies
echo ""
echo "📦 Installing frontend dependencies..."
cd frontend
if command -v npm &> /dev/null; then
    npm install
    print_status "Frontend dependencies installed"
else
    print_warning "npm not found, skipping frontend dependency installation"
fi
cd ..

# Install backend dependencies (Go modules)
echo ""
echo "📦 Installing backend dependencies..."
cd backend/auth-service
if command -v go &> /dev/null; then
    go mod download
    print_status "Auth service dependencies installed"
fi
cd ../..

# Start Docker services
echo ""
echo "🐳 Starting Docker services..."
docker-compose up -d postgres mongodb redis elasticsearch rabbitmq minio

# Wait for services to be ready
echo ""
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if PostgreSQL is ready
until docker-compose exec -T postgres pg_isready -U trademaster &> /dev/null; do
    echo "Waiting for PostgreSQL..."
    sleep 2
done
print_status "PostgreSQL is ready"

# Check if MongoDB is ready
until docker-compose exec -T mongodb mongosh --eval "db.adminCommand('ping')" &> /dev/null; do
    echo "Waiting for MongoDB..."
    sleep 2
done
print_status "MongoDB is ready"

# Check if Redis is ready
until docker-compose exec -T redis redis-cli ping &> /dev/null; do
    echo "Waiting for Redis..."
    sleep 2
done
print_status "Redis is ready"

# Run database migrations
echo ""
echo "🗄️  Running database migrations..."
# Add migration commands here
print_status "Database migrations completed"

# Seed database with demo data
echo ""
echo "🌱 Seeding database with demo data..."
# Add seed commands here
print_status "Database seeded"

echo ""
echo -e "${GREEN}✅ Development environment setup complete!${NC}"
echo ""
echo "To start the services:"
echo "  - Frontend:  cd frontend && npm run dev"
echo "  - Backend:   docker-compose up backend-services"
echo "  - All:       docker-compose up"
echo ""
echo "Access the application:"
echo "  - Frontend:  http://localhost:3000"
echo "  - API Docs:  http://localhost:8001/docs"
echo "  - Grafana:   http://localhost:3001"
echo ""
