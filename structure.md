# Project Architecture & Code Organization Guide

This document serves as a comprehensive guide for AI assistants working on this SvelteKit + .NET full-stack application. It explains the architectural patterns, code organization principles, and where different types of code should live.

## 🏗️ Overall Architecture

This is a **DDD-inspired Clean Architecture** full-stack application using:
- **Frontend**: SvelteKit 2.22.0 with Svelte 5.0.0
- **Backend**: .NET 9.0 Web API with DDD tactical patterns
- **Communication**: Type-safe remote functions (experimental SvelteKit feature)
- **Database**: Entity Framework Core (currently in-memory for development)

### DDD Characteristics Present ✅
- **Layered Architecture**: Domain → Application → Infrastructure → Presentation
- **Rich Domain Entities**: Business rules and validation in entities
- **Repository Pattern**: Interfaces in Application, implementations in Infrastructure
- **Application Services**: Orchestrate domain operations and business workflows
- **Ubiquitous Language**: Business concepts (User, Pokemon, Weather) over technical terms
- **Bounded Contexts**: Separate contexts for Authentication, Pokemon, Weather
- **Dependency Inversion**: Domain depends on abstractions defined in Application layer

### DDD Patterns Missing ❌
- **Aggregates & Aggregate Roots**: No explicit aggregate boundaries
- **Value Objects**: Email/Username could be immutable value objects
- **Domain Events**: No domain events for business processes
- **Domain Services**: Business logic in entities/application services, no pure domain services
- **Specifications**: No specification pattern for complex queries
- **Factories**: No domain object factories for complex creation logic

## 📁 Directory Structure

```
svelte-NET-Test/
├── MySvelteApp.Client/          # SvelteKit Frontend
├── MySvelteApp.Server/          # .NET Backend
├── CLAUDE.md                    # Development guide for AI assistants
└── structure.md                 # This file
```

## 🎯 Key Architectural Principles

### 1. **Clean Architecture Layers** (Backend)
- **Domain**: Core business entities and rules
- **Application**: Use cases, services, DTOs
- **Infrastructure**: External concerns (database, APIs, file system)
- **Presentation**: API controllers and web interface

### 2. **Remote Functions Pattern** (Frontend ↔ Backend)
- **Query**: Read-only operations with automatic caching
- **Command**: Write operations with optimistic updates
- **Form**: Form submissions with validation
- **Prerender**: Static data for build-time generation

### 3. **Type Safety First**
- All types match Prisma schema
- Zod schemas generated from OpenAPI spec
- End-to-end type safety between client and server

## 🔧 Backend Code Organization (.NET)

### Domain Layer (`MySvelteApp.Server/Domain/`)
```csharp
// 📍 PUT: Business entities and core rules
MySvelteApp.Server/Domain/
├── Entities/           # Core business entities
│   ├── User.cs
│   ├── Product.cs
│   └── Order.cs
└── ValueObjects/       # Domain value objects
    ├── Email.cs
    └── Money.cs
```

**Guidelines:**
- ✅ Pure business logic, no external dependencies
- ✅ Entities should encapsulate business rules
- ✅ Value objects should be immutable
- ✅ No database or UI concerns

### Application Layer (`MySvelteApp.Server/Application/`)
```csharp
// 📍 PUT: Use cases and business logic interfaces
MySvelteApp.Server/Application/
├── FeatureName/
│   ├── DTOs/              # Data Transfer Objects
│   │   ├── RequestDto.cs
│   │   └── ResponseDto.cs
│   ├── IFeatureService.cs # Service interfaces
│   └── Commands/          # CQRS command objects
└── Common/
    └── Exceptions/        # Custom business exceptions
```

**Guidelines:**
- ✅ Service interfaces define contracts
- ✅ DTOs for data transfer between layers
- ✅ No infrastructure dependencies
- ✅ Business rules and validation

### Infrastructure Layer (`MySvelteApp.Server/Infrastructure/`)
```csharp
// 📍 PUT: External implementations
MySvelteApp.Server/Infrastructure/
├── Persistence/          # Database implementations
│   ├── AppDbContext.cs
│   └── Repositories/
│       └── FeatureRepository.cs
├── External/             # External API clients
│   └── ExternalApiService.cs
├── Authentication/       # Auth implementations
└── Security/            # Security utilities
```

**Guidelines:**
- ✅ Concrete implementations of Application interfaces
- ✅ Database operations via Entity Framework
- ✅ External API integrations
- ✅ Infrastructure concerns (logging, caching, etc.)

