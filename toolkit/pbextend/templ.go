package pbextend

const templ = `
/*
* 本代码由pbtool工具生成，请勿手动修改
*/

package {{.GetPkgName}}

import (
	"github.com/golang/protobuf/proto"
)

{{range $st := .GetAllEnum -}}
func (d {{$st.Name}}) Integer() uint32 {
	return uint32(d.Number())
}

{{end}}

{{range $st := .GetAllRsp -}}
{{range $field := $st.Members -}}
{{if eq $field.Type.Name "*RspHead" -}}
func (d *{{$st.Name}}) SetRspHead(code int32, msg string) {
	d.{{$field.Name}} = &{{index $field.Type.Elements 0}}{Code:code, Msg:msg}
}

func (d *{{$st.Name}}) GetRspHead() (int32, string) {
	return d.{{$field.Name}}.Code, d.{{$field.Name}}.Msg
}
{{end}}
{{end}}
{{end}}

{{range $st := .GetAllStruct -}}
func(d *{{$st.Name}}) ToDB() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return proto.Marshal(d)
}

func(d *{{$st.Name}}) FromDB(val []byte) error {
	if len(val) <= 0 {
		return nil
	}
	return proto.Unmarshal(val, d)
}

{{end}}
`
