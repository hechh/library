package redispool

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type IClient interface {
	Init(cfg *DbConfig) error
	Close() error
	GetRealKey(key string) string
	Ctx() context.Context
	Pipeline() redis.Pipeliner
	TxPipeline() redis.Pipeliner
	Run(script *redis.Script, key string, values ...any) (any, error)
	ClusterNodes() (string, error)
	Ping() (string, error)
	Publish(channel string, message any) (int64, error)
	Del(keys ...string) (int64, error)
	Exists(key string) (int64, error)
	Expire(key string, expiration time.Duration) (bool, error)
	TTL(key string) (time.Duration, error)
	Get(key string) (string, error)
	Set(key string, val any, expiration time.Duration) error
	SetNX(key string, val any, expiration time.Duration) (bool, error)
	SetEX(key string, val any, expiration time.Duration) error
	Incr(key string) (int64, error)
	IncrBy(key string, val int64) (int64, error)
	Decr(key string) (int64, error)
	DecrBy(key string, value int64) (int64, error)
	MGet(keys ...string) ([]any, error)
	MSet(args ...any) error
	SAdd(key string, members ...any) (int64, error)
	SRem(key string, members ...any) (int64, error)
	SMembers(key string) ([]string, error)
	SIsMember(key string, member any) (bool, error)
	SCard(key string) (int64, error)
	SRandMemberN(key string, count int64) ([]string, error)
	ZAdd(key string, members ...*redis.Z) (int64, error)
	ZRemRangeByRank(key string, start, stop int64) (int64, error)
	ZRem(key string, members ...any) (int64, error)
	ZCard(key string) (int64, error)
	ZScore(key, member string) (float64, error)
	ZRange(key string, start, stop int64) ([]string, error)
	ZRevRange(key string, start, stop int64) ([]string, error)
	ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
	ZRevRangeByScore(key string, opt *redis.ZRangeBy) ([]string, error)
	ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) ([]redis.Z, error)
	ZRank(key, member string) (int64, error)
	ZRevRank(key, member string) (int64, error)
	LPush(key string, values ...any) (int64, error)
	RPush(key string, values ...any) (int64, error)
	LPop(key string) (string, error)
	RPop(key string) (string, error)
	LLen(key string) (int64, error)
	LRange(key string, start, stop int64) ([]string, error)
	LTrim(key string, start, stop int64) error
	LRem(key string, count int64, value any) (int64, error)
	HGet(key string, field string) (string, error)
	HSet(key string, field string, val any) error
	HMGet(key string, fields ...string) ([]any, error)
	HMSet(key string, vals ...any) error
	HGetAll(key string) (map[string]string, error)
	HDel(key string, fields ...string) (int64, error)
	HExists(key, field string) (bool, error)
	HIncrBy(key, field string, incr int64) (int64, error)
	HKeys(key string) ([]string, error)
	HLen(key string) (int64, error)
	HSetNX(key, field string, value any) (bool, error)
}

// DbConfig 数据库分片配置
type DbConfig struct {
	ShardsId uint32              `yaml:"shards_id,omitempty"` // 分片ID
	DbName   string              `yaml:"dbname,omitempty"`    // 数据库名称
	Db       uint32              `yaml:"db,omitempty"`        // 数据库编号
	User     string              `yaml:"user,omitempty"`      // 数据库用户名
	Password string              `yaml:"password,omitempty"`  // 数据库密码
	Ip       string              `yaml:"ip,omitempty"`        // 数据库IP地址
	Port     uint32              `yaml:"port,omitempty"`      // 数据库端口
	Prefix   string              `yaml:"prefix,omitempty"`    // key 前缀
	Slaves   map[int32]*DbConfig `yaml:"slaves,omitempty"`    // 从库配置列表
}

type Config struct {
	UidModValue uint64      `yaml:"uid_mod_value,omitempty"` // 用户ID取模值
	ShardsSize  int64       `yaml:"shards_size,omitempty"`   // 分片数量
	Globals     []*DbConfig `yaml:"globals,omitempty"`       // 全局 Redis 配置列表
	Shards      []*DbConfig `yaml:"shards,omitempty"`        // 分片 Redis 配置列表
}

type RedisPool struct {
	uidModValue uint64             // 开始 Uid 值，用于路由
	shardsSize  uint64             // 分片数量
	newFunc     func() IClient     // new函数
	globals     map[string]IClient // 全局数据库连接池
	shards      map[uint32]IClient // 分片数据库连接池
}

func NewRedisPool[T IClient](f func() T) *RedisPool {
	return &RedisPool{
		newFunc: func() IClient { return f() },
		globals: make(map[string]IClient),
		shards:  make(map[uint32]IClient),
	}
}

func (d *RedisPool) Init(cfg *Config) error {
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

func (d *RedisPool) GetByName(name string) IClient {
	return d.globals[name]
}

func (d *RedisPool) GetById(id uint32) IClient {
	return d.shards[id]
}

func (d *RedisPool) GetByUid(uid uint64) IClient {
	return d.shards[d.GetShardsId(uid)]
}

func (d *RedisPool) GetShardsId(uid uint64) uint32 {
	if d.uidModValue == 0 {
		return uint32(uid/d.shardsSize) + 1
	}
	return uint32((uid%d.uidModValue)/d.shardsSize) + 1
}
