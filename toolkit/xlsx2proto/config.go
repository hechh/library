package xlsx2proto

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hechh/library/toolkit"
)

// struct 结构
type Field struct {
	class    string
	name     string
	position int32
	desc     string
}

// struct 结构
type StructDescriptor struct {
	name string
	list []*Field
	data map[string]*Field
}

type ConfigDescriptor struct {
	StructDescriptor
}

func NewStructDescriptor(name string) *StructDescriptor {
	return &StructDescriptor{
		name: name,
		data: make(map[string]*Field),
	}
}

func NewConfigDescriptor(name string) *ConfigDescriptor {
	return &ConfigDescriptor{
		StructDescriptor: StructDescriptor{
			name: name,
			data: make(map[string]*Field),
		},
	}
}

func (d *StructDescriptor) Put(pos int32, name, tname, desc string) {
	item := d.parse(tname)
	item.name = name
	item.position = pos
	item.desc = desc
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
	return &Field{class: toolkit.Target(str)}
}

func (d *StructDescriptor) String() string {
	sort.Slice(d.list, func(i int, j int) bool {
		return d.list[i].position < d.list[j].position
	})
	strs := []string{}
	for _, item := range d.list {
		strs = append(strs, fmt.Sprintf("\t%s %s = %d;\t// %s",
			item.class, item.name, item.position, item.desc))
	}
	return fmt.Sprintf("message %s {\n%s\n}\n\n", d.name, strings.Join(strs, "\n"))
}

func (d *ConfigDescriptor) String() string {
	sort.Slice(d.list, func(i int, j int) bool {
		return d.list[i].position < d.list[j].position
	})
	strs := []string{}
	for _, item := range d.list {
		strs = append(strs, fmt.Sprintf("\t%s %s = %d;\t// %s",
			item.class, item.name, item.position, item.desc))
	}
	return fmt.Sprintf("message %s {\n%s\n}\n\nmessage %sAry {\nrepeated %s Ary = 1;\n}\n\n", d.name, strings.Join(strs, "\n"), d.name, d.name)
}
