package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type config struct {
	Port           string
	DSN            string
	AllowedOrigins []string
	RateLimit      int
}

func loadConfig() (*config, error) {
	port := envOr("PORT", "8080")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	origins := strings.Split(envOr("ALLOWED_ORIGINS", "http://localhost:3000"), ",")

	rl, err := strconv.Atoi(envOr("RATE_LIMIT_RPS", "100"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_RPS: %w", err)
	}

	return &config{
		Port:           port,
		DSN:            dsn,
		AllowedOrigins: origins,
		RateLimit:      rl,
	}, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
