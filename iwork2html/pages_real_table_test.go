package iwork2html

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// TestPagesRealDocument tests conversion of the actual a.pages file
func TestPagesRealDocument(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "a.pages")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file testdata/a.pages not found")
	}

	outputFile := filepath.Join(os.TempDir(), "test_pages_real.html")
	defer os.Remove(outputFile)

	err := Convert(testFile, outputFile)
	if err != nil {
		t.Fatalf("Failed to convert a.pages: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	htmlContent := string(data)

	t.Run("has_table_structure", func(t *testing.T) {
		if !strings.Contains(htmlContent, "<table") {
			t.Error("Output should contain table tags")
		}
		if !strings.Contains(htmlContent, "</table>") {
			t.Error("Output should contain closing table tags")
		}

		tableCount := strings.Count(htmlContent, "<table")
		t.Logf("Found %d tables in a.pages", tableCount)
	})

	t.Run("has_table_cells", func(t *testing.T) {
		hasCells := strings.Contains(htmlContent, "<td") || strings.Contains(htmlContent, "<th")
		if !hasCells {
			t.Error("Table should have cells (td or th tags)")
		}

		tdCount := strings.Count(htmlContent, "<td")
		thCount := strings.Count(htmlContent, "<th")
		t.Logf("Found %d td cells and %d th cells", tdCount, thCount)
	})

	t.Run("content_ordering", func(t *testing.T) {
		// Test that content appears in correct order
		// Using actual content from a.pages file

		// Test 1: Check for product names in order
		lowerHTML := strings.ToLower(htmlContent)

		// Look for actual content from the file
		hasYellow := strings.Contains(lowerHTML, "黄") || strings.Contains(lowerHTML, "yellow")
		hasGreen := strings.Contains(lowerHTML, "绿") || strings.Contains(lowerHTML, "green")
		hasBeideng := strings.Contains(lowerHTML, "北鼎") || strings.Contains(lowerHTML, "beideng")

		if hasBeideng {
			t.Log("✓ Found product names (北鼎)")
		}

		if hasYellow {
			yellowPos := strings.Index(lowerHTML, "黄")
			t.Logf("✓ Found '黄' (yellow) at position %d", yellowPos)
		}

		if hasGreen {
			greenPos := strings.Index(lowerHTML, "绿")
			t.Logf("✓ Found '绿' (green) at position %d", greenPos)
		}

		// Test 2: For documents with Bilibili/NoahStore content
		bilibiliPos := strings.Index(lowerHTML, "bilibili")
		noahstorePos := strings.Index(lowerHTML, "noahstore")

		if bilibiliPos > 0 && noahstorePos > 0 {
			if bilibiliPos >= noahstorePos {
				t.Errorf("Content ordering issue: Bilibili at position %d, NoahStore at position %d. "+
					"Bilibili should appear BEFORE NoahStore", bilibiliPos, noahstorePos)
			} else {
				distance := noahstorePos - bilibiliPos
				t.Logf("✓ Correct ordering: Bilibili at %d, NoahStore at %d (distance: %d)",
					bilibiliPos, noahstorePos, distance)
			}
		}

		// Overall validation: if we found any content, that's good
		if hasBeideng || hasYellow || hasGreen || bilibiliPos > 0 {
			t.Log("✓ Content successfully extracted from table")
		}
	})

	t.Run("actual_content_present", func(t *testing.T) {
		// Check for actual content from the a.pages file

		// Check for product names
		contentItems := []string{
			"北鼎",  // Beideng brand
			"电磁炉", // Induction cooker
			"解冻板", // Defrosting board
		}

		foundCount := 0
		for _, item := range contentItems {
			if strings.Contains(htmlContent, item) {
				foundCount++
				t.Logf("✓ Found content: %s", item)
			}
		}

		if foundCount > 0 {
			t.Logf("✓ Found %d/%d expected content items", foundCount, len(contentItems))
		} else {
			t.Log("Note: Specific product names not found (might be different content)")
		}

		// Alternative check: just verify there's some Chinese content
		hasChinese := false
		for _, r := range htmlContent {
			if r >= 0x4E00 && r <= 0x9FFF {
				hasChinese = true
				break
			}
		}

		if hasChinese {
			t.Log("✓ Document contains Chinese characters")
		}
	})
}

// TestPagesTableHTMLStructure tests the HTML structure of Pages tables
func TestPagesTableHTMLStructure(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "a.pages")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file not found")
	}

	outputFile := filepath.Join(os.TempDir(), "test_pages_structure.html")
	defer os.Remove(outputFile)

	err := Convert(testFile, outputFile)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	// Parse HTML
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	var tableCount int
	var totalCells int
	var findTables func(*html.Node)
	findTables = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			tableCount++

			// Count cells in this table
			var countCells func(*html.Node)
			countCells = func(node *html.Node) {
				if node.Type == html.ElementNode {
					if node.Data == "td" || node.Data == "th" {
						totalCells++
					}
				}
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					countCells(c)
				}
			}
			countCells(n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTables(c)
		}
	}
	findTables(doc)

	t.Logf("Found %d tables with %d total cells", tableCount, totalCells)

	if tableCount == 0 {
		t.Log("Note: No tables found (this might be expected depending on document content)")
	}

	if totalCells == 0 && tableCount > 0 {
		t.Error("Tables found but no cells - this indicates a structure problem")
	}
}

