// Package config provides configuration management for sub2api.
// It handles loading and validation of application settings from
// environment variables with sensible defaults.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration values.
type Config struct {
	// Server settings
	Host string
	Port int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Subscription settings
	SubURL        string
	RefreshInterval time.Duration
	UserAgent     string

	// Cache settings
	CacheEnabled bool
	CacheTTL     time.Duration

	// Logging
	LogLevel string
}

// Load reads configuration from environment variables and returns a Config.
// Missing required values will result in an error.
func Load() (*Config, error) {
	port, err := getEnvInt("PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	// I prefer a longer refresh interval to reduce outbound requests
	refreshInterval, err := getEnvDuration("REFRESH_INTERVAL", 60*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_INTERVAL: %w", err)
	}

	cacheTTL, err := getEnvDuration("CACHE_TTL", 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_TTL: %w", err)
	}

	cfg := &Config{
		Host:            getEnv("HOST", "0.0.0.0"),
		Port:            port,
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		SubURL:          getEnv("SUB_URL", ""),
		RefreshInterval: refreshInterval,
		UserAgent:       getEnv("USER_AGENT", "sub2api/1.0"),
		CacheEnabled:    getEnvBool("CACHE_ENABLED", true),
		CacheTTL:        cacheTTL,
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that required configuration values are present and valid.
func (c *Config) validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}
	return nil
}

// Addr returns the full address string for the HTTP server.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// getEnv returns the value of the environment variable named by key,
// or defaultVal if the variable is not set.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvInt returns the integer value of the named environment variable,
// or defaultVal if the variable is not set.
func getEnvInt(key string, defaultVal int) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(val)
}

// getEnvBool returns the boolean value of the named environment variable,
// or defaultVal if the variable is not set.
func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}

// getEnvDuration returns the duration value of the named environment variable,
// or defaultVal if the variable is not set.
func getEnvDuration(key string, defaultVal time.Duration) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	return time.ParseDuration(val)
}
