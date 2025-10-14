package iwork2html

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/orcastor/iwork-converter/index"
	"github.com/orcastor/iwork-converter/proto/KN"
	"github.com/orcastor/iwork-converter/proto/TN"
	"github.com/orcastor/iwork-converter/proto/TP"
	"github.com/orcastor/iwork-converter/proto/TSCE"
	"github.com/orcastor/iwork-converter/proto/TSD"
	"github.com/orcastor/iwork-converter/proto/TSP"
	"github.com/orcastor/iwork-converter/proto/TST"
	"github.com/orcastor/iwork-converter/proto/TSWP"

	"golang.org/x/net/html"
)

// Global debug mode flag
var debugMode bool = true

// SetDebugMode sets the global debug mode
func SetDebugMode(debug bool) {
	debugMode = debug
}

// T is a helper function for building html text nodes.
func T(value string) *html.Node {
	return &html.Node{Type: html.TextNode, Data: cleanText(value)}
}

// cleanText removes object replacement characters and other unwanted Unicode characters
func cleanText(text string) string {
	// Remove object replacement character (U+FFFC)
	text = strings.ReplaceAll(text, "\uFFFC", "")
	// Remove other common unwanted characters that might appear
	text = strings.ReplaceAll(text, "\uFEFF", "") // BOM (Byte Order Mark)
	return text
}

// E is a helper function for building HTML nodes.
func E(tag string, children ...interface{}) *html.Node {
	node := &html.Node{}
	node.Data = tag
	node.Type = html.ElementNode

	for _, child := range children {
		switch child.(type) {
		case string:
			node.AppendChild(T(child.(string)))
		case *html.Node:
			node.AppendChild(child.(*html.Node))
		case []string:
			args := child.([]string)
			for i := 0; i < len(args)-1; i += 2 {
				node.Attr = append(node.Attr, html.Attribute{Key: args[i], Val: args[i+1]})
			}
		default:
			log.Fatalf("Unhandled type %T in E", child)
		}
	}
	return node
}

type Context struct {
	styles    map[string]string
	imgs      map[string]uint64
	ix        *index.Index
	zr        *zip.ReadCloser
	fontScale float64
	// Slide dimensions for font scaling
	slideWidth  float64
	slideHeight float64
}

// Control whether to output table cell debug logs
var debugTableCells = false

// Control whether to output image processing debug logs
var debugImages = true

type Attachment struct {
	pos  uint32
	node *html.Node
}

func (ctx *Context) processImage(image *TSD.ImageArchive) *html.Node {
	if debugMode && debugImages {
		fmt.Printf("DEBUG IMG: processing image, dataId=%d\n", *image.Data.Identifier)
	}
	dataId := *image.Data.Identifier
	meta := ctx.ix.Records[2].(*TSP.PackageMetadata)
	var src string
	for _, data := range meta.Datas {
		if dataId == *data.Identifier {
			if data.FileName != nil {
				src = *data.FileName
				if debugMode && debugImages {
					fmt.Printf("DEBUG IMG: use FileName=%s for dataId=%d\n", src, dataId)
				}
			} else {
				if debugMode {
					fmt.Printf("No filename: %#v\n", data)
				}
				src = *data.PreferredFileName
				if debugMode && debugImages {
					fmt.Printf("DEBUG IMG: use PreferredFileName=%s for dataId=%d\n", src, dataId)
				}
			}
		}
	}
	ctx.imgs["Data/"+src] = dataId
	// not sure if this is px or pt.  It's px on the html side.
	width := fmt.Sprintf("%f", *image.OriginalSize.Width)
	height := fmt.Sprintf("%f", *image.OriginalSize.Height)
	if debugMode && debugImages {
		fmt.Printf("DEBUG IMG: original size width=%s height=%s for src=%s\n", width, height, src)
	}
	return E("img", []string{"src", "", "width", width, "height", height, "class", "img_" + fmt.Sprint(dataId)})
}

var LE = binary.LittleEndian

func popcount(v uint16) int {
	var c int
	for ; v != 0; c++ {
		v &= v - 1
	}
	return c
}

// applyCellStyle applies cell styles
func (ctx *Context) applyCellStyle(tm *TST.TableModelArchive, key uint32) string {
	style := ""

	if tm.BaseDataStore != nil && tm.BaseDataStore.StyleTable != nil {
		styleTableRef := ctx.ix.Deref(tm.BaseDataStore.StyleTable)
		if tdl, ok := styleTableRef.(*TST.TableDataList); ok {
			if debugTableCells {
				fmt.Printf("DEBUG: Looking for style with key %d in StyleTable with %d entries\n", key, len(tdl.Entries))
				// Print all StyleTable keys for debugging
				fmt.Printf("DEBUG: StyleTable keys: ")
				for i, entry := range tdl.Entries {
					if entry.Key != nil {
						fmt.Printf("%d", *entry.Key)
						if i < len(tdl.Entries)-1 {
							fmt.Printf(", ")
						}
					}
				}
				fmt.Printf("\n")
			}
			for _, entry := range tdl.Entries {
				if debugTableCells && entry.Key != nil {
					fmt.Printf("DEBUG: Checking StyleTable entry key %d against target key %d\n", *entry.Key, key)
				}
				if *entry.Key == key && entry.Reference != nil {
					entryRef := ctx.ix.Deref(entry.Reference)
					if debugTableCells {
						fmt.Printf("DEBUG: Found style entry for key %d, type: %T\n", key, entryRef)
					}
					if csa, ok := entryRef.(*TST.CellStyleArchive); ok {
						if csa.CellProperties != nil {
							// Handle background fill
							if csa.CellProperties.CellFill != nil {
								if css := colorToCSS(csa.CellProperties.CellFill.GetColor()); css != "" {
									applyBackgroundColor(&style, css, true)
								}
							}
							// Handle borders and rounded corners
							processCellBorders(&style, csa.CellProperties)
							// Handle font size
							processCellFont(&style, csa.CellProperties)
						}
					} else if psa, ok := entryRef.(*TSWP.ParagraphStyleArchive); ok {
						// Handle paragraph styles as cell styles
						if debugMode && debugTableCells {
							fmt.Printf("DEBUG: Processing ParagraphStyleArchive for key %d\n", key)
							fmt.Printf("DEBUG: ParaProperties: %v\n", psa.ParaProperties != nil)
							if psa.ParaProperties != nil {
								fmt.Printf("DEBUG: Fill: %v, Stroke: %v\n", psa.ParaProperties.Fill != nil, psa.ParaProperties.Stroke != nil)
								// Check if there are other properties that might contain background color
								if psa.ParaProperties.Fill != nil {
									fmt.Printf("DEBUG: Fill details: %+v\n", psa.ParaProperties.Fill)
								}
							}
						}
						if psa.ParaProperties != nil {
							// Handle paragraph background fill
							if psa.ParaProperties.Fill != nil {
								if css := colorToCSS(psa.ParaProperties.Fill); css != "" {
									applyBackgroundColor(&style, css, false)
									if debugTableCells {
										fmt.Printf("DEBUG: Applied paragraph background color: %s\n", css)
									}
								} else {
									if debugTableCells {
										fmt.Printf("DEBUG: Paragraph background color conversion failed\n")
									}
								}
							} else {
								if debugTableCells {
									fmt.Printf("DEBUG: No paragraph background fill found\n")
								}
							}
							// Handle paragraph borders
							if psa.ParaProperties.Stroke != nil {
								col := colorToCSS(psa.ParaProperties.Stroke.Color)
								w := 1.0
								if psa.ParaProperties.Stroke.Width != nil {
									w = float64(*psa.ParaProperties.Stroke.Width)
								}
								if col != "" {
									style += fmt.Sprintf("border: %.2fpx solid %s; box-sizing: border-box;", w, col)
								}
							}
						}
					} else {
						// Print unrecognized style types
						fmt.Printf("DEBUG: Unrecognized style type: %T\n", entryRef)
					}
					break
				}
			}
		} else {
			// Print unrecognized StyleTable types
			fmt.Printf("DEBUG: Unrecognized StyleTable type: %T\n", styleTableRef)
		}
	}

	return style
}

// applyNumbersCellStyle applies Numbers-specific cell styles from dedicated style references
func (ctx *Context) applyNumbersCellStyle(styleRef *TSP.Reference) string {
	style := ""

	if styleRef == nil {
		return style
	}

	styleObj := ctx.ix.Deref(styleRef)
	if styleObj == nil {
		if debugTableCells {
			fmt.Printf("DEBUG: Failed to deref Numbers style reference\n")
		}
		return style
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Numbers style object type: %T\n", styleObj)
	}

	// Handle CellStyleArchive
	if csa, ok := styleObj.(*TST.CellStyleArchive); ok {
		if csa.CellProperties != nil {
			// Handle background fill
			if csa.CellProperties.CellFill != nil {
				if css := colorToCSS(csa.CellProperties.CellFill.GetColor()); css != "" {
					applyBackgroundColor(&style, css, true)
					if debugTableCells {
						fmt.Printf("DEBUG: Applied Numbers cell background: %s\n", css)
					}
				}
			}
			// Handle borders and rounded corners
			processCellBorders(&style, csa.CellProperties)
			// Handle font size
			processCellFont(&style, csa.CellProperties)
		}
	} else if psa, ok := styleObj.(*TSWP.ParagraphStyleArchive); ok {
		// Handle paragraph styles as cell styles
		if psa.ParaProperties != nil {
			// Handle paragraph background fill
			if psa.ParaProperties.Fill != nil {
				if css := colorToCSS(psa.ParaProperties.Fill); css != "" {
					applyBackgroundColor(&style, css, false)
					if debugTableCells {
						fmt.Printf("DEBUG: Applied Numbers paragraph background: %s\n", css)
					}
				}
			}
			// Handle paragraph borders
			if psa.ParaProperties.Stroke != nil {
				col := colorToCSS(psa.ParaProperties.Stroke.Color)
				w := 1.0
				if psa.ParaProperties.Stroke.Width != nil {
					w = float64(*psa.ParaProperties.Stroke.Width)
				}
				if col != "" {
					style += fmt.Sprintf("border: %.2fpx solid %s; box-sizing: border-box;", w, col)
					if debugTableCells {
						fmt.Printf("DEBUG: Applied Numbers paragraph border: %.2fpx solid %s\n", w, col)
					}
				}
			}
		}
	} else {
		if debugTableCells {
			fmt.Printf("DEBUG: Unrecognized Numbers style type: %T\n", styleObj)
		}
	}

	return style
}

// applyBackgroundColor unified background color application
func applyBackgroundColor(style *string, css string, important bool) {
	if css == "" {
		return
	}

	// Always apply background colors - let the user decide what they want
	// The previous filtering was too aggressive and removed legitimate header colors
	if important {
		*style += fmt.Sprintf("background-color: %s !important;", css)
	} else {
		*style += fmt.Sprintf("background-color: %s;", css)
	}
}

// mergeParentStyles recursively processes parent style inheritance
func (ctx *Context) mergeParentStyles(child, parent *TSWP.ParagraphStyleArchive) {
	if parent.Super.Parent != nil {
		grandParent := ctx.ix.Deref(parent.Super.Parent).(*TSWP.ParagraphStyleArchive)
		// First process grandparent styles to parent styles
		mergeCharProps(parent.CharProperties, grandParent.CharProperties)
		mergeParaProps(parent.ParaProperties, grandParent.ParaProperties)
		// Recursively process deeper inheritance
		ctx.mergeParentStyles(parent, grandParent)
		// Then apply the processed parent styles to child styles
		mergeCharProps(child.CharProperties, parent.CharProperties)
		mergeParaProps(child.ParaProperties, parent.ParaProperties)
	}
}

// mergeParentCharStyles recursively processes parent character style inheritance
func (ctx *Context) mergeParentCharStyles(child, parent *TSWP.CharacterStyleArchive) {
	if parent.Super.Parent != nil {
		grandParentRef := ctx.ix.Deref(parent.Super.Parent)
		if grandParent, ok := grandParentRef.(*TSWP.CharacterStyleArchive); ok {
			// First process grandparent styles to parent styles
			mergeCharProps(parent.CharProperties, grandParent.CharProperties)
			// Recursively process deeper inheritance
			ctx.mergeParentCharStyles(parent, grandParent)
			// Then apply the processed parent styles to child styles
			mergeCharProps(child.CharProperties, parent.CharProperties)
		} else {
			// Print unrecognized grandparent character style type
			fmt.Printf("DEBUG: Unrecognized grandparent character style type: %T\n", grandParentRef)
		}
	}
}

// applyPositionBasedStyle applies styles based on cell position - now uses cell's own style from StyleTable
func (ctx *Context) applyPositionBasedStyle(tm *TST.TableModelArchive, globalRow, c int, shouldTreatFirstRowAsHeader bool) string {
	// Return empty string to rely on CSS styles
	// Individual cell styles will be applied from StyleTable in the cell processing loop
	return ""
}

