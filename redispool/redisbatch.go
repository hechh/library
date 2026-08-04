package redispool

import (
	"github.com/hechh/library/base/safe"
)

type mget struct {
	values []*Value
	args   []string
}

func (d *mget) Save() error {
	results, err := d.values[0].MGet(d.args...)
	if err != nil {
		return err
	}
	for i, val := range results {
		var err error
		switch vv := val.(type) {
		case string:
			err = d.values[i].UnmarshalVT(safe.StringToBytes(vv))
		case []byte:
			err = d.values[i].UnmarshalVT(vv)
		default:
			err = d.values[i].UnmarshalVT(nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type hmget struct {
	key    string
	values []*Value
	args   []string
}

func (d *hmget) Save() error {
	results, err := d.values[0].HMGet(d.key, d.args...)
	if err != nil {
		return err
	}
	for i, val := range results {
		var err error
		switch vv := val.(type) {
		case string:
			err = d.values[i].UnmarshalVT(safe.StringToBytes(vv))
		case []byte:
			err = d.values[i].UnmarshalVT(vv)
		default:
			err = d.values[i].UnmarshalVT(nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func MGet(args ...*Value) error {
	datas := map[uint32][]*Value{}
	for _, item := range args {
		id := item.UniqueId()
		datas[id] = append(datas[id], item)
	}

	for _, values := range datas {
		args := make([]string, 0, len(values))
		for _, item := range values {
			args = append(args, item.Key())
		}

		// 加载数据
		results, err := values[0].MGet(args...)
		if err != nil {
			return err
		}

		// 解析数据
		for i, val := range results {
			var err error
			switch vv := val.(type) {
			case string:
				err = values[i].UnmarshalVT(safe.StringToBytes(vv))
			case []byte:
				err = values[i].UnmarshalVT(vv)
			default:
				err = values[i].UnmarshalVT(nil)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MSet(args ...*Value) error {
	datas := map[uint32][]*Value{}
	for _, item := range args {
		id := item.UniqueId()
		datas[id] = append(datas[id], item)
	}

	for _, values := range datas {
		// 构造参数
		args := make([]any, 0, len(values)*2)
		for _, item := range values {
			data, err := item.MarshalVT()
			if err != nil {
				return err
			}
			args = append(args, item.Key(), safe.BytesToString(data))
		}
		// 保存数据
		if err := values[0].MSet(args...); err != nil {
			return err
		}
	}
	return nil
}

func HMGet(values ...*Value) error {
	datas := map[string][]*Value{}
	for _, item := range values {
		datas[item.Key()] = append(datas[item.Key()], item)
	}

	for key, values := range datas {
		args := make([]string, 0, len(values))
		for _, item := range values {
			args = append(args, item.Field())
		}

		// 加载数据
		results, err := values[0].HMGet(key, args...)
		if err != nil {
			return err
		}

		// 解析数据
		for i, value := range results {
			var err error
			switch vv := value.(type) {
			case string:
				err = values[i].UnmarshalVT(safe.StringToBytes(vv))
			case []byte:
				err = values[i].UnmarshalVT(vv)
			default:
				err = values[i].UnmarshalVT(nil)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func HMSet(args ...*Value) error {
	datas := map[string][]*Value{}
	for _, item := range args {
		datas[item.Key()] = append(datas[item.Key()], item)
	}

	for key, values := range datas {
		// 构造参数
		args := make([]any, 0, len(values)*2)
		for _, item := range values {
			data, err := item.MarshalVT()
			if err != nil {
				return err
			}
			args = append(args, item.Field(), safe.BytesToString(data))
		}
		// 保存数据
		if err := values[0].HMSet(key, args...); err != nil {
			return err
		}
	}
	return nil
}
