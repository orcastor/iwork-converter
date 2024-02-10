// Code generated by protoc-gen-go.
// source: TSKArchives.proto
// DO NOT EDIT!

/*
Package TSK is a generated protocol buffer package.

It is generated from these files:
	TSKArchives.proto

It has these top-level messages:
	TreeNode
	CommandHistory
	DocumentArchive
	DocumentSupportArchive
	ViewStateArchive
	CommandArchive
	CommandGroupArchive
	CommandContainerArchive
	ReplaceAllChildCommandArchive
	ReplaceAllCommandArchive
	ShuffleMappingArchive
	ProgressiveCommandGroupArchive
	CommandSelectionBehaviorHistoryArchive
	UndoRedoStateCommandSelectionBehaviorArchive
	FormatStructArchive
	CustomFormatArchive
	AnnotationAuthorArchive
	DeprecatedChangeAuthorArchive
	AnnotationAuthorStorageArchive
	AddAnnotationAuthorCommandArchive
	SetAnnotationAuthorColorCommandArchive
*/
package TSK

import proto "github.com/golang/protobuf/proto"
import math "math"
import "github.com/orcastor/iwork-converter/proto/TSP"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type TreeNode struct {
	Name             *string          `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Children         []*TSP.Reference `protobuf:"bytes,2,rep,name=children" json:"children,omitempty"`
	Object           *TSP.Reference   `protobuf:"bytes,3,opt,name=object" json:"object,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *TreeNode) Reset()         { *m = TreeNode{} }
func (m *TreeNode) String() string { return proto.CompactTextString(m) }
func (*TreeNode) ProtoMessage()    {}

func (m *TreeNode) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *TreeNode) GetChildren() []*TSP.Reference {
	if m != nil {
		return m.Children
	}
	return nil
}

func (m *TreeNode) GetObject() *TSP.Reference {
	if m != nil {
		return m.Object
	}
	return nil
}

type CommandHistory struct {
	UndoCount               *uint32          `protobuf:"varint,1,req,name=undo_count" json:"undo_count,omitempty"`
	Commands                []*TSP.Reference `protobuf:"bytes,2,rep,name=commands" json:"commands,omitempty"`
	MarkedRedoCommands      []*TSP.Reference `protobuf:"bytes,3,rep,name=marked_redo_commands" json:"marked_redo_commands,omitempty"`
	PendingPreflightCommand *TSP.Reference   `protobuf:"bytes,4,opt,name=pending_preflight_command" json:"pending_preflight_command,omitempty"`
	FixedRadar_13365177     *bool            `protobuf:"varint,10,opt,name=fixed_radar_13365177" json:"fixed_radar_13365177,omitempty"`
	XXX_unrecognized        []byte           `json:"-"`
}

func (m *CommandHistory) Reset()         { *m = CommandHistory{} }
func (m *CommandHistory) String() string { return proto.CompactTextString(m) }
func (*CommandHistory) ProtoMessage()    {}

func (m *CommandHistory) GetUndoCount() uint32 {
	if m != nil && m.UndoCount != nil {
		return *m.UndoCount
	}
	return 0
}

func (m *CommandHistory) GetCommands() []*TSP.Reference {
	if m != nil {
		return m.Commands
	}
	return nil
}

func (m *CommandHistory) GetMarkedRedoCommands() []*TSP.Reference {
	if m != nil {
		return m.MarkedRedoCommands
	}
	return nil
}

func (m *CommandHistory) GetPendingPreflightCommand() *TSP.Reference {
	if m != nil {
		return m.PendingPreflightCommand
	}
	return nil
}

func (m *CommandHistory) GetFixedRadar_13365177() bool {
	if m != nil && m.FixedRadar_13365177 != nil {
		return *m.FixedRadar_13365177
	}
	return false
}

type DocumentArchive struct {
	LocaleIdentifier        *string        `protobuf:"bytes,4,opt,name=locale_identifier" json:"locale_identifier,omitempty"`
	AnnotationAuthorStorage *TSP.Reference `protobuf:"bytes,7,opt,name=annotation_author_storage" json:"annotation_author_storage,omitempty"`
	XXX_unrecognized        []byte         `json:"-"`
}

func (m *DocumentArchive) Reset()         { *m = DocumentArchive{} }
func (m *DocumentArchive) String() string { return proto.CompactTextString(m) }
func (*DocumentArchive) ProtoMessage()    {}

func (m *DocumentArchive) GetLocaleIdentifier() string {
	if m != nil && m.LocaleIdentifier != nil {
		return *m.LocaleIdentifier
	}
	return ""
}

func (m *DocumentArchive) GetAnnotationAuthorStorage() *TSP.Reference {
	if m != nil {
		return m.AnnotationAuthorStorage
	}
	return nil
}

