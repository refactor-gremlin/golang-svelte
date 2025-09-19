package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// LogLevel represents an application logging level.
type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

// Config holds logger configuration flags.
type Config struct {
	Level     LogLevel
	Format    string // "json" or "text"
	Output    string // "stdout", "stderr", or file path
	AddSource bool
}

// DefaultConfig returns the baseline logger setup.
func DefaultConfig() Config {
	return Config{
		Level:     Info,
		Format:    "text",
		Output:    "stdout",
		AddSource: false,
	}
}

// NewLogger builds a slog.Logger configured according to the provided Config.
func NewLogger(config Config) *slog.Logger {
	var writer io.Writer
	switch strings.ToLower(config.Output) {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
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

	handlerOpts := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.AddSource,
	}

	var handler slog.Handler
	if strings.EqualFold(config.Format, "json") {
		handler = slog.NewJSONHandler(writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(writer, handlerOpts)
	}

	return slog.New(handler)
}

// NewDefaultLogger returns a slog.Logger using the default configuration.
func NewDefaultLogger() *slog.Logger {
	return NewLogger(DefaultConfig())
}