func (ctx *Context) processTable(tm *TST.TableModelArchive) *html.Node {
	rows := uint32(0)
	cols := uint32(0)
	if tm.NumberOfRows != nil {
		rows = *tm.NumberOfRows
	}
	if tm.NumberOfColumns != nil {
		cols = *tm.NumberOfColumns
	}
	if debugTableCells {
		fmt.Printf("DEBUG: Processing table with %d rows, %d columns (pointers: %p, %p)\n", rows, cols, tm.NumberOfRows, tm.NumberOfColumns)
		if tm.NumberOfHeaderRows != nil {
			fmt.Printf("DEBUG: NumberOfHeaderRows = %d\n", *tm.NumberOfHeaderRows)
		} else {
			fmt.Printf("DEBUG: NumberOfHeaderRows is nil\n")
		}

		// Check data storage structure
		if tm.BaseDataStore != nil && tm.BaseDataStore.Tiles != nil {
			fmt.Printf("DEBUG: DataStore has %d tiles\n", len(tm.BaseDataStore.Tiles.Tiles))
			for i, tinfo := range tm.BaseDataStore.Tiles.Tiles {
				tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
				fmt.Printf("DEBUG: Tile %d has %d rows\n", i, len(tile.RowInfos))
			}
		}
	}
	// Extract string and rich text tables
	var stringTable []*TST.TableDataList_ListEntry
	var richTable []*TST.TableDataList_ListEntry

	// Extract merge region map
	var mergeRegions []*TST.CellRange
	if debugTableCells {
		if tm.BaseDataStore != nil {
			if tm.BaseDataStore.MergeRegionMap != nil {
				fmt.Printf("DEBUG: MergeRegionMap reference found: ID=%d\n", *tm.BaseDataStore.MergeRegionMap.Identifier)
				// Try to deref and see what type it is
				derefed := ctx.ix.Deref(tm.BaseDataStore.MergeRegionMap)
				fmt.Printf("DEBUG: MergeRegionMap dereferenced type: %T\n", derefed)
			} else {
				fmt.Printf("DEBUG: No MergeRegionMap reference\n")
			}
		}
	}
	if tm.BaseDataStore != nil && tm.BaseDataStore.MergeRegionMap != nil {
		if mrm, ok := ctx.ix.Deref(tm.BaseDataStore.MergeRegionMap).(*TST.MergeRegionMapArchive); ok {
			mergeRegions = mrm.CellRange
			if debugTableCells {
				fmt.Printf("DEBUG: MergeRegionMap loaded with %d regions\n", len(mergeRegions))
				for i, region := range mergeRegions {
					if region.Origin != nil && region.Size != nil {
						// Parse CellID from PackedData (row<<16 | col)
						var fromRow, fromCol uint32
						if region.Origin.PackedData != nil {
							packed := *region.Origin.PackedData
							fromRow = packed >> 16
							fromCol = packed & 0xFFFF
						}
						toRow := fromRow + *region.Size.NumRows - 1
						toCol := fromCol + *region.Size.NumColumns - 1
						fmt.Printf("  MergeRegion[%d]: Row %d-%d, Col %d-%d (size: %d×%d)\n",
							i, fromRow, toRow, fromCol, toCol,
							*region.Size.NumRows, *region.Size.NumColumns)
					}
				}
			}
		}
	}

	if tm.BaseDataStore != nil {
		if debugTableCells {
			fmt.Printf("DEBUG: DataStore found\n")
		}
		if tm.BaseDataStore.StringTable != nil {
			if debugTableCells {
				fmt.Printf("DEBUG: StringTable reference found\n")
			}
			if tdl, ok := ctx.ix.Deref(tm.BaseDataStore.StringTable).(*TST.TableDataList); ok {
				stringTable = tdl.Entries
				if debugTableCells {
					fmt.Printf("DEBUG: StringTable loaded with %d entries\n", len(stringTable))
					for i, e := range stringTable {
						if e != nil && e.Key != nil && e.String_ != nil {
							preview := *e.String_
							if len(preview) > 30 {
								preview = preview[:30] + "..."
							}
							fmt.Printf("  StringTable[%d]: key=%d -> %q\n", i, *e.Key, preview)
						}
					}
				}
			} else {
				if debugTableCells {
					fmt.Printf("DEBUG: Failed to deref StringTable\n")
				}
			}
		} else {
			if debugTableCells {
				fmt.Printf("DEBUG: No StringTable reference\n")
			}
		}
		if tm.BaseDataStore.RichTextTable != nil {
			if debugTableCells {
				fmt.Printf("DEBUG: RichTextTable reference found\n")
			}
			if tdl, ok := ctx.ix.Deref(tm.BaseDataStore.RichTextTable).(*TST.TableDataList); ok {
				richTable = tdl.Entries
				if debugTableCells {
					fmt.Printf("DEBUG: RichTextTable loaded with %d entries\n", len(richTable))
					for i, entry := range richTable {
						fmt.Printf("DEBUG: RichTextTable[%d]: key=%d\n", i, *entry.Key)
						if entry.RichTextPayload != nil {
							if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
								if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
									if len(st.Text) > 0 {
										preview := st.Text[0]
										if len(preview) > 200 {
											preview = preview[:200] + "..."
										}
										fmt.Printf("DEBUG: RichTextTable[%d] content: %s\n", i, preview)
									} else {
										fmt.Printf("DEBUG: RichTextTable[%d] has no text\n", i)
									}
								} else {
									fmt.Printf("DEBUG: RichTextTable[%d] failed to deref Storage\n", i)
								}
							} else {
								fmt.Printf("DEBUG: RichTextTable[%d] failed to deref RichTextPayload\n", i)
							}
						} else {
							fmt.Printf("DEBUG: RichTextTable[%d] has no RichTextPayload\n", i)
						}
					}
				}
			} else {
				if debugTableCells {
					fmt.Printf("DEBUG: Failed to deref RichTextTable\n")
				}
			}
		} else {
			if debugTableCells {
				fmt.Printf("DEBUG: No RichTextTable reference\n")
			}
		}

		// Load format table for date cells (cellType=5)
		if tm.BaseDataStore.FormatTable != nil {
			if debugTableCells {
				fmt.Printf("DEBUG: FormatTable reference found\n")
			}
			refObj := ctx.ix.Deref(tm.BaseDataStore.FormatTable)
			if debugTableCells {
				fmt.Printf("DEBUG: FormatTable deref result type: %T\n", refObj)
			}
			if tdl, ok := refObj.(*TST.TableDataList); ok {
				if debugTableCells {
					fmt.Printf("DEBUG: FormatTable loaded with %d entries\n", len(tdl.Entries))
					for i, entry := range tdl.Entries {
						if entry != nil && entry.Key != nil {
							fmt.Printf("DEBUG: FormatTable[%d]: key=%d\n", i, *entry.Key)
						}
					}
				}
			} else {
				if debugTableCells {
					fmt.Printf("DEBUG: Failed to deref FormatTable\n")
				}
			}
		} else {
			if debugTableCells {
				fmt.Printf("DEBUG: No FormatTable reference\n")
			}
		}

		// Load PopUpMenuModel for control cells (cellType=0)
		if tm.BaseDataStore.MultipleChoiceListFormatTable != nil {
			if debugTableCells {
				fmt.Printf("DEBUG: MultipleChoiceListFormatTable reference found\n")
			}
			refObj := ctx.ix.Deref(tm.BaseDataStore.MultipleChoiceListFormatTable)
			if debugTableCells {
				fmt.Printf("DEBUG: MultipleChoiceListFormatTable deref result type: %T\n", refObj)
			}
			if tdl, ok := refObj.(*TST.TableDataList); ok {
				if debugTableCells {
					fmt.Printf("DEBUG: TableDataList loaded with %d entries\n", len(tdl.Entries))
				}
				// Look for PopUpMenuModel in the entries
				for i, entry := range tdl.Entries {
					if entry != nil && entry.Key != nil {
						if popup, ok := ctx.ix.Deref(entry.Reference).(*TST.PopUpMenuModel); ok {
							if debugTableCells {
								fmt.Printf("DEBUG: Found PopUpMenuModel at entry %d with %d items\n", i, len(popup.Item))
								for j, item := range popup.Item {
									if item != nil && item.CellValueType != nil {
										fmt.Printf("DEBUG: PopUpMenuModel[%d][%d]: type=%d\n", i, j, *item.CellValueType)
										// Try to extract string value
										if item.StringValue != nil && item.StringValue.Value != nil {
											fmt.Printf("DEBUG: PopUpMenuModel[%d][%d] string: %s\n", i, j, *item.StringValue.Value)
										}
									}
								}
							}
						}
					}
				}
			} else {
				if debugTableCells {
					fmt.Printf("DEBUG: Failed to deref TableDataList\n")
				}
			}
		} else {
			if debugTableCells {
				fmt.Printf("DEBUG: No MultipleChoiceListFormatTable reference\n")
			}
		}
	} else {
		if debugTableCells {
			fmt.Printf("DEBUG: No DataStore found\n")
		}
	}

	cc := int(*tm.NumberOfColumns)

	// Simplified: always show all columns, no active column scanning
	activeColumns := make([]bool, cc)
	for i := 0; i < cc; i++ {
		activeColumns[i] = true
	}

	// Calculate the number of columns with content
	activeColumnCount := 0
	for _, active := range activeColumns {
		if active {
			activeColumnCount++
		}
	}

	// Smart table structure detection: if only a few columns have content, it might be a table structure error
	if debugTableCells {
		fmt.Printf("DEBUG: Active columns: %d out of %d total columns\n", activeColumnCount, cc)
	}

	// Note: Smart restructuring is disabled - we now use proper column detection instead

	table := E("table")

	// Don't display table titles, remove all table title displays

	// Generate column definitions - intelligent width allocation
	colgroup := E("colgroup")

	// Check which columns actually have content by analyzing the data
	hasContentColumns := make([]bool, cc)

	// Only detect empty columns for Pages documents
	// For Numbers, show all columns regardless of content
	if ctx.ix.Type == "pages" {
		// First, try to get column widths from columnHeaders
		columnWidths := make([]float32, cc)
		hasWidthInfo := false

		if tm.BaseDataStore != nil && tm.BaseDataStore.ColumnHeaders != nil {
			colHeadersRef := ctx.ix.Deref(tm.BaseDataStore.ColumnHeaders)
			if debugTableCells {
				fmt.Printf("DEBUG: ColumnHeaders type: %T\n", colHeadersRef)
			}

			// ColumnHeaders can be either HeaderStorage or HeaderStorageBucket
			if colHeaders, ok := colHeadersRef.(*TST.HeaderStorage); ok {
				// HeaderStorage: contains multiple buckets
				for bucketIdx, bucketRef := range colHeaders.Buckets {
					if bucket, ok := ctx.ix.Deref(bucketRef).(*TST.HeaderStorageBucket); ok {
						if debugTableCells {
							fmt.Printf("DEBUG: Bucket %d has %d headers\n", bucketIdx, len(bucket.Headers))
						}
						for _, header := range bucket.Headers {
							if header.Index != nil && header.Size != nil {
								idx := *header.Index
								if int(idx) < len(columnWidths) {
									columnWidths[idx] = *header.Size
									hasWidthInfo = true
									if debugTableCells {
										fmt.Printf("DEBUG: Column %d width=%.2f\n", idx, *header.Size)
									}
								}
							}
						}
					}
				}
			} else if bucket, ok := colHeadersRef.(*TST.HeaderStorageBucket); ok {
				// HeaderStorageBucket: single bucket with headers
				if debugTableCells {
					fmt.Printf("DEBUG: Single bucket has %d headers\n", len(bucket.Headers))
				}
				for _, header := range bucket.Headers {
					if header.Index != nil && header.Size != nil {
						idx := *header.Index
						if int(idx) < len(columnWidths) {
							columnWidths[idx] = *header.Size
							hasWidthInfo = true
							if debugTableCells {
								fmt.Printf("DEBUG: Column %d width=%.2f\n", idx, *header.Size)
							}
						}
					}
				}
			}
		} else if debugTableCells {
			fmt.Printf("DEBUG: No ColumnHeaders available\n")
		}

		// Filter columns based on width and merge patterns
		if hasWidthInfo {
			// Check if columns have similar widths (indicating potential merged cells)
			// and filter based on content patterns
			columnUniqueKeys := make(map[int]map[uint32]bool)
			for i := 0; i < int(cc); i++ {
				columnUniqueKeys[i] = make(map[uint32]bool)
			}

			if tm.BaseDataStore != nil && tm.BaseDataStore.Tiles != nil {
				for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
					tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
					for _, rinfo := range tile.RowInfos {
						var offsets []uint16
						if len(rinfo.CellOffsets) > 0 {
							offsets = make([]uint16, len(rinfo.CellOffsets)/2)
							binary.Read(bytes.NewBuffer(rinfo.CellOffsets), LE, offsets)
						}

						cellStorageBuf := rinfo.CellStorageBuffer
						if len(cellStorageBuf) == 0 && len(rinfo.CellStorageBufferPreBnc) > 0 {
							cellStorageBuf = rinfo.CellStorageBufferPreBnc
						}

						for c := 0; c < int(cc) && c < len(offsets); c++ {
							offset := offsets[c]
							if offset != 65535 && len(cellStorageBuf) >= int(offset)+16 {
								key := LE.Uint32(cellStorageBuf[offset+12 : offset+16])
								if key != 0 {
									columnUniqueKeys[c][key] = true
								}
							}
						}
					}
				}
			}

			// Column detection strategy:
			// 1. Keep columns with width >= 1.0 AND diverse content
			// 2. Skip columns that have only merged cell keys (low diversity)
			for c := 0; c < int(cc); c++ {
				uniqueKeyCount := len(columnUniqueKeys[c])

				// Keep column if:
				// - Width >= 1.0
				// - AND (first column OR has highly diverse content: >70% unique keys)
				if columnWidths[c] >= 1.0 {
					if c == 0 {
						// Always keep first column
						hasContentColumns[c] = true
					} else if uniqueKeyCount > int(rows)*7/10 && uniqueKeyCount > 1 {
						// Keep columns with significant content diversity (>70% of rows have unique keys)
						// This filters out columns with mostly merged cells
						hasContentColumns[c] = true
					}
				}

				if debugTableCells {
					fmt.Printf("DEBUG: Column %d width=%.2f, unique_keys=%d, keep=%v\n", c, columnWidths[c], uniqueKeyCount, hasContentColumns[c])
				}
			}

			// Removed: automatic second column addition
			// User feedback indicates only Column 0 should be shown
		} else {
			// Fallback: use content-based detection
			columnUniqueKeys := make(map[int]map[uint32]bool)
			for i := 0; i < int(cc); i++ {
				columnUniqueKeys[i] = make(map[uint32]bool)
			}

			if tm.BaseDataStore != nil && tm.BaseDataStore.Tiles != nil {
				for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
					tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
					for _, rinfo := range tile.RowInfos {
						var offsets []uint16
						if len(rinfo.CellOffsets) > 0 {
							offsets = make([]uint16, len(rinfo.CellOffsets)/2)
							binary.Read(bytes.NewBuffer(rinfo.CellOffsets), LE, offsets)
						}

						cellStorageBuf := rinfo.CellStorageBuffer
						if len(cellStorageBuf) == 0 && len(rinfo.CellStorageBufferPreBnc) > 0 {
							cellStorageBuf = rinfo.CellStorageBufferPreBnc
						}

						for c := 0; c < int(cc) && c < len(offsets); c++ {
							offset := offsets[c]
							if offset != 65535 && len(cellStorageBuf) >= int(offset)+16 {
								key := LE.Uint32(cellStorageBuf[offset+12 : offset+16])
								if key != 0 {
									columnUniqueKeys[c][key] = true
								}
							}
						}
					}
				}

				// Keep columns with diverse content
				for c := 0; c < int(cc); c++ {
					uniqueKeyCount := len(columnUniqueKeys[c])
					if c == 0 {
						hasContentColumns[c] = uniqueKeyCount > 0
					} else if uniqueKeyCount >= int(rows)*3/10 && uniqueKeyCount > 2 {
						hasContentColumns[c] = true
					}
				}

				// Removed: automatic second column addition
				// User feedback indicates only diverse columns should be shown
			} else {
				for i := 0; i < cc; i++ {
					hasContentColumns[i] = true
				}
			}
		}
	} else {
		// For Numbers and other types, show all columns
		for i := 0; i < cc; i++ {
			hasContentColumns[i] = true
		}
	}

	// Calculate the number of columns with content
	contentColumnCount := 0
	for _, hasContent := range hasContentColumns {
		if hasContent {
			contentColumnCount++
		}
	}

	// Only create columns that have content, skip empty columns entirely
	actualColumnCount := 0
	for i := 0; i < cc; i++ {
		if hasContentColumns[i] {
			actualColumnCount++
		}
	}

	// Only add colgroup if we have content columns
	if actualColumnCount > 0 {
		for i := 0; i < cc; i++ {
			if hasContentColumns[i] {
				var col *html.Node
				if actualColumnCount == 1 {
					// Single column, let it auto-size naturally
					col = E("col")
				} else {
					// Multiple columns, distribute width evenly
					width := 100.0 / float64(actualColumnCount)
					col = E("col", []string{"style", fmt.Sprintf("width: %.1f%%;", width)})
				}
				colgroup.AppendChild(col)
			}
		}
	}
	table.AppendChild(colgroup)

	if debugTableCells {
		fmt.Printf("DEBUG: Column width allocation - Content columns: %v, Count: %d\n", hasContentColumns, contentColumnCount)
	}

	// Construct thead/tbody, put header rows into thead
	thead := E("thead")
	tbody := E("tbody")
	table.AppendChild(thead)
	table.AppendChild(tbody)

	// Use global row numbers across tiles to determine header/footer, and correctly parse/fill each cell
	// Note: Remove key usage tracking because the same key may need to be used in multiple cells
	// usedKeys := make(map[uint32]bool) // Commented out to avoid blocking duplicate content

	// Determine if first row should be treated as header
	// Default to true for tables with multiple rows (common case)
	shouldTreatFirstRowAsHeader := true
	headerRowCount := 1

	if tm.NumberOfHeaderRows != nil {
		headerRowCount = int(*tm.NumberOfHeaderRows)
		shouldTreatFirstRowAsHeader = headerRowCount > 0
		if debugTableCells {
			fmt.Printf("DEBUG: NumberOfHeaderRows explicitly set to %d\n", headerRowCount)
		}
	} else {
		// Default: treat first row as header for multi-row tables
		if debugTableCells {
			fmt.Printf("DEBUG: NumberOfHeaderRows is nil - defaulting to 1 header row\n")
		}
	}

	// Use rowTileTree to get correct row order
	var orderedRows []struct {
		RowIndex uint32
		Tile     *TST.Tile
		RowInfo  *TST.TileRowInfo
	}

	// Note: Same key can be used in multiple cells (e.g., merged cells)

	// Build a map of tile ID to tile for quick lookup
	tileMap := make(map[uint32]*TST.Tile)
	for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
		tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
		tileMap[*tinfo.Tileid] = tile

		// Debug: show physical order of rows in tile
		if debugTableCells && len(tile.RowInfos) > 10 {
			fmt.Printf("DEBUG: Tile physical order (first 15 rows):\n")
			for i, rinfo := range tile.RowInfos {
				if i < 15 && rinfo.TileRowIndex != nil {
					fmt.Printf("DEBUG:   Physical[%d] -> TileRowIndex=%d\n", i, *rinfo.TileRowIndex)
				}
			}
		}
	}

	// Try to get row order from RowHeaders (similar to ColumnHeaders)
	// RowHeaders contains the actual display order with index field
	rowOrderMap := make(map[uint32]uint32) // maps TileRowIndex to display order
	hasRowHeaders := false

	if tm.BaseDataStore != nil && tm.BaseDataStore.RowHeaders != nil {
		// RowHeaders is a HeaderStorage, need to read buckets
		for _, bucketRef := range tm.BaseDataStore.RowHeaders.Buckets {
			if bucket, ok := ctx.ix.Deref(bucketRef).(*TST.HeaderStorageBucket); ok {
				if debugTableCells {
					fmt.Printf("DEBUG: RowHeaders bucket has %d headers\n", len(bucket.Headers))
				}
				for _, header := range bucket.Headers {
					if header.Index != nil {
						idx := *header.Index
						// Index field represents the display order
						rowOrderMap[idx] = idx
						hasRowHeaders = true
						if debugTableCells && idx < 15 {
							fmt.Printf("DEBUG: RowHeader[%d]: index=%d\n", idx, idx)
						}
					}
				}
			}
		}
	}

	if debugTableCells && hasRowHeaders {
		fmt.Printf("DEBUG: Found RowHeaders with %d row mappings\n", len(rowOrderMap))
	}

	// Use rowTileTree to get correct row order
	if tm.BaseDataStore.RowTileTree != nil {
		if debugTableCells {
			fmt.Printf("DEBUG: RowTileTree has %d nodes\n", len(tm.BaseDataStore.RowTileTree.Nodes))
			for i, node := range tm.BaseDataStore.RowTileTree.Nodes {
				fmt.Printf("DEBUG: RowTileTree[%d]: rowIndex=%d, tileId=%d\n", i, *node.Key, *node.Value)
			}
		}

		// Sort nodes by key (rowIndex) to get correct order
		nodes := make([]*TST.TableRBTree_Node, len(tm.BaseDataStore.RowTileTree.Nodes))
		copy(nodes, tm.BaseDataStore.RowTileTree.Nodes)
		sort.Slice(nodes, func(i, j int) bool {
			return *nodes[i].Key < *nodes[j].Key
		})

		if debugTableCells {
			fmt.Printf("DEBUG: Sorted RowTileTree nodes by rowIndex:\n")
			for i, node := range nodes {
				fmt.Printf("DEBUG: Sorted[%d]: rowIndex=%d, tileId=%d\n", i, *node.Key, *node.Value)
			}
		}

		for _, node := range nodes {
			rowIndex := *node.Key
			tileId := *node.Value
			if tile, exists := tileMap[tileId]; exists {
				// Find the corresponding row info within this tile
				for _, rinfo := range tile.RowInfos {
					if rinfo.TileRowIndex != nil && *rinfo.TileRowIndex == rowIndex {
						orderedRows = append(orderedRows, struct {
							RowIndex uint32
							Tile     *TST.Tile
							RowInfo  *TST.TileRowInfo
						}{rowIndex, tile, rinfo})
						break
					}
				}
			}
		}
	} else {
		if debugTableCells {
			fmt.Printf("DEBUG: RowTileTree is nil\n")
		}
	}

	// If rowTileTree is not available or empty, fall back to original method
	if len(orderedRows) == 0 {
		if debugTableCells {
			fmt.Printf("DEBUG: RowTileTree is empty or nil, using fallback method with original tile order\n")
		}

		// Fallback: Use original tile order (as stored in the file)
		for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
			tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
			for _, rinfo := range tile.RowInfos {
				// Use TileRowIndex as the row index for ordering
				rowIndex := uint32(0)
				if rinfo.TileRowIndex != nil {
					rowIndex = *rinfo.TileRowIndex
				}
				orderedRows = append(orderedRows, struct {
					RowIndex uint32
					Tile     *TST.Tile
					RowInfo  *TST.TileRowInfo
				}{rowIndex, tile, rinfo})
			}
		}

		// Sort by TileRowIndex to maintain consistent order
		sort.Slice(orderedRows, func(i, j int) bool {
			if orderedRows[i].RowInfo.TileRowIndex != nil && orderedRows[j].RowInfo.TileRowIndex != nil {
				return *orderedRows[i].RowInfo.TileRowIndex < *orderedRows[j].RowInfo.TileRowIndex
			}
			return i < j // Fallback to insertion order
		})

		if debugTableCells {
			fmt.Printf("DEBUG: Fallback method processed %d rows\n", len(orderedRows))
		}
	} else {
		// Check if rowTileTree contains all rows, if not, supplement with remaining rows
		expectedRows := 0
		for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
			tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
			expectedRows += len(tile.RowInfos)
		}

		if len(orderedRows) < expectedRows {
			if debugTableCells {
				fmt.Printf("DEBUG: RowTileTree only has %d rows, expected %d rows, supplementing with remaining rows\n", len(orderedRows), expectedRows)
			}

			// Create a map of already processed rows
			processedRows := make(map[uint32]bool)
			for _, row := range orderedRows {
				if row.RowInfo.TileRowIndex != nil {
					processedRows[*row.RowInfo.TileRowIndex] = true
				}
			}

			// Add remaining rows
			for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
				tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
				for _, rinfo := range tile.RowInfos {
					if rinfo.TileRowIndex != nil && !processedRows[*rinfo.TileRowIndex] {
						orderedRows = append(orderedRows, struct {
							RowIndex uint32
							Tile     *TST.Tile
							RowInfo  *TST.TileRowInfo
						}{*rinfo.TileRowIndex, tile, rinfo})
					}
				}
			}

			// Important: Sort all rows (including supplemented ones) by TileRowIndex
			sort.Slice(orderedRows, func(i, j int) bool {
				if orderedRows[i].RowInfo.TileRowIndex != nil && orderedRows[j].RowInfo.TileRowIndex != nil {
					return *orderedRows[i].RowInfo.TileRowIndex < *orderedRows[j].RowInfo.TileRowIndex
				}
				return i < j
			})

			if debugTableCells {
				fmt.Printf("DEBUG: After supplementing, sorted %d rows by TileRowIndex\n", len(orderedRows))
			}
		}
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Processing %d rows in correct order\n", len(orderedRows))
		for i, row := range orderedRows {
			fmt.Printf("DEBUG: Row %d: rowIndex=%d, tileRowIndex=%d\n", i, row.RowIndex, *row.RowInfo.TileRowIndex)
		}
	}

	for r, rowData := range orderedRows {
		rinfo := rowData.RowInfo

		// Special case: Skip rows with small buffers (20 bytes)
		// These appear to be placeholder/reference rows that don't contain actual displayable content
		if len(rinfo.CellStorageBuffer) == 20 {
			if debugTableCells && len(orderedRows) == 18 {
				fmt.Printf("DEBUG: Row %2d SKIPPED (small buffer=20 bytes, appears to be a reference row)\n", r)
			}
			continue // Skip this row entirely
		}

		tr := E("tr")
		// Header row determined by headerRowCount
		isHeaderRow := r < headerRowCount

		// Mark last row as summary/total row
		isSummaryRow := (r == len(orderedRows)-1)
		if isSummaryRow {
			tr.Attr = append(tr.Attr, html.Attribute{Key: "class", Val: "summary-row"})
		}

		if isHeaderRow {
			thead.AppendChild(tr)
		} else {
			tbody.AppendChild(tr)
		}

		// Decode the row's column -> offset mapping
		offsets := make([]uint16, len(rinfo.CellOffsets)/2)
		binary.Read(bytes.NewBuffer(rinfo.CellOffsets), LE, offsets)

		// Use original content offsets
		contentOffsets := offsets

		if debugTableCells && r < 3 {
			fmt.Printf("DEBUG: Row %d - Raw offsets: %v\n", r, offsets)
		}

		// First pass: read all keys to detect merged cells
		rowKeys := make([]uint32, cc)
		for c := 0; c < cc; c++ {
			if !hasContentColumns[c] {
				continue
			}

			var offset uint16 = 65535
			if c < len(contentOffsets) {
				offset = contentOffsets[c]
			}

			// Read key from buffer
			if offset != 65535 && len(rinfo.CellStorageBuffer) >= int(offset)+16 {
				rowKeys[c] = LE.Uint32(rinfo.CellStorageBuffer[offset+12 : offset+16])
			}
		}

		// Debug output for 18-row table (second table)
		if debugTableCells && len(orderedRows) == 18 && r < 18 {
			bufferSize := len(rinfo.CellStorageBuffer)
			tileRowIdx := uint32(0)
			if rinfo.TileRowIndex != nil {
				tileRowIdx = *rinfo.TileRowIndex
			}
			// Show all column keys for this row
			fmt.Printf("DEBUG: Row %2d all keys: %v (bufferSize=%d, TileRowIndex=%d)\n", r, rowKeys[:min(7, len(rowKeys))], bufferSize, tileRowIdx)
		}

		// Special handling for summary/total rows:
		// If col 2 onwards have the same key, merge col 1 onwards
		if r == len(orderedRows)-1 { // Last row
			allSameFromCol2 := true
			if len(rowKeys) > 2 {
				firstKey := rowKeys[2]
				for c := 3; c < len(rowKeys); c++ {
					if rowKeys[c] != firstKey {
						allSameFromCol2 = false
						break
					}
				}
				if allSameFromCol2 && firstKey != 0 {
					// Merge: set col 2-N to have the same key as col 1
					for c := 2; c < len(rowKeys); c++ {
						rowKeys[c] = rowKeys[1]
					}
					if debugTableCells {
						fmt.Printf("DEBUG: Row %d detected as summary row, merging col 1-%d\n", r, len(rowKeys)-1)
						fmt.Printf("DEBUG: Row %d adjusted keys: %v\n", r, rowKeys)
					}
				}
			}
		}

		// Only process columns that have content
		actualColIndex := 0
		for c := 0; c < cc; c++ {
			// Skip empty columns entirely
			if !hasContentColumns[c] {
				continue
			}

			// Check if it's a header cell
			var cellTag string
			if isHeaderRow {
				cellTag = "th"
				// Debug: print header style information
				if debugTableCells {
					fmt.Printf("DEBUG: Header row %d, col %d (actual col %d) - checking for header styles\n", r, c, actualColIndex)
				}
			} else {
				cellTag = "td"
			}

			td := E(cellTag)

			// Calculate colspan for merged cells
			// Check if this cell and next cells have the same non-zero key
			colspan := 1
			if rowKeys[c] != 0 {
				for nextCol := c + 1; nextCol < cc; nextCol++ {
					if !hasContentColumns[nextCol] {
						break
					}
					if rowKeys[nextCol] == rowKeys[c] {
						colspan++
					} else {
						break
					}
				}
			}

			// Apply colspan if > 1
			if colspan > 1 {
				td.Attr = append(td.Attr, html.Attribute{Key: "colspan", Val: fmt.Sprintf("%d", colspan)})
				if debugTableCells && r == 10 {
					fmt.Printf("DEBUG: Cell[%d,%d] colspan=%d (key=%d)\n", r, c, colspan, rowKeys[c])
				}
			}

			// Don't add extra header style classes, keep it simple
			// if isHeaderRow {
			//	td.Attr = append(td.Attr, html.Attribute{Key: "class", Val: "table-header"})
			// }
			// Add debug attributes
			td.Attr = append(td.Attr, html.Attribute{Key: "data-row", Val: fmt.Sprintf("%d", r)})
			// Use actualColIndex for the final table structure
			td.Attr = append(td.Attr, html.Attribute{Key: "data-col", Val: fmt.Sprintf("%d", actualColIndex)})
			tr.AppendChild(td)

			// Increment actualColIndex after processing this cell
			actualColIndex++

			// Apply basic positioning styles
			basicStyle := ctx.applyPositionBasedStyle(tm, r, c, shouldTreatFirstRowAsHeader)
			if basicStyle != "" {
				td.Attr = append(td.Attr, html.Attribute{Key: "style", Val: basicStyle})
			}

			// Use rearranged content offsets
			var offset uint16 = 65535 // Default to empty cell
			if len(contentOffsets) > 0 && c < len(contentOffsets) {
				offset = contentOffsets[c]
			}

			// Since all offsets are 65535 (empty), we need to use StringTable-based content allocation
			// This is the correct approach for Numbers tables

			// Redesigned: simplified cell content retrieval logic
			// No longer relies on complex cellType judgment, directly try all possible content retrieval methods
			if debugTableCells {
				fmt.Printf("DEBUG: Processing cell at row %d, col %d, offset=%d\n", r, c, offset)
				if int(offset)+16 <= len(rinfo.CellStorageBuffer) {
					fmt.Printf("DEBUG: Buffer[%d:%d] = %v\n", offset, offset+16, rinfo.CellStorageBuffer[offset:offset+16])
				} else {
					fmt.Printf("DEBUG: Buffer[%d:%d] = out of range (buffer len=%d)\n", offset, offset+16, len(rinfo.CellStorageBuffer))
				}
			}

			// Try to extract actual key value from buffer
			var key uint32
			contentFound := false

			// Parse cell type and key from buffer
			var cellType uint8

			// Determine which buffer to use: Regular or PreBnc
			// Strategy: Use Regular buffer by default, only use PreBnc as fallback
			usePreBnc := false
			var activeBuffer []byte
			var preBncOffset uint16 = 65535

			// First, try Regular buffer
			activeBuffer = rinfo.CellStorageBuffer

			// Only fallback to PreBnc if Regular buffer is empty/invalid for this cell
			if offset == 65535 || len(rinfo.CellStorageBuffer) < int(offset)+16 {
				// Check if PreBnc version exists and is valid for this cell
				if len(rinfo.CellOffsetsPreBnc) > 0 && len(rinfo.CellStorageBufferPreBnc) > 0 {
					preBncOffsets := make([]uint16, len(rinfo.CellOffsetsPreBnc)/2)
					binary.Read(bytes.NewBuffer(rinfo.CellOffsetsPreBnc), LE, preBncOffsets)
					if c < len(preBncOffsets) {
						preBncOffset = preBncOffsets[c]
						if preBncOffset != 65535 && len(rinfo.CellStorageBufferPreBnc) >= int(preBncOffset)+16 {
							usePreBnc = true
							activeBuffer = rinfo.CellStorageBufferPreBnc
							offset = preBncOffset
						}
					}
				}
			}

			// Debug output for first column
			if debugTableCells && c == 0 && r <= 10 {
				fmt.Printf("\nDEBUG: Cell[%d,%d] buffer analysis:\n", r, c)
				fmt.Printf("  Regular offset: %d\n", contentOffsets[c])
				fmt.Printf("  PreBnc offset:  %d\n", preBncOffset)
				fmt.Printf("  Using: %s\n", map[bool]string{true: "PreBnc", false: "Regular"}[usePreBnc])
			}

			// Only try buffer parsing if offset is not 65535
			if offset != 65535 && len(activeBuffer) >= int(offset)+16 {
				// Debug: show complete 16-byte structure for first column
				if debugTableCells && c == 0 && r <= 10 {
					buf := activeBuffer[offset : offset+16]
					fmt.Printf("  16-byte buffer: %v\n", buf)
					fmt.Printf("  offset+0:  %d (0x%08x)\n", LE.Uint32(buf[0:4]), LE.Uint32(buf[0:4]))
					fmt.Printf("  offset+4:  %d (0x%08x)\n", LE.Uint32(buf[4:8]), LE.Uint32(buf[4:8]))
					fmt.Printf("  offset+8:  %d (0x%08x)\n", LE.Uint32(buf[8:12]), LE.Uint32(buf[8:12]))
					fmt.Printf("  offset+12: %d (0x%08x) <-- key position\n", LE.Uint32(buf[12:16]), LE.Uint32(buf[12:16]))
					fmt.Printf("  byte[0]=5, byte[1]=%d (cellType indicator)\n", buf[1])
				}

				// Parse cellType from byte[1] (when byte[0]=5)
				if activeBuffer[offset] == 5 {
					cellType = activeBuffer[offset+1]
				}

				// For Pages tables: key is at offset+12 (last 4 bytes of 16-byte cell structure)
				key = LE.Uint32(activeBuffer[offset+12 : offset+16])

				if debugTableCells && r <= 10 {
					fmt.Printf("  Row %d, Col %d: cellType=%d, key=%d\n", r, c, cellType, key)
				}

				if key != 0 && key != 65535 {
					contentFound = true
				}
			} else {
				if debugTableCells && r < 6 {
					fmt.Printf("DEBUG: Cell[%d,%d] is empty (offset=65535 or insufficient buffer)\n", r, c)
				}
			}

			// If no valid key was found from buffer, this cell is empty
			if !contentFound || key == 0 {
				if debugTableCells && r < 6 {
					fmt.Printf("DEBUG: Cell at row %d, col %d is empty (no valid key found)\n", r, c)
				}
				// Add empty placeholder to ensure cell has minimum height
				emptyDiv := E("div")
				emptyDiv.Attr = append(emptyDiv.Attr, html.Attribute{Key: "style", Val: "min-height: 1.2em; line-height: 1.2;"})
				td.AppendChild(emptyDiv)
				continue
			}

			// Debug: print the key that will be used for content lookup (for column 0 of all rows in debug mode)
			if debugTableCells && c == 0 {
				fmt.Printf("DEBUG: Cell at row %d, col %d will use key %d for content lookup\n", r, c, key)
			}

			// Debug: check Numbers style fields
			if debugTableCells && r == 0 && c == 0 {
				fmt.Printf("DEBUG: Document type: %s\n", ctx.ix.Type)
				if ctx.ix.Type == "numbers" {
					fmt.Printf("DEBUG: Numbers style fields - HeaderRowStyle: %v, BodyCellStyle: %v, HeaderColumnStyle: %v, FooterRowStyle: %v\n",
						tm.HeaderRowStyle != nil, tm.BodyCellStyle != nil, tm.HeaderColumnStyle != nil, tm.FooterRowStyle != nil)
				}
			}

			// Apply cell-specific styles
			var cellStyle string

			// For Numbers documents, use the dedicated style fields
			if ctx.ix.Type == "numbers" {
				if r == 0 && tm.HeaderRowStyle != nil {
					// Header row style
					cellStyle = ctx.applyNumbersCellStyle(tm.HeaderRowStyle)
				} else if c == 0 && tm.HeaderColumnStyle != nil {
					// Header column style
					cellStyle = ctx.applyNumbersCellStyle(tm.HeaderColumnStyle)
				} else if tm.NumberOfRows != nil && r == int(*tm.NumberOfRows)-1 && tm.FooterRowStyle != nil {
					// Footer row style
					cellStyle = ctx.applyNumbersCellStyle(tm.FooterRowStyle)
				} else if tm.BodyCellStyle != nil {
					// Body cell style
					cellStyle = ctx.applyNumbersCellStyle(tm.BodyCellStyle)
				}
			} else {
				// For Pages documents, use the old StyleTable approach
				styleKey := key
				if r == 0 && tm.BaseDataStore != nil && tm.BaseDataStore.StyleTable != nil {
					// Try to find a style with background color for header row
					styleTableRef := ctx.ix.Deref(tm.BaseDataStore.StyleTable)
					if tdl, ok := styleTableRef.(*TST.TableDataList); ok {
						for _, entry := range tdl.Entries {
							if entry.Key != nil && entry.Reference != nil {
								entryRef := ctx.ix.Deref(entry.Reference)
								if csa, ok := entryRef.(*TST.CellStyleArchive); ok {
									if csa.CellProperties != nil && csa.CellProperties.CellFill != nil {
										// Found a style with background color, use it for header
										styleKey = *entry.Key
										if debugTableCells {
											fmt.Printf("DEBUG: Using style key %d for header row (has background color)\n", styleKey)
										}
										break
									}
								}
							}
						}
					}
				}

				// For non-header rows, ensure they don't use background color styles
				if r > 0 && tm.BaseDataStore != nil && tm.BaseDataStore.StyleTable != nil {
					styleTableRef := ctx.ix.Deref(tm.BaseDataStore.StyleTable)
					if tdl, ok := styleTableRef.(*TST.TableDataList); ok {
						// Try to find a plain text style (ParagraphStyleArchive) for content rows
						for _, entry := range tdl.Entries {
							if entry.Key != nil && entry.Reference != nil {
								entryRef := ctx.ix.Deref(entry.Reference)
								if _, ok := entryRef.(*TSWP.ParagraphStyleArchive); ok {
									// Found a paragraph style, use it for content rows
									styleKey = *entry.Key
									if debugTableCells {
										fmt.Printf("DEBUG: Using style key %d for content row %d (plain text style)\n", styleKey, r)
									}
									break
								}
							}
						}
					}
				}

				if styleKey != 0 {
					cellStyle = ctx.applyCellStyle(tm, styleKey)
				}
			}

			if cellStyle != "" {
				// Merge with existing basic style
				existingStyle := ""
				for _, attr := range td.Attr {
					if attr.Key == "style" {
						existingStyle = attr.Val
						break
					}
				}
				if existingStyle != "" {
					cellStyle = existingStyle + ";" + cellStyle
				}

				// Update or add style attribute
				styleFound := false
				for i, attr := range td.Attr {
					if attr.Key == "style" {
						td.Attr[i].Val = cellStyle
						styleFound = true
						break
					}
				}
				if !styleFound {
					td.Attr = append(td.Attr, html.Attribute{Key: "style", Val: cellStyle})
				}
				if debugTableCells {
					fmt.Printf("DEBUG: Applied cell-specific style for key %d: %s\n", key, cellStyle)
				}
			} else {
				if debugTableCells {
					fmt.Printf("DEBUG: No cell-specific style found for key %d\n", key)
				}
			}

			// Render content
			contentFound = false

			if debugTableCells {
				fmt.Printf("DEBUG: Looking for content with key %d for cell at row %d, col %d\n", key, r, c)
			}

			// Normal content rendering
			// Strategy depends on cellType:
			// - cellType=9: richText content, look in richTable first
			// - cellType=2: may be a simple number, try stringTable first, then check if key is a literal number
			// - cellType=3: string content, look in stringTable

			if cellType == 9 {
				// Rich text: try richTable first
				for _, entry := range richTable {
					if *entry.Key == key {
						if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
							if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
								if debugTableCells {
									fmt.Printf("DEBUG: Found rich text content for key %d\n", key)
								}
								ctx.storageToNodeForTable(st, td)
								contentFound = true
							}
						}
						break
					}
				}
				// Fallback to stringTable
				if !contentFound {
					for _, entry := range stringTable {
						if *entry.Key == key {
							td.AppendChild(T(*entry.String_))
							if debugTableCells {
								fmt.Printf("DEBUG: Found string content for key %d: %s\n", key, *entry.String_)
							}
							contentFound = true
							break
						}
					}
				}
			} else if cellType == 2 {
				// cellType=2: simple number/sequence cell
				// For key values: use key as literal number
				// Wrap in a div with proper styling for consistency with richText cells
				numDiv := E("div")
				numDiv.Attr = append(numDiv.Attr, html.Attribute{
					Key: "style",
					Val: "min-height: 1.2em; line-height: 1.2;",
				})
				numDiv.AppendChild(T(fmt.Sprintf("%d", key)))
				td.AppendChild(numDiv)

				if debugTableCells {
					fmt.Printf("DEBUG: cellType=2, using key as literal number: %d\n", key)
				}
				contentFound = true
			} else if cellType == 0 {
				// cellType=0: control cell (popup menu, checkbox, etc.)
				// Try to get actual option text from PopUpMenuModel
				optionText := ""

				// Try to get PopUpMenuModel from DataStore
				if tm.BaseDataStore != nil && tm.BaseDataStore.MultipleChoiceListFormatTable != nil {
					// MultipleChoiceListFormatTable points to TableDataList, not directly to PopUpMenuModel
					if tdl, ok := ctx.ix.Deref(tm.BaseDataStore.MultipleChoiceListFormatTable).(*TST.TableDataList); ok {
						// Look for PopUpMenuModel in the TableDataList entries
						for _, entry := range tdl.Entries {
							if entry != nil && entry.Key != nil && entry.Reference != nil {
								if popup, ok := ctx.ix.Deref(entry.Reference).(*TST.PopUpMenuModel); ok {
									if key > 0 && int(key) <= len(popup.Item) {
										// key is 1-based index into PopUpMenuModel.Item
										itemIndex := int(key) - 1
										if itemIndex >= 0 && itemIndex < len(popup.Item) {
											item := popup.Item[itemIndex]
											if item != nil && item.StringValue != nil && item.StringValue.Value != nil {
												optionText = *item.StringValue.Value
												if debugTableCells {
													fmt.Printf("DEBUG: cellType=0 found option text from PopUpMenuModel: %s (key=%d, index=%d)\n", optionText, key, itemIndex)
												}
												break // Found the option, no need to check other entries
											}
										}
									}
								}
							}
						}
					}
				}

				controlDiv := E("div")
				controlDiv.Attr = append(controlDiv.Attr, html.Attribute{
					Key: "style",
					Val: "min-height: 1.2em; line-height: 1.2;",
				})
				controlDiv.AppendChild(T(optionText))
				td.AppendChild(controlDiv)

				if debugTableCells {
					fmt.Printf("DEBUG: cellType=0 (control cell), key=%d, text=%s\n", key, optionText)
				}
				contentFound = true
			} else if cellType == 5 {
				// cellType=5: date/time cell
				// Parse date from cell storage buffer using iWork format
				dateText := ""

				// Try to extract date from key value or buffer
				if key != 0 {
					// Method 1: Try to find DateCellValueArchive by key in format table
					if tm.BaseDataStore != nil && tm.BaseDataStore.FormatTable != nil {
						if tdl, ok := ctx.ix.Deref(tm.BaseDataStore.FormatTable).(*TST.TableDataList); ok {
							for _, entry := range tdl.Entries {
								if entry != nil && entry.Key != nil && *entry.Key == key {
									if entry.Reference != nil {
										// Try to deref as DateCellValueArchive
										if dateVal, ok := ctx.ix.Deref(entry.Reference).(*TSCE.DateCellValueArchive); ok {
											if dateVal.Value != nil {
												// Convert Apple timestamp to Unix timestamp
												value := *dateVal.Value + 978307200 // Apple to unix epoch
												tm := time.Unix(int64(value), 0)
												dateText = tm.Format("2006-01-02")
												if debugTableCells {
													fmt.Printf("DEBUG: cellType=5 found date from FormatTable: %s (key=%d, value=%f)\n", dateText, key, *dateVal.Value)
												}
												break
											}
										} else {
											if debugTableCells {
												fmt.Printf("DEBUG: Unrecognized type: %T\n", ctx.ix.Deref(entry.Reference))
											}
										}
									}
								} else {
									if debugTableCells {
										fmt.Printf("DEBUG: entry.Reference is nil for key %d=%d, %T\n", *entry.Key, key, ctx.ix.Deref(entry.Reference))
									}
								}
							}
						} else {
							if debugTableCells {
								fmt.Printf("DEBUG: Unrecognized type: %T\n", ctx.ix.Deref(tm.BaseDataStore.FormatTable))
							}
						}
					} else {
						if debugTableCells {
							fmt.Printf("DEBUG: Unrecognized type: %T\n", tm.BaseDataStore.FormatTable)
						}
					}

					// Method 2: Fallback to buffer extraction
					if dateText == "" && len(activeBuffer) >= int(offset)+20 {
						value := math.Float64frombits(LE.Uint64(activeBuffer[offset+12 : offset+20]))
						if value != 0 {
							tm := time.Unix(int64(value+978307200), 0) // Apple to unix epoch 2001-01-01
							dateText = tm.Format("2006-01-02")

							if debugTableCells {
								fmt.Printf("DEBUG: cellType=5 found date from buffer: %s (key=%d, value=%f)\n", dateText, key, value)
							}
						}
					}
				}

				dateDiv := E("div")
				dateDiv.Attr = append(dateDiv.Attr, html.Attribute{
					Key: "style",
					Val: "min-height: 1.2em; line-height: 1.2;",
				})
				dateDiv.AppendChild(T(dateText))
				td.AppendChild(dateDiv)

				if debugTableCells {
					fmt.Printf("DEBUG: cellType=5 (date cell), key=%d, text=%s\n", key, dateText)
				}
				contentFound = true
			} else if cellType == 3 {
				// cellType=3: string content from tables
				// Try stringTable first
				for _, entry := range stringTable {
					if *entry.Key == key {
						td.AppendChild(T(*entry.String_))
						if debugTableCells {
							fmt.Printf("DEBUG: Found string content for key %d: %s\n", key, *entry.String_)
						}
						contentFound = true
						break
					}
				}

				// Fallback to richTable
				if !contentFound {
					for _, entry := range richTable {
						if *entry.Key == key {
							if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
								if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
									if debugTableCells {
										fmt.Printf("DEBUG: Found rich text content for key %d\n", key)
									}
									ctx.storageToNodeForTable(st, td)
									contentFound = true
								}
							}
							break
						}
					}
				}
			} else {
				// Unknown cellType: try both tables
				for _, entry := range richTable {
					if *entry.Key == key {
						if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
							if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
								if debugTableCells {
									fmt.Printf("DEBUG: Found rich text content for key %d\n", key)
								}
								ctx.storageToNodeForTable(st, td)
								contentFound = true
							}
						}
						break
					}
				}
				if !contentFound {
					for _, entry := range stringTable {
						if *entry.Key == key {
							td.AppendChild(T(*entry.String_))
							if debugTableCells {
								fmt.Printf("DEBUG: Found string content for key %d: %s\n", key, *entry.String_)
							}
							contentFound = true
							break
						}
					}
				}
			}

			if !contentFound {
				if debugTableCells {
					fmt.Printf("DEBUG: No content found for cell at row %d, col %d with key %d\n", r, c, key)
				}
			}

			// Skip merged cells (cells with same key as current cell)
			if colspan > 1 {
				c += (colspan - 1)
				if debugTableCells && r == 10 {
					fmt.Printf("DEBUG: Skipping %d merged cells, next c=%d\n", colspan-1, c+1)
				}
			}
		}
	}
	rval := E("div")
	if tm.TableName != nil {
		// rval.AppendChild(E("h3", *tm.TableName)) // no need to show table name
	}
	rval.AppendChild(table)
	return rval
}

