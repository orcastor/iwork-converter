// Code generated by protoc-gen-go.
// source: TNArchives.proto
// DO NOT EDIT!

/*
Package TN is a generated protocol buffer package.

It is generated from these files:
	TNArchives.proto

It has these top-level messages:
	SheetUIStateArchive
	SheetUIStateDictionaryEntryArchive
	UIStateArchive
	SheetSelectionArchive
	UndoRedoStateArchive
	DocumentArchive
	PlaceholderArchive
	SheetArchive
	FormBasedSheetArchive
	ThemeArchive
	ChartMediatorFormulaStorage
	ChartMediatorArchive
	ChartSelectionArchive
*/
package TN

import proto "code.google.com/p/goprotobuf/proto"
import math "math"
import "github.com/dunhamsteve/iwork/proto/TSP"

// discarding unused import TSK "TSKArchives.pb"
import "github.com/dunhamsteve/iwork/proto/TSCH2"
import "github.com/dunhamsteve/iwork/proto/TSCE"
import "github.com/dunhamsteve/iwork/proto/TSS"
import "github.com/dunhamsteve/iwork/proto/TSD"
import "github.com/dunhamsteve/iwork/proto/TSWP"
import "github.com/dunhamsteve/iwork/proto/TSA"
import "github.com/dunhamsteve/iwork/proto/TST1"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type SheetPageOrder int32

const (
	SheetPageOrder_SheetPageOrderTopToBottom SheetPageOrder = 0
	SheetPageOrder_SheetPageOrderLeftToRight SheetPageOrder = 1
)

var SheetPageOrder_name = map[int32]string{
	0: "SheetPageOrderTopToBottom",
	1: "SheetPageOrderLeftToRight",
}
var SheetPageOrder_value = map[string]int32{
	"SheetPageOrderTopToBottom": 0,
	"SheetPageOrderLeftToRight": 1,
}

func (x SheetPageOrder) Enum() *SheetPageOrder {
	p := new(SheetPageOrder)
	*p = x
	return p
}
func (x SheetPageOrder) String() string {
	return proto.EnumName(SheetPageOrder_name, int32(x))
}
func (x *SheetPageOrder) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SheetPageOrder_value, data, "SheetPageOrder")
	if err != nil {
		return err
	}
	*x = SheetPageOrder(value)
	return nil
}

type UIStateArchive_InspectorPaneViewMode int32

const (
	UIStateArchive_kInspectorPaneViewModeFormat UIStateArchive_InspectorPaneViewMode = 0
	UIStateArchive_kInspectorPaneViewModeFilter UIStateArchive_InspectorPaneViewMode = 1
)

var UIStateArchive_InspectorPaneViewMode_name = map[int32]string{
	0: "kInspectorPaneViewModeFormat",
	1: "kInspectorPaneViewModeFilter",
}
var UIStateArchive_InspectorPaneViewMode_value = map[string]int32{
	"kInspectorPaneViewModeFormat": 0,
	"kInspectorPaneViewModeFilter": 1,
}

func (x UIStateArchive_InspectorPaneViewMode) Enum() *UIStateArchive_InspectorPaneViewMode {
	p := new(UIStateArchive_InspectorPaneViewMode)
	*p = x
	return p
}
func (x UIStateArchive_InspectorPaneViewMode) String() string {
	return proto.EnumName(UIStateArchive_InspectorPaneViewMode_name, int32(x))
}
func (x *UIStateArchive_InspectorPaneViewMode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(UIStateArchive_InspectorPaneViewMode_value, data, "UIStateArchive_InspectorPaneViewMode")
	if err != nil {
		return err
	}
	*x = UIStateArchive_InspectorPaneViewMode(value)
	return nil
}

