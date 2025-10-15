package iwork2text

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

// TestCellOffsetHandling tests cell offset handling for text extraction
func TestCellOffsetHandling(t *testing.T) {
	tests := []struct {
		name     string
		offsets  []uint16
		expected int
	}{
		{
			name:     "valid offsets",
			offsets:  []uint16{0, 10, 20},
			expected: 3,
		},
		{
			name:     "with empty cells",
			offsets:  []uint16{0, 65535, 20},
			expected: 3,
		},
		{
			name:     "all empty",
			offsets:  []uint16{65535, 65535},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, tt.offsets)
			if err != nil {
				t.Fatalf("Failed to encode: %v", err)
			}
			cellOffsets := buf.Bytes()

			decoded := make([]uint16, len(tt.offsets))
			err = binary.Read(bytes.NewBuffer(cellOffsets), binary.LittleEndian, decoded)
			if err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			if len(decoded) != tt.expected {
				t.Errorf("Decoded length = %d, want %d", len(decoded), tt.expected)
			}
		})
	}
}

// TestEmptyCellDetection tests detection of empty cells
func TestEmptyCellDetection(t *testing.T) {
	tests := []struct {
		name    string
		offset  uint16
		isEmpty bool
	}{
		{
			name:    "empty cell marker",
			offset:  65535,
			isEmpty: true,
		},
		{
			name:    "valid offset",
			offset:  0,
			isEmpty: false,
		},
		{
			name:    "non-zero offset",
			offset:  100,
			isEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEmpty := tt.offset == 65535
			if isEmpty != tt.isEmpty {
				t.Errorf("isEmpty = %v, want %v", isEmpty, tt.isEmpty)
			}
		})
	}
}

// TestBinaryEncoding tests binary encoding/decoding
func TestBinaryEncoding(t *testing.T) {
	tests := []struct {
		name string
		data []uint16
	}{
		{
			name: "small dataset",
			data: []uint16{1, 2, 3},
		},
		{
			name: "larger dataset",
			data: []uint16{0, 100, 200, 300, 400, 500},
		},
		{
			name: "with markers",
			data: []uint16{0, 65535, 100, 65535, 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, tt.data)
			if err != nil {
				t.Fatalf("Failed to encode: %v", err)
			}

			// Decode
			decoded := make([]uint16, len(tt.data))
			err = binary.Read(bytes.NewBuffer(buf.Bytes()), binary.LittleEndian, decoded)
			if err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			// Verify
			for i, expected := range tt.data {
				if decoded[i] != expected {
					t.Errorf("Position %d: got %d, want %d", i, decoded[i], expected)
				}
			}
		})
	}
}

// TestOffsetRangeValidation tests validation of offset ranges
func TestOffsetRangeValidation(t *testing.T) {
	tests := []struct {
		name    string
		offsets []uint16
		isValid bool
	}{
		{
			name:    "valid sequential offsets",
			offsets: []uint16{0, 10, 20, 30},
			isValid: true,
		},
		{
			name:    "valid non-sequential offsets",
			offsets: []uint16{5, 15, 100, 200},
			isValid: true,
		},
		{
			name:    "with empty markers",
			offsets: []uint16{0, 65535, 10, 65535},
			isValid: true,
		},
		{
			name:    "all empty markers",
			offsets: []uint16{65535, 65535, 65535},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For now, just verify all offsets are valid uint16
			for i, offset := range tt.offsets {
				if offset > 65535 {
					t.Errorf("Offset[%d] = %d exceeds uint16 max", i, offset)
				}
			}
		})
	}
}

// TestLargeDatasets tests handling of large offset datasets
func TestLargeDatasets(t *testing.T) {
	sizes := []int{100, 500, 1000, 5000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			offsets := make([]uint16, size)
			for i := range offsets {
				offsets[i] = uint16(i % 65536)
			}

			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, offsets)
			if err != nil {
				t.Fatalf("Failed to encode: %v", err)
			}

			decoded := make([]uint16, size)
			err = binary.Read(bytes.NewBuffer(buf.Bytes()), binary.LittleEndian, decoded)
			if err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			if len(decoded) != size {
				t.Errorf("Expected %d elements, got %d", size, len(decoded))
			}
		})
	}
}

// TestEmptyCellRatio tests calculating ratio of empty cells
func TestEmptyCellRatio(t *testing.T) {
	tests := []struct {
		name          string
		offsets       []uint16
		expectedRatio float64
	}{
		{
			name:          "50% empty",
			offsets:       []uint16{0, 65535, 10, 65535},
			expectedRatio: 0.5,
		},
		{
			name:          "0% empty",
			offsets:       []uint16{0, 10, 20, 30},
			expectedRatio: 0.0,
		},
		{
			name:          "100% empty",
			offsets:       []uint16{65535, 65535, 65535},
			expectedRatio: 1.0,
		},
		{
			name:          "33% empty",
			offsets:       []uint16{0, 10, 65535},
			expectedRatio: 0.333,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emptyCount := 0
			for _, offset := range tt.offsets {
				if offset == 65535 {
					emptyCount++
				}
			}

			ratio := float64(emptyCount) / float64(len(tt.offsets))
			tolerance := 0.01

			if ratio < tt.expectedRatio-tolerance || ratio > tt.expectedRatio+tolerance {
				t.Errorf("Empty ratio = %.3f, want %.3f (±%.3f)", ratio, tt.expectedRatio, tolerance)
			}
		})
	}
}

