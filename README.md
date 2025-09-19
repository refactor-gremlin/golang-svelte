# Svelte-NET Template

A modern, type-safe full-stack template featuring **ASP.NET Core 9.0** as Backend for Frontend (BFF) with **SvelteKit 2.22** (Svelte 5) frontend, implementing Clean Architecture and DDD-inspired patterns.

## 📦 What is this template?

This template demonstrates best practices for building full-stack applications using:
- **ASP.NET Core** as a type-safe Backend for Frontend (BFF)
- **SvelteKit** with modern Svelte 5 features for the frontend
- **Clean Architecture** with DDD-inspired patterns
- **Type-safe communication** between frontend and backend
- **Modern development workflow** with Docker and CI/CD

Perfect for developers looking to build scalable, maintainable full-stack applications.

## ✨ What Makes This Template Special

- **🛡️ BFF Architecture**: ASP.NET Core as dedicated backend for frontend
- **🔒 Type Safety**: End-to-end type safety with generated API clients
- **🎯 Modern Svelte 5**: Latest Svelte features with runes and async patterns
- **🏛️ Clean Architecture**: DDD-inspired patterns for maintainability
- **🐳 Docker Ready**: Complete containerization for development and production
- **⚡ Remote Functions**: Experimental SvelteKit feature for seamless client-server communication
- **📊 Observability**: Built-in monitoring with Grafana, Loki, and Prometheus
- **🔄 CI/CD**: GitHub Actions pipeline with automated testing and deployment

## 🏗️ Architecture

This application follows a **DDD-inspired Clean Architecture** with clear separation of concerns:

```
📁 Project Structure
├── MySvelteApp.Client/     # SvelteKit Frontend (Svelte 5)
├── MySvelteApp.Server/     # .NET 9.0 Backend (Clean Architecture)
├── CLAUDE.md              # AI Assistant Documentation
├── structure.md           # Code Organization Guide
└── README.md              # This file
```

### Backend Architecture (.NET 9.0)
- **Domain Layer**: Core business entities and rules
- **Application Layer**: Use cases, services, and DTOs
- **Infrastructure Layer**: External implementations (database, APIs)
- **Presentation Layer**: API controllers and web interface

### Frontend Architecture (SvelteKit 5)
- **Remote Functions**: Type-safe client-server communication
- **Component Organization**: Reusable UI components with shadcn/ui
- **Route-based Architecture**: File-system routing with layouts
- **Modern Reactivity**: Svelte 5 runes (`$state`, `$derived`)

## 🛠️ Technology Stack

### Frontend
- **Framework**: SvelteKit 2.22.0 with Svelte 5.0.0
- **Language**: TypeScript
- **Build Tool**: Vite 7.0.4
- **Styling**: Tailwind CSS 4.0
- **UI Components**: shadcn/ui + Tailwind
- **Testing**: Vitest (unit) + Playwright (E2E)
- **API Client**: Generated from OpenAPI spec

### Backend
- **Framework**: .NET 9.0 Web API
- **Architecture**: Clean Architecture + DDD patterns
- **Database**: Entity Framework Core (in-memory for dev)
- **Authentication**: JWT with HMACSHA512
- **Testing**: xUnit with FluentAssertions
- **API Documentation**: Swagger/OpenAPI

### DevOps & Tools
- **Containerization**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Code Quality**: ESLint + Prettier + TypeScript
- **Observability**: Grafana + Loki + Prometheus
- **Version Control**: Git with conventional commits

## 🚀 Quick Start

### Prerequisites
- **Node.js** 20+ (with npm)
- **.NET 9.0** SDK
- **Docker** & Docker Compose (for full-stack development)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/refactor-gremlin/svelte-NET-Test.git
   cd svelte-NET-Test
   ```

2. **Install dependencies**
   ```bash
   # Install backend dependencies
   dotnet restore

   # Install frontend dependencies
   npm ci --prefix MySvelteApp.Client
   ```

3. **Start development servers**
   ```bash
   # Start both frontend and backend concurrently
   npm run dev

   # Or start individually:
   # Backend (port 7216)
   dotnet run --project MySvelteApp.Server

   # Frontend (port 5173)
   npm run dev --prefix MySvelteApp.Client
   ```

4. **Open your browser**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:7216
   - API Documentation: http://localhost:7216/swagger

### Docker Development

For full-stack development with all services:

```bash
# Start all services
npm run docker:dev

