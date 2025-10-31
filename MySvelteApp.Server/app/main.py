"""Main FastAPI application entry point."""

from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import RedirectResponse, JSONResponse

from app.config import Settings, load_settings
from app.database import Database
from app.modules.auth.api.routes import router as auth_router
from app.modules.pokemon.api.routes import lifespan as pokemon_lifespan, router as pokemon_router
from app.platform.logging import get_logger
from app.platform.tracing import TracingProvider, instrument_fastapi

logger = get_logger()
_settings: Settings | None = None


def get_settings() -> Settings:
    """Get application settings."""
    global _settings
    if _settings is None:
        _settings = load_settings()
    return _settings


@asynccontextmanager
async def app_lifespan(app: FastAPI):
    """Manage application lifecycle."""
    settings = get_settings()
    
    # Initialize database
    database = Database(settings)
    database.create_tables()
    app.state.database = database
    
    # Initialize tracing
    tracing_provider = TracingProvider(
        service_name=settings.service_name,
        service_version=settings.service_version,
        environment=settings.environment,
    )
    tracing_provider.setup()
    app.state.tracing_provider = tracing_provider
    
    # Setup Pokemon service lifecycle
    async with pokemon_lifespan(app):
        yield
    
    # Shutdown tracing
    if hasattr(app.state, "tracing_provider"):
        app.state.tracing_provider.shutdown()
    
    # Shutdown database connections
    database.engine.dispose()


def create_app() -> FastAPI:
    """Create and configure FastAPI application."""
    settings = get_settings()
    
    app = FastAPI(
        title="MySvelteApp Server API",
        description="Python implementation of the MySvelteApp backend.",
        version="1.0.0",
        lifespan=app_lifespan,
    )
    
    # CORS middleware
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],  # Configure appropriately for production
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )
    
    # OpenTelemetry instrumentation
    instrument_fastapi(app)
    
    # Register routes
    app.include_router(auth_router)
    app.include_router(pokemon_router)
    
    # Compatibility routes for Swagger endpoints expected by client
    @app.get("/swagger/index.html")
    async def swagger_redirect():
        """Redirect to FastAPI docs for compatibility with Go backend."""
        return RedirectResponse(url="/docs")
    
    @app.get("/swagger/v1/swagger.json")
    async def swagger_json():
        """Return OpenAPI JSON at the path expected by client tools."""
        return JSONResponse(content=app.openapi())
    
    @app.get("/swagger/swagger.json")
    async def swagger_json_alt():
        """Alternative Swagger JSON path."""
        return JSONResponse(content=app.openapi())
    
    return app


app = create_app()


if __name__ == "__main__":
    import uvicorn
    
    settings = get_settings()
    logger.info(f"Starting server on port {settings.port}")
    logger.info(f"Service: {settings.service_name} v{settings.service_version}")
    
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=int(settings.port),
        reload=True,
    )

