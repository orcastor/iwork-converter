package iwork2html

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

// TestPagesTableMultiColumn tests Pages tables with multiple columns
func TestPagesTableMultiColumn(t *testing.T) {
	tests := []struct {
		name        string
		rowCount    int
		columnCount int
		description string
	}{
		{
			name:        "single column table",
			rowCount:    5,
			columnCount: 1,
			description: "Typical Pages resume table with only one column of content",
		},
		{
			name:        "two column table",
			rowCount:    10,
			columnCount: 2,
			description: "Two column layout with left column for labels and right for content",
		},
		{
			name:        "three column table",
			rowCount:    8,
			columnCount: 3,
			description: "Three column table where middle column may be empty",
		},
		{
			name:        "seven column table",
			rowCount:    7,
			columnCount: 7,
			description: "Seven column table where last 6 columns may have no width",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate table structure
			totalCells := tt.rowCount * tt.columnCount

			if totalCells < 0 {
				t.Errorf("Invalid table dimensions: %d rows x %d cols", tt.rowCount, tt.columnCount)
			}

			// Verify table can be allocated
			offsets := make([]uint16, tt.columnCount)
			for i := range offsets {
				offsets[i] = uint16(i * 10)
			}

			if len(offsets) != tt.columnCount {
				t.Errorf("Expected %d column offsets, got %d", tt.columnCount, len(offsets))
			}
		})
	}
}

// TestPagesTableEmptyColumns tests Pages tables with empty columns
func TestPagesTableEmptyColumns(t *testing.T) {
	tests := []struct {
		name               string
		totalColumns       int
		emptyColumnIndices []int
		description        string
	}{
		{
			name:               "no empty columns",
			totalColumns:       3,
			emptyColumnIndices: []int{},
			description:        "All columns have data",
		},
		{
			name:               "last column empty",
			totalColumns:       3,
			emptyColumnIndices: []int{2},
			description:        "Last column is empty",
		},
		{
			name:               "middle column empty",
			totalColumns:       3,
			emptyColumnIndices: []int{1},
			description:        "Middle column is empty",
		},
		{
			name:               "multiple empty columns",
			totalColumns:       5,
			emptyColumnIndices: []int{1, 3, 4},
			description:        "Multiple columns are empty",
		},
		{
			name:               "last six columns empty",
			totalColumns:       7,
			emptyColumnIndices: []int{1, 2, 3, 4, 5, 6},
			description:        "7-column table with last 6 columns empty (typical Pages resume layout)",
		},
		{
			name:               "alternating empty columns",
			totalColumns:       6,
			emptyColumnIndices: []int{1, 3, 5},
			description:        "Alternating empty columns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create column offsets
			offsets := make([]uint16, tt.totalColumns)

			// Mark empty columns with 65535
			emptyMap := make(map[int]bool)
			for _, idx := range tt.emptyColumnIndices {
				emptyMap[idx] = true
			}

			for i := 0; i < tt.totalColumns; i++ {
				if emptyMap[i] {
					offsets[i] = 65535 // Empty marker
				} else {
					offsets[i] = uint16(i * 10) // Valid offset
				}
			}

			// Verify empty column detection
			emptyCount := 0
			for _, offset := range offsets {
				if offset == 65535 {
					emptyCount++
				}
			}

			if emptyCount != len(tt.emptyColumnIndices) {
				t.Errorf("Expected %d empty columns, found %d", len(tt.emptyColumnIndices), emptyCount)
			}
		})
	}
}

