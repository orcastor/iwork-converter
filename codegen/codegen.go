package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

var foo = `
package index
import (
    {{if eq .FunctionName "Common"}}"errors"
    "fmt"
    {{end}}"github.com/golang/protobuf/proto"
    "github.com/orcastor/iwork-converter/proto/TN"
    "github.com/orcastor/iwork-converter/proto/TST"
    "github.com/orcastor/iwork-converter/proto/TSWP"
)

func decode{{.FunctionName}}(typ uint32, payload []byte) (interface{}, error) {
    switch typ {
        {{range $key, $value := .}}
        {{if ne $key "FunctionName"}}
        case {{$key}}:
            var value = &{{$value}}{}
            err := proto.Unmarshal(payload, value)
            return value, err
        {{end}}
        {{end}}
        default:
            {{if ne .FunctionName "Common"}}
            // 兜底逻辑：尝试使用decodeCommon
            return decodeCommon(typ, payload)
            {{else}}
            return nil,errors.New(fmt.Sprintf("Unknown type %d", typ))
            {{end}}
    }
}
`

func main() {
	if len(os.Args) < 3 {
		panic("Usage: codegen <json_file> <function_name>")
	}

	data, err := ioutil.ReadFile(os.Args[1])
	must(err)
	var info map[string]string

	must(json.Unmarshal(data, &info))

	// Add function name to the data
	info["FunctionName"] = os.Args[2]

	tmpl, err := template.New("test").Funcs(template.FuncMap{
		"hasPrefix": strings.HasPrefix,
	}).Parse(foo)
	must(err)

	must(tmpl.Execute(os.Stdout, info))
}
