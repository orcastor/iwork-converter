package iwork2text

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
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
)

// Global debug mode flag
var (
	debugMode       bool
	debugTableCells bool
)

// SetDebugMode sets the global debug mode
func SetDebugMode(debug bool) {
	debugMode = debug
	debugTableCells = debug
}

var LE = binary.LittleEndian

// popcount returns the number of set bits in a uint16
func popcount(x uint16) int {
	count := 0
	for x != 0 {
		count++
		x &= x - 1
	}
	return count
}

type Context struct {
	ix *index.Index
	zr *zip.ReadCloser
}

type Attachment struct {
	pos uint32
	doc string
}

func (ctx *Context) processImage(image *TSD.ImageArchive, ocr func(io.Reader) (string, error)) string {
	dataId := *image.Data.Identifier
	meta := ctx.ix.Records[2].(*TSP.PackageMetadata)
	var src string
	for _, data := range meta.Datas {
		if dataId == *data.Identifier {
			if data.FileName != nil {
				src = *data.FileName
			} else {
				fmt.Printf("No filename: %#v\n", data)
				src = *data.PreferredFileName
			}
		}
	}

	for _, f := range ctx.zr.File {
		if f.Name != "Data/"+src {
			continue
		}
		rc, _ := f.Open()
		defer rc.Close()

		imageBytes, _ := io.ReadAll(rc)
		if d, err := ocr(bytes.NewReader(imageBytes)); err == nil {
			return d
		}
		break
	}
	return ""
}

