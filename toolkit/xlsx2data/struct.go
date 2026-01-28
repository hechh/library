package xlsx2data

import (
	"strings"

	"github.com/hechh/library/toolkit"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type Token int32

const (
	IDENT   Token = 0
	POINTER Token = 1
	ARRAY   Token = 2
	MAP     Token = 3
	GROUP   Token = 4
)

// 结构类型
type Field struct {
	fieldType protoreflect.FieldDescriptor
	token     Token
	class     string
	name      string
	position  int32
}
type StructDescriptor struct {
	aryType protoreflect.MessageDescriptor
	cfgType protoreflect.MessageDescriptor
	name    string
	list    []*Field
	data    map[string]*Field
	rows    [][]string
}

func NewStructDescriptor(name string, ary, cfg protoreflect.MessageDescriptor, rows [][]string) *StructDescriptor {
	return &StructDescriptor{
		aryType: ary,
		cfgType: cfg,
		name:    name,
		data:    make(map[string]*Field),
		rows:    rows,
	}
}

func (d *StructDescriptor) Put(pos int32, name, tname string) {
	item := d.parse(tname)
	item.name = name
	item.position = pos
	item.fieldType = d.cfgType.Fields().ByName(protoreflect.Name(name))
	d.list = append(d.list, item)
	d.data[item.name] = item
}

func (d *StructDescriptor) parse(str string) *Field {
	if strings.HasPrefix(str, "[]") {
		item := d.parse(strings.TrimPrefix(str, "[]"))
		return &Field{class: item.class, token: ARRAY}
	}
	if strings.HasPrefix(str, "&") {
		item := d.parse(strings.TrimPrefix(str, "&"))
		return &Field{class: item.class, token: IDENT}
	}
	if strings.HasPrefix(str, "*") {
		item := d.parse(strings.TrimPrefix(str, "*"))
		return &Field{class: item.class, token: POINTER}
	}
	return &Field{class: str, token: IDENT}
}

func (d *StructDescriptor) Marshal() ([]byte, error) {
	ary := dynamicpb.NewMessage(d.aryType)
	list := ary.Mutable(d.aryType.Fields().ByName("Ary")).List()
	for _, line := range d.rows {
		cfg := dynamicpb.NewMessage(d.cfgType)
		for pos, field := range d.list {
			if pos+1 > len(line) {
				break
			}
			switch field.token {
			case IDENT, POINTER:
				cfg.Set(field.fieldType, convert(field, line[field.position-1]))
			case ARRAY:
				fieldList := cfg.Mutable(field.fieldType).List()
				for _, vv := range strings.Split(line[field.position-1], "|") {
					fieldList.Append(convert(field, vv))
				}
			}
		}
		list.Append(protoreflect.ValueOf(cfg))
	}

	marshaler := prototext.MarshalOptions{Multiline: true}
	return marshaler.Marshal(ary)
}

func convert(field *Field, val string) protoreflect.Value {
	value := toolkit.Convert(field.class, val)
	switch field.fieldType.Kind() {
	case protoreflect.BoolKind:
		return protoreflect.ValueOf(value)
	case protoreflect.EnumKind:
		return protoreflect.ValueOfEnum(protoreflect.EnumNumber(value.(int32)))
	case protoreflect.Int32Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.Sint32Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.Uint32Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.Int64Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.Sint64Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.Uint64Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.Sfixed32Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.Fixed32Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.FloatKind:
		return protoreflect.ValueOf(value)
	case protoreflect.Sfixed64Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.Fixed64Kind:
		return protoreflect.ValueOf(value)
	case protoreflect.DoubleKind:
		return protoreflect.ValueOf(value)
	case protoreflect.StringKind:
		return protoreflect.ValueOf(value)
	case protoreflect.BytesKind:
		return protoreflect.ValueOf(value)
	case protoreflect.MessageKind:
		return protoreflect.ValueOf(value)
	case protoreflect.GroupKind:
		return protoreflect.ValueOf(value)
	}
	return protoreflect.ValueOf(nil)
}
