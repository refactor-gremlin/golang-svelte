# MySvelteApp Server - Python Backend

Python implementation of the MySvelteApp backend using FastAPI, SQLAlchemy, and OpenTelemetry.

## Features

- **FastAPI**: Modern, fast web framework with automatic OpenAPI documentation
- **SQLAlchemy**: SQL toolkit and ORM for database operations
- **JWT Authentication**: Secure token-based authentication
- **OpenTelemetry**: Distributed tracing support
- **Clean Architecture**: Modular structure with domain/app/infra/api layers

## Project Structure

```
app/
├── main.py                 # FastAPI app initialization
├── config.py               # Configuration management
├── database.py             # SQLAlchemy setup
├── modules/
│   ├── auth/               # Authentication module
│   │   ├── domain/         # Domain entities
│   │   ├── app/            # Business logic
│   │   ├── infra/          # Infrastructure (persistence, security, tokens)
│   │   └── api/            # HTTP routes and schemas
│   └── pokemon/            # Pokemon module
│       ├── app/            # Business logic
│       └── api/            # HTTP routes
└── platform/               # Platform services (logging, tracing)
```

## Prerequisites

- Python 3.11+
- pip or poetry

## Installation

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Set up environment variables (copy `.env.example` to `.env` and configure):
```bash
cp .env.example .env
```

3. Run the server:
```bash
python -m app.main
```

Or using uvicorn directly:
```bash
uvicorn app.main:app --reload --port 8080
```

## Configuration

Environment variables:
- `SERVER_PORT`: Server port (default: 8080)
- `DATABASE_DSN`: Database connection string (default: sqlite:///./mysvelteapp.db)
- `JWT_KEY`: JWT signing key (base64 encoded)
- `JWT_ISSUER`: JWT issuer
- `JWT_AUDIENCE`: JWT audience
- `JWT_ACCESS_TOKEN_LIFETIME_HOURS`: Token lifetime in hours
- `OTEL_SERVICE_NAME`: OpenTelemetry service name
- `OTEL_SERVICE_VERSION`: Service version
- `ENVIRONMENT`: Environment (development/production)
- `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`: OpenTelemetry endpoint

## API Documentation

Once the server is running, visit:
- Swagger UI: http://localhost:8080/docs
- ReDoc: http://localhost:8080/redoc

## Database

The application uses SQLite by default. The database schema is automatically created on startup.

## Development

Run with hot reload:
```bash
uvicorn app.main:app --reload
```

## Docker

Build and run with Docker:
```bash
docker build -t mysvelteapp-server .
docker run -p 8080:8080 mysvelteapp-server
```

