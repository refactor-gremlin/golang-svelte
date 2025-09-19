```mermaid
flowchart TD
    Frontend[Svelte 5 app] --> BFF[SvelteKit remote functions<br/>BFF layer]
    BFF --> Client[OpenAPI TS client]
    Client --> Gin[Gin router]
    Gin --> Handler[Module api: Gin handler]
    Handler -->|validate + map DTO| UseCase[Module app: use-case]
    UseCase -->|ports| Ports[Ports (interfaces)]
    Ports -->|implemented by| Infra[Infra adapters]
    Infra --> DB[(Database)]
    Infra --> JWT[JWT]
    Infra --> Ext[External HTTP APIs]
    UseCase --> Domain[Domain entities/VOs]
    Domain --> UseCase
    UseCase -->|result or typed error| Handler
    Handler -->|HTTP DTO + status| Client
    Client --> BFF
    BFF -->|transformed data| Frontend
```

### Corrected data flow
```text
Svelte 5 app -> SvelteKit remote functions (BFF layer)
  -> OpenAPI TS client -> Gin router -> Module api: Gin handler
    -> Module app: use-case + ports -> Infra adapters (DB/JWT/HTTP)
      -> Domain entities/VOs -> back up to use-case
    -> HTTP DTO + status -> OpenAPI client
  -> SvelteKit BFF transforms/formats -> Svelte 5 app
```

### Goals
- **Clarity**: A single, repeatable request flow for every endpoint with proper BFF separation.
- **Separation**: `api -> app -> domain` with `infra` implementing `app` ports; `platform` wires.
- **Scalability**: Add modules without growing cross-coupling or central routers.

### Principles
- **One-way dependencies**: `api -> app -> domain`. `infra -> app` (adapters implement ports). `platform` only composes and wires.
- **Thin handlers**: HTTP-only (bind/validate, call one use-case, map errors/status, return DTOs).
- **Explicit use-cases**: Commands/queries orchestrate via ports; no Gin/GORM imports.
- **Pure domain**: Entities/VOs/invariants; no HTTP/DB/external knowledge.
- **Replaceable infra**: DB/JWT/HTTP clients live behind interfaces defined in `app`.
- **Uniform error mapping**: Typed errors from `app` mapped to consistent HTTP statuses in `api`.

### Phased plan

1) Scaffold platform and module structure
- Create `internal/modules/{auth,pokemon}/{api,app,domain,infra}`.
- Create `internal/platform/{httpserver,persistence,logging,config,security}`.
- Add Gin + gin-swagger to `cmd/server` via `platform/httpserver`.
- Keep existing endpoints working; introduce a temporary compatibility route group if needed.

2) Vertical-slice: pokemon (first mover)
- Move `pokeapi_service.go` to `modules/pokemon/infra` (adapter). Define a `RandomPokemonPort` in `modules/pokemon/app`.
- Move current pokemon `dtos.go` and `service.go` into `modules/pokemon/app` and shape as a use-case.
- Add `modules/pokemon/api` with Gin handlers + router; HTTP DTOs live here.
- Wire in `platform/httpserver` (e.g., `pokemonapi.Mount(r, handlers)`). Keep behavior identical.

3) Vertical-slice: auth
- Define ports in `modules/auth/app/ports.go`: `UserRepo`, `TokenGen`, `Hasher`.
- Move JWT/password implementations to `modules/auth/infra`; move user repository to `modules/auth/infra` (backed by `platform/persistence`).
- Implement `modules/auth/app` use-cases (`Login`, `Register`). Return typed errors (`ErrInvalidCredentials`, etc.).
- Add `modules/auth/api` handlers + router, with uniform error mapping. Keep Swagger annotations here.

4) Consolidate routing & middlewares
- Replace ad-hoc mux wiring with `platform/httpserver` that mounts all module routers.
- Centralize middlewares (recovery, logging, request-id, CORS) at `platform/httpserver`.
- Keep `cmd/server` responsible for wiring: construct infra adapters, build app services, inject into api handlers, mount routers.

5) Cross-cutting hardening
- `platform/config`: central env loading; pass typed config down to wiring.
- `platform/logging`: structured logger; inject into middlewares and optionally adapters.
- `platform/persistence`: DB initialization/migrations; hand out `*gorm.DB` (or alternative) to infra.
- `platform/security`: shared helpers (only if truly cross-cutting).

6) Tests and documentation
- Unit tests colocated per module (near code). Integration/e2e in `tests/` (HTTP + DB).
- Per-module `README.md` describing public use-cases and domain invariants.
- Swagger tags per module; keep HTTP DTOs in `api/models.go` 1:1 with Swagger schemas.

7) Cleanup and deprecation
- Remove `internal/presentation`, `internal/application` (old), and `internal/infrastructure` entries as modules take over.
- Ensure imports comply with boundaries; run `go vet`, linters, and e2e.

### Wiring snapshot (Gin + gin-swagger)
```go
// cmd/server/main.go (sketch)
r := gin.New()
r.Use(gin.Recovery(), gin.Logger())

// platform wiring
db := persistence.NewAppDB(...)
tokenGen := authentication.NewJWTTokenGenerator(...)
hasher := security.NewHMACPasswordHasher()

// auth module
userRepo := authinfra.NewGormUserRepo(db)
authSvc := authapp.NewService(userRepo, tokenGen, hasher)
authHandlers := authapi.NewHandlers(authSvc)

// pokemon module
pokeAdapter := pokemoninfra.NewPokeAPIAdapter(http.DefaultClient)
pokemonSvc := pokemonapp.NewService(pokeAdapter)
pokemonHandlers := pokemonapi.NewHandlers(pokemonSvc)

// routing
httpserver.MountAll(r, httpserver.Deps{
    Auth:    authHandlers,
    Pokemon: pokemonHandlers,
})

// swagger (consumed by SvelteKit BFF)
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
_ = r.Run(":" + port)
```

### SvelteKit BFF layer integration
The OpenAPI client generated from Swagger is used by SvelteKit remote functions:
```ts
// src/routes/auth/login/+page.server.ts (SvelteKit remote function)
import { client } from '$lib/api'; // OpenAPI client

export const actions = {
  login: async ({ request }) => {
    const data = await request.formData();
    const email = data.get('email');
    const password = data.get('password');

    const response = await client.auth.login({ email, password });
    // Transform/format response for frontend consumption
    return { success: true, token: response.accessToken };
  }
};
```

### Error mapping pattern (uniform)
```go
func mapLoginError(err error) (int, any) {
    switch {
    case app.IsErrInvalidCredentials(err):
        return http.StatusUnauthorized, gin.H{"error": "invalid_credentials"}
    case domain.IsErrRuleViolated(err):
        return http.StatusUnprocessableEntity, gin.H{"error": "domain_violation"}
    default:
        return http.StatusInternalServerError, gin.H{"error": "internal_error"}
    }
}
```

### Definition of done (per module)
- **API**: Gin router + handlers; request/response DTOs; Swagger annotations; uniform error mapping.
- **App**: Commands/queries/use-cases; ports defined; no HTTP/DB imports; typed errors.
- **Domain**: Entities/VOs/invariants; pure and tested.
- **Infra**: Adapters for ports (DB/JWT/external); no handler/use-case logic.
- **Wiring**: Registered in `platform/httpserver`; constructed in `cmd/server`.
- **Tests**: Unit tests near code; integration/e2e pass.

This structure keeps the data flow explicit and repeatable for every endpoint while allowing the codebase to scale safely as features and endpoints grow.


