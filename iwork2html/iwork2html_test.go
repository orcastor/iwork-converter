package iwork2html

import (
	"bytes"
	"encoding/binary"
	"testing"

	"golang.org/x/net/html"
)

// TestElementCreation tests the E() function for creating HTML elements
func TestElementCreation(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want string
	}{
		{
			name: "div element",
			tag:  "div",
			want: "div",
		},
		{
			name: "span element",
			tag:  "span",
			want: "span",
		},
		{
			name: "table element",
			tag:  "table",
			want: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := E(tt.tag)
			if node == nil {
				t.Fatal("E() returned nil")
			}
			if node.Data != tt.want {
				t.Errorf("E(%s) node.Data = %s, want %s", tt.tag, node.Data, tt.want)
			}
			if node.Type != html.ElementNode {
				t.Errorf("E(%s) node.Type = %v, want ElementNode", tt.tag, node.Type)
			}
		})
	}
}

// TestTextNodeCreation tests the T() function for creating text nodes
func TestTextNodeCreation(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "simple text",
			text: "Hello World",
			want: "Hello World",
		},
		{
			name: "empty text",
			text: "",
			want: "",
		},
		{
			name: "text with special chars",
			text: "Hello & <World>",
			want: "Hello & <World>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := T(tt.text)
			if node == nil {
				t.Fatal("T() returned nil")
			}
			if node.Data != tt.want {
				t.Errorf("T(%q) node.Data = %q, want %q", tt.text, node.Data, tt.want)
			}
			if node.Type != html.TextNode {
				t.Errorf("T(%q) node.Type = %v, want TextNode", tt.text, node.Type)
			}
		})
	}
}

// TestDebugMode tests the SetDebugMode function
func TestDebugMode(t *testing.T) {
	tests := []struct {
		name  string
		debug bool
	}{
		{
			name:  "enable debug",
			debug: true,
		},
		{
			name:  "disable debug",
			debug: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			SetDebugMode(tt.debug)
		})
	}
}

// TestCellOffsetParsing tests cell offset encoding and decoding
func TestCellOffsetParsing(t *testing.T) {
	tests := []struct {
		name     string
		offsets  []uint16
		expected []uint16
	}{
		{
			name:     "valid offsets",
			offsets:  []uint16{0, 10, 20, 30},
			expected: []uint16{0, 10, 20, 30},
		},
		{
			name:     "empty cells (65535 marker)",
			offsets:  []uint16{65535, 65535, 0, 10},
			expected: []uint16{65535, 65535, 0, 10},
		},
		{
			name:     "all empty",
			offsets:  []uint16{65535, 65535, 65535},
			expected: []uint16{65535, 65535, 65535},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode offsets to bytes
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, tt.offsets)
			cellOffsets := buf.Bytes()

			// Decode back
			decoded := make([]uint16, len(tt.offsets))
			binary.Read(bytes.NewBuffer(cellOffsets), binary.LittleEndian, decoded)

			for i, expected := range tt.expected {
				if decoded[i] != expected {
					t.Errorf("Offset[%d] = %d, want %d", i, decoded[i], expected)
				}
			}
		})
	}
}

// TestEmptyCellMarker tests the empty cell marker (65535)
func TestEmptyCellMarker(t *testing.T) {
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
			name:    "valid offset zero",
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
				t.Errorf("isEmpty = %v, want %v for offset %d", isEmpty, tt.isEmpty, tt.offset)
			}
		})
	}
}

// TestHTMLNodeStructure tests HTML node parent-child relationships
func TestHTMLNodeStructure(t *testing.T) {
	t.Run("parent with child", func(t *testing.T) {
		parent := E("div")
		child := E("span")
		parent.AppendChild(child)

		if parent.FirstChild != child {
			t.Error("Parent's first child doesn't match appended child")
		}
		if child.Parent != parent {
			t.Error("Child's parent doesn't match parent node")
		}
	})

	t.Run("multiple children", func(t *testing.T) {
		parent := E("div")
		child1 := E("span")
		child2 := E("p")

		parent.AppendChild(child1)
		parent.AppendChild(child2)

		if parent.FirstChild != child1 {
			t.Error("Parent's first child is incorrect")
		}
		if child1.NextSibling != child2 {
			t.Error("Child siblings are not properly linked")
		}
	})
}

