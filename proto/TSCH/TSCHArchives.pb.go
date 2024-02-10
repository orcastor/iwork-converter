// Code generated by protoc-gen-go.
// source: TSCHArchives.proto
// DO NOT EDIT!

/*
Package TSCH is a generated protocol buffer package.

It is generated from these files:
	TSCHArchives.proto

It has these top-level messages:
	ChartDrawableArchive
	ChartArchive
	ChartPasteboardAdditionsArchive
	ChartGridArchive
	ChartMediatorArchive
	ChartStylePreset
	ChartPresetsArchive
	PropertyValueStorageContainerArchive
	StylePasteboardDataArchive
	ChartSelectionPathTypeArchive
	ChartAxisIDArchive
	ChartSelectionPathArgumentArchive
	ChartSelectionPathArchive
	ChartSelectionArchive
	ChartUIState
	ChartFormatStructExtensions
*/
package TSCH

import proto "github.com/golang/protobuf/proto"
import math "math"
import "github.com/orcastor/iwork-converter/proto/TSP"
import "github.com/orcastor/iwork-converter/proto/TSK"
import "github.com/orcastor/iwork-converter/proto/TSD"
import "github.com/orcastor/iwork-converter/proto/TSS"

// discarding unused import TSCH_Generated "TSCHArchives.GEN.pb"
// discarding unused import TSCH1 "TSCH3DArchives.pb"
// discarding unused import TSCH_PreUFF "TSCHPreUFFArchives.pb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type ChartDrawableArchive struct {
	Super            *TSD.DrawableArchive      `protobuf:"bytes,1,opt,name=super" json:"super,omitempty"`
	XXX_extensions   map[int32]proto.Extension `json:"-"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *ChartDrawableArchive) Reset()         { *m = ChartDrawableArchive{} }
func (m *ChartDrawableArchive) String() string { return proto.CompactTextString(m) }
func (*ChartDrawableArchive) ProtoMessage()    {}

var extRange_ChartDrawableArchive = []proto.ExtensionRange{
	{10000, 536870911},
}

func (*ChartDrawableArchive) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_ChartDrawableArchive
}
func (m *ChartDrawableArchive) ExtensionMap() map[int32]proto.Extension {
	if m.XXX_extensions == nil {
		m.XXX_extensions = make(map[int32]proto.Extension)
	}
	return m.XXX_extensions
}

func (m *ChartDrawableArchive) GetSuper() *TSD.DrawableArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

type ChartArchive struct {
	ChartType                           *ChartType                   `protobuf:"varint,1,opt,name=chart_type,enum=ChartType" json:"chart_type,omitempty"`
	ScatterFormat                       *ScatterFormat               `protobuf:"varint,2,opt,name=scatter_format,enum=ScatterFormat" json:"scatter_format,omitempty"`
	LegendFrame                         *RectArchive                 `protobuf:"bytes,3,opt,name=legend_frame" json:"legend_frame,omitempty"`
	Preset                              *TSP.Reference               `protobuf:"bytes,4,opt,name=preset" json:"preset,omitempty"`
	SeriesDirection                     *SeriesDirection             `protobuf:"varint,5,opt,name=series_direction,enum=SeriesDirection" json:"series_direction,omitempty"`
	ContainsDefaultData                 *bool                        `protobuf:"varint,6,opt,name=contains_default_data" json:"contains_default_data,omitempty"`
	Grid                                *ChartGridArchive            `protobuf:"bytes,7,opt,name=grid" json:"grid,omitempty"`
	Mediator                            *TSP.Reference               `protobuf:"bytes,8,opt,name=mediator" json:"mediator,omitempty"`
	ChartStyle                          *TSP.Reference               `protobuf:"bytes,9,opt,name=chart_style" json:"chart_style,omitempty"`
	ChartNonStyle                       *TSP.Reference               `protobuf:"bytes,10,opt,name=chart_non_style" json:"chart_non_style,omitempty"`
	LegendStyle                         *TSP.Reference               `protobuf:"bytes,11,opt,name=legend_style" json:"legend_style,omitempty"`
	LegendNonStyle                      *TSP.Reference               `protobuf:"bytes,12,opt,name=legend_non_style" json:"legend_non_style,omitempty"`
	ValueAxisStyles                     []*TSP.Reference             `protobuf:"bytes,13,rep,name=value_axis_styles" json:"value_axis_styles,omitempty"`
	ValueAxisNonstyles                  []*TSP.Reference             `protobuf:"bytes,14,rep,name=value_axis_nonstyles" json:"value_axis_nonstyles,omitempty"`
	CategoryAxisStyles                  []*TSP.Reference             `protobuf:"bytes,15,rep,name=category_axis_styles" json:"category_axis_styles,omitempty"`
	CategoryAxisNonstyles               []*TSP.Reference             `protobuf:"bytes,16,rep,name=category_axis_nonstyles" json:"category_axis_nonstyles,omitempty"`
	SeriesThemeStyles                   []*TSP.Reference             `protobuf:"bytes,17,rep,name=series_theme_styles" json:"series_theme_styles,omitempty"`
	SeriesPrivateStyles                 *SparseReferenceArrayArchive `protobuf:"bytes,18,opt,name=series_private_styles" json:"series_private_styles,omitempty"`
	SeriesNonStyles                     *SparseReferenceArrayArchive `protobuf:"bytes,19,opt,name=series_non_styles" json:"series_non_styles,omitempty"`
	ParagraphStyles                     []*TSP.Reference             `protobuf:"bytes,20,rep,name=paragraph_styles" json:"paragraph_styles,omitempty"`
	MultidatasetIndex                   *uint32                      `protobuf:"varint,21,opt,name=multidataset_index" json:"multidataset_index,omitempty"`
	NeedsCalcEngineDeferredImportAction *bool                        `protobuf:"varint,22,opt,name=needs_calc_engine_deferred_import_action" json:"needs_calc_engine_deferred_import_action,omitempty"`
	XXX_extensions                      map[int32]proto.Extension    `json:"-"`
	XXX_unrecognized                    []byte                       `json:"-"`
}

func (m *ChartArchive) Reset()         { *m = ChartArchive{} }
func (m *ChartArchive) String() string { return proto.CompactTextString(m) }
func (*ChartArchive) ProtoMessage()    {}

var extRange_ChartArchive = []proto.ExtensionRange{
	{10000, 536870911},
}

func (*ChartArchive) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_ChartArchive
}
func (m *ChartArchive) ExtensionMap() map[int32]proto.Extension {
	if m.XXX_extensions == nil {
		m.XXX_extensions = make(map[int32]proto.Extension)
	}
	return m.XXX_extensions
}

func (m *ChartArchive) GetChartType() ChartType {
	if m != nil && m.ChartType != nil {
		return *m.ChartType
	}
	return ChartType_undefinedChartType
}

func (m *ChartArchive) GetScatterFormat() ScatterFormat {
	if m != nil && m.ScatterFormat != nil {
		return *m.ScatterFormat
	}
	return ScatterFormat_scatter_format_unknown
}

func (m *ChartArchive) GetLegendFrame() *RectArchive {
	if m != nil {
		return m.LegendFrame
	}
	return nil
}

func (m *ChartArchive) GetPreset() *TSP.Reference {
	if m != nil {
		return m.Preset
	}
	return nil
}

func (m *ChartArchive) GetSeriesDirection() SeriesDirection {
	if m != nil && m.SeriesDirection != nil {
		return *m.SeriesDirection
	}
	return SeriesDirection_series_direction_unknown
}

func (m *ChartArchive) GetContainsDefaultData() bool {
	if m != nil && m.ContainsDefaultData != nil {
		return *m.ContainsDefaultData
	}
	return false
}

func (m *ChartArchive) GetGrid() *ChartGridArchive {
	if m != nil {
		return m.Grid
	}
	return nil
}

func (m *ChartArchive) GetMediator() *TSP.Reference {
	if m != nil {
		return m.Mediator
	}
	return nil
}

func (m *ChartArchive) GetChartStyle() *TSP.Reference {
	if m != nil {
		return m.ChartStyle
	}
	return nil
}

func (m *ChartArchive) GetChartNonStyle() *TSP.Reference {
	if m != nil {
		return m.ChartNonStyle
	}
	return nil
}

func (m *ChartArchive) GetLegendStyle() *TSP.Reference {
	if m != nil {
		return m.LegendStyle
	}
	return nil
}

func (m *ChartArchive) GetLegendNonStyle() *TSP.Reference {
	if m != nil {
		return m.LegendNonStyle
	}
	return nil
}

func (m *ChartArchive) GetValueAxisStyles() []*TSP.Reference {
	if m != nil {
		return m.ValueAxisStyles
	}
	return nil
}

func (m *ChartArchive) GetValueAxisNonstyles() []*TSP.Reference {
	if m != nil {
		return m.ValueAxisNonstyles
	}
	return nil
}

func (m *ChartArchive) GetCategoryAxisStyles() []*TSP.Reference {
	if m != nil {
		return m.CategoryAxisStyles
	}
	return nil
}

func (m *ChartArchive) GetCategoryAxisNonstyles() []*TSP.Reference {
	if m != nil {
		return m.CategoryAxisNonstyles
	}
	return nil
}

func (m *ChartArchive) GetSeriesThemeStyles() []*TSP.Reference {
	if m != nil {
		return m.SeriesThemeStyles
	}
	return nil
}

func (m *ChartArchive) GetSeriesPrivateStyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.SeriesPrivateStyles
	}
	return nil
}

func (m *ChartArchive) GetSeriesNonStyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.SeriesNonStyles
	}
	return nil
}

func (m *ChartArchive) GetParagraphStyles() []*TSP.Reference {
	if m != nil {
		return m.ParagraphStyles
	}
	return nil
}

func (m *ChartArchive) GetMultidatasetIndex() uint32 {
	if m != nil && m.MultidatasetIndex != nil {
		return *m.MultidatasetIndex
	}
	return 0
}

func (m *ChartArchive) GetNeedsCalcEngineDeferredImportAction() bool {
	if m != nil && m.NeedsCalcEngineDeferredImportAction != nil {
		return *m.NeedsCalcEngineDeferredImportAction
	}
	return false
}

var E_ChartArchive_Unity = &proto.ExtensionDesc{
	ExtendedType:  (*ChartDrawableArchive)(nil),
	ExtensionType: (*ChartArchive)(nil),
	Field:         10000,
	Name:          "ChartArchive.unity",
	Tag:           "bytes,10000,opt,name=unity",
}

type ChartPasteboardAdditionsArchive struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *ChartPasteboardAdditionsArchive) Reset()         { *m = ChartPasteboardAdditionsArchive{} }
func (m *ChartPasteboardAdditionsArchive) String() string { return proto.CompactTextString(m) }
func (*ChartPasteboardAdditionsArchive) ProtoMessage()    {}

var E_ChartPasteboardAdditionsArchive_PresetIndexForPasteboard = &proto.ExtensionDesc{
	ExtendedType:  (*ChartArchive)(nil),
	ExtensionType: (*uint32)(nil),
	Field:         10000,
	Name:          "ChartPasteboardAdditionsArchive.preset_index_for_pasteboard",
	Tag:           "varint,10000,opt,name=preset_index_for_pasteboard",
}

var E_ChartPasteboardAdditionsArchive_PresetUuidForPasteboard = &proto.ExtensionDesc{
	ExtendedType:  (*ChartArchive)(nil),
	ExtensionType: ([]byte)(nil),
	Field:         10001,
	Name:          "ChartPasteboardAdditionsArchive.preset_uuid_for_pasteboard",
	Tag:           "bytes,10001,opt,name=preset_uuid_for_pasteboard",
}

type ChartGridArchive struct {
	RowName          []string                    `protobuf:"bytes,1,rep,name=row_name" json:"row_name,omitempty"`
	ColumnName       []string                    `protobuf:"bytes,2,rep,name=column_name" json:"column_name,omitempty"`
	GridRow          []*ChartGridArchive_GridRow `protobuf:"bytes,3,rep,name=grid_row" json:"grid_row,omitempty"`
	XXX_unrecognized []byte                      `json:"-"`
}

func (m *ChartGridArchive) Reset()         { *m = ChartGridArchive{} }
func (m *ChartGridArchive) String() string { return proto.CompactTextString(m) }
func (*ChartGridArchive) ProtoMessage()    {}

func (m *ChartGridArchive) GetRowName() []string {
	if m != nil {
		return m.RowName
	}
	return nil
}

func (m *ChartGridArchive) GetColumnName() []string {
	if m != nil {
		return m.ColumnName
	}
	return nil
}

func (m *ChartGridArchive) GetGridRow() []*ChartGridArchive_GridRow {
	if m != nil {
		return m.GridRow
	}
	return nil
}

type ChartGridArchive_GridRow struct {
	Value            []*ChartGridArchive_GridRow_GridValue `protobuf:"bytes,1,rep,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte                                `json:"-"`
}