### Presentation Layer (`MySvelteApp.Server/Presentation/`)
```csharp
// 📍 PUT: API controllers and models
MySvelteApp.Server/Presentation/
├── Controllers/
│   └── FeatureController.cs
└── Models/
    └── ApiModels.cs
```

**Guidelines:**
- ✅ Minimal logic, delegate to Application services
- ✅ HTTP concerns (routing, status codes, serialization)
- ✅ Input validation and error handling
- ✅ API documentation attributes

## 🎨 Frontend Code Organization (SvelteKit)

### Remote Functions (`MySvelteApp.Client/src/routes/**/feature.remote.ts`)
```typescript
// 📍 PUT: Server-side data operations
import { query, command, form } from '$app/server';
import { getRequestEvent } from '$app/server';

// Query functions for reading data
export const getData = query(async () => {
    const { locals } = getRequestEvent();
    // Access to server context, database, etc.
    return data;
});

// Command functions for mutations
export const updateData = command(async (params) => {
    // Side effects, database writes
    return result;
});

// Form functions for form submissions
export const submitForm = form(async (formData) => {
    // Form processing with validation
    return result;
});
```

#### Advanced Remote Function Features

**Multiple Form Actions with buttonProps:**
```typescript
// One form, multiple actions
<form {...loginForm}>
    <input name="username" />
    <input name="password" type="password" />

    <button>Login</button>
    <button {...registerForm.buttonProps}>Register</button>
</form>
```

**Custom Validation Error Handling:**
```typescript
// src/hooks.server.js
export function handleValidationError({ issues }) {
    return {
        message: 'Custom validation error message'
    };
}
```

**Server-Side Tracing (Experimental):**
```javascript
// svelte.config.js
export default {
    kit: {
        experimental: {
            tracing: { server: true },
            instrumentation: { server: true }
        }
    }
};
```

**Guidelines:**
- ✅ One `.remote.ts` file per feature/route
- ✅ Use appropriate function types (query/command/form)
- ✅ Access server context via `getRequestEvent()`
- ✅ Return serializable data only
- ✅ Use `buttonProps` for multiple actions in one form
- ✅ Implement `handleValidationError` for custom error messages
- ✅ Enable tracing for debugging server operations

### Page Components (`MySvelteApp.Client/src/routes/**/+page.svelte`)
```svelte
<!-- 📍 PUT: UI components with data fetching -->
<script>
  import { getData } from './data.remote';

  // Svelte 5 reactive state (replaces let)
  let count = $state(0);

  // Svelte 5 derived values (replaces $: reactive statements)
  const doubled = $derived(count * 2);

  // Async data with derived (correct Svelte 5 pattern)
  const data = $derived(await getData());

  // Manual state management for more control
  const dataQuery = getData();
</script>

<!-- Svelte 5 boundary with pending snippet (replaces {#await}) -->
<svelte:boundary>
  <DataDisplay {data} />

  {#snippet pending()}
    <LoadingSpinner />
  {/snippet}

  {#snippet error(err)}
    <ErrorDisplay {err} />
  {/snippet}
</svelte:boundary>

<!-- Manual state management -->
{#if dataQuery.loading}
  <LoadingSpinner />
{:else if dataQuery.error}
  <ErrorDisplay {dataQuery.error} />
{:else}
  <DataDisplay {dataQuery.current} />
{/if}
```

#### Svelte 5 Best Practices for Async Operations

**Reactive State with $state (Replaces `let`):**
```svelte
<script>
  // Svelte 5: Use $state() for reactive variables
  let count = $state(0);
  let user = $state(null);

  // Don't use: let count = 0; (not reactive in Svelte 5)
</script>
```

**Derived Values with $derived (Replaces `$:` statements):**
```svelte
<script>
  let count = $state(0);
  let items = $state([]);

  // Svelte 5: Use $derived() for computed values
  const doubled = $derived(count * 2);
  const total = $derived(items.length);

  // Don't use: $: doubled = count * 2; (legacy syntax)
</script>
```

**Async Data with $derived:**
```svelte
<script>
  import { getData } from './data.remote';

  // Correct Svelte 5 pattern: Use $derived with await
  const data = $derived(await getData());

  // Alternative: Direct await in template (simpler)
  // See template example below
</script>

<!-- Svelte 5: Use boundary for async operations -->
<svelte:boundary>
  <DataDisplay {data} />

  {#snippet pending()}
    <LoadingSpinner />
  {/snippet}
</svelte:boundary>
```