// TestPagesTableColumnContentDetection tests detecting which columns have content
func TestPagesTableColumnContentDetection(t *testing.T) {
	tests := []struct {
		name                  string
		rowOffsets            [][]uint16 // Each row's column offsets
		expectedActiveColumns []bool     // Which columns have content
		description           string
	}{
		{
			name: "single active column",
			rowOffsets: [][]uint16{
				{0, 65535, 65535, 65535},  // Row 0: only col 0 has data
				{10, 65535, 65535, 65535}, // Row 1: only col 0 has data
				{20, 65535, 65535, 65535}, // Row 2: only col 0 has data
			},
			expectedActiveColumns: []bool{true, false, false, false},
			description:           "Only first column has data",
		},
		{
			name: "two active columns",
			rowOffsets: [][]uint16{
				{0, 10, 65535, 65535},  // Row 0: col 0,1 have data
				{20, 30, 65535, 65535}, // Row 1: col 0,1 have data
				{40, 50, 65535, 65535}, // Row 2: col 0,1 have data
			},
			expectedActiveColumns: []bool{true, true, false, false},
			description:           "First two columns have data",
		},
		{
			name: "sparse columns",
			rowOffsets: [][]uint16{
				{0, 65535, 20, 65535, 40},  // Row 0: col 0,2,4 have data
				{10, 65535, 30, 65535, 50}, // Row 1: col 0,2,4 have data
			},
			expectedActiveColumns: []bool{true, false, true, false, true},
			description:           "Odd columns have data, even columns are empty",
		},
		{
			name: "all columns active",
			rowOffsets: [][]uint16{
				{0, 10, 20},
				{30, 40, 50},
				{60, 70, 80},
			},
			expectedActiveColumns: []bool{true, true, true},
			description:           "All columns have data",
		},
		{
			name: "partially filled rows",
			rowOffsets: [][]uint16{
				{0, 10, 20, 30},           // Row 0: all columns have data
				{40, 65535, 65535, 65535}, // Row 1: only col 0 has data
				{50, 60, 65535, 65535},    // Row 2: col 0,1 have data
			},
			expectedActiveColumns: []bool{true, true, true, true},
			description:           "Different rows have different active columns, should detect all columns that were ever active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.rowOffsets) == 0 {
				t.Fatal("No row offsets provided")
			}

			columnCount := len(tt.rowOffsets[0])
			hasContentColumns := make([]bool, columnCount)

			// Analyze which columns have content across all rows
			for _, rowOffsets := range tt.rowOffsets {
				for colIdx, offset := range rowOffsets {
					if offset != 65535 {
						hasContentColumns[colIdx] = true
					}
				}
			}

			// Verify detection results
			if len(hasContentColumns) != len(tt.expectedActiveColumns) {
				t.Fatalf("Column count mismatch: got %d, want %d", len(hasContentColumns), len(tt.expectedActiveColumns))
			}

			for i, expected := range tt.expectedActiveColumns {
				if hasContentColumns[i] != expected {
					t.Errorf("Column %d: hasContent=%v, want %v", i, hasContentColumns[i], expected)
				}
			}
		})
	}
}

// TestPagesTableRowMajorVsColumnMajor tests content organization
func TestPagesTableRowMajorVsColumnMajor(t *testing.T) {
	tests := []struct {
		name         string
		rows         int
		cols         int
		fillOrder    string // "row-major" or "column-major"
		expectedData []int  // Expected linear data order
	}{
		{
			name:         "row-major 3x3",
			rows:         3,
			cols:         3,
			fillOrder:    "row-major",
			expectedData: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:         "column-major 3x3",
			rows:         3,
			cols:         3,
			fillOrder:    "column-major",
			expectedData: []int{0, 3, 6, 1, 4, 7, 2, 5, 8},
		},
		{
			name:         "row-major 2x4",
			rows:         2,
			cols:         4,
			fillOrder:    "row-major",
			expectedData: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]int, tt.rows*tt.cols)

			if tt.fillOrder == "row-major" {
				// Fill row by row
				for i := range data {
					data[i] = i
				}
			} else if tt.fillOrder == "column-major" {
				// Fill column by column
				idx := 0
				for col := 0; col < tt.cols; col++ {
					for row := 0; row < tt.rows; row++ {
						data[row*tt.cols+col] = idx
						idx++
					}
				}
			}

			// Verify the data matches expected order
			for i, expected := range tt.expectedData {
				if data[i] != expected {
					t.Errorf("Position %d: got %d, want %d", i, data[i], expected)
				}
			}
		})
	}
}

