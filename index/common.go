package index

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/orcastor/iwork-converter/proto/TP"
	"github.com/orcastor/iwork-converter/proto/TSA"
	"github.com/orcastor/iwork-converter/proto/TSCE"
	"github.com/orcastor/iwork-converter/proto/TSCH"
	TSCH_PreUFF "github.com/orcastor/iwork-converter/proto/TSCH/PreUFF"
	"github.com/orcastor/iwork-converter/proto/TSD"
	"github.com/orcastor/iwork-converter/proto/TSK"
	"github.com/orcastor/iwork-converter/proto/TSP"
	"github.com/orcastor/iwork-converter/proto/TSS"
	"github.com/orcastor/iwork-converter/proto/TST"
	"github.com/orcastor/iwork-converter/proto/TSWP"
)

func decodeCommon(typ uint32, payload []byte) (interface{}, error) {
	switch typ {

	case 0:
		var value = &TSP.Reference{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10016:
		var value = &TP.SectionTemplateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10024:
		var value = &TP.PageTemplateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10166:
		var value = &TP.TopicNumberHintsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11000:
		var value = &TSP.PasteboardObject{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11006:
		var value = &TSP.PackageMetadata{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11007:
		var value = &TSP.PasteboardMetadata{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11008:
		var value = &TSP.ObjectContainer{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11011:
		var value = &TSP.Reference{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11014:
		var value = &TSP.Reference{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11015:
		var value = &TSP.Reference{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 12050:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 200:
		var value = &TSK.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2001:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2002:
		var value = &TSWP.SelectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2003:
		var value = &TSWP.DrawableAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2004:
		var value = &TSWP.TextualAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2005:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2006:
		var value = &TSWP.UIGraphicalAttachment{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2007:
		var value = &TSWP.TextualAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2008:
		var value = &TSWP.FootnoteReferenceAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2009:
		var value = &TSWP.TextualAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 201:
		var value = &TSK.LocalCommandHistory{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2010:
		var value = &TSWP.TSWPTOCPageNumberAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2011:
		var value = &TSWP.ShapeInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2013:
		var value = &TSWP.HighlightArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2014:
		var value = &TSWP.CommentInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 202:
		var value = &TSK.CommandGroupArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2021:
		var value = &TSWP.CharacterStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2022:
		var value = &TSWP.ParagraphStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2023:
		var value = &TSWP.ListStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2024:
		var value = &TSWP.ColumnStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2025:
		var value = &TSWP.ShapeStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2026:
		var value = &TSWP.TOCEntryStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 203:
		var value = &TSK.CommandContainerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2031:
		var value = &TSWP.PlaceholderSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2032:
		var value = &TSWP.HyperlinkFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2033:
		var value = &TSWP.FilenameSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2034:
		var value = &TSWP.DateTimeSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2035:
		var value = &TSWP.BookmarkFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2036:
		var value = &TSWP.MergeSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2037:
		var value = &TSWP.CitationRecordArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2038:
		var value = &TSWP.CitationSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2039:
		var value = &TSWP.UnsupportedHyperlinkFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 204:
		var value = &TSK.CommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2040:
		var value = &TSWP.BibliographySmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2041:
		var value = &TSWP.TOCSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2042:
		var value = &TSWP.RubyFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2043:
		var value = &TSWP.NumberAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 205:
		var value = &TSK.TreeNode{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2050:
		var value = &TSWP.TextStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2051:
		var value = &TSWP.TOCSettingsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2052:
		var value = &TSWP.TOCEntryInstanceArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 206:
		var value = &TSK.ProgressiveCommandGroupArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2060:
		var value = &TSWP.ChangeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2061:
		var value = &TSK.DeprecatedChangeAuthorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2062:
		var value = &TSWP.ChangeSessionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 208:
		var value = &TSK.CommandBehaviorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 209:
		var value = &TSK.CommandSelectionBehaviorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 210:
		var value = &TSK.ViewStateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2101:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2102:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2104:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2105:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2107:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2108:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 211:
		var value = &TSK.DocumentSupportArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2113:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2114:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2115:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2116:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2117:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2118:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2119:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 212:
		var value = &TSK.AnnotationAuthorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2120:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2121:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2122:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 213:
		var value = &TSK.AnnotationAuthorStorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 214:
		var value = &TSK.AnnotationAuthorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 215:
		var value = &TSK.AnnotationAuthorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 219:
		var value = &TSP.Reference{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2206:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2207:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 222:
		var value = &TSP.Reference{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2231:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2232:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2240:
		var value = &TSWP.TOCInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2241:
		var value = &TSWP.TOCAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2242:
		var value = &TSWP.TOCLayoutHintArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2400:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2401:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2402:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2403:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2404:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2405:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2406:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2411:
		var value = &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 258:
		var value = &TSK.CommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3002:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3003:
		var value = &TSD.ContainerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3004:
		var value = &TSD.ShapeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3005:
		var value = &TSD.ImageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3006:
		var value = &TSD.MaskArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3007:
		var value = &TSD.MovieArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3008:
		var value = &TSD.GroupArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3009:
		var value = &TSD.ConnectionLineArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3015:
		var value = &TSD.ShapeStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3016:
		var value = &TSD.MediaStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3020:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3021:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3022:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3023:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3024:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3025:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3026:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3027:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3028:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3030:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3031:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3032:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3033:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3034:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3035:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3036:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3037:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3038:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3039:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3040:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3041:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3042:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3043:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3045:
		var value = &TSD.CanvasSelectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3046:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3047:
		var value = &TSD.GuideStorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3048:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3049:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3050:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3051:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3052:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3053:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3054:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3055:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3056:
		var value = &TSD.CommentStorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3057:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3058:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3059:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3060:
		var value = &TSD.CommentStorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3091:
		var value = &TSP.Reference{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3097:
		var value = &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 400:
		var value = &TSS.StyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4000:
		var value = &TSCE.FormulaArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4001:
		var value = &TSCE.FormulaArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4002:
		var value = &TSCE.RangeReferenceArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4003:
		var value = &TSCE.FormulaArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4004:
		var value = &TSCE.RangeReferenceArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4005:
		var value = &TSCE.RangeReferenceArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4008:
		var value = &TSCE.CellValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4009:
		var value = &TSCE.NumberCellValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 401:
		var value = &TSS.StylesheetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4010:
		var value = &TSCE.StringCellValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4011:
		var value = &TSCE.BooleanCellValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4012:
		var value = &TSCE.DateCellValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4013:
		var value = &TSCE.ErrorCellValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 402:
		var value = &TSS.ThemeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 410:
		var value = &TSS.ApplyThemeCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 411:
		var value = &TSS.ApplyThemeChildCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 412:
		var value = &TSS.StyleUpdatePropertyMapCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 413:
		var value = &TSS.ThemeReplacePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 414:
		var value = &TSS.ThemeAddStylePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 415:
		var value = &TSS.ThemeRemoveStylePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 416:
		var value = &TSS.ThemeReplaceColorPresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 417:
		var value = &TSS.ThemeMovePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 418:
		var value = &TSS.ThemeReplaceColorPresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5000:
		var value = &TSCH_PreUFF.ChartInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5002:
		var value = &TSCH_PreUFF.ChartGridArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5004:
		var value = &TSCH.ChartMediatorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5010:
		var value = &TSCH.ChartStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5011:
		var value = &TSCH.ChartSeriesStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5012:
		var value = &TSCH.ChartAxisStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5013:
		var value = &TSCH.LegendStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5014:
		var value = &TSCH.ChartNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5015:
		var value = &TSCH.ChartSeriesNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5016:
		var value = &TSCH.ChartAxisNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5017:
		var value = &TSCH.LegendNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5020:
		var value = &TSCH.ChartStylePreset{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5021:
		var value = &TSCH.ChartDrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5022:
		var value = &TSCH.ChartStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5023:
		var value = &TSCH.ChartNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5024:
		var value = &TSCH.LegendStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5025:
		var value = &TSCH.LegendNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5026:
		var value = &TSCH.ChartAxisStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5027:
		var value = &TSCH.ChartAxisNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5028:
		var value = &TSCH.ChartSeriesStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5029:
		var value = &TSCH.ChartSeriesNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5030:
		var value = &TSP.Reference{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5103:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5104:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5105:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5107:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5108:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5109:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5110:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5113:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5114:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5115:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5116:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5117:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5118:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5119:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5120:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5121:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5122:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5123:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5124:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5125:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5126:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5127:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5129:
		var value = &TSCH.StylePasteboardDataArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5130:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5131:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5132:
		var value = &TSCH.ChartArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 600:
		var value = &TSA.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6000:
		var value = &TST.TableInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6001:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6002:
		var value = &TST.Tile{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6003:
		var value = &TST.TableStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6004:
		var value = &TST.CellStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6005:
		var value = &TST.TableDataList{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6006:
		var value = &TST.HeaderStorageBucket{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6007:
		var value = &TST.WPTableInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6008:
		var value = &TST.TableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6009:
		var value = &TST.TableStrokePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 601:
		var value = &TSA.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6010:
		var value = &TST.ConditionalStyleSetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 602:
		var value = &TSA.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6100:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6101:
		var value = &TST.TableInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6102:
		var value = &TST.TableStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6103:
		var value = &TST.CellStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6104:
		var value = &TST.TableDataList{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6105:
		var value = &TST.CellMapArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6106:
		var value = &TST.CellListArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6107:
		var value = &TST.MergeRegionMapArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6108:
		var value = &TST.FormulaArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6109:
		var value = &TST.ExpressionNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6110:
		var value = &TST.BooleanNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6111:
		var value = &TST.NumberNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6112:
		var value = &TST.StringNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6113:
		var value = &TST.ArrayNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6114:
		var value = &TST.ListNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6115:
		var value = &TST.OperatorNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6116:
		var value = &TST.FunctionNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6117:
		var value = &TST.DateNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6118:
		var value = &TST.ReferenceNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6119:
		var value = &TST.DurationNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6120:
		var value = &TST.ConditionalStyleSetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6121:
		var value = &TST.FilterSetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6122:
		var value = &TST.UniqueIndexArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6123:
		var value = &TST.HiddenStatesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6124:
		var value = &TST.ExpandCollapseStateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6125:
		var value = &TST.TokenAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6126:
		var value = &TST.IdentifierNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6127:
		var value = &TST.PostfixOperatorNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6128:
		var value = &TST.PrefixOperatorNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6129:
		var value = &TST.FunctionEndNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6130:
		var value = &TST.ArgumentPlaceholderNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6131:
		var value = &TST.EmptyExpressionNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6132:
		var value = &TST.LetNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6134:
		var value = &TST.InNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6135:
		var value = &TST.VariableNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6136:
		var value = &TST.LayoutHintArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6137:
		var value = &TST.CompletionTokenAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6138:
		var value = &TST.HiddenStateFormulaOwnerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6139:
		var value = &TST.FormulaStoreArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6141:
		var value = &TST.MergeOperationArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6142:
		var value = &TST.MergeOwnerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6143:
		var value = &TST.PencilAnnotationArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6144:
		var value = &TST.MergeRegionMapArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6145:
		var value = &TST.PencilAnnotationOwnerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6146:
		var value = &TST.AccumulatorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6147:
		var value = &TST.GroupColumnArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6148:
		var value = &TST.GroupColumnListArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6179:
		var value = &TST.FormulaEqualsTokenAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6181:
		var value = &TST.TokenAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6182:
		var value = &TST.ExpressionNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6183:
		var value = &TST.BooleanNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6184:
		var value = &TST.NumberNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6185:
		var value = &TST.StringNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6186:
		var value = &TST.ArrayNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6187:
		var value = &TST.ListNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6188:
		var value = &TST.OperatorNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6189:
		var value = &TST.FunctionNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6190:
		var value = &TST.DateNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6191:
		var value = &TST.ReferenceNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6192:
		var value = &TST.DurationNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6193:
		var value = &TST.ArgumentPlaceholderNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6194:
		var value = &TST.PostfixOperatorNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6195:
		var value = &TST.PrefixOperatorNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6196:
		var value = &TST.FunctionEndNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6197:
		var value = &TST.EmptyExpressionNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6198:
		var value = &TST.LayoutHintArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6199:
		var value = &TST.CompletionTokenAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6200:
		var value = &TST.TableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6201:
		var value = &TST.TableDataList{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6202:
		var value = &TST.ColumnAggregateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6203:
		var value = &TST.ColumnAggregateListArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6204:
		var value = &TST.HiddenStateFormulaOwnerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6205:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6206:
		var value = &TST.PopUpMenuModel{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6207:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6208:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6209:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6210:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6211:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6212:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6213:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6214:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6215:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6216:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6217:
		var value = &TST.TableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6218:
		var value = &TST.RichTextPayloadArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6219:
		var value = &TST.TableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6220:
		var value = &TST.FilterSetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6221:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6222:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6223:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6224:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6225:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6226:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6227:
		var value = &TST.TableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6228:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6229:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6231:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6232:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6233:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6234:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6235:
		var value = &TST.IdentifierNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6236:
		var value = &TST.TableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6237:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6238:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6239:
		var value = &TST.TableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6240:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6241:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6242:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6244:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6245:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6246:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6247:
		var value = &TST.TableStyleNetworkArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6248:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6249:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6250:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6251:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6252:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6253:
		var value = &TST.TableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6254:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6255:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6256:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6267:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6305:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6306:
		var value = &TST.TableStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6316:
		var value = &TST.DataStore{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6317:
		var value = &TST.TableStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6318:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6365:
		var value = &TST.ControlCellSelectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6366:
		var value = &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6372:
		var value = &TST.TableStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6373:
		var value = &TST.TableStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6383:
		var value = &TST.TableStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	default:

		return nil, errors.New(fmt.Sprintf("Unknown type %d", typ))

	}
}
