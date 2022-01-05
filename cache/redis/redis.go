package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisCluster struct {
	Client *redis.ClusterClient
	Ctx    context.Context
}

func NewRedisCluster(addrs []string, password string) *RedisCluster {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		Password:     password,
		PoolSize:     10,
		MinIdleConns: 5,
		IdleTimeout:  -1,
	})

	return &RedisCluster{
		Client: client,
		Ctx:    context.Background(),
	}
}

func IsRedisErrKeyNotExist(err error) bool {
	return err == redis.Nil
}
