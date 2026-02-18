package config

import (
	"os"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			// URL: envOr("DB_CONN", "postgres://root:12312@localhost:5432/test_db"),
			URL: envOr("DB_CONN", "postgres://postgres:postgres@localhost:5432/postgres"),
		},
		Server: ServerConfig{
			Port: envOr("PORT", ":8000"),
		},
	}
}

func envOr(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}
