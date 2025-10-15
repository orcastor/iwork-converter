# Testing Guide for iWork Converter

This document provides comprehensive information about the test suite for the iWork Converter project.

## Table of Contents

- [Test Structure](#test-structure)
- [Running Tests](#running-tests)
- [Test Coverage](#test-coverage)
- [Test Files](#test-files)
- [Writing Tests](#writing-tests)
- [Continuous Integration](#continuous-integration)

## Test Structure

The project uses Go's built-in testing framework and includes:

1. **Unit Tests**: Test individual functions and components in isolation
2. **Integration Tests**: Test the full conversion pipeline with real iWork files
3. **Benchmarks**: Performance testing for critical operations

### Test Organization

```
iwork-converter/
├── integration_test.go          # Integration tests (root package)
├── run_tests.sh                 # Test runner script
├── testdata/                    # Test data directory
│   └── README.md                # Test data documentation
├── index/
│   └── index_test.go            # Index package unit tests
├── iwork2html/
│   └── iwork2html_test.go       # HTML conversion unit tests
└── iwork2text/
    └── iwork2text_test.go       # Text conversion unit tests
```

## Running Tests

### Quick Start

Run all unit tests (excluding integration tests):
```bash
./run_tests.sh -s
```

Run all tests including integration tests:
```bash
./run_tests.sh
```

### Using Go Test Directly

Run all unit tests:
```bash
go test -short ./...
```

Run all tests including integration tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./iwork2html
go test ./iwork2text
go test ./index
```

Run tests with verbose output:
```bash
go test -v ./...
```

### Test Runner Script Options

The `run_tests.sh` script provides several options:

```bash
./run_tests.sh [options]

Options:
  -b, --bench      Run benchmarks
  -c, --coverage   Generate coverage report
  -s, --short      Run only unit tests (skip integration tests)
  -v, --verbose    Verbose output
  -h, --help       Show help message
```

Examples:
```bash
# Run unit tests with coverage
./run_tests.sh -s -c

# Run all tests with verbose output
./run_tests.sh -v

# Run benchmarks
./run_tests.sh -b

# Combine options
./run_tests.sh -v -c -b
```

## Test Coverage

### Generating Coverage Reports

Generate coverage report:
```bash
./run_tests.sh -c
```

Or manually:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

View coverage in terminal:
```bash
go test -cover ./...
```

View detailed coverage by function:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Current Coverage

The test suite currently covers:
- Core HTML generation functions
- Cell offset parsing and decoding
- Binary encoding/decoding
- Document type detection
- File handling and I/O operations

## Test Files

### Unit Tests

#### `iwork2html/iwork2html_test.go`
Tests for HTML conversion functionality:
- `TestElementCreation`: HTML element creation (E function)
- `TestTextNodeCreation`: Text node creation (T function)
- `TestDebugMode`: Debug mode toggling
- `TestCellOffsetParsing`: Cell offset encoding/decoding
- `TestEmptyCellMarker`: Empty cell detection (65535 marker)
- `TestHTMLNodeStructure`: Parent-child relationships in HTML nodes

#### `iwork2text/iwork2text_test.go`
Tests for text conversion functionality:
- `TestCellOffsetHandling`: Cell offset handling for text extraction
- `TestEmptyCellDetection`: Empty cell detection
- `TestBinaryEncoding`: Binary data encoding/decoding

#### `index/index_test.go`
Tests for index/parsing functionality:
- `TestOpen`: Index creation from files
- `TestDetectDocumentType`: Document type detection
- `TestRecordLoading`: Record loading and management
- `TestRecordIdentifierHandling`: Record identifier operations
- `TestFileExtensionDetection`: File extension detection

### Integration Tests

#### `integration_test.go`
End-to-end tests with real iWork files:
- `TestPagesConversion`: Pages document conversion
- `TestNumbersConversion`: Numbers spreadsheet conversion
- `TestKeynoteConversion`: Keynote presentation conversion
- `TestDocumentTypeDetection`: Type detection for all formats
- `TestHTMLValidation`: HTML output validation
- `TestTableProcessing`: Table structure processing

**Note**: Integration tests require test files (`a.pages`, `a.numbers`, `a.key`) to be present in the project root. They are automatically skipped if files are not found.

### Benchmarks

Performance benchmarks are included for critical operations:
- `BenchmarkElementCreation`: HTML element creation speed
- `BenchmarkTextNodeCreation`: Text node creation speed
- `BenchmarkCellOffsetDecoding`: Cell offset decoding performance
- `BenchmarkPagesConversion`: Full Pages document conversion
- `BenchmarkNumbersConversion`: Full Numbers document conversion

Run benchmarks:
```bash
go test -bench=. -benchmem ./...
```

## Writing Tests

### Test Naming Conventions

- Test functions: `TestXxx`
- Benchmark functions: `BenchmarkXxx`
- Example functions: `ExampleXxx`

### Test Structure

```go
func TestFeatureName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
    }{
        {
            name:     "descriptive test case name",
            input:    testInput,
            expected: expectedOutput,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionUnderTest(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Best Practices

1. **Use Table-Driven Tests**: Group related test cases together
2. **Test Edge Cases**: Empty inputs, nil values, boundary conditions
3. **Clear Test Names**: Use descriptive names that explain what is being tested
4. **Isolate Tests**: Each test should be independent
5. **Test Error Handling**: Include tests for error conditions
6. **Use Subtests**: Use `t.Run()` for better organization
7. **Skip When Appropriate**: Use `t.Skip()` for tests requiring specific conditions

### Adding Integration Tests

Integration tests should:
1. Check if test files exist
2. Skip gracefully if files are not present
3. Use the `-short` flag to skip in CI environments
4. Clean up temporary files

Example:
```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    testFile := "test.pages"
    if _, err := os.Stat(testFile); os.IsNotExist(err) {
        t.Skip("Test file not found")
    }

    // Test implementation...
}
```

## Continuous Integration

### GitHub Actions (Recommended)

Create `.github/workflows/test.yml`:

```yaml
name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run tests
      run: go test -short -v ./...
    
    - name: Generate coverage
      run: go test -short -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

### Pre-commit Hook

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
echo "Running tests..."
go test -short ./...
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

## Troubleshooting

### Common Issues

#### Tests Fail Due to Missing Files
- Integration tests require real iWork files
- Run with `-short` flag to skip: `go test -short ./...`

#### Import Errors
- Ensure module name matches: `github.com/orcastor/iwork-converter`
- Run `go mod tidy` to fix dependencies

#### Coverage Too Low
- Focus on testing exported functions first
- Add edge case tests
- Test error handling paths

#### Slow Tests
- Use `-short` flag for quick feedback
- Run integration tests separately
- Use benchmarks to identify bottlenecks

## Additional Resources

- [Go Testing Package Documentation](https://pkg.go.dev/testing)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Table Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

## Contributing

When contributing tests:
1. Follow existing test patterns
2. Include both positive and negative test cases
3. Update this documentation if adding new test types
4. Ensure all tests pass before submitting PR
5. Aim for meaningful test coverage, not just high percentages

---

For questions or issues with the test suite, please open an issue on the project repository.

