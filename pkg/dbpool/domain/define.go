package domain

import "xorm.io/xorm"

// DbConfig 数据库分片配置
type DbConfig struct {
	ShardsId uint32              `yaml:"shards_id,omitempty"` // 分片ID
	DbName   string              `yaml:"db_name,omitempty"`   // 数据库名称
	Db       uint32              `yaml:"db,omitempty"`        // 数据库编号
	User     string              `yaml:"user,omitempty"`      // 数据库用户名
	Password string              `yaml:"password,omitempty"`  // 数据库密码
	Ip       string              `yaml:"ip,omitempty"`        // 数据库IP地址
	Port     uint32              `yaml:"port,omitempty"`      // 数据库端口
	Slaves   map[int32]*DbConfig `yaml:"slaves,omitempty"`    // 从库配置列表
}

// Config 数据库连接池配置
type Config struct {
	UidModValue uint64      `yaml:"uid_mod_value,omitempty"` // 用户ID取模值
	ShardsSize  int64       `yaml:"shards_size,omitempty"`   // 分片数量
	Globals     []*DbConfig `yaml:"globals,omitempty"`       // 全局数据库配置列表
	Shards      []*DbConfig `yaml:"shards,omitempty"`        // 分片数据库配置列表
}

type IClient interface {
	Init(*DbConfig, ...any) error
	Close() error
	Connect() error
	Ping() error
	IsAlive() bool
	Engine() *xorm.EngineGroup
	NewSession() *xorm.Session
}
