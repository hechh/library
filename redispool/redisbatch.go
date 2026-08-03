package redispool

import (
	"github.com/hechh/library/base/logic"
	"github.com/hechh/library/base/safe"
	"github.com/hechh/library/base/templ"
	"github.com/hechh/library/base/tuple"
)

func Load(ca ICache, list ...IData) (map[string]Message, error) {
	type data struct {
		client IClient
		key    string
		list   []IData
		args   []string
	}
	list = ca.GetTypes(list...)
	result := make(map[string]Message, len(list))
	items := make(map[tuple.Tuple3[uint32, uint32, uint32]]*data)
	for _, item := range list {
		key := item.GetKey()
		field := item.GetField()
		cacheKey := key + field
		if value, ok := ca.GetCache(cacheKey); ok {
			result[cacheKey] = value.(Message)
			continue
		}
		mask := item.GetMask()
		cid := item.GetClient().UniqueId()
		kk := tuple.T3(mask, cid, templ.Or(logic.Has(mask, HASH_FLAG), item.UniqueId(), 0))
		vv, ok := items[kk]
		if !ok {
			vv = &data{client: item.GetClient()}
			items[kk] = vv
		}
		vv.list = append(vv.list, item)
		if logic.Has(item.GetMask(), STRING_FLAG) {
			vv.args = append(vv.args, key)
		} else if logic.Has(item.GetMask(), HASH_FLAG) {
			vv.args = append(vv.args, field)
			vv.key = key
		}
	}
	for kk, vv := range items {
		var values []any
		var err error
		if logic.Has(kk.V1, STRING_FLAG) {
			values, err = vv.client.MGet(vv.args...)
		} else if logic.Has(kk.V1, HASH_FLAG) {
			values, err = vv.client.HMGet(vv.key, vv.args...)
		}
		if err != nil {
			return nil, err
		}
		for i, value := range values {
			obj, err := unmarshal(vv.list[i], value)
			if err != nil {
				return nil, err
			}
			hkey := templ.Or(logic.Has(kk.V1, STRING_FLAG), vv.args[i], vv.key+vv.args[i])
			result[hkey] = obj
			ca.SetCache(hkey, obj, vv.list[i].GetMask())
		}
	}
	return result, nil
}

func Save(ca ICache, list ...IData) (reterr error) {
	type data struct {
		client IClient
		key    string
		args   []any
	}
	list = ca.GetTypes(list...)
	items := make(map[tuple.Tuple3[uint32, uint32, uint32]]*data)
	for _, item := range list {
		key := item.GetKey()
		field := item.GetField()
		var buff []byte
		cacheKey := key + field
		if !ca.IsChanged(cacheKey) {
			continue
		}
		val, _ := ca.GetCache(cacheKey)
		buff, reterr = item.Marshal(val.(Message))
		if reterr != nil {
			return
		}
		mask := item.GetMask()
		cid := item.GetClient().UniqueId()
		kk := tuple.T3(mask, cid, templ.Or(logic.Has(mask, HASH_FLAG), item.UniqueId(), 0))
		vv, ok := items[kk]
		if !ok {
			vv = &data{client: item.GetClient()}
			items[kk] = vv
		}
		if logic.Has(kk.V1, STRING_FLAG) {
			vv.args = append(vv.args, key, safe.BytesToString(buff))
		} else if logic.Has(kk.V1, HASH_FLAG) {
			vv.args = append(vv.args, field, safe.BytesToString(buff))
			vv.key = key
		}
	}
	// 批量保存 hash field 数据
	for _, vv := range items {
		if err := vv.client.HMSet(vv.key, vv.args...); err != nil {
			reterr = err
		}
	}
	return nil
}

func MGet(ca ICache, list ...IData) (map[string]Message, error) {
	type data struct {
		client IClient
		list   []IData
		args   []string
	}
	list = ca.GetTypes(list...)
	result := make(map[string]Message, len(list))
	items := make(map[tuple.Tuple2[uint32, uint32]]*data)
	for _, item := range list {
		mask := item.GetMask()
		if logic.Has(mask, HASH_FLAG) {
			continue
		}
		key := item.GetKey()
		if value, ok := ca.GetCache(key); ok {
			result[key] = value.(Message)
			continue
		}
		kk := tuple.T2(mask, item.GetClient().UniqueId())
		vv, ok := items[kk]
		if !ok {
			vv = &data{client: item.GetClient()}
			items[kk] = vv
		}
		vv.list = append(vv.list, item)
		vv.args = append(vv.args, key)
	}
	// 批量加载数据
	for _, vv := range items {
		values, err := vv.client.MGet(vv.args...)
		if err != nil {
			return nil, err
		}
		// 解析数据
		for i, value := range values {
			obj, err := unmarshal(vv.list[i], value)
			if err != nil {
				return nil, err
			}
			result[vv.args[i]] = obj
			ca.SetCache(vv.args[i], obj, vv.list[i].GetMask())
		}
	}
	return result, nil
}