type DocumentSupportArchive struct {
	CommandHistory                  *TSP.Reference `protobuf:"bytes,1,opt,name=command_history" json:"command_history,omitempty"`
	CommandSelectionBehaviorHistory *TSP.Reference `protobuf:"bytes,2,opt,name=command_selection_behavior_history" json:"command_selection_behavior_history,omitempty"`
	UndoCount                       *uint32        `protobuf:"varint,4,opt,name=undo_count" json:"undo_count,omitempty"`
	RedoCount                       *uint32        `protobuf:"varint,5,opt,name=redo_count" json:"redo_count,omitempty"`
	UndoActionString                *string        `protobuf:"bytes,6,opt,name=undo_action_string" json:"undo_action_string,omitempty"`
	RedoActionString                *string        `protobuf:"bytes,7,opt,name=redo_action_string" json:"redo_action_string,omitempty"`
	WebState                        *TSP.Reference `protobuf:"bytes,8,opt,name=web_state" json:"web_state,omitempty"`
	XXX_unrecognized                []byte         `json:"-"`
}

func (m *DocumentSupportArchive) Reset()         { *m = DocumentSupportArchive{} }
func (m *DocumentSupportArchive) String() string { return proto.CompactTextString(m) }
func (*DocumentSupportArchive) ProtoMessage()    {}

func (m *DocumentSupportArchive) GetCommandHistory() *TSP.Reference {
	if m != nil {
		return m.CommandHistory
	}
	return nil
}

func (m *DocumentSupportArchive) GetCommandSelectionBehaviorHistory() *TSP.Reference {
	if m != nil {
		return m.CommandSelectionBehaviorHistory
	}
	return nil
}

func (m *DocumentSupportArchive) GetUndoCount() uint32 {
	if m != nil && m.UndoCount != nil {
		return *m.UndoCount
	}
	return 0
}

func (m *DocumentSupportArchive) GetRedoCount() uint32 {
	if m != nil && m.RedoCount != nil {
		return *m.RedoCount
	}
	return 0
}

func (m *DocumentSupportArchive) GetUndoActionString() string {
	if m != nil && m.UndoActionString != nil {
		return *m.UndoActionString
	}
	return ""
}

func (m *DocumentSupportArchive) GetRedoActionString() string {
	if m != nil && m.RedoActionString != nil {
		return *m.RedoActionString
	}
	return ""
}

func (m *DocumentSupportArchive) GetWebState() *TSP.Reference {
	if m != nil {
		return m.WebState
	}
	return nil
}

