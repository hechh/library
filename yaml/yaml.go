package yaml

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SlaveConfig struct {
	DbName   string `yaml:"dbname"`
	Db       int    `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type DbConfig struct {
	DbName   string                 `yaml:"dbname"`
	Db       int32                  `yaml:"db"`
	Prefix   string                 `yaml:"prefix"`
	User     string                 `yaml:"user"`
	Password string                 `yaml:"password"`
	Host     string                 `yaml:"host"`
	Slave    map[int32]*SlaveConfig `yaml:"slave"`
}

type EtcdConfig struct {
	Topic     string   `yaml:"topic"`
	Endpoints []string `yaml:"endpoints"`
}

type NatsConfig struct {
	Topic     string `yaml:"topic"`
	Endpoints string `yaml:"endpoints"`
}

type NodeConfig struct {
	RouterExpire   int64  `yaml:"router_expire"`
	RegisterExpire int64  `yaml:"register_expire"`
	LogLevel       string `yaml:"log_level"`
	LogPath        string `yaml:"log_path"`
	Ip             string `yaml:"ip"`
	Port           int    `yaml:"port"`
}

type CommonConfig struct {
	Env         int32  `yaml:"env"`
	Mode        int32  `yaml:"mode"`
	IsOpenPprof bool   `yaml:"is_open_pprof"`
	TablePath   string `yaml:"table_path"`
	TokenKey    string `yaml:"token_key"`
	AesKey      string `yaml:"aes_key"`
}

func Load(filename string, cfg any) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, cfg)
}