func (ctx *Context) processDrawable(ref *TSP.Reference) *html.Node {
	item := ctx.ix.Deref(ref)
	switch item.(type) {
	case *TSD.ImageArchive:
		img := item.(*TSD.ImageArchive)
		node := ctx.processImage(img)
		if img.Super != nil && img.Super.Geometry != nil {
			// Detect full-page background image by geometry ≈ canvas
			canvasW := 1920.0
			canvasH := 1080.0
			for _, rec := range ctx.ix.Records {
				if sh, ok := rec.(*KN.ShowArchive); ok {
					if sh.Size != nil && sh.Size.Width != nil && sh.Size.Height != nil {
						canvasW = float64(*sh.Size.Width)
						canvasH = float64(*sh.Size.Height)
					}
					break
				}
			}
			isFull := false
			if g := img.Super.Geometry; g != nil && g.Position != nil && g.Size != nil &&
				g.Position.X != nil && g.Position.Y != nil && g.Size.Width != nil && g.Size.Height != nil {
				x := float64(*g.Position.X)
				y := float64(*g.Position.Y)
				w := float64(*g.Size.Width)
				h := float64(*g.Size.Height)
				if math.Abs(x) < 1e-2 && math.Abs(y) < 1e-2 &&
					math.Abs(w-canvasW) < 1e-1 && math.Abs(h-canvasH) < 1e-1 {
					isFull = true
				}
			}
			if isFull {
				// Mark as slide background image
				node.Attr = append(node.Attr, html.Attribute{Key: "class", Val: "background-img"})
				return node
			}
			return ctx.wrapWithGeometry(node, img.Super.Geometry, "", false)
		}
		return node
	case *TST.WPTableInfoArchive:
		table := item.(*TST.WPTableInfoArchive)
		tm := ctx.ix.Deref(table.Super.TableModel).(*TST.TableModelArchive)
		return ctx.processTable(tm)
	case *TST.TableInfoArchive:
		if debugTableCells {
			fmt.Printf("DEBUG: Processing Numbers TableInfoArchive\n")
		}
		tm := ctx.ix.Deref(item.(*TST.TableInfoArchive).TableModel).(*TST.TableModelArchive)
		return ctx.processTable(tm)
	case *TSWP.ShapeInfoArchive:
		sia := item.(*TSWP.ShapeInfoArchive)
		node := ctx.processShapeInfo(sia)
		fillCSS := ""
		strokeCSS := ""
		if sia.Super != nil && sia.Super.Style != nil {
			styleAny := ctx.ix.Deref(sia.Super.Style)
			if ss, ok := styleAny.(*TSD.ShapeStyleArchive); ok {
				if ss.ShapeProperties != nil && ss.ShapeProperties.Fill != nil {
					if css := colorToCSS(ss.ShapeProperties.Fill.GetColor()); css != "" {
						fillCSS = css
					} else if g := ss.ShapeProperties.Fill.Gradient; g != nil {
						// simple linear-gradient from first->last stop
						stops := g.GetStops()
						if len(stops) >= 2 {
							c1 := colorToCSS(stops[0].GetColor())
							c2 := colorToCSS(stops[len(stops)-1].GetColor())
							if c1 != "" && c2 != "" {
								fillCSS = fmt.Sprintf("linear-gradient(%s, %s)", c1, c2)
							}
						}
					} else if img := ss.ShapeProperties.Fill.Image; img != nil && img.Imagedata != nil && img.Imagedata.Identifier != nil {
						// mark to receive background image later on wrapper
						ctx.imgs[fmt.Sprintf("Data/%d", *img.Imagedata.Identifier)] = *img.Imagedata.Identifier
						fillCSS = fmt.Sprintf("url(#bgimg_%d)", *img.Imagedata.Identifier)
					}
				}
				if ss.ShapeProperties != nil && ss.ShapeProperties.Stroke != nil {
					if col := colorToCSS(ss.ShapeProperties.Stroke.Color); col != "" {
						w := 1.0
						if ss.ShapeProperties.Stroke.Width != nil {
							w = float64(*ss.ShapeProperties.Stroke.Width)
						}
						strokeCSS = fmt.Sprintf("border: %.2fpx solid %s; box-sizing: border-box;", w, col)

						// Check if there are rounded corners
						if ss.ShapeProperties.Stroke.Join != nil && *ss.ShapeProperties.Stroke.Join == TSD.LineJoin_RoundJoin {
							strokeCSS += "border-radius: 4px;"
						}
					}
				}
			} else if swp, ok := styleAny.(*TSWP.ShapeStyleArchive); ok {
				if swp.GetSuper() != nil && swp.GetSuper().ShapeProperties != nil && swp.GetSuper().ShapeProperties.Fill != nil {
					if css := colorToCSS(swp.GetSuper().ShapeProperties.Fill.GetColor()); css != "" {
						fillCSS = css
					}
				}
				if sps := swp.GetSuper().ShapeProperties; sps != nil && sps.Stroke != nil {
					if col := colorToCSS(sps.Stroke.Color); col != "" {
						w := 1.0
						if sps.Stroke.Width != nil {
							w = float64(*sps.Stroke.Width)
						}
						strokeCSS = fmt.Sprintf("border: %.2fpx solid %s; box-sizing: border-box;", w, col)

						// Check if there are rounded corners
						if sps.Stroke.Join != nil && *sps.Stroke.Join == TSD.LineJoin_RoundJoin {
							strokeCSS += "border-radius: 4px;"
						}
					}
				}
			} else {
				// Print unrecognized shape style type
				fmt.Printf("DEBUG: Unrecognized shape style type: %T\n", styleAny)
			}
		}
		// Pages text boxes sometimes provide fill/stroke through paragraph styles
		if fillCSS == "" || strokeCSS == "" {
			if sia.OwnedStorage != nil {
				if stor, ok := ctx.ix.Deref(sia.OwnedStorage).(*TSWP.StorageArchive); ok && stor.TableParaStyle != nil && len(stor.TableParaStyle.Entries) > 0 {
					if stor.TableParaStyle.Entries[0].Object != nil {
						if psa, ok := ctx.ix.Deref(stor.TableParaStyle.Entries[0].Object).(*TSWP.ParagraphStyleArchive); ok && psa.ParaProperties != nil {
							if fillCSS == "" && psa.ParaProperties.Fill != nil {
								if css := colorToCSS(psa.ParaProperties.Fill); css != "" {
									fillCSS = css
								}
							}
							if strokeCSS == "" && psa.ParaProperties.Stroke != nil {
								col := colorToCSS(psa.ParaProperties.Stroke.Color)
								w := 1.0
								if psa.ParaProperties.Stroke.Width != nil {
									w = float64(*psa.ParaProperties.Stroke.Width)
								}
								if col != "" {
									strokeCSS = fmt.Sprintf("border: %.2fpx solid %s; box-sizing: border-box;", w, col)
								}
							}
						}
					}
					// Fallback: if still no fill, try first character style background color as text box background
					if fillCSS == "" && stor.TableCharStyle != nil && len(stor.TableCharStyle.Entries) > 0 {
						if stor.TableCharStyle.Entries[0].Object != nil {
							if csa, ok := ctx.ix.Deref(stor.TableCharStyle.Entries[0].Object).(*TSWP.CharacterStyleArchive); ok && csa.CharProperties != nil {
								if bc := csa.CharProperties.GetBackgroundColor(); bc != nil {
									if css := colorToCSS(bc); css != "" {
										fillCSS = css
									}
								}
							}
						}
					}
				}
			}
		}
		if sia.Super != nil && sia.Super.Super != nil && sia.Super.Super.Geometry != nil {
			extra := ""
			var bgClass string
			if fillCSS != "" {
				if strings.HasPrefix(fillCSS, "url(#bgimg_") {
					bgClass = "bgimg_" + strings.TrimSuffix(strings.TrimPrefix(fillCSS, "url(#bgimg_"), ")")
					extra += "background-size:cover;background-position:center;"
				} else {
					// Apply background color for text boxes
					extra += "background:" + fillCSS + ";"
				}
			}
			if strokeCSS != "" {
				extra += strokeCSS
			}
			wrapper := ctx.wrapWithGeometry(node, sia.Super.Super.Geometry, extra, true)
			if bgClass != "" {
				wrapper.Attr = append(wrapper.Attr, html.Attribute{Key: "class", Val: bgClass})
			}
			return wrapper
		}
		return node
	case *TSD.GroupArchive:
		return ctx.processDrawableArchive(item.(*TSD.GroupArchive).Super)
	case *KN.PlaceholderArchive:
		return ctx.processShapeInfo(item.(*KN.PlaceholderArchive).Super)
	default:
		msg := fmt.Sprintf("*** Unhandled attachment type %T\n", item)
		fmt.Println(msg)
		return E("div", msg)
	}
}