**Pure Load Functions (Recommended):**
```typescript
// MySvelteApp.Client/src/routes/+page.server.ts
export async function load({ fetch }) {
    const response = await fetch('/api/user');
    return {
        user: await response.json() // Return data, don't mutate state
    };
}
```

**Combining Server + Universal Load:**
```typescript
// +page.server.ts - serializable data
export async function load() {
    return { serverData: 'from server' };
}

// +page.ts - can use server data + client logic
export async function load({ data }) {
    return {
        serverData: data.serverData,
        clientData: 'from client'
    };
}
```

**Manual Invalidation:**
```typescript
// Control when data refreshes
export async function load({ fetch, depends }) {
    depends('app:custom-dependency'); // Manual refresh trigger

    const response = await fetch('/api/data');
    return { data: await response.json() };
}

// Later: invalidate('app:custom-dependency') to refresh
```

**Guidelines:**
- ✅ Use `$state()` for reactive variables (replaces `let`)
- ✅ Use `$derived()` for computed values (replaces `$:` statements)
- ✅ Use `$derived(await func())` for async reactive values
- ✅ Use `<svelte:boundary>` with `{#snippet pending()}` for loading states
- ✅ Keep load functions pure - return data, don't mutate state
- ✅ Combine server + universal loads for optimal SSR/hydration
- ✅ Use `depends()` and `invalidate()` for manual refresh control
- ✅ Handle async operations with proper error boundaries

### Route Layouts (`MySvelteApp.Client/src/routes/**/+layout.svelte`)
```svelte
<!-- 📍 PUT: Shared layout components -->
<script>
  // Layout-specific logic
</script>

<main>
  <slot />
</main>
```

### Server Layouts (`MySvelteApp.Client/src/routes/**/+layout.server.ts`)
```typescript
// 📍 PUT: Server-side layout logic
import { redirect } from '@sveltejs/kit';

export async function load({ cookies, url }) {
    // Authentication checks
    // Data preloading for layout
    return {
        user: authenticatedUser
    };
}
```

**Guidelines:**
- ✅ Authentication checks (not in hooks.server.ts)
- ✅ Layout-specific data loading
- ✅ Redirect logic for protected routes

## 🔗 API Client Generation

### Generated Files (`MySvelteApp.Client/src/api/schema/`)
```typescript
// 📍 AUTO-GENERATED: Do not edit directly
// Generated from OpenAPI spec via @hey-api/openapi-ts
├── sdk.gen.ts      # API client functions
├── types.gen.ts    # TypeScript types
├── zod.gen.ts      # Zod validation schemas
└── client.gen.ts   # HTTP client configuration
```

**Guidelines:**
- ✅ Never edit generated files directly
- ✅ Regenerate after API changes: `npm run generate-api-classes`
- ✅ Use generated types and schemas for type safety
- ✅ Custom API logic goes in `MySvelteApp.Client/src/api/` (not schema/)

### Custom API Logic (`MySvelteApp.Client/src/api/`)
```typescript
// 📍 PUT: Custom API extensions and utilities
import { generatedApiFunction } from '$api/schema/sdk.gen';

export const customApiCall = async (params) => {
    // Custom logic, error handling, retries
    const response = await generatedApiFunction(params);
    return transformedResponse;
};
```

## 🧩 Component Organization

### UI Components (`MySvelteApp.Client/src/lib/components/ui/`)
```svelte
<!-- 📍 PUT: Reusable UI components (shadcn/ui style) -->
<script lang="ts">
  import { cn } from '$lib/utils';

  // Svelte 5: Use $state() for internal reactive state
  let isLoading = $state(false);

  interface Props {
    variant?: 'default' | 'destructive';
    size?: 'sm' | 'md' | 'lg';
    loading?: boolean;
  }

  // Svelte 5: Use $derived() for computed values
  const classes = $derived(
    cn(baseClasses, variantClasses, sizeClasses, {
      'opacity-50 cursor-not-allowed': isLoading
    })
  );

  // Props are automatically reactive in Svelte 5
  $: if (loading !== undefined) {
    isLoading = loading;
  }
</script>

<button {classes} disabled={isLoading}>
  {#if isLoading}
    <Spinner />
  {:else}
    <slot />
  {/if}
</button>
```

**Guidelines:**
- ✅ Use shadcn/ui component structure
- ✅ Use `$state()` for internal reactive state (replaces `let`)
- ✅ Use `$derived()` for computed class names/styles
- ✅ Props are automatically reactive in Svelte 5
- ✅ TypeScript interfaces for props
- ✅ `cn()` utility for class merging
- ✅ Consistent with design system

