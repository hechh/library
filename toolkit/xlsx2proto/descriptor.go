package xlsx2proto

import (
	"sort"
	"strings"

	"github.com/hechh/library/toolkit"
)

// enum 结构
type Value struct {
	Class string
	Name  string
	Value int32
	Desc  string
}
type EnumDescriptor struct {
	Name string
	List []*Value
	data map[string]*Value
}

func NewEnumDescriptor(Name string) *EnumDescriptor {
	return &EnumDescriptor{
		Name: Name,
		data: make(map[string]*Value),
	}
}

// E|游戏类型-德州NORMAL|GameType|Normal|1
func (d *EnumDescriptor) Put(val int32, Name string, gameType string, Desc string) {
	item := &Value{
		Class: gameType,
		Name:  Name,
		Value: val,
		Desc:  Desc,
	}
	d.List = append(d.List, item)
	d.data[item.Desc] = item
}

func (d *EnumDescriptor) Sort() {
	sort.Slice(d.List, func(i int, j int) bool {
		return d.List[i].Value < d.List[j].Value
	})
}

// struct 结构
type Field struct {
	Class    string
	Name     string
	Position int32
	Desc     string
}

// struct 结构
type StructDescriptor struct {
	Name string
	List []*Field
	data map[string]*Field
}

func NewStructDescriptor(Name string) *StructDescriptor {
	return &StructDescriptor{
		Name: Name,
		data: make(map[string]*Field),
	}
}

func (d *StructDescriptor) Put(pos int32, Name, tname, Desc string) {
	item := d.parse(tname)
	item.Name = Name
	item.Position = pos
	item.Desc = Desc
	d.List = append(d.List, item)
	d.data[item.Name] = item
}

func (d *StructDescriptor) parse(str string) *Field {
	if strings.HasPrefix(str, "[]") {
		item := d.parse(strings.TrimPrefix(str, "[]"))
		return &Field{Class: "repeated " + item.Class}
	}
	if strings.HasPrefix(str, "&") {
		item := d.parse(strings.TrimPrefix(str, "&"))
		return &Field{Class: item.Class}
	}
	if strings.HasPrefix(str, "*") {
		item := d.parse(strings.TrimPrefix(str, "*"))
		return &Field{Class: item.Class}
	}
	return &Field{Class: toolkit.Target(str)}
}

func (d *StructDescriptor) Sort() {
	sort.Slice(d.List, func(i int, j int) bool {
		return d.List[i].Position < d.List[j].Position
	})
}
