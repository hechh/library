package database

import (
	"time"

	"github.com/hechh/library/async"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/yaml"
)

const (
	MysqlDriver      = "mysql"
	PostgreSqlDriver = "postgres"
)

var (
	clients = make(map[string]*Client)
	tables  = make(map[string][]interface{})
	exit    = make(chan struct{})
)

func Register(dbname string, tabs ...interface{}) {
	if len(tabs) > 0 {
		tables[dbname] = append(tables[dbname], tabs...)
	}
}

func Init(driverName string, cfgs map[int32]*yaml.DbConfig) error {
	for _, cfg := range cfgs {
		cli := NewClient(driverName, cfg)
		if err := cli.Connect(tables[cfg.DbName]...); err != nil {
			return err
		}
		clients[cfg.DbName] = cli
	}
	async.Go(check)
	return nil
}

func Close() {
	close(exit)
	for _, cli := range clients {
		cli.Close()
	}
}

func Get(dbname string) *Client {
	client, ok := clients[dbname]
	if ok && client.IsAlive() {
		return client
	}
	return nil
}

func check() {
	tt := time.NewTicker(5 * time.Second)
	defer tt.Stop()
	for {
		select {
		case <-tt.C:
			// 检测连通信
			for _, client := range clients {
				if err := client.Ping(); err == nil {
					continue
				} else {
					mlog.Errorf("mysql连接异常断开: %v", err)
				}
				// 重新连接
				if err := client.Connect(); err != nil {
					mlog.Errorf("mysql重新连接失败, error:%v", err)
				}
			}
		case <-exit:
			return
		}
	}
}
