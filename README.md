# Go + SvelteKit Starter

A full-stack application that pairs a Go 1.23 API (Gin + GORM) with a modern SvelteKit 2 frontend. The stack includes generated OpenAPI clients, experimental Svelte remote functions, JWT authentication, and optional observability via OpenTelemetry, Grafana, Tempo, Loki, and Promtail.

- Go backend exposes REST endpoints for auth and a Pokémon demo, backed by SQLite and instrumented with OpenTelemetry.
- SvelteKit frontend consumes the API through generated, type-safe clients and remote functions.
- Dockerfiles and Compose files cover local development, production builds, and observability tooling.

## Project Structure

```
📁 Project Root
├── MySvelteApp.Client/      # SvelteKit application (Svelte 5, Tailwind, experimental remote functions)
├── MySvelteApp.Server/      # Go service (Gin, GORM, JWT, Swagger, OpenTelemetry)
├── observability/           # Grafana + Loki + Tempo + Promtail provisioning
├── docker-compose*.yml      # Compose files for dev, prod, and observability
├── Dockerfile               # Multi-stage build that bundles client + server
├── package.json             # Root scripts (concurrently run client + server)
└── README.md                # You are here
```

## Tech Highlights

- **Backend**: Go 1.23, Gin, GORM (SQLite), JWT auth, Swagger docs, modular architecture under `internal/`
- **Frontend**: SvelteKit 2, Svelte 5 runes, Tailwind CSS 4, shadcn/bits UI, experimental remote functions (queries/forms/commands)
- **Integration**: `@hey-api/openapi-ts` generates a typed SDK consumed by remote functions such as `src/routes/(app)/pokemon/data.remote.ts`
- **Observability**: OpenTelemetry tracing (SDK on both Go and SvelteKit), optional Grafana/Loki/Tempo stack via `observability.compose.yml`
- **Tooling**: Vitest, Playwright, ESLint, Prettier, Husky, Go tests with `go test ./...`

## Prerequisites

- Node.js 20+
- Go 1.23+
- Docker & Docker Compose (for container workflows or observability stack)

## Install Dependencies

```bash
# Install root dev tooling (concurrently, husky, etc.)
npm install

# Install SvelteKit dependencies
npm install --prefix MySvelteApp.Client

# Download Go modules (optional if go build/test runs)
(cd MySvelteApp.Server && go mod download)
```

## Local Development

### Run everything together

```bash
npm run dev
```

- SvelteKit dev server: http://localhost:5173
- Go API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html

Logs for both services stream in the same terminal (via `concurrently`).

### Run servers independently

```bash
# API
cd MySvelteApp.Server
go run ./cmd/server

# Frontend (new terminal)
npm run dev --prefix MySvelteApp.Client
```

### Regenerate the API client

The client SDK under `MySvelteApp.Client/src/routes` is generated from the Go Swagger spec. Regenerate after backend changes:

```bash
npm run generate-api-classes --prefix MySvelteApp.Client
```

## Testing & Quality

```bash
# Go unit tests
(cd MySvelteApp.Server && go test ./...)

# Frontend types + lint + unit tests
npm run check --prefix MySvelteApp.Client
npm run lint --prefix MySvelteApp.Client
npm run test:unit --prefix MySvelteApp.Client

# Playwright E2E (requires browsers installed)
npm run test:e2e --prefix MySvelteApp.Client
```

## Docker Workflows

- `npm run docker:dev` → starts `docker-compose.dev.yml` with hot-reload mounts (frontend on 5173, API on 8080)
- `npm run docker:prod` → builds production images via `docker-compose.yml` (frontend served on 3000 through the Node adapter, API on 8080)

The root `Dockerfile` is multi-stage: it builds the SvelteKit client, bundles assets into the Go image, and produces a single binary container exposing port 8080.

## Observability Stack

1. Launch the collector stack:
   ```bash
   docker compose -f observability.compose.yml up -d
   ```
2. Configure the API (`MySvelteApp.Server/internal/platform/config`) and SvelteKit (`src/instrumentation.server.js`) with OTLP endpoints if you want to ship traces/logs to Tempo/Loki.
3. Grafana will be available at http://localhost:3000 (anonymous admin enabled by default).

## Environment Variables

Backend defaults live in `MySvelteApp.Server/internal/platform/config/config.go`.

| Variable | Default | Purpose |
| --- | --- | --- |
| `SERVER_PORT` | `8080` | Port for the Go HTTP server |
| `DATABASE_DSN` | `file:mysvelteapp.db?cache=shared&_fk=1` | SQLite DSN (file stored next to the binary) |
| `JWT_KEY` | sample key | HMAC secret for JWT signing |
| `JWT_ISSUER` / `JWT_AUDIENCE` | `mysvelteapp` | JWT metadata |
| `JWT_ACCESS_TOKEN_LIFETIME_HOURS` | `24` | Override token TTL |
| `OTEL_SERVICE_NAME` | `mysvelteapp-server` | OpenTelemetry service name |
| `OTEL_SERVICE_VERSION` | `1.0.0` | Service version tag |
| `ENVIRONMENT` | `development` | Environment label |

Frontend environment values go into `MySvelteApp.Client/.env` and support entries like:

```
VITE_API_URL=http://localhost:8080
PUBLIC_API_ENDPOINT=http://localhost:8080
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318/v1/traces
```

## API Overview

| Route | Method | Description |
| --- | --- | --- |
| `/auth/register` | POST | Register a new user (username, email, password) |
| `/auth/login` | POST | Authenticate and receive a JWT |
| `/RandomPokemon` | GET | Fetch a random Pokémon demo payload |
| `/swagger/index.html` | GET | Interactive API reference |

Auth handlers issue JWTs stored as HTTP-only cookies on the frontend (`src/routes/(auth)/auth.remote.ts`). Passwords are hashed with an HMAC-based password hasher before persistence.

## Frontend Features

- Layout separation for authenticated vs application routes under `src/routes/(auth)` and `src/routes/(app)`
- Experimental SvelteKit remote functions (`query`, `form`, `command`) for type-safe, server-executed logic
- UI built with shadcn components in `src/lib/components/ui`
- Demo Pokémon page showcasing async runes, suspense boundaries, and transitions (`src/routes/(app)/pokemon/+page.svelte`)
- OpenTelemetry Node SDK initialized in `src/instrumentation.server.js`

## Backend Highlights

- Modular packages for auth and Pokémon features (`internal/modules/...`)
- Gin router configured via `internal/platform/httpserver`
- GORM-backed SQLite database setup in `internal/platform/persistence`
- JWT generation/validation with configurable lifetimes (`internal/modules/auth/infra/token`)
- Swagger docs generated under `internal/docs`
- Graceful shutdown and structured logging via the `logging` package

## Useful Commands

```bash
# Format frontend code
npm run format --prefix MySvelteApp.Client

# Build production artifacts
npm run build

# Run Go binary with a custom port
SERVER_PORT=9000 go run ./cmd/server (from MySvelteApp.Server)
```

## Next Steps

- Configure real databases by swapping the SQLite driver string in `DATABASE_DSN`
- Extend the OpenAPI spec and regenerate the client before consuming new endpoints
- Hook up authentication in the frontend routes under `src/routes/(auth)` to match your user flows

