#!/bin/bash

# Test runner script for iwork-converter
# This script runs all tests with various configurations

set -e

echo "========================================="
echo "iWork Converter Test Suite"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print section headers
print_header() {
    echo -e "${GREEN}>>> $1${NC}"
    echo ""
}

# Function to print warnings
print_warning() {
    echo -e "${YELLOW}WARNING: $1${NC}"
    echo ""
}

# Function to print errors
print_error() {
    echo -e "${RED}ERROR: $1${NC}"
    echo ""
}

# Get project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go first."
    exit 1
fi

echo "Go version: $(go version)"
echo ""

# Parse arguments
RUN_BENCHMARKS=false
RUN_COVERAGE=false
RUN_SHORT=false
VERBOSE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -b|--bench)
            RUN_BENCHMARKS=true
            shift
            ;;
        -c|--coverage)
            RUN_COVERAGE=true
            shift
            ;;
        -s|--short)
            RUN_SHORT=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  -b, --bench      Run benchmarks"
            echo "  -c, --coverage   Generate coverage report"
            echo "  -s, --short      Run only unit tests (skip integration tests)"
            echo "  -v, --verbose    Verbose output"
            echo "  -h, --help       Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                  # Run all tests"
            echo "  $0 -s               # Run only unit tests"
            echo "  $0 -c               # Run tests with coverage"
            echo "  $0 -b               # Run benchmarks"
            echo "  $0 -v -c            # Run tests with verbose output and coverage"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Check for test files
TEST_FILES_PRESENT=false
if [ -f "a.pages" ] || [ -f "a.numbers" ] || [ -f "a.key" ]; then
    TEST_FILES_PRESENT=true
fi

if [ "$TEST_FILES_PRESENT" = false ] && [ "$RUN_SHORT" = false ]; then
    print_warning "No test files (a.pages, a.numbers, a.key) found in root directory."
    print_warning "Integration tests will be skipped."
    echo ""
fi

# Build test flags
TEST_FLAGS=""
if [ "$VERBOSE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -v"
fi

if [ "$RUN_SHORT" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -short"
    print_header "Running Unit Tests Only (Short Mode)"
else
    print_header "Running All Tests (Unit + Integration)"
fi

# Run tests
echo "Test flags: $TEST_FLAGS"
echo ""

if go test $TEST_FLAGS ./...; then
    echo ""
    echo -e "${GREEN}✓ All tests passed${NC}"
    echo ""
else
    echo ""
    print_error "Some tests failed"
    exit 1
fi

# Run tests with coverage if requested
if [ "$RUN_COVERAGE" = true ]; then
    print_header "Generating Coverage Report"
    
    COVERAGE_FLAGS="-coverprofile=coverage.out -covermode=atomic"
    if [ "$RUN_SHORT" = true ]; then
        COVERAGE_FLAGS="$COVERAGE_FLAGS -short"
    fi
    
    go test $COVERAGE_FLAGS ./...
    
    echo ""
    echo "Coverage summary:"
    go tool cover -func=coverage.out | tail -1
    
    echo ""
    echo "Generating HTML coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    
    echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"
    echo ""
fi

# Run benchmarks if requested
if [ "$RUN_BENCHMARKS" = true ]; then
    print_header "Running Benchmarks"
    
    if go test -bench=. -benchmem ./...; then
        echo ""
        echo -e "${GREEN}✓ Benchmarks completed${NC}"
        echo ""
    else
        echo ""
        print_error "Benchmark execution failed"
        exit 1
    fi
fi

# Summary
echo "========================================="
echo "Test Summary"
echo "========================================="
echo ""
echo "Unit tests:        ✓ Passed"
if [ "$RUN_SHORT" = false ]; then
    if [ "$TEST_FILES_PRESENT" = true ]; then
        echo "Integration tests: ✓ Passed"
    else
        echo "Integration tests: ⊘ Skipped (no test files)"
    fi
fi
if [ "$RUN_COVERAGE" = true ]; then
    echo "Coverage report:   ✓ Generated (coverage.html)"
fi
if [ "$RUN_BENCHMARKS" = true ]; then
    echo "Benchmarks:        ✓ Completed"
fi
echo ""
echo -e "${GREEN}All requested tests completed successfully!${NC}"