type ViewStateArchive struct {
	ViewStateRoot    *TSP.Reference `protobuf:"bytes,1,req,name=view_state_root" json:"view_state_root,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *ViewStateArchive) Reset()         { *m = ViewStateArchive{} }
func (m *ViewStateArchive) String() string { return proto.CompactTextString(m) }
func (*ViewStateArchive) ProtoMessage()    {}

func (m *ViewStateArchive) GetViewStateRoot() *TSP.Reference {
	if m != nil {
		return m.ViewStateRoot
	}
	return nil
}

type CommandArchive struct {
	UndoRedoState    *TSP.Reference `protobuf:"bytes,1,opt,name=undoRedoState" json:"undoRedoState,omitempty"`
	UndoCollection   *TSP.Reference `protobuf:"bytes,2,opt,name=undoCollection" json:"undoCollection,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *CommandArchive) Reset()         { *m = CommandArchive{} }
func (m *CommandArchive) String() string { return proto.CompactTextString(m) }
func (*CommandArchive) ProtoMessage()    {}

func (m *CommandArchive) GetUndoRedoState() *TSP.Reference {
	if m != nil {
		return m.UndoRedoState
	}
	return nil
}

func (m *CommandArchive) GetUndoCollection() *TSP.Reference {
	if m != nil {
		return m.UndoCollection
	}
	return nil
}

type CommandGroupArchive struct {
	Super            *CommandArchive  `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	Commands         []*TSP.Reference `protobuf:"bytes,2,rep,name=commands" json:"commands,omitempty"`
	ProcessResults   *TSP.IndexSet    `protobuf:"bytes,3,opt,name=process_results" json:"process_results,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *CommandGroupArchive) Reset()         { *m = CommandGroupArchive{} }
func (m *CommandGroupArchive) String() string { return proto.CompactTextString(m) }
func (*CommandGroupArchive) ProtoMessage()    {}

func (m *CommandGroupArchive) GetSuper() *CommandArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *CommandGroupArchive) GetCommands() []*TSP.Reference {
	if m != nil {
		return m.Commands
	}
	return nil
}

func (m *CommandGroupArchive) GetProcessResults() *TSP.IndexSet {
	if m != nil {
		return m.ProcessResults
	}
	return nil
}

type CommandContainerArchive struct {
	Commands         []*TSP.Reference `protobuf:"bytes,1,rep,name=commands" json:"commands,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *CommandContainerArchive) Reset()         { *m = CommandContainerArchive{} }
func (m *CommandContainerArchive) String() string { return proto.CompactTextString(m) }
func (*CommandContainerArchive) ProtoMessage()    {}

func (m *CommandContainerArchive) GetCommands() []*TSP.Reference {
	if m != nil {
		return m.Commands
	}
	return nil
}

type ReplaceAllChildCommandArchive struct {
	Super            *CommandArchive `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *ReplaceAllChildCommandArchive) Reset()         { *m = ReplaceAllChildCommandArchive{} }
func (m *ReplaceAllChildCommandArchive) String() string { return proto.CompactTextString(m) }
func (*ReplaceAllChildCommandArchive) ProtoMessage()    {}

func (m *ReplaceAllChildCommandArchive) GetSuper() *CommandArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

type ReplaceAllCommandArchive struct {
	Super            *CommandArchive  `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	Commands         []*TSP.Reference `protobuf:"bytes,2,rep,name=commands" json:"commands,omitempty"`
	FindString       *string          `protobuf:"bytes,3,req,name=find_string" json:"find_string,omitempty"`
	ReplaceString    *string          `protobuf:"bytes,4,req,name=replace_string" json:"replace_string,omitempty"`
	Options          *uint32          `protobuf:"varint,5,req,name=options" json:"options,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *ReplaceAllCommandArchive) Reset()         { *m = ReplaceAllCommandArchive{} }
func (m *ReplaceAllCommandArchive) String() string { return proto.CompactTextString(m) }
func (*ReplaceAllCommandArchive) ProtoMessage()    {}

func (m *ReplaceAllCommandArchive) GetSuper() *CommandArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *ReplaceAllCommandArchive) GetCommands() []*TSP.Reference {
	if m != nil {
		return m.Commands
	}
	return nil
}

func (m *ReplaceAllCommandArchive) GetFindString() string {
	if m != nil && m.FindString != nil {
		return *m.FindString
	}
	return ""
}

func (m *ReplaceAllCommandArchive) GetReplaceString() string {
	if m != nil && m.ReplaceString != nil {
		return *m.ReplaceString
	}
	return ""
}

func (m *ReplaceAllCommandArchive) GetOptions() uint32 {
	if m != nil && m.Options != nil {
		return *m.Options
	}
	return 0
}

type ShuffleMappingArchive struct {
	StartIndex              *uint32                        `protobuf:"varint,1,req,name=start_index" json:"start_index,omitempty"`
	EndIndex                *uint32                        `protobuf:"varint,2,req,name=end_index" json:"end_index,omitempty"`
	Entries                 []*ShuffleMappingArchive_Entry `protobuf:"bytes,3,rep,name=entries" json:"entries,omitempty"`
	IsVertical              *bool                          `protobuf:"varint,4,opt,name=is_vertical,def=1" json:"is_vertical,omitempty"`
	IsMoveOperation         *bool                          `protobuf:"varint,5,opt,name=is_move_operation,def=0" json:"is_move_operation,omitempty"`
	FirstMovedIndex         *uint32                        `protobuf:"varint,6,opt,name=first_moved_index,def=0" json:"first_moved_index,omitempty"`
	DestinationIndexForMove *uint32                        `protobuf:"varint,7,opt,name=destination_index_for_move,def=0" json:"destination_index_for_move,omitempty"`
	NumberOfIndicesMoved    *uint32                        `protobuf:"varint,8,opt,name=number_of_indices_moved,def=0" json:"number_of_indices_moved,omitempty"`
	XXX_unrecognized        []byte                         `json:"-"`
}

func (m *ShuffleMappingArchive) Reset()         { *m = ShuffleMappingArchive{} }
func (m *ShuffleMappingArchive) String() string { return proto.CompactTextString(m) }
func (*ShuffleMappingArchive) ProtoMessage()    {}

const Default_ShuffleMappingArchive_IsVertical bool = true
const Default_ShuffleMappingArchive_IsMoveOperation bool = false
const Default_ShuffleMappingArchive_FirstMovedIndex uint32 = 0
const Default_ShuffleMappingArchive_DestinationIndexForMove uint32 = 0
const Default_ShuffleMappingArchive_NumberOfIndicesMoved uint32 = 0

func (m *ShuffleMappingArchive) GetStartIndex() uint32 {
	if m != nil && m.StartIndex != nil {
		return *m.StartIndex
	}
	return 0
}

func (m *ShuffleMappingArchive) GetEndIndex() uint32 {
	if m != nil && m.EndIndex != nil {
		return *m.EndIndex
	}
	return 0
}

func (m *ShuffleMappingArchive) GetEntries() []*ShuffleMappingArchive_Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *ShuffleMappingArchive) GetIsVertical() bool {
	if m != nil && m.IsVertical != nil {
		return *m.IsVertical
	}
	return Default_ShuffleMappingArchive_IsVertical
}

func (m *ShuffleMappingArchive) GetIsMoveOperation() bool {
	if m != nil && m.IsMoveOperation != nil {
		return *m.IsMoveOperation
	}
	return Default_ShuffleMappingArchive_IsMoveOperation
}

func (m *ShuffleMappingArchive) GetFirstMovedIndex() uint32 {
	if m != nil && m.FirstMovedIndex != nil {
		return *m.FirstMovedIndex
	}
	return Default_ShuffleMappingArchive_FirstMovedIndex
}

func (m *ShuffleMappingArchive) GetDestinationIndexForMove() uint32 {
	if m != nil && m.DestinationIndexForMove != nil {
		return *m.DestinationIndexForMove
	}
	return Default_ShuffleMappingArchive_DestinationIndexForMove
}

