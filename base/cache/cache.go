package cache

type ICache interface {
	WalkCache(func(string, any, uint32, uint32))
	SetCache(string, any, uint32)
	GetCache(string) (any, bool)
	IsChanged(string) bool
	Change(string)
	Reset(string)
}

type Value struct {
	data  any
	times uint32
	mask  uint32
}

type Cache struct {
	values map[string]*Value
}

func New() *Cache {
	return &Cache{values: make(map[string]*Value)}
}

func (d *Cache) WalkCache(f func(k string, v any, times uint32, mask uint32)) {
	for k, v := range d.values {
		f(k, v.data, v.times, v.mask)
	}
}

func (d *Cache) SetCache(key string, value any, flag uint32) {
	if val, ok := d.values[key]; ok {
		val.data = value
		val.mask = flag
		return
	}
	d.values[key] = &Value{data: value, mask: flag}
}

func (d *Cache) GetCache(key string) (any, bool) {
	if v, ok := d.values[key]; ok {
		return v.data, ok
	}
	return nil, false
}

func (d *Cache) IsChanged(key string) bool {
	if v, ok := d.values[key]; ok {
		return v.times > 0
	}
	return false
}

func (d *Cache) Change(key string) {
	if v, ok := d.values[key]; ok {
		v.times++
	}
}

func (d *Cache) Reset(key string) {
	if v, ok := d.values[key]; ok {
		v.times = 0
	}
}
