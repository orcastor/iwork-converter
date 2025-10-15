# Test Data Directory

This directory contains test files for the iwork-converter project.

## Test Files

Place your iWork test files here:

- `.pages` files - Apple Pages documents
- `.numbers` files - Apple Numbers spreadsheets  
- `.key` files - Apple Keynote presentations

## Running Tests

### Unit Tests

Run all unit tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./iwork2html
go test ./iwork2text
go test ./index
```

### Integration Tests

Run integration tests (requires test files in root directory):
```bash
go test -v
```

Skip integration tests (run only unit tests):
```bash
go test -short ./...
```

### Benchmarks

Run benchmarks:
```bash
go test -bench=. -benchmem
```

Run benchmarks for a specific package:
```bash
go test -bench=. -benchmem ./iwork2html
```

## Test Coverage

Generate test coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Generate coverage for specific package:
```bash
go test -coverprofile=coverage.out ./iwork2html
go tool cover -html=coverage.out
```

## Adding New Tests

When adding new test files:

1. Place test `.pages`, `.numbers`, or `.key` files in this directory
2. Update integration tests if needed
3. Ensure all tests pass before committing

## Test File Naming

Test files should follow Go conventions:
- `*_test.go` - Test files
- Functions starting with `Test` - Test functions
- Functions starting with `Benchmark` - Benchmark functions
- Functions starting with `Example` - Example functions