func (ctx *Context) processShapeInfo(sia *TSWP.ShapeInfoArchive) *html.Node {
	if debugMode {
		fmt.Printf("DEBUG: Processing ShapeInfo\n")
	}
	containedStorageRef := ctx.ix.Deref(sia.OwnedStorage)
	if cs, ok := containedStorageRef.(*TSWP.StorageArchive); ok {
		if debugMode {
			fmt.Printf("DEBUG: Found ContainedStorage with text: %s\n", cs.Text)
		}
		div := E("div")
		if ctx.storageToNode(cs, div) == nil {
			return div
		}
	} else {
		// Print unrecognized ContainedStorage type
		fmt.Printf("DEBUG: Unrecognized ContainedStorage type: %T\n", containedStorageRef)
	}
	return ctx.processDrawableArchive(sia.Super.Super)
}

func (ctx *Context) processDrawableArchive(da *TSD.DrawableArchive) *html.Node {
	if da == nil {
		return nil
	}

	// Process drawable object geometry information
	if da.Geometry != nil {
		// Create basic drawable element container
		container := E("div", []string{"class", "drawable-archive"})

		// DrawableArchive has no direct style fields, styles are handled through other means

		// Apply geometry styles
		return ctx.wrapWithGeometry(container, da.Geometry, "", false)
	}

	return nil
}

// processNumberAttachment processes number attachments
func (ctx *Context) processNumberAttachment(na *TSWP.NumberAttachmentArchive) *html.Node {
	if na.Super == nil {
		return nil
	}

	// Create number display element
	numberNode := E("span", []string{"class", "number-attachment"})

	// If there's a string value, display it
	if na.StringValue != nil {
		numberNode.AppendChild(T(*na.StringValue))
	} else {
		numberNode.AppendChild(T("0"))
	}

	return numberNode
}

// wrapWithGeometry wraps a child node with an absolutely positioned container based on TSD.GeometryArchive.
// extraStyle can include any CSS declarations, e.g. "background:rgba(...);border:1px solid red; display:flex;"
// isTextBox indicates if this is a text box that needs y-position adjustment based on font size
func (ctx *Context) wrapWithGeometry(child *html.Node, geom *TSD.GeometryArchive, extraStyle string, isTextBox bool) *html.Node {
	if geom == nil || geom.Position == nil || geom.Size == nil {
		return child
	}

	// Determine target CSS size and canvas size for scaling
	var targetWidthCSS float64
	var canvasW float64
	var canvasH float64
	if ctx.ix.Type == "key" {
		targetWidthCSS = 1200.0
		canvasW = 1920.0
		canvasH = 1080.0
		for _, rec := range ctx.ix.Records {
			if sh, ok := rec.(*KN.ShowArchive); ok {
				if sh.Size != nil && sh.Size.Width != nil && sh.Size.Height != nil {
					canvasW = float64(*sh.Size.Width)
					canvasH = float64(*sh.Size.Height)
				}
				break
			}
		}
	} else {
		// Pages: approximate A4 points (72dpi): 595 x 842
		targetWidthCSS = 900.0
		canvasW = 595.0
		canvasH = 842.0
	}
	slideHeightCSS := targetWidthCSS * canvasH / canvasW

	sx := targetWidthCSS / canvasW
	sy := slideHeightCSS / canvasH

	x := float64(0)
	y := float64(0)
	w := float64(0)
	h := float64(0)
	angle := float64(0)
	if geom.Position.X != nil {
		x = float64(*geom.Position.X) * sx
	}
	if geom.Position.Y != nil {
		y = float64(*geom.Position.Y) * sy

		// For text boxes, only apply a very small adjustment for large fonts
		if isTextBox && child != nil {
			fontSize := ctx.getTextBoxFontSize(child)
			// Only adjust for very large fonts (>50pt) to avoid over-correction
			if fontSize > 50 {
				adjustment := fontSize * 0.05
				y -= adjustment
			}
		}
	}
	if geom.Size.Width != nil {
		w = float64(*geom.Size.Width) * sx
	}
	if geom.Size.Height != nil {
		h = float64(*geom.Size.Height) * sy
	}
	if geom.Angle != nil {
		angle = float64(*geom.Angle)
	}

	style := fmt.Sprintf("position:absolute; left:%.2fpx; top:%.2fpx; width:%.2fpx; height:%.2fpx;", x, y, w, h)
	if angle != 0 {
		style += fmt.Sprintf(" transform: rotate(%.6frad); transform-origin: 0 0;", angle)
	}
	if extraStyle != "" {
		if !strings.HasSuffix(extraStyle, ";") {
			extraStyle += ";"
		}
		style += " " + extraStyle
	}

	wrapper := E("div", []string{"style", style, "class", "geom"})
	if child != nil {
		if child.Type == html.ElementNode && child.Data == "img" {
			// Make placed images scale to the geometry box while preserving aspect
			// Inline style to override any width/height attributes
			child.Attr = append(child.Attr, html.Attribute{Key: "style", Val: "width:100%;height:100%;object-fit:contain;display:block;"})
		}
		wrapper.AppendChild(child)
	}
	return wrapper
}

