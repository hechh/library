package xlsx2proto

import (
	"fmt"
	"sort"
	"strings"
)

// enum 结构
type Value struct {
	class string
	name  string
	value int32
	desc  string
}
type EnumDescriptor struct {
	name string
	list []*Value
	data map[string]*Value
}

func NewEnumDescriptor(name string) *EnumDescriptor {
	return &EnumDescriptor{
		name: name,
		data: make(map[string]*Value),
	}
}

// E|游戏类型-德州NORMAL|GameType|Normal|1
func (d *EnumDescriptor) Put(val int32, name string, gameType string, Desc string) {
	item := &Value{
		class: gameType,
		name:  name,
		value: val,
		desc:  Desc,
	}
	d.list = append(d.list, item)
	d.data[item.desc] = item
}

func (d *EnumDescriptor) String() string {
	sort.Slice(d.list, func(i int, j int) bool {
		return d.list[i].value < d.list[j].value
	})
	strs := []string{}
	for _, item := range d.list {
		strs = append(strs, fmt.Sprintf("\t%s\t=\t%d;\t// %s", item.name, item.value, item.desc))
	}
	return fmt.Sprintf("enum %s {\n%s\n}\n\n", d.name, strings.Join(strs, "\n"))
}