func (ctx *Context) processTable(tm *TST.TableModelArchive, ocr func(io.Reader) (string, error)) string {
	var doc string

	stringTable := ctx.ix.Deref(tm.BaseDataStore.StringTable).(*TST.TableDataList).Entries
	richTable := ctx.ix.Deref(tm.BaseDataStore.RichTextTable).(*TST.TableDataList).Entries

	// Debug: print stringTable contents
	if debugTableCells {
		fmt.Printf("DEBUG: stringTable has %d entries\n", len(stringTable))
		for i, entry := range stringTable {
			if i < 10 { // Only print first 10 entries
				fmt.Printf("DEBUG: stringTable[%d]: key=%d, string=%s\n", i, *entry.Key, *entry.String_)
			}
		}
	}

	// Format table is accessed directly in cellType=5 processing

	// rc := *tm.NumberOfRows
	cc := *tm.NumberOfColumns

	// I found some hints at http://stingrayreader.sourceforge.net/workbook/numbers_13.html about how
	// this works, but I'm still flying blind.

	// so for now we assume at most one tile per row, and rows are in the right order.  I suspect long rows (more than
	// 255 columns) will have multiple tiles, however.  This would likely only happen in a spreadsheet.

	for _, tinfo := range tm.BaseDataStore.Tiles.Tiles {
		tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
		for rowIndex, rinfo := range tile.RowInfos {
			// Try using cell_offsets first, fallback to cell_offsets_pre_bnc
			cellOffsets := rinfo.CellOffsets
			cellStorageBuffer := rinfo.CellStorageBuffer
			if len(cellOffsets) == 0 && len(rinfo.CellOffsetsPreBnc) > 0 {
				cellOffsets = rinfo.CellOffsetsPreBnc
				cellStorageBuffer = rinfo.CellStorageBufferPreBnc
			}

			offsets := make([]uint16, len(cellOffsets)/2)
			binary.Read(bytes.NewBuffer(cellOffsets), LE, offsets)
			// Navely assuming that the index is the column number, per the stringrayreader code.
			// FIXME - figure out the right way to determine column number.
			for c, offset := range offsets {
				if uint32(c) >= cc {
					break
				}

				// 0xffff is an empty cell (This only occurs at the end in my sample document.)
				if offset == 65535 {
					continue
				}

				var cellType int
				// Parse cellType from byte[1] (when byte[0]=5) - same as iwork2html
				if cellStorageBuffer[offset] == 5 {
					cellType = int(cellStorageBuffer[offset+1])
				}

				// Try both key calculation methods
				// Method 1: Simple offset+12 (like iwork2html)
				key := LE.Uint32(cellStorageBuffer[offset+12 : offset+16])

				// Method 2: Original popcount method (fallback)
				if key == 0 {
					flags := LE.Uint16(cellStorageBuffer[offset+4 : offset+6])
					o := popcount(flags)*4 + 8 + int(offset)
					key = LE.Uint32(cellStorageBuffer[o : o+4])
				}

				if debugTableCells {
					fmt.Printf("DEBUG: Row %d, Col %d: cellType=%d, key=%d, offset=%d, buffer[offset]=%d\n", rowIndex, c, cellType, key, offset, cellStorageBuffer[offset])
					if cellType == 0 {
						fmt.Printf("DEBUG: Found cellType=0 cell at Row %d, Col %d with key=%d\n", rowIndex, c, key)
					}
				}
				// fmt.Printf("P %d %x %d %x\n", cellType, flags, popcount(flags), cellStorageBuffer[o:o+4])
				// version := LE.Uint16(cellStorageBuffer[offset : offset+2])
				// fmt.Println("XXX", c, version, cellType, hex.EncodeToString(cellStorageBuffer[offset:]))
				switch cellType {
				case 0:
					// cellType=0: control cell (popup menu, checkbox, etc.)
					// Only process if this is actually a control cell (not a misidentified cell)
					optionText := ""
					found := false

					// Only try PopUpMenuModel if we have a valid key and MultipleChoiceListFormatTable
					if key > 0 && tm.BaseDataStore != nil && tm.BaseDataStore.MultipleChoiceListFormatTable != nil {
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
													found = true
													break // Found the option, no need to check other entries
												}
											}
										}
									}
								}
							}
						}
					}

					// Only add content if we found something from PopUpMenuModel
					// Don't fallback to other tables for cellType=0 to avoid duplicates
					if found {
						if debugTableCells {
							fmt.Printf("DEBUG: cellType=0 (control cell), key=%d, text=%s\n", key, optionText)
						}
						doc += " " + optionText
					} else if debugTableCells {
						fmt.Printf("DEBUG: cellType=0 (control cell), key=%d, no content found\n", key)
					}
				case 2: // number
					if debugTableCells {
						fmt.Printf("DEBUG: cellType=2, using key as literal number: %d\n", key)
					}
					doc += " " + fmt.Sprint(key)
				case 5: // date
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
						if dateText == "" && len(cellStorageBuffer) >= int(offset)+20 {
							value := math.Float64frombits(LE.Uint64(cellStorageBuffer[offset+12 : offset+20]))
							if value != 0 {
								tm := time.Unix(int64(value+978307200), 0) // Apple to unix epoch 2001-01-01
								dateText = tm.Format("2006-01-02")

								if debugTableCells {
									fmt.Printf("DEBUG: cellType=5 found date from buffer: %s (key=%d, value=%f)\n", dateText, key, value)
								}
							}
						}
					}

					if debugTableCells {
						fmt.Printf("DEBUG: cellType=5 (date cell), key=%d, text=%s\n", key, dateText)
					}
					doc += " " + dateText
				case 6: // boolean
					// Try simple offset+12 first, then fallback to popcount method
					var value float64
					if len(cellStorageBuffer) >= int(offset)+20 {
						value = math.Float64frombits(LE.Uint64(cellStorageBuffer[offset+12 : offset+20]))
					} else {
						// Fallback to popcount method
						flags := LE.Uint16(cellStorageBuffer[offset+4 : offset+6])
						o := popcount(flags)*4 + 8 + int(offset)
						if len(cellStorageBuffer) >= o+8 {
							value = math.Float64frombits(LE.Uint64(cellStorageBuffer[o : o+8]))
						}
					}
					label := "???"
					if value == 0 {
						label = "FALSE"
					} else if value == 0xf03f {
						label = "TRUE"
					}
					doc += " " + label
				case 3:
					// cellType=3: string content from tables
					// Try stringTable first
					found := false
					if debugTableCells {
						fmt.Printf("DEBUG: Looking for string content with key %d in stringTable with %d entries\n", key, len(stringTable))
					}
					for _, entry := range stringTable {
						if *entry.Key == key {
							doc += " " + *entry.String_
							if debugTableCells {
								fmt.Printf("DEBUG: Found string content for key %d: %s, doc so far: %s\n", key, *entry.String_, doc)
							}
							found = true
							break
						}
					}

					// Fallback to richTable
					if !found {
						for _, entry := range richTable {
							if *entry.Key == key {
								if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
									if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
										if debugTableCells {
											fmt.Printf("DEBUG: Found rich text content for key %d\n", key)
										}
										if d, err := ctx.storageToNode(st, ocr); err == nil {
											doc += " " + d
										}
									}
								}
								break
							}
						}
					}
				case 9:
					for _, entry := range richTable {
						if *entry.Key == key {
							rt := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive)
							st := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive)
							if d, err := ctx.storageToNode(st, ocr); err == nil {
								doc += " " + d
							}
							break
						}
					}
				default:
					// Unknown cellType: try both tables
					found := false
					for _, entry := range richTable {
						if *entry.Key == key {
							if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
								if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
									if debugTableCells {
										fmt.Printf("DEBUG: Found rich text content for key %d\n", key)
									}
									if d, err := ctx.storageToNode(st, ocr); err == nil {
										doc += " " + d
									}
									found = true
								}
							}
							break
						}
					}
					if !found {
						for _, entry := range stringTable {
							if *entry.Key == key {
								doc += " " + *entry.String_
								if debugTableCells {
									fmt.Printf("DEBUG: Found string content for key %d: %s\n", key, *entry.String_)
								}
								found = true
								break
							}
						}
					}
					if !found && debugTableCells {
						fmt.Printf("DEBUG: No content found for cell with key %d, cellType=%d\n", key, cellType)
					}
				}
			}
			// Add newline after each row
			doc += "\n"
		}
	}
	return doc
}

