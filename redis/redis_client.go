package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	*redis.Client
	prefix string
}

func NewClient(cli *redis.Client, prefix string) *Client {
	return &Client{Client: cli, prefix: prefix}
}

func (d *Client) GetRealKey(key string) string {
	if len(d.prefix) > 0 {
		return key
	}
	return fmt.Sprintf("%s_%s", d.prefix, key)
}

func (d *Client) Del(keys ...string) (int64, error) {
	for i, key := range keys {
		keys[i] = d.GetRealKey(key)
	}
	flag, err := d.Client.Del(context.Background(), keys...).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) Exists(key string) (int64, error) {
	flag, err := d.Client.Exists(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) Expire(key string, expiration time.Duration) (bool, error) {
	flag, err := d.Client.Expire(context.Background(), d.GetRealKey(key), expiration).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) TTL(key string) (time.Duration, error) {
	ttl, err := d.Client.TTL(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return ttl, err
}

// --------------------------------string操作--------------------------------
func (d *Client) Get(key string) (str string, err error) {
	str, err = d.Client.Get(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *Client) Set(key string, val any, expiration time.Duration) (err error) {
	_, err = d.Client.Set(context.Background(), d.GetRealKey(key), val, expiration).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// 不存在key时，设置该key的值未val
func (d *Client) SetNX(key string, val any, expiration time.Duration) (exist bool, err error) {
	exist, err = d.Client.SetNX(context.Background(), d.GetRealKey(key), val, expiration).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// 设置key的过期时间
func (d *Client) SetEX(key string, val any, expiration time.Duration) (err error) {
	_, err = d.Client.SetEX(context.Background(), d.GetRealKey(key), val, expiration).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *Client) Incr(key string) (ret int64, err error) {
	ret, err = d.Client.Incr(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *Client) IncrBy(key string, val int64) (ret int64, err error) {
	ret, err = d.Client.IncrBy(context.Background(), d.GetRealKey(key), val).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *Client) Decr(key string) (int64, error) {
	flag, err := d.Client.Decr(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) DecrBy(key string, value int64) (int64, error) {
	flag, err := d.Client.DecrBy(context.Background(), d.GetRealKey(key), value).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

// 批量读取操作
func (d *Client) MGet(keys ...string) (rets []any, err error) {
	args := []string{}
	for i := 0; i < len(keys); i++ {
		args = append(args, d.GetRealKey(keys[i]))
	}
	rets, err = d.Client.MGet(context.Background(), args...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// 批量设置键值对
func (d *Client) MSet(args ...any) (err error) {
	_, err = d.Client.MSet(context.Background(), args...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// --------------------------------hash操作--------------------------------
func (d *Client) HGet(key string, field string) (ret string, err error) {
	ret, err = d.Client.HGet(context.Background(), d.GetRealKey(key), field).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *Client) HSet(key string, field string, val interface{}) (err error) {
	_, err = d.Client.HSet(context.Background(), d.GetRealKey(key), field, val).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *Client) HMGet(key string, fields ...string) ([]interface{}, error) {
	rets, err := d.Client.HMGet(context.Background(), key, fields...).Result()
	if err == redis.Nil {
		err = nil
	}
	return rets, err
}

func (d *Client) HMSet(key string, vals ...interface{}) (err error) {
	_, err = d.Client.HMSet(context.Background(), d.GetRealKey(key), vals...).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *Client) HGetAll(key string) (ret map[string]string, err error) {
	ret, err = d.Client.HGetAll(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

func (d *Client) HDel(key string, fields ...string) (int64, error) {
	flag, err := d.Client.HDel(context.Background(), d.GetRealKey(key), fields...).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) HExists(key, field string) (bool, error) {
	flag, err := d.Client.HExists(context.Background(), d.GetRealKey(key), field).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) HIncrBy(key, field string, incr int64) (int64, error) {
	flag, err := d.Client.HIncrBy(context.Background(), d.GetRealKey(key), field, incr).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) HKeys(key string) ([]string, error) {
	rets, err := d.Client.HKeys(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return rets, err
}

func (d *Client) HLen(key string) (int64, error) {
	flag, err := d.Client.HLen(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) HSetNX(key, field string, value interface{}) (bool, error) {
	flag, err := d.Client.HSetNX(context.Background(), d.GetRealKey(key), field, value).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

// --------------------------------List操作--------------------------------
func (d *Client) LPush(key string, values ...interface{}) (int64, error) {
	flag, err := d.Client.LPush(context.Background(), d.GetRealKey(key), values...).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) RPush(key string, values ...interface{}) (int64, error) {
	flag, err := d.Client.RPush(context.Background(), d.GetRealKey(key), values...).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) LPop(key string) (string, error) {
	flag, err := d.Client.LPop(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) RPop(key string) (string, error) {
	flag, err := d.Client.RPop(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) LLen(key string) (int64, error) {
	flag, err := d.Client.LLen(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) LRange(key string, start, stop int64) ([]string, error) {
	flag, err := d.Client.LRange(context.Background(), d.GetRealKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) LTrim(key string, start, stop int64) (string, error) {
	flag, err := d.Client.LTrim(context.Background(), d.GetRealKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) LRem(key string, count int64, value interface{}) (int64, error) {
	flag, err := d.Client.LRem(context.Background(), d.GetRealKey(key), count, value).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

// --------------------------------set操作--------------------------------
func (d *Client) SAdd(key string, members ...interface{}) (int64, error) {
	flag, err := d.Client.SAdd(context.Background(), d.GetRealKey(key), members...).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) SRem(key string, members ...interface{}) (int64, error) {
	flag, err := d.Client.SRem(context.Background(), d.GetRealKey(key), members...).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) SMembers(key string) ([]string, error) {
	flag, err := d.Client.SMembers(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) SIsMember(key string, member interface{}) (bool, error) {
	flag, err := d.Client.SIsMember(context.Background(), d.GetRealKey(key), member).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) SCard(key string) (int64, error) {
	flag, err := d.Client.SCard(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

// --------------------------------sorted set操作--------------------------------
func (d *Client) ZAdd(key string, members ...*redis.Z) (int64, error) {
	flag, err := d.Client.ZAdd(context.Background(), d.GetRealKey(key), members...).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) ZRem(key string, members ...interface{}) (int64, error) {
	flag, err := d.Client.ZRem(context.Background(), d.GetRealKey(key), members...).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) ZCard(key string) (int64, error) {
	flag, err := d.Client.ZCard(context.Background(), d.GetRealKey(key)).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) ZScore(key, member string) (float64, error) {
	flag, err := d.Client.ZScore(context.Background(), d.GetRealKey(key), member).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) ZRange(key string, start, stop int64) ([]string, error) {
	flag, err := d.Client.ZRange(context.Background(), d.GetRealKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) ZRevRange(key string, start, stop int64) ([]string, error) {
	flag, err := d.Client.ZRevRange(context.Background(), d.GetRealKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	flag, err := d.Client.ZRangeWithScores(context.Background(), d.GetRealKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	flag, err := d.Client.ZRevRangeWithScores(context.Background(), d.GetRealKey(key), start, stop).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) ZRank(key, member string) (int64, error) {
	flag, err := d.Client.ZRank(context.Background(), d.GetRealKey(key), member).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) ZRevRank(key, member string) (int64, error) {
	flag, err := d.Client.ZRevRank(context.Background(), d.GetRealKey(key), member).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) Run(script *redis.Script, key string, values ...any) (any, error) {
	return script.Run(context.Background(), d.Client, []string{d.GetRealKey(key)}, values...).Result()
}

// 发布订阅
func (d *Client) Publish(channel string, message interface{}) (int64, error) {
	flag, err := d.Client.Publish(context.Background(), d.GetRealKey(channel), message).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

// Pipeline支持
func (d *Client) Pipeline() redis.Pipeliner {
	return d.Client.Pipeline()
}

// 事务支持
func (d *Client) TxPipeline() redis.Pipeliner {
	return d.Client.TxPipeline()
}

// 集群支持
func (d *Client) ClusterNodes() (string, error) {
	flag, err := d.Client.ClusterNodes(context.Background()).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

// 连接管理
func (d *Client) Ping() (string, error) {
	flag, err := d.Client.Ping(context.Background()).Result()
	if err == redis.Nil {
		err = nil
	}
	return flag, err
}

func (d *Client) Close() error {
	err := d.Client.Close()
	if err == redis.Nil {
		err = nil
	}
	return err
}
