package rdb

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	universalClient redis.UniversalClient

	ErrAddressInvalid = errors.New("err_redis_address_invalid")
)

// UniversalClient 获取原始的redis 虚拟连接实例
func UniversalClient() redis.UniversalClient {
	return universalClient
}

type RedisClientOps struct {
	Addr     []string
	Cluster  bool
	Username string
	Password string
	DB       int
}

func NewClient(ops RedisClientOps) error {
	if len(ops.Addr) == 0 {
		return ErrAddressInvalid
	}

	switch {
	case ops.Cluster:
		universalClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          ops.Addr,
			Username:       ops.Username,
			Password:       ops.Password,
			MaxIdleConns:   20,
			MaxActiveConns: 50,
		})
	default:
		universalClient = redis.NewClient(&redis.Options{
			Addr:           ops.Addr[0],
			Username:       ops.Username,
			Password:       ops.Password, // no password set
			DB:             ops.DB,       // use default db
			MaxIdleConns:   20,
			MaxActiveConns: 50,
		})
	}

	if err := universalClient.Ping(context.Background()).Err(); err != nil {
		return err
	}

	client = &Client{
		cli: universalClient,
	}
	return nil
}
