# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Important Instructions

**Always use Context7 for unfamiliar topics**: When encountering concepts you might not be familiar with, such as:
- shadcn svelte components and patterns
- Svelte component syntax and features
- Remote functions and async Svelte
- SvelteKit experimental features
- Any modern JavaScript/TypeScript frameworks or libraries

Use the Context7 tools to research and get up-to-date documentation before providing guidance or making changes.

## Development Commands

### Common Development Tasks
- `npm run dev` - Start both client and server in development mode (concurrently)
- `npm run build` - Build for production (client + server)
- `npm run test` - Run all tests (root level, currently not configured)
- `npm run docker:dev` - Start development containers
- `npm run docker:prod` - Start production containers

### Client-Specific Commands (run from `/MySvelteApp.Client/`)
- `npm run dev` - Start SvelteKit dev server (port 5173)
- `npm run build` - Build SvelteKit for production
- `npm run check` - TypeScript type checking
- `npm run lint` - Run ESLint and Prettier checks
- `npm run format` - Format code with Prettier
- `npm run test:unit` - Run Vitest unit tests
- `npm run test:e2e` - Run Playwright E2E tests
- `npm run generate-api-classes` - Generate TypeScript API client from OpenAPI spec

### Server-Specific Commands (run from `/MySvelteApp.Server/`)
- `dotnet run` - Start .NET Web API (port 7216)
- `dotnet build` - Build .NET project
- `dotnet test` - Run .NET tests

## Architecture Overview

This is a full-stack web application with SvelteKit frontend and .NET 9.0 backend:

### Frontend (SvelteKit)
- **Location**: `/MySvelteApp.Client/`
- **Framework**: SvelteKit 2.22.0 with Svelte 5.0.0
- **Build**: Vite 7.0.4
- **Styling**: Tailwind CSS 4.0 with shadcn/ui components
- **Testing**: Vitest (unit) + Playwright (E2E)
- **API Client**: Generated from OpenAPI spec using `@hey-api/openapi-ts`

### Backend (.NET Web API)
- **Location**: `/MySvelteApp.Server/`
- **Framework**: .NET 9.0 with ASP.NET Core Web API
- **Database**: Entity Framework Core with in-memory database
- **Authentication**: JWT-based with HMACSHA512 password hashing
- **API Documentation**: Swagger/OpenAPI integration

## Key Project Structure

### Routes
- `(auth)/` - Authentication routes (login, register)
- `(app)/` - Protected application routes (dashboard, pokemon)
- API endpoints: `/Auth/*`, `/WeatherForecast`, `/RandomPokemon`, `/TestAuth`

### Component Organization
- `/src/lib/components/ui/` - shadcn/ui components (30+ components)
- `/src/lib/components/` - Custom components (sidebar, navigation, etc.)
- `/src/api/` - Generated API client code
- `/src/hooks/` - SvelteKit server hooks

### Server Structure
- `/Controllers/` - API controllers with authentication
- `/Models/` - Data models and DTOs
- `/Services/` - Business logic services
- `/Data/` - Database context

## Development Environment

### Ports
- **Client**: 5173 (Vite dev server)
- **API**: 7216 (ASP.NET Core)
- **Grafana**: 3000 (Observability)
- **Loki**: 3100 (Log aggregation)

### Dev Container
- Pre-configured with .NET 9.0 and Node.js
- VS Code extensions and debugging setup
- Integrated observability stack

## Authentication System

- JWT-based authentication with secure password hashing
- Global authorization policy
- Public endpoints marked with `[AllowAnonymous]`
- Protected routes in `(app)` layout group
- Server-side authentication in `/src/routes/(auth)/+layout.server.ts`

## Build & Deployment

### Production Build Process
1. Build SvelteKit client (`npm run build` in client directory)
2. Build .NET server with static assets
3. Multi-stage Docker image with Nginx reverse proxy
4. Container orchestration with Docker Compose

### API Client Generation
- Uses `@hey-api/openapi-ts` to generate TypeScript client from OpenAPI spec
- Configuration in `/MySvelteApp.Client/openapi-ts.config.ts`
- Run `npm run generate-api-classes` after API changes
- Generates Zod schemas for runtime validation
- Includes fetch client and SDK utilities

## Remote Functions and Async Svelte

The project uses SvelteKit's experimental remote functions and async Svelte features for type-safe client-server communication and effective loading states.

### Remote Functions Configuration

**Enabled in `svelte.config.js`:**
```javascript
kit: {
  experimental: {
    remoteFunctions: true
  }
},
compilerOptions: {
  experimental: {
    async: true
  }
}
```

### Remote Functions Types

SvelteKit remote functions come in four flavors:

#### 1. Query Functions (`$app/server`)
For reading dynamic data from the server:
```javascript
// src/lib/data.remote.js
import { query } from '$app/server';
import * as v from 'valibot';

export const getPosts = query(async () => {
  const posts = await db.sql`SELECT * FROM posts`;
  return posts;
});

export const getPost = query(v.string(), async (slug) => {
  const [post] = await db.sql`SELECT * FROM posts WHERE slug = ${slug}`;
  return post;
});
```

#### 2. Command Functions
For server-side data manipulation:
```javascript
export const addLike = command(v.string(), async (postId) => {
  await db.sql`UPDATE posts SET likes = likes + 1 WHERE id = ${postId}`;
});
```

#### 3. Form Functions
For handling form submissions:
```javascript
export const createPost = form(async (data) => {
  const title = data.get('title');
  const content = data.get('content');
  
  await db.sql`INSERT INTO posts (title, content) VALUES (${title}, ${content})`;
  
  return { success: true };
});
```

