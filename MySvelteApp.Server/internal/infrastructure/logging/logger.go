package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// LogLevel represents the logging level
type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

// Config holds the logger configuration
type Config struct {
	Level     LogLevel
	Format    string // "json" or "text"
	Output    string // "stdout", "stderr", or file path
	AddSource bool   // Include source file and line info
}

// DefaultConfig returns a default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:     Info,
		Format:    "text",
		Output:    "stdout",
		AddSource: false,
	}
}

// NewLogger creates a new structured logger with the given configuration
func NewLogger(config Config) *slog.Logger {
	var writer io.Writer
	switch strings.ToLower(config.Output) {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		// Assume it's a file path
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// Fall back to stdout if file can't be opened
			writer = os.Stdout
		} else {
			writer = file
		}
	}

	var level slog.Level
	switch config.Level {
	case Debug:
		level = slog.LevelDebug
	case Info:
		level = slog.LevelInfo
	case Warn:
		level = slog.LevelWarn
	case Error:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	handlerOpts := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.AddSource,
	}

	if strings.ToLower(config.Format) == "json" {
		handler = slog.NewJSONHandler(writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(writer, handlerOpts)
	}

	return slog.New(handler)
}

// NewDefaultLogger creates a logger with default configuration
func NewDefaultLogger() *slog.Logger {
	return NewLogger(DefaultConfig())
}