func MSet(ca ICache, list ...IData) (reterr error) {
	type data struct {
		client IClient
		args   []any
	}
	list = ca.GetTypes(list...)
	items := make(map[tuple.Tuple2[uint32, uint32]]*data)
	for _, item := range list {
		mask := item.GetMask()
		if logic.Has(mask, HASH_FLAG) {
			continue
		}
		key := item.GetKey()
		if !ca.IsChanged(key) {
			continue
		}
		val, _ := ca.GetCache(key)
		buff, err := item.Marshal(val.(Message))
		if err != nil {
			return err
		}
		kk := tuple.T2(mask, item.GetClient().UniqueId())
		vv, ok := items[kk]
		if !ok {
			vv = &data{client: item.GetClient()}
			items[kk] = vv
		}
		vv.args = append(vv.args, key, safe.BytesToString(buff))
	}
	// 批量保存数据
	for _, vv := range items {
		if err := vv.client.MSet(vv.args...); err != nil {
			reterr = err
		}
	}
	return
}

func HMGet(ca ICache, list ...IData) (map[string]Message, error) {
	type data struct {
		client IClient
		key    string
		list   []IData
		args   []string
	}
	list = ca.GetTypes(list...)
	result := make(map[string]Message, len(list))
	items := make(map[tuple.Tuple3[uint32, uint32, uint32]]*data)
	for _, item := range list {
		mask := item.GetMask()
		if logic.Has(mask, STRING_FLAG) {
			continue
		}
		key := item.GetKey()
		field := item.GetField()
		cacheKey := key + field
		if value, ok := ca.GetCache(cacheKey); ok {
			result[cacheKey] = value.(Message)
			continue
		}
		kk := tuple.T3(mask, item.GetClient().UniqueId(), item.UniqueId())
		vv, ok := items[kk]
		if !ok {
			vv = &data{client: item.GetClient(), key: key}
			items[kk] = vv
		}
		vv.list = append(vv.list, item)
		vv.args = append(vv.args, field)
	}
	// 批量加载 hash field 数据
	for _, vv := range items {
		values, err := vv.client.HMGet(vv.key, vv.args...)
		if err != nil {
			return nil, err
		}
		// 解析数据
		for i, value := range values {
			obj, err := unmarshal(vv.list[i], value)
			if err != nil {
				return nil, err
			}
			cacheKey := vv.key + vv.args[i]
			result[cacheKey] = obj
			ca.SetCache(cacheKey, obj, vv.list[i].GetMask())
		}
	}
	return result, nil
}

func HMSet(ca ICache, list ...IData) (reterr error) {
	type data struct {
		client IClient
		key    string
		args   []any
	}
	list = ca.GetTypes(list...)
	items := make(map[tuple.Tuple3[uint32, uint32, uint32]]*data)
	for _, item := range list {
		mask := item.GetMask()
		if logic.Has(mask, STRING_FLAG) {
			continue
		}
		key := item.GetKey()
		field := item.GetField()
		cacheKey := key + field
		if !ca.IsChanged(cacheKey) {
			continue
		}
		val, _ := ca.GetCache(cacheKey)
		buff, err := item.Marshal(val.(Message))
		if err != nil {
			return err
		}
		kk := tuple.T3(mask, item.GetClient().UniqueId(), item.UniqueId())
		vv, ok := items[kk]
		if !ok {
			vv = &data{client: item.GetClient(), key: key}
			items[kk] = vv
		}
		vv.args = append(vv.args, field, safe.BytesToString(buff))
	}
	// 批量保存 hash field 数据
	for _, vv := range items {
		if err := vv.client.HMSet(vv.key, vv.args...); err != nil {
			reterr = err
		}
	}
	return
}
