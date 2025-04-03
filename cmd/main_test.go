package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRateLimiterIT(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := redis.Run(ctx, "redis:7")
	require.NoError(t, err)
	defer redisContainer.Terminate(ctx)

	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	require.NotEmpty(t, redisPort)
	redisContainerHost, _ := redisContainer.Host(ctx)
	fmt.Printf("Redis container host: %s:%s\n", redisContainerHost, redisPort.Port())

	rateLimitIpBlockDurationStr := "6s"
	rateLimitIpBlockDuration, _ := time.ParseDuration(rateLimitIpBlockDurationStr)
	rateLimitTokenBlockDurationStr := "5s"
	rateLimitTokenBlockDuration, _ := time.ParseDuration(rateLimitTokenBlockDurationStr)
	rateLimitIpMaxRequest := 5
	rateLimitTokenMaxRequest := 10

	input := ApplicationDependenciesInputDto{
		RedisHost:                   fmt.Sprintf("%s:%s", redisContainerHost, redisPort.Port()),
		RateLimitIpBlockDuration:    rateLimitIpBlockDurationStr,
		RateLimitTokenBlockDuration: rateLimitTokenBlockDurationStr,
		RateLimitIpMaxRequest:       rateLimitIpMaxRequest,
		RateLimitTokenMaxRequest:    rateLimitTokenMaxRequest,
	}

	dependencies := ApplicationDependenciesFactory(input)
	require.NotNil(t, dependencies)

	router := RouterFactory(dependencies.rateLimiterMiddleware)

	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("Given a token", func(t *testing.T) {
		t.Run("Should rate limit requests", func(t *testing.T) {
			for i := 1; i <= rateLimitTokenMaxRequest+1; i++ {
				resp := makeRequestWithToken(t, server.URL, "f2bb9a91-2c93-4613-ab54-7728792a6280")

				if i <= rateLimitTokenMaxRequest {
					assert.Equal(t, http.StatusOK, resp.StatusCode)
				} else {
					assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
				}
			}
		})

		t.Run("Should reset rate limit after block duration", func(t *testing.T) {
			token := "bab0e627-c1f7-474f-b0cf-7b78747f2f40"
			batchTokenRequest(t, server.URL, rateLimitTokenMaxRequest, token)

			resp := makeRequestWithToken(t, server.URL, token)
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

			time.Sleep(rateLimitTokenBlockDuration)

			resp = makeRequestWithToken(t, server.URL, token)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("Given an ip", func(t *testing.T) {
		t.Run("Should rate limit requests", func(t *testing.T) {
			for i := 1; i <= rateLimitIpMaxRequest+1; i++ {
				resp := makeRequestWithIp(t, server.URL, "192.168.0.1")

				if i <= rateLimitIpMaxRequest {
					assert.Equal(t, http.StatusOK, resp.StatusCode)
				} else {
					assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
				}
			}
		})

		t.Run("Should reset rate limit after block duration", func(t *testing.T) {
			ip := "192.168.0.1"
			batchIpRequest(t, server.URL, rateLimitIpMaxRequest, ip)

			resp := makeRequestWithIp(t, server.URL, ip)
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

			time.Sleep(rateLimitIpBlockDuration)

			resp = makeRequestWithIp(t, server.URL, ip)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

}

func makeRequestWithIp(t *testing.T, url string, ip string) *http.Response {
	return makeRequest(t, url, "X-Forwarded-For", ip)
}

func makeRequestWithToken(t *testing.T, url string, token string) *http.Response {
	return makeRequest(t, url, "API_KEY", token)
}

func batchIpRequest(t *testing.T, url string, batchSize int, ip string) {
	for i := 0; i < batchSize; i++ {
		makeRequestWithIp(t, url, ip)
	}
}

func batchTokenRequest(t *testing.T, url string, batchSize int, token string) {
	for i := 0; i < batchSize; i++ {
		makeRequestWithToken(t, url, token)
	}
}

func makeRequest(t *testing.T, url string, header, value string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	req.Header.Set(header, value)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}
