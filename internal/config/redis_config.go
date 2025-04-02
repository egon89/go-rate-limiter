package config

import (
	"github.com/redis/go-redis/v9"
)

func RedisClient(host string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: host,
	})
}