// getTextBoxFontSize extracts font size from text box content
func (ctx *Context) getTextBoxFontSize(node *html.Node) float64 {
	if node == nil {
		return 0
	}

	// Look for font-size in style attributes
	for _, attr := range node.Attr {
		if attr.Key == "style" {
			// Parse CSS for font-size
			if strings.Contains(attr.Val, "font-size:") {
				// Extract font-size value using regex
				re := regexp.MustCompile(`font-size:\s*([0-9.]+)pt`)
				matches := re.FindStringSubmatch(attr.Val)
				if len(matches) > 1 {
					if size, err := strconv.ParseFloat(matches[1], 64); err == nil {
						return size
					}
				}
			}
		}
	}

	// Look in child nodes for font-size
	var child *html.Node
	for child = node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			for _, attr := range child.Attr {
				if attr.Key == "class" {
					// Check if this class has font-size defined
					if style, exists := ctx.styles[attr.Val]; exists {
						if strings.Contains(style, "font-size:") {
							re := regexp.MustCompile(`font-size:\s*([0-9.]+)pt`)
							matches := re.FindStringSubmatch(style)
							if len(matches) > 1 {
								if size, err := strconv.ParseFloat(matches[1], 64); err == nil {
									return size
								}
							}
						}
					}
				}
			}
		}
	}

	return 0
}

// getTextBoxLineCount counts the number of lines in text box content
func (ctx *Context) getTextBoxLineCount(node *html.Node) int {
	if node == nil {
		return 1
	}

	lineCount := 1
	var child *html.Node
	for child = node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			if child.Data == "p" {
				lineCount++
			} else if child.Data == "br" {
				lineCount++
			}
		} else if child.Type == html.TextNode {
			// Count newlines in text content
			text := child.Data
			lineCount += strings.Count(text, "\n")
		}
	}

	return lineCount
}

// getFloatValue safely gets float value from pointer
func getFloatValue(f *float32) float64 {
	if f == nil {
		return 0.0
	}
	return float64(*f)
}

// mapBlueToGray maps blue theme colors to gray equivalents
func mapBlueToGray(r, g, b, a float64) (float64, float64, float64, float64) {
	// Check if this matches the blue header color
	if abs(r-0.357) < 0.001 && abs(g-0.608) < 0.001 && abs(b-0.835) < 0.001 {
		// Convert to gray header color (medium gray)
		return 0.6, 0.6, 0.6, a
	}
	// Check if this matches the light blue body color
	if abs(r-0.817) < 0.001 && abs(g-0.869) < 0.001 && abs(b-0.938) < 0.001 {
		// Convert to light gray body color
		return 0.9, 0.9, 0.9, a
	}
	// Return original colors if no match
	return r, g, b, a
}

// abs returns absolute value of float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// colorToCSS converts TSP.Color to CSS rgba() string.
func colorToCSS(c *TSP.Color) string {
	if c == nil {
		return ""
	}

	// Check color model
	model := TSP.Color_rgb // default
	if c.Model != nil {
		model = *c.Model
	}

	var r, g, b, a float64 = 0, 0, 0, 1

	switch model {
	case TSP.Color_rgb:
		// RGB model
		if c.R != nil {
			r = float64(*c.R)
		}
		if c.G != nil {
			g = float64(*c.G)
		}
		if c.B != nil {
			b = float64(*c.B)
		}
		if c.A != nil {
			a = float64(*c.A)
		}
	case TSP.Color_white:
		// Grayscale model - use W value for all RGB components
		if c.W != nil {
			gray := float64(*c.W)
			r, g, b = gray, gray, gray
		}
		if c.A != nil {
			a = float64(*c.A)
		}
	case TSP.Color_cmyk:
		// CMYK to RGB conversion
		var cy, ma, ye, k float64 = 0, 0, 0, 0
		if c.C != nil {
			cy = float64(*c.C)
		}
		if c.M != nil {
			ma = float64(*c.M)
		}
		if c.Y != nil {
			ye = float64(*c.Y)
		}
		if c.K != nil {
			k = float64(*c.K)
		}
		if c.A != nil {
			a = float64(*c.A)
		}

		// Convert CMYK to RGB
		r = (1 - cy) * (1 - k)
		g = (1 - ma) * (1 - k)
		b = (1 - ye) * (1 - k)
	}

	// Apply blue-to-gray mapping
	r, g, b, a = mapBlueToGray(r, g, b, a)

	// Clamp 0..1 then convert to 0..255
	clamp := func(v float64) float64 {
		if v < 0 {
			return 0
		}
		if v > 1 {
			return 1
		}
		return v
	}
	r255 := clamp(r) * 255
	g255 := clamp(g) * 255
	b255 := clamp(b) * 255
	return fmt.Sprintf("rgba(%d,%d,%d,%.3f)", int(r255+0.5), int(g255+0.5), int(b255+0.5), clamp(a))
}

// storageToNode populates a html node with the contents of a StorageArchive. This happens with both the
// main body of the document and rich text table cells.
// storageToNodeForTable processes table cell content without creating paragraph elements
func (ctx *Context) storageToNodeForTable(bs *TSWP.StorageArchive, td *html.Node) error {
	texts := bs.Text

	if len(texts) == 0 {
		// Add an empty placeholder to ensure cell has content and calculate height
		emptyDiv := E("div")
		emptyDiv.Attr = append(emptyDiv.Attr, html.Attribute{Key: "style", Val: "min-height: 1.2em; line-height: 1.2;"})
		td.AppendChild(emptyDiv)
		return nil
	}

	// Process multiple text fragments
	var text string
	if len(texts) == 1 {
		text = texts[0]
	} else {
		// Merge multiple text fragments
		for _, t := range texts {
			text += t
		}
	}

	// Check if text is empty or contains only whitespace
	if strings.TrimSpace(text) == "" {
		// Add an empty placeholder to ensure cell has content and calculate height
		emptyDiv := E("div")
		emptyDiv.Attr = append(emptyDiv.Attr, html.Attribute{Key: "style", Val: "min-height: 1.2em; line-height: 1.2;"})
		td.AppendChild(emptyDiv)
		return nil
	}

	// Offsets are in terms of unicode runes, so we have to convert to runes
	rr := []rune(text)

	// <p>
	parStyles := bs.TableParaStyle.Entries

	var attachments []Attachment
	if bs.TableAttachment != nil {
		for _, entry := range bs.TableAttachment.Entries {
			pos := *entry.CharacterIndex
			switch ctx.ix.Deref(entry.Object).(type) {
			case *TSWP.DrawableAttachmentArchive:
				archive := ctx.ix.Deref(entry.Object).(*TSWP.DrawableAttachmentArchive)
				node := ctx.processDrawable(archive.Drawable)
				if node != nil {
					attachments = append(attachments, Attachment{pos, node})
				}
			case *TSWP.NumberAttachmentArchive:
				archive := ctx.ix.Deref(entry.Object).(*TSWP.NumberAttachmentArchive)
				node := ctx.processNumberAttachment(archive)
				if node != nil {
					attachments = append(attachments, Attachment{pos, node})
				}
			}
		}
	}
	// bs.TableListStyle - seems to change on headings, look into it.

	// build content without paragraphs
	if len(parStyles) == 0 {
		// If no paragraph styles, directly add text content
		td.AppendChild(T(text))
		return nil
	}

	// Process paragraphs, wrap consecutive list items in appropriate lists
	var currentList *html.Node
	var currentListType string
	var inList bool

	for i, e := range parStyles {
		pos := *e.CharacterIndex
		end := uint32(len(rr))
		if i+1 < len(parStyles) {
			end = *parStyles[i+1].CharacterIndex
		}

		// Check if current paragraph is a list item
		var isListItem bool
		var listType string
		if e.Object != nil {
			objRef := ctx.ix.Deref(e.Object)
			if psa, ok := objRef.(*TSWP.ParagraphStyleArchive); ok {
				if psa != nil && psa.ParaProperties != nil {
					// Check ListStyleNull field
					if psa.ParaProperties.ListStyleNull == nil || !*psa.ParaProperties.ListStyleNull {
						// ListStyle is not null, check if there's a ListStyle
						if psa.ParaProperties.ListStyle != nil {
							// Check if ListStyle has label types, consistent with normal paragraph processing
							if ls, ok := ctx.ix.Deref(psa.ParaProperties.ListStyle).(*TSWP.ListStyleArchive); ok {
								if len(ls.LabelTypes) > 0 {
									isListItem = true
									labelType := ls.LabelTypes[0]
									switch labelType {
									case TSWP.ListStyleArchive_kNumber:
										listType = "ol"
									case TSWP.ListStyleArchive_kString, TSWP.ListStyleArchive_kImage:
										listType = "ul"
									default:
										listType = "ul"
									}
									if debugTableCells {
										fmt.Printf("DEBUG: storageToNodeForTable detected list item with type: %s\n", listType)
									}
								}
							}
						} else {
							if debugTableCells {
								fmt.Printf("DEBUG: storageToNodeForTable no ListStyle found\n")
							}
						}
					} else {
						if debugTableCells {
							fmt.Printf("DEBUG: storageToNodeForTable ListStyleNull is true\n")
						}
					}
				} else {
					if debugTableCells {
						fmt.Printf("DEBUG: storageToNodeForTable no ParaProperties found\n")
					}
				}
			} else {
				// Print unrecognized paragraph style type
				fmt.Printf("DEBUG: Unrecognized paragraph style type: %T\n", objRef)
			}
		} else {
			if debugTableCells {
				fmt.Printf("DEBUG: storageToNodeForTable no Object found\n")
			}
		}

		// Render paragraph content
		paragraphNode := ctx.processTableCellParagraph(rr[pos:end], e, bs, pos, end)

		if isListItem {
			// If it's a list item
			if !inList || currentListType != listType {
				// Start new list or switch to different list type
				if inList {
					// End current list
					inList = false
					currentList = nil
				}
				// Create new list
				currentList = E(listType)
				currentList.Attr = append(currentList.Attr, html.Attribute{Key: "style", Val: "margin: 0; padding-left: 20px;"})
				td.AppendChild(currentList)
				currentListType = listType
				inList = true
			}
			currentList.AppendChild(paragraphNode)
		} else {
			// If not a list item
			if inList {
				// End current list
				inList = false
				currentList = nil
				currentListType = ""
			}
			td.AppendChild(paragraphNode)
		}

		// After paragraph rendering, if there are attachments within the paragraph range, attach them to the cell (non-inline)
		for len(attachments) > 0 && attachments[0].pos < end {
			td.AppendChild(attachments[0].node)
			attachments = attachments[1:]
		}
	}

	return nil
}

// processTableCellParagraph processes paragraph content in table cells, maintaining format but simplifying structure
func (ctx *Context) processTableCellParagraph(text []rune, paraStyle *TSWP.ObjectAttributeTable_ObjectAttribute, bs *TSWP.StorageArchive, globalStart uint32, globalEnd uint32) *html.Node {
	// Get paragraph style
	var psa *TSWP.ParagraphStyleArchive
	if paraStyle.Object != nil {
		objRef := ctx.ix.Deref(paraStyle.Object)
		if psaRef, ok := objRef.(*TSWP.ParagraphStyleArchive); ok {
			psa = psaRef
		} else {
			// Print unrecognized paragraph style type
			fmt.Printf("DEBUG: processTableCellParagraph unrecognized paragraph style type: %T\n", objRef)
		}
	}

	// Check if it's a list item - consistent with normal paragraph processing
	var isListItem bool
	var listType string
	if psa != nil && psa.ParaProperties != nil {
		// Check ListStyleNull field
		if psa.ParaProperties.ListStyleNull == nil || !*psa.ParaProperties.ListStyleNull {
			// ListStyle is not null, check if there's a ListStyle
			if psa.ParaProperties.ListStyle != nil {
				// Check if ListStyle has label types, consistent with normal paragraph processing
				if ls, ok := ctx.ix.Deref(psa.ParaProperties.ListStyle).(*TSWP.ListStyleArchive); ok {
					if len(ls.LabelTypes) > 0 {
						isListItem = true
						labelType := ls.LabelTypes[0]
						switch labelType {
						case TSWP.ListStyleArchive_kNumber:
							listType = "ol"
						case TSWP.ListStyleArchive_kString, TSWP.ListStyleArchive_kImage:
							listType = "ul"
						default:
							listType = "ul"
						}
						if debugTableCells {
							fmt.Printf("DEBUG: processTableCellParagraph detected list item with type: %s\n", listType)
						}
					}
				}
			}
		}
	}

	// Create container element
	var container *html.Node
	if isListItem {
		// Create list item element
		container = E("li")
	} else {
		// Create normal div element
		container = E("div")
	}

	// Apply paragraph styles to container
	if psa != nil {
		className := fmt.Sprintf("ps%d", *paraStyle.Object.Identifier)
		container.Attr = append(container.Attr, html.Attribute{Key: "class", Val: className})

		// Apply paragraph styles
		if psa.ParaProperties != nil {
			style := translateParaProps(ctx.ix, psa.ParaProperties)
			if isListItem && psa.ParaProperties.ListStyle != nil {
				// Add list styles for list items
				listCSS := translateListStyle(ctx.ix, psa.ParaProperties.ListStyle)
				if listCSS != "" {
					if style != "" {
						style += " " + listCSS
					} else {
						style = listCSS
					}
				}
			}
			if style != "" {
				container.Attr = append(container.Attr, html.Attribute{Key: "style", Val: style})
			}
		}
	}

	// Process character styles - output segment by segment within paragraph (no smart merging)
	if bs.TableCharStyle != nil {
		// Paragraph entries, mapped to local index
		var charStyles []*TSWP.ObjectAttributeTable_ObjectAttribute
		for _, entry := range bs.TableCharStyle.Entries {
			if entry.CharacterIndex == nil {
				continue
			}
			ci := *entry.CharacterIndex
			if ci < globalStart || ci >= globalEnd {
				continue
			}
			ne := *entry
			local := ci - globalStart
			ne.CharacterIndex = &local
			charStyles = append(charStyles, &ne)
		}

		// Ensure character styles are sorted by position in ascending order to avoid interval calculation errors causing duplicate text
		if len(charStyles) > 1 {
			for i := 0; i < len(charStyles)-1; i++ {
				for j := i + 1; j < len(charStyles); j++ {
					if *charStyles[i].CharacterIndex > *charStyles[j].CharacterIndex {
						charStyles[i], charStyles[j] = charStyles[j], charStyles[i]
					}
				}
			}
		}

		var pos uint32 = 0
		endLocal := uint32(len(text))
		for i, e := range charStyles {
			cs := *e.CharacterIndex
			if cs < pos {
				continue
			}
			if cs >= endLocal {
				break
			}

			ce := endLocal
			if i+1 < len(charStyles) {
				ce = *charStyles[i+1].CharacterIndex
			}
			if ce > endLocal {
				ce = endLocal
			}
			if cs > pos {
				container.AppendChild(T(string(text[pos:cs])))
				pos = cs
			}

			if e.Object != nil {
				objRef := ctx.ix.Deref(e.Object)
				if ref, ok := objRef.(*TSWP.CharacterStyleArchive); ok {
					key := fmt.Sprintf("ss%d", *e.Object.Identifier)

					if ref.Super.Parent != nil {
						parentRef := ctx.ix.Deref(ref.Super.Parent)
						if parent, ok := parentRef.(*TSWP.CharacterStyleArchive); ok {
							mergeCharProps(ref.CharProperties, parent.CharProperties)
						} else {
							// Print unrecognized parent character style type
							fmt.Printf("DEBUG: Unrecognized parent character style type: %T\n", parentRef)
						}
					}

					style := translateCharProps(ctx, ref.CharProperties)

					props := ref.CharProperties
					if props != nil && props.Bold != nil && *props.Bold &&
						(props.Italic == nil || !*props.Italic) &&
						props.FontSize == nil && props.FontName == nil {
						container.AppendChild(E("b", string(text[cs:ce])))
					} else if props != nil && props.Italic != nil && *props.Italic &&
						(props.Bold == nil || !*props.Bold) &&
						props.FontSize == nil && props.FontName == nil {
						container.AppendChild(E("em", string(text[cs:ce])))
					} else if style != "" {
						ctx.styles[key] = style
						container.AppendChild(E("span", []string{"class", key}, string(text[cs:ce])))
					} else {
						container.AppendChild(T(string(text[cs:ce])))
					}
				} else {
					// Print unrecognized character style type
					fmt.Printf("DEBUG: Unrecognized character style type: %T\n", objRef)
					container.AppendChild(T(string(text[cs:ce])))
				}
			} else {
				container.AppendChild(T(string(text[cs:ce])))
			}
			pos = ce
		}
		if pos < endLocal {
			container.AppendChild(T(string(text[pos:endLocal])))
		}
	} else {
		container.AppendChild(T(string(text)))
	}

	return container
}

// processTableCellText processes text and character styles in table cells
func (ctx *Context) processTableCellText(text []rune, paraStyle *TSWP.ParagraphStyleArchive, bs *TSWP.StorageArchive) *html.Node {
	// Create container element
	container := E("span")

	// Process character styles
	if bs.TableCharStyle != nil {
		charStyles := bs.TableCharStyle.Entries
		pos := 0

		// Improved character style processing logic
		// First sort character style entries by position
		sortedStyles := make([]struct {
			index uint32
			entry *TSWP.ObjectAttributeTable_ObjectAttribute
		}, len(charStyles))

		for i, e := range charStyles {
			sortedStyles[i] = struct {
				index uint32
				entry *TSWP.ObjectAttributeTable_ObjectAttribute
			}{*e.CharacterIndex, e}
		}

		// Sort by position
		for i := 0; i < len(sortedStyles)-1; i++ {
			for j := i + 1; j < len(sortedStyles); j++ {
				if sortedStyles[i].index > sortedStyles[j].index {
					sortedStyles[i], sortedStyles[j] = sortedStyles[j], sortedStyles[i]
				}
			}
		}

		// Process sorted styles
		for i, styleEntry := range sortedStyles {
			cs := styleEntry.index
			if cs < uint32(pos) {
				continue
			}
			if cs >= uint32(len(text)) {
				break
			}

			// Calculate end position of current style
			ce := uint32(len(text)) // Default to end of text

			// Find next style position
			for j := i + 1; j < len(sortedStyles); j++ {
				nextCs := sortedStyles[j].index
				if nextCs > cs {
					ce = nextCs
					break
				}
			}

			// Ensure ce doesn't exceed text length
			if ce > uint32(len(text)) {
				ce = uint32(len(text))
			}

			// Add previous text (part without styles)
			if cs > uint32(pos) {
				container.AppendChild(T(string(text[pos:cs])))
			}

			// Process current character style
			if styleEntry.entry.Object != nil {
				objRef := ctx.ix.Deref(styleEntry.entry.Object)
				if csa, ok := objRef.(*TSWP.CharacterStyleArchive); ok {
					span := E("span")

					// Apply character styles
					if csa.CharProperties != nil {
						style := translateCharProps(ctx, csa.CharProperties)
						if style != "" {
							span.Attr = append(span.Attr, html.Attribute{Key: "style", Val: style})
						}
					}

					// Add text content
					span.AppendChild(T(string(text[cs:ce])))
					container.AppendChild(span)
				} else {
					// Print unrecognized character style type
					fmt.Printf("DEBUG: Unrecognized character style type: %T\n", objRef)
					// If no character styles, directly add text
					container.AppendChild(T(string(text[cs:ce])))
				}
			} else {
				// If no character styles, directly add text
				container.AppendChild(T(string(text[cs:ce])))
			}

			pos = int(ce)
		}

		// Add remaining text
		if pos < len(text) {
			container.AppendChild(T(string(text[pos:])))
		}
	} else {
		// No character styles, directly add text
		container.AppendChild(T(string(text)))
	}

	return container
}

