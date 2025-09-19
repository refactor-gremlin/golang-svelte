# GitHub Actions CI/CD

This repository uses GitHub Actions for continuous integration and automated testing.

## Workflows

### CI Pipeline (`ci.yml`)
Runs on every push and pull request to `main` and `develop` branches.

**What it does:**
- ✅ Sets up .NET 9.0 environment
- ✅ Sets up Node.js 20 environment
- ✅ Installs dependencies for both backend and frontend
- ✅ Builds the Svelte client application
- ✅ Builds the .NET solution
- ✅ Runs .NET unit tests (xUnit)
- ✅ Runs client linting (ESLint) — in code-quality.yml
- ✅ Runs client type checking (TypeScript)
- ✅ Runs client unit tests (Vitest)
- 📊 Uploads test results as artifacts

### Code Quality (`code-quality.yml`)
Runs on every push and pull request to `main` and `develop` branches.

**What it does:**
- ✅ Runs ESLint for code style checking
- ✅ Runs TypeScript type checking
- ✅ Runs Prettier format checking

## Local Development

To run tests locally before pushing:

```bash
# Backend tests
dotnet test svelte-NET-Test.sln

# Frontend tests
npm ci --prefix MySvelteApp.Client
npm run test:unit --prefix MySvelteApp.Client

# Frontend linting
npm run lint --prefix MySvelteApp.Client

# Frontend type checking
npm run check --prefix MySvelteApp.Client
```

### Skip Client Build

For faster backend-only development and CI runs, you can skip the client build by setting the `SkipClientBuild=true` environment variable:

```bash
# Build .NET without client
dotnet build svelte-NET-Test.sln /p:SkipClientBuild=true

# Run tests without building client
dotnet test svelte-NET-Test.sln /p:SkipClientBuild=true
```

## Test Results

Test results are automatically uploaded as artifacts when tests fail, allowing you to:
- Download and analyze test results
- View detailed test logs
- Debug failing tests

## Branch Protection

Consider setting up branch protection rules that require:
- ✅ All CI checks to pass
- ✅ Code review approval
- ✅ Up-to-date branches

This ensures code quality and prevents breaking changes from being merged.
