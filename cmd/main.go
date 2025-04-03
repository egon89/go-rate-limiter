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

type ApplicationDependenciesInputDto struct {
	RedisHost                   string
	RateLimitIpBlockDuration    string
	RateLimitTokenBlockDuration string
	RateLimitIpCount            int
	RateLimitTokenCount         int
}

type ApplicationDependenciesOutputDto struct {
	rateLimiterMiddleware *middlewares.RateLimiterMiddleware
}

func main() {
	config.LoadEnv()

	input := ApplicationDependenciesInputDto{
		RedisHost:                   config.RedisHost,
		RateLimitIpBlockDuration:    config.RateLimitIpBlockDuration,
		RateLimitTokenBlockDuration: config.RateLimitTokenBlockDuration,
		RateLimitIpCount:            config.RateLimitIpCount,
		RateLimitTokenCount:         config.RateLimitTokenCount,
	}
	dependencies := ApplicationDependenciesFactory(input)

	router := RouterFactory(dependencies.rateLimiterMiddleware)

	log.Println("Starting server on port " + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

func RouterFactory(rateLimitMiddleware *middlewares.RateLimiterMiddleware) *chi.Mux {
	r := chi.NewRouter()
	r.Use(rateLimitMiddleware.Intercept)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	return r
}

func ApplicationDependenciesFactory(input ApplicationDependenciesInputDto) *ApplicationDependenciesOutputDto {
	redisAdapter := adapters.NewRedisAdapter(input.RedisHost)
	ipBlockDuration, err := time.ParseDuration(input.RateLimitIpBlockDuration)
	if err != nil {
		log.Fatalf("parse rate limit ip block duration: %v", err)
	}
	tokenBlockDuration, err := time.ParseDuration(input.RateLimitTokenBlockDuration)
	if err != nil {
		log.Fatalf("parse rate limit token block duration: %v", err)
	}

	rateLimitIpService := services.NewRateLimiterIpService(redisAdapter, input.RateLimitIpCount, ipBlockDuration)
	rateLimitTokenService := services.NewRateLimiterTokenService(redisAdapter, input.RateLimitTokenCount, tokenBlockDuration)

	rateLimitStrategySelector := selectors.NewRateLimiterStrategySelector(
		[]services.RateLimiterService{
			rateLimitIpService,
			rateLimitTokenService,
		},
	)

	rateLimitMiddleware := middlewares.NewRateLimiterMiddleware(rateLimitStrategySelector)

	return &ApplicationDependenciesOutputDto{
		rateLimiterMiddleware: rateLimitMiddleware,
	}
}