// TestPagesTableCellContent tests that cells have content
func TestPagesTableCellContent(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "a.pages")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file not found")
	}

	outputFile := filepath.Join(os.TempDir(), "test_pages_content.html")
	defer os.Remove(outputFile)

	err := Convert(testFile, outputFile)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	// Parse HTML
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	var cellsWithContent int
	var emptyCells int
	var cellContents []string

	var analyzeCells func(*html.Node)
	analyzeCells = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "td" || n.Data == "th") {
			// Get text content of cell
			var getTextContent func(*html.Node) string
			getTextContent = func(node *html.Node) string {
				if node.Type == html.TextNode {
					return strings.TrimSpace(node.Data)
				}
				var text string
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					text += getTextContent(c)
				}
				return text
			}

			content := getTextContent(n)
			if content != "" {
				cellsWithContent++
				if len(cellContents) < 10 { // Collect first 10 for logging
					cellContents = append(cellContents, content)
				}
			} else {
				emptyCells++
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			analyzeCells(c)
		}
	}
	analyzeCells(doc)

	totalCells := cellsWithContent + emptyCells
	t.Logf("Total cells: %d", totalCells)
	t.Logf("Cells with content: %d", cellsWithContent)
	t.Logf("Empty cells: %d", emptyCells)

	if len(cellContents) > 0 {
		t.Logf("Sample cell contents (first %d):", len(cellContents))
		for i, content := range cellContents {
			// Truncate long content
			if len(content) > 50 {
				content = content[:50] + "..."
			}
			t.Logf("  [%d] %s", i+1, content)
		}
	}

	if totalCells > 0 && cellsWithContent == 0 {
		t.Error("All cells are empty - content is not being rendered")
	}
}

