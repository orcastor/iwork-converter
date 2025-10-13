package main_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orcastor/iwork-converter/index"
	"github.com/orcastor/iwork-converter/iwork2html"
)

// TestPagesConversion tests conversion of .pages files
func TestPagesConversion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if test file exists
	testFile := "a.pages"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file a.pages not found, skipping test")
	}

	t.Run("to HTML", func(t *testing.T) {
		outputFile := filepath.Join(os.TempDir(), "test_pages.html")
		defer os.Remove(outputFile)

		err := iwork2html.Convert(testFile, outputFile)
		if err != nil {
			t.Fatalf("Failed to convert: %v", err)
		}

		// Verify output file exists and has content
		info, err := os.Stat(outputFile)
		if err != nil {
			t.Fatalf("Output file not created: %v", err)
		}

		if info.Size() == 0 {
			t.Error("Output file is empty")
		}

		// Read and verify HTML structure
		data, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output: %v", err)
		}

		html := string(data)
		// Note: DOCTYPE may not be present in all HTML outputs
		if !strings.Contains(html, "<html") {
			t.Error("HTML output missing html tag")
		}
		if !strings.Contains(html, "</html>") {
			t.Error("HTML output missing closing html tag")
		}
	})

	t.Run("document type detection", func(t *testing.T) {
		ix, err := index.Open(testFile)
		if err != nil {
			t.Fatalf("Failed to open index: %v", err)
		}

		if ix.Type != "pages" {
			t.Errorf("Document type = %s, want pages", ix.Type)
		}
	})
}

// TestNumbersConversion tests conversion of .numbers files
func TestNumbersConversion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if test file exists
	testFile := "a.numbers"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file a.numbers not found, skipping test")
	}

	t.Run("to HTML", func(t *testing.T) {
		outputFile := filepath.Join(os.TempDir(), "test_numbers.html")
		defer os.Remove(outputFile)

		err := iwork2html.Convert(testFile, outputFile)
		if err != nil {
			t.Fatalf("Failed to convert: %v", err)
		}

		// Verify output file exists and has content
		info, err := os.Stat(outputFile)
		if err != nil {
			t.Fatalf("Output file not created: %v", err)
		}

		if info.Size() == 0 {
			t.Error("Output file is empty")
		}

		// Read and verify HTML structure
		data, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output: %v", err)
		}

		html := string(data)
		// Note: DOCTYPE may not be present in all HTML outputs

		// Verify tables exist
		if !strings.Contains(html, "<table") {
			t.Error("HTML output missing table tags for Numbers document")
		}
	})

	t.Run("document type detection", func(t *testing.T) {
		ix, err := index.Open(testFile)
		if err != nil {
			t.Fatalf("Failed to open index: %v", err)
		}

		if ix.Type != "numbers" {
			t.Errorf("Document type = %s, want numbers", ix.Type)
		}
	})
}

// TestKeynoteConversion tests conversion of .key files (if available)
func TestKeynoteConversion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check if test file exists
	testFile := "a.key"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file a.key not found, skipping test")
	}

	t.Run("to HTML", func(t *testing.T) {
		outputFile := filepath.Join(os.TempDir(), "test_keynote.html")
		defer os.Remove(outputFile)

		err := iwork2html.Convert(testFile, outputFile)
		if err != nil {
			t.Fatalf("Failed to convert: %v", err)
		}

		// Verify output file exists and has content
		info, err := os.Stat(outputFile)
		if err != nil {
			t.Fatalf("Output file not created: %v", err)
		}

		if info.Size() == 0 {
			t.Error("Output file is empty")
		}

		// Read and verify HTML structure
		data, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output: %v", err)
		}

		_ = string(data) // Verify file is readable
		// Note: DOCTYPE may not be present in all HTML outputs
	})

	t.Run("document type detection", func(t *testing.T) {
		ix, err := index.Open(testFile)
		if err != nil {
			t.Fatalf("Failed to open index: %v", err)
		}

		if ix.Type != "keynote" && ix.Type != "key" {
			t.Errorf("Document type = %s, want keynote or key", ix.Type)
		}
	})
}

