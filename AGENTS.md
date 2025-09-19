# Repository Guidelines

## Project Structure & Module Organization
- `MySvelteApp.Client/`: SvelteKit 5 frontend; remote functions live under `src/routes/**/feature.remote.ts`, UI components in `src/lib/components/`, utilities in `src/lib/utils.ts`.
- `MySvelteApp.Server/`: .NET 9 backend organized by Clean Architecture (`Domain/`, `Application/`, `Infrastructure/`, `Presentation/`); shared test helpers in `Tests/TestUtilities/`.
- `observability/`: telemetry collectors and docker assets; keep instrumentation configs here.
- Root `Dockerfile` and compose files wire the full stack; avoid editing generated code under `MySvelteApp.Client/src/api/schema/`.

## Build, Test, and Development Commands
- `npm run dev`: starts SvelteKit dev server and .NET API concurrently.
- `npm run build`: compiles the Svelte client (`vite build`) and builds the API (`dotnet build`).
- `npm run docker:dev` / `npm run docker:prod`: spin up full stack via docker compose.
- `dotnet test MySvelteApp.Server/Tests/MySvelteApp.Server.Tests.csproj`: executes backend unit tests; append `--collect:"XPlat Code Coverage"` for coverage.
- `npm run test:unit`, `npm run test:e2e`, `npm run lint`, `npm run check` (run from `MySvelteApp.Client/` or with `--prefix`).

## Coding Style & Naming Conventions
- Follow Prettier (default Svelte 2-space indentation) and ESLint rules; run `npm run format` + `npm run lint` before committing.
- Use Svelte 5 runes (`$state`, `$derived`) and alias imports (`$lib`, `$api`).
- TypeScript interfaces in `camelCase` files; domain classes in `PascalCase` under `Domain/Entities/`.
- .NET code follows nullable-enabled C# 12 with guard clauses; keep domain logic inside entities/services.

## Testing Guidelines
- Frontend: Vitest for units, Playwright for E2E; name specs `*.test.ts` or `*.spec.ts` near the code under test.
- Backend: xUnit with FluentAssertions and Moq; mirror application namespace in `Tests/` and keep fixture builders in `TestUtilities/`.
- Target meaningful coverage for application and domain layers; add regression tests with each bug fix.

## Commit & Pull Request Guidelines
- Follow existing imperative, descriptive commit style (e.g., `Refactor ESLint configuration ...`).
- Keep commits scoped to a single concern; include tooling runs (format/lint/test) in the same commit when relevant.
- Pull requests should explain motivation, summarize architectural impact, list verification steps, and reference related issues or tickets.
- Attach screenshots or console output for UI or API behavior changes; note any follow-up tasks or TODOs explicitly.

## Architecture Reminders
- Preserve the Domain → Application → Infrastructure dependency flow; never call repositories from controllers directly.
- Prefer remote functions over manual fetches; regenerate API clients with `npm run generate-api-classes` after contract changes.
- Store environment-sensitive values in `appsettings.*.json` and `.env` files, not in source.
