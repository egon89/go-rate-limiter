package selectors

import (
	"fmt"
	"log"

	"github.com/egon89/go-rate-limiter/internal/services"
)

type RateLimiterStrategySelector struct {
	Services []services.RateLimiterService
}

func NewRateLimiterStrategySelector(services []services.RateLimiterService) *RateLimiterStrategySelector {
	return &RateLimiterStrategySelector{
		Services: services,
	}
}

func (r *RateLimiterStrategySelector) Select(rateType services.RateLimiterType) (services.RateLimiterService, error) {
	for _, service := range r.Services {
		if service.Select(rateType) {
			log.Printf("[rate-limiter-strategy-selector] rate limiter service selected by type: %d\n", rateType)

			return service, nil
		}
	}
	return nil, fmt.Errorf("rate limiter service not found")
}