#### 4. Prerender Functions
For static data at build time:
```javascript
export const getStaticData = prerender(async () => {
  return dataThatChangesInfrequently;
});
```

### Async Svelte Component Usage

#### Direct Await Usage
```svelte
<script>
  import { getPosts } from './data.remote';
</script>

{#await getPosts() as posts}
  <ul>
    {#each posts as post}
      <li>{post.title}</li>
    {/each}
  </ul>
{:catch error}
  <p>Error: {error.message}</p>
{/await}
```

#### Manual State Management
```svelte
<script>
  import { getPosts } from './data.remote';
  
  const postsQuery = getPosts();
</script>

{#if postsQuery.loading}
  <p>Loading...</p>
{:else if postsQuery.error}
  <p>Error: {postsQuery.error.message}</p>
{:else}
  <ul>
    {#each postsQuery.current as post}
      <li>{post.title}</li>
    {/each}
  </ul>
{/if}
```

### Form Integration

#### Basic Form Usage
```svelte
<script>
  import { createPost } from './data.remote';
</script>

<form {...createPost}>
  <input name="title" />
  <textarea name="content"></textarea>
  <button>Create Post</button>
</form>
```

#### Enhanced Form with Custom Logic
```svelte
<script>
  import { createPost } from './data.remote';
  import { showToast } from '$lib/toast';
</script>

<form {...createPost.enhance(async ({ form, data, submit }) => {
  try {
    await submit();
    form.reset();
    showToast('Post created successfully!');
  } catch (error) {
    showToast('Failed to create post');
  }
})}>
  <!-- form fields -->
</form>
```

### Optimistic Updates and Query Refreshing

#### Optimistic Updates
```svelte
<script>
  import { getLikes, addLike } from './likes.remote';
  
  async function handleLike() {
    await addLike(postId).updates(
      getLikes(postId).withOverride(current => current + 1)
    );
  }
</script>
```

#### Query Refreshing
```svelte
<script>
  import { getPosts } from './data.remote';
  
  function refreshPosts() {
    getPosts().refresh();
  }
</script>
```

### Type Safety and Validation

Remote functions automatically provide type safety and can include validation:

```javascript
import * as v from 'valibot';

export const getUser = query(
  v.string(), // Validates input
  async (userId) => {
    // Function implementation
    return user;
  }
);
```

### Access to Server Context

Remote functions can access server-only modules and context:

```javascript
import { getRequestEvent } from '$app/server';

export const getProfile = query(async () => {
  const { locals, cookies } = getRequestEvent();
  const user = await locals.getUser();
  return user;
});
```

### Benefits

- **Type Safety**: End-to-end type safety between client and server
- **Server-Only Access**: Safe access to databases, environment variables, and other server resources
- **Simplified Data Flow**: Direct function calls instead of manual fetch requests
- **Loading States**: Built-in loading and error state management
- **Optimistic Updates**: Support for immediate UI updates
- **Form Integration**: Seamless form handling with progressive enhancement

## Component Libraries

### shadcn/ui Components
- 30+ UI components including button, input, card, dialog, etc.
- Located in `/src/lib/components/ui/`
- Integrated with Tailwind CSS and Svelte 5

### Additional Libraries
- **bits-ui**: Advanced UI components
- **formsnap**: Form handling utilities
- **svelte-sonner**: Toast notifications
- **mode-watcher**: Dark/light theme support

## Important Configuration Files

- `svelte.config.js` - SvelteKit configuration with remote functions enabled
- `vite.config.ts` - Vite build configuration with testing setup
- `tsconfig.json` - TypeScript strict mode
- `openapi-ts.config.ts` - API client generation configuration
- `docker-compose.yml` - Production containers
- `docker-compose.dev.yml` - Development containers

### Example Implementations

The project includes example implementations of remote functions:

- **Weather Dashboard** (`/src/lib/components/WeatherDashboard.svelte`) - Complete example using query, command, and form functions
- **Weather Remote Functions** (`/src/lib/weather.remote.js`) - Server-side functions for weather data management

These examples demonstrate:
- Async/await usage in Svelte components
- Loading and error state management
- Form handling with progressive enhancement
- Optimistic updates and query refreshing
- Type-safe client-server communication

## CORS Configuration

Configured for multiple frontend origins:
- `http://localhost:5173` (development)
- `http://localhost:3000` (production preview)
- `http://web:3000` (Docker production)

## Database Configuration

Currently uses in-memory Entity Framework Core database for development. For production, consider migrating to PostgreSQL or another persistent database.

## Testing Configuration

- **Unit Tests**: Vitest with browser testing support
- **E2E Tests**: Playwright with web server integration
- **Client Testing**: Browser-based tests for Svelte components
- **Server Testing**: Node.js environment for utility functions

## Code Quality & Linting

- **ESLint**: With Svelte plugin and TypeScript support
- **Prettier**: With Svelte and Tailwind plugins
- **TypeScript**: Strict mode enabled
- **Testing**: Comprehensive unit and E2E test setup
- **ast-grep**: Structural search and lint tool for catching Svelte syntax issues

### Important: Svelte 5 Validation

**After editing any Svelte file (.svelte), always run:**
```bash
ast-grep scan
```

This command helps catch common Svelte 4 to 5 migration issues and ensures proper Svelte 5 syntax. The ast-grep tool performs structural pattern matching to identify:

- Deprecated Svelte 4 syntax patterns
- Incorrect reactive statements
- Invalid component structures
- Type safety issues specific to Svelte 5
- Proper usage of new Svelte 5 features

This is crucial for maintaining code quality and ensuring compatibility with Svelte 5's new features and syntax requirements.