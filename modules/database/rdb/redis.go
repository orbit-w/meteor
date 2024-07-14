package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ctx    = context.Background()
	client *Client
)

type Client struct {
	cli redis.UniversalClient
}

// RedisClient 获取包装过的redis 虚拟连接实例
func RedisClient() *Client {
	return client
}

func (c *Client) Get(key string) *redis.StringCmd {
	return c.cli.Get(ctx, key)
}

func (c *Client) Set(key string, value any, expiration time.Duration) *redis.StatusCmd {
	return c.cli.Set(ctx, key, value, expiration)
}

func (c *Client) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.cli.SetNX(ctx, key, value, expiration)
}

func (c *Client) Append(key, value string) *redis.IntCmd {
	return c.cli.Append(ctx, key, value)
}

func (c *Client) Decr(key string) *redis.IntCmd {
	return c.cli.Decr(ctx, key)
}

func (c *Client) DecrBy(key string, decrement int64) *redis.IntCmd {
	return c.cli.DecrBy(ctx, key, decrement)
}

func (c *Client) GetRange(key string, start, end int64) *redis.StringCmd {
	return c.cli.GetRange(ctx, key, start, end)
}

func (c *Client) GetSet(key string, value interface{}) *redis.StringCmd {
	return c.cli.GetSet(ctx, key, value)
}

func (c *Client) GetEx(key string, expiration time.Duration) *redis.StringCmd {
	return c.cli.GetEx(ctx, key, expiration)
}

func (c *Client) GetDel(key string) *redis.StringCmd {
	return c.cli.GetDel(ctx, key)
}

func (c *Client) Incr(key string) *redis.IntCmd {
	return c.cli.Incr(ctx, key)
}

func (c *Client) IncrBy(key string, value int64) *redis.IntCmd {
	return c.cli.IncrBy(ctx, key, value)
}

func (c *Client) IncrByFloat(key string, value float64) *redis.FloatCmd {
	return c.cli.IncrByFloat(ctx, key, value)
}

func (c *Client) LCS(q *redis.LCSQuery) *redis.LCSCmd {
	return c.cli.LCS(ctx, q)
}

func (c *Client) MGet(keys ...string) *redis.SliceCmd {
	return c.cli.MGet(ctx, keys...)
}

func (c *Client) MSet(values ...interface{}) *redis.StatusCmd {
	return c.cli.MSet(ctx, values...)
}

func (c *Client) MSetNX(values ...interface{}) *redis.BoolCmd {
	return c.cli.MSetNX(ctx, values...)
}

func (c *Client) SetArgs(key string, value interface{}, a redis.SetArgs) *redis.StatusCmd {
	return c.cli.SetArgs(ctx, key, value, a)
}

func (c *Client) SetEx(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.cli.SetEx(ctx, key, value, expiration)
}

func (c *Client) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return c.cli.SetXX(ctx, key, value, expiration)
}

func (c *Client) SetRange(key string, offset int64, value string) *redis.IntCmd {
	return c.cli.SetRange(ctx, key, offset, value)
}

func (c *Client) StrLen(key string) *redis.IntCmd {
	return c.cli.StrLen(ctx, key)
}
