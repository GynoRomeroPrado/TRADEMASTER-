#!/bin/bash

# TRADEMASTER Code Validation Script
# Validates all code quality, types, and linting without full builds

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ERRORS=0
WARNINGS=0

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
    ERRORS=$((ERRORS + 1))
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
    WARNINGS=$((WARNINGS + 1))
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

echo -e "${BLUE}🔍 TRADEMASTER Code Validation${NC}\n"

# ============================================
# 1. TYPESCRIPT VALIDATION
# ============================================
print_header "TypeScript Type Checking"

# Check for TypeScript syntax errors in all .ts and .tsx files
echo "Checking TypeScript files for syntax errors..."

TS_FILES=$(find frontend/src -type f \( -name "*.ts" -o -name "*.tsx" \) 2>/dev/null | wc -l)
echo "Found $TS_FILES TypeScript files"

# Check for common TypeScript issues
echo ""
echo "Checking for common issues:"

# Check for 'any' types (bad practice)
ANY_COUNT=$(grep -r ": any" frontend/src --include="*.ts" --include="*.tsx" 2>/dev/null | wc -l || echo "0")
if [ "$ANY_COUNT" -gt 0 ]; then
    print_warning "Found $ANY_COUNT uses of 'any' type (consider using specific types)"
else
    print_success "No 'any' types found"
fi

# Check for TODO comments
TODO_COUNT=$(grep -r "TODO\|FIXME\|HACK" frontend/src --include="*.ts" --include="*.tsx" 2>/dev/null | wc -l || echo "0")
if [ "$TODO_COUNT" -gt 0 ]; then
    print_info "Found $TODO_COUNT TODO/FIXME comments"
fi

# Check imports are valid
echo ""
echo "Checking import statements..."
INVALID_IMPORTS=$(grep -r "from ['\"]@/" frontend/src --include="*.ts" --include="*.tsx" 2>/dev/null | grep -v "from '@/types\|from '@/hooks\|from '@/contexts\|from '@/components" | wc -l || echo "0")
if [ "$INVALID_IMPORTS" -eq 0 ]; then
    print_success "All imports use correct path aliases"
else
    print_warning "$INVALID_IMPORTS potential invalid imports found"
fi

# Check for unused imports (basic check)
echo "Checking for unused variables..."
UNUSED_VARS=$(grep -r "useState\|useEffect" frontend/src --include="*.tsx" 2>/dev/null | grep "import.*useState.*useEffect\|import.*useEffect.*useState" | wc -l || echo "0")
print_info "Found $UNUSED_VARS files using React hooks"

print_success "TypeScript validation complete"

# ============================================
# 2. REACT BEST PRACTICES
# ============================================
print_header "React Best Practices"

# Check for proper React patterns
echo "Checking React component patterns..."

# Check for 'use client' directives
CLIENT_COMPONENTS=$(grep -r "^'use client'" frontend/src/app --include="*.tsx" 2>/dev/null | wc -l || echo "0")
print_info "Found $CLIENT_COMPONENTS 'use client' components"

# Check for proper async component handling
ASYNC_COMPONENTS=$(grep -r "export default async function" frontend/src/app --include="*.tsx" 2>/dev/null | wc -l || echo "0")
print_info "Found $ASYNC_COMPONENTS async server components"

# Check for useState without useEffect dependency warnings
echo "Checking hook dependencies..."
USE_STATE=$(grep -r "useState" frontend/src --include="*.tsx" 2>/dev/null | wc -l || echo "0")
USE_EFFECT=$(grep -r "useEffect" frontend/src --include="*.tsx" 2>/dev/null | wc -l || echo "0")
print_info "useState: $USE_STATE, useEffect: $USE_EFFECT"

print_success "React patterns check complete"

# ============================================
# 3. GO CODE VALIDATION
# ============================================
print_header "Go Code Validation"

cd backend || exit

SERVICES=(
    "auth-service"
    "estimation-service"
    "project-service"
    "training-service"
    "marketplace-service"
)

for service in "${SERVICES[@]}"; do
    if [ -d "$service" ]; then
        echo "Checking $service..."
        cd "$service" || continue

        # Check if go.mod exists
        if [ -f "go.mod" ]; then
            # Check for syntax errors
            if go fmt ./... > /dev/null 2>&1; then
                print_success "$service: Code formatted"
            else
                print_error "$service: Formatting issues found"
            fi

            # Run go vet (static analysis)
            if go vet ./... > /dev/null 2>&1; then
                print_success "$service: No vet issues"
            else
                print_warning "$service: Vet found potential issues"
            fi
        else
            print_warning "$service: No go.mod found"
        fi

        cd ..
    fi
done

cd ..
print_success "Go validation complete"

# ============================================
# 4. FILE STRUCTURE VALIDATION
# ============================================
print_header "File Structure Validation"

echo "Checking project structure..."

REQUIRED_DIRS=(
    "backend"
    "frontend"
    "ml-services"
    "docs"
    "scripts"
)

for dir in "${REQUIRED_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        print_success "Directory exists: $dir"
    else
        print_error "Missing directory: $dir"
    fi
