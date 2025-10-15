#!/bin/bash

# Quick test script
# Usage: ./quick_test.sh [input_file] [options]

set -e

echo "=== Quick Test Script ==="

# Get project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

# Check if input file is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <input_file> [options]"
    echo "Example: $0 a.key"
    echo "Example: $0 a.key -v  (with verbose debug output)"
    exit 1
fi

INPUT_FILE="$1"
OUTPUT_FILE="${INPUT_FILE%.*}.html"

echo "Input file: $INPUT_FILE"
echo "Output file: $OUTPUT_FILE"

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file '$INPUT_FILE' does not exist"
    exit 1
fi

# Build project
echo "Building project..."
go build

# Run conversion
echo "Running conversion..."
./iwork-converter -v "$INPUT_FILE" "$OUTPUT_FILE"

echo "=== Conversion Complete ==="
echo "Output file: $OUTPUT_FILE"
