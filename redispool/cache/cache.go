package cache

import (
	"github.com/hechh/library/redispool"
)

type Value struct {
	value any
	mask  uint32
	times uint32
}

type Cache struct {
	values map[string]*Value
	types  map[uint32]redispool.IData
}

func (d *Value) Get() any      { return d.value }
func (d *Value) Mask() uint32  { return d.mask }
func (d *Value) Times() uint32 { return d.times }

// NewCache 创建空缓存
func New(vals map[string]any, typs map[uint32]redispool.IData) *Cache {
	values := make(map[string]*Value)
	for kk, vv := range vals {
		values[kk] = &Value{value: vv}
	}
	if typs == nil {
		typs = make(map[uint32]redispool.IData)
	}
	return &Cache{values: values, types: typs}
}

// GetTypes 返回已注册的数据类型描述（供 MGet/MSet 批量加载）
func (d *Cache) GetTypes(items ...redispool.IData) []redispool.IData {
	filter := map[uint32]struct{}{}
	rets := make([]redispool.IData, 0, len(items)+len(d.types))
	for _, v := range d.types {
		id := v.UniqueId()
		if _, ok := filter[id]; !ok {
			filter[id] = struct{}{}
			rets = append(rets, v)
		}
	}
	return rets
}

// AddType 注册数据类型描述
func (d *Cache) AddType(t redispool.IData) {
	d.types[t.UniqueId()] = t
}

// SetCache 写入缓存值，flag 记录数据来源标志（已存在则更新值并覆盖标志）
func (d *Cache) SetCache(key string, value any, flag uint32) {
	if val, ok := d.values[key]; ok {
		val.value = value
		val.mask = flag
		return
	}
	d.values[key] = &Value{value: value, mask: flag}
}

// GetCache 读取缓存值
func (d *Cache) GetCache(key string) (any, bool) {
	if v, ok := d.values[key]; ok {
		return v.value, ok
	}
	return nil, false
}

// IsChanged 判断 key 是否被标记为已变更
func (d *Cache) IsChanged(key string) bool {
	if v, ok := d.values[key]; ok {
		return v.times > 0
	}
	return false
}

// Change 标记 key 为已变更（计数累加，供持久化层判定脏数据）
func (d *Cache) Change(key string) {
	if v, ok := d.values[key]; ok {
		v.times++
	}
}

// Reset 清除 key 的变更标记
func (d *Cache) Reset(key string) {
	if v, ok := d.values[key]; ok {
		v.times = 0
	}
}
