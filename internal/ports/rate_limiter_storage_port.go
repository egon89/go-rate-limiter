package ports

import (
	"context"
	"time"
)

type RateLimiterStorage interface {
	GetRequestCount(ctx context.Context, key string) (int, error)
	Increment(ctx context.Context, key string) (int, error)
	Expire(ctx context.Context, key string, window time.Duration) error
}
