package xlsx2proto

const headTempl = `
syntax = "proto3";

package %s;

option go_package = "%s";

`

const importTempl = `
{{- range $item := .}}
import "{{$item}}";
{{end -}}
`

const enumTempl = `
{{- range $enum := .}}
enum {{$enum.Name}} {
{{- range $item := $enum.List}}
	{{$item.Name}} = {{$item.Value}}; // {{$item.Desc}}
{{- end}}
}
{{end -}}

`

const structTempl = `
{{- range $st := .}}
message {{$st.Name}} {
{{- range $item := $st.List}}
	{{$item.Class}} {{$item.Name}} = {{$item.Position}}; // ${{$item.Desc}}
{{- end}}
}
{{end}}

`

const configTempl = `
{{- range $cfg := .}}
message {{$cfg.Name}} {
{{- range $item := $cfg.List}}
	{{$item.Class}} {{$item.Name}} = {{$item.Position}}; // ${{$item.Desc}}
{{- end}}
}

message {{$cfg.Name}}Ary {
	repeated {{$cfg.Name}} Ary = 1;
}
{{- end}}

`