func (ctx *Context) storageToNode(bs *TSWP.StorageArchive, body *html.Node) error {
	ix := ctx.ix
	texts := bs.Text

	if len(texts) == 0 {
		return fmt.Errorf("no text content found")
	}

	// Process multiple text fragments
	var text string
	if len(texts) == 1 {
		text = texts[0]
	} else {
		// Merge multiple text fragments
		for _, t := range texts {
			text += t
		}
	}

	// Offsets are in terms of unicode runes, so we have to convert to runes
	rr := []rune(text)

	// <p>
	parStyles := bs.TableParaStyle.Entries

	var attachments []Attachment
	if bs.TableAttachment != nil {
		for _, entry := range bs.TableAttachment.Entries {
			pos := *entry.CharacterIndex
			switch ctx.ix.Deref(entry.Object).(type) {
			case *TSWP.DrawableAttachmentArchive:
				archive := ctx.ix.Deref(entry.Object).(*TSWP.DrawableAttachmentArchive)
				node := ctx.processDrawable(archive.Drawable)
				if node != nil {
					attachments = append(attachments, Attachment{pos, node})
				}
			case *TSWP.NumberAttachmentArchive:
				archive := ctx.ix.Deref(entry.Object).(*TSWP.NumberAttachmentArchive)
				node := ctx.processNumberAttachment(archive)
				if node != nil {
					attachments = append(attachments, Attachment{pos, node})
				}
			}
		}
	}
	// bs.TableListStyle - seems to change on headings, look into it.

	// A null style seems to imply "use the previous class," so this is declared outside the loop.
	var className string

	// List state tracking
	var currentList *html.Node
	var currentListType string

	// build paragraphs
	for i, e := range parStyles {

		pos := *e.CharacterIndex
		end := uint32(len(rr))
		if i+1 < len(parStyles) {
			end = *parStyles[i+1].CharacterIndex
		}

		for len(attachments) > 0 && attachments[0].pos < end {
			// Handle attachment position: if attachment is not at paragraph start, insert placeholder at paragraph start
			if attachments[0].pos != pos {
				// Insert text between paragraph start and attachment position
				if attachments[0].pos > pos {
					body.AppendChild(T(string(rr[pos:attachments[0].pos])))
				}
			}
			body.AppendChild(attachments[0].node)
			attachments = attachments[1:]
		}

		tag := "p"

		// Get style, change tag if appropriate.
		if e.Object != nil {
			ref := ix.Deref(e.Object).(*TSWP.ParagraphStyleArchive)
			className = fmt.Sprintf("ps%d", *e.Object.Identifier)

			// Pages: insert a page-break marker before this paragraph when requested by style
			if ctx.ix.Type == "pages" && ref.ParaProperties != nil && ref.ParaProperties.GetPageBreakBefore() {
				body.AppendChild(E("hr", []string{"class", "page-break"}))
			}

			// Some properties are inherited (e.g. if you apply a style and then tweak it.)
			// We can't just include both because FirstLineIndent in parent can combine with LeftIndent in child
			// to produce a css text-indent.
			if ref.Super.Parent != nil {
				parent := ix.Deref(ref.Super.Parent).(*TSWP.ParagraphStyleArchive)
				mergeCharProps(ref.CharProperties, parent.CharProperties)
				mergeParaProps(ref.ParaProperties, parent.ParaProperties)

				// Recursively process parent style inheritance
				ctx.mergeParentStyles(ref, parent)
			}

			style := translateParaProps(ctx.ix, ref.ParaProperties) + translateCharProps(ctx, ref.CharProperties)
			if debugMode {
				fmt.Printf("DEBUG: Style for class %s: %s\n", className, style)
			}
			ctx.styles[className] = style

			if ref.ParaProperties.OutlineLevel != nil {
				level := *ref.ParaProperties.OutlineLevel
				if level < 7 {
					tag = fmt.Sprintf("h%d", level)
				}
			}
		}

		// Check if it's a list item
		var isListItem bool
		var listType string
		if e.Object != nil {
			ref := ix.Deref(e.Object).(*TSWP.ParagraphStyleArchive)
			if ref.ParaProperties != nil && (ref.ParaProperties.ListStyleNull == nil || !*ref.ParaProperties.ListStyleNull) && ref.ParaProperties.ListStyle != nil {
				// Only treat as list item when ListStyle has label types; otherwise keep as normal paragraph
				if ls, ok := ix.Deref(ref.ParaProperties.ListStyle).(*TSWP.ListStyleArchive); ok {
					if len(ls.LabelTypes) > 0 {
						isListItem = true
						labelType := ls.LabelTypes[0]
						switch labelType {
						case TSWP.ListStyleArchive_kNumber:
							listType = "ol"
						case TSWP.ListStyleArchive_kString, TSWP.ListStyleArchive_kImage:
							listType = "ul"
						default:
							listType = "ul"
						}
					}
				}
			}
		}

		// Skip empty list items: if paragraph is list-styled but content is empty, do not render
		if isListItem {
			if strings.TrimSpace(string(rr[pos:end])) == "" {
				continue
			}
		}

		// Process list items
		var p *html.Node
		if isListItem {
			// Check if need to create new list container
			if currentList == nil || currentListType != listType {
				// Close current list (if any)
				if currentList != nil {
					body.AppendChild(currentList)
				}

				// Create new list
				currentList = E(listType)
				currentListType = listType
			}

			// Create list item
			p = E("li", []string{"class", className})
			// Apply list styles
			if e.Object != nil {
				ref := ix.Deref(e.Object).(*TSWP.ParagraphStyleArchive)
				if ref.ParaProperties != nil && ref.ParaProperties.ListStyle != nil {
					listCSS := translateListStyle(ix, ref.ParaProperties.ListStyle)
					if listCSS != "" {
						p.Attr = append(p.Attr, html.Attribute{Key: "style", Val: listCSS})
					}
				}
			}
		} else {
			// Non-list item, close current list (if any)
			if currentList != nil {
				body.AppendChild(currentList)
				currentList = nil
				currentListType = ""
			}
			p = E(tag, []string{"class", className})
		}

		// <span> <em> and <b> - paragraph character style processing
		if bs.TableCharStyle != nil {
			charStyles := bs.TableCharStyle.Entries

			// Only process character styles within current paragraph range
			for i, e := range charStyles {
				cs := *e.CharacterIndex
				if cs < pos {
					continue
				}
				if cs >= end {
					break
				}

				// Calculate end position of current style
				ce := uint32(len(rr))
				if i+1 < len(charStyles) {
					ce = *charStyles[i+1].CharacterIndex
				}
				// Limit to current paragraph range
				if ce > end {
					ce = end
				}

				// Add text before style
				if cs > pos {
					p.AppendChild(T(string(rr[pos:cs])))
					pos = cs
				}

				// Apply character styles
				if e.Object != nil {
					ref := ix.Deref(e.Object).(*TSWP.CharacterStyleArchive)
					key := fmt.Sprintf("ss%d", *e.Object.Identifier)

					if ref.Super.Parent != nil {
						parent := ix.Deref(ref.Super.Parent).(*TSWP.CharacterStyleArchive)
						mergeCharProps(ref.CharProperties, parent.CharProperties)
						// Recursively process parent style inheritance
						ctx.mergeParentCharStyles(ref, parent)
					}

					style := translateCharProps(ctx, ref.CharProperties)

					// Check specific style properties to decide which tag to use
					props := ref.CharProperties
					if props != nil {
						// Check if only bold
						if props.Bold != nil && *props.Bold &&
							(props.Italic == nil || !*props.Italic) &&
							props.FontSize == nil && props.FontName == nil {
							p.AppendChild(E("b", string(rr[cs:ce])))
						} else if props.Italic != nil && *props.Italic &&
							(props.Bold == nil || !*props.Bold) &&
							props.FontSize == nil && props.FontName == nil {
							p.AppendChild(E("em", string(rr[cs:ce])))
						} else if style != "" {
							// Has complex styles, use span
							ctx.styles[key] = style
							p.AppendChild(E("span", []string{"class", key}, string(rr[cs:ce])))
						} else {
							// No styles, directly add text
							p.AppendChild(T(string(rr[cs:ce])))
						}
					} else {
						// No style properties, directly add text
						p.AppendChild(T(string(rr[cs:ce])))
					}
				} else {
					p.AppendChild(T(string(rr[cs:ce])))
				}
				pos = ce
			}
		}

		// Add paragraph content, but skip empty paragraphs
		paragraphContent := string(rr[pos:end])
		if strings.TrimSpace(paragraphContent) != "" {
			p.AppendChild(T(paragraphContent))
		}

		// Process list items
		if isListItem {
			// Add list item to current list
			currentList.AppendChild(p)
		} else {
			// Add normal paragraph
			body.AppendChild(p)
		}
		body.AppendChild(T("\n"))
	}

	// Process any remaining lists
	if currentList != nil {
		body.AppendChild(currentList)
	}

	return nil
}

type Style map[string]interface{}

func Convert(in, out string) error {
	fmt.Println("Processing", in)

	var err error
	var ctx Context
	ctx.styles = make(map[string]string)
	ctx.imgs = make(map[string]uint64)
	if ctx.ix, err = index.Open(in); err != nil {
		return err
	}
	if ctx.zr, err = zip.OpenReader(in); err != nil {
		return err
	}
	defer ctx.zr.Close()

	if debugMode {
		fmt.Println("Read", len(ctx.ix.Records), "records")
	}

	var doc *html.Node
	// set global for font scaling hook
	switch ctx.ix.Type {
	case "pages":
		doc = ctx.processPages()
	case "numbers":
		doc = ctx.processNumbers()
	case "key":
		doc = ctx.processKeynote()
	}

	w, err := os.Create(out)
	if err != nil {
		return err
	}
	defer w.Close()
	fmt.Println("Writing", out)
	if strings.HasSuffix(out, ".json") {
		out, err := json.MarshalIndent(ctx.ix, "", "  ")
		if err != nil {
			return err
		}
		if _, err = w.Write(out); err != nil {
			return err
		}
	} else {
		// Ensure HTML output uses UTF-8 encoding
		html.Render(w, doc)
	}
	return nil
}

func (ctx *Context) renderImgData() *html.Node {
	if len(ctx.imgs) <= 0 {
		return nil
	}

	s := ""
	for _, f := range ctx.zr.File {
		if id, ok := ctx.imgs[f.Name]; ok {
			rc, _ := f.Open()
			defer rc.Close()

			imageBytes, _ := io.ReadAll(rc)
			s += "document.querySelectorAll('.img_" + fmt.Sprint(id) + "').forEach(function(e) {e.src='data:image/" + filepath.Ext(f.Name)[1:] + ";base64," + base64.StdEncoding.EncodeToString(imageBytes) + "';});"
			s += "document.querySelectorAll('.bgimg_" + fmt.Sprint(id) + "').forEach(function(e) {e.style.backgroundImage='url(data:image/" + filepath.Ext(f.Name)[1:] + ";base64," + base64.StdEncoding.EncodeToString(imageBytes) + ")'; e.style.backgroundSize='cover'; e.style.backgroundPosition='center';});"
		}
	}

	script := E("script")
	script.AppendChild(T(s))
	return script
}

