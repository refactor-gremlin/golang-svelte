package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultPort             = "8080"
	defaultDatabaseDSN      = "file:mysvelteapp.db?cache=shared&_fk=1"
	defaultJWTKey           = "base64:YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="
	defaultJWTIssuer        = "mysvelteapp"
	defaultJWTAudience      = "mysvelteapp"
	defaultJWTLifetimeHours = 24
)

// Server holds runtime configuration needed to start the API server.
type Server struct {
	Port                   string
	DatabaseDSN            string
	JWTKey                 string
	JWTIssuer              string
	JWTAudience            string
	JWTAccessLifetimeHours int
}

// Load reads configuration from environment variables, applying defaults where required.
func Load() (Server, error) {
	cfg := Server{
		Port:                   getEnv("SERVER_PORT", defaultPort),
		DatabaseDSN:            getEnv("DATABASE_DSN", defaultDatabaseDSN),
		JWTKey:                 getEnv("JWT_KEY", defaultJWTKey),
		JWTIssuer:              getEnv("JWT_ISSUER", defaultJWTIssuer),
		JWTAudience:            getEnv("JWT_AUDIENCE", defaultJWTAudience),
		JWTAccessLifetimeHours: defaultJWTLifetimeHours,
	}

	if lifetimeStr := os.Getenv("JWT_ACCESS_TOKEN_LIFETIME_HOURS"); lifetimeStr != "" {
		parsed, err := strconv.Atoi(lifetimeStr)
		if err != nil {
			return Server{}, fmt.Errorf("parse JWT_ACCESS_TOKEN_LIFETIME_HOURS: %w", err)
		}
		cfg.JWTAccessLifetimeHours = parsed
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