func (ctx *Context) processDrawable(ref *TSP.Reference, ocr func(io.Reader) (string, error)) string {
	item := ctx.ix.Deref(ref)
	switch item.(type) {
	case *TSD.ImageArchive:
		return ctx.processImage(item.(*TSD.ImageArchive), ocr)
	case *TST.WPTableInfoArchive:
		table := item.(*TST.WPTableInfoArchive)
		tm := ctx.ix.Deref(table.Super.TableModel).(*TST.TableModelArchive)
		return ctx.processTable(tm, ocr)
	case *TST.TableInfoArchive:
		tm := ctx.ix.Deref(item.(*TST.TableInfoArchive).TableModel).(*TST.TableModelArchive)
		return ctx.processTable(tm, ocr)
	case *TSWP.ShapeInfoArchive:
		return ctx.processShapeInfo(item.(*TSWP.ShapeInfoArchive), ocr)
	case *TSD.GroupArchive:
		return ctx.processDrawableArchive(item.(*TSD.GroupArchive).Super, ocr)
	case *KN.PlaceholderArchive:
		return ctx.processShapeInfo(item.(*KN.PlaceholderArchive).Super, ocr)
	default:
		msg := fmt.Sprintf("*** Unhandled attachment type %T\n", item)
		fmt.Println(msg)
		return ""
	}
}

func (ctx *Context) processShapeInfo(sia *TSWP.ShapeInfoArchive, ocr func(io.Reader) (string, error)) string {
	if cs := ctx.ix.Deref(sia.OwnedStorage).(*TSWP.StorageArchive); cs != nil {
		if doc, err := ctx.storageToNode(cs, ocr); err == nil {
			return doc
		}
	}
	return ctx.processDrawableArchive(sia.Super.Super, ocr)
}

func (ctx *Context) processDrawableArchive(da *TSD.DrawableArchive, ocr func(io.Reader) (string, error)) string {
	// do nothing...
	return ""
}

