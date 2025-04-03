package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/egon89/go-rate-limiter/internal/selectors"
	"github.com/egon89/go-rate-limiter/internal/services"
	"github.com/egon89/go-rate-limiter/internal/utils"
)

type RateLimiterMiddleware struct {
	rateLimiterServiceStrategySelector *selectors.RateLimiterStrategySelector
}

func NewRateLimiterMiddleware(rateLimiterServiceStrategySelector *selectors.RateLimiterStrategySelector) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		rateLimiterServiceStrategySelector: rateLimiterServiceStrategySelector,
	}
}

func (rlm *RateLimiterMiddleware) Intercept(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("[middleware] rate limiter middleware")

		value, rateLimiterService, err := rlm.getStrategyValues(r)
		if err != nil {
			http.Error(rw, "No rate limiter service found", http.StatusInternalServerError)
			return
		}

		allowed, err := rateLimiterService.Allow(r.Context(), value)
		if err != nil {
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			rw.WriteHeader(http.StatusTooManyRequests)
			rw.Header().Set("Content-Type", "application/json")
			rw.Write([]byte(`{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"}`))
			return
		}

		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(fn)
}

func (r *RateLimiterMiddleware) getStrategyValues(request *http.Request) (string, services.RateLimiterService, error) {
	token := strings.TrimSpace(request.Header.Get("API_KEY"))
	if token != "" {
		log.Printf("[middleware] asking for token rate limiter service: %s\n", token)
		service, err := r.rateLimiterServiceStrategySelector.Select(services.TOKEN)

		return token, service, err
	}

	ip, _ := utils.GetIpAddress(request)
	if ip != "" {
		log.Printf("[middleware] asking for ip rate limiter service: %s\n", ip)
		service, err := r.rateLimiterServiceStrategySelector.Select(services.IP)

		return ip, service, err
	}

	return "", nil, fmt.Errorf("rate limiter service not found")
}
