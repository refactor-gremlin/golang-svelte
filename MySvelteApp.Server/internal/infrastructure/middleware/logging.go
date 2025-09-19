package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware creates a middleware that logs HTTP requests and responses
type LoggingMiddleware struct {
	logger *slog.Logger
}

// NewLoggingMiddleware creates a new logging middleware instance
func NewLoggingMiddleware(logger *slog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // default status
		body:           bytes.NewBuffer(nil),
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

// Middleware returns the HTTP handler middleware function
func (lm *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture response details
		rw := newResponseWriter(w)

		// Get client IP
		clientIP := getClientIP(r)

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Determine log level and status message
		logLevel, statusMsg := getStatusInfo(rw.statusCode)

		// Simple, concise logging
		if rw.statusCode >= 400 {
			// Log errors with more details
			lm.logger.Log(nil, logLevel, "API Request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"message", statusMsg,
				"duration_ms", duration.Milliseconds(),
				"client_ip", clientIP,
			)
		} else {
			// Success responses - minimal logging
			lm.logger.Info("API Request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration_ms", duration.Milliseconds(),
			)
		}
	})
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if idx := bytes.IndexByte([]byte(xff), ','); idx > 0 {
			return xff[:idx]
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// getStatusInfo returns the appropriate log level and status message for a given status code
func getStatusInfo(statusCode int) (slog.Level, string) {
	switch {
	case statusCode >= 500:
		return slog.LevelError, "Server Error"
	case statusCode == 404:
		return slog.LevelWarn, "Not Found"
	case statusCode == 405:
		return slog.LevelWarn, "Method Not Allowed"
	case statusCode == 401:
		return slog.LevelWarn, "Unauthorized"
	case statusCode == 403:
		return slog.LevelWarn, "Forbidden"
	case statusCode >= 400:
		return slog.LevelWarn, "Client Error"
	case statusCode >= 300:
		return slog.LevelInfo, "Redirect"
	default:
		return slog.LevelInfo, ""
	}
}
