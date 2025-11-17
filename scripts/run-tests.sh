#!/bin/bash

# TRADEMASTER Test Runner Script
# Runs all tests across backend, frontend, and ML services

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

print_header() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Parse command line arguments
RUN_BACKEND=true
RUN_FRONTEND=true
RUN_ML=true
COVERAGE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --backend-only)
            RUN_FRONTEND=false
            RUN_ML=false
            shift
            ;;
        --frontend-only)
            RUN_BACKEND=false
            RUN_ML=false
            shift
            ;;
        --ml-only)
            RUN_BACKEND=false
            RUN_FRONTEND=false
            shift
            ;;
        --coverage)
            COVERAGE=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--backend-only|--frontend-only|--ml-only] [--coverage]"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🧪 TRADEMASTER Test Suite${NC}\n"

# Backend Tests (Go)
if [ "$RUN_BACKEND" = true ]; then
    print_header "Backend Services Tests (Go)"

    BACKEND_SERVICES=(
        "auth-service"
        "estimation-service"
        "project-service"
        "training-service"
        "marketplace-service"
    )

    for service in "${BACKEND_SERVICES[@]}"; do
        echo "Testing $service..."
        cd "backend/$service"

        if [ "$COVERAGE" = true ]; then
            if go test -v -cover -coverprofile=coverage.out ./...; then
                print_success "$service tests passed"
                PASSED_TESTS=$((PASSED_TESTS + 1))

                # Generate coverage report
                go tool cover -html=coverage.out -o coverage.html
                COVERAGE_PERCENT=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
                echo "  Coverage: $COVERAGE_PERCENT"
            else
                print_error "$service tests failed"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
        else
            if go test -v ./...; then
                print_success "$service tests passed"
                PASSED_TESTS=$((PASSED_TESTS + 1))
            else
                print_error "$service tests failed"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
        fi

        TOTAL_TESTS=$((TOTAL_TESTS + 1))
        cd ../..
        echo ""
    done
fi

# Frontend Tests (Jest/React Testing Library)
if [ "$RUN_FRONTEND" = true ]; then
    print_header "Frontend Tests (Jest)"

    cd frontend

    if [ "$COVERAGE" = true ]; then
        if npm test -- --coverage --watchAll=false; then
            print_success "Frontend tests passed"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            print_error "Frontend tests failed"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    else
        if npm test -- --watchAll=false; then
            print_success "Frontend tests passed"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            print_error "Frontend tests failed"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    fi

    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    cd ..
    echo ""
fi

# ML Services Tests (pytest)
if [ "$RUN_ML" = true ]; then
    print_header "ML Services Tests (Python/pytest)"

    ML_SERVICES=(
        "computer-vision"
        "cost-prediction"
        "training-coach"
        "inventory-optimization"
    )

    for service in "${ML_SERVICES[@]}"; do
        echo "Testing $service..."
        cd "ml-services/$service"

        if [ "$COVERAGE" = true ]; then
            if python -m pytest --cov=src --cov-report=html --cov-report=term; then
                print_success "$service tests passed"
                PASSED_TESTS=$((PASSED_TESTS + 1))
            else
                print_error "$service tests failed"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
        else
            if python -m pytest -v; then
                print_success "$service tests passed"
                PASSED_TESTS=$((PASSED_TESTS + 1))
            else
                print_error "$service tests failed"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
        fi

        TOTAL_TESTS=$((TOTAL_TESTS + 1))
        cd ../..
        echo ""
    done
fi

# Print summary
print_header "Test Summary"

echo "Total test suites: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"

if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"
    echo ""
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
else
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"
    echo ""
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
fi