func (m *ChartGridArchive_GridRow) Reset()         { *m = ChartGridArchive_GridRow{} }
func (m *ChartGridArchive_GridRow) String() string { return proto.CompactTextString(m) }
func (*ChartGridArchive_GridRow) ProtoMessage()    {}

func (m *ChartGridArchive_GridRow) GetValue() []*ChartGridArchive_GridRow_GridValue {
	if m != nil {
		return m.Value
	}
	return nil
}

type ChartGridArchive_GridRow_GridValue struct {
	NumericValue     *float64 `protobuf:"fixed64,1,opt,name=numeric_value" json:"numeric_value,omitempty"`
	DateValue        *float64 `protobuf:"fixed64,2,opt,name=date_value" json:"date_value,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *ChartGridArchive_GridRow_GridValue) Reset()         { *m = ChartGridArchive_GridRow_GridValue{} }
func (m *ChartGridArchive_GridRow_GridValue) String() string { return proto.CompactTextString(m) }
func (*ChartGridArchive_GridRow_GridValue) ProtoMessage()    {}

func (m *ChartGridArchive_GridRow_GridValue) GetNumericValue() float64 {
	if m != nil && m.NumericValue != nil {
		return *m.NumericValue
	}
	return 0
}

func (m *ChartGridArchive_GridRow_GridValue) GetDateValue() float64 {
	if m != nil && m.DateValue != nil {
		return *m.DateValue
	}
	return 0
}

type ChartMediatorArchive struct {
	Info                *TSP.Reference `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
	LocalSeriesIndexes  []uint32       `protobuf:"varint,2,rep,name=local_series_indexes" json:"local_series_indexes,omitempty"`
	RemoteSeriesIndexes []uint32       `protobuf:"varint,3,rep,name=remote_series_indexes" json:"remote_series_indexes,omitempty"`
	XXX_unrecognized    []byte         `json:"-"`
}

