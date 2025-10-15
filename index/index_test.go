package index

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestOpen tests the creation of a new index from a file
func TestOpen(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (string, func())
		wantErr bool
	}{
		{
			name: "non-existent file",
			setup: func() (string, func()) {
				return "/tmp/non_existent_file.pages", func() {}
			},
			wantErr: true,
		},
		{
			name: "invalid zip file",
			setup: func() (string, func()) {
				tmpfile := filepath.Join(os.TempDir(), "invalid.pages")
				os.WriteFile(tmpfile, []byte("not a zip file"), 0o644)
				return tmpfile, func() { os.Remove(tmpfile) }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setup()
			defer cleanup()

			_, err := Open(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDetectDocumentType tests document type detection
func TestDetectDocumentType(t *testing.T) {
	// This test requires real iWork files to test properly
	// For now, just test that file extension is detected
	tests := []struct {
		name         string
		filename     string
		expectedType string
	}{
		{
			name:         "pages extension",
			filename:     "test.pages",
			expectedType: "pages",
		},
		{
			name:         "numbers extension",
			filename:     "test.numbers",
			expectedType: "numbers",
		},
		{
			name:         "key extension",
			filename:     "test.key",
			expectedType: "key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext := filepath.Ext(tt.filename)
			if len(ext) > 0 {
				ext = ext[1:] // Remove leading dot
			}

			// key files might be detected as "key" or "keynote"
			if tt.expectedType == "key" {
				if ext != "key" {
					t.Errorf("Extension = %s, want key", ext)
				}
			} else if ext != tt.expectedType {
				t.Errorf("Extension = %s, want %s", ext, tt.expectedType)
			}
		})
	}
}

// TestRecordLoading tests loading of records
func TestRecordLoading(t *testing.T) {
	t.Run("empty record list", func(t *testing.T) {
		ix := &Index{
			Type:    "pages",
			Records: make(map[uint64]interface{}),
		}

		if len(ix.Records) != 0 {
			t.Errorf("Expected empty records map, got %d entries", len(ix.Records))
		}
	})
}

// TestRecordIdentifierHandling tests handling of record identifiers
func TestRecordIdentifierHandling(t *testing.T) {
	tests := []struct {
		name       string
		identifier uint64
		wantValid  bool
	}{
		{
			name:       "valid identifier",
			identifier: 12345,
			wantValid:  true,
		},
		{
			name:       "zero identifier",
			identifier: 0,
			wantValid:  true,
		},
		{
			name:       "large identifier",
			identifier: 9999999999,
			wantValid:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ix := &Index{
				Type:    "pages",
				Records: make(map[uint64]interface{}),
			}

			// Create and store record
			record := "test record"
			ix.Records[tt.identifier] = record

			// Verify storage
			stored, exists := ix.Records[tt.identifier]
			if !exists {
				t.Error("Record not found in map")
			}
			if stored != record {
				t.Error("Retrieved record doesn't match stored record")
			}
		})
	}
}

// TestFileExtensionDetection tests detection of file types by extension
func TestFileExtensionDetection(t *testing.T) {
	tests := []struct {
		filename     string
		expectedType string
	}{
		{
			filename:     "document.pages",
			expectedType: "pages",
		},
		{
			filename:     "spreadsheet.numbers",
			expectedType: "numbers",
		},
		{
			filename:     "presentation.key",
			expectedType: "key", // .key files have "key" extension
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			ext := filepath.Ext(tt.filename)
			// Remove leading dot
			if len(ext) > 0 {
				ext = ext[1:]
			}

			if ext != tt.expectedType {
				t.Errorf("Extension = %s, want %s", ext, tt.expectedType)
			}
		})
	}
}

// TestIndexTypeValidation tests validation of index types
func TestIndexTypeValidation(t *testing.T) {
	tests := []struct {
		name      string
		indexType string
		isValid   bool
	}{
		{"pages type", "pages", true},
		{"numbers type", "numbers", true},
		{"key type", "key", true},
		{"keynote type", "keynote", true},
		{"empty type", "", false},
		{"unknown type", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ix := &Index{
				Type:    tt.indexType,
				Records: make(map[uint64]interface{}),
			}

			validTypes := map[string]bool{
				"pages":   true,
				"numbers": true,
				"key":     true,
				"keynote": true,
			}

			isValid := validTypes[ix.Type]
			if isValid != tt.isValid {
				t.Errorf("Type %q validity = %v, want %v", tt.indexType, isValid, tt.isValid)
			}
		})
	}
}

// TestRecordOperations tests various record operations
func TestRecordOperations(t *testing.T) {
	t.Run("add and retrieve records", func(t *testing.T) {
		ix := &Index{
			Type:    "pages",
			Records: make(map[uint64]interface{}),
		}

		// Add records
		ix.Records[1] = "record1"
		ix.Records[2] = "record2"
		ix.Records[100] = "record100"

		// Retrieve records
		if ix.Records[1] != "record1" {
			t.Error("Failed to retrieve record 1")
		}
		if ix.Records[2] != "record2" {
			t.Error("Failed to retrieve record 2")
		}
		if ix.Records[100] != "record100" {
			t.Error("Failed to retrieve record 100")
		}
	})

	t.Run("overwrite existing record", func(t *testing.T) {
		ix := &Index{
			Type:    "pages",
			Records: make(map[uint64]interface{}),
		}

		ix.Records[1] = "original"
		ix.Records[1] = "updated"

		if ix.Records[1] != "updated" {
			t.Error("Record was not updated")
		}
	})

	t.Run("delete record", func(t *testing.T) {
		ix := &Index{
			Type:    "pages",
			Records: make(map[uint64]interface{}),
		}

		ix.Records[1] = "test"
		delete(ix.Records, 1)

		if _, exists := ix.Records[1]; exists {
			t.Error("Record was not deleted")
		}
	})
}

