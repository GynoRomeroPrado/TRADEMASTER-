#!/bin/bash

echo "=== TRADEMASTER Development Deployment ==="

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
  echo -e "${YELLOW}Docker is not running. Please start Docker and try again.${NC}"
  exit 1
fi

# Pull latest code
echo -e "${GREEN}Pulling latest code...${NC}"
git pull origin develop

# Start infrastructure
echo -e "${GREEN}Starting infrastructure services...${NC}"
docker-compose up -d postgres mongodb redis elasticsearch rabbitmq

# Wait for databases
echo -e "${GREEN}Waiting for databases to be ready...${NC}"
sleep 10

# Run migrations
echo -e "${GREEN}Running database migrations...${NC}"
cd backend && make migrate
cd ..

# Seed database (optional)
read -p "Seed database with demo data? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  echo -e "${GREEN}Seeding database...${NC}"
  cd backend && make seed
  cd ..
fi

# Build backend services
echo -e "${GREEN}Building backend services...${NC}"
cd backend && make build
cd ..

# Install frontend dependencies
echo -e "${GREEN}Installing frontend dependencies...${NC}"
cd frontend && npm install
cd ..

# Start all services
echo -e "${GREEN}Starting all services...${NC}"
make dev

echo -e "${GREEN}=== Deployment Complete! ===${NC}"
echo ""
echo "Services available at:"
echo "  Frontend: http://localhost:3000"
echo "  API Gateway: http://localhost:8080"
echo "  ML Services: http://localhost:8001-8004"
echo ""
echo "To stop all services: make docker-down"