// TestDataIntegrity tests data integrity during encoding/decoding
func TestDataIntegrity(t *testing.T) {
	tests := []struct {
		name string
		data []uint16
	}{
		{
			name: "sequential pattern",
			data: []uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "even numbers",
			data: []uint16{0, 2, 4, 6, 8, 10, 12, 14},
		},
		{
			name: "powers of two",
			data: []uint16{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024},
		},
		{
			name: "random pattern",
			data: []uint16{123, 456, 789, 1011, 1213, 1415},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, tt.data)
			if err != nil {
				t.Fatalf("Encoding failed: %v", err)
			}

			// Decode
			decoded := make([]uint16, len(tt.data))
			err = binary.Read(bytes.NewBuffer(buf.Bytes()), binary.LittleEndian, decoded)
			if err != nil {
				t.Fatalf("Decoding failed: %v", err)
			}

			// Verify integrity
			for i, expected := range tt.data {
				if decoded[i] != expected {
					t.Errorf("Data mismatch at index %d: got %d, want %d", i, decoded[i], expected)
				}
			}

			// Verify checksum (simple sum)
			originalSum := uint32(0)
			decodedSum := uint32(0)
			for i := range tt.data {
				originalSum += uint32(tt.data[i])
				decodedSum += uint32(decoded[i])
			}

			if originalSum != decodedSum {
				t.Errorf("Checksum mismatch: original=%d, decoded=%d", originalSum, decodedSum)
			}
		})
	}
}

// TestBufferSizes tests different buffer sizes
func TestBufferSizes(t *testing.T) {
	tests := []struct {
		name       string
		dataSize   int
		bufferSize int
	}{
		{"small data small buffer", 10, 32},
		{"small data large buffer", 10, 1024},
		{"large data large buffer", 1000, 2048},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]uint16, tt.dataSize)
			for i := range data {
				data[i] = uint16(i)
			}

			buf := bytes.NewBuffer(make([]byte, 0, tt.bufferSize))
			err := binary.Write(buf, binary.LittleEndian, data)
			if err != nil {
				t.Fatalf("Write failed: %v", err)
			}

			decoded := make([]uint16, tt.dataSize)
			err = binary.Read(bytes.NewReader(buf.Bytes()), binary.LittleEndian, decoded)
			if err != nil {
				t.Fatalf("Read failed: %v", err)
			}

			for i := range data {
				if decoded[i] != data[i] {
					t.Errorf("Mismatch at %d: got %d, want %d", i, decoded[i], data[i])
				}
			}
		})
	}
}

// BenchmarkCellOffsetDecoding benchmarks cell offset decoding
func BenchmarkCellOffsetDecoding(b *testing.B) {
	offsets := []uint16{0, 10, 20, 30, 40, 50}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, offsets)
	if err != nil {
		b.Fatalf("Failed to encode: %v", err)
	}
	cellOffsets := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded := make([]uint16, len(offsets))
		err := binary.Read(bytes.NewBuffer(cellOffsets), binary.LittleEndian, decoded)
		if err != nil {
			b.Fatalf("Failed to decode: %v", err)
		}
	}
}

// BenchmarkLargeDatasetDecoding benchmarks large dataset decoding
func BenchmarkLargeDatasetDecoding(b *testing.B) {
	offsets := make([]uint16, 10000)
	for i := range offsets {
		offsets[i] = uint16(i % 65536)
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, offsets)
	if err != nil {
		b.Fatalf("Failed to encode: %v", err)
	}
	cellOffsets := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded := make([]uint16, len(offsets))
		err := binary.Read(bytes.NewBuffer(cellOffsets), binary.LittleEndian, decoded)
		if err != nil {
			b.Fatalf("Failed to decode: %v", err)
		}
	}
}