# Or use docker-compose directly
docker-compose -f docker-compose.dev.yml up
```

## 📚 Available Scripts

### Root Level Scripts
```bash
npm run dev          # Start both client and server concurrently
npm run build        # Build for production
npm run docker:dev   # Start development containers
npm run docker:prod  # Start production containers
```

### Client Scripts (MySvelteApp.Client/)
```bash
npm run dev          # Start dev server (port 5173)
npm run build        # Production build
npm run check        # TypeScript type checking
npm run lint         # ESLint code quality
npm run format       # Prettier code formatting
npm run test:unit    # Vitest unit tests
npm run test:e2e     # Playwright E2E tests
npm run generate-api-classes  # Regenerate API client
```

### Server Scripts (MySvelteApp.Server/)
```bash
dotnet run           # Start development server (port 7216)
dotnet build         # Build the solution
dotnet test          # Run unit tests
dotnet restore       # Restore NuGet packages
```

## 🔐 Authentication & BFF Pattern

This template implements the **Backend for Frontend (BFF) pattern** with ASP.NET Core serving as a dedicated API layer for the Svelte frontend.

### Authentication Features
- **JWT-based authentication** with secure password hashing (HMACSHA512)
- **Registration/Login** with business rule validation
- **HTTP-only cookies** for secure token storage
- **Route protection** with automatic redirects
- **Session management** with proper logout handling

### BFF Benefits
- **Type Safety**: End-to-end type safety between frontend and backend
- **Simplified Frontend**: Frontend focuses on UI/UX, backend handles business logic
- **API Optimization**: Backend can optimize responses for specific frontend needs
- **Security**: Centralized authentication and authorization
- **Performance**: Reduced round trips and optimized data fetching

### Authentication Flow
1. User interacts with Svelte frontend
2. Frontend calls remote functions (type-safe API calls)
3. ASP.NET Core BFF validates requests and business rules
4. JWT tokens managed securely via HTTP-only cookies
5. Protected routes automatically redirect unauthenticated users

## 🧪 Testing

### Backend Tests
```bash
# Run all backend tests
dotnet test MySvelteApp.Server/Tests/

# Run with code coverage
dotnet test MySvelteApp.Server/Tests/ --collect:"XPlat Code Coverage"
```

### Frontend Tests
```bash
# Unit tests
npm run test:unit --prefix MySvelteApp.Client

# End-to-end tests
npm run test:e2e --prefix MySvelteApp.Client

# Type checking
npm run check --prefix MySvelteApp.Client
```

## 🚀 Deployment

### Production Build

1. **Build the application**
   ```bash
   npm run build
   ```

2. **Using Docker**
   ```bash
   # Production containers
   npm run docker:prod

   # Or manually
   docker-compose -f docker-compose.yml up -d
   ```

3. **Manual deployment**
   - Backend: Deploy .NET application to your server
   - Frontend: Deploy static files to CDN/web server
   - Database: Configure production database connection

### Environment Variables

Create `.env` files for different environments:

```bash
# Frontend (.env)
VITE_API_URL=http://localhost:7216

# For production, update to your deployed API URL
# VITE_API_URL=https://your-api-domain.com
```

```json
// Backend (appsettings.Production.json)
{
  "ConnectionStrings": {
    "DefaultConnection": "Server=your-server;Database=your-db;..."
  },
  "Jwt": {
    "Secret": "your-256-bit-secret-key-here",
    "Issuer": "your-app-name",
    "Audience": "your-app-name"
  }
}
```

## 📖 API Documentation

The backend provides comprehensive API documentation:

- **Swagger UI**: http://localhost:7216/swagger (development)
- **OpenAPI Spec**: Automatically generated from controllers
- **TypeScript Client**: Generated from OpenAPI spec for type safety

### Key Endpoints

- `POST /Auth/login` - User authentication
- `POST /Auth/register` - User registration
- `GET /RandomPokemon` - Get random Pokemon data
- `GET /WeatherForecast` - Sample weather data

## 🔧 Development Workflow

### Code Quality
```bash
# Run all quality checks
npm run lint --prefix MySvelteApp.Client
npm run check --prefix MySvelteApp.Client
dotnet test MySvelteApp.Server/Tests/
```

### Adding New Features

1. **Backend First**: Implement domain logic and API
2. **Generate API Client**: `npm run generate-api-classes`
3. **Implement Frontend**: Use generated types and client
4. **Add Tests**: Unit tests for both backend and frontend
5. **Update Documentation**: Keep README and structure.md current

### Remote Functions Pattern

This project uses SvelteKit's experimental remote functions for type-safe client-server communication:

```typescript
// Remote function (server-side)
export const getData = query(async () => {
    // Server-side logic
    return data;
});

