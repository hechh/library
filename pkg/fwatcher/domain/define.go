package domain

import (
	"hash"

	"google.golang.org/protobuf/proto"
)

type EtcdConfig struct {
	PrefixTopic string   `yaml:"prefix_topic,omitempty"` // etcd 前缀主题
	Endpoints   []string `yaml:"endpoints,omitempty"`    // etcd 节点地址列表
	KeepAlive   int64    `yaml:"keep_alive,omitempty"`   // 保活时间（秒）
}

type Config struct {
	DataPath string      `yaml:"data_path,omitempty"` // 数据文件目录
	XlsxPath string      `yaml:"xlsx_path,omitempty"` // Excel 文件目录
	Ext      string      `yaml:"ext,omitempty"`       // 文件扩展名
	Etcd     *EtcdConfig `yaml:"etcd,omitempty"`      // etcd 配置
}

// 配置解析接口
type IParser interface {
	RegisterChange(...func())
	Sheet() string
	New([]byte) (proto.Message, error)
	GetValue() string
	Parse(hash.Hash, []byte) error
}

// 配置同步接口
type ISync interface {
	Init(*Config) error
	Close()
	Put(string, []byte) error
	Update(string, []byte) error
	Delete(string) error
	Watch(func(string, []byte)) error
}