func (m *ChartMediatorArchive) Reset()         { *m = ChartMediatorArchive{} }
func (m *ChartMediatorArchive) String() string { return proto.CompactTextString(m) }
func (*ChartMediatorArchive) ProtoMessage()    {}

func (m *ChartMediatorArchive) GetInfo() *TSP.Reference {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *ChartMediatorArchive) GetLocalSeriesIndexes() []uint32 {
	if m != nil {
		return m.LocalSeriesIndexes
	}
	return nil
}

func (m *ChartMediatorArchive) GetRemoteSeriesIndexes() []uint32 {
	if m != nil {
		return m.RemoteSeriesIndexes
	}
	return nil
}

type ChartStylePreset struct {
	ChartStyle         *TSP.Reference   `protobuf:"bytes,1,opt,name=chart_style" json:"chart_style,omitempty"`
	LegendStyle        *TSP.Reference   `protobuf:"bytes,2,opt,name=legend_style" json:"legend_style,omitempty"`
	ValueAxisStyles    []*TSP.Reference `protobuf:"bytes,3,rep,name=value_axis_styles" json:"value_axis_styles,omitempty"`
	CategoryAxisStyles []*TSP.Reference `protobuf:"bytes,4,rep,name=category_axis_styles" json:"category_axis_styles,omitempty"`
	SeriesStyles       []*TSP.Reference `protobuf:"bytes,5,rep,name=series_styles" json:"series_styles,omitempty"`
	ParagraphStyles    []*TSP.Reference `protobuf:"bytes,6,rep,name=paragraph_styles" json:"paragraph_styles,omitempty"`
	Uuid               []byte           `protobuf:"bytes,7,opt,name=uuid" json:"uuid,omitempty"`
	XXX_unrecognized   []byte           `json:"-"`
}

