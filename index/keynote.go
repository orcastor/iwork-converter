package index

import (
	"github.com/golang/protobuf/proto"
	"github.com/orcastor/iwork-converter/proto/KN"
	"github.com/orcastor/iwork-converter/proto/TSWP"
)

func decodeKeynote(typ uint32, payload []byte) (interface{}, error) {
	switch typ {
	case 1:
		value := &KN.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10:
		value := &KN.ThemeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 100:
		value := &KN.CommandBuildSetValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10011:
		value := &TSWP.SectionPlaceholderArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 101:
		value := &KN.CommandShowInsertSlideArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 102:
		value := &KN.CommandShowMoveSlideArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 103:
		value := &KN.CommandShowRemoveSlideArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 104:
		value := &KN.CommandSlideInsertDrawablesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 105:
		value := &KN.CommandSlideRemoveDrawableArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 106:
		value := &KN.CommandSlideNodeSetPropertyArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 107:
		value := &KN.CommandSlideInsertBuildArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 108:
		value := &KN.CommandSlideMoveBuildWithoutMovingChunksArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 109:
		value := &KN.CommandSlideRemoveBuildArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 11:
		value := &KN.PasteboardNativeStorageArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 110:
		value := &KN.CommandSlideInsertBuildChunkArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 111:
		value := &KN.CommandSlideMoveBuildChunkArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 112:
		value := &KN.CommandSlideRemoveBuildChunkArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 113:
		value := &KN.CommandSlideSetValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 114:
		value := &KN.CommandTransitionSetValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 115:
		value := &KN.UIStateCommandGroupArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 116:
		value := &KN.CommandSlidePasteDrawablesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 117:
		value := &KN.CommandSlideApplyThemeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 118:
		value := &KN.CommandSlideMoveDrawableZOrderArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 119:
		value := &KN.CommandChangeMasterSlideArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 12:
		value := &KN.PlaceholderArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 123:
		value := &KN.CommandShowSetSlideNumberVisibilityArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 124:
		value := &KN.CommandShowSetValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 128:
		value := &KN.CommandShowMarkOutOfSyncRecordingArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 129:
		value := &KN.CommandShowRemoveRecordingArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 130:
		value := &KN.CommandShowReplaceRecordingArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 131:
		value := &KN.CommandShowSetSoundtrack{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 132:
		value := &KN.CommandSoundtrackSetValue{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 133:
		value := &KN.CommandMasterRescaleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 134:
		value := &KN.CommandMoveMastersArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 135:
		value := &KN.CommandInsertMasterArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 136:
		value := &KN.CommandSlideSetStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 137:
		value := &KN.CommandSlideSetPlaceholdersForTagsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 138:
		value := &KN.CommandBuildChunkSetValueArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 139:
		value := &KN.CommandSlideMoveBuildChunksArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 14:
		value := &TSWP.TextualAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 140:
		value := &KN.CommandRemoveMasterArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 141:
		value := &KN.CommandRenameMasterArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 142:
		value := &KN.CommandMasterSetThumbnailTextArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 143:
		value := &KN.CommandShowChangeThemeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 144:
		value := &KN.CommandSlidePrimitiveSetMasterArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 145:
		value := &KN.CommandMasterSetBodyStylesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 146:
		value := &KN.CommandSlideReapplyMasterArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 147:
		value := &KN.SlideCollectionCommandSelectionBehaviorArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 148:
		value := &KN.ChartInfoGeometryCommandArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 15:
		value := &KN.NoteArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 16:
		value := &KN.RecordingArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 17:
		value := &KN.RecordingEventTrackArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 18:
		value := &KN.RecordingMovieTrackArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 19:
		value := &KN.ClassicStylesheetRecordArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 2:
		value := &KN.ShowArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 20:
		value := &KN.ClassicThemeRecordArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 21:
		value := &KN.Soundtrack{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 22:
		value := &KN.SlideNumberAttachmentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 23:
		value := &KN.DesktopUILayoutArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 24:
		value := &KN.CanvasSelectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 25:
		value := &KN.SlideCollectionSelectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 3:
		value := &KN.UIStateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 4:
		value := &KN.SlideNodeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 5:
		value := &KN.SlideArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 6:
		value := &KN.SlideArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 7:
		value := &KN.PlaceholderArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 8:
		value := &KN.BuildArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 9:
		value := &KN.SlideStyleArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	default:
		return decodeCommon(typ, payload)
	}
}
