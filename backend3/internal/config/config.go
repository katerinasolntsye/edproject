package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Path string
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret            string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}

func Load() *Config {
	accessExp := parseDuration(envOr("JWT_ACCESS_EXPIRATION", "15"), time.Minute)
	refreshExp := parseDuration(envOr("JWT_REFRESH_EXPIRATION", "7"), 24*time.Hour)

	return &Config{
		Database: DatabaseConfig{
			Path: envOr("DB_PATH", "./data.db"),
		},
		Server: ServerConfig{
			Port: envOr("PORT", ":8000"),
		},
		JWT: JWTConfig{
			Secret:            envOr("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessExpiration:  accessExp,
			RefreshExpiration: refreshExp,
		},
	}
}

func envOr(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(value string, unit time.Duration) time.Duration {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 15 * time.Minute // default fallback
	}
	return time.Duration(i) * unit
}