func (m *ShuffleMappingArchive) GetNumberOfIndicesMoved() uint32 {
	if m != nil && m.NumberOfIndicesMoved != nil {
		return *m.NumberOfIndicesMoved
	}
	return Default_ShuffleMappingArchive_NumberOfIndicesMoved
}

type ShuffleMappingArchive_Entry struct {
	From             *uint32 `protobuf:"varint,1,req,name=from" json:"from,omitempty"`
	To               *uint32 `protobuf:"varint,2,req,name=to" json:"to,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ShuffleMappingArchive_Entry) Reset()         { *m = ShuffleMappingArchive_Entry{} }
func (m *ShuffleMappingArchive_Entry) String() string { return proto.CompactTextString(m) }
func (*ShuffleMappingArchive_Entry) ProtoMessage()    {}

func (m *ShuffleMappingArchive_Entry) GetFrom() uint32 {
	if m != nil && m.From != nil {
		return *m.From
	}
	return 0
}

func (m *ShuffleMappingArchive_Entry) GetTo() uint32 {
	if m != nil && m.To != nil {
		return *m.To
	}
	return 0
}

type ProgressiveCommandGroupArchive struct {
	Super            *CommandGroupArchive `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *ProgressiveCommandGroupArchive) Reset()         { *m = ProgressiveCommandGroupArchive{} }
func (m *ProgressiveCommandGroupArchive) String() string { return proto.CompactTextString(m) }
func (*ProgressiveCommandGroupArchive) ProtoMessage()    {}

func (m *ProgressiveCommandGroupArchive) GetSuper() *CommandGroupArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

type CommandSelectionBehaviorHistoryArchive struct {
	Entries          []*CommandSelectionBehaviorHistoryArchive_Entry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
	XXX_unrecognized []byte                                          `json:"-"`
}

func (m *CommandSelectionBehaviorHistoryArchive) Reset() {
	*m = CommandSelectionBehaviorHistoryArchive{}
}
func (m *CommandSelectionBehaviorHistoryArchive) String() string { return proto.CompactTextString(m) }
func (*CommandSelectionBehaviorHistoryArchive) ProtoMessage()    {}

func (m *CommandSelectionBehaviorHistoryArchive) GetEntries() []*CommandSelectionBehaviorHistoryArchive_Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

type CommandSelectionBehaviorHistoryArchive_Entry struct {
	Command                  *TSP.Reference `protobuf:"bytes,1,req,name=command" json:"command,omitempty"`
	CommandSelectionBehavior *TSP.Reference `protobuf:"bytes,2,req,name=command_selection_behavior" json:"command_selection_behavior,omitempty"`
	XXX_unrecognized         []byte         `json:"-"`
}

func (m *CommandSelectionBehaviorHistoryArchive_Entry) Reset() {
	*m = CommandSelectionBehaviorHistoryArchive_Entry{}
}
func (m *CommandSelectionBehaviorHistoryArchive_Entry) String() string {
	return proto.CompactTextString(m)
}
func (*CommandSelectionBehaviorHistoryArchive_Entry) ProtoMessage() {}

func (m *CommandSelectionBehaviorHistoryArchive_Entry) GetCommand() *TSP.Reference {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *CommandSelectionBehaviorHistoryArchive_Entry) GetCommandSelectionBehavior() *TSP.Reference {
	if m != nil {
		return m.CommandSelectionBehavior
	}
	return nil
}

