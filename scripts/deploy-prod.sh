#!/bin/bash

echo "=== TRADEMASTER Production Deployment ==="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check environment
if [ "$ENVIRONMENT" != "production" ]; then
  echo -e "${RED}Error: ENVIRONMENT must be set to 'production'${NC}"
  exit 1
fi

# Confirmation
echo -e "${YELLOW}WARNING: You are about to deploy to PRODUCTION${NC}"
read -p "Are you sure? (type 'yes' to continue) " -r
if [ "$REPLY" != "yes" ]; then
  echo "Deployment cancelled"
  exit 1
fi

# Backup database
echo -e "${GREEN}Creating database backup...${NC}"
./scripts/backup-db.sh

# Pull latest code
echo -e "${GREEN}Pulling latest stable code...${NC}"
git checkout main
git pull origin main

# Run tests
echo -e "${GREEN}Running tests...${NC}"
make test
if [ $? -ne 0 ]; then
  echo -e "${RED}Tests failed! Aborting deployment.${NC}"
  exit 1
fi

# Build Docker images
echo -e "${GREEN}Building Docker images...${NC}"
docker-compose -f docker-compose.prod.yml build

# Run migrations
echo -e "${GREEN}Running database migrations...${NC}"
docker-compose -f docker-compose.prod.yml run --rm backend-migrate

# Deploy with zero-downtime
echo -e "${GREEN}Deploying services...${NC}"
docker-compose -f docker-compose.prod.yml up -d --no-deps --build

# Health checks
echo -e "${GREEN}Performing health checks...${NC}"
sleep 10

# Check services
for service in frontend backend-api ml-services; do
  if docker-compose -f docker-compose.prod.yml ps | grep -q "$service.*Up"; then
    echo -e "${GREEN}✓ $service is running${NC}"
  else
    echo -e "${RED}✗ $service failed to start${NC}"
    docker-compose -f docker-compose.prod.yml logs $service
    exit 1
  fi
done

# Clean up old images
echo -e "${GREEN}Cleaning up old images...${NC}"
docker image prune -f

echo -e "${GREEN}=== Production Deployment Complete! ===${NC}"
echo ""
echo "Monitoring: docker-compose -f docker-compose.prod.yml logs -f"
echo "Rollback: ./scripts/rollback.sh"