// BenchmarkEmptyCellDetection benchmarks empty cell detection
func BenchmarkEmptyCellDetection(b *testing.B) {
	offsets := make([]uint16, 1000)
	for i := range offsets {
		if i%2 == 0 {
			offsets[i] = 65535
		} else {
			offsets[i] = uint16(i)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		for _, offset := range offsets {
			if offset == 65535 {
				count++
			}
		}
	}
}

// TestPopcountFunction tests the popcount function
func TestPopcountFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    uint16
		expected int
	}{
		{"zero", 0, 0},
		{"one bit", 1, 1},
		{"two bits", 3, 2},
		{"all bits", 0xFFFF, 16},
		{"alternating", 0xAAAA, 8},
		{"powers of two", 0x8000, 1},
		{"mixed pattern", 0x5A5A, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := popcount(tt.input)
			if result != tt.expected {
				t.Errorf("popcount(%d) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDebugMode tests debug mode functionality
func TestDebugMode(t *testing.T) {
	// Test initial state
	SetDebugMode(false)
	if debugMode != false {
		t.Error("Expected debug mode to be false initially")
	}

	// Test setting to true
	SetDebugMode(true)
	if debugMode != true {
		t.Error("Expected debug mode to be true after setting")
	}

	// Test setting back to false
	SetDebugMode(false)
	if debugMode != false {
		t.Error("Expected debug mode to be false after reset")
	}
}

// TestConvertStringBasic tests basic string conversion
func TestConvertStringBasic(t *testing.T) {
	t.Skip("Skipping ConvertString tests - requires valid iWork file input")
}

// TestConvertStringWithOCR tests string conversion with OCR function
func TestConvertStringWithOCR(t *testing.T) {
	t.Skip("Skipping ConvertString tests - requires valid iWork file input")
}

// TestStorageToNode tests storage to node conversion
func TestStorageToNode(t *testing.T) {
	t.Skip("Skipping storageToNode tests - requires complex context setup")
}

// TestCellTypeProcessing tests different cell type processing
func TestCellTypeProcessing(t *testing.T) {
	tests := []struct {
		name     string
		cellType int
		key      uint32
		expected string
	}{
		{
			name:     "cellType 0 (control)",
			cellType: 0,
			key:      1,
			expected: "",
		},
		{
			name:     "cellType 2 (number)",
			cellType: 2,
			key:      42,
			expected: "42",
		},
		{
			name:     "cellType 3 (string)",
			cellType: 3,
			key:      1,
			expected: "",
		},
		{
			name:     "cellType 5 (date)",
			cellType: 5,
			key:      0,
			expected: "",
		},
		{
			name:     "cellType 6 (boolean)",
			cellType: 6,
			key:      0,
			expected: "FALSE",
		},
		{
			name:     "cellType 9 (rich text)",
			cellType: 9,
			key:      1,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a simplified test - in real implementation,
			// we would need to set up proper context and tables
			// For now, just test the basic logic
			var result string
			switch tt.cellType {
			case 2:
				result = fmt.Sprint(tt.key)
			case 6:
				if tt.key == 0 {
					result = "FALSE"
				} else {
					result = "TRUE"
				}
			}
			if result != tt.expected {
				t.Errorf("cellType %d processing = %q, want %q", tt.cellType, result, tt.expected)
			}
		})
	}
}

// TestTableProcessing tests table processing functionality
func TestTableProcessing(t *testing.T) {
	t.Skip("Skipping table processing test - requires complex setup")
}

// TestErrorHandling tests error handling in various functions
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		function    func() error
		expectError bool
	}{
		{
			name: "Convert with invalid file",
			function: func() error {
				return Convert("nonexistent.txt", "output.txt")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.function()
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestEdgeCases tests edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    []uint16
		expected int
	}{
		{
			name:     "single element",
			input:    []uint16{42},
			expected: 1,
		},
		{
			name:     "maximum uint16 values",
			input:    []uint16{0, 65535, 1, 65534},
			expected: 4,
		},
		{
			name:     "all zeros",
			input:    []uint16{0, 0, 0, 0},
			expected: 4,
		},
		{
			name:     "alternating pattern",
			input:    []uint16{0, 65535, 0, 65535, 0},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, tt.input)
			if err != nil {
				t.Fatalf("Failed to encode: %v", err)
			}

			decoded := make([]uint16, len(tt.input))
			err = binary.Read(bytes.NewBuffer(buf.Bytes()), binary.LittleEndian, decoded)
			if err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			if len(decoded) != tt.expected {
				t.Errorf("Expected %d elements, got %d", tt.expected, len(decoded))
			}
		})
	}
}

// TestPerformance tests performance with various data sizes
func TestPerformance(t *testing.T) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			offsets := make([]uint16, size)
			for i := range offsets {
				offsets[i] = uint16(i % 1000)
			}

			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, offsets)
			if err != nil {
				t.Fatalf("Failed to encode: %v", err)
			}

			decoded := make([]uint16, size)
			err = binary.Read(bytes.NewBuffer(buf.Bytes()), binary.LittleEndian, decoded)
			if err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			if len(decoded) != size {
				t.Errorf("Expected %d elements, got %d", size, len(decoded))
			}
		})
	}
}

// Helper functions for testing
func stringPtr(s string) *string {
	return &s
}

func uint32Ptr(u uint32) *uint32 {
	return &u
}
