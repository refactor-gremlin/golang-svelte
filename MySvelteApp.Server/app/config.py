"""Configuration management for the application."""

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
        extra="ignore",
    )

    # Server configuration
    server_port: str = "8080"
    database_dsn: str = "sqlite:///./mysvelteapp.db"

    # JWT configuration
    jwt_key: str = "base64:YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="
    jwt_issuer: str = "mysvelteapp"
    jwt_audience: str = "mysvelteapp"
    jwt_access_token_lifetime_hours: int = 24

    # OpenTelemetry configuration
    otel_service_name: str = "mysvelteapp-server"
    otel_service_version: str = "1.0.0"
    environment: str = "development"
    otel_exporter_otlp_traces_endpoint: str = "localhost:4317"

    @property
    def port(self) -> str:
        """Get server port."""
        return self.server_port

    @property
    def service_name(self) -> str:
        """Get service name for OpenTelemetry."""
        return self.otel_service_name

    @property
    def service_version(self) -> str:
        """Get service version for OpenTelemetry."""
        return self.otel_service_version


def load_settings() -> Settings:
    """Load application settings from environment."""
    return Settings()

