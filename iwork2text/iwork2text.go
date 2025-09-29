package iwork2text

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"time"

	"github.com/orcastor/iwork-converter/index"
	"github.com/orcastor/iwork-converter/proto/KN"
	"github.com/orcastor/iwork-converter/proto/TN"
	"github.com/orcastor/iwork-converter/proto/TP"
	"github.com/orcastor/iwork-converter/proto/TSD"
	"github.com/orcastor/iwork-converter/proto/TSP"
	"github.com/orcastor/iwork-converter/proto/TST"
	"github.com/orcastor/iwork-converter/proto/TSWP"
)

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

var LE = binary.LittleEndian

func popcount(v uint16) int {
	var c int
	for ; v != 0; c++ {
		v &= v - 1
	}
	return c
}

func (ctx *Context) processTable(tm *TST.TableModelArchive, ocr func(io.Reader) (string, error)) string {
	var doc string

	stringTable := ctx.ix.Deref(tm.DataStore.StringTable).(*TST.TableDataList).Entries
	richTable := ctx.ix.Deref(tm.DataStore.RichTextPayloadTable).(*TST.TableDataList).Entries

	// rc := *tm.NumberOfRows
	cc := *tm.NumberOfColumns

	// I found some hints at http://stingrayreader.sourceforge.net/workbook/numbers_13.html about how
	// this works, but I'm still flying blind.

	// so for now we assume at most one tile per row, and rows are in the right order.  I suspect long rows (more than
	// 255 columns) will have multiple tiles, however.  This would likely only happen in a spreadsheet.

	for _, tinfo := range tm.DataStore.Tiles.Tiles {
		tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
		for r, rinfo := range tile.RowInfos {
			offsets := make([]uint16, len(rinfo.CellOffsets)/2)
			binary.Read(bytes.NewBuffer(rinfo.CellOffsets), LE, offsets)
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
				// this has changed since I first wrote the code.  There is now a 4 in the first byte and the type in the next
				// the "stingrayreader" site says there is a halfword "version" and then the type, which I think worked at one
				// point, but I no longer have the file.
				if rinfo.CellStorageBuffer[offset] == 4 {
					cellType = int(rinfo.CellStorageBuffer[offset+1])
				} else {
					cellType = int(rinfo.CellStorageBuffer[offset+2])
				}

				// As far as I can tell, the records are variable length, with the pointer into the string/rich table at
				// the end, but this field seems to contain one bit per uint32 before the pointer to the string table
				// I suspect they are flags indicating which numbers/fields follow.
				flags := LE.Uint16(rinfo.CellStorageBuffer[offset+4 : offset+6])
				o := popcount(flags)*4 + 8 + int(offset)

				key := LE.Uint32(rinfo.CellStorageBuffer[o : o+4])

				// fmt.Printf("P %d %x %d %x\n", cellType, flags, popcount(flags), rinfo.CellStorageBuffer[o:o+4])
				// version := LE.Uint16(rinfo.CellStorageBuffer[offset : offset+2])
				// fmt.Println("XXX", c, version, cellType, hex.EncodeToString(rinfo.CellStorageBuffer[offset:]))
				switch cellType {
				case 0:
					// blank cells are type 0
				case 2: // number
					value := math.Float64frombits(LE.Uint64(rinfo.CellStorageBuffer[o : o+8]))
					doc += " " + fmt.Sprint(value)
				case 5: // date
					value := math.Float64frombits(LE.Uint64(rinfo.CellStorageBuffer[o : o+8]))
					value += 978307200 // Apple to unix epoch
					tm := time.Unix(int64(value), 0)
					doc += " " + fmt.Sprint(tm)
				case 6: // boolean
					value := math.Float64frombits(LE.Uint64(rinfo.CellStorageBuffer[o : o+8]))
					label := "???"
					if value == 0 {
						label = "FALSE"
					} else if value == 0xf03f {
						label = "TRUE"
					}
					doc += " " + label
				case 3:
					for _, entry := range stringTable {
						if *entry.Key == key {
							doc += " " + *entry.String_
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
						}
					}
				default:
					fmt.Printf("P %d %x %d %x\n", cellType, flags, popcount(flags), rinfo.CellStorageBuffer[o:o+8])
					fmt.Printf("CELL %d:%d type %d %s\n", r, c, cellType, hex.EncodeToString(rinfo.CellStorageBuffer[offset:]))
				}
			}
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
	if cs := ctx.ix.Deref(sia.ContainedStorage).(*TSWP.StorageArchive); cs != nil {
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

// storageToNode populates a html node with the contents of a StorageArchive. This happens with both the
// main body of the document and rich text table cells.
func (ctx *Context) storageToNode(bs *TSWP.StorageArchive, ocr func(io.Reader) (string, error)) (string, error) {
	var doc string

	texts := bs.Text
	if len(texts) != 1 {
		return doc, fmt.Errorf("FIXME - Expecting exactly one text, got %d", len(texts))
	}
	text := texts[0]

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
				d := ctx.processDrawable(archive.Drawable, ocr)
				if d != "" {
					attachments = append(attachments, Attachment{pos, d})
				}
			case *TSWP.NumberAttachmentArchive:
				// do nothing...
			}
		}
	}
	// bs.TableListStyle - seems to change on headings, look into it.

	// build paragraphs
	for i, e := range parStyles {

		pos := *e.CharacterIndex
		end := uint32(len(rr))
		if i+1 < len(parStyles) {
			end = *parStyles[i+1].CharacterIndex
		}

		for len(attachments) > 0 && attachments[0].pos < end {
			if attachments[0].pos != pos {
				fmt.Printf("FIXME - attachment not at start of paragraph - pstart=%d pend=%d att=%d par=%#v\n",
					pos, end, attachments[0].pos, string(rr[pos:end]))
			}
			doc += attachments[0].doc
			attachments = attachments[1:]
		}

		// <span> <em> and <b>
		if bs.TableCharStyle != nil {
			charStyles := bs.TableCharStyle.Entries
			for i, e := range charStyles { // build any span/em/b as needed
				cs := *e.CharacterIndex
				if cs < pos {
					continue
				}
				if cs >= end {
					break
				}
				ce := uint32(len(rr))
				if i+1 < len(charStyles) {
					ce = *charStyles[i+1].CharacterIndex
				}
				if ce > end {
					if e.Object != nil {
						// fmt.Println("ERR? ce > end", ce, end, e.Object)
					}
					ce = end
				}
				if cs > pos {
					doc += string(rr[pos:cs])
					pos = cs
				}
				doc += string(rr[cs:ce])
				pos = ce
			}
		}
		doc += string(rr[pos:end])
		doc += string("\n")
	}

	return doc, nil
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

	fmt.Println("Read", len(ctx.ix.Records), "records")

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
		fmt.Println(`WARNING - 
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

	meta := ctx.ix.Records[2].(*TSP.PackageMetadata)
	ids := []uint64{}
	for _, comp := range meta.Components {
		if *comp.PreferredLocator == "Slide" {
			ids = append(ids, *comp.Identifier)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	for _, id := range ids {
		slide := ctx.ix.Records[id].(*KN.SlideArchive)
		for _, d := range append([]*TSP.Reference{slide.BodyPlaceholder}, slide.Drawables...) {
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