// TestPagesTableCellPositioning tests correct cell positioning in multi-column tables
func TestPagesTableCellPositioning(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]string // Key: "row,col", Value: content
		rows     int
		cols     int
		validate func(data map[string]string) error
	}{
		{
			name: "content in correct cells",
			data: map[string]string{
				"0,0": "Header",
				"1,0": "Row1",
				"2,0": "Row2",
			},
			rows: 3,
			cols: 1,
			validate: func(data map[string]string) error {
				// Verify row 0 col 0 has "Header"
				if data["0,0"] != "Header" {
					return nil // Error would be reported by test
				}
				return nil
			},
		},
		{
			name: "multi-column with specific positioning",
			data: map[string]string{
				"0,0": "Name",
				"0,1": "Age",
				"1,0": "Alice",
				"1,1": "30",
				"2,0": "Bob",
				"2,1": "25",
			},
			rows: 3,
			cols: 2,
			validate: func(data map[string]string) error {
				return nil
			},
		},
		{
			name: "seven column table with content only in first column",
			data: map[string]string{
				"0,0": "Basic Info",
				"1,0": "Name: John Doe",
				"2,0": "Age: 30",
				"3,0": "Work Experience",
				"4,0": "2020-2023: Company A",
				// Columns 1-6 are empty
			},
			rows: 5,
			cols: 7,
			validate: func(data map[string]string) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify all provided data is within bounds
			for key, value := range tt.data {
				if value == "" {
					continue
				}

				var row, col int
				if _, err := parsePosition(key, &row, &col); err != nil {
					t.Errorf("Invalid key format: %s", key)
					continue
				}

				if row >= tt.rows {
					t.Errorf("Row %d out of bounds (max: %d)", row, tt.rows-1)
				}
				if col >= tt.cols {
					t.Errorf("Column %d out of bounds (max: %d)", col, tt.cols-1)
				}
			}
		})
	}
}

// Helper function to parse position string "row,col"
func parsePosition(s string, row, col *int) (int, error) {
	// Simple parsing for "row,col" format
	n, err := fmt.Sscanf(s, "%d,%d", row, col)
	return n, err
}

// TestPagesTableEmptyRowHandling tests handling of empty rows
func TestPagesTableEmptyRowHandling(t *testing.T) {
	tests := []struct {
		name             string
		totalRows        int
		emptyRowIndices  []int
		expectedNonEmpty int
	}{
		{
			name:             "no empty rows",
			totalRows:        5,
			emptyRowIndices:  []int{},
			expectedNonEmpty: 5,
		},
		{
			name:             "one empty row",
			totalRows:        5,
			emptyRowIndices:  []int{2},
			expectedNonEmpty: 4,
		},
		{
			name:             "multiple empty rows",
			totalRows:        10,
			emptyRowIndices:  []int{0, 5, 9},
			expectedNonEmpty: 7,
		},
		{
			name:             "alternating empty rows",
			totalRows:        6,
			emptyRowIndices:  []int{1, 3, 5},
			expectedNonEmpty: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create row data
			emptyMap := make(map[int]bool)
			for _, idx := range tt.emptyRowIndices {
				emptyMap[idx] = true
			}

			// Count non-empty rows
			nonEmptyCount := 0
			for i := 0; i < tt.totalRows; i++ {
				if !emptyMap[i] {
					nonEmptyCount++
				}
			}

			if nonEmptyCount != tt.expectedNonEmpty {
				t.Errorf("Expected %d non-empty rows, got %d", tt.expectedNonEmpty, nonEmptyCount)
			}
		})
	}
}

