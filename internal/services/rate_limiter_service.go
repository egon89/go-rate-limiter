package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/egon89/go-rate-limiter/internal/ports"
)

type RateLimiterType int

const (
	IP RateLimiterType = iota
	TOKEN
)

type AllowRequestInput struct {
	Value string
	Type  RateLimiterType
}

type RateLimiterService interface {
	Allow(ctx context.Context, key string) (bool, error)
	Select(rateType RateLimiterType) bool
}

type rateLimiterBaseService struct {
	storage                ports.RateLimiterStorage
	rateLimitMaxRequest    int
	rateLimitBlockDuration time.Duration
}

func NewRateLimiterBaseService(storage ports.RateLimiterStorage, rateLimitMaxRequest int, rateLimitBlockDuration time.Duration) *rateLimiterBaseService {
	return &rateLimiterBaseService{
		storage:                storage,
		rateLimitMaxRequest:    rateLimitMaxRequest,
		rateLimitBlockDuration: rateLimitBlockDuration,
	}
}

func (r *rateLimiterBaseService) Allow(ctx context.Context, key string) (bool, error) {
	log.Printf("[rate-limit-service] limit count: %d, block duration: %v\n", r.rateLimitMaxRequest, r.rateLimitBlockDuration)
	log.Printf("[rate-limit-service] key: %s\n", key)

	count, err := r.storage.Increment(ctx, key)
	if err != nil {
		return false, fmt.Errorf("increment request count error for key %s: %v", key, err)
	}

	log.Printf("[rate-limit-service] the key %s has been requested %d time(s)\n", key, count)

	if count == 1 {
		log.Printf("[rate-limit-service] setting expiration for key %s: %v\n", key, r.rateLimitBlockDuration)

		err = r.storage.Expire(ctx, key, r.rateLimitBlockDuration)
		if err != nil {
			return false, fmt.Errorf("set expiration for key %s: %v", key, err)
		}
	}

	if count > r.rateLimitMaxRequest {
		log.Printf("[rate-limit-service] the key %s has reached the limit count\n", key)

		return false, nil
	}

	return true, nil
}