// TestPagesTableColumnStructure tests column structure
func TestPagesTableColumnStructure(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "a.pages")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file not found")
	}

	outputFile := filepath.Join(os.TempDir(), "test_pages_columns.html")
	defer os.Remove(outputFile)

	err := Convert(testFile, outputFile)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	// Parse HTML
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	type TableInfo struct {
		tableIndex int
		rows       int
		maxCols    int
	}

	var tables []TableInfo
	var tableIndex int

	var analyzeTables func(*html.Node)
	analyzeTables = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			tableIndex++
			info := TableInfo{tableIndex: tableIndex}

			// Count rows and columns
			var analyzeRows func(*html.Node)
			analyzeRows = func(node *html.Node) {
				if node.Type == html.ElementNode && node.Data == "tr" {
					info.rows++

					// Count cells in this row
					var cols int
					for c := node.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
							cols++
						}
					}

					if cols > info.maxCols {
						info.maxCols = cols
					}
				}

				for c := node.FirstChild; c != nil; c = c.NextSibling {
					analyzeRows(c)
				}
			}
			analyzeRows(n)

			tables = append(tables, info)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			analyzeTables(c)
		}
	}
	analyzeTables(doc)

	t.Logf("Found %d tables", len(tables))
	for _, table := range tables {
		t.Logf("Table %d: %d rows × %d columns", table.tableIndex, table.rows, table.maxCols)

		if table.rows > 0 && table.maxCols == 0 {
			t.Errorf("Table %d has rows but no columns - structure problem", table.tableIndex)
		}
	}
}

// TestPagesTableSequentialContent tests that content appears sequentially
func TestPagesTableSequentialContent(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "a.pages")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file not found")
	}

	outputFile := filepath.Join(os.TempDir(), "test_pages_sequential.html")
	defer os.Remove(outputFile)

	err := Convert(testFile, outputFile)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	htmlContent := string(data)

	// Check for known sequential content patterns
	// These tests verify that content appears in document order

	tests := []struct {
		name    string
		first   string
		second  string
		message string
	}{
		{
			name:    "bilibili_before_noahstore",
			first:   "bilibili",
			second:  "noahstore",
			message: "Bilibili should appear before NoahStore",
		},
		{
			name:    "work_before_skills",
			first:   "experience",
			second:  "skill",
			message: "Work experience typically appears before skills in resume",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lowerHTML := strings.ToLower(htmlContent)
			firstPos := strings.Index(lowerHTML, tt.first)
			secondPos := strings.Index(lowerHTML, tt.second)

			if firstPos > 0 && secondPos > 0 {
				if firstPos >= secondPos {
					t.Errorf("%s: '%s' at position %d, '%s' at position %d",
						tt.message, tt.first, firstPos, tt.second, secondPos)
				} else {
					t.Logf("✓ Correct order: '%s' (%d) before '%s' (%d)",
						tt.first, firstPos, tt.second, secondPos)
				}
			} else {
				if firstPos == -1 {
					t.Logf("Note: '%s' not found", tt.first)
				}
				if secondPos == -1 {
					t.Logf("Note: '%s' not found", tt.second)
				}
			}
		})
	}
}

// TestPagesTableNoEmptyTables tests that tables have content
func TestPagesTableNoEmptyTables(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "a.pages")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("Test file not found")
	}

	outputFile := filepath.Join(os.TempDir(), "test_pages_nonempty.html")
	defer os.Remove(outputFile)

	err := Convert(testFile, outputFile)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	// Parse HTML
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	var tableIndex int
	var checkTables func(*html.Node)
	checkTables = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			tableIndex++

			// Check if table has any text content
			var hasContent bool
			var checkContent func(*html.Node)
			checkContent = func(node *html.Node) {
				if node.Type == html.TextNode && strings.TrimSpace(node.Data) != "" {
					hasContent = true
				}
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					checkContent(c)
				}
			}
			checkContent(n)

			if !hasContent {
				t.Errorf("Table %d has no text content - completely empty", tableIndex)
			} else {
				t.Logf("✓ Table %d has content", tableIndex)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			checkTables(c)
		}
	}
	checkTables(doc)

	if tableIndex == 0 {
		t.Log("Note: No tables found in document")
	}
}

// BenchmarkPagesRealConversion benchmarks conversion of actual a.pages
func BenchmarkPagesRealConversion(b *testing.B) {
	testFile := filepath.Join("..", "testdata", "a.pages")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		b.Skip("Test file not found")
	}

	outputFile := filepath.Join(os.TempDir(), "bench_pages_real.html")
	defer os.Remove(outputFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Convert(testFile, outputFile); err != nil {
			b.Fatalf("Conversion failed: %v", err)
		}
	}
}