type SheetUIStateArchive struct {
	ViewScale                   *float32   `protobuf:"fixed32,1,req,name=view_scale" json:"view_scale,omitempty"`
	ScrollPosition              *TSP.Point `protobuf:"bytes,2,req,name=scroll_position" json:"scroll_position,omitempty"`
	PreviousViewScale           *float32   `protobuf:"fixed32,3,opt,name=previous_view_scale" json:"previous_view_scale,omitempty"`
	ScrollPositionIsUnscaled    *bool      `protobuf:"varint,4,opt,name=scroll_position_is_unscaled" json:"scroll_position_is_unscaled,omitempty"`
	PreviousScrollPosition      *TSP.Point `protobuf:"bytes,5,opt,name=previous_scroll_position" json:"previous_scroll_position,omitempty"`
	ScrollPositionValid         *bool      `protobuf:"varint,6,opt,name=scroll_position_valid" json:"scroll_position_valid,omitempty"`
	PreviousScrollPositionValid *bool      `protobuf:"varint,7,opt,name=previous_scroll_position_valid" json:"previous_scroll_position_valid,omitempty"`
	VisibleSize                 *TSP.Size  `protobuf:"bytes,8,opt,name=visible_size" json:"visible_size,omitempty"`
	PreviousVisibleSize         *TSP.Size  `protobuf:"bytes,9,opt,name=previous_visible_size" json:"previous_visible_size,omitempty"`
	DeviceIdiom                 *uint32    `protobuf:"varint,10,opt,name=device_idiom" json:"device_idiom,omitempty"`
	FormFocusedRecordIndex      *uint32    `protobuf:"varint,11,opt,name=form_focused_record_index" json:"form_focused_record_index,omitempty"`
	FormFocusedFieldIndex       *uint32    `protobuf:"varint,12,opt,name=form_focused_field_index" json:"form_focused_field_index,omitempty"`
	XXX_unrecognized            []byte     `json:"-"`
}

func (m *SheetUIStateArchive) Reset()         { *m = SheetUIStateArchive{} }
func (m *SheetUIStateArchive) String() string { return proto.CompactTextString(m) }
func (*SheetUIStateArchive) ProtoMessage()    {}

func (m *SheetUIStateArchive) GetViewScale() float32 {
	if m != nil && m.ViewScale != nil {
		return *m.ViewScale
	}
	return 0
}

func (m *SheetUIStateArchive) GetScrollPosition() *TSP.Point {
	if m != nil {
		return m.ScrollPosition
	}
	return nil
}

func (m *SheetUIStateArchive) GetPreviousViewScale() float32 {
	if m != nil && m.PreviousViewScale != nil {
		return *m.PreviousViewScale
	}
	return 0
}

func (m *SheetUIStateArchive) GetScrollPositionIsUnscaled() bool {
	if m != nil && m.ScrollPositionIsUnscaled != nil {
		return *m.ScrollPositionIsUnscaled
	}
	return false
}

func (m *SheetUIStateArchive) GetPreviousScrollPosition() *TSP.Point {
	if m != nil {
		return m.PreviousScrollPosition
	}
	return nil
}

func (m *SheetUIStateArchive) GetScrollPositionValid() bool {
	if m != nil && m.ScrollPositionValid != nil {
		return *m.ScrollPositionValid
	}
	return false
}

func (m *SheetUIStateArchive) GetPreviousScrollPositionValid() bool {
	if m != nil && m.PreviousScrollPositionValid != nil {
		return *m.PreviousScrollPositionValid
	}
	return false
}

func (m *SheetUIStateArchive) GetVisibleSize() *TSP.Size {
	if m != nil {
		return m.VisibleSize
	}
	return nil
}

func (m *SheetUIStateArchive) GetPreviousVisibleSize() *TSP.Size {
	if m != nil {
		return m.PreviousVisibleSize
	}
	return nil
}

func (m *SheetUIStateArchive) GetDeviceIdiom() uint32 {
	if m != nil && m.DeviceIdiom != nil {
		return *m.DeviceIdiom
	}
	return 0
}

func (m *SheetUIStateArchive) GetFormFocusedRecordIndex() uint32 {
	if m != nil && m.FormFocusedRecordIndex != nil {
		return *m.FormFocusedRecordIndex
	}
	return 0
}

func (m *SheetUIStateArchive) GetFormFocusedFieldIndex() uint32 {
	if m != nil && m.FormFocusedFieldIndex != nil {
		return *m.FormFocusedFieldIndex
	}
	return 0
}