type UndoRedoStateCommandSelectionBehaviorArchive struct {
	UndoRedoState    *TSP.Reference `protobuf:"bytes,2,opt,name=undo_redo_state" json:"undo_redo_state,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *UndoRedoStateCommandSelectionBehaviorArchive) Reset() {
	*m = UndoRedoStateCommandSelectionBehaviorArchive{}
}
func (m *UndoRedoStateCommandSelectionBehaviorArchive) String() string {
	return proto.CompactTextString(m)
}
func (*UndoRedoStateCommandSelectionBehaviorArchive) ProtoMessage() {}

func (m *UndoRedoStateCommandSelectionBehaviorArchive) GetUndoRedoState() *TSP.Reference {
	if m != nil {
		return m.UndoRedoState
	}
	return nil
}

type FormatStructArchive struct {
	FormatType                     *uint32                   `protobuf:"varint,1,req,name=format_type" json:"format_type,omitempty"`
	DecimalPlaces                  *uint32                   `protobuf:"varint,2,opt,name=decimal_places" json:"decimal_places,omitempty"`
	CurrencyCode                   *string                   `protobuf:"bytes,3,opt,name=currency_code" json:"currency_code,omitempty"`
	NegativeStyle                  *uint32                   `protobuf:"varint,4,opt,name=negative_style" json:"negative_style,omitempty"`
	ShowThousandsSeparator         *bool                     `protobuf:"varint,5,opt,name=show_thousands_separator" json:"show_thousands_separator,omitempty"`
	UseAccountingStyle             *bool                     `protobuf:"varint,6,opt,name=use_accounting_style" json:"use_accounting_style,omitempty"`
	DurationStyle                  *uint32                   `protobuf:"varint,7,opt,name=duration_style" json:"duration_style,omitempty"`
	Base                           *uint32                   `protobuf:"varint,8,opt,name=base" json:"base,omitempty"`
	BasePlaces                     *uint32                   `protobuf:"varint,9,opt,name=base_places" json:"base_places,omitempty"`
	BaseUseMinusSign               *bool                     `protobuf:"varint,10,opt,name=base_use_minus_sign" json:"base_use_minus_sign,omitempty"`
	FractionAccuracy               *uint32                   `protobuf:"varint,11,opt,name=fraction_accuracy" json:"fraction_accuracy,omitempty"`
	SuppressDateFormat             *bool                     `protobuf:"varint,12,opt,name=suppress_date_format" json:"suppress_date_format,omitempty"`
	SuppressTimeFormat             *bool                     `protobuf:"varint,13,opt,name=suppress_time_format" json:"suppress_time_format,omitempty"`
	DateTimeFormat                 *string                   `protobuf:"bytes,14,opt,name=date_time_format" json:"date_time_format,omitempty"`
	DurationUnitLargest            *uint32                   `protobuf:"varint,15,opt,name=duration_unit_largest" json:"duration_unit_largest,omitempty"`
	DurationUnitSmallest           *uint32                   `protobuf:"varint,16,opt,name=duration_unit_smallest" json:"duration_unit_smallest,omitempty"`
	CustomId                       *uint32                   `protobuf:"varint,17,opt,name=custom_id" json:"custom_id,omitempty"`
	CustomFormatString             *string                   `protobuf:"bytes,18,opt,name=custom_format_string" json:"custom_format_string,omitempty"`
	ScaleFactor                    *float64                  `protobuf:"fixed64,19,opt,name=scale_factor" json:"scale_factor,omitempty"`
	RequiresFractionReplacement    *bool                     `protobuf:"varint,20,opt,name=requires_fraction_replacement" json:"requires_fraction_replacement,omitempty"`
	ControlMinimum                 *float64                  `protobuf:"fixed64,21,opt,name=control_minimum" json:"control_minimum,omitempty"`
	ControlMaximum                 *float64                  `protobuf:"fixed64,22,opt,name=control_maximum" json:"control_maximum,omitempty"`
	ControlIncrement               *float64                  `protobuf:"fixed64,23,opt,name=control_increment" json:"control_increment,omitempty"`
	ControlFormatType              *uint32                   `protobuf:"varint,24,opt,name=control_format_type" json:"control_format_type,omitempty"`
	SliderOrientation              *uint32                   `protobuf:"varint,25,opt,name=slider_orientation" json:"slider_orientation,omitempty"`
	SliderPosition                 *uint32                   `protobuf:"varint,26,opt,name=slider_position" json:"slider_position,omitempty"`
	DecimalWidth                   *uint32                   `protobuf:"varint,27,opt,name=decimal_width" json:"decimal_width,omitempty"`
	MinIntegerWidth                *uint32                   `protobuf:"varint,28,opt,name=min_integer_width" json:"min_integer_width,omitempty"`
	NumNonspaceIntegerDigits       *uint32                   `protobuf:"varint,29,opt,name=num_nonspace_integer_digits" json:"num_nonspace_integer_digits,omitempty"`
	NumNonspaceDecimalDigits       *uint32                   `protobuf:"varint,30,opt,name=num_nonspace_decimal_digits" json:"num_nonspace_decimal_digits,omitempty"`
	IndexFromRightLastInteger      *uint32                   `protobuf:"varint,31,opt,name=index_from_right_last_integer" json:"index_from_right_last_integer,omitempty"`
	InterstitialStrings            []string                  `protobuf:"bytes,32,rep,name=interstitial_strings" json:"interstitial_strings,omitempty"`
	IntersStrInsertionIndexes      *TSP.IndexSet             `protobuf:"bytes,33,opt,name=inters_str_insertion_indexes" json:"inters_str_insertion_indexes,omitempty"`
	NumHashDecimalDigits           *uint32                   `protobuf:"varint,34,opt,name=num_hash_decimal_digits" json:"num_hash_decimal_digits,omitempty"`
	TotalNumDecimalDigits          *uint32                   `protobuf:"varint,35,opt,name=total_num_decimal_digits" json:"total_num_decimal_digits,omitempty"`
	IsComplex                      *bool                     `protobuf:"varint,36,opt,name=is_complex" json:"is_complex,omitempty"`
	ContainsIntegerToken           *bool                     `protobuf:"varint,37,opt,name=contains_integer_token" json:"contains_integer_token,omitempty"`
	MultipleChoiceListInitialValue *uint32                   `protobuf:"varint,38,opt,name=multiple_choice_list_initial_value" json:"multiple_choice_list_initial_value,omitempty"`
	MultipleChoiceListId           *uint32                   `protobuf:"varint,39,opt,name=multiple_choice_list_id" json:"multiple_choice_list_id,omitempty"`
	UseAutomaticDurationUnits      *bool                     `protobuf:"varint,40,opt,name=use_automatic_duration_units" json:"use_automatic_duration_units,omitempty"`
	XXX_extensions                 map[int32]proto.Extension `json:"-"`
	XXX_unrecognized               []byte                    `json:"-"`
}

func (m *FormatStructArchive) Reset()         { *m = FormatStructArchive{} }
func (m *FormatStructArchive) String() string { return proto.CompactTextString(m) }
func (*FormatStructArchive) ProtoMessage()    {}

var extRange_FormatStructArchive = []proto.ExtensionRange{
	{10000, 19999},
}

func (*FormatStructArchive) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_FormatStructArchive
}
func (m *FormatStructArchive) ExtensionMap() map[int32]proto.Extension {
	if m.XXX_extensions == nil {
		m.XXX_extensions = make(map[int32]proto.Extension)
	}
	return m.XXX_extensions
}

func (m *FormatStructArchive) GetFormatType() uint32 {
	if m != nil && m.FormatType != nil {
		return *m.FormatType
	}
	return 0
}

func (m *FormatStructArchive) GetDecimalPlaces() uint32 {
	if m != nil && m.DecimalPlaces != nil {
		return *m.DecimalPlaces
	}
	return 0
}

func (m *FormatStructArchive) GetCurrencyCode() string {
	if m != nil && m.CurrencyCode != nil {
		return *m.CurrencyCode
	}
	return ""
}

func (m *FormatStructArchive) GetNegativeStyle() uint32 {
	if m != nil && m.NegativeStyle != nil {
		return *m.NegativeStyle
	}
	return 0
}

func (m *FormatStructArchive) GetShowThousandsSeparator() bool {
	if m != nil && m.ShowThousandsSeparator != nil {
		return *m.ShowThousandsSeparator
	}
	return false
}

func (m *FormatStructArchive) GetUseAccountingStyle() bool {
	if m != nil && m.UseAccountingStyle != nil {
		return *m.UseAccountingStyle
	}
	return false
}

func (m *FormatStructArchive) GetDurationStyle() uint32 {
	if m != nil && m.DurationStyle != nil {
		return *m.DurationStyle
	}
	return 0
}

func (m *FormatStructArchive) GetBase() uint32 {
	if m != nil && m.Base != nil {
		return *m.Base
	}
	return 0
}

func (m *FormatStructArchive) GetBasePlaces() uint32 {
	if m != nil && m.BasePlaces != nil {
		return *m.BasePlaces
	}
	return 0
}

func (m *FormatStructArchive) GetBaseUseMinusSign() bool {
	if m != nil && m.BaseUseMinusSign != nil {
		return *m.BaseUseMinusSign
	}
	return false
}

func (m *FormatStructArchive) GetFractionAccuracy() uint32 {
	if m != nil && m.FractionAccuracy != nil {
		return *m.FractionAccuracy
	}
	return 0
}

func (m *FormatStructArchive) GetSuppressDateFormat() bool {
	if m != nil && m.SuppressDateFormat != nil {
		return *m.SuppressDateFormat
	}
	return false
}

func (m *FormatStructArchive) GetSuppressTimeFormat() bool {
	if m != nil && m.SuppressTimeFormat != nil {
		return *m.SuppressTimeFormat
	}
	return false
}

func (m *FormatStructArchive) GetDateTimeFormat() string {
	if m != nil && m.DateTimeFormat != nil {
		return *m.DateTimeFormat
	}
	return ""
}

func (m *FormatStructArchive) GetDurationUnitLargest() uint32 {
	if m != nil && m.DurationUnitLargest != nil {
		return *m.DurationUnitLargest
	}
	return 0
}

func (m *FormatStructArchive) GetDurationUnitSmallest() uint32 {
	if m != nil && m.DurationUnitSmallest != nil {
		return *m.DurationUnitSmallest
	}
	return 0
}

func (m *FormatStructArchive) GetCustomId() uint32 {
	if m != nil && m.CustomId != nil {
		return *m.CustomId
	}
	return 0
}

func (m *FormatStructArchive) GetCustomFormatString() string {
	if m != nil && m.CustomFormatString != nil {
		return *m.CustomFormatString
	}
	return ""
}

func (m *FormatStructArchive) GetScaleFactor() float64 {
	if m != nil && m.ScaleFactor != nil {
		return *m.ScaleFactor
	}
	return 0
}

func (m *FormatStructArchive) GetRequiresFractionReplacement() bool {
	if m != nil && m.RequiresFractionReplacement != nil {
		return *m.RequiresFractionReplacement
	}
	return false
}

func (m *FormatStructArchive) GetControlMinimum() float64 {
	if m != nil && m.ControlMinimum != nil {
		return *m.ControlMinimum
	}
	return 0
}

func (m *FormatStructArchive) GetControlMaximum() float64 {
	if m != nil && m.ControlMaximum != nil {
		return *m.ControlMaximum
	}
	return 0
}

func (m *FormatStructArchive) GetControlIncrement() float64 {
	if m != nil && m.ControlIncrement != nil {
		return *m.ControlIncrement
	}
	return 0
}

func (m *FormatStructArchive) GetControlFormatType() uint32 {
	if m != nil && m.ControlFormatType != nil {
		return *m.ControlFormatType
	}
	return 0
}

func (m *FormatStructArchive) GetSliderOrientation() uint32 {
	if m != nil && m.SliderOrientation != nil {
		return *m.SliderOrientation
	}
	return 0
}

func (m *FormatStructArchive) GetSliderPosition() uint32 {
	if m != nil && m.SliderPosition != nil {
		return *m.SliderPosition
	}
	return 0
}

func (m *FormatStructArchive) GetDecimalWidth() uint32 {
	if m != nil && m.DecimalWidth != nil {
		return *m.DecimalWidth
	}
	return 0
}

func (m *FormatStructArchive) GetMinIntegerWidth() uint32 {
	if m != nil && m.MinIntegerWidth != nil {
		return *m.MinIntegerWidth
	}
	return 0
}

func (m *FormatStructArchive) GetNumNonspaceIntegerDigits() uint32 {
	if m != nil && m.NumNonspaceIntegerDigits != nil {
		return *m.NumNonspaceIntegerDigits
	}
	return 0
}

func (m *FormatStructArchive) GetNumNonspaceDecimalDigits() uint32 {
	if m != nil && m.NumNonspaceDecimalDigits != nil {
		return *m.NumNonspaceDecimalDigits
	}
	return 0
}

func (m *FormatStructArchive) GetIndexFromRightLastInteger() uint32 {
	if m != nil && m.IndexFromRightLastInteger != nil {
		return *m.IndexFromRightLastInteger
	}
	return 0
}

func (m *FormatStructArchive) GetInterstitialStrings() []string {
	if m != nil {
		return m.InterstitialStrings
	}
	return nil
}

func (m *FormatStructArchive) GetIntersStrInsertionIndexes() *TSP.IndexSet {
	if m != nil {
		return m.IntersStrInsertionIndexes
	}
	return nil
}

func (m *FormatStructArchive) GetNumHashDecimalDigits() uint32 {
	if m != nil && m.NumHashDecimalDigits != nil {
		return *m.NumHashDecimalDigits
	}
	return 0
}

func (m *FormatStructArchive) GetTotalNumDecimalDigits() uint32 {
	if m != nil && m.TotalNumDecimalDigits != nil {
		return *m.TotalNumDecimalDigits
	}
	return 0
}

func (m *FormatStructArchive) GetIsComplex() bool {
	if m != nil && m.IsComplex != nil {
		return *m.IsComplex
	}
	return false
}

func (m *FormatStructArchive) GetContainsIntegerToken() bool {
	if m != nil && m.ContainsIntegerToken != nil {
		return *m.ContainsIntegerToken
	}
	return false
}

func (m *FormatStructArchive) GetMultipleChoiceListInitialValue() uint32 {
	if m != nil && m.MultipleChoiceListInitialValue != nil {
		return *m.MultipleChoiceListInitialValue
	}
	return 0
}

func (m *FormatStructArchive) GetMultipleChoiceListId() uint32 {
	if m != nil && m.MultipleChoiceListId != nil {
		return *m.MultipleChoiceListId
	}
	return 0
}

func (m *FormatStructArchive) GetUseAutomaticDurationUnits() bool {
	if m != nil && m.UseAutomaticDurationUnits != nil {
		return *m.UseAutomaticDurationUnits
	}
	return false
}

type CustomFormatArchive struct {
	Name             *string                          `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	FormatType       *uint32                          `protobuf:"varint,2,req,name=format_type" json:"format_type,omitempty"`
	DefaultFormat    *FormatStructArchive             `protobuf:"bytes,3,req,name=default_format" json:"default_format,omitempty"`
	Conditions       []*CustomFormatArchive_Condition `protobuf:"bytes,4,rep,name=conditions" json:"conditions,omitempty"`
	XXX_unrecognized []byte                           `json:"-"`
}

