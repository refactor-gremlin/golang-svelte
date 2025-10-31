# Module Architecture Standards

This document defines the standardized structure and conventions that should be followed across all API modules in this project.

## Standard Module Structure

Every module should follow this consistent directory structure:

```
module-name/
├── api/                 # HTTP Layer
│   ├── models.go       # Request/response models, JSON contracts
│   ├── handlers.go     # HTTP handlers, request processing
│   └── routes.go       # Route definitions, middleware setup
├── app/                 # Business Logic Layer
│   ├── service.go      # Core business logic, use cases
│   ├── ports.go        # Interface definitions (if needed)
│   ├── commands.go     # Command/Query DTOs (if complex)
│   └── errors.go       # Custom error types (if needed)
├── domain/              # Core Business Entities (if complex)
│   └── entities.go     # Domain models, business rules
└── infra/               # Infrastructure Implementations
    ├── persistence/    # Database repositories (if needed)
    ├── external/       # External API clients (if needed)
    └── security/       # Security implementations (if needed)
```

## Data Flow Patterns

### Complex Business Logic (e.g., Authentication)
```
HTTP Request → API Layer → Service Layer → Domain Layer → Infrastructure Layer
HTTP Response ← API Layer ← Service Layer ← Domain Layer ← Infrastructure Layer
```

### Simple Data Retrieval (e.g., External API calls)
```
HTTP Request → API Layer → Service Layer → External API
HTTP Response ← API Layer ← Service Layer ← External API Response
```

## File Naming Conventions

### Standard Files (Always Present)
- `models.go` - API models and response structures
- `handlers.go` - HTTP request handlers
- `routes.go` - Route definitions
- `service.go` - Business logic service

### Conditional Files (Include when needed)
- `ports.go` - Interface definitions for dependencies
- `commands.go` - Complex request/response DTOs
- `queries.go` - Query-specific DTOs
- `errors.go` - Custom error types
- `entities.go` - Domain entities
- `repository.go` - Data access implementations
- `client.go` - External API clients
- `mapper.go` - Data transformation logic

## Layer Responsibilities

### API Layer (`api/`)
- HTTP request/response handling
- JSON serialization/deserialization
- HTTP status codes and error responses
- Input validation and binding
- Swagger documentation
- **Do NOT**: Business logic, database operations

### Service Layer (`app/`)
- Business logic and use case orchestration
- Coordination between infrastructure components
- Transaction management
- Input validation (business rules)
- **Do NOT**: HTTP concerns, direct database access

### Domain Layer (`domain/`)
- Core business entities and invariants
- Domain-specific business rules
- Entity state management
- **Include only**: When complex business rules exist

### Infrastructure Layer (`infra/`)
- External service integrations
- Database persistence
- Security implementations
- External API clients
- **Do NOT**: Business logic

## Documentation Standards

### Package Documentation
Each package must include comprehensive package-level documentation explaining:
- Purpose and responsibilities
- Data flow patterns
- Key design decisions
- Dependencies and interfaces

### Function Documentation
All public functions must include:
- Clear description of purpose
- Data flow explanation
- Error conditions and handling
- External dependencies
- Security considerations (if applicable)

### Type Documentation
All public types must include:
- Purpose and usage context
- Field explanations for complex structures
- Design rationale for architectural decisions

## Error Handling Standards

### Custom Error Types
Define custom error types in `errors.go`:
```go
type ValidationError struct {
    Message string
}

func (e ValidationError) Error() string {
    return e.Message
}
```

### Error Mapping
Map application errors to HTTP status codes in handlers:
- `ValidationError` → 400 Bad Request
- `NotFoundError` → 404 Not Found  
- `ConflictError` → 409 Conflict
- `UnauthorizedError` → 401 Unauthorized
- `InternalError` → 500 Internal Server Error

## Testing Standards

### Test Organization
- Unit tests for each layer
- Integration tests for service layer
- Mock interfaces for external dependencies
- Test coverage for all error paths

### Test File Naming
- `*_test.go` in same package as source
- `mock_*.go` for mock implementations
- `testdata/` directory for test fixtures

## Dependencies and Interfaces

### Dependency Injection
- Use constructor injection for all dependencies
- Define interfaces in the consuming layer
- Provide implementations from the outside
- Enable easy testing with mocks

### Interface Naming
- Use descriptive names: `UserRepository`, `PasswordHasher`
- Keep interfaces focused and small
- Define interfaces in the consuming package

## When to Simplify Architecture

### Use Full 5-Layer Architecture When:
- Complex business rules exist
- Multiple data sources are involved
- Rich domain model is required
- Transaction management is needed

### Use Simplified 3-Layer Architecture When:
- Simple data retrieval operations
- Single external API integration
- No complex business logic
- Read-only operations

## Examples in This Codebase

### Auth Module (Complex - 5 layers)
- Business logic: validation, password hashing, JWT generation
- Persistence: User repository with database operations
- Domain: User entity with business invariants

### Pokemon Module (Simple - 3 layers)  
- Business logic: Simple data transformation
- External integration: Single API call to PokeAPI
- No persistence or complex domain rules

## Quality Checklists

### Pokemon Module
- Seed pseudo-random generators exactly once and guard zero/negative counts before calling `rand.Intn`.
- Cache slow-changing metadata (e.g., total Pokémon count) with a sensible TTL and invalidate gracefully on failures.
- Centralize HTTP response decoding into typed helpers to reduce repetition and ease testing.
- Surface external API failures as `503 Service Unavailable` equivalents and include trace/logger hints for observability.
- Cover success, retry, and error paths with unit tests that use mockable HTTP clients and deterministic random sources.

### Auth Module
- Keep request/response DTOs in `api/models.go` and reuse them across handlers and Swagger annotations to avoid drift.
- Enforce username/email invariants inside the domain factory while sharing validation helpers with the service layer to stay DRY.
- Let the database enforce uniqueness and translate `gorm.ErrDuplicatedKey` (and driver-specific variants) into `ConflictError`.
- Ensure password hashing and comparison remain constant time and centralize configuration (salt length, hashing algorithm) for reuse.
- Add integration tests to cover concurrent registration/login flows and confirm correct error mapping.

Follow these standards to ensure consistency, maintainability, and clear data flow across all modules.
