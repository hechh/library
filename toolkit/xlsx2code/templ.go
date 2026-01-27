package xlsx2code

const Templ = `
/*
* 本代码由cfgtool工具生成，请勿手动修改
*/

{{$type := .Name}}

package {{ToSnake $type}}

import (
	pb "%s"	
	"sync/atomic"
{{if .List}}
	"github.com/hechh/library/util"
{{end}}
	"github.com/golang/protobuf/proto"
	"github.com/hechh/library/fwatcher"
	"google.golang.org/protobuf/encoding/prototext"
)

const (
	SHEET_NAME = "{{$type}}"
)

var obj = atomic.Pointer[{{$type}}Data]{}


type {{$type}}Data struct {
	list []*pb.{{$type}}
{{- range $index := .List}}
	{{ToLowerCamel $index.Name}} {{$index.Type}}[{{$index.GetType}}, *pb.{{$type}}]
{{- end}}
}

func init() {
	fwatcher.Register(SHEET_NAME, parse)
}

func Change(f func()) {
	fwatcher.Listen(SHEET_NAME, f)
}

func DeepCopy(item *pb.{{$type}}) *pb.{{$type}} {
	buf, _ := proto.Marshal(item)
	ret := &pb.{{$type}}{}
	proto.Unmarshal(buf, ret)
	return ret
}

func parse(buf []byte) error {
	ary := &pb.{{$type}}Ary{}
	if err := prototext.Unmarshal(buf, ary); err != nil {
		return err	
	}

	data := &{{$type}}Data{
{{- range $index := .List}}
		{{ToLowerCamel $index.Name}}: make({{$index.Type}}[{{$index.GetType}}, *pb.{{$type}}]),
{{- end}}
	}
	for _, item := range ary.Ary {
		data.list = append(data.list, item)
	{{- range $index := .List}}
		data.{{ToLowerCamel $index.Name}}.Put({{$index.GetValue "item"}}, item)
	{{- end}}
	}
	obj.Store(data)
	return nil
}

func SGet(pos int) *pb.{{$type}} {
	if pos < 0 {
		pos = 0	
	}
	list := obj.Load().list
	if ll := len(list); ll-1 < pos {
		pos = ll-1	
	}
	return list[pos]
}

func LGet() (rets []*pb.{{$type}}) {
	list := obj.Load().list
	rets = make([]*pb.{{$type}}, len(list))
	copy(rets, list)
	return
}

func Walk(f func(*pb.{{$type}})bool) {
	for _, item := range obj.Load().list {
		if !f(item) {
			return
		}	
	}
}

{{range $index := .List}}
{{if eq $index.Kind 3}}		{{/* map类型 */}}
func MGet{{$index.Name}}({{$index.GetArg}}) *pb.{{$type}} {
	data := obj.Load().{{ToLowerCamel $index.Name}}
	if value, ok := data.Get({{$index.GetValue ""}}); ok {
		return value
	}
	return nil
}

{{else if eq $index.Kind 4}}	{{/* group类型 */}}
func GGet{{$index.Name}}({{$index.GetArg}}) []*pb.{{$type}} {
	data := obj.Load().{{ToLowerCamel $index.Name}}
	if value, ok := data.Get({{$index.GetValue ""}}); ok {
		return value
	}
	return nil
}

func GWalk{{$index.Name}}({{$index.GetArg}}, f func(*pb.{{$type}})bool) {
	data := obj.Load().{{ToLowerCamel $index.Name}}
	if values, ok := data.Get({{$index.GetValue ""}}); ok {
		for _, item := range values {
			if !f(item) {
				return	
			}	
		}
	}
}

func GRange{{$index.Name}}(f func([]*pb.{{$type}})bool) {
	data := obj.Load().{{ToLowerCamel $index.Name}}
	for _, values := range data {
		if !f(values) {
			return	
		}	
	}	
}
{{end}}
{{end}}

`