func (m *CustomFormatArchive) Reset()         { *m = CustomFormatArchive{} }
func (m *CustomFormatArchive) String() string { return proto.CompactTextString(m) }
func (*CustomFormatArchive) ProtoMessage()    {}

func (m *CustomFormatArchive) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *CustomFormatArchive) GetFormatType() uint32 {
	if m != nil && m.FormatType != nil {
		return *m.FormatType
	}
	return 0
}

func (m *CustomFormatArchive) GetDefaultFormat() *FormatStructArchive {
	if m != nil {
		return m.DefaultFormat
	}
	return nil
}

func (m *CustomFormatArchive) GetConditions() []*CustomFormatArchive_Condition {
	if m != nil {
		return m.Conditions
	}
	return nil
}

type CustomFormatArchive_Condition struct {
	ConditionType     *uint32              `protobuf:"varint,1,req,name=condition_type" json:"condition_type,omitempty"`
	ConditionValue    *float32             `protobuf:"fixed32,2,opt,name=condition_value" json:"condition_value,omitempty"`
	ConditionFormat   *FormatStructArchive `protobuf:"bytes,3,req,name=condition_format" json:"condition_format,omitempty"`
	ConditionValueDbl *float64             `protobuf:"fixed64,4,opt,name=condition_value_dbl" json:"condition_value_dbl,omitempty"`
	XXX_unrecognized  []byte               `json:"-"`
}

