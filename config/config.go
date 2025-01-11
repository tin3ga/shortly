package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/tin3ga/shortly/utils"
)

// Config func to get env value
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Error loading .env file")
	}
	return os.Getenv(key)
}

// Config holds the application configuration.
type ConfigParams struct {
	Port                   string
	DatabaseURL            string
	URL                    string
	RedisAddr              string
	RedisPassword          string
	RedisDB                string
	EnableCaching          bool
	CacheTTL               time.Duration
	EnableRateLimiting     bool
	MaxConnectionsLimit    int
	RateLimitExpiration    time.Duration
	APIKey                 string
	SkipFailedRequests     bool
	SkipSuccessfulRequests bool
	Title                  string
	FontURL                string
	JWTSecret              string
}

func InitializeConfig() *ConfigParams {
	enableCaching, _ := strconv.ParseBool(Config("caching_enabled"))
	ttl, _ := utils.ConvertStr(Config("cache_TTL"))
	enableRateLimiting, _ := strconv.ParseBool(Config("rate_limiting_enabled"))
	maxConnections, _ := utils.ConvertStr(Config("max_connections_limit"))
	expiration, _ := utils.ConvertStr(Config("expiration"))
	skip_failed_requests, _ := strconv.ParseBool(Config("skip_failed_requests"))
	skip_successful_requests, _ := strconv.ParseBool(Config("skip_successful_requests"))

	return &ConfigParams{
		Port:                   Config("PORT"),
		DatabaseURL:            Config("DATABASE_URL"),
		URL:                    Config("URL"),
		RedisAddr:              Config("REDIS_ADDR"),
		RedisPassword:          Config("REDIS_PASSWORD"),
		RedisDB:                Config("REDIS_DB"),
		EnableCaching:          enableCaching,
		CacheTTL:               time.Duration(ttl) * time.Minute,
		EnableRateLimiting:     enableRateLimiting,
		MaxConnectionsLimit:    maxConnections,
		RateLimitExpiration:    time.Duration(expiration) * time.Minute,
		APIKey:                 Config("apiKey"),
		SkipFailedRequests:     skip_failed_requests,
		SkipSuccessfulRequests: skip_successful_requests,
		Title:                  Config("metrics_title"),
		FontURL:                Config("metrics_font_URL"),
		JWTSecret:              Config("jwt_secret"),
	}
}