type SheetUIStateDictionaryEntryArchive struct {
	Sheet            *TSP.Reference       `protobuf:"bytes,1,req,name=sheet" json:"sheet,omitempty"`
	SheetUistate     *SheetUIStateArchive `protobuf:"bytes,2,req,name=sheet_uistate" json:"sheet_uistate,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *SheetUIStateDictionaryEntryArchive) Reset()         { *m = SheetUIStateDictionaryEntryArchive{} }
func (m *SheetUIStateDictionaryEntryArchive) String() string { return proto.CompactTextString(m) }
func (*SheetUIStateDictionaryEntryArchive) ProtoMessage()    {}

func (m *SheetUIStateDictionaryEntryArchive) GetSheet() *TSP.Reference {
	if m != nil {
		return m.Sheet
	}
	return nil
}

func (m *SheetUIStateDictionaryEntryArchive) GetSheetUistate() *SheetUIStateArchive {
	if m != nil {
		return m.SheetUistate
	}
	return nil
}

type UIStateArchive struct {
	ActiveSheetIndex                    *uint32                               `protobuf:"varint,1,req,name=active_sheet_index" json:"active_sheet_index,omitempty"`
	SelectedInfo                        []*TSP.Reference                      `protobuf:"bytes,2,rep,name=selected_info" json:"selected_info,omitempty"`
	SheetUistateDictionaryEntry         []*SheetUIStateDictionaryEntryArchive `protobuf:"bytes,3,rep,name=sheet_uistate_dictionary_entry" json:"sheet_uistate_dictionary_entry,omitempty"`
	TableSelection                      *TST1.SelectionArchive                `protobuf:"bytes,4,opt,name=table_selection" json:"table_selection,omitempty"`
	EditingSheetIndex                   *uint32                               `protobuf:"varint,5,opt,name=editing_sheet_index" json:"editing_sheet_index,omitempty"`
	DocumentMode                        *int32                                `protobuf:"varint,6,opt,name=document_mode" json:"document_mode,omitempty"`
	EditModeSheetUistateDictionaryEntry []*SheetUIStateDictionaryEntryArchive `protobuf:"bytes,7,rep,name=edit_mode_sheet_uistate_dictionary_entry" json:"edit_mode_sheet_uistate_dictionary_entry,omitempty"`
	TableEditingMode                    *int32                                `protobuf:"varint,8,opt,name=table_editing_mode" json:"table_editing_mode,omitempty"`
	FormFocusedRecordIndex              *uint32                               `protobuf:"varint,9,opt,name=form_focused_record_index" json:"form_focused_record_index,omitempty"`
	FormFocusedFieldIndex               *uint32                               `protobuf:"varint,10,opt,name=form_focused_field_index" json:"form_focused_field_index,omitempty"`
	InChartMode                         *bool                                 `protobuf:"varint,11,opt,name=in_chart_mode" json:"in_chart_mode,omitempty"`
	ChartSelection                      *ChartSelectionArchive                `protobuf:"bytes,12,opt,name=chart_selection" json:"chart_selection,omitempty"`
	SheetSelection                      *TSP.Reference                        `protobuf:"bytes,13,opt,name=sheet_selection" json:"sheet_selection,omitempty"`
	InspectorPaneVisible                *bool                                 `protobuf:"varint,14,opt,name=inspector_pane_visible,def=1" json:"inspector_pane_visible,omitempty"`
	InspectorPaneViewMode               *UIStateArchive_InspectorPaneViewMode `protobuf:"varint,15,opt,name=inspector_pane_view_mode,enum=TN.UIStateArchive_InspectorPaneViewMode,def=0" json:"inspector_pane_view_mode,omitempty"`
	SelectedQuickCalcFunctions          []uint32                              `protobuf:"varint,16,rep,name=selected_quick_calc_functions" json:"selected_quick_calc_functions,omitempty"`
	RemovedAllQuickCalcFunctions        *bool                                 `protobuf:"varint,17,opt,name=removed_all_quick_calc_functions" json:"removed_all_quick_calc_functions,omitempty"`
	ShowCanvasGuides                    *bool                                 `protobuf:"varint,18,opt,name=show_canvas_guides" json:"show_canvas_guides,omitempty"`
	ShowsComments                       *bool                                 `protobuf:"varint,19,opt,name=shows_comments" json:"shows_comments,omitempty"`
	XXX_unrecognized                    []byte                                `json:"-"`
}

func (m *UIStateArchive) Reset()         { *m = UIStateArchive{} }
func (m *UIStateArchive) String() string { return proto.CompactTextString(m) }
func (*UIStateArchive) ProtoMessage()    {}

const Default_UIStateArchive_InspectorPaneVisible bool = true
const Default_UIStateArchive_InspectorPaneViewMode UIStateArchive_InspectorPaneViewMode = UIStateArchive_kInspectorPaneViewModeFormat

func (m *UIStateArchive) GetActiveSheetIndex() uint32 {
	if m != nil && m.ActiveSheetIndex != nil {
		return *m.ActiveSheetIndex
	}
	return 0
}

func (m *UIStateArchive) GetSelectedInfo() []*TSP.Reference {
	if m != nil {
		return m.SelectedInfo
	}
	return nil
}

func (m *UIStateArchive) GetSheetUistateDictionaryEntry() []*SheetUIStateDictionaryEntryArchive {
	if m != nil {
		return m.SheetUistateDictionaryEntry
	}
	return nil
}

func (m *UIStateArchive) GetTableSelection() *TST1.SelectionArchive {
	if m != nil {
		return m.TableSelection
	}
	return nil
}

func (m *UIStateArchive) GetEditingSheetIndex() uint32 {
	if m != nil && m.EditingSheetIndex != nil {
		return *m.EditingSheetIndex
	}
	return 0
}

func (m *UIStateArchive) GetDocumentMode() int32 {
	if m != nil && m.DocumentMode != nil {
		return *m.DocumentMode
	}
	return 0
}

func (m *UIStateArchive) GetEditModeSheetUistateDictionaryEntry() []*SheetUIStateDictionaryEntryArchive {
	if m != nil {
		return m.EditModeSheetUistateDictionaryEntry
	}
	return nil
}

func (m *UIStateArchive) GetTableEditingMode() int32 {
	if m != nil && m.TableEditingMode != nil {
		return *m.TableEditingMode
	}
	return 0
}

func (m *UIStateArchive) GetFormFocusedRecordIndex() uint32 {
	if m != nil && m.FormFocusedRecordIndex != nil {
		return *m.FormFocusedRecordIndex
	}
	return 0
}

func (m *UIStateArchive) GetFormFocusedFieldIndex() uint32 {
	if m != nil && m.FormFocusedFieldIndex != nil {
		return *m.FormFocusedFieldIndex
	}
	return 0
}

func (m *UIStateArchive) GetInChartMode() bool {
	if m != nil && m.InChartMode != nil {
		return *m.InChartMode
	}
	return false
}

func (m *UIStateArchive) GetChartSelection() *ChartSelectionArchive {
	if m != nil {
		return m.ChartSelection
	}
	return nil
}

func (m *UIStateArchive) GetSheetSelection() *TSP.Reference {
	if m != nil {
		return m.SheetSelection
	}
	return nil
}

func (m *UIStateArchive) GetInspectorPaneVisible() bool {
	if m != nil && m.InspectorPaneVisible != nil {
		return *m.InspectorPaneVisible
	}
	return Default_UIStateArchive_InspectorPaneVisible
}

func (m *UIStateArchive) GetInspectorPaneViewMode() UIStateArchive_InspectorPaneViewMode {
	if m != nil && m.InspectorPaneViewMode != nil {
		return *m.InspectorPaneViewMode
	}
	return Default_UIStateArchive_InspectorPaneViewMode
}

func (m *UIStateArchive) GetSelectedQuickCalcFunctions() []uint32 {
	if m != nil {
		return m.SelectedQuickCalcFunctions
	}
	return nil
}

func (m *UIStateArchive) GetRemovedAllQuickCalcFunctions() bool {
	if m != nil && m.RemovedAllQuickCalcFunctions != nil {
		return *m.RemovedAllQuickCalcFunctions
	}
	return false
}

func (m *UIStateArchive) GetShowCanvasGuides() bool {
	if m != nil && m.ShowCanvasGuides != nil {
		return *m.ShowCanvasGuides
	}
	return false
}

func (m *UIStateArchive) GetShowsComments() bool {
	if m != nil && m.ShowsComments != nil {
		return *m.ShowsComments
	}
	return false
}

type SheetSelectionArchive struct {
	Sheet            *TSP.Reference `protobuf:"bytes,1,opt,name=sheet" json:"sheet,omitempty"`
	Paginated        *bool          `protobuf:"varint,2,opt,name=paginated" json:"paginated,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *SheetSelectionArchive) Reset()         { *m = SheetSelectionArchive{} }
func (m *SheetSelectionArchive) String() string { return proto.CompactTextString(m) }
func (*SheetSelectionArchive) ProtoMessage()    {}

func (m *SheetSelectionArchive) GetSheet() *TSP.Reference {
	if m != nil {
		return m.Sheet
	}
	return nil
}

func (m *SheetSelectionArchive) GetPaginated() bool {
	if m != nil && m.Paginated != nil {
		return *m.Paginated
	}
	return false
}

type UndoRedoStateArchive struct {
	UiState          *UIStateArchive `protobuf:"bytes,1,req,name=ui_state" json:"ui_state,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *UndoRedoStateArchive) Reset()         { *m = UndoRedoStateArchive{} }
func (m *UndoRedoStateArchive) String() string { return proto.CompactTextString(m) }
func (*UndoRedoStateArchive) ProtoMessage()    {}

func (m *UndoRedoStateArchive) GetUiState() *UIStateArchive {
	if m != nil {
		return m.UiState
	}
	return nil
}

type DocumentArchive struct {
	Sheets            []*TSP.Reference     `protobuf:"bytes,1,rep,name=sheets" json:"sheets,omitempty"`
	Super             *TSA.DocumentArchive `protobuf:"bytes,8,req,name=super" json:"super,omitempty"`
	CalculationEngine *TSP.Reference       `protobuf:"bytes,3,opt,name=calculation_engine" json:"calculation_engine,omitempty"`
	Stylesheet        *TSP.Reference       `protobuf:"bytes,4,req,name=stylesheet" json:"stylesheet,omitempty"`
	SidebarOrder      *TSP.Reference       `protobuf:"bytes,5,req,name=sidebar_order" json:"sidebar_order,omitempty"`
	Theme             *TSP.Reference       `protobuf:"bytes,6,req,name=theme" json:"theme,omitempty"`
	Uistate           *UIStateArchive      `protobuf:"bytes,7,opt,name=uistate" json:"uistate,omitempty"`
	CustomFormatList  *TSP.Reference       `protobuf:"bytes,9,opt,name=custom_format_list" json:"custom_format_list,omitempty"`
	PrinterId         *string              `protobuf:"bytes,10,opt,name=printer_id" json:"printer_id,omitempty"`
	PaperId           *string              `protobuf:"bytes,11,opt,name=paper_id" json:"paper_id,omitempty"`
	PageSize          *TSP.Size            `protobuf:"bytes,12,opt,name=page_size" json:"page_size,omitempty"`
	XXX_unrecognized  []byte               `json:"-"`
}

func (m *DocumentArchive) Reset()         { *m = DocumentArchive{} }
func (m *DocumentArchive) String() string { return proto.CompactTextString(m) }
func (*DocumentArchive) ProtoMessage()    {}

func (m *DocumentArchive) GetSheets() []*TSP.Reference {
	if m != nil {
		return m.Sheets
	}
	return nil
}

func (m *DocumentArchive) GetSuper() *TSA.DocumentArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *DocumentArchive) GetCalculationEngine() *TSP.Reference {
	if m != nil {
		return m.CalculationEngine
	}
	return nil
}

func (m *DocumentArchive) GetStylesheet() *TSP.Reference {
	if m != nil {
		return m.Stylesheet
	}
	return nil
}

func (m *DocumentArchive) GetSidebarOrder() *TSP.Reference {
	if m != nil {
		return m.SidebarOrder
	}
	return nil
}

func (m *DocumentArchive) GetTheme() *TSP.Reference {
	if m != nil {
		return m.Theme
	}
	return nil
}

func (m *DocumentArchive) GetUistate() *UIStateArchive {
	if m != nil {
		return m.Uistate
	}
	return nil
}

func (m *DocumentArchive) GetCustomFormatList() *TSP.Reference {
	if m != nil {
		return m.CustomFormatList
	}
	return nil
}

func (m *DocumentArchive) GetPrinterId() string {
	if m != nil && m.PrinterId != nil {
		return *m.PrinterId
	}
	return ""
}

func (m *DocumentArchive) GetPaperId() string {
	if m != nil && m.PaperId != nil {
		return *m.PaperId
	}
	return ""
}

func (m *DocumentArchive) GetPageSize() *TSP.Size {
	if m != nil {
		return m.PageSize
	}
	return nil
}

type PlaceholderArchive struct {
	Super            *TSWP.ShapeInfoArchive `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *PlaceholderArchive) Reset()         { *m = PlaceholderArchive{} }
func (m *PlaceholderArchive) String() string { return proto.CompactTextString(m) }
func (*PlaceholderArchive) ProtoMessage()    {}

func (m *PlaceholderArchive) GetSuper() *TSWP.ShapeInfoArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

type SheetArchive struct {
	Name                      *string                `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	DrawableInfos             []*TSP.Reference       `protobuf:"bytes,2,rep,name=drawable_infos" json:"drawable_infos,omitempty"`
	InPortraitPageOrientation *bool                  `protobuf:"varint,3,opt,name=in_portrait_page_orientation" json:"in_portrait_page_orientation,omitempty"`
	ShowRepeatingHeaders      *bool                  `protobuf:"varint,4,opt,name=show_repeating_headers" json:"show_repeating_headers,omitempty"`
	ShowPageNumbers           *bool                  `protobuf:"varint,5,opt,name=show_page_numbers" json:"show_page_numbers,omitempty"`
	IsAutofitOn               *bool                  `protobuf:"varint,6,opt,name=is_autofit_on" json:"is_autofit_on,omitempty"`
	ContentScale              *float32               `protobuf:"fixed32,7,opt,name=content_scale" json:"content_scale,omitempty"`
	PageOrder                 *SheetPageOrder        `protobuf:"varint,8,opt,name=page_order,enum=TN.SheetPageOrder" json:"page_order,omitempty"`
	PrintMargins              *TSD.EdgeInsetsArchive `protobuf:"bytes,10,opt,name=print_margins" json:"print_margins,omitempty"`
	UsingStartPageNumber      *bool                  `protobuf:"varint,11,opt,name=using_start_page_number" json:"using_start_page_number,omitempty"`
	StartPageNumber           *int32                 `protobuf:"varint,12,opt,name=start_page_number" json:"start_page_number,omitempty"`
	PageHeaderInset           *float32               `protobuf:"fixed32,13,opt,name=page_header_inset" json:"page_header_inset,omitempty"`
	PageFooterInset           *float32               `protobuf:"fixed32,14,opt,name=page_footer_inset" json:"page_footer_inset,omitempty"`
	HeaderStorage             *TSP.Reference         `protobuf:"bytes,15,opt,name=header_storage" json:"header_storage,omitempty"`
	FooterStorage             *TSP.Reference         `protobuf:"bytes,16,opt,name=footer_storage" json:"footer_storage,omitempty"`
	UserDefinedGuideStorage   *TSP.Reference         `protobuf:"bytes,17,opt,name=userDefinedGuideStorage" json:"userDefinedGuideStorage,omitempty"`
	XXX_unrecognized          []byte                 `json:"-"`
}

func (m *SheetArchive) Reset()         { *m = SheetArchive{} }
func (m *SheetArchive) String() string { return proto.CompactTextString(m) }
func (*SheetArchive) ProtoMessage()    {}

func (m *SheetArchive) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *SheetArchive) GetDrawableInfos() []*TSP.Reference {
	if m != nil {
		return m.DrawableInfos
	}
	return nil
}

func (m *SheetArchive) GetInPortraitPageOrientation() bool {
	if m != nil && m.InPortraitPageOrientation != nil {
		return *m.InPortraitPageOrientation
	}
	return false
}

func (m *SheetArchive) GetShowRepeatingHeaders() bool {
	if m != nil && m.ShowRepeatingHeaders != nil {
		return *m.ShowRepeatingHeaders
	}
	return false
}

func (m *SheetArchive) GetShowPageNumbers() bool {
	if m != nil && m.ShowPageNumbers != nil {
		return *m.ShowPageNumbers
	}
	return false
}

func (m *SheetArchive) GetIsAutofitOn() bool {
	if m != nil && m.IsAutofitOn != nil {
		return *m.IsAutofitOn
	}
	return false
}

func (m *SheetArchive) GetContentScale() float32 {
	if m != nil && m.ContentScale != nil {
		return *m.ContentScale
	}
	return 0
}

func (m *SheetArchive) GetPageOrder() SheetPageOrder {
	if m != nil && m.PageOrder != nil {
		return *m.PageOrder
	}
	return SheetPageOrder_SheetPageOrderTopToBottom
}

func (m *SheetArchive) GetPrintMargins() *TSD.EdgeInsetsArchive {
	if m != nil {
		return m.PrintMargins
	}
	return nil
}

func (m *SheetArchive) GetUsingStartPageNumber() bool {
	if m != nil && m.UsingStartPageNumber != nil {
		return *m.UsingStartPageNumber
	}
	return false
}

func (m *SheetArchive) GetStartPageNumber() int32 {
	if m != nil && m.StartPageNumber != nil {
		return *m.StartPageNumber
	}
	return 0
}

func (m *SheetArchive) GetPageHeaderInset() float32 {
	if m != nil && m.PageHeaderInset != nil {
		return *m.PageHeaderInset
	}
	return 0
}

func (m *SheetArchive) GetPageFooterInset() float32 {
	if m != nil && m.PageFooterInset != nil {
		return *m.PageFooterInset
	}
	return 0
}

func (m *SheetArchive) GetHeaderStorage() *TSP.Reference {
	if m != nil {
		return m.HeaderStorage
	}
	return nil
}

func (m *SheetArchive) GetFooterStorage() *TSP.Reference {
	if m != nil {
		return m.FooterStorage
	}
	return nil
}

func (m *SheetArchive) GetUserDefinedGuideStorage() *TSP.Reference {
	if m != nil {
		return m.UserDefinedGuideStorage
	}
	return nil
}

type FormBasedSheetArchive struct {
	Super            *SheetArchive       `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	TableId          *TSCE.CFUUIDArchive `protobuf:"bytes,2,opt,name=table_id" json:"table_id,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *FormBasedSheetArchive) Reset()         { *m = FormBasedSheetArchive{} }
func (m *FormBasedSheetArchive) String() string { return proto.CompactTextString(m) }
func (*FormBasedSheetArchive) ProtoMessage()    {}

func (m *FormBasedSheetArchive) GetSuper() *SheetArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *FormBasedSheetArchive) GetTableId() *TSCE.CFUUIDArchive {
	if m != nil {
		return m.TableId
	}
	return nil
}

type ThemeArchive struct {
	Super            *TSS.ThemeArchive `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	Prototypes       []*TSP.Reference  `protobuf:"bytes,2,rep,name=prototypes" json:"prototypes,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *ThemeArchive) Reset()         { *m = ThemeArchive{} }
func (m *ThemeArchive) String() string { return proto.CompactTextString(m) }
func (*ThemeArchive) ProtoMessage()    {}

func (m *ThemeArchive) GetSuper() *TSS.ThemeArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *ThemeArchive) GetPrototypes() []*TSP.Reference {
	if m != nil {
		return m.Prototypes
	}
	return nil
}

type ChartMediatorFormulaStorage struct {
	DataFormulae                   []*TSCE.FormulaArchive `protobuf:"bytes,1,rep,name=data_formulae" json:"data_formulae,omitempty"`
	RowLabelFormulae               []*TSCE.FormulaArchive `protobuf:"bytes,3,rep,name=row_label_formulae" json:"row_label_formulae,omitempty"`
	ColLabelFormulae               []*TSCE.FormulaArchive `protobuf:"bytes,4,rep,name=col_label_formulae" json:"col_label_formulae,omitempty"`
	Direction                      *int32                 `protobuf:"varint,5,opt,name=direction" json:"direction,omitempty"`
	ErrorCustomPosFormulae         []*TSCE.FormulaArchive `protobuf:"bytes,6,rep,name=error_custom_pos_formulae" json:"error_custom_pos_formulae,omitempty"`
	ErrorCustomNegFormulae         []*TSCE.FormulaArchive `protobuf:"bytes,7,rep,name=error_custom_neg_formulae" json:"error_custom_neg_formulae,omitempty"`
	ErrorCustomPosScatterXFormulae []*TSCE.FormulaArchive `protobuf:"bytes,8,rep,name=error_custom_pos_scatterX_formulae" json:"error_custom_pos_scatterX_formulae,omitempty"`
	ErrorCustomNegScatterXFormulae []*TSCE.FormulaArchive `protobuf:"bytes,9,rep,name=error_custom_neg_scatterX_formulae" json:"error_custom_neg_scatterX_formulae,omitempty"`
	XXX_unrecognized               []byte                 `json:"-"`
}

func (m *ChartMediatorFormulaStorage) Reset()         { *m = ChartMediatorFormulaStorage{} }
func (m *ChartMediatorFormulaStorage) String() string { return proto.CompactTextString(m) }
func (*ChartMediatorFormulaStorage) ProtoMessage()    {}

func (m *ChartMediatorFormulaStorage) GetDataFormulae() []*TSCE.FormulaArchive {
	if m != nil {
		return m.DataFormulae
	}
	return nil
}

func (m *ChartMediatorFormulaStorage) GetRowLabelFormulae() []*TSCE.FormulaArchive {
	if m != nil {
		return m.RowLabelFormulae
	}
	return nil
}

func (m *ChartMediatorFormulaStorage) GetColLabelFormulae() []*TSCE.FormulaArchive {
	if m != nil {
		return m.ColLabelFormulae
	}
	return nil
}

func (m *ChartMediatorFormulaStorage) GetDirection() int32 {
	if m != nil && m.Direction != nil {
		return *m.Direction
	}
	return 0
}

func (m *ChartMediatorFormulaStorage) GetErrorCustomPosFormulae() []*TSCE.FormulaArchive {
	if m != nil {
		return m.ErrorCustomPosFormulae
	}
	return nil
}

func (m *ChartMediatorFormulaStorage) GetErrorCustomNegFormulae() []*TSCE.FormulaArchive {
	if m != nil {
		return m.ErrorCustomNegFormulae
	}
	return nil
}

func (m *ChartMediatorFormulaStorage) GetErrorCustomPosScatterXFormulae() []*TSCE.FormulaArchive {
	if m != nil {
		return m.ErrorCustomPosScatterXFormulae
	}
	return nil
}

func (m *ChartMediatorFormulaStorage) GetErrorCustomNegScatterXFormulae() []*TSCE.FormulaArchive {
	if m != nil {
		return m.ErrorCustomNegScatterXFormulae
	}
	return nil
}

type ChartMediatorArchive struct {
	Super                      *TSCH2.ChartMediatorArchive  `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	EntityId                   *string                      `protobuf:"bytes,2,req,name=entity_id" json:"entity_id,omitempty"`
	Formulas                   *ChartMediatorFormulaStorage `protobuf:"bytes,3,opt,name=formulas" json:"formulas,omitempty"`
	ColumnsAreSeries           *bool                        `protobuf:"varint,4,opt,name=columns_are_series" json:"columns_are_series,omitempty"`
	IsRegisteredWithCalcEngine *bool                        `protobuf:"varint,5,opt,name=is_registered_with_calc_engine" json:"is_registered_with_calc_engine,omitempty"`
	XXX_unrecognized           []byte                       `json:"-"`
}

func (m *ChartMediatorArchive) Reset()         { *m = ChartMediatorArchive{} }
func (m *ChartMediatorArchive) String() string { return proto.CompactTextString(m) }
func (*ChartMediatorArchive) ProtoMessage()    {}

func (m *ChartMediatorArchive) GetSuper() *TSCH2.ChartMediatorArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *ChartMediatorArchive) GetEntityId() string {
	if m != nil && m.EntityId != nil {
		return *m.EntityId
	}
	return ""
}

