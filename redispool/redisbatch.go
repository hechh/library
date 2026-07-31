package redispool

import (
	"github.com/hechh/library/base/safe"
)

func MGet(list ...IString) (map[string]Message, error) {
	items := make(map[uint32][]IString)
	for _, item := range list {
		uuid := item.GetClient().UniqueId()
		items[uuid] = append(items[uuid], item)
	}

	result := make(map[string]Message, len(list))
	for _, vals := range items {
		client := vals[0].GetClient()

		// 构造MGet参数
		args := make([]string, 0, len(vals))
		for _, item := range vals {
			args = append(args, item.GetKey())
		}

		// 批量加载数据
		values, err := client.MGet(args...)
		if err != nil {
			return nil, err
		}

		// 解析数据
		for i, value := range values {
			var obj Message
			var err error
			switch v := value.(type) {
			case string:
				obj, err = vals[i].Unmarshal(safe.StringToBytes(v))
			case []byte:
				obj, err = vals[i].Unmarshal(v)
			case nil:
				obj, err = vals[i].Unmarshal(nil)
			}
			if err != nil {
				return nil, err
			}
			result[args[i]] = obj
		}
	}
	return result, nil
}

func MSet(list []IString, data map[string]Message) error {
	items := make(map[uint32][]IString)
	for _, item := range list {
		uuid := item.GetClient().UniqueId()
		items[uuid] = append(items[uuid], item)
	}

	var reterr error
	for _, vals := range items {
		client := vals[0].GetClient()

		// 构造MGet参数
		args := make([]any, 0, len(vals)*2)
		for _, item := range vals {
			key := item.GetKey()
			val, ok := data[key]
			if !ok {
				continue
			}

			buff, err := item.Marshal(val)
			if err != nil {
				return err
			}
			args = append(args, key, safe.BytesToString(buff))
		}

		// 批量保存数据
		if err := client.MSet(args...); err != nil {
			reterr = err
		}
	}
	return reterr
}

func HMGet(list ...IHash) (map[string]Message, error) {
	// 按 (UniqueId, Key) 分组，同一 hash key 下的多个 field 可合并为单次 HMGet
	type groupKey struct {
		uuid uint32
		key  string
	}
	items := make(map[groupKey][]IHash)
	for _, item := range list {
		gk := groupKey{uuid: item.GetClient().UniqueId(), key: item.GetKey()}
		items[gk] = append(items[gk], item)
	}

	result := make(map[string]Message, len(list))
	for _, vals := range items {
		client := vals[0].GetClient()
		key := vals[0].GetKey()

		// 构造 HMGet 的 fields 参数
		fields := make([]string, 0, len(vals))
		for _, item := range vals {
			fields = append(fields, item.GetField())
		}

		// 批量加载 hash field 数据
		values, err := client.HMGet(key, fields...)
		if err != nil {
			return nil, err
		}

		// 解析数据
		for i, value := range values {
			var obj Message
			var err error
			switch v := value.(type) {
			case string:
				obj, err = vals[i].Unmarshal(safe.StringToBytes(v))
			case []byte:
				obj, err = vals[i].Unmarshal(v)
			case nil:
				obj, err = vals[i].Unmarshal(nil)
			}
			if err != nil {
				return nil, err
			}
			result[fields[i]] = obj
		}
	}
	return result, nil
}

func HMSet(list []IHash, data map[string]Message) error {
	// 按 (UniqueId, Key) 分组，同一 hash key 下的多个 field 可合并为单次 HMSet
	type groupKey struct {
		uuid uint32
		key  string
	}
	items := make(map[groupKey][]IHash)
	for _, item := range list {
		gk := groupKey{uuid: item.GetClient().UniqueId(), key: item.GetKey()}
		items[gk] = append(items[gk], item)
	}

	var reterr error
	for _, vals := range items {
		client := vals[0].GetClient()
		key := vals[0].GetKey()

		// 构造 HMSet 的 field-value 参数对
		args := make([]any, 0, len(vals)*2)
		for _, item := range vals {
			field := item.GetField()
			val, ok := data[field]
			if !ok {
				continue
			}

			buff, err := item.Marshal(val)
			if err != nil {
				return err
			}
			args = append(args, field, safe.BytesToString(buff))
		}

		// 批量保存 hash field 数据
		if err := client.HMSet(key, args...); err != nil {
			reterr = err
		}
	}
	return reterr
}
