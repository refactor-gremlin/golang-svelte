"""OpenTelemetry tracing setup."""

import os
from contextlib import asynccontextmanager
from typing import Optional

from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.semconv.resource import ResourceAttributes

from app.platform.logging import get_logger

logger = get_logger("tracing")


def get_env(key: str, fallback: str) -> str:
    """Get environment variable with fallback."""
    return os.getenv(key, fallback)


class TracingProvider:
    """OpenTelemetry tracing provider."""

    def __init__(self, service_name: str, service_version: str, environment: str = "development"):
        """Initialize tracing provider."""
        self.service_name = service_name
        self.service_version = service_version
        self.environment = environment
        self.provider: Optional[TracerProvider] = None

    def setup(self):
        """Set up OpenTelemetry tracing."""
        # Create resource with service information
        resource = Resource.create(
            {
                ResourceAttributes.SERVICE_NAME: self.service_name,
                ResourceAttributes.SERVICE_VERSION: self.service_version,
                ResourceAttributes.DEPLOYMENT_ENVIRONMENT: self.environment,
            }
        )

        # Create OTLP exporter
        endpoint = get_env("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "localhost:4317")
        exporter = OTLPSpanExporter(
            endpoint=endpoint,
            insecure=True,  # Use for local development
        )
        logger.info(f"Using OTLP trace exporter (direct to Tempo)", extra={"endpoint": endpoint})

        # Create tracer provider
        self.provider = TracerProvider(resource=resource)
        self.provider.add_span_processor(BatchSpanProcessor(exporter))

        # Set global tracer provider
        trace.set_tracer_provider(self.provider)

    def shutdown(self):
        """Shutdown tracing provider."""
        if self.provider:
            self.provider.shutdown()


def instrument_fastapi(app):
    """Instrument FastAPI application with OpenTelemetry."""
    FastAPIInstrumentor.instrument_app(app)

