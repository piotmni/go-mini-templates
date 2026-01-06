package config

import (
	"os"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	URL string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DB_URL", "postgres://pastebin:pastebin@localhost:5432/pastebin?sslmode=disable"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