// processPages translates a pages file.
func (ctx *Context) processPages() *html.Node {
	// Root of output document
	head, body := E("head", "\n", E("meta", []string{"charset", "utf-8"}), "\n"), E("body", "\n")
	doc := E("", E("html", head, "\n", body))
	doc.Type = html.DocumentNode

	da := ctx.ix.Records[1].(*TP.DocumentArchive)
	bs := ctx.ix.Deref(da.BodyStorage).(*TSWP.StorageArchive)

	fda := ctx.ix.Deref(da.FloatingDrawables).(*TP.FloatingDrawablesArchive)
	if len(fda.PageGroups) != 0 {
		fmt.Println(`WARNING - 
            This document has floating drawables (e.g. floating images/tables/text blocks) which we don't handle in HTML
            conversion.
            
            Figuring out where to place them in the document would probably be tricky.
            
`)
	}

	// Render into a temporary container to post-process page breaks
	temp := E("div")
	// Pages font units are consistent with canvas units, if page size exists, font scaling can be set here
	// For Pages documents, don't use font scaling
	ctx.fontScale = 1.0 // Don't use font scaling, maintain original font size
	fmt.Printf("*** Setting font scale factor: %.2f\n", ctx.fontScale)
	ctx.storageToNode(bs, temp)

	// Split content by hr.page-break into pages
	container := E("div", []string{"class", "page-container"})
	page := E("div", []string{"class", "page"})
	for n := temp.FirstChild; n != nil; {
		next := n.NextSibling
		isBreak := false
		if n.Type == html.ElementNode && n.Data == "hr" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "page-break" {
					isBreak = true
					break
				}
			}
		}
		// Detach from temp before reparenting/moving
		temp.RemoveChild(n)
		if isBreak {
			container.AppendChild(page)
			page = E("div", []string{"class", "page"})
		} else {
			page.AppendChild(n)
		}
		n = next
	}
	container.AppendChild(page)

	// Place floating drawables per page
	if da != nil {
		if fdaRef := da.FloatingDrawables; fdaRef != nil {
			if fda, ok := ctx.ix.Deref(fdaRef).(*TP.FloatingDrawablesArchive); ok {
				// Ensure there are enough page nodes for all indices
				maxIdx := -1
				for _, pg := range fda.GetPageGroups() {
					if int(pg.GetPageIndex()) > maxIdx {
						maxIdx = int(pg.GetPageIndex())
					}
				}
				// Count existing pages
				existing := 0
				for n := container.FirstChild; n != nil; n = n.NextSibling {
					if n.Type == html.ElementNode && n.Data == "div" {
						existing++
					}
				}
				for existing <= maxIdx {
					container.AppendChild(E("div", []string{"class", "page"}))
					existing++
				}

				// Build slice of page nodes for index lookup
				pages := []*html.Node{}
				for n := container.FirstChild; n != nil; n = n.NextSibling {
					if n.Type == html.ElementNode && n.Data == "div" {
						pages = append(pages, n)
					}
				}
				for _, pg := range fda.GetPageGroups() {
					idx := int(pg.GetPageIndex())
					if idx >= 0 && idx < len(pages) {
						host := pages[idx]
						// add drawables
						add := func(entries []*TP.FloatingDrawablesArchive_DrawableEntry) {
							for _, de := range entries {
								if de == nil || de.GetDrawable() == nil {
									continue
								}
								if node := ctx.processDrawable(de.GetDrawable()); node != nil {
									host.AppendChild(node)
								}
							}
						}
						add(pg.GetBackgroundDrawables())
						add(pg.GetDrawables())
						add(pg.GetForegroundDrawables())
					}
				}
			}
		}
	}

	body.AppendChild(container)

	paginate := E("script")
	paginate.AppendChild(T(
		"document.addEventListener('DOMContentLoaded', function() {\n" +
			"  var container = document.querySelector('.page-container');\n" +
			"  if (!container) return;\n" +
			"  var pages = Array.from(container.querySelectorAll('.page'));\n" +
			"  if (pages.length === 0) return;\n" +
			"  var firstPage = pages[0];\n" +
			"  function getUsableHeight(el){\n" +
			"    var cs = getComputedStyle(el);\n" +
			"    var pt = parseFloat(cs.paddingTop)||0; var pb = parseFloat(cs.paddingBottom)||0;\n" +
			"    return el.clientHeight - pt - pb;\n" +
			"  }\n" +
			"  function getTableHeight(table) {\n" +
			"    var tempDiv = document.createElement('div');\n" +
			"    tempDiv.style.position = 'absolute';\n" +
			"    tempDiv.style.visibility = 'hidden';\n" +
			"    tempDiv.style.width = table.offsetWidth + 'px';\n" +
			"    tempDiv.style.top = '-9999px';\n" +
			"    tempDiv.style.left = '-9999px';\n" +
			"    tempDiv.appendChild(table.cloneNode(true));\n" +
			"    document.body.appendChild(tempDiv);\n" +
			"    var height = tempDiv.offsetHeight;\n" +
			"    document.body.removeChild(tempDiv);\n" +
			"    return height;\n" +
			"  }\n" +
			"  \n" +
			"  function getTableRowHeight(row) {\n" +
			"    var tempDiv = document.createElement('div');\n" +
			"    tempDiv.style.position = 'absolute';\n" +
			"    tempDiv.style.visibility = 'hidden';\n" +
			"    tempDiv.style.width = '100%';\n" +
			"    tempDiv.style.top = '-9999px';\n" +
			"    tempDiv.style.left = '-9999px';\n" +
			"    tempDiv.style.borderCollapse = 'collapse';\n" +
			"    tempDiv.style.tableLayout = 'fixed';\n" +
			"    \n" +
			"    \n" +
			"    var tempTable = document.createElement('table');\n" +
			"    tempTable.style.width = '100%';\n" +
			"    tempTable.style.borderCollapse = 'collapse';\n" +
			"    tempTable.style.tableLayout = 'fixed';\n" +
			"    \n" +
			"    \n" +
			"    var colgroup = document.querySelector('colgroup');\n" +
			"    if (colgroup) {\n" +
			"      tempTable.appendChild(colgroup.cloneNode(true));\n" +
			"    }\n" +
			"    \n" +
			"    tempTable.appendChild(row.cloneNode(true));\n" +
			"    tempDiv.appendChild(tempTable);\n" +
			"    document.body.appendChild(tempDiv);\n" +
			"    \n" +
			"    var height = tempTable.offsetHeight;\n" +
			"    document.body.removeChild(tempDiv);\n" +
			"    return height;\n" +
			"  }\n" +
			"  \n" +
			"  function getTableRowsHeight(rows, maxRows, colgroup, thead) {\n" +
			"    if (rows.length === 0) return 0;\n" +
			"    \n" +
			"    var actualRows = rows.slice(0, Math.min(maxRows, rows.length));\n" +
			"    \n" +
			"    \n" +
			"    var tempDiv = document.createElement('div');\n" +
			"    tempDiv.style.position = 'absolute';\n" +
			"    tempDiv.style.visibility = 'hidden';\n" +
			"    tempDiv.style.width = '100%';\n" +
			"    tempDiv.style.top = '-9999px';\n" +
			"    tempDiv.style.left = '-9999px';\n" +
			"    \n" +
			"    var tempTable = document.createElement('table');\n" +
			"    tempTable.style.width = '100%';\n" +
			"    tempTable.style.borderCollapse = 'collapse';\n" +
			"    tempTable.style.tableLayout = 'fixed';\n" +
			"    \n" +
			"    \n" +
			"    if (colgroup) {\n" +
			"      tempTable.appendChild(colgroup.cloneNode(true));\n" +
			"    }\n" +
			"    \n" +
			"    \n" +
			"    if (thead) {\n" +
			"      tempTable.appendChild(thead.cloneNode(true));\n" +
			"    }\n" +
			"    \n" +
			"    \n" +
			"    var tbody = document.createElement('tbody');\n" +
			"    actualRows.forEach(function(row) {\n" +
			"      tbody.appendChild(row.cloneNode(true));\n" +
			"    });\n" +
			"    tempTable.appendChild(tbody);\n" +
			"    \n" +
			"    tempDiv.appendChild(tempTable);\n" +
			"    document.body.appendChild(tempDiv);\n" +
			"    \n" +
			"    var height = tempTable.offsetHeight;\n" +
			"    document.body.removeChild(tempDiv);\n" +
			"    return height;\n" +
			"  }\n" +
			"  \n" +
			"  function calculateTablePageBreak(rows, availableHeight, colgroup, thead) {\n" +
			"    if (rows.length === 0) return 0;\n" +
			"    \n" +
			"    var maxRows = rows.length;\n" +
			"    var minRows = 1;\n" +
			"    var bestFit = 0;\n" +
			"    \n" +
			"    \n" +
			"    while (minRows <= maxRows) {\n" +
			"      var mid = Math.floor((minRows + maxRows) / 2);\n" +
			"      var height = getTableRowsHeight(rows, mid, colgroup, thead);\n" +
			"      \n" +
			"      if (height <= availableHeight) {\n" +
			"        bestFit = mid;\n" +
			"        minRows = mid + 1;\n" +
			"      } else {\n" +
			"        maxRows = mid - 1;\n" +
			"      }\n" +
			"    }\n" +
			"    \n" +
			"    \n" +
			"    return Math.max(1, bestFit);\n" +
			"  }\n" +
			"  function isElementOverflowing(element) {\n" +
			"    return element.scrollHeight > element.clientHeight;\n" +
			"  }\n" +
			"  function isLikelyTextBox(n){\n" +
			"    if (!(n instanceof HTMLElement)) return false;\n" +
			"    var hasNonText = n.querySelector('table,img,canvas,svg,video,figure') != null;\n" +
			"    if (hasNonText) return false;\n" +
			"    var textLen = (n.textContent||'').trim().length;\n" +
			"    if (textLen === 0) return false;\n" +
			"    return true;\n" +
			"  }\n" +
			"  var flow = [];\n" +
			"  pages.forEach(function(page){\n" +
			"    Array.from(page.childNodes).forEach(function(n){\n" +
			"      if (!(n instanceof HTMLElement)) return;\n" +
			"      if (n.tagName === 'HR' && n.classList.contains('page-break')) { flow.push('FORCE_BREAK'); page.removeChild(n); return; }\n" +
			"      var pos = getComputedStyle(n).position;\n" +
			"      if (pos === 'absolute' || pos === 'fixed') {\n" +
			"        if (!isLikelyTextBox(n)) { return; }\n" +
			"        \n" +
			"        n.dataset.flowText = '1';\n" +
			"        n.style.position = 'static';\n" +
			"        n.style.left = ''; n.style.top = '';\n" +
			"        n.style.width = '100%';\n" +
			"        n.style.whiteSpace = 'normal';\n" +
			"        n.style.display = 'block';\n" +
			"      }\n" +
			"      flow.push(n); page.removeChild(n);\n" +
			"    });\n" +
			"  });\n" +
			"  for (var i=pages.length-1;i>=1;i--) { container.removeChild(pages[i]); }\n" +
			"  pages = [firstPage];\n" +
			"  function newPage(){ var p = document.createElement('div'); p.className='page'; container.appendChild(p); pages.push(p); return p; }\n" +
			"  var current = firstPage;\n" +
			"  var limit = getUsableHeight(current);\n" +
			"  flow.forEach(function(node){\n" +
			"    if (node === 'FORCE_BREAK'){ current = newPage(); limit = getUsableHeight(current); return; }\n" +
			"    \n" +
			"    \n" +
			"    var table = null;\n" +
			"    if (node.tagName === 'TABLE') {\n" +
			"      table = node;\n" +
			"    } else if (node.querySelector && node.querySelector('table')) {\n" +
			"      table = node.querySelector('table');\n" +
			"    }\n" +
			"    \n" +
			"    if (table) {\n" +
			"      \n" +
			"      var original = node;\n" +
			"      var colgroup = table.querySelector('colgroup');\n" +
			"      var thead = table.querySelector('thead');\n" +
			"      var tbody = table.querySelector('tbody') || table;\n" +
			"      var rows = Array.from(tbody.querySelectorAll('tr'));\n" +
			"      \n" +
			"      console.log('Processing table, row count:', rows.length, 'table node:', table);\n" +
			"      \n" +
			"      \n" +
			"      var headerRows = [];\n" +
			"      if (thead) { \n" +
			"        headerRows = Array.from(thead.querySelectorAll('tr'));\n" +
			"        console.log('Found header row count:', headerRows.length);\n" +
			"      }\n" +
			"      \n" +
			"      if (headerRows.length === 0 && rows.length > 0) {\n" +
			"        var firstRow = rows[0];\n" +
			"        \n" +
			"        // Check if it's a th element or if it's a single-column table with centered text\n" +
			"        var hasThElement = firstRow.querySelector('th');\n" +
			"        var isSingleColumn = table.querySelectorAll('colgroup col').length === 1;\n" +
			"        var hasHeaderStyle = false;\n" +
			"        \n" +
			"        // Check for header-like styling (centered text, specific classes)\n" +
			"        if (isSingleColumn && firstRow.cells.length > 0) {\n" +
			"          var firstCell = firstRow.cells[0];\n" +
			"          var hasCenterAlign = firstCell.innerHTML.includes('text-align: center');\n" +
			"          hasHeaderStyle = hasCenterAlign;\n" +
			"        }\n" +
			"        \n" +
			"        if (hasThElement || hasHeaderStyle) {\n" +
			"          headerRows = [firstRow];\n" +
			"          rows = rows.slice(1);\n" +
			"          console.log('Identifying first row as header row (th=' + hasThElement + ', style=' + hasHeaderStyle + ')');\n" +
			"        }\n" +
			"      }\n" +
			"      \n" +
			"      \n" +
			"      function createTableShell(){\n" +
			"        var t = document.createElement('table');\n" +
			"        t.className = 'table-paginated';\n" +
			"        if (colgroup) t.appendChild(colgroup.cloneNode(true));\n" +
			"        var thd = document.createElement('thead');\n" +
			"        if (headerRows.length > 0) { headerRows.forEach(function(hr){ thd.appendChild(hr.cloneNode(true)); }); }\n" +
			"        t.appendChild(thd);\n" +
			"        var tbd = document.createElement('tbody');\n" +
			"        t.appendChild(tbd);\n" +
			"        return {table:t, body:tbd};\n" +
			"      }\n" +
			"      \n" +
			"      \n" +
			"      if (rows.length <= 2) {\n" +
			"        console.log('Small table processing, rows:', rows.length);\n" +
			"        \n" +
			"        var hasContent = current.children.length > 0;\n" +
			"        if (hasContent) {\n" +
			"          var tempDiv = document.createElement('div');\n" +
			"          tempDiv.appendChild(original.cloneNode(true));\n" +
			"          current.appendChild(tempDiv);\n" +
			"          if (isElementOverflowing(current)) {\n" +
			"            current.removeChild(tempDiv);\n" +
			"            current = newPage();\n" +
			"            current.appendChild(original);\n" +
			"          } else {\n" +
			"            current.removeChild(tempDiv);\n" +
			"            current.appendChild(original);\n" +
			"          }\n" +
			"        } else {\n" +
			"          current.appendChild(original);\n" +
			"        }\n" +
			"        return;\n" +
			"      }\n" +
			"      \n" +
			"      \n" +
			"      console.log('Large table smart pagination processing, row count:', rows.length);\n" +
			"      \n" +
			"      var remainingRows = rows.slice();\n" +
			"      var currentPage = current;\n" +
			"      \n" +
			"      while (remainingRows.length > 0) {\n" +
			"        \n" +
			"        var availableHeight = getUsableHeight(currentPage);\n" +
			"        \n" +
			"        \n" +
			"        var rowsToFit = calculateTablePageBreak(remainingRows, availableHeight, colgroup, thead);\n" +
			"        \n" +
			"        console.log('Current page available height:', availableHeight, 'px, can fit rows:', rowsToFit);\n" +
			"        \n" +
			"        \n" +
			"        var part = createTableShell();\n" +
			"        currentPage.appendChild(part.table);\n" +
			"        \n" +
			"        \n" +
			"        for (var i = 0; i < rowsToFit && i < remainingRows.length; i++) {\n" +
			"          part.body.appendChild(remainingRows[i].cloneNode(true));\n" +
			"        }\n" +
			"        \n" +
			"        \n" +
			"        if (isElementOverflowing(currentPage)) {\n" +
			"          console.log('Validation failed, reducing row count');\n" +
			"          \n" +
			"          while (isElementOverflowing(currentPage) && part.body.children.length > 0) {\n" +
			"            part.body.removeChild(part.body.lastChild);\n" +
			"          }\n" +
			"          \n" +
			"          \n" +
			"          if (part.body.children.length === 0 && remainingRows.length > 0) {\n" +
			"            part.body.appendChild(remainingRows[0].cloneNode(true));\n" +
			"            remainingRows = remainingRows.slice(1);\n" +
			"          } else {\n" +
			"            \n" +
			"            var actualPlaced = part.body.children.length;\n" +
			"            remainingRows = remainingRows.slice(actualPlaced);\n" +
			"          }\n" +
			"        } else {\n" +
			"          \n" +
			"          remainingRows = remainingRows.slice(rowsToFit);\n" +
			"        }\n" +
			"        \n" +
			"        \n" +
			"        if (remainingRows.length > 0) {\n" +
			"          console.log('Still', remainingRows.length, 'rows to process, creating new page');\n" +
			"          currentPage = newPage();\n" +
			"        }\n" +
			"      }\n" +
			"      \n" +
			"      console.log('Table smart pagination completed, total pages:', pages.length);\n" +
			"    } else {\n" +
			"      \n" +
			"      if (node.dataset && node.dataset.flowText === '1') {\n" +
			"        var children = Array.from(node.childNodes).filter(function(c){ return c instanceof HTMLElement; });\n" +
			"        var shell = document.createElement('div');\n" +
			"        shell.style.width = '100%'; shell.style.whiteSpace = 'normal'; shell.style.display = 'block';\n" +
			"        current.appendChild(shell);\n" +
			"        children.forEach(function(c){\n" +
			"          shell.appendChild(c);\n" +
			"          if (isElementOverflowing(current)) {\n" +
			"            shell.removeChild(c);\n" +
			"            current = newPage();\n" +
			"            shell = document.createElement('div');\n" +
			"            shell.style.width = '100%'; shell.style.whiteSpace = 'normal'; shell.style.display = 'block';\n" +
			"            current.appendChild(shell);\n" +
			"            shell.appendChild(c);\n" +
			"          }\n" +
			"        });\n" +
			"      } else {\n" +
			"        current.appendChild(node);\n" +
			"        if (isElementOverflowing(current)) {\n" +
			"          current.removeChild(node);\n" +
			"          current = newPage();\n" +
			"          current.appendChild(node);\n" +
			"        }\n" +
			"      }\n" +
			"    }\n" +
			"  });\n" +
			"});\n"))
	body.AppendChild(paginate)

	if img := ctx.renderImgData(); img != nil {
		body.AppendChild(img)
	}

	style := E("style")
	style.AppendChild(T(
		"html, body { height: 100%; background: transparent; }\n" +
			"body { margin: 0; font-family: -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,'Noto Sans','Liberation Sans',sans-serif; color: #111; }\n" +
			".page-container { display: flex; flex-direction: column; align-items: center; gap: 20px; padding: 20px 12px 40px; box-sizing: border-box; }\n" +
			".page { position: relative; width: min(900px, calc(100vw - 48px)); aspect-ratio: 210 / 297; background: #fff; overflow: hidden; box-shadow: 0 10px 30px rgba(0,0,0,0.25); box-sizing: border-box; padding: 20px; }\n" +
			".page * { box-sizing: border-box; }\n" +
			"p { margin: 0; line-height: 1.2; }\n" +
			"/* Table Styles */\n" +
			".page table { width: 100%; margin: 16px auto; border-collapse: collapse; border: 1px solid #ddd; }\n" +
			".page td, .page th { padding: 8px 12px; vertical-align: middle; border: 1px solid #ddd; }\n" +
			".page tbody tr:hover { background-color: #fafafa; }\n" +
			"/* Table pagination optimization styles */\n" +
			".page table.table-paginated { page-break-inside: auto; }\n" +
			".page table.table-paginated thead { display: table-header-group; }\n" +
			".page table.table-paginated tbody { display: table-row-group; }\n" +
			"/* Table pagination indicator */\n" +
			".page .table-page-indicator { \n" +
			"  font-size: 0.8em; color: #666; text-align: center; \n" +
			"  margin: 0.5em 0; padding: 0.25em; \n" +
			"  border-top: 1px dashed #ccc; \n" +
			"}\n"))
	for k, v := range ctx.styles {
		style.AppendChild(T(fmt.Sprintf(".%s {\n%s}\n", k, v)))
	}
	head.AppendChild(style)
	return doc
}

func (ctx *Context) processNumbers() *html.Node {
	if debugTableCells {
		fmt.Printf("DEBUG: processNumbers called, debugTableCells=%v\n", debugTableCells)
	}
	// Root of output document
	head, body := E("head", "\n", E("meta", []string{"charset", "utf-8"}), "\n"), E("body", "\n")
	doc := E("", E("html"), "\n", E("html", head, "\n", body))
	doc.Type = html.DocumentNode
	doc.FirstChild.Type = html.DoctypeNode
	da := ctx.ix.Records[1].(*TN.DocumentArchive)

	for _, ref := range da.Sheets {
		sheet := ctx.ix.Deref(ref).(*TN.SheetArchive)
		section := E("section", E("h2", "Sheet - ", *sheet.Name))
		body.AppendChild(section)
		if debugTableCells {
			fmt.Printf("DEBUG: Processing Numbers sheet with %d drawables\n", len(sheet.DrawableInfos))
		}
		for i, ref := range sheet.DrawableInfos {
			if debugTableCells {
				fmt.Printf("DEBUG: Processing drawable %d\n", i)
			}
			// if this cast throws there are other kinds of drawables...
			e := ctx.processDrawable(ref)
			if e != nil {
				section.AppendChild(e)
			}
		}
	}

	if img := ctx.renderImgData(); img != nil {
		body.AppendChild(img)
	}

	style := E("style")
	style.AppendChild(T("\np { margin: 0; }\n")) // reset paragraphs
	for k, v := range ctx.styles {
		style.AppendChild(T(fmt.Sprintf(".%s {\n%s}\n", k, v)))
	}
	head.AppendChild(style)

	return doc
}

