package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host             string
	Port             int
	DatabaseURL      string
	RedisURL         string
	JWTSecret        string
	JWTIssuer        string
	JWTAudience      string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	ApprovalMode     bool
	LogLevel         string
	AllowedOriginURL string
}

func Load() (Config, error) {
	port, err := intFromEnv("GATEWAY_PORT", 8080)
	if err != nil {
		return Config{}, err
	}

	accessMinutes, err := intFromEnv("JWT_ACCESS_TTL_MINUTES", 15)
	if err != nil {
		return Config{}, err
	}

	refreshHours, err := intFromEnv("JWT_REFRESH_TTL_HOURS", 168)
	if err != nil {
		return Config{}, err
	}

	secret := getEnv("JWT_SECRET", "")
	if len(secret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	cfg := Config{
		Host:             getEnv("GATEWAY_HOST", "0.0.0.0"),
		Port:             port,
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://nexus:nexus@localhost:5432/nexus_mcp?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:        secret,
		JWTIssuer:        getEnv("JWT_ISSUER", "nexus-mcp"),
		JWTAudience:      getEnv("JWT_AUDIENCE", "nexus-mcp-clients"),
		AccessTokenTTL:   time.Duration(accessMinutes) * time.Minute,
		RefreshTokenTTL:  time.Duration(refreshHours) * time.Hour,
		ApprovalMode:     boolFromEnv("APPROVAL_MODE", true),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		AllowedOriginURL: getEnv("ALLOWED_ORIGIN", "http://localhost:3000"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func intFromEnv(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer for %s: %w", key, err)
	}
	return parsed, nil
}

func boolFromEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