func (m *CustomFormatArchive_Condition) Reset()         { *m = CustomFormatArchive_Condition{} }
func (m *CustomFormatArchive_Condition) String() string { return proto.CompactTextString(m) }
func (*CustomFormatArchive_Condition) ProtoMessage()    {}

func (m *CustomFormatArchive_Condition) GetConditionType() uint32 {
	if m != nil && m.ConditionType != nil {
		return *m.ConditionType
	}
	return 0
}

func (m *CustomFormatArchive_Condition) GetConditionValue() float32 {
	if m != nil && m.ConditionValue != nil {
		return *m.ConditionValue
	}
	return 0
}

func (m *CustomFormatArchive_Condition) GetConditionFormat() *FormatStructArchive {
	if m != nil {
		return m.ConditionFormat
	}
	return nil
}

func (m *CustomFormatArchive_Condition) GetConditionValueDbl() float64 {
	if m != nil && m.ConditionValueDbl != nil {
		return *m.ConditionValueDbl
	}
	return 0
}

type AnnotationAuthorArchive struct {
	Name             *string    `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Color            *TSP.Color `protobuf:"bytes,2,opt,name=color" json:"color,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *AnnotationAuthorArchive) Reset()         { *m = AnnotationAuthorArchive{} }
func (m *AnnotationAuthorArchive) String() string { return proto.CompactTextString(m) }
func (*AnnotationAuthorArchive) ProtoMessage()    {}

func (m *AnnotationAuthorArchive) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *AnnotationAuthorArchive) GetColor() *TSP.Color {
	if m != nil {
		return m.Color
	}
	return nil
}

