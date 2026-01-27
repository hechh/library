package xlsx2proto

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hechh/library/convertor"
)

// struct 结构
type Field struct {
	class    string
	name     string
	Position int32
	Desc     string
}

type StructDescriptor struct {
	name string
	list []*Field
	data map[string]*Field
}

func NewStructDescriptor(name string) *StructDescriptor {
	return &StructDescriptor{
		name: name,
		data: make(map[string]*Field),
	}
}

func (d *StructDescriptor) Put(pos int32, name, tname, Desc string) {
	item := d.parse(tname)
	item.name = name
	item.Position = pos
	item.Desc = Desc
	d.list = append(d.list, item)
	d.data[item.name] = item
}

func (d *StructDescriptor) parse(str string) *Field {
	if strings.HasPrefix(str, "[]") {
		item := d.parse(strings.TrimPrefix(str, "[]"))
		return &Field{class: "repeated " + item.class}
	}
	if strings.HasPrefix(str, "&") {
		item := d.parse(strings.TrimPrefix(str, "&"))
		return &Field{class: item.class}
	}
	if strings.HasPrefix(str, "*") {
		item := d.parse(strings.TrimPrefix(str, "*"))
		return &Field{class: item.class}
	}
	return &Field{class: convertor.Target(str)}
}

func (d *StructDescriptor) String() string {
	sort.Slice(d.list, func(i int, j int) bool {
		return d.list[i].Position < d.list[j].Position
	})
	strs := []string{}
	for _, item := range d.list {
		strs = append(strs, fmt.Sprintf("\t%s %s = %d;\t// %s",
			item.class, item.name, item.Position, item.Desc))
	}
	return fmt.Sprintf("message %s {\n%s\n}\n\nmessage %sAry {\nrepeated %s Ary = 1;\n}\n\n", d.name, strings.Join(strs, "\n"), d.name, d.name)
}