### Feature Components (`MySvelteApp.Client/src/lib/components/`)
```svelte
<!-- 📍 PUT: Feature-specific components -->
<script>
  import { getData } from '$lib/data.remote';
  import DataDisplay from './ui/DataDisplay.svelte';

  // Svelte 5: Use $state() for local reactive state
  let refreshKey = $state(0);

  // Svelte 5: Use $derived() for computed values
  const data = $derived(await getData());

  // Function to refresh data
  function refresh() {
    refreshKey++;
  }
</script>

<!-- Svelte 5: Use boundary for async operations -->
<svelte:boundary>
  <DataDisplay {data} {refresh} />

  {#snippet pending()}
    <div class="flex items-center justify-center p-4">
      <Spinner class="w-6 h-6" />
      <span class="ml-2">Loading...</span>
    </div>
  {/snippet}

  {#snippet error(err)}
    <div class="text-red-500 p-4">
      <p>Error: {err.message}</p>
      <button onclick={refresh} class="mt-2 px-4 py-2 bg-red-500 text-white rounded">
        Try Again
      </button>
    </div>
  {/snippet}
</svelte:boundary>
```

**Guidelines:**
- ✅ Use `$state()` for reactive local state
- ✅ Use `$derived()` for computed reactive values
- ✅ Use `<svelte:boundary>` for async operations (replaces `{#await}`)
- ✅ Feature-specific business logic
- ✅ Composition of UI components
- ✅ Remote function integration
- ✅ Feature boundaries respected
- ✅ Proper error handling with retry functionality

## 🗂️ Utility Organization

### Core Utilities (`MySvelteApp.Client/src/lib/utils.ts`)
```typescript
// 📍 PUT: Core utility functions
import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

// Class name merging utility
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// Other core utilities...
```

### Feature Utilities (`MySvelteApp.Client/src/lib/feature/`)
```typescript
// 📍 PUT: Feature-specific utilities (create directory as needed)
export const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(amount);
};
```

## 🧪 Testing Organization

### Unit Tests (`MySvelteApp.Client/src/tests/unit/` or `MySvelteApp.Client/tests/unit/`)
```typescript
// 📍 PUT: Unit tests for utilities and components (create directory as needed)
import { describe, it, expect } from 'vitest';
import { formatCurrency } from '$lib/utils';

describe('formatCurrency', () => {
  it('formats USD correctly', () => {
    expect(formatCurrency(1234.56)).toBe('$1,234.56');
  });
});
```

### E2E Tests (`MySvelteApp.Client/e2e/`)
```typescript
// 📍 PUT: End-to-end tests
import { test, expect } from '@playwright/test';

test('user can login', async ({ page }) => {
  await page.goto('/login');
  // Test interactions...
});
```

### Server Tests (`MySvelteApp.Server/Tests/`)
```csharp
// 📍 PUT: .NET unit and integration tests
[TestClass]
public class AuthServiceTests
{
    [TestMethod]
    public async Task Login_ValidCredentials_ReturnsSuccess()
    {
        // Test implementation
    }
}
```

## 🔐 Authentication & Security

### Authentication Logic
```typescript
// 📍 PUT: In MySvelteApp.Client/src/routes/(auth)/auth.remote.ts
export const login = form(async (formData) => {
  // JWT token handling
  // Cookie management
  // Server-side validation
});
```

```csharp
// 📍 PUT: In MySvelteApp.Server/Infrastructure/Authentication/
// JWT token generation
// Password hashing (HMACSHA512)
// Authentication middleware
```

### Route Protection
```typescript
// 📍 PUT: In MySvelteApp.Client/src/routes/(app)/+layout.server.ts
export async function load({ cookies }) {
  const token = cookies.get('auth_token');
  if (!token) throw redirect(302, '/login');
  // Validate token and return user
}
```

## 📊 Data Flow Patterns

### Read Operations (Queries)
1. **Component** calls remote `query()` function
2. **Remote function** accesses server context via `getRequestEvent()`
3. **Service** (Application layer) orchestrates business logic
4. **Repository** (Infrastructure) executes database queries
5. **Data** flows back through layers to component

### Write Operations (Commands)
1. **Component** calls remote `command()` function
2. **Remote function** processes input and calls service
3. **Service** validates and orchestrates business logic
4. **Repository** executes database mutations
5. **Optimistic updates** via `.updates()` and `.withOverride()`

### Form Submissions
1. **Form component** uses `{...remoteForm}` spread
2. **Remote form function** validates input with Zod
3. **Service** processes validated data
4. **Repository** persists changes
5. **Success/error** handled automatically