type DeprecatedChangeAuthorArchive struct {
	Name             *string    `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	ChangeColor      *TSP.Color `protobuf:"bytes,2,opt,name=change_color" json:"change_color,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *DeprecatedChangeAuthorArchive) Reset()         { *m = DeprecatedChangeAuthorArchive{} }
func (m *DeprecatedChangeAuthorArchive) String() string { return proto.CompactTextString(m) }
func (*DeprecatedChangeAuthorArchive) ProtoMessage()    {}

func (m *DeprecatedChangeAuthorArchive) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *DeprecatedChangeAuthorArchive) GetChangeColor() *TSP.Color {
	if m != nil {
		return m.ChangeColor
	}
	return nil
}

type AnnotationAuthorStorageArchive struct {
	AnnotationAuthor []*TSP.Reference `protobuf:"bytes,1,rep,name=annotation_author" json:"annotation_author,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *AnnotationAuthorStorageArchive) Reset()         { *m = AnnotationAuthorStorageArchive{} }
func (m *AnnotationAuthorStorageArchive) String() string { return proto.CompactTextString(m) }
func (*AnnotationAuthorStorageArchive) ProtoMessage()    {}

func (m *AnnotationAuthorStorageArchive) GetAnnotationAuthor() []*TSP.Reference {
	if m != nil {
		return m.AnnotationAuthor
	}
	return nil
}

type AddAnnotationAuthorCommandArchive struct {
	Super            *CommandArchive `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	DocumentRoot     *TSP.Reference  `protobuf:"bytes,2,opt,name=document_root" json:"document_root,omitempty"`
	AnnotationAuthor *TSP.Reference  `protobuf:"bytes,3,opt,name=annotation_author" json:"annotation_author,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *AddAnnotationAuthorCommandArchive) Reset()         { *m = AddAnnotationAuthorCommandArchive{} }
func (m *AddAnnotationAuthorCommandArchive) String() string { return proto.CompactTextString(m) }
func (*AddAnnotationAuthorCommandArchive) ProtoMessage()    {}

func (m *AddAnnotationAuthorCommandArchive) GetSuper() *CommandArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *AddAnnotationAuthorCommandArchive) GetDocumentRoot() *TSP.Reference {
	if m != nil {
		return m.DocumentRoot
	}
	return nil
}

func (m *AddAnnotationAuthorCommandArchive) GetAnnotationAuthor() *TSP.Reference {
	if m != nil {
		return m.AnnotationAuthor
	}
	return nil
}

type SetAnnotationAuthorColorCommandArchive struct {
	Super            *CommandArchive `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	AnnotationAuthor *TSP.Reference  `protobuf:"bytes,2,opt,name=annotation_author" json:"annotation_author,omitempty"`
	Color            *TSP.Color      `protobuf:"bytes,3,opt,name=color" json:"color,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *SetAnnotationAuthorColorCommandArchive) Reset() {
	*m = SetAnnotationAuthorColorCommandArchive{}
}
func (m *SetAnnotationAuthorColorCommandArchive) String() string { return proto.CompactTextString(m) }
func (*SetAnnotationAuthorColorCommandArchive) ProtoMessage()    {}

func (m *SetAnnotationAuthorColorCommandArchive) GetSuper() *CommandArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *SetAnnotationAuthorColorCommandArchive) GetAnnotationAuthor() *TSP.Reference {
	if m != nil {
		return m.AnnotationAuthor
	}
	return nil
}

func (m *SetAnnotationAuthorColorCommandArchive) GetColor() *TSP.Color {
	if m != nil {
		return m.Color
	}
	return nil
}

func init() {
}
