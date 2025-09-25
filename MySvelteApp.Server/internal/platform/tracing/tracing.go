package tracing

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Provider manages the OpenTelemetry tracing provider
type Provider struct {
	provider trace.TracerProvider
	shutdown func(context.Context) error
	logger   *slog.Logger
}

// New creates a new tracing provider with the given configuration
func New(serviceName, serviceVersion string, logger *slog.Logger) (*Provider, error) {
	ctx := context.Background()

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			attribute.String("environment", getEnv("ENVIRONMENT", "development")),
		),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Determine exporter based on environment
	var exporter sdktrace.SpanExporter
	if isDevelopment() {
		// Use stdout exporter for local development
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
		}
		logger.Info("using stdout trace exporter for local development")
	} else {
		// Use OTLP exporter for production/staging
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(getEnv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "localhost:4317")),
			otlptracegrpc.WithInsecure(), // Use for local development, secure for production
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
		}
		logger.Info("using OTLP trace exporter", "endpoint", getEnv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "localhost:4317"))
	}

	// Create tracer provider
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	// Set global tracer provider
	otel.SetTracerProvider(provider)

	// Create shutdown function
	shutdown := func(ctx context.Context) error {
		if err := provider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
		return nil
	}

	return &Provider{
		provider: provider,
		shutdown: shutdown,
		logger:   logger,
	}, nil
}

// Tracer returns a tracer for the given name
func (p *Provider) Tracer(name string) trace.Tracer {
	return p.provider.Tracer(name)
}

// Shutdown gracefully shuts down the tracing provider
func (p *Provider) Shutdown(ctx context.Context) error {
	if p.shutdown != nil {
		return p.shutdown(ctx)
	}
	return nil
}

// isDevelopment checks if we're running in development mode
func isDevelopment() bool {
	return getEnv("ENVIRONMENT", "development") == "development"
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
