package config

import (
	"os"
	"strconv"
	"time"
)

// Config 应用配置结构
type Config struct {
	DatabaseURL       string
	SourceURL         string
	SyncEnabled      bool
	ScrapeInterval    time.Duration
	MaxRetries        int
	RequestTimeout    time.Duration
	UserAgent         string
	RateLimitDelay    time.Duration
	// 隐私发帖配置
	ProxyEnabled      bool
	ProxyURL          string
}

// Load 加载配置
func Load() *Config {
	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "data.db"),
		SourceURL:         getEnv("SOURCE_URL", "https://example-tree-hole.com"),
		SyncEnabled:      getEnv("SYNC_ENABLED", "false") == "true",
		MaxRetries:        getIntEnv("MAX_RETRIES", 3),
		RequestTimeout:    getDurationEnv("REQUEST_TIMEOUT", 30*time.Second),
		UserAgent:         getEnv("USER_AGENT", "TreeHoleMirror/1.0"),
		RateLimitDelay:    getDurationEnv("RATE_LIMIT_DELAY", 1*time.Second),
		// 隐私发帖配置
		ProxyEnabled:      getEnv("PROXY_ENABLED", "false") == "true",
		ProxyURL:          getEnv("PROXY_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
