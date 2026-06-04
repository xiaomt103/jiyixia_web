package config

import (
	"os"
	"time"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RedisURL      string
	AdminUsername string
	AdminPassword string
	SessionTTL    time.Duration
}

func Load() Config {
	return Config{
		Port:          getenv("PORT", "8080"),
		DatabaseURL:   getenv("DATABASE_URL", "postgres://members:members@localhost:5432/members?sslmode=disable"),
		RedisURL:      getenv("REDIS_URL", "redis://localhost:6379/0"),
		AdminUsername: getenv("ADMIN_USERNAME", "admin"),
		AdminPassword: getenv("ADMIN_PASSWORD", "admin123"),
		SessionTTL:    12 * time.Hour,
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
