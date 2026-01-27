package xlsx2code

import (
	"fmt"
	"strings"

	"github.com/hechh/library/toolkit"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Token int32

const (
	IDENT   Token = 0
	POINTER Token = 1
	ARRAY   Token = 2
	MAP     Token = 3
	GROUP   Token = 4
)

type Field struct {
	protoreflect.FieldDescriptor
	Name string
	Type string
}

type Index struct {
	Kind Token
	Name string
	Type string
	List []*Field
}

type StructDescriptor struct {
	protoreflect.MessageDescriptor
	Name string
	Data map[string]*Index
	List []*Index
	data map[string]*Field
	list []*Field
}

func NewStructDescriptor(name string, msgType protoreflect.MessageDescriptor) *StructDescriptor {
	return &StructDescriptor{
		MessageDescriptor: msgType,
		Name:              name,
		Data:              make(map[string]*Index),
		data:              make(map[string]*Field),
	}
}

func (d *StructDescriptor) AddIndex(rule string) {
	strs := strings.Split(rule, ":")
	item := &Index{Name: strings.ReplaceAll(strs[1], ",", "")}
	for _, name := range strings.Split(strs[1], ",") {
		item.List = append(item.List, d.data[name])
	}
	switch strings.ToLower(strs[0]) {
	case "map":
		item.Kind = MAP
		item.Type = fmt.Sprintf("util.Map%d", len(item.List))
	case "group":
		item.Kind = GROUP
		item.Type = fmt.Sprintf("util.Group%d", len(item.List))
	}
	d.List = append(d.List, item)
	d.Data[item.Name] = item
}

func (d *StructDescriptor) Put(pos int32, name, tname string) {
	item := d.parse(tname)
	item.Name = name
	item.FieldDescriptor = d.Fields().ByName(protoreflect.Name(name))
	d.list = append(d.list, item)
	d.data[item.Name] = item
}

func (d *StructDescriptor) parse(str string) *Field {
	if strings.HasPrefix(str, "[]") {
		item := d.parse(strings.TrimPrefix(str, "[]"))
		return &Field{Type: item.Type}
	}
	if strings.HasPrefix(str, "&") {
		item := d.parse(strings.TrimPrefix(str, "&"))
		return &Field{Type: item.Type}
	}
	if strings.HasPrefix(str, "*") {
		item := d.parse(strings.TrimPrefix(str, "*"))
		return &Field{Type: item.Type}
	}
	return &Field{Type: toolkit.Target(str)}
}

func (d *Index) GetArg() string {
	strs := []string{}
	for _, val := range d.List {
		if val.Kind() == protoreflect.EnumKind {
			strs = append(strs, val.Name+" pb."+val.Type)
		} else {
			strs = append(strs, val.Name+" "+val.Type)
		}
	}
	return strings.Join(strs, ", ")
}

func (d *Index) GetType() string {
	strs := []string{}
	for _, val := range d.List {
		if val.Kind() == protoreflect.EnumKind {
			strs = append(strs, "pb."+val.Type)
		} else {
			strs = append(strs, val.Type)
		}
	}
	return strings.Join(strs, ", ")
}

func (d *Index) GetValue(ref string) string {
	strs := []string{}
	for _, val := range d.List {
		if len(ref) > 0 {
			strs = append(strs, ref+"."+val.Name)
		} else {
			strs = append(strs, val.Name)
		}
	}
	return strings.Join(strs, ", ")
}
