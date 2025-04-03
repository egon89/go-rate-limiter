package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Port                        string
	RateLimitIpMaxRequest       int
	RateLimitIpBlockDuration    string
	RateLimitTokenMaxRequest    int
	RateLimitTokenBlockDuration string
	RedisHost                   string
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	Port = GetEnv("PORT", "8080")
	RateLimitIpMaxRequest = GetEnvAsInt("RATE_LIMIT_IP_MAX_REQUEST", 5)
	RateLimitIpBlockDuration = GetEnv("RATE_LIMIT_IP_BLOCK_DURATION", "300s")
	RateLimitTokenMaxRequest = GetEnvAsInt("RATE_LIMIT_TOKEN_MAX_REQUEST", 10)
	RateLimitTokenBlockDuration = GetEnv("RATE_LIMIT_TOKEN_BLOCK_DURATION", "300s")
	RedisHost = GetEnv("REDIS_HOST", "")
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return intValue
	}
	return fallback
}
