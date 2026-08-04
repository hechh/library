package redispool

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hechh/library/base/templ"
	"google.golang.org/protobuf/proto"
)

/*
const (
	HASH_FLAG      = 1 << 0 // hash数据类型
	STRING_FLAG    = 1 << 1 // string数据类型
	GLOBAL_FLAG    = 1 << 2 // 全局数据库
	SHARDS_FLAG    = 1 << 3 // 分片数据库
	TEMP_FLAG      = 1 << 4 // 临时数据
	PERMANENT_FLAG = 1 << 5 // 常驻数据
)
*/

type IClient interface {
	UniqueId() uint32
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

type Message interface {
	CloneMessageVT() proto.Message
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

type Value struct {
	IClient
	Message
	key   string
	field string
	times uint32
}

func (d *Value) Key() string     { return d.key }
func (d *Value) Field() string   { return d.field }
func (d *Value) IsChanged() bool { return d.times > 0 }
func (d *Value) Change()         { d.times++ }
func (d *Value) Reset()          { d.times = 0 }
func (d *Value) Get() any        { return d.Message }

func (d *Value) Clone() *Value {
	return &Value{
		IClient: d.IClient,
		Message: d.CloneMessageVT().(Message),
		key:     d.key,
		field:   d.field,
	}
}

func NewValue(cli IClient, obj Message, args ...string) *Value {
	return &Value{
		IClient: cli,
		Message: obj,
		key:     templ.Index(args, 0, ""),
		field:   templ.Index(args, 1, ""),
	}
}