// processKeynote translates a keynote file.
func (ctx *Context) processKeynote() *html.Node {
	if debugMode {
		fmt.Printf("DEBUG: processKeynote called with %d records\n", len(ctx.ix.Records))
	}

	// Root of output document
	head, body := E("head", "\n", E("meta", []string{"charset", "utf-8"}), "\n"), E("body", "\n")
	doc := E("", E("html"), "\n", E("html", head, "\n", body))
	doc.Type = html.DocumentNode
	doc.FirstChild.Type = html.DoctypeNode

	// First try to use SlideTree.Slides order
	var slideTreeSlides []*TSP.Reference
	for _, rec := range ctx.ix.Records {
		if sh, ok := rec.(*KN.ShowArchive); ok {
			if sh.SlideTree != nil && sh.SlideTree.Slides != nil {
				slideTreeSlides = sh.SlideTree.Slides
				break
			}
		}
	}

	ids := []uint64{}
	if slideTreeSlides != nil {
		// Use SlideTree.Slides order and check IsSkipped field
		if debugMode {
			fmt.Printf("DEBUG: Using SlideTree.Slides order with %d slides\n", len(slideTreeSlides))
		}
		for _, slideRef := range slideTreeSlides {
			if slideRef != nil && slideRef.Identifier != nil {
				slideNodeId := *slideRef.Identifier
				if slideNode, ok := ctx.ix.Records[slideNodeId].(*KN.SlideNodeArchive); ok {
					// Check if slide is skipped
					if slideNode.Slide != nil && slideNode.Slide.Identifier != nil {
						actualSlideId := *slideNode.Slide.Identifier
						ids = append(ids, actualSlideId)
						if debugMode {
							fmt.Printf("DEBUG: Added slide %d to display list\n", actualSlideId)
						}
					}
				}
			}
		}
	} else {
		// Fallback to original method
		if debugMode {
			fmt.Printf("DEBUG: Using fallback method\n")
		}
		for key, rec := range ctx.ix.Records {
			if _, ok := rec.(*KN.SlideArchive); ok {
				ids = append(ids, key)
				if debugMode {
					fmt.Printf("DEBUG: Found slide at record %d\n", key)
				}
			}
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	}
	if debugMode {
		fmt.Printf("DEBUG: Final slide count for display: %d\n", len(ids))
	}

	var canvasW, canvasH float64
	foundSize := false
	showArchiveCount := 0
	for _, rec := range ctx.ix.Records {
		if sh, ok := rec.(*KN.ShowArchive); ok {
			showArchiveCount++
			if sh.Size != nil && sh.Size.Width != nil && sh.Size.Height != nil {
				canvasW = float64(*sh.Size.Width)
				canvasH = float64(*sh.Size.Height)
				foundSize = true
				break
			}
		}
	}

	if !foundSize {
		canvasW = 1920.0
		canvasH = 1080.0
	}

	ctx.slideWidth = canvasW
	ctx.slideHeight = canvasH

	if debugMode {
		fmt.Printf("DEBUG: Using actual slide dimensions: %.0f x %.0f\n", canvasW, canvasH)
	}

	container := E("container", []string{"class", "slide-container"})
	for _, id := range ids {
		slide := ctx.ix.Records[id].(*KN.SlideArchive)
		div := E("div", []string{"class", "slide", "style", fmt.Sprintf("aspect-ratio: %.0f / %.0f;", canvasW, canvasH)})
		for _, d := range append([]*TSP.Reference{slide.BodyPlaceholder}, slide.OwnedDrawables...) {
			if d == nil {
				continue
			}
			e := ctx.processDrawable(d)
			if e != nil {
				div.AppendChild(e)
			}
		}
		container.AppendChild(div)
	}
	body.AppendChild(container)

	if img := ctx.renderImgData(); img != nil {
		body.AppendChild(img)
	}

	style := E("style")
	style.AppendChild(T(
		"html, body { height: 100%; background: transparent; }\n" +
			"body { margin: 0; font-family: -apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Arial,'Noto Sans','Liberation Sans',sans-serif; color: #111; }\n" +
			".slide-container { display: flex; flex-direction: column; align-items: center; gap: 28px; padding: 40px 24px 80px; box-sizing: border-box; }\n" +
			".slide { position: relative; width: min(1200px, calc(100vw - 48px)); aspect-ratio: 16 / 9; background: transparent; overflow: hidden; box-shadow: 0 10px 30px rgba(0,0,0,0.25); box-sizing: border-box; }\n" +
			".slide:hover { box-shadow: 0 14px 40px rgba(0,0,0,0.35); }\n" +
			"p { margin: 0; }\n" +
			".slide img { max-width: 100%; max-height: 100%; height: auto; object-fit: contain; display: block; position: relative; z-index: 0; }\n" +
			".slide .background-img { position: absolute; inset: 0; width: 100%; height: 100%; object-fit: cover; object-position: center; z-index: 0; }\n" +
			".slide > :not(.background-img) { position: relative; z-index: 1; }\n" +
			".slide p, .slide span, .slide b, .slide em, .slide h1, .slide h2, .slide h3, .slide h4 { position: relative; z-index: 2; }\n"))
	for k, v := range ctx.styles {
		style.AppendChild(T(fmt.Sprintf(".%s {\n%s}\n", k, v)))
	}
	head.AppendChild(style)

	// Add JavaScript to dynamically calculate slide width and scale fonts
	script := E("script")
	script.AppendChild(T(fmt.Sprintf(`
		(function() {
			var originalSlideWidth = %.0f;
			var scaleFactor = 1;
			
			function updateFontScale() {
				var slides = document.querySelectorAll('.slide');
				if (slides.length > 0) {
					var renderWidth = slides[0].offsetWidth;
					var viewportWidth = window.innerWidth;
					var expectedWidth = Math.min(1200, viewportWidth - 48);
					
					scaleFactor = Math.min(1, renderWidth / originalSlideWidth);
					document.documentElement.style.setProperty('--slide-scale', scaleFactor);
					
				}
			}

			updateFontScale();
			window.addEventListener('resize', updateFontScale);
			
			window.addEventListener('load', updateFontScale);
		})();
	`, ctx.slideWidth)))
	head.AppendChild(script)

	return doc
}

func mergeCharProps(props *TSWP.CharacterStylePropertiesArchive, parent *TSWP.CharacterStylePropertiesArchive) {
	if parent != nil {
		if props.Bold == nil {
			props.Bold = parent.Bold
		}
		if props.Italic == nil {
			props.Italic = parent.Italic
		}
		if props.FontSize == nil {
			props.FontSize = parent.FontSize
		}
		if props.FontName == nil {
			props.FontName = parent.FontName
		}
	}
}

// translateCharProps converts a TSWP.CharacterStylePropertiesArchive into CSS
func translateCharProps(ctx *Context, props *TSWP.CharacterStylePropertiesArchive) string {
	if props == nil {
		return ""
	}

	rval := ""
	if props.Bold != nil && *props.Bold {
		rval += "font-weight: bold;"
	}
	if props.Italic != nil && *props.Italic {
		rval += "font-style: italic;"
	}
	if props.FontSize != nil {
		fs := float64(*props.FontSize)
		rval += fmt.Sprintf("font-size: calc(var(--slide-scale, 1) * %.2fpt);", fs)
	}
	if props.FontName != nil {
		rval += fmt.Sprintf("font-family: '%s';", *props.FontName)
	}
	return rval
}

func mergeParaProps(props *TSWP.ParagraphStylePropertiesArchive, parent *TSWP.ParagraphStylePropertiesArchive) {
	if parent != nil {
		if props.LeftIndent == nil {
			props.LeftIndent = parent.LeftIndent
		}
		if props.RightIndent == nil {
			props.RightIndent = parent.RightIndent
		}
		if props.SpaceBefore == nil {
			props.SpaceBefore = parent.SpaceBefore
		}
		if props.SpaceAfter == nil {
			props.SpaceAfter = parent.SpaceAfter
		}
		if props.FirstLineIndent == nil {
			props.FirstLineIndent = parent.FirstLineIndent
		}
	}
}

// translateParaProps converts a TSWP.ParagraphStylePropertiesArchive into CSS.
func translateParaProps(ix *index.Index, props *TSWP.ParagraphStylePropertiesArchive) string {
	rval := ""
	// text alignment
	if props.Alignment != nil {
		switch *props.Alignment {
		case TSWP.ParagraphStylePropertiesArchive_TATvalue0:
			rval += "  text-align: left;\n"
		case TSWP.ParagraphStylePropertiesArchive_TATvalue1:
			rval += "  text-align: right;\n"
		case TSWP.ParagraphStylePropertiesArchive_TATvalue2:
			rval += "  text-align: center;\n"
		case TSWP.ParagraphStylePropertiesArchive_TATvalue3:
			rval += "  text-align: justify;\n"
		}
	}

	if props.LeftIndent != nil && *props.LeftIndent != 0. {
		rval += fmt.Sprintf("  margin-left: %fpt;\n", *props.LeftIndent)
	}

	if props.FirstLineIndent != nil {
		textIndent := *props.FirstLineIndent
		if props.LeftIndent != nil {
			textIndent = textIndent - *props.LeftIndent
		}
		if textIndent != 0. {
			rval += fmt.Sprintf("  text-indent: %fpt;\n", textIndent)
		}
	}

	if props.RightIndent != nil && *props.RightIndent > 0. {
		rval += fmt.Sprintf("  margin-right: %fpt;\n", *props.RightIndent)
	}
	if props.SpaceBefore != nil && *props.SpaceBefore > 0. {
		rval += fmt.Sprintf("  margin-top: %fpt;\n", *props.SpaceBefore)
	}
	if props.SpaceAfter != nil && *props.SpaceAfter > 0. {
		rval += fmt.Sprintf("  margin-bottom: %fpt;\n", *props.SpaceAfter)
	}

	// Paragraph background fill
	if props.Fill != nil {
		if css := colorToCSS(props.Fill); css != "" {
			if strings.HasPrefix(css, "rgba(") || strings.HasPrefix(css, "rgb(") || strings.HasPrefix(css, "#") {
				rval += fmt.Sprintf("  background-color:%s;\n", css)
			} else {
				rval += fmt.Sprintf("  background:%s;\n", css)
			}
		}
	}

	// List style processing - re-enable simple list style processing
	if props.ListStyleNull == nil || !*props.ListStyleNull {
		if props.ListStyle != nil {
			// Process list styles
			listCSS := translateListStyle(ix, props.ListStyle)
			if listCSS != "" {
				rval += listCSS
			}
		}
	}

	return rval
}

func writejson(foo interface{}, fn string) {
	a, err := json.MarshalIndent(foo, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}
	ioutil.WriteFile(fn, a, 0o644)
}

func dumpjson(foo interface{}) {
	a, err := json.MarshalIndent(foo, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(a))
}

func dump(foo interface{}) {
	fmt.Printf("%#v\n", foo)
}

func ptype(x interface{}) {
	fmt.Printf("type %T\n", x)
}

// translateListStyle processes list styles and returns CSS styles
func translateListStyle(ix *index.Index, listStyleRef *TSP.Reference) string {
	if listStyleRef == nil {
		return ""
	}

	// Get ListStyleArchive from reference
	ls, ok := ix.Deref(listStyleRef).(*TSWP.ListStyleArchive)
	if !ok {
		fmt.Printf("*** List style is not ListStyleArchive type: %T\n", ix.Deref(listStyleRef))
		return ""
	}

	rval := ""

	// Improved list style processing
	if len(ls.LabelTypes) > 0 {
		labelType := ls.LabelTypes[0]
		switch labelType {
		case TSWP.ListStyleArchive_kNumber:
			// Numbered list
			rval += "list-style-type: decimal;"
		case TSWP.ListStyleArchive_kString:
			// String list
			rval += "list-style-type: disc;"
		case TSWP.ListStyleArchive_kImage:
			// Image list
			rval += "list-style-type: disc;"
		default:
			// Default to disc
			rval += "list-style-type: disc;"
		}
	} else {
		// Default to disc
		rval += "list-style-type: disc;"
	}

	// Process indentation
	if len(ls.Indents) > 0 {
		indent := ls.Indents[0]
		if rval != "" {
			rval += " "
		}
		rval += fmt.Sprintf("margin-left: %fpt;", indent)
	}

	// Add basic list styles
	if rval != "" {
		rval += " "
	}
	rval += "margin-bottom: 6pt;"

	return rval
}

// processCellBorders processes cell borders and rounded corner styles
func processCellBorders(style *string, props *TST.CellStylePropertiesArchive) {
	// Process four-side borders
	processCellStroke(style, "border-top", props.TopStroke)
	processCellStroke(style, "border-right", props.RightStroke)
	processCellStroke(style, "border-bottom", props.BottomStroke)
	processCellStroke(style, "border-left", props.LeftStroke)

	// Check for rounded corners (by checking border Join property)
	hasRoundJoin := false
	if props.TopStroke != nil && props.TopStroke.Join != nil && *props.TopStroke.Join == TSD.LineJoin_RoundJoin {
		hasRoundJoin = true
	}
	if props.RightStroke != nil && props.RightStroke.Join != nil && *props.RightStroke.Join == TSD.LineJoin_RoundJoin {
		hasRoundJoin = true
	}
	if props.BottomStroke != nil && props.BottomStroke.Join != nil && *props.BottomStroke.Join == TSD.LineJoin_RoundJoin {
		hasRoundJoin = true
	}
	if props.LeftStroke != nil && props.LeftStroke.Join != nil && *props.LeftStroke.Join == TSD.LineJoin_RoundJoin {
		hasRoundJoin = true
	}

	if hasRoundJoin {
		*style += "border-radius: 4px;"
	}
}

// processCellStroke processes a single border
func processCellStroke(style *string, borderProp string, stroke *TSD.StrokeArchive) {
	if stroke != nil && stroke.Color != nil {
		col := colorToCSS(stroke.Color)
		w := 1.0
		if stroke.Width != nil {
			w = float64(*stroke.Width)
		}
		if col != "" {
			*style += fmt.Sprintf("%s: %.2fpx solid %s !important;", borderProp, w, col)
		}
	}
}

// processCellFont processes cell font styles
func processCellFont(style *string, props *TST.CellStylePropertiesArchive) {
	// Process cell font styles
	if props == nil {
		return
	}

	// CellStylePropertiesArchive mainly handles cell layout and borders
	// Font styles are handled through character styles, here we only handle cell-level font settings

	// Process text wrapping
	if props.TextWrap != nil && *props.TextWrap {
		*style += "white-space: normal;"
	} else {
		*style += "white-space: nowrap;"
	}

	// Process vertical alignment
	if props.VerticalAlignment != nil {
		switch *props.VerticalAlignment {
		case 0: // Top alignment
			*style += "vertical-align: top;"
		case 1: // Middle alignment
			*style += "vertical-align: middle;"
		case 2: // Bottom alignment
			*style += "vertical-align: bottom;"
		}
	}

	// Process padding
	if props.Padding != nil {
		if props.Padding.Top != nil {
			*style += fmt.Sprintf("padding-top: %.2fpt;", float64(*props.Padding.Top))
		}
		if props.Padding.Right != nil {
			*style += fmt.Sprintf("padding-right: %.2fpt;", float64(*props.Padding.Right))
		}
		if props.Padding.Bottom != nil {
			*style += fmt.Sprintf("padding-bottom: %.2fpt;", float64(*props.Padding.Bottom))
		}
		if props.Padding.Left != nil {
			*style += fmt.Sprintf("padding-left: %.2fpt;", float64(*props.Padding.Left))
		}
	}
}

// detectEmptyColumns detects which columns are empty (have no content)
func (ctx *Context) detectEmptyColumns(tm *TST.TableModelArchive, stringTable []*TST.TableDataList_ListEntry, richTable []*TST.TableDataList_ListEntry) []bool {
	cc := int(*tm.NumberOfColumns)
	emptyColumns := make([]bool, cc)

	// Initialize all columns as empty
	for i := 0; i < cc; i++ {
		emptyColumns[i] = true
	}

	// Iterate through all tiles and rows to check if each column has content
	if tm.BaseDataStore != nil && tm.BaseDataStore.Tiles != nil {
		for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
			tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
			for _, rinfo := range tile.RowInfos {
				// Decode the column -> offset mapping for this row
				offsets := make([]uint16, len(rinfo.CellOffsets)/2)
				binary.Read(bytes.NewBuffer(rinfo.CellOffsets), LE, offsets)

				// Check each column
				for c := 0; c < cc && c < len(offsets); c++ {
					offset := offsets[c]
					if debugTableCells {
						fmt.Printf("DEBUG: Checking column %d, offset=%d\n", c, offset)
					}
					if offset != 65535 { // Not an empty cell
						// Check if this cell has actual content
						hasContent := ctx.hasCellContent(rinfo, offset, stringTable, richTable)
						if debugTableCells {
							fmt.Printf("DEBUG: Column %d has content: %v\n", c, hasContent)
						}
						if hasContent {
							emptyColumns[c] = false
						}
					} else {
						if debugTableCells {
							fmt.Printf("DEBUG: Column %d is empty (offset=65535)\n", c)
						}
					}
				}
			}
		}
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Empty columns detection: ")
		for i, empty := range emptyColumns {
			if empty {
				fmt.Printf("Column %d (empty) ", i)
			}
		}
		fmt.Printf("\n")
	}

	return emptyColumns
}

// hasCellContent checks if the specified cell has actual content
func (ctx *Context) hasCellContent(rinfo *TST.TileRowInfo, offset uint16, stringTable []*TST.TableDataList_ListEntry, richTable []*TST.TableDataList_ListEntry) bool {
	// Try to read key value from buffer
	if len(rinfo.CellStorageBuffer) > int(offset)+8 {
		possibleKey := LE.Uint32(rinfo.CellStorageBuffer[offset+4 : offset+8])

		// Check richTable
		if len(richTable) > 0 {
			for _, entry := range richTable {
				if *entry.Key == possibleKey {
					// Check if rich text has actual content
					if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
						if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
							// Check if storage has text content
							return ctx.hasStorageContent(st)
						}
					}
					return false
				}
			}
		}

		// Check stringTable
		if len(stringTable) > 0 {
			for _, entry := range stringTable {
				if *entry.Key == possibleKey {
					// Check if string is not empty
					return entry.String_ != nil && len(strings.TrimSpace(*entry.String_)) > 0
				}
			}
		}
	}

	// If no matching key value is found in buffer, this cell has no actual content
	// Because if there really is content, it should be found in stringTable or richTable with corresponding key value
	if debugTableCells {
		fmt.Printf("DEBUG: No matching key found in tables, treating as empty\n")
	}
	return false
}

// hasStorageContent checks if storage has actual text content
func (ctx *Context) hasStorageContent(st *TSWP.StorageArchive) bool {
	if st == nil {
		return false
	}

	// Check if there is text content
	if st.Text != nil && len(st.Text) > 0 {
		for _, text := range st.Text {
			if len(strings.TrimSpace(text)) > 0 {
				return true
			}
		}
	}

	return false
}

// calculateColumnWidths calculates column widths based on empty column information
func (ctx *Context) calculateColumnWidths(totalColumns int, emptyColumns []bool) []float64 {
	widths := make([]float64, totalColumns)

	// Calculate the number of non-empty columns
	nonEmptyCount := 0
	for _, empty := range emptyColumns {
		if !empty {
			nonEmptyCount++
		}
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Total columns: %d, Non-empty columns: %d\n", totalColumns, nonEmptyCount)
	}

	// If all columns are empty, first column takes 100%
	if nonEmptyCount == 0 {
		widths[0] = 100.0
		for i := 1; i < totalColumns; i++ {
			widths[i] = 0.0
		}
		if debugTableCells {
			fmt.Printf("DEBUG: All columns empty, first column gets 100%%\n")
		}
	} else {
		// Non-empty columns share width equally, empty columns have 0 width
		widthPerNonEmpty := 100.0 / float64(nonEmptyCount)
		for i := 0; i < totalColumns; i++ {
			if emptyColumns[i] {
				widths[i] = 0.0
			} else {
				widths[i] = widthPerNonEmpty
			}
		}
		if debugTableCells {
			fmt.Printf("DEBUG: Non-empty columns get %.6f%% each\n", widthPerNonEmpty)
		}
	}

	return widths
}

// rearrangeContentToNonEmptyColumns rearranges content to fill all content into non-empty columns
func (ctx *Context) rearrangeContentToNonEmptyColumns(offsets []uint16, emptyColumns []bool, currentRow int) []uint16 {
	// Calculate the number of non-empty columns
	nonEmptyCount := 0
	for _, empty := range emptyColumns {
		if !empty {
			nonEmptyCount++
		}
	}

	// Create new offsets array
	newOffsets := make([]uint16, len(offsets))

	// Collect all content offsets (regardless of whether they are in empty columns)
	contentOffsets := make([]uint16, 0)
	for _, offset := range offsets {
		if offset != 65535 {
			contentOffsets = append(contentOffsets, offset)
		}
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Row %d - Found %d content offsets: %v\n", currentRow, len(contentOffsets), contentOffsets)
	}

	// If all columns are empty, fill all content into the first column
	if nonEmptyCount == 0 {
		if len(contentOffsets) > 0 {
			// Put all content in the first column, keep other columns empty
			newOffsets[0] = contentOffsets[0]
			for i := 1; i < len(newOffsets); i++ {
				newOffsets[i] = 65535 // Empty cell
			}
		}
	} else {
		// If there are non-empty columns, fill content into non-empty columns
		contentIndex := 0
		for i := 0; i < len(newOffsets); i++ {
			if i < len(emptyColumns) && !emptyColumns[i] && contentIndex < len(contentOffsets) {
				newOffsets[i] = contentOffsets[contentIndex]
				contentIndex++
			} else {
				newOffsets[i] = 65535 // Empty cell
			}
		}
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Row %d - Rearranged content: original offsets %v -> new offsets %v\n", currentRow, offsets, newOffsets)
	}

	return newOffsets
}

// collectAllContentOffsets collects content offsets for all columns of all rows for rearrangement
func (ctx *Context) collectAllContentOffsets(tm *TST.TableModelArchive) [][]uint16 {
	var allContentOffsets []uint16

	// First collect all content offsets
	if tm.BaseDataStore != nil && tm.BaseDataStore.Tiles != nil {
		for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
			tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
			for _, rinfo := range tile.RowInfos {
				// Decode the column -> offset mapping for this row
				offsets := make([]uint16, len(rinfo.CellOffsets)/2)
				binary.Read(bytes.NewBuffer(rinfo.CellOffsets), LE, offsets)

				// Collect all non-empty offsets, but check if they are within buffer range
				for _, offset := range offsets {
					if offset != 65535 && int(offset) < len(rinfo.CellStorageBuffer) {
						allContentOffsets = append(allContentOffsets, offset)
					}
				}
			}
		}
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Collected %d content offsets: %v\n", len(allContentOffsets), allContentOffsets)
	}

	// Calculate total number of rows
	totalRows := 0
	if tm.BaseDataStore != nil && tm.BaseDataStore.Tiles != nil {
		for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
			tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
			totalRows += len(tile.RowInfos)
		}
	}

	// Calculate number of columns
	cc := int(*tm.NumberOfColumns)

	// Rearrange content: fill all content sequentially into the first column
	var rearrangedOffsets [][]uint16
	contentIndex := 0

	for row := 0; row < totalRows; row++ {
		rowOffsets := make([]uint16, cc)

		// Fill content in first column, keep other columns empty
		if contentIndex < len(allContentOffsets) {
			rowOffsets[0] = allContentOffsets[contentIndex]
			contentIndex++
		} else {
			rowOffsets[0] = 65535
		}

		// Keep other columns empty
		for c := 1; c < cc; c++ {
			rowOffsets[c] = 65535
		}

		rearrangedOffsets = append(rearrangedOffsets, rowOffsets)
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Rearranged into %d rows, %d columns\n", len(rearrangedOffsets), cc)
	}

	return rearrangedOffsets
}