## 🚀 Deployment & Configuration

### Environment Configuration
```typescript
// 📍 PUT: In MySvelteApp.Client/src/api/config.ts (or create MySvelteApp.Client/src/lib/config.ts)
export const config = {
  apiUrl: import.meta.env.VITE_API_URL,
  environment: import.meta.env.MODE
};
```

### SvelteKit Experimental Features (Already Enabled)
```javascript
// svelte.config.js - Already configured in your project
export default {
  kit: {
    experimental: {
      remoteFunctions: true,  // ✅ Enabled - Type-safe client-server communication
      async: true             // ✅ Enabled - Svelte 5 async/await syntax
    }
  },
  compilerOptions: {
    experimental: {
      async: true             // ✅ Enabled - Async compiler features
    }
  }
};
```

### Docker Configuration
```dockerfile
# 📍 PUT: In MySvelteApp.Client/Dockerfile
# Multi-stage build for SvelteKit
# Nginx configuration for static assets
```

## 📝 Development Workflow

### Adding New Features
1. **Backend First**: Implement domain entities and application services
2. **API Layer**: Create controllers and DTOs
3. **Client Generation**: Run `npm run generate-api-classes`
4. **Remote Functions**: Implement client-side remote functions
5. **Components**: Build UI components using remote functions
6. **Testing**: Add unit and E2E tests

### Code Generation
- API clients: `npm run generate-api-classes`
- Database migrations: `dotnet ef migrations add`
- Type checking: `npm run check`

## 🔍 Common Patterns & Anti-Patterns

### ✅ DOs
- Use alias imports (`$lib`, `$api`) instead of relative paths
- Prefer remote functions over manual fetch calls
- Keep business logic in Application/Infrastructure layers
- Use generated types and Zod schemas
- Follow Clean Architecture principles
- Use Svelte 5 runes and async syntax

### ❌ DON'Ts
- Don't put database calls in controllers
- Don't use `any` types (full type safety required)
- Don't bypass service/repository layers
- Don't put authentication logic in hooks.server.ts
- Don't edit generated API client files
- Don't use relative imports in components

## 🎯 Quick Reference

| What | Where | Example |
|------|-------|---------|
| Business entities | `MySvelteApp.Server/Domain/Entities/` | `User.cs` |
| Service interfaces | `MySvelteApp.Server/Application/Feature/` | `IUserService.cs` |
| Service implementations | `MySvelteApp.Server/Infrastructure/Authentication/` | `AuthService.cs` |
| API controllers | `MySvelteApp.Server/Presentation/Controllers/` | `UserController.cs` |
| **Remote functions** | `MySvelteApp.Client/src/routes/**/feature.remote.ts` | `user.remote.ts` |
| **Page components** | `MySvelteApp.Client/src/routes/**/+page.svelte` | `+page.svelte` |
| **Server layouts** | `MySvelteApp.Client/src/routes/**/+layout.server.ts` | Auth checks |
| UI components | `MySvelteApp.Client/src/lib/components/ui/` | `Button.svelte` |
| Feature components | `MySvelteApp.Client/src/lib/components/` | `UserProfile.svelte` |
| Utilities | `MySvelteApp.Client/src/lib/utils.ts` | `cn()` function |
| API client | `MySvelteApp.Client/src/api/schema/sdk.gen.ts` | Generated functions |
| **Load functions** | `MySvelteApp.Client/src/routes/**/+page.server.ts` | Pure data loading |
| **Validation errors** | `MySvelteApp.Client/src/hooks.server.js` | `handleValidationError` |

## 🔧 Key Experimental Features (Already Enabled)

- **Remote Functions**: Type-safe client-server communication
- **Async Svelte**: `$derived(await func())` and `<svelte:boundary>` support
- **Svelte 5 Runes**: `$state()`, `$derived()`, and other modern reactivity
- **buttonProps**: Multiple form actions in one form
- **Tracing**: Server-side operation debugging (optional)

## ⚡ Svelte 5 Migration Notes

**From Svelte 4 to Svelte 5:**
- `let count = 0` → `let count = $state(0)`
- `$: doubled = count * 2` → `const doubled = $derived(count * 2)`
- `{#await promise}` → `<svelte:boundary>{#snippet pending()}`
- Reactive statements → `$derived()` expressions
- Top-level variables → Explicit `$state()` wrapping

This structure ensures maintainable, type-safe, and scalable code that follows modern full-stack development best practices.