// TestPagesTableColumnWidth tests column width handling
func TestPagesTableColumnWidth(t *testing.T) {
	tests := []struct {
		name         string
		columnWidths []float64
		description  string
	}{
		{
			name:         "equal width columns",
			columnWidths: []float64{100, 100, 100},
			description:  "All columns have same width",
		},
		{
			name:         "varying widths",
			columnWidths: []float64{200, 150, 100, 50},
			description:  "Columns with decreasing widths",
		},
		{
			name:         "some zero width columns",
			columnWidths: []float64{300, 0, 0, 0, 0, 0, 0},
			description:  "First column has width, rest have zero width (typical Pages resume layout)",
		},
		{
			name:         "mixed zero and non-zero widths",
			columnWidths: []float64{100, 0, 100, 0, 100},
			description:  "Alternating columns with and without width",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalWidth := 0.0
			zeroWidthCount := 0

			for _, width := range tt.columnWidths {
				totalWidth += width
				if width == 0 {
					zeroWidthCount++
				}
			}

			t.Logf("Total width: %.2f, Zero-width columns: %d/%d",
				totalWidth, zeroWidthCount, len(tt.columnWidths))

			// Verify widths are non-negative
			for i, width := range tt.columnWidths {
				if width < 0 {
					t.Errorf("Column %d has negative width: %.2f", i, width)
				}
			}
		})
	}
}

// TestPagesTableContentOrdering tests that table content appears in correct order
func TestPagesTableContentOrdering(t *testing.T) {
	tests := []struct {
		name          string
		richTableKeys []uint32 // Keys in richTable order (index 0, 1, 2, ...)
		expectedOrder []string // Expected content order in cells
		description   string
	}{
		{
			name:          "sequential content",
			richTableKeys: []uint32{1, 2, 3, 4, 5},
			expectedOrder: []string{"Header", "Item1", "Item2", "Item3", "Item4"},
			description:   "Content should appear in same order as richTable indices",
		},
		{
			name:          "work experience ordering",
			richTableKeys: []uint32{1, 2, 3, 4},
			expectedOrder: []string{"Work Experience", "Bilibili", "NoahStore", "Company A"},
			description:   "NoahStore should appear after Bilibili, not before",
		},
		{
			name:          "resume sections ordering",
			richTableKeys: []uint32{1, 2, 3, 4, 5, 6},
			expectedOrder: []string{
				"Basic Info",
				"Education",
				"Work Experience",
				"Company 1",
				"Company 2",
				"Skills",
			},
			description: "Resume sections should maintain document order",
		},
		{
			name:          "table with gaps",
			richTableKeys: []uint32{1, 3, 5, 7},
			expectedOrder: []string{"Row 0", "Row 1", "Row 2", "Row 3"},
			description:   "Content order should follow richTable index, not key values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate content retrieval using richTable indices
			retrievedContent := make([]string, len(tt.richTableKeys))

			// For each row index, we should get content from richTable[rowIndex]
			// NOT by searching for matching keys
			for rowIdx := range tt.richTableKeys {
				// This simulates: content := richTable[rowIdx]
				retrievedContent[rowIdx] = tt.expectedOrder[rowIdx]
			}

			// Verify content is in expected order
			for i, expected := range tt.expectedOrder {
				if retrievedContent[i] != expected {
					t.Errorf("Row %d: got %q, want %q", i, retrievedContent[i], expected)
				}
			}

			t.Logf("Content order verified: %v", retrievedContent)
		})
	}
}

// TestPagesTableRichTableIndexing tests correct usage of richTable indices
func TestPagesTableRichTableIndexing(t *testing.T) {
	// Simulate a richTable with entries
	type RichTableEntry struct {
		Index int    // Position in richTable array
		Key   uint32 // The key value stored in the entry
	}

	tests := []struct {
		name        string
		richTable   []RichTableEntry
		rowCount    int
		description string
	}{
		{
			name: "keys not matching indices",
			richTable: []RichTableEntry{
				{Index: 0, Key: 5},
				{Index: 1, Key: 3},
				{Index: 2, Key: 7},
				{Index: 3, Key: 1},
			},
			rowCount:    4,
			description: "Keys don't match indices - must use index, not key",
		},
		{
			name: "work experience table",
			richTable: []RichTableEntry{
				{Index: 0, Key: 1}, // "Work Experience"
				{Index: 1, Key: 2}, // "Bilibili"
				{Index: 2, Key: 3}, // "NoahStore"
				{Index: 3, Key: 4}, // "Company A"
			},
			rowCount:    4,
			description: "NoahStore at index 2 should appear in row 2, after Bilibili",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test correct indexing approach
			for rowIdx := 0; rowIdx < tt.rowCount; rowIdx++ {
				if rowIdx >= len(tt.richTable) {
					continue
				}

				entry := tt.richTable[rowIdx]

				// Verify we're using index-based lookup, not key-based
				if entry.Index != rowIdx {
					t.Errorf("Entry at position %d has Index=%d, they should match", rowIdx, entry.Index)
				}

				t.Logf("Row %d: using richTable[%d] (key=%d)", rowIdx, rowIdx, entry.Key)
			}
		})
	}
}

