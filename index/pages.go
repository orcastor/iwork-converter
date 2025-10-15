package index

import (
	"github.com/golang/protobuf/proto"
	"github.com/orcastor/iwork-converter/proto/TP"
)

func decodePages(typ uint32, payload []byte) (interface{}, error) {
	switch typ {

	case 10000:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10001:
		var value = &TP.ThemeArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10010:
		var value = &TP.FloatingDrawablesArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10011:
		var value = &TP.SectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10012:
		var value = &TP.SettingsArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10015:
		var value = &TP.DrawablesZOrderArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10101:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10102:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10108:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10109:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10110:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10111:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10112:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10113:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10114:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10115:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10116:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10117:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10118:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10119:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10120:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10121:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10125:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10126:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10127:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10128:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10130:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10131:
		var value = &TP.LayoutStateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10132:
		var value = &TP.CanvasSelectionArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10133:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10134:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10140:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10141:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10142:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10143:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10147:
		var value = &TP.UIStateArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10148:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10149:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10150:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10151:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10152:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10153:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10154:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10155:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10156:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 10157:
		var value = &TP.DocumentArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	case 7:
		var value = &TP.PlaceholderArchive{}
		err := proto.Unmarshal(payload, value)
		return value, err

	default:

		// 兜底逻辑：尝试使用decodeCommon
		return decodeCommon(typ, payload)

	}
}