func (m *ChartStylePreset) Reset()         { *m = ChartStylePreset{} }
func (m *ChartStylePreset) String() string { return proto.CompactTextString(m) }
func (*ChartStylePreset) ProtoMessage()    {}

func (m *ChartStylePreset) GetChartStyle() *TSP.Reference {
	if m != nil {
		return m.ChartStyle
	}
	return nil
}

func (m *ChartStylePreset) GetLegendStyle() *TSP.Reference {
	if m != nil {
		return m.LegendStyle
	}
	return nil
}

func (m *ChartStylePreset) GetValueAxisStyles() []*TSP.Reference {
	if m != nil {
		return m.ValueAxisStyles
	}
	return nil
}

func (m *ChartStylePreset) GetCategoryAxisStyles() []*TSP.Reference {
	if m != nil {
		return m.CategoryAxisStyles
	}
	return nil
}

func (m *ChartStylePreset) GetSeriesStyles() []*TSP.Reference {
	if m != nil {
		return m.SeriesStyles
	}
	return nil
}

func (m *ChartStylePreset) GetParagraphStyles() []*TSP.Reference {
	if m != nil {
		return m.ParagraphStyles
	}
	return nil
}

func (m *ChartStylePreset) GetUuid() []byte {
	if m != nil {
		return m.Uuid
	}
	return nil
}

