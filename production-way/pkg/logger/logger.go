// Package logger provides a structured slog.Logger configured for the environment.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a production-ready structured logger.
// level can be "debug", "info", "warn", "error" (case-insensitive). Defaults to "info".
func New(level string) *slog.Logger {
	var l slog.Level
	switch strings.ToLower(level) {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     l,
		AddSource: l == slog.LevelDebug, // show file:line only in debug mode
	}

	// JSON in production, text locally (detected via LOG_FORMAT env var).
	var handler slog.Handler
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
