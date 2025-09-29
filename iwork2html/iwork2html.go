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
	"sort"
	"strings"

	"github.com/orcastor/iwork-converter/index"
	"github.com/orcastor/iwork-converter/proto/KN"
	"github.com/orcastor/iwork-converter/proto/TN"
	"github.com/orcastor/iwork-converter/proto/TP"
	"github.com/orcastor/iwork-converter/proto/TSD"
	"github.com/orcastor/iwork-converter/proto/TSP"
	"github.com/orcastor/iwork-converter/proto/TST"
	"github.com/orcastor/iwork-converter/proto/TSWP"

	"golang.org/x/net/html"
)

// T is a helper function for building html text nodes.
func T(value string) *html.Node {
	return &html.Node{Type: html.TextNode, Data: value}
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
}

// 控制是否输出表格单元格的调试日志
var debugTableCells = true

// 控制是否输出图片处理调试日志
var debugImages = true

// 控制是否启用智能表头检测
var enableSmartHeaderDetection = false

type Attachment struct {
	pos  uint32
	node *html.Node
}

func (ctx *Context) processImage(image *TSD.ImageArchive) *html.Node {
	if debugImages {
		fmt.Printf("DEBUG IMG: processing image, dataId=%d\n", *image.Data.Identifier)
	}
	dataId := *image.Data.Identifier
	meta := ctx.ix.Records[2].(*TSP.PackageMetadata)
	var src string
	for _, data := range meta.Datas {
		if dataId == *data.Identifier {
			if data.FileName != nil {
				src = *data.FileName
				if debugImages {
					fmt.Printf("DEBUG IMG: use FileName=%s for dataId=%d\n", src, dataId)
				}
			} else {
				fmt.Printf("No filename: %#v\n", data)
				src = *data.PreferredFileName
				if debugImages {
					fmt.Printf("DEBUG IMG: use PreferredFileName=%s for dataId=%d\n", src, dataId)
				}
			}
		}
	}
	ctx.imgs["Data/"+src] = dataId
	// not sure if this is px or pt.  It's px on the html side.
	width := fmt.Sprintf("%f", *image.OriginalSize.Width)
	height := fmt.Sprintf("%f", *image.OriginalSize.Height)
	if debugImages {
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

// applyCellStyle 应用单元格样式
func (ctx *Context) applyCellStyle(tm *TST.TableModelArchive, key uint32) string {
	style := ""

	if tm.DataStore != nil && tm.DataStore.StyleTable != nil {
		styleTableRef := ctx.ix.Deref(tm.DataStore.StyleTable)
		if tdl, ok := styleTableRef.(*TST.TableDataList); ok {
			for _, entry := range tdl.Entries {
				if *entry.Key == key && entry.Reference != nil {
					entryRef := ctx.ix.Deref(entry.Reference)
					if csa, ok := entryRef.(*TST.CellStyleArchive); ok {
						if csa.CellProperties != nil {
							// 处理背景填充
							if csa.CellProperties.CellFill != nil {
								if css := colorToCSS(csa.CellProperties.CellFill.GetColor()); css != "" {
									applyBackgroundColor(&style, css, true)
								}
							}
							// 处理边框和圆角
							processCellBorders(&style, csa.CellProperties)
							// 处理字体大小
							processCellFont(&style, csa.CellProperties)
						}
					} else if psa, ok := entryRef.(*TSWP.ParagraphStyleArchive); ok {
						// 处理段落样式作为单元格样式
						if psa.ParaProperties != nil {
							// 处理段落背景填充
							if psa.ParaProperties.Fill != nil {
								if css := colorToCSS(psa.ParaProperties.Fill); css != "" {
									applyBackgroundColor(&style, css, false)
								}
							}
							// 处理段落边框
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
						// 打印不认识的样式类型
						fmt.Printf("DEBUG: 不认识的样式类型: %T\n", entryRef)
					}
					break
				}
			}
		} else {
			// 打印不认识的StyleTable类型
			fmt.Printf("DEBUG: 不认识的StyleTable类型: %T\n", styleTableRef)
		}
	}

	return style
}

// applyBackgroundColor 统一处理背景色应用
func applyBackgroundColor(style *string, css string, important bool) {
	if css == "" {
		return
	}

	*style += "background:" + css
	if important {
		*style += " !important"
	}
	*style += ";"

	if strings.HasPrefix(css, "rgba(") || strings.HasPrefix(css, "rgb(") || strings.HasPrefix(css, "#") {
		*style += "background-color:" + css
		if important {
			*style += " !important"
		}
		*style += ";"
	}
}

// mergeParentStyles 递归处理父样式继承
func (ctx *Context) mergeParentStyles(child, parent *TSWP.ParagraphStyleArchive) {
	if parent.Super.Parent != nil {
		grandParent := ctx.ix.Deref(parent.Super.Parent).(*TSWP.ParagraphStyleArchive)
		// 先处理祖辈样式到父样式
		mergeCharProps(parent.CharProperties, grandParent.CharProperties)
		mergeParaProps(parent.ParaProperties, grandParent.ParaProperties)
		// 递归处理更深层的继承
		ctx.mergeParentStyles(parent, grandParent)
		// 然后将处理后的父样式应用到子样式
		mergeCharProps(child.CharProperties, parent.CharProperties)
		mergeParaProps(child.ParaProperties, parent.ParaProperties)
	}
}

// mergeParentCharStyles 递归处理父字符样式继承
func (ctx *Context) mergeParentCharStyles(child, parent *TSWP.CharacterStyleArchive) {
	if parent.Super.Parent != nil {
		grandParentRef := ctx.ix.Deref(parent.Super.Parent)
		if grandParent, ok := grandParentRef.(*TSWP.CharacterStyleArchive); ok {
			// 先处理祖辈样式到父样式
			mergeCharProps(parent.CharProperties, grandParent.CharProperties)
			// 递归处理更深层的继承
			ctx.mergeParentCharStyles(parent, grandParent)
			// 然后将处理后的父样式应用到子样式
			mergeCharProps(child.CharProperties, parent.CharProperties)
		} else {
			// 打印不认识的祖辈字符样式类型
			fmt.Printf("DEBUG: 不认识的祖辈字符样式类型: %T\n", grandParentRef)
		}
	}
}

// applyPositionBasedStyle 根据单元格位置应用样式
func (ctx *Context) applyPositionBasedStyle(tm *TST.TableModelArchive, globalRow, c int, shouldTreatFirstRowAsHeader bool) string {
	style := ""

	// 检查单元格位置类型
	var cellStyleRef *TSP.Reference
	var isHeaderRow, isHeaderColumn, isFooterRow bool

	// 检查是否为表头行
	if tm.NumberOfHeaderRows != nil && globalRow < int(*tm.NumberOfHeaderRows) {
		isHeaderRow = true
	} else if shouldTreatFirstRowAsHeader && globalRow == 0 {
		isHeaderRow = true
	}
	// 检查是否为表头列
	if tm.NumberOfHeaderColumns != nil && c < int(*tm.NumberOfHeaderColumns) {
		isHeaderColumn = true
	}
	// 检查是否为表尾行
	if tm.NumberOfFooterRows != nil && tm.NumberOfRows != nil && globalRow >= int(*tm.NumberOfRows-*tm.NumberOfFooterRows) {
		isFooterRow = true
	}

	// 优先级：表头行 > 表头列 > 表尾行 > 默认
	if isHeaderRow && tm.HeaderRowStyle != nil {
		cellStyleRef = tm.HeaderRowStyle
		if debugTableCells {
			fmt.Printf("DEBUG: Using HeaderRowStyle for row %d, col %d\n", globalRow, c)
		}
	} else if isHeaderColumn && tm.HeaderColumnStyle != nil {
		cellStyleRef = tm.HeaderColumnStyle
		if debugTableCells {
			fmt.Printf("DEBUG: Using HeaderColumnStyle for row %d, col %d\n", globalRow, c)
		}
	} else if isFooterRow && tm.FooterRowStyle != nil {
		cellStyleRef = tm.FooterRowStyle
		if debugTableCells {
			fmt.Printf("DEBUG: Using FooterRowStyle for row %d, col %d\n", globalRow, c)
		}
	} else if tm.BodyCellStyle != nil {
		cellStyleRef = tm.BodyCellStyle
		if debugTableCells {
			fmt.Printf("DEBUG: Using BodyCellStyle for row %d, col %d\n", globalRow, c)
		}
	}

	// 应用选中的样式
	if cellStyleRef != nil {
		cellStyleRefObj := ctx.ix.Deref(cellStyleRef)
		if debugTableCells {
			fmt.Printf("DEBUG: CellStyleRef type: %T\n", cellStyleRefObj)
		}
		if csp, ok := cellStyleRefObj.(*TST.CellStylePropertiesArchive); ok {
			// 处理背景填充
			if csp.CellFill != nil {
				if css := colorToCSS(csp.CellFill.GetColor()); css != "" {
					applyBackgroundColor(&style, css, true)
				}
			}
			// 处理边框
			processCellBorders(&style, csp)
		} else if csp, ok := cellStyleRefObj.(*TST.CellStyleArchive); ok {
			// 处理CellStyleArchive类型
			if csp.CellProperties != nil {
				// 处理背景填充
				if csp.CellProperties.CellFill != nil {
					if css := colorToCSS(csp.CellProperties.CellFill.GetColor()); css != "" {
						applyBackgroundColor(&style, css, true)
					}
				}
				// 处理边框
				processCellBorders(&style, csp.CellProperties)
				// 处理字体大小
				processCellFont(&style, csp.CellProperties)
			}
		} else {
			// 打印不认识的单元格样式属性类型
			fmt.Printf("DEBUG: 不认识的单元格样式属性类型: %T\n", cellStyleRefObj)
		}
	}

	return style
}

// analyzeFirstRowAsHeader 分析第一行内容，判断是否应该作为标题行
func (ctx *Context) analyzeFirstRowAsHeader(tm *TST.TableModelArchive, stringTable []*TST.TableDataList_ListEntry, richTable []*TST.TableDataList_ListEntry) bool {
	if tm.DataStore == nil || tm.DataStore.Tiles == nil || len(tm.DataStore.Tiles.Tiles) == 0 {
		return false
	}

	// 获取第一个tile的第一行
	firstTileInfo := tm.DataStore.Tiles.Tiles[0]
	firstTile := ctx.ix.Deref(firstTileInfo.Tile).(*TST.Tile)
	if len(firstTile.RowInfos) == 0 {
		return false
	}

	firstRowInfo := firstTile.RowInfos[0]

	// 解码第一行的 column -> offset 映射
	offsets := make([]uint16, len(firstRowInfo.CellOffsets)/2)
	binary.Read(bytes.NewBuffer(firstRowInfo.CellOffsets), LE, offsets)

	// 分析第一行的内容特征
	nonEmptyCells := 0
	textCells := 0
	hasHeaderStyle := false

	for _, offset := range offsets {
		if offset == 65535 { // 空单元格
			continue
		}

		nonEmptyCells++

		// 检查是否为文本内容
		if ctx.isTextCell(offset, stringTable, richTable) {
			textCells++
		}
	}

	// 检查是否有表头行样式定义
	if tm.HeaderRowStyle != nil {
		hasHeaderStyle = true
		if debugTableCells {
			fmt.Printf("DEBUG: Table has HeaderRowStyle defined\n")
		}
	}

	// 对于单列表格，如果第一行有任何内容，都认为是标题行
	if len(offsets) == 1 && nonEmptyCells > 0 {
		if debugTableCells {
			fmt.Printf("DEBUG: Single column table with content in first row - treating as header\n")
		}
		return true
	}

	// 更保守的标题行判断条件：
	// 只有在明确的标题行特征时才判断为标题行
	if nonEmptyCells > 0 {
		threshold := 0.8 // 提高阈值，减少误判
		// 如果第一行的内容都是文本，并且满足阈值要求，才认为是标题行
		if float64(textCells)/float64(nonEmptyCells) >= threshold {
			if debugTableCells {
				fmt.Printf("DEBUG: First row analysis - nonEmpty: %d, text: %d, ratio: %.2f (threshold: %.2f, hasHeaderStyle: %v)\n",
					nonEmptyCells, textCells, float64(textCells)/float64(nonEmptyCells), threshold, hasHeaderStyle)
			}
			return true
		}
	}

	return false
}

// isTextCell 检查给定offset的单元格是否包含文本内容
func (ctx *Context) isTextCell(offset uint16, stringTable []*TST.TableDataList_ListEntry, richTable []*TST.TableDataList_ListEntry) bool {
	// 尝试在字符串表中查找
	for _, entry := range stringTable {
		if entry.Key != nil && *entry.Key == uint32(offset) && entry.String_ != nil && *entry.String_ != "" {
			return true
		}
	}

	// 尝试在富文本表中查找
	for _, entry := range richTable {
		if entry.Key != nil && *entry.Key == uint32(offset) && entry.RichTextPayload != nil {
			return true
		}
	}

	return false
}

// analyzeTableStructure 分析表格结构，判断第一行是否应该作为标题行
func (ctx *Context) analyzeTableStructure(tm *TST.TableModelArchive, stringTable []*TST.TableDataList_ListEntry, richTable []*TST.TableDataList_ListEntry) bool {
	if tm.DataStore == nil || tm.DataStore.Tiles == nil || len(tm.DataStore.Tiles.Tiles) == 0 {
		return false
	}

	firstTile := ctx.ix.Deref(tm.DataStore.Tiles.Tiles[0].Tile).(*TST.Tile)
	if len(firstTile.RowInfos) < 2 {
		return false // 需要至少两行才能比较
	}

	// 分析前两行的内容密度差异
	firstRowContent := ctx.getRowContentDensity(firstTile.RowInfos[0], stringTable, richTable)
	secondRowContent := ctx.getRowContentDensity(firstTile.RowInfos[1], stringTable, richTable)

	if debugTableCells {
		fmt.Printf("DEBUG: Structure analysis - first row density: %.2f, second row density: %.2f\n",
			firstRowContent, secondRowContent)
	}

	// 更严格的条件：只有在第一行内容密度明显高于第二行时才判断为标题行
	if firstRowContent > 0.5 && firstRowContent > secondRowContent*2.0 {
		return true
	}

	// 更严格的条件：第一行内容丰富且第二行基本为空
	if firstRowContent > 0.7 && secondRowContent < 0.1 {
		return true
	}

	return false
}

// getRowContentDensity 计算行的内容密度（非空单元格比例）
func (ctx *Context) getRowContentDensity(rowInfo *TST.TileRowInfo, stringTable []*TST.TableDataList_ListEntry, richTable []*TST.TableDataList_ListEntry) float64 {
	offsets := make([]uint16, len(rowInfo.CellOffsets)/2)
	binary.Read(bytes.NewBuffer(rowInfo.CellOffsets), LE, offsets)

	nonEmptyCells := 0
	for _, offset := range offsets {
		if offset != 65535 && ctx.isTextCell(offset, stringTable, richTable) {
			nonEmptyCells++
		}
	}

	if len(offsets) == 0 {
		return 0
	}

	return float64(nonEmptyCells) / float64(len(offsets))
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
		fmt.Printf("DEBUG: Smart header detection enabled: %v\n", enableSmartHeaderDetection)

		// 检查数据存储结构
		if tm.DataStore != nil && tm.DataStore.Tiles != nil {
			fmt.Printf("DEBUG: DataStore has %d tiles\n", len(tm.DataStore.Tiles.Tiles))
			for i, tinfo := range tm.DataStore.Tiles.Tiles {
				tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
				fmt.Printf("DEBUG: Tile %d has %d rows\n", i, len(tile.RowInfos))
			}
		}
	}
	// 提取字符串和富文本表
	var stringTable []*TST.TableDataList_ListEntry
	var richTable []*TST.TableDataList_ListEntry
	if tm.DataStore != nil {
		if debugTableCells {
			fmt.Printf("DEBUG: DataStore found\n")
		}
		if tm.DataStore.StringTable != nil {
			if debugTableCells {
				fmt.Printf("DEBUG: StringTable reference found\n")
			}
			if tdl, ok := ctx.ix.Deref(tm.DataStore.StringTable).(*TST.TableDataList); ok {
				stringTable = tdl.Entries
				if debugTableCells {
					fmt.Printf("DEBUG: StringTable loaded with %d entries\n", len(stringTable))
					for i, entry := range stringTable {
						if i < 5 { // 只显示前5个条目
							fmt.Printf("DEBUG: StringTable[%d]: key=%d, value=%s\n", i, *entry.Key, *entry.String_)
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
		if tm.DataStore.RichTextPayloadTable != nil {
			if debugTableCells {
				fmt.Printf("DEBUG: RichTextPayloadTable reference found\n")
			}
			if tdl, ok := ctx.ix.Deref(tm.DataStore.RichTextPayloadTable).(*TST.TableDataList); ok {
				richTable = tdl.Entries
				if debugTableCells {
					fmt.Printf("DEBUG: RichTextPayloadTable loaded with %d entries\n", len(richTable))
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
					fmt.Printf("DEBUG: Failed to deref RichTextPayloadTable\n")
				}
			}
		} else {
			if debugTableCells {
				fmt.Printf("DEBUG: No RichTextPayloadTable reference\n")
			}
		}
	} else {
		if debugTableCells {
			fmt.Printf("DEBUG: No DataStore found\n")
		}
	}

	cc := int(*tm.NumberOfColumns)

	// 简化：始终显示所有列，不进行活跃列扫描
	activeColumns := make([]bool, cc)
	for i := 0; i < cc; i++ {
		activeColumns[i] = true
	}

	// 计算有内容的列数
	activeColumnCount := 0
	for _, active := range activeColumns {
		if active {
			activeColumnCount++
		}
	}

	// 智能表格结构检测：如果只有很少的列有内容，可能是表格结构错误
	if debugTableCells {
		fmt.Printf("DEBUG: Active columns: %d out of %d total columns\n", activeColumnCount, cc)
		for i, active := range activeColumns {
			if active {
				fmt.Printf("DEBUG: Column %d is active\n", i)
			}
		}
	}

	// 关闭智能重构：始终按原始列数渲染
	shouldRestructure := false

	table := E("table")

	// 生成列定义 - 智能宽度分配
	colgroup := E("colgroup")

	// 检查哪些列实际有内容
	// 根据实际的内容分配策略：只有第一列分配内容，其他列保持为空
	hasContentColumns := make([]bool, cc)
	hasContentColumns[0] = true // 第一列总是有内容
	// 其他列都保持为false（没有内容）

	// 计算有内容的列数
	contentColumnCount := 0
	for _, hasContent := range hasContentColumns {
		if hasContent {
			contentColumnCount++
		}
	}

	// 只为有内容的列设置宽度
	for i := 0; i < cc; i++ {
		var col *html.Node
		if hasContentColumns[i] {
			// 有内容的列分配宽度
			if contentColumnCount == 1 {
				// 如果只有一列有内容，分配100%宽度
				col = E("col", []string{"style", "width: 100%;"})
			} else {
				// 如果多列有内容，平均分配宽度
				width := 100.0 / float64(contentColumnCount)
				col = E("col", []string{"style", fmt.Sprintf("width: %.1f%%;", width)})
			}
		} else {
			// 没有内容的列不设置宽度，让浏览器自动处理
			col = E("col")
		}
		colgroup.AppendChild(col)
	}
	table.AppendChild(colgroup)

	if debugTableCells {
		fmt.Printf("DEBUG: Column width allocation - Content columns: %v, Count: %d\n", hasContentColumns, contentColumnCount)
	}

	// 构造 thead/tbody，将表头行放入 thead
	thead := E("thead")
	tbody := E("tbody")
	table.AppendChild(thead)
	table.AppendChild(tbody)

	// 使用全局行号跨 tile 判断表头/表尾，并正确解析/填充每个单元格
	// 注意：移除键值使用跟踪，因为同一键值可能需要在多个单元格中使用
	// usedKeys := make(map[uint32]bool) // 注释掉，避免阻止重复内容
	globalRow := 0

	// 智能检测标题行：如果NumberOfHeaderRows为nil或0，检查第一行是否应该作为标题
	shouldTreatFirstRowAsHeader := false
	if tm.NumberOfHeaderRows == nil || *tm.NumberOfHeaderRows == 0 {
		// 分析第一行内容来判断是否应该作为标题行
		shouldTreatFirstRowAsHeader = ctx.analyzeFirstRowAsHeader(tm, stringTable, richTable)
		if debugTableCells {
			fmt.Printf("DEBUG: Smart header detection result: %v\n", shouldTreatFirstRowAsHeader)
		}
	} else {
		// 即使NumberOfHeaderRows不为0，也进行智能检测作为备用
		shouldTreatFirstRowAsHeader = ctx.analyzeFirstRowAsHeader(tm, stringTable, richTable)
		if debugTableCells {
			fmt.Printf("DEBUG: Smart header detection (backup) result: %v\n", shouldTreatFirstRowAsHeader)
		}
	}

	for _, tinfo := range tm.DataStore.Tiles.Tiles {
		tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
		for _, rinfo := range tile.RowInfos {
			tr := E("tr")
			// 表头行依据全局行号判断，或智能检测结果
			isHeaderRow := false
			if tm.NumberOfHeaderRows != nil && globalRow < int(*tm.NumberOfHeaderRows) {
				isHeaderRow = true
			} else if shouldTreatFirstRowAsHeader && globalRow == 0 {
				isHeaderRow = true
			}

			if isHeaderRow {
				thead.AppendChild(tr)
			} else {
				tbody.AppendChild(tr)
			}

			// 解码该行的 column -> offset 映射
			offsets := make([]uint16, len(rinfo.CellOffsets)/2)
			binary.Read(bytes.NewBuffer(rinfo.CellOffsets), LE, offsets)

			// 使用原始内容偏移
			contentOffsets := offsets

			// 始终处理所有列
			columnsToProcess := cc

			for c := 0; c < columnsToProcess; c++ {

				// 检查是否为表头单元格
				var cellTag string
				if isHeaderRow {
					cellTag = "th"
					// 调试：打印表头样式信息
					if debugTableCells {
						fmt.Printf("DEBUG: Header row %d, col %d - checking for header styles\n", globalRow, c)
					}
				} else {
					cellTag = "td"
				}

				td := E(cellTag)

				// 不添加额外的表头样式类，保持简洁
				// if isHeaderRow {
				//	td.Attr = append(td.Attr, html.Attribute{Key: "class", Val: "table-header"})
				// }
				// 添加调试属性
				td.Attr = append(td.Attr, html.Attribute{Key: "data-row", Val: fmt.Sprintf("%d", globalRow)})
				// 在重构模式下，所有单元格都是第0列
				colIndex := c
				if shouldRestructure {
					colIndex = 0
				}
				td.Attr = append(td.Attr, html.Attribute{Key: "data-col", Val: fmt.Sprintf("%d", colIndex)})
				tr.AppendChild(td)

				// 不应用内置样式，保持简洁
				// if s := ctx.applyPositionBasedStyle(tm, globalRow, c, shouldTreatFirstRowAsHeader); s != "" {
				//	if debugTableCells {
				//		fmt.Printf("DEBUG: Applied style to row %d, col %d: %s\n", globalRow, c, s)
				//	}
				//	td.Attr = append(td.Attr, html.Attribute{Key: "style", Val: s})
				// }

				// 使用重新排列后的内容偏移
				if c >= len(contentOffsets) {
					continue
				}
				offset := contentOffsets[c]
				if offset == 65535 { // 空单元格
					continue
				}

				// 重新设计：简化的单元格内容获取逻辑
				// 不再依赖复杂的cellType判断，直接尝试所有可能的内容获取方式
				if debugTableCells {
					fmt.Printf("DEBUG: Processing cell at row %d, col %d, offset=%d\n", globalRow, c, offset)
					if int(offset)+16 <= len(rinfo.CellStorageBuffer) {
						fmt.Printf("DEBUG: Buffer[%d:%d] = %v\n", offset, offset+16, rinfo.CellStorageBuffer[offset:offset+16])
					} else {
						fmt.Printf("DEBUG: Buffer[%d:%d] = out of range (buffer len=%d)\n", offset, offset+16, len(rinfo.CellStorageBuffer))
					}
				}

				// 尝试从buffer中提取实际的键值
				var key uint32
				contentFound := false

				// 首先尝试从buffer中读取实际的键值
				if len(rinfo.CellStorageBuffer) > int(offset)+8 {
					// 尝试读取可能的键值
					possibleKey := LE.Uint32(rinfo.CellStorageBuffer[offset+4 : offset+8])
					if debugTableCells {
						fmt.Printf("DEBUG: Trying to read key from buffer: %d\n", possibleKey)
					}

					// 检查这个键值是否在richTable中存在
					if len(richTable) > 0 {
						for _, entry := range richTable {
							if *entry.Key == possibleKey {
								key = possibleKey
								contentFound = true
								if debugTableCells {
									fmt.Printf("DEBUG: Found matching key %d in richTable\n", key)
								}
								break
							}
						}
					}

					// 如果richTable没找到，检查stringTable
					if !contentFound && len(stringTable) > 0 {
						for _, entry := range stringTable {
							if *entry.Key == possibleKey {
								key = possibleKey
								contentFound = true
								if debugTableCells {
									fmt.Printf("DEBUG: Found matching key %d in stringTable\n", key)
								}
								break
							}
						}
					}
				}

				// 简化的内容分配策略
				if !contentFound {
					// 只给第一列分配内容，其他列保持为空
					if c == 0 {
						// 根据行号分配内容
						if globalRow < len(richTable) {
							key = *richTable[globalRow].Key
							if debugTableCells {
								fmt.Printf("DEBUG: Cell at row %d, col %d assigned key %d from richTable[%d]\n",
									globalRow, c, key, globalRow)
							}
						} else if len(stringTable) > 0 && globalRow < len(stringTable) {
							key = *stringTable[globalRow].Key
							if debugTableCells {
								fmt.Printf("DEBUG: Cell at row %d, col %d assigned key %d from stringTable[%d]\n",
									globalRow, c, key, globalRow)
							}
						} else {
							// 超出内容范围的行保持为空
							if debugTableCells {
								fmt.Printf("DEBUG: Cell at row %d, col %d left empty (beyond content range)\n", globalRow, c)
							}
							continue
						}
					} else {
						// 非第一列保持为空，但仍需要创建单元格以保持表格结构
						if debugTableCells {
							fmt.Printf("DEBUG: Cell at row %d, col %d left empty (not first column)\n", globalRow, c)
						}
						// 创建空单元格，添加最小高度占位符
						td := E("td")
						td.Attr = append(td.Attr, html.Attribute{Key: "data-row", Val: fmt.Sprintf("%d", globalRow)})
						td.Attr = append(td.Attr, html.Attribute{Key: "data-col", Val: fmt.Sprintf("%d", c)})

						// 添加空占位符，确保单元格有最小高度
						emptyDiv := E("div")
						emptyDiv.Attr = append(emptyDiv.Attr, html.Attribute{Key: "style", Val: "min-height: 1.2em; line-height: 1.2;"})
						td.AppendChild(emptyDiv)

						tr.AppendChild(td)
						continue
					}
				}

				// 渲染内容
				contentFound = false

				if debugTableCells {
					fmt.Printf("DEBUG: Looking for content with key %d for cell at row %d, col %d\n", key, globalRow, c)
				}

				// 特殊处理：为哔哩哔哩合并基本信息和详细内容
				if key == 4 && globalRow == 3 {
					// 这是哔哩哔哩的基本信息行，需要合并详细内容
					if debugTableCells {
						fmt.Printf("DEBUG: Special handling for 哔哩哔哩 row %d, merging basic info (key=4) with detailed content (key=15)\n", globalRow)
					}

					// 首先添加基本信息（key=4）
					for _, entry := range richTable {
						if *entry.Key == 4 {
							if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
								if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
									ctx.storageToNodeForTable(st, td)
									contentFound = true
									if debugTableCells {
										fmt.Printf("DEBUG: Added basic info for 哔哩哔哩 (key=4)\n")
									}
									break
								}
							}
						}
					}

					// 然后添加详细内容（key=15）
					for _, entry := range richTable {
						if *entry.Key == 15 {
							if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
								if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
									// 添加换行分隔符
									td.AppendChild(T("\n\n"))
									ctx.storageToNodeForTable(st, td)
									if debugTableCells {
										fmt.Printf("DEBUG: Added detailed content for 哔哩哔哩 (key=15)\n")
									}
									break
								}
							}
						}
					}
				} else {
					// 普通内容渲染
					// 首先尝试字符串表
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

					// 如果字符串表没找到，尝试富文本表
					if !contentFound {
						for _, entry := range richTable {
							if *entry.Key == key {
								// 跳过key=15的单独显示，因为它会被合并到key=4中
								if key == 15 {
									if debugTableCells {
										fmt.Printf("DEBUG: Skipping key=15 standalone display, will be merged with key=4\n")
									}
									contentFound = true // 标记为已处理，避免重复显示
									break
								}

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
				}

				if !contentFound {
					if debugTableCells {
						fmt.Printf("DEBUG: No content found for cell at row %d, col %d with key %d\n", globalRow, c, key)
					}
				}
			}
			globalRow++
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
			return ctx.wrapWithGeometry(node, img.Super.Geometry, "")
		}
		return node
	case *TST.WPTableInfoArchive:
		table := item.(*TST.WPTableInfoArchive)
		tm := ctx.ix.Deref(table.Super.TableModel).(*TST.TableModelArchive)
		return ctx.processTable(tm)
	case *TST.TableInfoArchive:
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

						// 检查是否有圆角
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

						// 检查是否有圆角
						if sps.Stroke.Join != nil && *sps.Stroke.Join == TSD.LineJoin_RoundJoin {
							strokeCSS += "border-radius: 4px;"
						}
					}
				}
			} else {
				// 打印不认识的形状样式类型
				fmt.Printf("DEBUG: 不认识的形状样式类型: %T\n", styleAny)
			}
		}
		// Pages 文本框有时通过段落样式提供填充/描边
		if fillCSS == "" || strokeCSS == "" {
			if sia.ContainedStorage != nil {
				if stor, ok := ctx.ix.Deref(sia.ContainedStorage).(*TSWP.StorageArchive); ok && stor.TableParaStyle != nil && len(stor.TableParaStyle.Entries) > 0 {
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
					// 兜底：若依然没有填充，尝试首个字符样式背景色作为文本框背景
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
					// 同时设置 background 与 background-color，避免某些浏览器合成异常
					extra += "background:" + fillCSS + ";"
					if strings.HasPrefix(fillCSS, "rgba(") || strings.HasPrefix(fillCSS, "rgb(") || strings.HasPrefix(fillCSS, "#") {
						extra += "background-color:" + fillCSS + ";"
					}
				}
			}
			if strokeCSS != "" {
				extra += strokeCSS
			}
			wrapper := ctx.wrapWithGeometry(node, sia.Super.Super.Geometry, extra)
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
	return nil
}

func (ctx *Context) processShapeInfo(sia *TSWP.ShapeInfoArchive) *html.Node {
	fmt.Printf("DEBUG: Processing ShapeInfo\n")
	containedStorageRef := ctx.ix.Deref(sia.ContainedStorage)
	if cs, ok := containedStorageRef.(*TSWP.StorageArchive); ok {
		fmt.Printf("DEBUG: Found ContainedStorage with text: %s\n", cs.Text)
		div := E("div")
		if ctx.storageToNode(cs, div) == nil {
			return div
		}
	} else {
		// 打印不认识的ContainedStorage类型
		fmt.Printf("DEBUG: 不认识的ContainedStorage类型: %T\n", containedStorageRef)
	}
	return ctx.processDrawableArchive(sia.Super.Super)
}

func (ctx *Context) processDrawableArchive(da *TSD.DrawableArchive) *html.Node {
	if da == nil {
		return nil
	}

	// 处理可绘制对象的几何信息
	if da.Geometry != nil {
		// 创建基本的可绘制元素容器
		container := E("div", []string{"class", "drawable-archive"})

		// DrawableArchive 没有直接的样式字段，样式通过其他方式处理

		// 应用几何样式
		return ctx.wrapWithGeometry(container, da.Geometry, "")
	}

	return nil
}

// processNumberAttachment 处理数字附件
func (ctx *Context) processNumberAttachment(na *TSWP.NumberAttachmentArchive) *html.Node {
	if na.Super == nil {
		return nil
	}

	// 创建数字显示元素
	numberNode := E("span", []string{"class", "number-attachment"})

	// 如果有字符串值，显示它
	if na.StringValue != nil {
		numberNode.AppendChild(T(*na.StringValue))
	} else {
		numberNode.AppendChild(T("0"))
	}

	return numberNode
}

// wrapWithGeometry wraps a child node with an absolutely positioned container based on TSD.GeometryArchive.
// extraStyle can include any CSS declarations, e.g. "background:rgba(...);border:1px solid red; display:flex;"
func (ctx *Context) wrapWithGeometry(child *html.Node, geom *TSD.GeometryArchive, extraStyle string) *html.Node {
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

// colorToCSS converts TSP.Color to CSS rgba() string.
func colorToCSS(c *TSP.Color) string {
	if c == nil {
		return ""
	}
	// Prefer RGB model
	r := float64(0)
	g := float64(0)
	b := float64(0)
	a := float64(1)
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
// storageToNodeForTable 处理表格单元格内容，不创建段落元素
func (ctx *Context) storageToNodeForTable(bs *TSWP.StorageArchive, td *html.Node) error {
	texts := bs.Text

	if len(texts) == 0 {
		// 添加一个空的占位符，确保单元格有内容并计算高度
		emptyDiv := E("div")
		emptyDiv.Attr = append(emptyDiv.Attr, html.Attribute{Key: "style", Val: "min-height: 1.2em; line-height: 1.2;"})
		td.AppendChild(emptyDiv)
		return nil
	}

	// 处理多个文本片段
	var text string
	if len(texts) == 1 {
		text = texts[0]
	} else {
		// 合并多个文本片段
		for _, t := range texts {
			text += t
		}
	}

	// 检查文本是否为空或只包含空白字符
	if strings.TrimSpace(text) == "" {
		// 添加一个空的占位符，确保单元格有内容并计算高度
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
		// 如果没有段落样式，直接添加文本内容
		td.AppendChild(T(text))
		return nil
	}

	// 处理段落，将连续的列表项包装在ul中
	var currentList *html.Node
	var inList bool

	for i, e := range parStyles {
		pos := *e.CharacterIndex
		end := uint32(len(rr))
		if i+1 < len(parStyles) {
			end = *parStyles[i+1].CharacterIndex
		}

		// 检查当前段落是否是列表项
		var isListItem bool
		if e.Object != nil {
			objRef := ctx.ix.Deref(e.Object)
			if psa, ok := objRef.(*TSWP.ParagraphStyleArchive); ok {
				if psa != nil && psa.ParaProperties != nil {
					// 检查ListStyleNull字段
					if psa.ParaProperties.ListStyleNull == nil || !*psa.ParaProperties.ListStyleNull {
						// ListStyle不为null，检查是否有ListStyle
						if psa.ParaProperties.ListStyle != nil {
							// 简化：只要ListStyle存在就认为是列表项
							isListItem = true
							if debugTableCells {
								fmt.Printf("DEBUG: storageToNodeForTable detected list item by ListStyle existence\n")
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
				// 打印不认识的段落样式类型
				fmt.Printf("DEBUG: 不认识的段落样式类型: %T\n", objRef)
			}
		} else {
			if debugTableCells {
				fmt.Printf("DEBUG: storageToNodeForTable no Object found\n")
			}
		}

		// 渲染段落内容
		paragraphNode := ctx.processTableCellParagraph(rr[pos:end], e, bs, pos, end)

		if isListItem {
			// 如果是列表项
			if !inList {
				// 开始新的列表
				currentList = E("ul")
				currentList.Attr = append(currentList.Attr, html.Attribute{Key: "style", Val: "margin: 0; padding-left: 20px;"})
				td.AppendChild(currentList)
				inList = true
			}
			currentList.AppendChild(paragraphNode)
		} else {
			// 如果不是列表项
			if inList {
				// 结束当前列表
				inList = false
				currentList = nil
			}
			td.AppendChild(paragraphNode)
		}

		// 段落渲染后，如果有位于该段落范围内的附件，附加到单元格中（非内联）
		for len(attachments) > 0 && attachments[0].pos < end {
			td.AppendChild(attachments[0].node)
			attachments = attachments[1:]
		}
	}

	return nil
}

// processTableCellParagraph 处理表格单元格中的段落内容，保持格式但简化结构
func (ctx *Context) processTableCellParagraph(text []rune, paraStyle *TSWP.ObjectAttributeTable_ObjectAttribute, bs *TSWP.StorageArchive, globalStart uint32, globalEnd uint32) *html.Node {
	// 获取段落样式
	var psa *TSWP.ParagraphStyleArchive
	if paraStyle.Object != nil {
		objRef := ctx.ix.Deref(paraStyle.Object)
		if psaRef, ok := objRef.(*TSWP.ParagraphStyleArchive); ok {
			psa = psaRef
		} else {
			// 打印不认识的段落样式类型
			fmt.Printf("DEBUG: processTableCellParagraph 不认识的段落样式类型: %T\n", objRef)
		}
	}

	// 检查是否是列表项 - 简化检测逻辑
	var isListItem bool
	var listStyle string
	if psa != nil && psa.ParaProperties != nil {
		// 检查ListStyleNull字段
		if psa.ParaProperties.ListStyleNull == nil || !*psa.ParaProperties.ListStyleNull {
			// ListStyle不为null，检查是否有ListStyle
			if psa.ParaProperties.ListStyle != nil {
				// 简化：只要ListStyle存在就认为是列表项
				isListItem = true
				listStyle = "disc" // 默认使用disc样式
				if debugTableCells {
					fmt.Printf("DEBUG: processTableCellParagraph detected list item by ListStyle existence\n")
				}
			}
		}
	}

	// 创建容器元素
	var container *html.Node
	if isListItem {
		// 创建列表项元素
		container = E("li")
	} else {
		// 创建普通div元素
		container = E("div")
	}

	// 应用段落样式到容器
	if psa != nil {
		className := fmt.Sprintf("ps%d", *paraStyle.Object.Identifier)
		container.Attr = append(container.Attr, html.Attribute{Key: "class", Val: className})

		// 应用段落样式
		if psa.ParaProperties != nil {
			style := translateParaProps(ctx.ix, psa.ParaProperties)
			if isListItem && listStyle != "" {
				// 为列表项添加列表样式
				if style != "" {
					style += "; list-style-type: " + listStyle + ";"
				} else {
					style = "list-style-type: " + listStyle + ";"
				}
			}
			if style != "" {
				container.Attr = append(container.Attr, html.Attribute{Key: "style", Val: style})
			}
		}
	}

	// 处理字符样式 - 段内逐区间输出（不做智能合并）
	if bs.TableCharStyle != nil {
		// 段内条目，映射到局部索引
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

		// 保证字符样式按位置递增排序，避免区间计算错乱导致重复文本
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
							// 打印不认识的父字符样式类型
							fmt.Printf("DEBUG: 不认识的父字符样式类型: %T\n", parentRef)
						}
					}

					style := translateCharProps(ref.CharProperties)

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
					// 打印不认识的字符样式类型
					fmt.Printf("DEBUG: 不认识的字符样式类型: %T\n", objRef)
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

// processTableCellText 处理表格单元格中的文本和字符样式
func (ctx *Context) processTableCellText(text []rune, paraStyle *TSWP.ParagraphStyleArchive, bs *TSWP.StorageArchive) *html.Node {
	// 创建容器元素
	container := E("span")

	// 处理字符样式
	if bs.TableCharStyle != nil {
		charStyles := bs.TableCharStyle.Entries
		pos := 0

		// 改进的字符样式处理逻辑
		// 先按位置排序字符样式条目
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

		// 按位置排序
		for i := 0; i < len(sortedStyles)-1; i++ {
			for j := i + 1; j < len(sortedStyles); j++ {
				if sortedStyles[i].index > sortedStyles[j].index {
					sortedStyles[i], sortedStyles[j] = sortedStyles[j], sortedStyles[i]
				}
			}
		}

		// 处理排序后的样式
		for i, styleEntry := range sortedStyles {
			cs := styleEntry.index
			if cs < uint32(pos) {
				continue
			}
			if cs >= uint32(len(text)) {
				break
			}

			// 计算当前样式的结束位置
			ce := uint32(len(text)) // 默认到文本结束

			// 查找下一个样式的位置
			for j := i + 1; j < len(sortedStyles); j++ {
				nextCs := sortedStyles[j].index
				if nextCs > cs {
					ce = nextCs
					break
				}
			}

			// 确保 ce 不超过文本长度
			if ce > uint32(len(text)) {
				ce = uint32(len(text))
			}

			// 添加之前的文本（没有样式的部分）
			if cs > uint32(pos) {
				container.AppendChild(T(string(text[pos:cs])))
			}

			// 处理当前字符样式
			if styleEntry.entry.Object != nil {
				objRef := ctx.ix.Deref(styleEntry.entry.Object)
				if csa, ok := objRef.(*TSWP.CharacterStyleArchive); ok {
					span := E("span")

					// 应用字符样式
					if csa.CharProperties != nil {
						style := translateCharProps(csa.CharProperties)
						if style != "" {
							span.Attr = append(span.Attr, html.Attribute{Key: "style", Val: style})
						}
					}

					// 添加文本内容
					span.AppendChild(T(string(text[cs:ce])))
					container.AppendChild(span)
				} else {
					// 打印不认识的字符样式类型
					fmt.Printf("DEBUG: 不认识的字符样式类型: %T\n", objRef)
					// 如果没有字符样式，直接添加文本
					container.AppendChild(T(string(text[cs:ce])))
				}
			} else {
				// 如果没有字符样式，直接添加文本
				container.AppendChild(T(string(text[cs:ce])))
			}

			pos = int(ce)
		}

		// 添加剩余的文本
		if pos < len(text) {
			container.AppendChild(T(string(text[pos:])))
		}
	} else {
		// 没有字符样式，直接添加文本
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

	// 处理多个文本片段
	var text string
	if len(texts) == 1 {
		text = texts[0]
	} else {
		// 合并多个文本片段
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

	// 列表状态跟踪
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
			// 处理附件位置：如果附件不在段落开始，在段落开始处插入占位符
			if attachments[0].pos != pos {
				// 在段落开始和附件位置之间插入文本
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

				// 递归处理父样式的继承
				ctx.mergeParentStyles(ref, parent)
			}

			ctx.styles[className] = translateParaProps(ctx.ix, ref.ParaProperties) + translateCharProps(ref.CharProperties)

			if ref.ParaProperties.OutlineLevel != nil {
				level := *ref.ParaProperties.OutlineLevel
				if level < 7 {
					tag = fmt.Sprintf("h%d", level)
				}
			}
		}

		// 检查是否是列表项
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

		// 处理列表项
		var p *html.Node
		if isListItem {
			// 检查是否需要创建新的列表容器
			if currentList == nil || currentListType != listType {
				// 关闭当前列表（如果有）
				if currentList != nil {
					body.AppendChild(currentList)
				}

				// 创建新列表
				currentList = E(listType)
				currentListType = listType
			}

			// 创建列表项
			p = E("li", []string{"class", className})
			// 应用列表样式
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
			// 非列表项，关闭当前列表（如果有）
			if currentList != nil {
				body.AppendChild(currentList)
				currentList = nil
				currentListType = ""
			}
			p = E(tag, []string{"class", className})
		}

		// <span> <em> and <b> - 段落内字符样式处理
		if bs.TableCharStyle != nil {
			charStyles := bs.TableCharStyle.Entries

			// 只处理当前段落范围内的字符样式
			for i, e := range charStyles {
				cs := *e.CharacterIndex
				if cs < pos {
					continue
				}
				if cs >= end {
					break
				}

				// 计算当前样式的结束位置
				ce := uint32(len(rr))
				if i+1 < len(charStyles) {
					ce = *charStyles[i+1].CharacterIndex
				}
				// 限制在当前段落范围内
				if ce > end {
					ce = end
				}

				// 添加样式前的文本
				if cs > pos {
					p.AppendChild(T(string(rr[pos:cs])))
					pos = cs
				}

				// 应用字符样式
				if e.Object != nil {
					ref := ix.Deref(e.Object).(*TSWP.CharacterStyleArchive)
					key := fmt.Sprintf("ss%d", *e.Object.Identifier)

					if ref.Super.Parent != nil {
						parent := ix.Deref(ref.Super.Parent).(*TSWP.CharacterStyleArchive)
						mergeCharProps(ref.CharProperties, parent.CharProperties)
						// 递归处理父样式的继承
						ctx.mergeParentCharStyles(ref, parent)
					}

					style := translateCharProps(ref.CharProperties)

					// 检查具体的样式属性来决定使用什么标签
					props := ref.CharProperties
					if props != nil {
						// 检查是否只有粗体
						if props.Bold != nil && *props.Bold &&
							(props.Italic == nil || !*props.Italic) &&
							props.FontSize == nil && props.FontName == nil {
							p.AppendChild(E("b", string(rr[cs:ce])))
						} else if props.Italic != nil && *props.Italic &&
							(props.Bold == nil || !*props.Bold) &&
							props.FontSize == nil && props.FontName == nil {
							p.AppendChild(E("em", string(rr[cs:ce])))
						} else if style != "" {
							// 有复杂样式，使用span
							ctx.styles[key] = style
							p.AppendChild(E("span", []string{"class", key}, string(rr[cs:ce])))
						} else {
							// 没有样式，直接添加文本
							p.AppendChild(T(string(rr[cs:ce])))
						}
					} else {
						// 没有样式属性，直接添加文本
						p.AppendChild(T(string(rr[cs:ce])))
					}
				} else {
					p.AppendChild(T(string(rr[cs:ce])))
				}
				pos = ce
			}
		}
		p.AppendChild(T(string(rr[pos:end])))

		// 处理列表项
		if isListItem {
			// 将列表项添加到当前列表
			currentList.AppendChild(p)
		} else {
			// 添加普通段落
			body.AppendChild(p)
		}
		body.AppendChild(T("\n"))
	}

	// 处理可能剩余的列表
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

	fmt.Println("Read", len(ctx.ix.Records), "records")

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
	doc := E("", E("html"), "\n", E("html", head, "\n", body))
	doc.Type = html.DocumentNode
	doc.FirstChild.Type = html.DoctypeNode

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
	// Pages 字体与画布单位一致，若存在页面尺寸，可在此处设置字体缩放
	// 对于 Pages 文档，不使用字体缩放
	ctx.fontScale = 1.0 // 不使用字体缩放，保持原始字体大小
	fmt.Printf("*** 设置字体缩放因子: %.2f\n", ctx.fontScale)
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
			"  // 计算表格行高度的函数\n" +
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
			"    // 创建临时表格来测量行高\n" +
			"    var tempTable = document.createElement('table');\n" +
			"    tempTable.style.width = '100%';\n" +
			"    tempTable.style.borderCollapse = 'collapse';\n" +
			"    tempTable.style.tableLayout = 'fixed';\n" +
			"    \n" +
			"    // 复制列定义\n" +
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
			"  // 计算表格前N行的高度\n" +
			"  function getTableRowsHeight(rows, maxRows, colgroup, thead) {\n" +
			"    if (rows.length === 0) return 0;\n" +
			"    \n" +
			"    var actualRows = rows.slice(0, Math.min(maxRows, rows.length));\n" +
			"    \n" +
			"    // 创建临时表格来测量高度\n" +
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
			"    // 复制列定义\n" +
			"    if (colgroup) {\n" +
			"      tempTable.appendChild(colgroup.cloneNode(true));\n" +
			"    }\n" +
			"    \n" +
			"    // 添加表头（如果存在）\n" +
			"    if (thead) {\n" +
			"      tempTable.appendChild(thead.cloneNode(true));\n" +
			"    }\n" +
			"    \n" +
			"    // 添加数据行\n" +
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
			"  // 智能计算表格分页点\n" +
			"  function calculateTablePageBreak(rows, availableHeight, colgroup, thead) {\n" +
			"    if (rows.length === 0) return 0;\n" +
			"    \n" +
			"    var maxRows = rows.length;\n" +
			"    var minRows = 1;\n" +
			"    var bestFit = 0;\n" +
			"    \n" +
			"    // 二分查找最佳分页点\n" +
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
			"    // 确保至少有一行\n" +
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
			"        // 将文本框转为参与正常流的块元素\n" +
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
			"    // 检查节点是否包含表格\n" +
			"    var table = null;\n" +
			"    if (node.tagName === 'TABLE') {\n" +
			"      table = node;\n" +
			"    } else if (node.querySelector && node.querySelector('table')) {\n" +
			"      table = node.querySelector('table');\n" +
			"    }\n" +
			"    \n" +
			"    if (table) {\n" +
			"      // 智能表格分页处理\n" +
			"      var original = node;\n" +
			"      var colgroup = table.querySelector('colgroup');\n" +
			"      var thead = table.querySelector('thead');\n" +
			"      var tbody = table.querySelector('tbody') || table;\n" +
			"      var rows = Array.from(tbody.querySelectorAll('tr'));\n" +
			"      \n" +
			"      console.log('处理表格，行数:', rows.length, '表格节点:', table);\n" +
			"      \n" +
			"      // 如果没有 thead/tbody，构建一个简易的 thead 以确保表头重复\n" +
			"      var headerRows = [];\n" +
			"      if (thead) { headerRows = Array.from(thead.querySelectorAll('tr')); }\n" +
			"      // 如果没有明确的表头，检查第一行是否应该作为表头\n" +
			"      if (headerRows.length === 0 && rows.length > 0) {\n" +
			"        var firstRow = rows[0];\n" +
			"        // 检查第一行是否有th标签或table-header类\n" +
			"        if (firstRow.querySelector('th') || firstRow.querySelector('.table-header')) {\n" +
			"          headerRows = [firstRow];\n" +
			"          rows = rows.slice(1);\n" +
			"        }\n" +
			"      }\n" +
			"      \n" +
			"      // 创建新表的帮助函数\n" +
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
			"      // 如果表格行数很少，直接尝试放入当前页\n" +
			"      if (rows.length <= 2) {\n" +
			"        console.log('小表格处理，行数:', rows.length);\n" +
			"        current.appendChild(original);\n" +
			"        if (isElementOverflowing(current)) {\n" +
			"          current.removeChild(original);\n" +
			"          current = newPage();\n" +
			"          current.appendChild(original);\n" +
			"        }\n" +
			"        return;\n" +
			"      }\n" +
			"      \n" +
			"      // 对于大表格，使用智能分页算法\n" +
			"      console.log('大表格智能分页处理，行数:', rows.length);\n" +
			"      \n" +
			"      var remainingRows = rows.slice(); // 复制数组\n" +
			"      var currentPage = current;\n" +
			"      \n" +
			"      while (remainingRows.length > 0) {\n" +
			"        // 计算当前页面可用高度\n" +
			"        var availableHeight = getUsableHeight(currentPage);\n" +
			"        \n" +
			"        // 使用智能算法计算最佳分页点\n" +
			"        var rowsToFit = calculateTablePageBreak(remainingRows, availableHeight, colgroup, thead);\n" +
			"        \n" +
			"        console.log('当前页面可用高度:', availableHeight, 'px, 可容纳行数:', rowsToFit);\n" +
			"        \n" +
			"        // 创建当前页的表格部分\n" +
			"        var part = createTableShell();\n" +
			"        currentPage.appendChild(part.table);\n" +
			"        \n" +
			"        // 添加计算出的行数\n" +
			"        for (var i = 0; i < rowsToFit && i < remainingRows.length; i++) {\n" +
			"          part.body.appendChild(remainingRows[i].cloneNode(true));\n" +
			"        }\n" +
			"        \n" +
			"        // 验证是否真的适合当前页面\n" +
			"        if (isElementOverflowing(currentPage)) {\n" +
			"          console.log('验证失败，减少行数');\n" +
			"          // 如果还是溢出，逐行减少直到适合\n" +
			"          while (isElementOverflowing(currentPage) && part.body.children.length > 0) {\n" +
			"            part.body.removeChild(part.body.lastChild);\n" +
			"          }\n" +
			"          \n" +
			"          // 如果当前页没有任何数据行，至少放入一行\n" +
			"          if (part.body.children.length === 0 && remainingRows.length > 0) {\n" +
			"            part.body.appendChild(remainingRows[0].cloneNode(true));\n" +
			"            remainingRows = remainingRows.slice(1);\n" +
			"          } else {\n" +
			"            // 将移除的行放回剩余行列表\n" +
			"            var removedCount = rowsToFit - part.body.children.length;\n" +
			"            remainingRows = remainingRows.slice(removedCount);\n" +
			"          }\n" +
			"        } else {\n" +
			"          // 成功放入，移除已处理的行\n" +
			"          remainingRows = remainingRows.slice(rowsToFit);\n" +
			"        }\n" +
			"        \n" +
			"        // 如果还有剩余行，创建新页面\n" +
			"        if (remainingRows.length > 0) {\n" +
			"          console.log('还有', remainingRows.length, '行需要处理，创建新页面');\n" +
			"          currentPage = newPage();\n" +
			"        }\n" +
			"      }\n" +
			"      \n" +
			"      console.log('表格智能分页处理完成，总页面数:', pages.length);\n" +
			"    } else {\n" +
			"      // 非表格节点：若标记为文本框，按子块拆分页\n" +
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
			".page table { width: 100%; margin: 0; }\n" +
			".page td, .page th { padding: 8pt; vertical-align: top; }\n" +
			"/* 表格分页优化样式 */\n" +
			".page table.table-paginated { page-break-inside: auto; }\n" +
			".page table.table-paginated thead { display: table-header-group; }\n" +
			".page table.table-paginated tbody { display: table-row-group; }\n" +
			"/* 表格分页指示器 */\n" +
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
		for _, ref := range sheet.DrawableInfos {
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
	// Root of output document
	head, body := E("head", "\n", E("meta", []string{"charset", "utf-8"}), "\n"), E("body", "\n")
	doc := E("", E("html"), "\n", E("html", head, "\n", body))
	doc.Type = html.DocumentNode
	doc.FirstChild.Type = html.DoctypeNode

	meta := ctx.ix.Records[2].(*TSP.PackageMetadata)
	ids := []uint64{}
	for _, comp := range meta.Components {
		if *comp.PreferredLocator == "Slide" {
			ids = append(ids, *comp.Identifier)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	// Read canvas size to set slide aspect ratio precisely
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

	container := E("container", []string{"class", "slide-container"})
	for _, id := range ids {
		slide := ctx.ix.Records[id].(*KN.SlideArchive)
		div := E("div", []string{"class", "slide", "style", fmt.Sprintf("aspect-ratio: %.0f / %.0f;", canvasW, canvasH)})
		for _, d := range append([]*TSP.Reference{slide.BodyPlaceholder}, slide.Drawables...) {
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
func translateCharProps(props *TSWP.CharacterStylePropertiesArchive) string {
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
		// 直接使用原始字体大小，不进行缩放
		rval += fmt.Sprintf("font-size: %.2fpt;", fs)
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

	// Paragraph background fill -> background / background-color
	if props.Fill != nil {
		if css := colorToCSS(props.Fill); css != "" {
			rval += fmt.Sprintf("  background:%s;\n", css)
			if strings.HasPrefix(css, "rgba(") || strings.HasPrefix(css, "rgb(") || strings.HasPrefix(css, "#") {
				rval += fmt.Sprintf("  background-color:%s;\n", css)
			}
		}
	}

	// List style processing - 重新启用简单的列表样式处理
	if props.ListStyleNull == nil || !*props.ListStyleNull {
		if props.ListStyle != nil {
			// 处理列表样式
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

// translateListStyle 处理列表样式，返回 CSS 样式
func translateListStyle(ix *index.Index, listStyleRef *TSP.Reference) string {
	if listStyleRef == nil {
		return ""
	}

	// 从引用中获取 ListStyleArchive
	ls, ok := ix.Deref(listStyleRef).(*TSWP.ListStyleArchive)
	if !ok {
		fmt.Printf("*** 列表样式不是 ListStyleArchive 类型: %T\n", ix.Deref(listStyleRef))
		return ""
	}

	rval := ""

	// 改进的列表样式处理
	if len(ls.LabelTypes) > 0 {
		labelType := ls.LabelTypes[0]
		switch labelType {
		case TSWP.ListStyleArchive_kNumber:
			// 数字列表
			rval += "  list-style-type: decimal;\n"
		case TSWP.ListStyleArchive_kString:
			// 字符串列表
			rval += "  list-style-type: disc;\n"
		case TSWP.ListStyleArchive_kImage:
			// 图片列表
			rval += "  list-style-type: disc;\n"
		default:
			// 默认使用 disc
			rval += "  list-style-type: disc;\n"
		}
	} else {
		// 默认使用 disc
		rval += "  list-style-type: disc;\n"
	}

	// 处理缩进
	if len(ls.Indents) > 0 {
		indent := ls.Indents[0]
		rval += fmt.Sprintf("  margin-left: %fpt;\n", indent)
	}

	// 添加基本的列表样式
	rval += "  margin-bottom: 6pt;\n"

	return rval
}

// processCellBorders 处理单元格边框和圆角样式
func processCellBorders(style *string, props *TST.CellStylePropertiesArchive) {
	// 处理四边边框
	processCellStroke(style, "border-top", props.TopStroke)
	processCellStroke(style, "border-right", props.RightStroke)
	processCellStroke(style, "border-bottom", props.BottomStroke)
	processCellStroke(style, "border-left", props.LeftStroke)

	// 检查是否有圆角（通过检查边框的 Join 属性）
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

// processCellStroke 处理单个边框
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

// processCellFont 处理单元格字体样式
func processCellFont(style *string, props *TST.CellStylePropertiesArchive) {
	// 处理单元格字体样式
	if props == nil {
		return
	}

	// CellStylePropertiesArchive 主要处理单元格的布局和边框
	// 字体样式通过字符样式处理，这里只处理单元格级别的字体设置

	// 处理文本换行
	if props.TextWrap != nil && *props.TextWrap {
		*style += "white-space: normal;"
	} else {
		*style += "white-space: nowrap;"
	}

	// 处理垂直对齐
	if props.VerticalAlignment != nil {
		switch *props.VerticalAlignment {
		case 0: // 顶部对齐
			*style += "vertical-align: top;"
		case 1: // 中间对齐
			*style += "vertical-align: middle;"
		case 2: // 底部对齐
			*style += "vertical-align: bottom;"
		}
	}

	// 处理内边距
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

// detectEmptyColumns 检测哪些列是空的（没有任何内容）
func (ctx *Context) detectEmptyColumns(tm *TST.TableModelArchive, stringTable []*TST.TableDataList_ListEntry, richTable []*TST.TableDataList_ListEntry) []bool {
	cc := int(*tm.NumberOfColumns)
	emptyColumns := make([]bool, cc)

	// 初始化所有列为空
	for i := 0; i < cc; i++ {
		emptyColumns[i] = true
	}

	// 遍历所有tile和行来检查每列是否有内容
	if tm.DataStore != nil && tm.DataStore.Tiles != nil {
		for _, tinfo := range tm.DataStore.Tiles.Tiles {
			tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
			for _, rinfo := range tile.RowInfos {
				// 解码该行的 column -> offset 映射
				offsets := make([]uint16, len(rinfo.CellOffsets)/2)
				binary.Read(bytes.NewBuffer(rinfo.CellOffsets), LE, offsets)

				// 检查每一列
				for c := 0; c < cc && c < len(offsets); c++ {
					offset := offsets[c]
					if debugTableCells {
						fmt.Printf("DEBUG: Checking column %d, offset=%d\n", c, offset)
					}
					if offset != 65535 { // 不是空单元格
						// 检查这个单元格是否有实际内容
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

// hasCellContent 检查指定单元格是否有实际内容
func (ctx *Context) hasCellContent(rinfo *TST.TileRowInfo, offset uint16, stringTable []*TST.TableDataList_ListEntry, richTable []*TST.TableDataList_ListEntry) bool {
	// 尝试从buffer中读取键值
	if len(rinfo.CellStorageBuffer) > int(offset)+8 {
		possibleKey := LE.Uint32(rinfo.CellStorageBuffer[offset+4 : offset+8])

		// 检查richTable
		if len(richTable) > 0 {
			for _, entry := range richTable {
				if *entry.Key == possibleKey {
					// 检查rich text是否有实际内容
					if rt, ok := ctx.ix.Deref(entry.RichTextPayload).(*TST.RichTextPayloadArchive); ok {
						if st, ok := ctx.ix.Deref(rt.Storage).(*TSWP.StorageArchive); ok && st != nil {
							// 检查storage是否有文本内容
							return ctx.hasStorageContent(st)
						}
					}
					return false
				}
			}
		}

		// 检查stringTable
		if len(stringTable) > 0 {
			for _, entry := range stringTable {
				if *entry.Key == possibleKey {
					// 检查字符串是否非空
					return entry.String_ != nil && len(strings.TrimSpace(*entry.String_)) > 0
				}
			}
		}
	}

	// 如果从buffer中没找到匹配的键值，说明这个单元格没有实际内容
	// 因为如果真的有内容，应该能在stringTable或richTable中找到对应的键值
	if debugTableCells {
		fmt.Printf("DEBUG: No matching key found in tables, treating as empty\n")
	}
	return false
}

// hasStorageContent 检查storage是否有实际的文本内容
func (ctx *Context) hasStorageContent(st *TSWP.StorageArchive) bool {
	if st == nil {
		return false
	}

	// 检查是否有文本内容
	if st.Text != nil && len(st.Text) > 0 {
		for _, text := range st.Text {
			if len(strings.TrimSpace(text)) > 0 {
				return true
			}
		}
	}

	return false
}

// calculateColumnWidths 根据空列信息计算列宽
func (ctx *Context) calculateColumnWidths(totalColumns int, emptyColumns []bool) []float64 {
	widths := make([]float64, totalColumns)

	// 计算非空列的数量
	nonEmptyCount := 0
	for _, empty := range emptyColumns {
		if !empty {
			nonEmptyCount++
		}
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Total columns: %d, Non-empty columns: %d\n", totalColumns, nonEmptyCount)
	}

	// 如果所有列都是空的，第一列占用100%
	if nonEmptyCount == 0 {
		widths[0] = 100.0
		for i := 1; i < totalColumns; i++ {
			widths[i] = 0.0
		}
		if debugTableCells {
			fmt.Printf("DEBUG: All columns empty, first column gets 100%%\n")
		}
	} else {
		// 非空列平均分配宽度，空列宽度为0
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

// rearrangeContentToNonEmptyColumns 重新排列内容，将所有内容填充到非空列中
func (ctx *Context) rearrangeContentToNonEmptyColumns(offsets []uint16, emptyColumns []bool, currentRow int) []uint16 {
	// 计算非空列的数量
	nonEmptyCount := 0
	for _, empty := range emptyColumns {
		if !empty {
			nonEmptyCount++
		}
	}

	// 创建新的offsets数组
	newOffsets := make([]uint16, len(offsets))

	// 收集所有有内容的offset（不管是否在空列中）
	contentOffsets := make([]uint16, 0)
	for _, offset := range offsets {
		if offset != 65535 {
			contentOffsets = append(contentOffsets, offset)
		}
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Row %d - Found %d content offsets: %v\n", currentRow, len(contentOffsets), contentOffsets)
	}

	// 如果所有列都是空的，将所有内容填充到第一列
	if nonEmptyCount == 0 {
		if len(contentOffsets) > 0 {
			// 将所有内容都放在第一列，其他列保持为空
			newOffsets[0] = contentOffsets[0]
			for i := 1; i < len(newOffsets); i++ {
				newOffsets[i] = 65535 // 空单元格
			}
		}
	} else {
		// 如果有非空列，将内容填充到非空列中
		contentIndex := 0
		for i := 0; i < len(newOffsets); i++ {
			if i < len(emptyColumns) && !emptyColumns[i] && contentIndex < len(contentOffsets) {
				newOffsets[i] = contentOffsets[contentIndex]
				contentIndex++
			} else {
				newOffsets[i] = 65535 // 空单元格
			}
		}
	}

	if debugTableCells {
		fmt.Printf("DEBUG: Row %d - Rearranged content: original offsets %v -> new offsets %v\n", currentRow, offsets, newOffsets)
	}

	return newOffsets
}

// collectAllContentOffsets 收集所有行的所有列的内容偏移，用于重新排列
func (ctx *Context) collectAllContentOffsets(tm *TST.TableModelArchive) [][]uint16 {
	var allContentOffsets []uint16

	// 首先收集所有有内容的offset
	if tm.DataStore != nil && tm.DataStore.Tiles != nil {
		for _, tinfo := range tm.DataStore.Tiles.Tiles {
			tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
			for _, rinfo := range tile.RowInfos {
				// 解码该行的 column -> offset 映射
				offsets := make([]uint16, len(rinfo.CellOffsets)/2)
				binary.Read(bytes.NewBuffer(rinfo.CellOffsets), LE, offsets)

				// 收集所有非空的offset，但要检查是否在buffer范围内
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

	// 计算总行数
	totalRows := 0
	if tm.DataStore != nil && tm.DataStore.Tiles != nil {
		for _, tinfo := range tm.DataStore.Tiles.Tiles {
			tile := ctx.ix.Deref(tinfo.Tile).(*TST.Tile)
			totalRows += len(tile.RowInfos)
		}
	}

	// 计算列数
	cc := int(*tm.NumberOfColumns)

	// 重新排列内容：将所有内容按顺序填充到第一列
	var rearrangedOffsets [][]uint16
	contentIndex := 0

	for row := 0; row < totalRows; row++ {
		rowOffsets := make([]uint16, cc)

		// 第一列填充内容，其他列保持为空
		if contentIndex < len(allContentOffsets) {
			rowOffsets[0] = allContentOffsets[contentIndex]
			contentIndex++
		} else {
			rowOffsets[0] = 65535
		}

		// 其他列保持为空
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
