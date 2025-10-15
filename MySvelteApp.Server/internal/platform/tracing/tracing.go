package tracing

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
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
			semconv.ServiceInstanceIDKey.String(getEnv("SERVICE_INSTANCE_ID", serviceName+"-"+getEnv("HOSTNAME", "localhost"))),
			semconv.TelemetrySDKLanguageGo,
			semconv.TelemetrySDKNameKey.String("opentelemetry"),
			semconv.TelemetrySDKVersionKey.String("1.38.0"),
			attribute.String("environment", getEnv("ENVIRONMENT", "development")),
			attribute.String("host.name", getEnv("HOSTNAME", "localhost")),
		),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Always use OTLP exporter for consistent tracing
	// Connect directly to Tempo instead of OTEL collector for simpler setup
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(getEnv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "localhost:4317")),
		otlptracegrpc.WithInsecure(), // Use for local development, secure for production
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}
	logger.Info("using OTLP trace exporter (direct to Tempo)", "endpoint", getEnv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "localhost:4317"))

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



// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
