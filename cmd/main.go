package main

import (
	"log"
	"net/http"
	"time"

	"github.com/egon89/go-rate-limiter/internal/adapters"
	"github.com/egon89/go-rate-limiter/internal/config"
	"github.com/egon89/go-rate-limiter/internal/middlewares"
	"github.com/egon89/go-rate-limiter/internal/selectors"
	"github.com/egon89/go-rate-limiter/internal/services"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.LoadEnv()

	redisAdapter := adapters.NewRedisAdapter(config.RedisHost)
	ipBlockDuration, err := time.ParseDuration(config.RateLimitIpBlockDuration)
	if err != nil {
		log.Fatalf("parse rate limit ip block duration: %v", err)
	}
	tokenBlockDuration, err := time.ParseDuration(config.RateLimitTokenBlockDuration)
	if err != nil {
		log.Fatalf("parse rate limit token block duration: %v", err)
	}

	rateLimitIpService := services.NewRateLimiterIpService(redisAdapter, config.RateLimitIpCount, ipBlockDuration)
	rateLimitTokenService := services.NewRateLimiterTokenService(redisAdapter, config.RateLimitTokenCount, tokenBlockDuration)

	rateLimitStrategySelector := selectors.NewRateLimiterStrategySelector(
		[]services.RateLimiterService{
			rateLimitIpService,
			rateLimitTokenService,
		},
	)

	rateLimitMiddleware := middlewares.NewRateLimiterMiddleware(rateLimitStrategySelector)

	r := chi.NewRouter()
	r.Use(rateLimitMiddleware.Intercept)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	log.Println("Starting server on port " + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}
