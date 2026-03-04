package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://analytics:analytics_pass@localhost:5433/analytics_db"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