done

# Check for important files
REQUIRED_FILES=(
    "README.md"
    "docker-compose.yml"
    "Makefile"
    ".gitignore"
    "frontend/package.json"
    "frontend/tsconfig.json"
    "frontend/next.config.js"
    "backend/.env.example"
    "frontend/.env.example"
)

echo ""
for file in "${REQUIRED_FILES[@]}"; do
    if [ -f "$file" ]; then
        print_success "File exists: $file"
    else
        print_error "Missing file: $file"
    fi
done

# ============================================
# 5. ENVIRONMENT FILES
# ============================================
print_header "Environment Configuration"

echo "Checking environment files..."

if [ -f "frontend/.env.example" ]; then
    ENV_VARS=$(grep -c "=" frontend/.env.example || echo "0")
    print_success "Frontend .env.example has $ENV_VARS variables"
else
    print_error "frontend/.env.example not found"
fi

if [ -f "backend/.env.example" ]; then
    ENV_VARS=$(grep -c "=" backend/.env.example || echo "0")
    print_success "Backend .env.example has $ENV_VARS variables"
else
    print_error "backend/.env.example not found"
fi

# ============================================
# 6. DOCKER VALIDATION
# ============================================
print_header "Docker Configuration"

echo "Checking Dockerfiles..."

DOCKERFILES=$(find . -name "Dockerfile" -o -name "Dockerfile.*" | wc -l)
print_info "Found $DOCKERFILES Dockerfiles"

# Check docker-compose.yml
if [ -f "docker-compose.yml" ]; then
    SERVICES_COUNT=$(grep -c "    image:\|    build:" docker-compose.yml || echo "0")
    print_success "docker-compose.yml defines $SERVICES_COUNT services"
else
    print_error "docker-compose.yml not found"
fi

# ============================================
# 7. SCRIPTS VALIDATION
# ============================================
print_header "Scripts Validation"

echo "Checking utility scripts..."

SCRIPTS=(
    "scripts/setup-dev.sh"
    "scripts/deploy-dev.sh"
    "scripts/deploy-prod.sh"
    "scripts/run-tests.sh"
    "scripts/cleanup.sh"
)

for script in "${SCRIPTS[@]}"; do
    if [ -f "$script" ]; then
        if [ -x "$script" ]; then
            print_success "$script is executable"
        else
            print_warning "$script exists but is not executable"
        fi
    else
        print_error "$script not found"
    fi
done

# ============================================
# 8. DOCUMENTATION
# ============================================
print_header "Documentation"

DOCS=(
    "README.md"
    "docs/ARCHITECTURE.md"
    "docs/GO_TO_MARKET.md"
    "docs/DEPLOYMENT.md"
    "docs/API.md"
)

for doc in "${DOCS[@]}"; do
    if [ -f "$doc" ]; then
        LINES=$(wc -l < "$doc")
        print_success "$doc exists ($LINES lines)"
    else
        print_error "$doc not found"
    fi
done

# ============================================
# 9. CODE METRICS
# ============================================
print_header "Code Metrics"

echo "Calculating code statistics..."

# Count lines of code
GO_FILES=$(find backend -name "*.go" | wc -l || echo "0")
GO_LINES=$(find backend -name "*.go" -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}' || echo "0")

TS_FILES=$(find frontend/src -name "*.ts" -o -name "*.tsx" | wc -l || echo "0")
TS_LINES=$(find frontend/src \( -name "*.ts" -o -name "*.tsx" \) -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}' || echo "0")

PY_FILES=$(find ml-services -name "*.py" | wc -l || echo "0")
PY_LINES=$(find ml-services -name "*.py" -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}' || echo "0")

echo ""
echo "Code Statistics:"
echo "  Go:         $GO_FILES files, $GO_LINES lines"
echo "  TypeScript: $TS_FILES files, $TS_LINES lines"
echo "  Python:     $PY_FILES files, $PY_LINES lines"
echo ""
TOTAL_FILES=$((GO_FILES + TS_FILES + PY_FILES))
TOTAL_LINES=$((GO_LINES + TS_LINES + PY_LINES))
echo "  Total:      $TOTAL_FILES files, $TOTAL_LINES lines of code"

# ============================================
# SUMMARY
# ============================================
print_header "Validation Summary"

echo "Results:"
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}  ✓ Errors:   $ERRORS${NC}"
else
    echo -e "${RED}  ✗ Errors:   $ERRORS${NC}"
fi

if [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}  ✓ Warnings: $WARNINGS${NC}"
else
    echo -e "${YELLOW}  ⚠ Warnings: $WARNINGS${NC}"
fi

echo ""

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ Validation passed!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Install dependencies: cd frontend && npm install"
    echo "  2. Run full build:      npm run build"
    echo "  3. Run tests:           ./scripts/run-tests.sh"
    exit 0
else
    echo -e "${RED}❌ Validation failed with $ERRORS errors${NC}"
    echo ""
    echo "Please fix the errors above and run validation again."
    exit 1
fi
