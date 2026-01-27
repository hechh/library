package xlsx2data

// 枚举类型
type value struct {
	class string
	name  string
	value int32
	desc  string
}
type EnumDescriptor struct {
	name string
	list []*value
	data map[string]*value
}

func NewEnumDescriptor(name string) *EnumDescriptor {
	return &EnumDescriptor{
		name: name,
		data: make(map[string]*value),
	}
}

// E|游戏类型-德州NORMAL|GameType|Normal|1
func (d *EnumDescriptor) Put(val int32, name string, gameType string, desc string) {
	item := &value{
		class: gameType,
		name:  name,
		value: val,
		desc:  desc,
	}
	d.list = append(d.list, item)
	d.data[item.desc] = item
}

func (d *EnumDescriptor) ToInt32(val string) int32 {
	if val, ok := d.data[val]; ok {
		return val.value
	}
	return 0
}
