package index

import (
	"fmt"

	"github.com/orcastor/iwork-converter/proto/TSA"
	"github.com/orcastor/iwork-converter/proto/TSCE"
	"github.com/orcastor/iwork-converter/proto/TSCH"
	"github.com/orcastor/iwork-converter/proto/TSCH/PreUFF"
	"github.com/orcastor/iwork-converter/proto/TSD"
	"github.com/orcastor/iwork-converter/proto/TSK"
	"github.com/orcastor/iwork-converter/proto/TSP"
	"github.com/orcastor/iwork-converter/proto/TSS"
	"github.com/orcastor/iwork-converter/proto/TST"
	"github.com/orcastor/iwork-converter/proto/TSWP"

	"github.com/golang/protobuf/proto"
)

func decodeCommon(typ uint32, payload []byte) (interface{}, error) {
	switch typ {

	case 11000:
		value := &TSP.PasteboardObject{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11006:
		value := &TSP.PackageMetadata{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11007:
		value := &TSP.PasteboardMetadata{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11008:
		value := &TSP.ObjectContainer{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11011:
		value := &TSP.ViewStateMetadata{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 200:
		value := &TSK.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2001:
		value := &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2002:
		value := &TSWP.SelectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2003:
		value := &TSWP.DrawableAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2004:
		value := &TSWP.TextualAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2005:
		value := &TSWP.StorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2006:
		value := &TSWP.UIGraphicalAttachment{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2007:
		value := &TSWP.TextualAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2008:
		value := &TSWP.FootnoteReferenceAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2009:
		value := &TSWP.TextualAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 201:
		value := &TSK.CommandHistory{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2010:
		value := &TSWP.TSWPTOCPageNumberAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2011:
		value := &TSWP.ShapeInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2013:
		value := &TSWP.HighlightArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2014:
		value := &TSWP.CommentInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 202:
		value := &TSK.CommandGroupArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2021:
		value := &TSWP.CharacterStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2022:
		value := &TSWP.ParagraphStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2023:
		value := &TSWP.ListStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2024:
		value := &TSWP.ColumnStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2025:
		value := &TSWP.ShapeStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2026:
		value := &TSWP.TOCEntryStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 203:
		value := &TSK.CommandContainerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2031:
		value := &TSWP.PlaceholderSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2032:
		value := &TSWP.HyperlinkFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2033:
		value := &TSWP.FilenameSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2034:
		value := &TSWP.DateTimeSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2035:
		value := &TSWP.BookmarkFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2036:
		value := &TSWP.MergeSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2037:
		value := &TSWP.CitationRecordArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2038:
		value := &TSWP.CitationSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2039:
		value := &TSWP.UnsupportedHyperlinkFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 204:
		value := &TSK.ReplaceAllCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2040:
		value := &TSWP.BibliographySmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2041:
		value := &TSWP.TOCSmartFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2042:
		value := &TSWP.RubyFieldArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2043:
		value := &TSWP.NumberAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 205:
		value := &TSK.TreeNode{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2050:
		value := &TSWP.TextStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2051:
		value := &TSWP.TOCSettingsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2052:
		value := &TSWP.TOCEntryInstanceArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 206:
		value := &TSK.ProgressiveCommandGroupArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2060:
		value := &TSWP.ChangeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2061:
		value := &TSK.DeprecatedChangeAuthorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2062:
		value := &TSWP.ChangeSessionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 208:
		value := &TSK.CommandSelectionBehaviorHistoryArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 209:
		value := &TSK.UndoRedoStateCommandSelectionBehaviorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 210:
		value := &TSK.ViewStateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2101:
		value := &TSWP.TextCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2102:
		value := &TSWP.InsertAttachmentCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2104:
		value := &TSWP.ReplaceAllTextCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2105:
		value := &TSWP.FormatTextCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2107:
		value := &TSWP.ApplyPlaceholderTextCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2108:
		value := &TSWP.ApplyHighlightTextCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 211:
		value := &TSK.DocumentSupportArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2113:
		value := &TSWP.CreateHyperlinkCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2114:
		value := &TSWP.RemoveHyperlinkCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2115:
		value := &TSWP.ModifyHyperlinkCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2116:
		value := &TSWP.ApplyRubyTextCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2117:
		value := &TSWP.RemoveRubyTextCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2118:
		value := &TSWP.ModifyRubyTextCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2119:
		value := &TSWP.UpdateDateTimeFieldCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 212:
		value := &TSK.AnnotationAuthorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2120:
		value := &TSWP.ModifyTOCSettingsBaseCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2121:
		value := &TSWP.ModifyTOCSettingsForTOCInfoCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2122:
		value := &TSWP.ModifyTOCSettingsPresetForThemeCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 213:
		value := &TSK.AnnotationAuthorStorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 214:
		value := &TSK.AddAnnotationAuthorCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 215:
		value := &TSK.SetAnnotationAuthorColorCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2206:
		value := &TSWP.AnchorAttachmentCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2207:
		value := &TSWP.TextApplyThemeCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2231:
		value := &TSWP.ShapeApplyPresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2232:
		value := &TSWP.ShapePasteStyleCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2240:
		value := &TSWP.TOCInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2241:
		value := &TSWP.TOCAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2242:
		value := &TSWP.TOCLayoutHintArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2400:
		value := &TSWP.StyleBaseCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2401:
		value := &TSWP.StyleCreateCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2402:
		value := &TSWP.StyleRenameCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2403:
		value := &TSWP.StyleUpdateCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2404:
		value := &TSWP.StyleDeleteCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2405:
		value := &TSWP.StyleReorderCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2406:
		value := &TSWP.StyleUpdatePropertyMapCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3002:
		value := &TSD.DrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3003:
		value := &TSD.ContainerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3004:
		value := &TSD.ShapeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3005:
		value := &TSD.ImageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3006:
		value := &TSD.MaskArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3007:
		value := &TSD.MovieArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3008:
		value := &TSD.GroupArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3009:
		value := &TSD.ConnectionLineArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3015:
		value := &TSD.ShapeStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3016:
		value := &TSD.MediaStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3020:
		value := &TSD.DrawablesCommandGroupArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3021:
		value := &TSD.InfoGeometryCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3022:
		value := &TSD.DrawablePathSourceCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3023:
		value := &TSD.ShapePathSourceFlipCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3024:
		value := &TSD.ImageMaskCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3025:
		value := &TSD.ImageMediaCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3026:
		value := &TSD.ImageReplaceCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3027:
		value := &TSD.MediaOriginalSizeCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3028:
		value := &TSD.ShapeStyleSetValueCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3030:
		value := &TSD.MediaStyleSetValueCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3031:
		value := &TSD.ShapeApplyPresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3032:
		value := &TSD.MediaApplyPresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3033:
		value := &TSD.DrawableApplyThemeCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3034:
		value := &TSD.MovieSetValueCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3035:
		value := &TSD.ShapeSetLineEndCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3036:
		value := &TSD.ExteriorTextWrapCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3037:
		value := &TSD.MediaFlagsCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3038:
		value := &TSD.GroupDrawablesCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3039:
		value := &TSD.UngroupGroupCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3040:
		value := &TSD.DrawableHyperlinkCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3041:
		value := &TSD.ConnectionLineConnectCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3042:
		value := &TSD.InstantAlphaCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3043:
		value := &TSD.DrawableLockCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3045:
		value := &TSD.CanvasSelectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3046:
		value := &TSD.CommandSelectionBehaviorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3047:
		value := &TSD.GuideStorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3048:
		value := &TSD.StyledInfoSetStyleCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3049:
		value := &TSD.DrawableInfoCommentCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3050:
		value := &TSD.GuideCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3051:
		value := &TSD.DrawableAspectRatioLockedCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3052:
		value := &TSD.ContainerRemoveChildrenCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3053:
		value := &TSD.ContainerInsertChildrenCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3054:
		value := &TSD.ContainerReorderChildrenCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3055:
		value := &TSD.ImageAdjustmentsCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3056:
		value := &TSD.CommentStorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3057:
		value := &TSD.ThemeReplaceFillPresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3058:
		value := &TSD.DrawableAccessibilityDescriptionCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3059:
		value := &TSD.PasteStyleCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3060:
		value := &TSD.CommentStorageApplyCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 400:
		value := &TSS.StyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4000:
		value := &TSCE.CalculationEngineArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4001:
		value := &TSCE.FormulaRewriteCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4002:
		value := &TSCE.TrackedReferencesRewriteCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4003:
		value := &TSCE.NamedReferenceManagerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4004:
		value := &TSCE.ReferenceTrackerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4005:
		value := &TSCE.TrackedReferenceArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 401:
		value := &TSS.StylesheetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 402:
		value := &TSS.ThemeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 410:
		value := &TSS.ApplyThemeCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 411:
		value := &TSS.ApplyThemeChildCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 412:
		value := &TSS.StyleUpdatePropertyMapCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 413:
		value := &TSS.ThemeReplacePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 414:
		value := &TSS.ThemeAddStylePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 415:
		value := &TSS.ThemeRemoveStylePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 416:
		value := &TSS.ThemeReplaceColorPresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 417:
		value := &TSS.ThemeMovePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 418:
		value := &TSS.ThemeReplaceStylePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5000:
		value := &PreUFF.ChartInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5002:
		value := &PreUFF.ChartGridArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5004:
		value := &TSCH.ChartMediatorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5010:
		value := &PreUFF.ChartStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5011:
		value := &PreUFF.ChartSeriesStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5012:
		value := &PreUFF.ChartAxisStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5013:
		value := &PreUFF.LegendStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5014:
		value := &PreUFF.ChartNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5015:
		value := &PreUFF.ChartSeriesNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5016:
		value := &PreUFF.ChartAxisNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5017:
		value := &PreUFF.LegendNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5020:
		value := &TSCH.ChartStylePreset{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5021:
		value := &TSCH.ChartDrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5022:
		value := &TSCH.ChartStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5023:
		value := &TSCH.ChartNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5024:
		value := &TSCH.LegendStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5025:
		value := &TSCH.LegendNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5026:
		value := &TSCH.ChartAxisStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5027:
		value := &TSCH.ChartAxisNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5028:
		value := &TSCH.ChartSeriesStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5029:
		value := &TSCH.ChartSeriesNonStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5103:
		value := &TSCH.CommandSetChartTypeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5104:
		value := &TSCH.CommandSetSeriesNameArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5105:
		value := &TSCH.CommandSetCategoryNameArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5107:
		value := &TSCH.CommandSetScatterFormatArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5108:
		value := &TSCH.CommandSetLegendFrameArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5109:
		value := &TSCH.CommandSetGridValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5110:
		value := &TSCH.CommandSetGridDirectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5113:
		value := &TSCH.SynchronousCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5114:
		value := &TSCH.CommandReplaceAllArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5115:
		value := &TSCH.CommandAddGridRowsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5116:
		value := &TSCH.CommandAddGridColumnsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5117:
		value := &TSCH.CommandSetPreviewLocArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5118:
		value := &TSCH.CommandMoveGridRowsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5119:
		value := &TSCH.CommandMoveGridColumnsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5120:
		value := &TSCH.CommandDeleteGridRowsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5121:
		value := &TSCH.CommandDeleteGridColumnsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5122:
		value := &TSCH.CommandSetPieWedgeExplosion{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5123:
		value := &TSCH.CommandStyleSwapArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5124:
		value := &TSCH.CommandChartApplyTheme{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5125:
		value := &TSCH.CommandChartApplyPreset{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5126:
		value := &TSCH.ChartCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5127:
		value := &TSCH.CommandReplaceGridValuesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5129:
		value := &TSCH.StylePasteboardDataArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5130:
		value := &TSCH.CommandSetMultiDataSetIndexArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5131:
		value := &TSCH.CommandReplaceThemePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5132:
		value := &TSCH.CommandInvalidateWPCaches{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 600:
		value := &TSA.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6000:
		value := &TST.TableInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6001:
		value := &TST.TableModelArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6002:
		value := &TST.Tile{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6003:
		value := &TST.TableStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6004:
		value := &TST.CellStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6005:
		value := &TST.TableDataList{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6006:
		value := &TST.HeaderStorageBucket{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6007:
		value := &TST.WPTableInfoArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6008:
		value := &TST.TableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6009:
		value := &TST.TableStrokePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 601:
		value := &TSA.FunctionBrowserStateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6010:
		value := &TST.ConditionalStyleSetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 602:
		value := &TSA.PropagatePresetCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6100:
		value := &TST.TableCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6101:
		value := &TST.CommandDeleteCellsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6102:
		value := &TST.CommandInsertColumnsOrRowsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6103:
		value := &TST.CommandRemoveColumnsOrRowsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6104:
		value := &TST.CommandResizeColumnOrRowArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6105:
		value := &TST.CommandSetCellArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6106:
		value := &TST.CommandSetNumberOfHeadersOrFootersArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6107:
		value := &TST.CommandSetTableNameArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6108:
		value := &TST.CommandStyleCellsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6109:
		value := &TST.CommandFillCellsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6110:
		value := &TST.CommandReplaceAllTextArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6111:
		value := &TST.CommandChangeFreezeHeaderStateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6112:
		value := &TST.CommandReplaceTextArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6113:
		value := &TST.CommandPasteArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6114:
		value := &TST.CommandSetTableNameEnabledArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6115:
		value := &TST.CommandMoveRowsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6116:
		value := &TST.CommandMoveColumnsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6117:
		value := &TST.CommandApplyTableStylePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6118:
		value := &TST.CommandApplyStrokePresetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6119:
		value := &TST.CommandSetExplicitFormatArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6120:
		value := &TST.CommandSetRepeatingHeaderEnabledArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6121:
		value := &TST.CommandApplyThemeToTableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6122:
		value := &TST.CommandApplyThemeChildForTableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6123:
		value := &TST.CommandSortArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6124:
		value := &TST.CommandToggleTextPropertyArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6125:
		value := &TST.CommandStyleTableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6126:
		value := &TST.CommandSetNumberOfDecimalPlacesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6127:
		value := &TST.CommandSetShowThousandsSeparatorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6128:
		value := &TST.CommandSetNegativeNumberStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6129:
		value := &TST.CommandSetFractionAccuracyArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6130:
		value := &TST.CommandSetSingleNumberFormatParameterArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6131:
		value := &TST.CommandSetCurrencyCodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6132:
		value := &TST.CommandSetUseAccountingStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6134:
		value := &TST.CommandRewriteFormulasForSortArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6135:
		value := &TST.CommandRewriteFormulasForTectonicShiftArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6136:
		value := &TST.CommandSetTableFontNameArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6137:
		value := &TST.CommandSetTableFontSizeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6138:
		value := &TST.CommandRewriteFormulasForMoveArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6139:
		value := &TST.CommandFixStylesInHeadersOrFootersArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6141:
		value := &TST.CommandResetFillPropertyToDefault{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6142:
		value := &TST.CommandSetTableNameHeightArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6143:
		value := &TST.CommandMergeUnmergeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6144:
		value := &TST.MergeRegionMapArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6145:
		value := &TST.CommandHideShowArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6146:
		value := &TST.CommandSetBaseArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6147:
		value := &TST.CommandSetBasePlacesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6148:
		value := &TST.CommandSetBaseUseMinusSignArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6179:
		value := &TST.FormulaEqualsTokenAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6181:
		value := &TST.TokenAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6182:
		value := &TST.ExpressionNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6183:
		value := &TST.BooleanNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6184:
		value := &TST.NumberNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6185:
		value := &TST.StringNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6186:
		value := &TST.ArrayNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6187:
		value := &TST.ListNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6188:
		value := &TST.OperatorNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6189:
		value := &TST.FunctionNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6190:
		value := &TST.DateNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6191:
		value := &TST.ReferenceNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6192:
		value := &TST.DurationNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6193:
		value := &TST.ArgumentPlaceholderNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6194:
		value := &TST.PostfixOperatorNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6195:
		value := &TST.PrefixOperatorNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6196:
		value := &TST.FunctionEndNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6197:
		value := &TST.EmptyExpressionNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6198:
		value := &TST.LayoutHintArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6199:
		value := &TST.CompletionTokenAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6200:
		value := &TST.FormulaEditingCommandGroupArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6201:
		value := &TST.TableDataList{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6202:
		value := &TST.CommandCoerceMultipleCellsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6203:
		value := &TST.CommandSetMultipleCellsCustomArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6204:
		value := &TST.HiddenStateFormulaOwnerArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6205:
		value := &TST.CommandSetAutomaticDurationUnitsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6206:
		value := &TST.PopUpMenuModel{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6207:
		value := &TST.CommandSetControlMinimumArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6208:
		value := &TST.CommandSetControlMaximumArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6209:
		value := &TST.CommandSetControlIncrementArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6210:
		value := &TST.CommandSetControlCellsDisplayNumberFormatArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6211:
		value := &TST.CommandSetMultipleCellsMultipleChoiceListArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6212:
		value := &TST.CommandSetMultipleChoiceListFormatForEditedItemArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6213:
		value := &TST.CommandSetMultipleChoiceListFormatForDeleteItemArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6214:
		value := &TST.CommandSetMultipleChoiceListFormatForReorderItemArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6215:
		value := &TST.CommandSetMultipleChoiceListFormatForInitialValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6216:
		value := &TST.CommandRewriteFormulasForCellMergeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6217:
		value := &TST.TableInfoGeometryCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6218:
		value := &TST.RichTextPayloadArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6219:
		value := &TST.EditingStateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6220:
		value := &TST.FilterSetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6221:
		value := &TST.CommandSetFiltersEnabledArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6222:
		value := &TST.CommandRewriteFilterFormulasForTectonicShiftArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6223:
		value := &TST.CommandRewriteFilterFormulasForSortArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6224:
		value := &TST.CommandRewriteFilterFormulasForTableResizeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6225:
		value := &TST.CommandSetAutomaticFormatArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6226:
		value := &TST.CommandTextPreflightInsertCellArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6227:
		value := &TST.FormulaEditingCommandSelectionBehaviorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6228:
		value := &TST.CommandDeleteCellContentsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6229:
		value := &TST.CommandPostflightSetCellArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6231:
		value := &TST.CommandRewriteConditionalStylesForTectonicShiftArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6232:
		value := &TST.CommandRewriteConditionalStylesForSortArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6233:
		value := &TST.CommandRewriteConditionalStylesForRangeMoveArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6234:
		value := &TST.CommandRewriteConditionalStylesForCellMergeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6235:
		value := &TST.IdentifierNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6236:
		value := &TST.UndoRedoStateCommandSelectionBehaviorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6237:
		value := &TST.CommandSetStyleApplyClearsAllFlagArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6238:
		value := &TST.CommandSetDateTimeFormatArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6239:
		value := &TST.TableCommandSelectionBehaviorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6240:
		value := &TST.CommandAddQuickFilterRulesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6241:
		value := &TST.CommandModifyFilterRuleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6242:
		value := &TST.CommandDeleteFilterRulesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6244:
		value := &TST.CommandApplyCellCommentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6245:
		value := &TST.CommandApplyConditionalStyleSetArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6246:
		value := &TST.CommandSetFormulaTokenizationArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6247:
		value := &TST.TableStyleNetworkArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6248:
		value := &TST.CommandSetFilterEnabledArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6249:
		value := &TST.CommandSetFilterRuleEnabledArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6250:
		value := &TST.CommandSetFilterSetTypeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6251:
		value := &TST.CommandSetStyleNetworkArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6252:
		value := &TST.CommandMutateCellsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6253:
		value := &TST.DisableTableNameSelectionBehaviorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6254:
		value := &TST.CommandDisableFilterRulesForColumnArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6255:
		value := &TST.CommandSetTextStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6256:
		value := &TST.CommandNotifyForTransformingArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	default:
		return nil, fmt.Errorf("Unknown type %d", typ)
	}
}
