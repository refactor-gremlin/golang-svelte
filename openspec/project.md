# Project Context

## Purpose
MySvelteApp is a modern full-stack application demonstrating best practices in web development with a Python backend and Svelte frontend. The project serves as a template and learning platform for building scalable, type-safe web applications with comprehensive observability.

## Tech Stack

### Backend (Python 3.11+)
- **Framework**: FastAPI (v0.115+)
- **Database**: SQLite with SQLAlchemy ORM (v2.0+)
- **Authentication**: JWT tokens with HMACSHA512 password hashing
- **API Documentation**: Swagger/OpenAPI with FastAPI automatic docs
- **Observability**: OpenTelemetry tracing with Tempo backend
- **Logging**: Structured logging with Python logging
- **Architecture**: Clean Architecture with DDD patterns

### Frontend (SvelteKit 5)
- **Framework**: SvelteKit 2.38+ with Svelte 5 (runes-based)
- **Language**: TypeScript with strict mode
- **Build Tool**: Vite 7.0+
- **Styling**: Tailwind CSS 4.0 with custom components
- **UI Components**: Bits UI, Lucide Svelte icons
- **API Client**: Auto-generated from OpenAPI spec (@hey-api/openapi-ts)
- **Testing**: Vitest (unit) + Playwright (E2E)
- **Code Quality**: ESLint + Prettier + Svelte Check

### DevOps & Infrastructure
- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Docker Compose for development and production
- **Observability Stack**: Grafana + Loki + Tempo + Promtail
- **OpenTelemetry**: Distributed tracing and metrics
- **Development**: Hot reload with concurrent frontend/backend dev

## Project Conventions

### Code Style

#### Python Backend
- **Naming**: Python conventions (snake_case for functions/variables, PascalCase for classes)
- **File Structure**: Clean Architecture layers (domain, app, infra, api)
- **Error Handling**: Exceptions with custom error types
- **Logging**: Structured logging with Python logging module
- **Comments**: Docstrings for all public functions and classes

#### TypeScript/Frontend
- **Naming**: camelCase for variables, PascalCase for components/types
- **File Naming**: kebab-case for files, PascalCase for component folders
- **Imports**: Grouped imports (external, internal, relative)
- **Type Safety**: Strict TypeScript mode enabled
- **Formatting**: Prettier with 2-space indentation

### Architecture Patterns

#### Backend Architecture
```
├── app/
│   ├── main.py         # Application entry point
│   ├── config.py       # Configuration management
│   ├── database.py     # Database setup
│   ├── modules/        # Domain modules (auth, pokemon, etc.)
│   │   ├── app/        # Application layer (use cases)
│   │   ├── api/        # Interface layer (routes, schemas)
│   │   ├── domain/     # Domain entities
│   │   └── infra/      # Infrastructure layer (repositories, external APIs)
│   └── platform/       # Cross-cutting concerns
│       ├── logging.py  # Logging setup
│       └── tracing.py  # OpenTelemetry setup
```

#### Frontend Architecture
- **SvelteKit 5**: Modern runes-based reactivity (`$state`, `$derived`)
- **Type-Safe APIs**: Auto-generated client from OpenAPI spec
- **Component Organization**: Reusable UI components with shadcn/ui patterns
- **Route-Based**: File-system routing with layouts

### Testing Strategy
- **Backend**: Unit tests with pytest
- **Frontend**: Vitest for unit testing, Playwright for E2E testing
- **API Testing**: Integration tests for endpoints
- **Coverage**: Aim for 80%+ code coverage
- **CI/CD**: Automated testing in development workflow

### Git Workflow
- **Branching**: Feature branches from main (`feature/auth-improvement`)
- **Commits**: Conventional commits format (`feat:`, `fix:`, `docs:`, etc.)
- **PR Process**: Code review required for all changes
- **Hooks**: Husky for pre-commit checks (linting, formatting)

## Domain Context

### Authentication & Authorization
- JWT-based authentication with secure HTTP-only cookies
- Password validation with complexity requirements
- User registration and login flows
- Session management with automatic refresh

### API Design
- RESTful API design principles
- OpenAPI 3.0 specification for documentation
- Type-safe request/response handling
- Error handling with appropriate HTTP status codes

### Pokemon Integration
- External Pokemon API integration
- Data transformation and caching
- Error handling for external service failures

## Important Constraints

### Technical Constraints
- **Database**: SQLite for development, configurable for production
- **Authentication**: Must use secure password hashing (HMACSHA512)
- **API Documentation**: Must be auto-generated from OpenAPI spec
- **Type Safety**: End-to-end type safety between frontend and backend
- **Performance**: Response times under 100ms for most operations

### Development Constraints
- **Environment Variables**: All configuration via environment variables
- **Hot Reload**: Development server must support hot reload
- **Containerization**: Must work in Docker for consistency
- **Observability**: Full tracing and logging required

## External Dependencies

### Development Dependencies
- **Python Packages**: All Python dependencies managed via requirements.txt and uv
- **NPM Packages**: Frontend dependencies via package.json
- **Docker Images**: Official images for all services

### Production Services
- **Pokemon API**: External Pokemon data API (pokeapi.co)
- **Observability**: Self-hosted Grafana, Loki, Tempo stack
- **Database**: SQLite (development), configurable for production

### Security Considerations
- **JWT Secret**: Must be properly secured in production
- **Environment Files**: Must not be committed to version control
- **HTTPS**: Required for production deployments
- **CORS**: Properly configured for API access
