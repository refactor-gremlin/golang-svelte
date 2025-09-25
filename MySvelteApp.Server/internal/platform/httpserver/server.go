package httpserver

import (
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// New constructs a gin.Engine with the baseline middlewares configured.
func New(logger *slog.Logger, serviceName string) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	if serviceName == "" {
		serviceName = "mysvelteapp-server"
	}
	engine.Use(otelgin.Middleware(serviceName))

	if logger != nil {
		engine.Use(loggingMiddleware(logger))
	}

	return engine
}

func loggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		clientIP := c.ClientIP()
		latency := time.Since(start)

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("request failed",
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"status", status,
					"duration_ms", latency.Milliseconds(),
					"client_ip", clientIP,
					"error", err.Error(),
				)
			}
			return
		}

		level, statusMsg := getStatusInfo(status)
		if status >= 400 {
			logger.Log(c, level, "request completed",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", status,
				"message", statusMsg,
				"duration_ms", latency.Milliseconds(),
				"client_ip", clientIP,
			)
			return
		}

		logger.Info("request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"duration_ms", latency.Milliseconds(),
			"client_ip", clientIP,
		)
	}
}

func getStatusInfo(statusCode int) (slog.Level, string) {
	switch {
	case statusCode >= 500:
		return slog.LevelError, "server error"
	case statusCode == 404:
		return slog.LevelWarn, "not found"
	case statusCode == 405:
		return slog.LevelWarn, "method not allowed"
	case statusCode == 401:
		return slog.LevelWarn, "unauthorized"
	case statusCode == 403:
		return slog.LevelWarn, "forbidden"
	case statusCode >= 400:
		return slog.LevelWarn, "client error"
	case statusCode >= 300:
		return slog.LevelInfo, "redirect"
	default:
		return slog.LevelInfo, ""
	}
}