// BenchmarkElementCreation benchmarks E() function
func BenchmarkElementCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = E("div")
	}
}

// BenchmarkTextNodeCreation benchmarks T() function
func BenchmarkTextNodeCreation(b *testing.B) {
	text := "Hello World"
	for i := 0; i < b.N; i++ {
		_ = T(text)
	}
}

// TestElementWithAttributes tests creating elements with attributes
func TestElementWithAttributes(t *testing.T) {
	tests := []struct {
		name  string
		tag   string
		attrs []html.Attribute
	}{
		{
			name: "div with class",
			tag:  "div",
			attrs: []html.Attribute{
				{Key: "class", Val: "container"},
			},
		},
		{
			name: "span with style",
			tag:  "span",
			attrs: []html.Attribute{
				{Key: "style", Val: "color: red;"},
			},
		},
		{
			name: "table with multiple attributes",
			tag:  "table",
			attrs: []html.Attribute{
				{Key: "class", Val: "data-table"},
				{Key: "border", Val: "1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := E(tt.tag)
			node.Attr = tt.attrs

			if len(node.Attr) != len(tt.attrs) {
				t.Errorf("Expected %d attributes, got %d", len(tt.attrs), len(node.Attr))
			}

			for i, attr := range tt.attrs {
				if node.Attr[i].Key != attr.Key || node.Attr[i].Val != attr.Val {
					t.Errorf("Attribute[%d] = {%s: %s}, want {%s: %s}",
						i, node.Attr[i].Key, node.Attr[i].Val, attr.Key, attr.Val)
				}
			}
		})
	}
}

// TestNestedElements tests creating nested element structures
func TestNestedElements(t *testing.T) {
	t.Run("three level nesting", func(t *testing.T) {
		grandparent := E("div")
		parent := E("ul")
		child := E("li")

		parent.AppendChild(child)
		grandparent.AppendChild(parent)

		if grandparent.FirstChild != parent {
			t.Error("Grandparent's first child should be parent")
		}
		if parent.FirstChild != child {
			t.Error("Parent's first child should be child")
		}
		if child.Parent != parent {
			t.Error("Child's parent should be parent")
		}
		if parent.Parent != grandparent {
			t.Error("Parent's parent should be grandparent")
		}
	})

	t.Run("multiple children at same level", func(t *testing.T) {
		parent := E("ul")
		child1 := E("li")
		child2 := E("li")
		child3 := E("li")

		parent.AppendChild(child1)
		parent.AppendChild(child2)
		parent.AppendChild(child3)

		// Verify chain
		if parent.FirstChild != child1 {
			t.Error("First child is incorrect")
		}
		if child1.NextSibling != child2 {
			t.Error("Child1's next sibling should be child2")
		}
		if child2.NextSibling != child3 {
			t.Error("Child2's next sibling should be child3")
		}
		if child3.NextSibling != nil {
			t.Error("Child3 should have no next sibling")
		}
	})
}

// TestTextNodeWithEmptyString tests text node with empty string
func TestTextNodeWithEmptyString(t *testing.T) {
	node := T("")
	if node == nil {
		t.Fatal("T(\"\") returned nil")
	}
	if node.Data != "" {
		t.Errorf("Expected empty string, got %q", node.Data)
	}
	if node.Type != html.TextNode {
		t.Error("Node should be TextNode type")
	}
}

// TestTextNodeWithUnicode tests text node with Unicode characters
func TestTextNodeWithUnicode(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"Chinese characters", "你好世界"},
		{"Japanese characters", "こんにちは"},
		{"Korean characters", "안녕하세요"},
		{"Emoji", "😀😁😂"},
		{"Mixed", "Hello 世界 🌍"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := T(tt.text)
			if node.Data != tt.text {
				t.Errorf("Expected %q, got %q", tt.text, node.Data)
			}
		})
	}
}

