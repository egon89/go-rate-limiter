package adapters

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/egon89/go-rate-limiter/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(host string) *RedisAdapter {
	return &RedisAdapter{
		client: config.RedisClient(host),
	}
}

func (r *RedisAdapter) getRateLimitKey(key string) string {
	return fmt.Sprintf("rate_limit:%s", key)
}

func (r *RedisAdapter) GetRequestCount(ctx context.Context, key string) (int, error) {
	limitKey := r.getRateLimitKey(key)

	count, err := r.client.Get(ctx, limitKey).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf("get request count: %v", err)
	}

	return count, nil
}

func (r *RedisAdapter) Increment(ctx context.Context, key string) (int, error) {
	count, err := r.client.Incr(ctx, r.getRateLimitKey(key)).Result()
	if err != nil {
		log.Printf("increment request count error for key %s: %v\n", key, err)

		return 0, err
	}

	return int(count), nil
}

func (r *RedisAdapter) Expire(ctx context.Context, key string, window time.Duration) error {
	log.Printf("setting expiration for key %s: %v\n", key, window)

	err := r.client.Expire(ctx, r.getRateLimitKey(key), window).Err()
	if err != nil {
		return fmt.Errorf("error setting expiration for key %s: %v", key, err)
	}

	return nil
}