// Component usage (client-side)
const data = $derived(await getData());
```

## 🤝 Contributing

### Development Setup
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes with proper tests
4. Ensure all checks pass: `npm run lint && npm run check && dotnet test`
5. Commit with conventional format: `git commit -m "feat: add new feature"`
6. Push and create a pull request

### Code Standards
- **Backend**: C# 12 with nullable types, guard clauses
- **Frontend**: TypeScript with strict mode, Svelte 5 patterns
- **Commits**: Conventional commits format
- **Tests**: 80%+ code coverage target
- **Documentation**: Update README and structure.md for changes

## 📄 Project Structure & Customization

### Template Structure
```
svelte-NET-Test/
├── MySvelteApp.Client/     # SvelteKit frontend (Svelte 5)
├── MySvelteApp.Server/     # ASP.NET Core BFF (Clean Architecture)
├── CLAUDE.md              # AI assistant development guide
├── structure.md           # Complete architecture documentation
├── docker-compose.yml     # Production deployment
├── docker-compose.dev.yml # Development environment
└── .github/               # CI/CD workflows
```

### Customization Guide

**For Your Project:**
1. **Rename Projects**: Update `MySvelteApp.Client` and `MySvelteApp.Server` to your app names
2. **Update Package Names**: Change namespaces and package IDs
3. **Configure Database**: Replace in-memory DB with your preferred database
4. **Environment Variables**: Set up your production environment variables
5. **Domain/Branding**: Update API URLs, app names, and branding

**Key Documentation:**
- **[structure.md](./structure.md)** - Complete code organization guide
- **[CLAUDE.md](./CLAUDE.md)** - AI assistant development guide
- **[.github/README.md](./.github/README.md)** - CI/CD pipeline documentation

### Architecture Decisions
- **Clean Architecture** for maintainability and testability
- **BFF Pattern** for optimal frontend-backend communication
- **Remote Functions** for type-safe client-server calls
- **DDD Patterns** for complex business logic organization
- **Docker** for consistent development and deployment environments

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋 Support

- **Issues**: [GitHub Issues](https://github.com/refactor-gremlin/svelte-NET-Test/issues)
- **Discussions**: [GitHub Discussions](https://github.com/refactor-gremlin/svelte-NET-Test/discussions)
- **Documentation**: See [structure.md](./structure.md) for detailed architecture info
- **Contributing**: See [CONTRIBUTING.md](./CONTRIBUTING.md) for contribution guidelines

---

**Professional full-stack template built with ❤️**

*Frontend: SvelteKit 5 + TypeScript + Tailwind CSS*
*Backend: ASP.NET Core 9.0 + Clean Architecture + DDD patterns*
*BFF Pattern: Type-safe communication between frontend and backend*
*DevOps: Docker + GitHub Actions + Observability*

## 🎯 Perfect For

- **🚀 Startups** needing rapid development with production-ready architecture
- **🏢 Enterprises** requiring scalable, maintainable full-stack solutions
- **👥 Development Teams** wanting type-safe frontend-backend communication
- **🎓 Educational Projects** demonstrating modern web development patterns
- **🔄 Migration Projects** moving from traditional MVC to modern SPA architectures

## 🏃‍♂️ Quick Start Your Project

```bash
# 1. Clone this template
git clone https://github.com/refactor-gremlin/svelte-NET-Test.git my-awesome-app
cd my-awesome-app

# 2. Customize for your needs
# - Rename projects and namespaces
# - Configure your database
# - Update environment variables
# - Add your business logic

# 3. Start developing!
npm run dev
```

**This template gives you a production-ready foundation so you can focus on building your application's unique features!** 🌟