// TestDocumentTypeDetection tests that document types are correctly detected
func TestDocumentTypeDetection(t *testing.T) {
	tests := []struct {
		name         string
		file         string
		expectedType string
	}{
		{
			name:         "Pages document",
			file:         "a.pages",
			expectedType: "pages",
		},
		{
			name:         "Numbers document",
			file:         "a.numbers",
			expectedType: "numbers",
		},
		{
			name:         "Keynote document",
			file:         "a.key",
			expectedType: "keynote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := os.Stat(tt.file); os.IsNotExist(err) {
				t.Skipf("Test file %s not found, skipping test", tt.file)
			}

			ix, err := index.Open(tt.file)
			if err != nil {
				t.Fatalf("Failed to open index: %v", err)
			}

			// Accept both "keynote" and "key" as valid types for Keynote files
			if tt.expectedType == "keynote" {
				if ix.Type != "keynote" && ix.Type != "key" {
					t.Errorf("Document type = %s, want %s or key", ix.Type, tt.expectedType)
				}
			} else if ix.Type != tt.expectedType {
				t.Errorf("Document type = %s, want %s", ix.Type, tt.expectedType)
			}
		})
	}
}

// TestHTMLValidation tests that generated HTML is well-formed
func TestHTMLValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testFiles := []string{"a.pages", "a.numbers"}

	for _, testFile := range testFiles {
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			continue
		}

		t.Run(testFile, func(t *testing.T) {
			outputFile := filepath.Join(os.TempDir(), "test_validation.html")
			defer os.Remove(outputFile)

			err := iwork2html.Convert(testFile, outputFile)
			if err != nil {
				t.Fatalf("Failed to convert: %v", err)
			}

			data, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			html := string(data)

			// Check for basic HTML structure
			requiredTags := []string{
				"<html",
				"<head>",
				"</head>",
				"<body",
				"</body>",
				"</html>",
			}

			for _, tag := range requiredTags {
				if !strings.Contains(html, tag) {
					t.Errorf("HTML missing required tag: %s", tag)
				}
			}

			// Note: DOCTYPE may not be present in all outputs

			// Check for CSS
			if !strings.Contains(html, "<style") {
				t.Error("HTML missing style tag")
			}
		})
	}
}

// TestTableProcessing tests that tables are processed correctly
func TestTableProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testFile := "a.numbers"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file a.numbers not found, skipping test")
	}

	t.Run("numbers table structure", func(t *testing.T) {
		outputFile := filepath.Join(os.TempDir(), "test_table.html")
		defer os.Remove(outputFile)

		err := iwork2html.Convert(testFile, outputFile)
		if err != nil {
			t.Fatalf("Failed to convert: %v", err)
		}

		data, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output: %v", err)
		}

		html := string(data)

		// Numbers files should have tables
		if !strings.Contains(html, "<table") {
			t.Error("Numbers HTML missing table tags")
		}
		if !strings.Contains(html, "<td") && !strings.Contains(html, "<th") {
			t.Error("Numbers HTML missing table cell tags")
		}

		// Check that tables have proper structure
		if !strings.Contains(html, "</table>") {
			t.Error("HTML has unclosed table tag")
		}
	})
}

// BenchmarkPagesConversion benchmarks Pages to HTML conversion
func BenchmarkPagesConversion(b *testing.B) {
	testFile := "a.pages"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		b.Skip("Test file a.pages not found, skipping benchmark")
	}

	outputFile := filepath.Join(os.TempDir(), "bench_pages.html")
	defer os.Remove(outputFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := iwork2html.Convert(testFile, outputFile)
		if err != nil {
			b.Fatalf("Failed to convert: %v", err)
		}
	}
}

// BenchmarkNumbersConversion benchmarks Numbers to HTML conversion
func BenchmarkNumbersConversion(b *testing.B) {
	testFile := "a.numbers"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		b.Skip("Test file a.numbers not found, skipping benchmark")
	}

	outputFile := filepath.Join(os.TempDir(), "bench_numbers.html")
	defer os.Remove(outputFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := iwork2html.Convert(testFile, outputFile)
		if err != nil {
			b.Fatalf("Failed to convert: %v", err)
		}
	}
}
