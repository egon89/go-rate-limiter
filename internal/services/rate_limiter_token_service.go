package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/egon89/go-rate-limiter/internal/ports"
)

type RateLimiterTokenService struct {
	rateLimiterBaseService
}

func NewRateLimiterTokenService(storage ports.RateLimiterStorage, rateLimitMaxRequest int, rateLimitBlockDuration time.Duration) *RateLimiterTokenService {
	return &RateLimiterTokenService{
		rateLimiterBaseService: *NewRateLimiterBaseService(storage, rateLimitMaxRequest, rateLimitBlockDuration),
	}
}

// override the Allow method for a custom behavior
func (r *RateLimiterTokenService) Allow(ctx context.Context, key string) (bool, error) {
	rateLimitBlockDuration := r.getRateLimitBlockDuration(key)

	log.Printf("[rate-limit-token-service] limit count: %d, block duration: %v\n", r.rateLimitMaxRequest, rateLimitBlockDuration)
	log.Printf("[rate-limit-token-service] key: %s\n", key)

	count, err := r.storage.Increment(ctx, key)
	if err != nil {
		return false, fmt.Errorf("increment request count error for key %s: %v", key, err)
	}

	log.Printf("[rate-limit-token-service] the key %s has been requested %d time(s)\n", key, count)

	if count == 1 {
		log.Printf("[rate-limit-token-service] setting expiration for key %s: %v\n", key, rateLimitBlockDuration)

		err = r.storage.Expire(ctx, key, rateLimitBlockDuration)
		if err != nil {
			return false, fmt.Errorf("set expiration for key %s: %v", key, err)
		}
	}

	if count > r.rateLimitMaxRequest {
		log.Printf("[rate-limit-token-service] the key %s has reached the limit count\n", key)

		return false, nil
	}

	return true, nil
}

func (r *RateLimiterTokenService) Select(rateType RateLimiterType) bool {
	return rateType == TOKEN
}

func (r *RateLimiterTokenService) getRateLimitBlockDuration(key string) time.Duration {
	if duration, ok := customTokenBlockDuration()[key]; ok {
		log.Printf("[rate-limit-token-service] custom block duration for key %s: %v\n", key, duration)
		return duration
	}

	return r.rateLimitBlockDuration
}

func customTokenBlockDuration() map[string]time.Duration {
	return map[string]time.Duration{
		"2c02b5ce-04d0-4c75-9810-c3e75c397956": 10 * time.Second,
		"a6b3fdef-c107-4970-8ecc-94817ed5968c": 30 * time.Second,
		"16a661d8-ce97-44b3-a405-a1400d705de8": 2 * time.Minute,
	}
}
