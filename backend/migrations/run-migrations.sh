#!/bin/bash

# TRADEMASTER Database Migration Runner
# Applies migrations to PostgreSQL database

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# Load environment variables
if [ -f "../.env" ]; then
    export $(cat ../.env | grep -v '^#' | xargs)
fi

# Default database URL
DB_URL=${DATABASE_URL:-"postgresql://trademaster:password@localhost:5432/trademaster"}

echo "🗄️  TRADEMASTER Database Migration"
echo ""

# Parse command
COMMAND=${1:-up}

case $COMMAND in
    up)
        print_info "Applying migrations..."

        # Apply migrations in order
        for file in *_*.up.sql; do
            if [ -f "$file" ]; then
                echo "  Running: $file"
                psql "$DB_URL" -f "$file" > /dev/null 2>&1
                if [ $? -eq 0 ]; then
                    print_success "Applied $file"
                else
                    print_error "Failed to apply $file"
                    exit 1
                fi
            fi
        done

        echo ""
        print_success "All migrations applied successfully"
        ;;

    down)
        print_info "Rolling back migrations..."

        # Rollback in reverse order
        for file in $(ls -r *_*.down.sql); do
            if [ -f "$file" ]; then
                echo "  Running: $file"
                psql "$DB_URL" -f "$file" > /dev/null 2>&1
                if [ $? -eq 0 ]; then
                    print_success "Rolled back $file"
                else
                    print_error "Failed to rollback $file"
                    exit 1
                fi
            fi
        done

        echo ""
        print_success "All migrations rolled back"
        ;;

    reset)
        print_info "Resetting database (down then up)..."

        # Run down
        ./run-migrations.sh down

        # Run up
        ./run-migrations.sh up

        print_success "Database reset complete"
        ;;

    status)
        print_info "Checking database connection..."

        if psql "$DB_URL" -c "SELECT 1" > /dev/null 2>&1; then
            print_success "Database connection successful"

            # Check if tables exist
            TABLE_COUNT=$(psql "$DB_URL" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")
            echo "  Tables: $TABLE_COUNT"
        else
            print_error "Database connection failed"
            exit 1
        fi
        ;;

    *)
        echo "Usage: $0 {up|down|reset|status}"
        echo ""
        echo "Commands:"
        echo "  up     - Apply all migrations"
        echo "  down   - Rollback all migrations"
        echo "  reset  - Rollback and reapply all migrations"
        echo "  status - Check database connection and status"
        exit 1
        ;;
esac
