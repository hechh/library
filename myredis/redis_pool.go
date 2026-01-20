package myredis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hechh/library/yaml"
)

var (
	pools = make(map[string]*Client)
)

func Init(cfgs map[int32]*yaml.DbConfig) error {
	for _, cfg := range cfgs {
		// 建立redis连接
		cli := redis.NewClient(&redis.Options{
			IdleTimeout:  1 * time.Minute,
			MinIdleConns: 100,
			DB:           int(cfg.Db),
			ReadTimeout:  -1,
			WriteTimeout: -1,
			Addr:         cfg.Host,
			Username:     cfg.User,
			Password:     cfg.Password,
			OnConnect:    func(ctx context.Context, cn *redis.Conn) error { return nil },
		})
		// 连接到redis服务器，测试连通性
		if _, err := cli.Ping(context.Background()).Result(); err != nil {
			return err
		}
		pools[cfg.DbName] = NewClient(cli, cfg.Prefix)
	}
	return nil
}

func Close() {
	for _, cli := range pools {
		cli.Close()
	}
}

func Get(db string) *Client {
	return pools[db]
}
