package pb2redis

const hashTempl = `
/*
* 本代码由pbtool工具生成，请勿手动修改
*/

package {{.Pkg}}

import (
	"fmt"
	pb "%s"
	
	"github.com/hechh/library/myredis"
	"github.com/hechh/library/uerror"
	"google.golang.org/protobuf/proto"
)


func GetKey({{GetArgs .Keys}}) string {
	return fmt.Sprintf("{{.KeyFmt}}" {{if .Keys}},{{end}} {{GetValues .Keys}})
}

func GetField({{GetArgs .Fields}}) string {
	return fmt.Sprintf("{{.FieldFmt}}" {{if .Fields}},{{end}} {{GetValues .Fields}})
}

func HGetAll({{GetArgs .Keys}}) (ret map[string]*pb.{{.Name}}, err error) {
	// 获取redis连接
	client := myredis.Get("{{.DbName}}")
	if client == nil {
		return nil, uerror.New(-1,"{{.DbName}}数据库不存在")
	}
	// 加载数据
	kvs, err := client.HGetAll(GetKey({{GetValues .Keys}}))
	if err != nil {
		return nil, err
	}
	// 解析数据
	ret = make(map[string]*pb.{{.Name}})
	for k, item := range kvs {
		if len(item) <= 0 {
			continue
		}
		data := &pb.{{$.Name}}{}
		if err := proto.Unmarshal([]byte(item), data); err != nil {
			return nil, err
		}
		ret[k] = data
	}
	return
}

func HMGet({{GetArgs .Keys}} {{if .Keys}},{{end}} fields ...string) (map[string]*pb.{{.Name}}, error) {
	if len(fields) <= 0 {
		return nil, nil	
	}
	// 获取redis连接
	client := myredis.Get("{{.DbName}}")
	if client == nil {
		return nil, uerror.New(-1,"{{.DbName}}数据库不存在")
	}
	// 批量获取
	results, err := client.HMGet(GetKey({{GetValues .Keys}}), fields...)
	if err != nil {
		return nil, err
	}
	// 解析数据
	ret := make(map[string]*pb.{{.Name}})
	for i, field := range fields {
		if results[i] == nil {
			continue
		}	
		var buf []byte
		switch vv := results[i].(type) {
		case string:
			buf = []byte(vv)	
		case []byte:
			buf = vv
		default:
			return nil, uerror.New(-1, "数据类型不支持")
		}
		item := &pb.{{.Name}}{}
		if err := proto.Unmarshal(buf, item); err != nil {
			return nil, err
		} else {
			ret[field] = item
		}
	}
	return ret, nil
}

func HMSet({{GetArgs .Keys}} {{if .Keys}},{{end}} data map[string]*pb.{{.Name}}) error {
	// 获取redis连接
	client := myredis.Get("{{.DbName}}")
	if client == nil {
		return uerror.New(-1,"{{.DbName}}数据库不存在")
	}
	// 设置数据
	vals := []interface{}{}
	for k, v := range data {
		buf, err := proto.Marshal(v)
		if err != nil {
			return err
		}
		vals = append(vals, k, buf)
	}
	return client.HMSet(GetKey({{GetValues .Keys}}), vals...)
}

func HGet({{GetArgs .Keys .Fields}}) (*pb.{{.Name}}, error) {
	// 获取redis连接
	client := myredis.Get("{{.DbName}}")
	if client == nil {
		return nil, uerror.New(-1,"{{.DbName}}数据库不存在")
	}
	// 加载数据
	key := GetKey({{GetValues .Keys}})
	field := GetField({{GetValues .Fields}})
	str, err := client.HGet(key, field)
	if err != nil {
		return nil, err
	}
	if len(str) <= 0 {
		return nil, nil
	}
	// 解析数据
	data := &pb.{{.Name}}{}
	err = proto.Unmarshal([]byte(str), data)
	return data, err 
}

func HSet({{GetArgs .Keys .Fields}}, data *pb.{{.Name}}) error {
	// 获取redis连接
	client := myredis.Get("{{.DbName}}")
	if client == nil {
		return uerror.New(-1,"{{.DbName}}数据库不存在")
	}
	// 序列化数据
	buf, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	key := GetKey({{GetValues .Keys}}) 
	field := GetField({{GetValues .Fields}})
	return client.HSet(key, field, buf)
}

func HDel({{GetArgs .Keys}} {{if .Keys}},{{end}} fields ...string) error {
	if len(fields) <= 0 {
		return nil	
	}
	// 获取redis连接
	client := myredis.Get("{{.DbName}}")
	if client == nil {
		return uerror.New(-1,"{{.DbName}}数据库不存在")
	}
	_, err := client.HDel(GetKey({{GetValues .Keys}}), fields...)
	return err
}


`
