// Code generated by protoc-gen-go.
// source: TSAArchives.proto
// DO NOT EDIT!

/*
Package TSA is a generated protocol buffer package.

It is generated from these files:
	TSAArchives.proto

It has these top-level messages:
	DocumentArchive
	FunctionBrowserStateArchive
	TestDocumentArchive
	PropagatePresetCommandArchive
*/
package TSA

import proto "github.com/golang/protobuf/proto"
import math "math"
import "github.com/orcastor/iwork-converter/proto/TSK"
import "github.com/orcastor/iwork-converter/proto/TSP"
import "github.com/orcastor/iwork-converter/proto/TSWP"

// discarding unused import TSS "TSSArchives.pb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type DocumentArchive struct {
	Super                          *TSK.DocumentArchive                 `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	TextPresetDisplayItems         []*TSWP.TextPresetDisplayItemArchive `protobuf:"bytes,2,rep,name=text_preset_display_items" json:"text_preset_display_items,omitempty"`
	CreationLanguage               *string                              `protobuf:"bytes,3,opt,name=creation_language" json:"creation_language,omitempty"`
	CalculationEngine              *TSP.Reference                       `protobuf:"bytes,4,opt,name=calculation_engine" json:"calculation_engine,omitempty"`
	ViewState                      *TSP.Reference                       `protobuf:"bytes,5,opt,name=view_state" json:"view_state,omitempty"`
	FunctionBrowserState           *TSP.Reference                       `protobuf:"bytes,6,opt,name=function_browser_state" json:"function_browser_state,omitempty"`
	TablesCustomFormatList         *TSP.Reference                       `protobuf:"bytes,7,opt,name=tables_custom_format_list" json:"tables_custom_format_list,omitempty"`
	NeedsMovieCompatibilityUpgrade *bool                                `protobuf:"varint,8,opt,name=needs_movie_compatibility_upgrade" json:"needs_movie_compatibility_upgrade,omitempty"`
	TemplateIdentifier             *string                              `protobuf:"bytes,9,opt,name=template_identifier" json:"template_identifier,omitempty"`
	XXX_unrecognized               []byte                               `json:"-"`
}

func (m *DocumentArchive) Reset()         { *m = DocumentArchive{} }
func (m *DocumentArchive) String() string { return proto.CompactTextString(m) }
func (*DocumentArchive) ProtoMessage()    {}

func (m *DocumentArchive) GetSuper() *TSK.DocumentArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *DocumentArchive) GetTextPresetDisplayItems() []*TSWP.TextPresetDisplayItemArchive {
	if m != nil {
		return m.TextPresetDisplayItems
	}
	return nil
}

func (m *DocumentArchive) GetCreationLanguage() string {
	if m != nil && m.CreationLanguage != nil {
		return *m.CreationLanguage
	}
	return ""
}

func (m *DocumentArchive) GetCalculationEngine() *TSP.Reference {
	if m != nil {
		return m.CalculationEngine
	}
	return nil
}

func (m *DocumentArchive) GetViewState() *TSP.Reference {
	if m != nil {
		return m.ViewState
	}
	return nil
}

func (m *DocumentArchive) GetFunctionBrowserState() *TSP.Reference {
	if m != nil {
		return m.FunctionBrowserState
	}
	return nil
}

func (m *DocumentArchive) GetTablesCustomFormatList() *TSP.Reference {
	if m != nil {
		return m.TablesCustomFormatList
	}
	return nil
}

func (m *DocumentArchive) GetNeedsMovieCompatibilityUpgrade() bool {
	if m != nil && m.NeedsMovieCompatibilityUpgrade != nil {
		return *m.NeedsMovieCompatibilityUpgrade
	}
	return false
}

func (m *DocumentArchive) GetTemplateIdentifier() string {
	if m != nil && m.TemplateIdentifier != nil {
		return *m.TemplateIdentifier
	}
	return ""
}

type FunctionBrowserStateArchive struct {
	RecentFunctions  []uint32 `protobuf:"varint,1,rep,name=recent_functions" json:"recent_functions,omitempty"`
	BackFunctions    []uint32 `protobuf:"varint,2,rep,name=back_functions" json:"back_functions,omitempty"`
	ForwardFunctions []uint32 `protobuf:"varint,3,rep,name=forward_functions" json:"forward_functions,omitempty"`
	CurrentFunction  *uint32  `protobuf:"varint,4,opt,name=current_function" json:"current_function,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *FunctionBrowserStateArchive) Reset()         { *m = FunctionBrowserStateArchive{} }
func (m *FunctionBrowserStateArchive) String() string { return proto.CompactTextString(m) }
func (*FunctionBrowserStateArchive) ProtoMessage()    {}

func (m *FunctionBrowserStateArchive) GetRecentFunctions() []uint32 {
	if m != nil {
		return m.RecentFunctions
	}
	return nil
}

func (m *FunctionBrowserStateArchive) GetBackFunctions() []uint32 {
	if m != nil {
		return m.BackFunctions
	}
	return nil
}

func (m *FunctionBrowserStateArchive) GetForwardFunctions() []uint32 {
	if m != nil {
		return m.ForwardFunctions
	}
	return nil
}

func (m *FunctionBrowserStateArchive) GetCurrentFunction() uint32 {
	if m != nil && m.CurrentFunction != nil {
		return *m.CurrentFunction
	}
	return 0
}

type TestDocumentArchive struct {
	Super            *DocumentArchive `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	Value            *string          `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *TestDocumentArchive) Reset()         { *m = TestDocumentArchive{} }
func (m *TestDocumentArchive) String() string { return proto.CompactTextString(m) }
func (*TestDocumentArchive) ProtoMessage()    {}

func (m *TestDocumentArchive) GetSuper() *DocumentArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func (m *TestDocumentArchive) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

type PropagatePresetCommandArchive struct {
	Super            *TSK.CommandArchive `protobuf:"bytes,1,req,name=super" json:"super,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *PropagatePresetCommandArchive) Reset()         { *m = PropagatePresetCommandArchive{} }
func (m *PropagatePresetCommandArchive) String() string { return proto.CompactTextString(m) }
func (*PropagatePresetCommandArchive) ProtoMessage()    {}

func (m *PropagatePresetCommandArchive) GetSuper() *TSK.CommandArchive {
	if m != nil {
		return m.Super
	}
	return nil
}

func init() {
}
