package cache

import "github.com/hechh/library/redispool"

type Wrapper struct {
	data map[string]any
}

func Wrap(data map[string]any) *Wrapper {
	return &Wrapper{data: data}
}

func (d *Wrapper) GetTypes() []redispool.IData { return nil }
func (d *Wrapper) AddType(redispool.IData)     {}
func (d *Wrapper) IsChanged(string) bool       { return true }
func (d *Wrapper) Change(string)               {}
func (d *Wrapper) Reset(string)                {}
func (d *Wrapper) SetCache(k string, v any, flag uint32) {
	d.data[k] = v
}
func (d *Wrapper) GetCache(k string) (any, bool) {
	val, ok := d.data[k]
	return val, ok
}
