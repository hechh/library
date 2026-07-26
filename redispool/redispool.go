package redispool

import (
	"github.com/hechh/library/redispool/domain"
)

type RedisPool struct {
	uidModValue uint64                    // 开始 Uid 值，用于路由
	shardsSize  uint64                    // 分片数量
	newFunc     func() domain.IClient     // new函数
	globals     map[string]domain.IClient // 全局数据库连接池
	shards      map[uint32]domain.IClient // 分片数据库连接池
}

func NewRedisPool[T any](f func() *T) *RedisPool {
	return &RedisPool{
		newFunc: func() domain.IClient { return any(f()).(domain.IClient) },
		globals: make(map[string]domain.IClient),
		shards:  make(map[uint32]domain.IClient),
	}
}

func (d *RedisPool) Init(cfg *domain.Config) error {
	// 初始化全局数据库
	for _, dbCfg := range cfg.Globals {
		cli := d.newFunc()
		if err := cli.Init(dbCfg); err != nil {
			d.Close()
			return err
		}
		d.globals[dbCfg.DbName] = cli
	}

	d.uidModValue = cfg.UidModValue
	d.shardsSize = uint64(cfg.ShardsSize)

	// 初始化分片数据库
	for _, dbCfg := range cfg.Shards {
		cli := d.newFunc()
		if err := cli.Init(dbCfg); err != nil {
			d.Close()
			return err
		}
		d.shards[dbCfg.ShardsId] = cli
	}
	return nil
}

func (d *RedisPool) Close() {
	for _, cli := range d.globals {
		cli.Close()
	}
	for _, cli := range d.shards {
		cli.Close()
	}
}

func (d *RedisPool) GetByName(name string) domain.IClient {
	return d.globals[name]
}

func (d *RedisPool) GetById(id uint32) domain.IClient {
	return d.shards[id]
}

func (d *RedisPool) GetByUid(uid uint64) domain.IClient {
	return d.shards[d.GetShardsId(uid)]
}

func (d *RedisPool) GetShardsId(uid uint64) uint32 {
	if d.uidModValue == 0 {
		return uint32(uid/d.shardsSize) + 1
	}
	return uint32((uid%d.uidModValue)/d.shardsSize) + 1
}
