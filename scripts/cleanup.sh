#!/bin/bash

# TRADEMASTER Cleanup Script
# Cleans up development environment, Docker containers, and build artifacts

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

echo "🧹 TRADEMASTER Cleanup Script"
echo ""

# Ask for confirmation
read -p "This will stop all containers and remove volumes. Continue? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cleanup cancelled"
    exit 0
fi

echo ""
echo "Stopping Docker containers..."
docker-compose down -v

if [ $? -eq 0 ]; then
    print_success "Docker containers stopped and removed"
else
    print_error "Failed to stop Docker containers"
fi

echo ""
echo "Removing build artifacts..."

# Remove Go build artifacts
echo "  - Backend build artifacts"
find backend -type f -name "*.exe" -delete
find backend -type f -name "*.test" -delete
find backend -type f -name "coverage.out" -delete
find backend -type f -name "coverage.html" -delete
print_success "Backend build artifacts removed"

# Remove frontend build artifacts
echo "  - Frontend build artifacts"
if [ -d "frontend/.next" ]; then
    rm -rf frontend/.next
fi
if [ -d "frontend/out" ]; then
    rm -rf frontend/out
fi
if [ -d "frontend/coverage" ]; then
    rm -rf frontend/coverage
fi
print_success "Frontend build artifacts removed"

# Remove ML service artifacts
echo "  - ML service artifacts"
find ml-services -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
find ml-services -type d -name ".pytest_cache" -exec rm -rf {} + 2>/dev/null || true
find ml-services -type d -name "htmlcov" -exec rm -rf {} + 2>/dev/null || true
find ml-services -type f -name "*.pyc" -delete 2>/dev/null || true
find ml-services -type f -name ".coverage" -delete 2>/dev/null || true
print_success "ML service artifacts removed"

# Remove log files
echo ""
echo "Removing log files..."
if [ -d "logs" ]; then
    rm -rf logs/*
    print_success "Log files removed"
fi

# Remove data directories (ask for confirmation)
echo ""
read -p "Remove local data directories (postgres, mongodb, redis)? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -d "data" ]; then
        rm -rf data/*
        print_success "Data directories removed"
    fi
else
    print_warning "Data directories preserved"
fi

# Remove node_modules (ask for confirmation)
echo ""
read -p "Remove node_modules? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -d "frontend/node_modules" ]; then
        rm -rf frontend/node_modules
        print_success "node_modules removed"
    fi
else
    print_warning "node_modules preserved"
fi

# Remove Docker images (ask for confirmation)
echo ""
read -p "Remove Docker images? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker-compose down --rmi all
    print_success "Docker images removed"
else
    print_warning "Docker images preserved"
fi

echo ""
echo -e "${GREEN}✅ Cleanup complete!${NC}"
echo ""
echo "To set up the development environment again, run:"
echo "  ./scripts/setup-dev.sh"
echo ""
