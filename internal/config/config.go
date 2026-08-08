// Package config loads runtime configuration from a .env file and/or the
// process environment.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration.
type Config struct {
	// DatabaseURL is the full Postgres DSN used by pgx and goose.
	// When empty, it is composed from the DB_* fields below.
	DatabaseURL string

	AppEnv   string
	HTTPAddr string

	// Individual connection fields (only used when DatabaseURL is empty).
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load reads the optional .env file at path (falling back to ".env") and the
// process environment, then returns a populated Config.
func Load(path string) (*Config, error) {
	if path == "" {
		path = ".env"
	}
	// Missing .env is fine; only surface real read errors.
	if err := godotenv.Load(path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load %s: %w", path, err)
	}

	c := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AppEnv:      getEnv("APP_ENV", "development"),
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "appdb"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
	}

	if c.DatabaseURL == "" {
		c.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
		)
	}
	return c, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
