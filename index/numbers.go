
package index
import (
    "github.com/golang/protobuf/proto"
    "github.com/orcastor/iwork-converter/proto/TN"
    "github.com/orcastor/iwork-converter/proto/TST"
    "github.com/orcastor/iwork-converter/proto/TSWP"
)

func decodeNumbers(typ uint32, payload []byte) (interface{}, error) {
    switch typ {
        
        
        case 1:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 10011:
            var value = &TSWP.SectionPlaceholderArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12002:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12003:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12004:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12005:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12006:
            var value = &TN.ChartMediatorArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12007:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12008:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12009:
            var value = &TN.ThemeArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12010:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12011:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12012:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12013:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12014:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12015:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12016:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12017:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12018:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12019:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12021:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12024:
            var value = &TN.UndoRedoStateArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12025:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12026:
            var value = &TN.UIStateArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12027:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12028:
            var value = &TN.SheetSelectionArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12029:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 12030:
            var value = &TN.DocumentArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 2:
            var value = &TN.SheetArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 3:
            var value = &TN.FormBasedSheetArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 3061:
            var value = &TST.TableModelArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        case 6030:
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
        
        
        
        case 7:
            var value = &TN.PlaceholderArchive{}
            err := proto.Unmarshal(payload, value)
            return value, err
        
        
        
        
        default:
            
            // 兜底逻辑：尝试使用decodeCommon
            return decodeCommon(typ, payload)
            
    }
}