// storageToNode extracts text content from a StorageArchive
func (ctx *Context) storageToNode(bs *TSWP.StorageArchive, ocr func(io.Reader) (string, error)) (string, error) {
	texts := bs.Text

	if len(texts) == 0 {
		return "", fmt.Errorf("no text content found")
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

	// Process paragraphs
	parStyles := bs.TableParaStyle.Entries
	var result string

	// Process each paragraph
	for i, e := range parStyles {
		pos := *e.CharacterIndex
		end := uint32(len(rr))
		if i+1 < len(parStyles) {
			end = *parStyles[i+1].CharacterIndex
		}

		// Extract paragraph text
		paragraphText := string(rr[pos:end])

		// Skip empty paragraphs
		if strings.TrimSpace(paragraphText) == "" {
			continue
		}

		// Add paragraph text to result
		result += paragraphText + "\n"
	}

	return result, nil
}

type Style map[string]interface{}

func ConvertString(in string, ocr func(io.Reader) (string, error)) (string, error) {
	if ocr == nil {
		ocr = func(r io.Reader) (string, error) { return "", nil }
	}
	fmt.Println("Processing", in)

	var err error
	var ctx Context
	if ctx.ix, err = index.Open(in); err != nil {
		return "", err
	}
	if ctx.zr, err = zip.OpenReader(in); err != nil {
		return "", err
	}
	defer ctx.zr.Close()

	if debugMode {
		fmt.Println("Read", len(ctx.ix.Records), "records")
	}

	var doc string
	switch ctx.ix.Type {
	case "pages":
		doc = ctx.processPages(ocr)
	case "numbers":
		doc = ctx.processNumbers(ocr)
	case "key":
		doc = ctx.processKeynote(ocr)
	}

	return doc, nil
}

func Convert(in, out string) error {
	doc, err := ConvertString(in, nil)
	if err != nil {
		return err
	}
	return os.WriteFile(out, []byte(doc), os.ModePerm)
}

// processPages translates a pages file.
func (ctx *Context) processPages(ocr func(io.Reader) (string, error)) string {
	var doc string

	da := ctx.ix.Records[1].(*TP.DocumentArchive)
	bs := ctx.ix.Deref(da.BodyStorage).(*TSWP.StorageArchive)

	fda := ctx.ix.Deref(da.FloatingDrawables).(*TP.FloatingDrawablesArchive)
	if len(fda.PageGroups) != 0 {
		fmt.Print(`WARNING - 
            This document has floating drawables (e.g. floating images/tables/text blocks) which we don't handle in HTML
            conversion.
            
            Figuring out where to place them in the document would probably be tricky.
`)
	}

	if d, err := ctx.storageToNode(bs, ocr); err == nil {
		doc += d
	}

	return doc
}

func (ctx *Context) processNumbers(ocr func(io.Reader) (string, error)) string {
	var doc string

	da := ctx.ix.Records[1].(*TN.DocumentArchive)

	for _, ref := range da.Sheets {
		sheet := ctx.ix.Deref(ref).(*TN.SheetArchive)
		for _, ref := range sheet.DrawableInfos {
			// if this cast throws there are other kinds of drawables...
			d := ctx.processDrawable(ref, ocr)
			if d != "" {
				doc += d
			}
		}
	}

	return doc
}

// processKeynote translates a keynote file.
func (ctx *Context) processKeynote(ocr func(io.Reader) (string, error)) string {
	var doc string

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
		for _, slideRef := range slideTreeSlides {
			if slideRef != nil && slideRef.Identifier != nil {
				slideNodeId := *slideRef.Identifier
				if slideNode, ok := ctx.ix.Records[slideNodeId].(*KN.SlideNodeArchive); ok {
					// Check if slide is skipped
					if slideNode.Slide != nil && slideNode.Slide.Identifier != nil {
						actualSlideId := *slideNode.Slide.Identifier
						ids = append(ids, actualSlideId)
					}
				}
			}
		}
	} else {
		// Fallback to original method
		for key, rec := range ctx.ix.Records {
			if _, ok := rec.(*KN.SlideArchive); ok {
				ids = append(ids, key)
			}
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	}

	for _, id := range ids {
		slide := ctx.ix.Records[id].(*KN.SlideArchive)
		for _, d := range append([]*TSP.Reference{slide.BodyPlaceholder}, slide.OwnedDrawables...) {
			if d == nil {
				continue
			}
			d := ctx.processDrawable(d, ocr)
			if d != "" {
				doc += d
			}
		}
	}

	return doc
}