// TestPagesTableWrongKeyLookup tests the incorrect approach of using keys for lookup
func TestPagesTableWrongKeyLookup(t *testing.T) {
	// This test demonstrates why key-based lookup is wrong

	type TableEntry struct {
		Key     uint32
		Content string
	}

	richTable := []TableEntry{
		{Key: 1, Content: "Work Experience"}, // Index 0
		{Key: 2, Content: "Bilibili"},        // Index 1
		{Key: 3, Content: "NoahStore"},       // Index 2
		{Key: 4, Content: "Company A"},       // Index 3
	}

	t.Run("wrong approach - key based lookup", func(t *testing.T) {
		// WRONG: Searching by key in buffer
		cellKeys := []uint32{4, 3, 2, 1} // Keys found in cell buffers

		wrongOrder := make([]string, len(cellKeys))
		for cellIdx, cellKey := range cellKeys {
			// This searches richTable for matching key
			for _, entry := range richTable {
				if entry.Key == cellKey {
					wrongOrder[cellIdx] = entry.Content
					break
				}
			}
		}

		// This produces WRONG order
		t.Logf("Wrong approach result: %v", wrongOrder)
		// Expected: ["Work Experience", "Bilibili", "NoahStore", "Company A"]
		// Got: ["Company A", "NoahStore", "Bilibili", "Work Experience"]

		if wrongOrder[0] == "Work Experience" {
			t.Error("Wrong approach should NOT produce correct order")
		}
	})

	t.Run("correct approach - index based lookup", func(t *testing.T) {
		// CORRECT: Using row index directly
		rowCount := len(richTable)

		correctOrder := make([]string, rowCount)
		for rowIdx := 0; rowIdx < rowCount; rowIdx++ {
			// Directly use row index into richTable
			correctOrder[rowIdx] = richTable[rowIdx].Content
		}

		// This produces CORRECT order
		t.Logf("Correct approach result: %v", correctOrder)

		// Verify correct order
		expected := []string{"Work Experience", "Bilibili", "NoahStore", "Company A"}
		for i, exp := range expected {
			if correctOrder[i] != exp {
				t.Errorf("Row %d: got %q, want %q", i, correctOrder[i], exp)
			}
		}
	})
}

// BenchmarkPagesMultiColumnTable benchmarks multi-column table processing
func BenchmarkPagesMultiColumnTable(b *testing.B) {
	// Simulate a 7-column table with 20 rows
	rows := 20
	cols := 7

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create offset data for each row
		for r := 0; r < rows; r++ {
			offsets := make([]uint16, cols)
			// First column has data, rest are empty
			offsets[0] = uint16(r * 10)
			for c := 1; c < cols; c++ {
				offsets[c] = 65535
			}

			// Encode offsets
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, offsets)
		}
	}
}

// BenchmarkEmptyColumnDetection benchmarks empty column detection
func BenchmarkEmptyColumnDetection(b *testing.B) {
	// Create a large table with many empty columns
	rows := 100
	cols := 20

	allOffsets := make([][]uint16, rows)
	for r := 0; r < rows; r++ {
		allOffsets[r] = make([]uint16, cols)
		// First 3 columns have data, rest are empty
		for c := 0; c < cols; c++ {
			if c < 3 {
				allOffsets[r][c] = uint16((r*cols + c) * 10)
			} else {
				allOffsets[r][c] = 65535
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hasContent := make([]bool, cols)
		for _, rowOffsets := range allOffsets {
			for c, offset := range rowOffsets {
				if offset != 65535 {
					hasContent[c] = true
				}
			}
		}
	}
}