// TestCellOffsetBoundaryValues tests boundary values for cell offsets
func TestCellOffsetBoundaryValues(t *testing.T) {
	tests := []struct {
		name    string
		offsets []uint16
	}{
		{
			name:    "zero values",
			offsets: []uint16{0, 0, 0, 0},
		},
		{
			name:    "max uint16",
			offsets: []uint16{65535, 65535, 65535},
		},
		{
			name:    "mixed boundary values",
			offsets: []uint16{0, 65535, 32768, 1, 65534},
		},
		{
			name:    "single offset",
			offsets: []uint16{100},
		},
		{
			name:    "large array",
			offsets: make([]uint16, 1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, tt.offsets)
			if err != nil {
				t.Fatalf("Failed to encode: %v", err)
			}

			decoded := make([]uint16, len(tt.offsets))
			err = binary.Read(bytes.NewBuffer(buf.Bytes()), binary.LittleEndian, decoded)
			if err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			for i, expected := range tt.offsets {
				if decoded[i] != expected {
					t.Errorf("Offset[%d] = %d, want %d", i, decoded[i], expected)
				}
			}
		})
	}
}

// TestEmptyCellPattern tests common empty cell patterns
func TestEmptyCellPattern(t *testing.T) {
	tests := []struct {
		name            string
		offsets         []uint16
		expectedEmpties int
	}{
		{
			name:            "all empty",
			offsets:         []uint16{65535, 65535, 65535, 65535},
			expectedEmpties: 4,
		},
		{
			name:            "alternating empty",
			offsets:         []uint16{0, 65535, 10, 65535, 20, 65535},
			expectedEmpties: 3,
		},
		{
			name:            "leading empties",
			offsets:         []uint16{65535, 65535, 0, 10, 20},
			expectedEmpties: 2,
		},
		{
			name:            "trailing empties",
			offsets:         []uint16{0, 10, 20, 65535, 65535},
			expectedEmpties: 2,
		},
		{
			name:            "no empties",
			offsets:         []uint16{0, 10, 20, 30, 40},
			expectedEmpties: 0,
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

			if emptyCount != tt.expectedEmpties {
				t.Errorf("Expected %d empty cells, found %d", tt.expectedEmpties, emptyCount)
			}
		})
	}
}

// TestDebugModeToggling tests debug mode can be toggled multiple times
func TestDebugModeToggling(t *testing.T) {
	modes := []bool{true, false, true, false, true}

	for i, mode := range modes {
		t.Run("toggle_"+string(rune('0'+i)), func(t *testing.T) {
			// Should not panic
			SetDebugMode(mode)
		})
	}
}

// TestElementTypes tests creation of various HTML element types
func TestElementTypes(t *testing.T) {
	elementTypes := []string{
		"div", "span", "p", "a", "table", "tr", "td", "th",
		"ul", "ol", "li", "h1", "h2", "h3", "h4", "h5", "h6",
		"header", "footer", "section", "article", "nav",
	}

	for _, tag := range elementTypes {
		t.Run("create_"+tag, func(t *testing.T) {
			node := E(tag)
			if node == nil {
				t.Fatal("E() returned nil")
			}
			if node.Data != tag {
				t.Errorf("Expected tag %s, got %s", tag, node.Data)
			}
			if node.Type != html.ElementNode {
				t.Error("Expected ElementNode type")
			}
		})
	}
}

// BenchmarkCellOffsetDecoding benchmarks cell offset decoding
func BenchmarkCellOffsetDecoding(b *testing.B) {
	offsets := []uint16{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, offsets)
	cellOffsets := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded := make([]uint16, len(offsets))
		binary.Read(bytes.NewBuffer(cellOffsets), binary.LittleEndian, decoded)
	}
}

// BenchmarkLargeOffsetArray benchmarks decoding large offset arrays
func BenchmarkLargeOffsetArray(b *testing.B) {
	offsets := make([]uint16, 1000)
	for i := range offsets {
		offsets[i] = uint16(i % 65536)
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, offsets)
	cellOffsets := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded := make([]uint16, len(offsets))
		binary.Read(bytes.NewBuffer(cellOffsets), binary.LittleEndian, decoded)
	}
}

// BenchmarkNestedElementCreation benchmarks creating nested elements
func BenchmarkNestedElementCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parent := E("div")
		for j := 0; j < 10; j++ {
			child := E("span")
			parent.AppendChild(child)
		}
	}
}