func (m *ChartMediatorArchive) GetFormulas() *ChartMediatorFormulaStorage {
	if m != nil {
		return m.Formulas
	}
	return nil
}

func (m *ChartMediatorArchive) GetColumnsAreSeries() bool {
	if m != nil && m.ColumnsAreSeries != nil {
		return *m.ColumnsAreSeries
	}
	return false
}

func (m *ChartMediatorArchive) GetIsRegisteredWithCalcEngine() bool {
	if m != nil && m.IsRegisteredWithCalcEngine != nil {
		return *m.IsRegisteredWithCalcEngine
	}
	return false
}

type ChartSelectionArchive struct {
	Reference        *TSCE.RangeReferenceArchive  `protobuf:"bytes,1,opt,name=reference" json:"reference,omitempty"`
	Super            *TSCH2.ChartSelectionArchive `protobuf:"bytes,2,opt,name=super" json:"super,omitempty"`
	XXX_unrecognized []byte                       `json:"-"`
}

func (m *ChartSelectionArchive) Reset()         { *m = ChartSelectionArchive{} }
func (m *ChartSelectionArchive) String() string { return proto.CompactTextString(m) }
func (*ChartSelectionArchive) ProtoMessage()    {}

func (m *ChartSelectionArchive) GetReference() *TSCE.RangeReferenceArchive {
	if m != nil {
		return m.Reference
	}
	return nil
}

func (m *ChartSelectionArchive) GetSuper() *TSCH2.ChartSelectionArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func init() {
	proto.RegisterEnum("TN.SheetPageOrder", SheetPageOrder_name, SheetPageOrder_value)
	proto.RegisterEnum("TN.UIStateArchive_InspectorPaneViewMode", UIStateArchive_InspectorPaneViewMode_name, UIStateArchive_InspectorPaneViewMode_value)
}