type ChartPresetsArchive struct {
	ChartPresets     []*TSP.Reference `protobuf:"bytes,1,rep,name=chart_presets" json:"chart_presets,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *ChartPresetsArchive) Reset()         { *m = ChartPresetsArchive{} }
func (m *ChartPresetsArchive) String() string { return proto.CompactTextString(m) }
func (*ChartPresetsArchive) ProtoMessage()    {}

func (m *ChartPresetsArchive) GetChartPresets() []*TSP.Reference {
	if m != nil {
		return m.ChartPresets
	}
	return nil
}

var E_ChartPresetsArchive_Extension = &proto.ExtensionDesc{
	ExtendedType:  (*TSS.ThemeArchive)(nil),
	ExtensionType: (*ChartPresetsArchive)(nil),
	Field:         120,
	Name:          "ChartPresetsArchive.extension",
	Tag:           "bytes,120,opt,name=extension",
}

type PropertyValueStorageContainerArchive struct {
	ChartStyle            *TSP.Reference               `protobuf:"bytes,1,opt,name=chart_style" json:"chart_style,omitempty"`
	ChartNonstyle         *TSP.Reference               `protobuf:"bytes,2,opt,name=chart_nonstyle" json:"chart_nonstyle,omitempty"`
	LegendStyle           *TSP.Reference               `protobuf:"bytes,3,opt,name=legend_style" json:"legend_style,omitempty"`
	LegendNonstyle        *TSP.Reference               `protobuf:"bytes,4,opt,name=legend_nonstyle" json:"legend_nonstyle,omitempty"`
	ValueAxisStyles       *SparseReferenceArrayArchive `protobuf:"bytes,5,opt,name=value_axis_styles" json:"value_axis_styles,omitempty"`
	ValueAxisNonstyles    *SparseReferenceArrayArchive `protobuf:"bytes,6,opt,name=value_axis_nonstyles" json:"value_axis_nonstyles,omitempty"`
	CategoryAxisStyles    *SparseReferenceArrayArchive `protobuf:"bytes,7,opt,name=category_axis_styles" json:"category_axis_styles,omitempty"`
	CategoryAxisNonstyles *SparseReferenceArrayArchive `protobuf:"bytes,8,opt,name=category_axis_nonstyles" json:"category_axis_nonstyles,omitempty"`
	SeriesThemeStyles     *SparseReferenceArrayArchive `protobuf:"bytes,9,opt,name=series_theme_styles" json:"series_theme_styles,omitempty"`
	SeriesPrivateStyles   *SparseReferenceArrayArchive `protobuf:"bytes,10,opt,name=series_private_styles" json:"series_private_styles,omitempty"`
	SeriesNonstyles       *SparseReferenceArrayArchive `protobuf:"bytes,11,opt,name=series_nonstyles" json:"series_nonstyles,omitempty"`
	ParagraphStyles       *SparseReferenceArrayArchive `protobuf:"bytes,12,opt,name=paragraph_styles" json:"paragraph_styles,omitempty"`
	XXX_unrecognized      []byte                       `json:"-"`
}

func (m *PropertyValueStorageContainerArchive) Reset()         { *m = PropertyValueStorageContainerArchive{} }
func (m *PropertyValueStorageContainerArchive) String() string { return proto.CompactTextString(m) }
func (*PropertyValueStorageContainerArchive) ProtoMessage()    {}

func (m *PropertyValueStorageContainerArchive) GetChartStyle() *TSP.Reference {
	if m != nil {
		return m.ChartStyle
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetChartNonstyle() *TSP.Reference {
	if m != nil {
		return m.ChartNonstyle
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetLegendStyle() *TSP.Reference {
	if m != nil {
		return m.LegendStyle
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetLegendNonstyle() *TSP.Reference {
	if m != nil {
		return m.LegendNonstyle
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetValueAxisStyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.ValueAxisStyles
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetValueAxisNonstyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.ValueAxisNonstyles
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetCategoryAxisStyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.CategoryAxisStyles
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetCategoryAxisNonstyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.CategoryAxisNonstyles
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetSeriesThemeStyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.SeriesThemeStyles
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetSeriesPrivateStyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.SeriesPrivateStyles
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetSeriesNonstyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.SeriesNonstyles
	}
	return nil
}

func (m *PropertyValueStorageContainerArchive) GetParagraphStyles() *SparseReferenceArrayArchive {
	if m != nil {
		return m.ParagraphStyles
	}
	return nil
}

type StylePasteboardDataArchive struct {
	Super                 *TSS.StyleArchive                     `protobuf:"bytes,1,opt,name=super" json:"super,omitempty"`
	StyleNetwork          *PropertyValueStorageContainerArchive `protobuf:"bytes,2,opt,name=style_network" json:"style_network,omitempty"`
	CopiedFromEntireChart *bool                                 `protobuf:"varint,3,opt,name=copied_from_entire_chart" json:"copied_from_entire_chart,omitempty"`
	XXX_unrecognized      []byte                                `json:"-"`
}

func (m *StylePasteboardDataArchive) Reset()         { *m = StylePasteboardDataArchive{} }
func (m *StylePasteboardDataArchive) String() string { return proto.CompactTextString(m) }
func (*StylePasteboardDataArchive) ProtoMessage()    {}

func (m *StylePasteboardDataArchive) GetSuper() *TSS.StyleArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *StylePasteboardDataArchive) GetStyleNetwork() *PropertyValueStorageContainerArchive {
	if m != nil {
		return m.StyleNetwork
	}
	return nil
}

func (m *StylePasteboardDataArchive) GetCopiedFromEntireChart() bool {
	if m != nil && m.CopiedFromEntireChart != nil {
		return *m.CopiedFromEntireChart
	}
	return false
}

type ChartSelectionPathTypeArchive struct {
	PathType         *string `protobuf:"bytes,1,opt,name=path_type" json:"path_type,omitempty"`
	PathName         *string `protobuf:"bytes,2,opt,name=path_name" json:"path_name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ChartSelectionPathTypeArchive) Reset()         { *m = ChartSelectionPathTypeArchive{} }
func (m *ChartSelectionPathTypeArchive) String() string { return proto.CompactTextString(m) }
func (*ChartSelectionPathTypeArchive) ProtoMessage()    {}

func (m *ChartSelectionPathTypeArchive) GetPathType() string {
	if m != nil && m.PathType != nil {
		return *m.PathType
	}
	return ""
}

func (m *ChartSelectionPathTypeArchive) GetPathName() string {
	if m != nil && m.PathName != nil {
		return *m.PathName
	}
	return ""
}

type ChartAxisIDArchive struct {
	AxisType         *AxisType `protobuf:"varint,1,opt,name=axis_type,enum=AxisType" json:"axis_type,omitempty"`
	Ordinal          *uint32   `protobuf:"varint,2,opt,name=ordinal" json:"ordinal,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *ChartAxisIDArchive) Reset()         { *m = ChartAxisIDArchive{} }
func (m *ChartAxisIDArchive) String() string { return proto.CompactTextString(m) }
func (*ChartAxisIDArchive) ProtoMessage()    {}

func (m *ChartAxisIDArchive) GetAxisType() AxisType {
	if m != nil && m.AxisType != nil {
		return *m.AxisType
	}
	return AxisType_axis_type_unknown
}

func (m *ChartAxisIDArchive) GetOrdinal() uint32 {
	if m != nil && m.Ordinal != nil {
		return *m.Ordinal
	}
	return 0
}

type ChartSelectionPathArgumentArchive struct {
	Number           *uint32             `protobuf:"varint,1,opt,name=number" json:"number,omitempty"`
	AxisId           *ChartAxisIDArchive `protobuf:"bytes,2,opt,name=axis_id" json:"axis_id,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *ChartSelectionPathArgumentArchive) Reset()         { *m = ChartSelectionPathArgumentArchive{} }
func (m *ChartSelectionPathArgumentArchive) String() string { return proto.CompactTextString(m) }
func (*ChartSelectionPathArgumentArchive) ProtoMessage()    {}

func (m *ChartSelectionPathArgumentArchive) GetNumber() uint32 {
	if m != nil && m.Number != nil {
		return *m.Number
	}
	return 0
}

func (m *ChartSelectionPathArgumentArchive) GetAxisId() *ChartAxisIDArchive {
	if m != nil {
		return m.AxisId
	}
	return nil
}

type ChartSelectionPathArchive struct {
	PathType         *ChartSelectionPathTypeArchive       `protobuf:"bytes,1,opt,name=path_type" json:"path_type,omitempty"`
	SubSelection     *ChartSelectionPathArchive           `protobuf:"bytes,2,opt,name=sub_selection" json:"sub_selection,omitempty"`
	Arguments        []*ChartSelectionPathArgumentArchive `protobuf:"bytes,3,rep,name=arguments" json:"arguments,omitempty"`
	XXX_unrecognized []byte                               `json:"-"`
}

func (m *ChartSelectionPathArchive) Reset()         { *m = ChartSelectionPathArchive{} }
func (m *ChartSelectionPathArchive) String() string { return proto.CompactTextString(m) }
func (*ChartSelectionPathArchive) ProtoMessage()    {}

func (m *ChartSelectionPathArchive) GetPathType() *ChartSelectionPathTypeArchive {
	if m != nil {
		return m.PathType
	}
	return nil
}

func (m *ChartSelectionPathArchive) GetSubSelection() *ChartSelectionPathArchive {
	if m != nil {
		return m.SubSelection
	}
	return nil
}

func (m *ChartSelectionPathArchive) GetArguments() []*ChartSelectionPathArgumentArchive {
	if m != nil {
		return m.Arguments
	}
	return nil
}

type ChartSelectionArchive struct {
	Chart            *TSP.Reference               `protobuf:"bytes,1,opt,name=chart" json:"chart,omitempty"`
	Paths            []*ChartSelectionPathArchive `protobuf:"bytes,2,rep,name=paths" json:"paths,omitempty"`
	XXX_unrecognized []byte                       `json:"-"`
}

func (m *ChartSelectionArchive) Reset()         { *m = ChartSelectionArchive{} }
func (m *ChartSelectionArchive) String() string { return proto.CompactTextString(m) }
func (*ChartSelectionArchive) ProtoMessage()    {}

func (m *ChartSelectionArchive) GetChart() *TSP.Reference {
	if m != nil {
		return m.Chart
	}
	return nil
}

func (m *ChartSelectionArchive) GetPaths() []*ChartSelectionPathArchive {
	if m != nil {
		return m.Paths
	}
	return nil
}

type ChartUIState struct {
	Chart              *TSP.Reference `protobuf:"bytes,1,opt,name=chart" json:"chart,omitempty"`
	CdeLastRowSelected *int32         `protobuf:"varint,2,opt,name=cde_last_row_selected" json:"cde_last_row_selected,omitempty"`
	CdeLastColSelected *int32         `protobuf:"varint,3,opt,name=cde_last_col_selected" json:"cde_last_col_selected,omitempty"`
	CdeLastRowCount    *int32         `protobuf:"varint,4,opt,name=cde_last_row_count" json:"cde_last_row_count,omitempty"`
	CdeLastColCount    *int32         `protobuf:"varint,5,opt,name=cde_last_col_count" json:"cde_last_col_count,omitempty"`
	XXX_unrecognized   []byte         `json:"-"`
}

func (m *ChartUIState) Reset()         { *m = ChartUIState{} }
func (m *ChartUIState) String() string { return proto.CompactTextString(m) }
func (*ChartUIState) ProtoMessage()    {}

func (m *ChartUIState) GetChart() *TSP.Reference {
	if m != nil {
		return m.Chart
	}
	return nil
}

func (m *ChartUIState) GetCdeLastRowSelected() int32 {
	if m != nil && m.CdeLastRowSelected != nil {
		return *m.CdeLastRowSelected
	}
	return 0
}

func (m *ChartUIState) GetCdeLastColSelected() int32 {
	if m != nil && m.CdeLastColSelected != nil {
		return *m.CdeLastColSelected
	}
	return 0
}

func (m *ChartUIState) GetCdeLastRowCount() int32 {
	if m != nil && m.CdeLastRowCount != nil {
		return *m.CdeLastRowCount
	}
	return 0
}

func (m *ChartUIState) GetCdeLastColCount() int32 {
	if m != nil && m.CdeLastColCount != nil {
		return *m.CdeLastColCount
	}
	return 0
}

type ChartFormatStructExtensions struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *ChartFormatStructExtensions) Reset()         { *m = ChartFormatStructExtensions{} }
func (m *ChartFormatStructExtensions) String() string { return proto.CompactTextString(m) }
func (*ChartFormatStructExtensions) ProtoMessage()    {}

var E_ChartFormatStructExtensions_Prefix = &proto.ExtensionDesc{
	ExtendedType:  (*TSK.FormatStructArchive)(nil),
	ExtensionType: (*string)(nil),
	Field:         10000,
	Name:          "ChartFormatStructExtensions.prefix",
	Tag:           "bytes,10000,opt,name=prefix",
}

var E_ChartFormatStructExtensions_Suffix = &proto.ExtensionDesc{
	ExtendedType:  (*TSK.FormatStructArchive)(nil),
	ExtensionType: (*string)(nil),
	Field:         10001,
	Name:          "ChartFormatStructExtensions.suffix",
	Tag:           "bytes,10001,opt,name=suffix",
}

var E_Scene3DSettingsConstantDepth = &proto.ExtensionDesc{
	ExtendedType:  (*ChartArchive)(nil),
	ExtensionType: (*bool)(nil),
	Field:         10002,
	Name:          "scene3d_settings_constant_depth",
	Tag:           "varint,10002,opt,name=scene3d_settings_constant_depth",
}

func init() {
	proto.RegisterExtension(E_ChartArchive_Unity)
	proto.RegisterExtension(E_ChartPasteboardAdditionsArchive_PresetIndexForPasteboard)
	proto.RegisterExtension(E_ChartPasteboardAdditionsArchive_PresetUuidForPasteboard)
	proto.RegisterExtension(E_ChartPresetsArchive_Extension)
	proto.RegisterExtension(E_ChartFormatStructExtensions_Prefix)
	proto.RegisterExtension(E_ChartFormatStructExtensions_Suffix)
	proto.RegisterExtension(E_Scene3DSettingsConstantDepth)
}