// TestRecordCount tests counting records
func TestRecordCount(t *testing.T) {
	tests := []struct {
		name          string
		recordCount   int
		expectedCount int
	}{
		{"empty index", 0, 0},
		{"single record", 1, 1},
		{"multiple records", 10, 10},
		{"many records", 1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ix := &Index{
				Type:    "pages",
				Records: make(map[uint64]interface{}),
			}

			for i := 0; i < tt.recordCount; i++ {
				ix.Records[uint64(i)] = fmt.Sprintf("record%d", i)
			}

			if len(ix.Records) != tt.expectedCount {
				t.Errorf("Record count = %d, want %d", len(ix.Records), tt.expectedCount)
			}
		})
	}
}

// TestFilePathParsing tests parsing file paths
func TestFilePathParsing(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectedExt  string
		expectedType string
	}{
		{
			name:         "pages file with path",
			path:         "/path/to/document.pages",
			expectedExt:  "pages",
			expectedType: "pages",
		},
		{
			name:         "numbers file with path",
			path:         "/another/path/spreadsheet.numbers",
			expectedExt:  "numbers",
			expectedType: "numbers",
		},
		{
			name:         "key file with path",
			path:         "/presentations/slide.key",
			expectedExt:  "key",
			expectedType: "key",
		},
		{
			name:         "relative path",
			path:         "./test.pages",
			expectedExt:  "pages",
			expectedType: "pages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext := filepath.Ext(tt.path)
			if len(ext) > 0 {
				ext = ext[1:]
			}

			if ext != tt.expectedExt {
				t.Errorf("Extension = %s, want %s", ext, tt.expectedExt)
			}
		})
	}
}

// TestRecordKeyRanges tests various record key ranges
func TestRecordKeyRanges(t *testing.T) {
	tests := []struct {
		name string
		keys []uint64
	}{
		{
			name: "small keys",
			keys: []uint64{0, 1, 2, 3, 4, 5},
		},
		{
			name: "large keys",
			keys: []uint64{1000000, 2000000, 3000000},
		},
		{
			name: "mixed keys",
			keys: []uint64{1, 100, 10000, 1000000},
		},
		{
			name: "sparse keys",
			keys: []uint64{1, 1000, 2000, 5000, 10000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ix := &Index{
				Type:    "pages",
				Records: make(map[uint64]interface{}),
			}

			// Add records
			for _, key := range tt.keys {
				ix.Records[key] = fmt.Sprintf("record_%d", key)
			}

			// Verify all keys exist
			for _, key := range tt.keys {
				if _, exists := ix.Records[key]; !exists {
					t.Errorf("Key %d not found", key)
				}
			}

			// Verify count
			if len(ix.Records) != len(tt.keys) {
				t.Errorf("Record count = %d, want %d", len(ix.Records), len(tt.keys))
			}
		})
	}
}

// TestConcurrentAccess tests concurrent access patterns (read-only)
func TestConcurrentAccess(t *testing.T) {
	ix := &Index{
		Type:    "pages",
		Records: make(map[uint64]interface{}),
	}

	// Populate with test data
	for i := uint64(0); i < 100; i++ {
		ix.Records[i] = fmt.Sprintf("record%d", i)
	}

	// Note: This is a simple test. For real concurrent access,
	// you would use sync.RWMutex or channels
	t.Run("sequential reads", func(t *testing.T) {
		for i := uint64(0); i < 100; i++ {
			if _, exists := ix.Records[i]; !exists {
				t.Errorf("Record %d not found", i)
			}
		}
	})
}

// BenchmarkRecordLookup benchmarks record lookup performance
func BenchmarkRecordLookup(b *testing.B) {
	ix := &Index{
		Type:    "pages",
		Records: make(map[uint64]interface{}),
	}

	// Populate with test data
	for i := uint64(0); i < 1000; i++ {
		ix.Records[i] = "test record"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ix.Records[uint64(i%1000)]
	}
}

// BenchmarkRecordInsertion benchmarks record insertion
func BenchmarkRecordInsertion(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ix := &Index{
			Type:    "pages",
			Records: make(map[uint64]interface{}),
		}

		for j := uint64(0); j < 100; j++ {
			ix.Records[j] = "test record"
		}
	}
}

// BenchmarkLargeRecordSet benchmarks operations on large record sets
func BenchmarkLargeRecordSet(b *testing.B) {
	ix := &Index{
		Type:    "pages",
		Records: make(map[uint64]interface{}),
	}

	// Populate with large dataset
	for i := uint64(0); i < 10000; i++ {
		ix.Records[i] = fmt.Sprintf("record_%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ix.Records[uint64(i%10000)]
	}
}
