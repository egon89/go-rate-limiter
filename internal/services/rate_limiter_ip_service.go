package services

import (
	"time"

	"github.com/egon89/go-rate-limiter/internal/ports"
)

type RateLimiterIpService struct {
	rateLimiterBaseService
}

func NewRateLimiterIpService(storage ports.RateLimiterStorage, rateLimitCount int, rateLimitBlockDuration time.Duration) *RateLimiterIpService {
	return &RateLimiterIpService{
		rateLimiterBaseService: *NewRateLimiterBaseService(storage, rateLimitCount, rateLimitBlockDuration),
	}
}

func (r *RateLimiterIpService) Select(rateType RateLimiterType) bool {
	return rateType == IP
}
