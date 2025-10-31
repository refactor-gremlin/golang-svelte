"""Pokemon API routes."""

from contextlib import asynccontextmanager

from fastapi import APIRouter, HTTPException, Request, status

from app.modules.pokemon.app.service import ExternalAPIError, PokemonService

router = APIRouter(tags=["pokemon"])


@asynccontextmanager
async def lifespan(app):
    """Manage Pokemon service lifecycle."""
    service = PokemonService()
    app.state.pokemon_service = service
    yield
    await service.close()


def get_pokemon_service(request: Request) -> PokemonService:
    """Get Pokemon service instance from app state."""
    return request.app.state.pokemon_service


@router.get("/RandomPokemon")
async def get_random_pokemon(request: Request):
    """Get a random Pokemon from external API."""
    service = get_pokemon_service(request)
    try:
        pokemon = await service.get_random_pokemon()
        return pokemon
    except ExternalAPIError:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="Pokemon service temporarily unavailable",
        )
    except Exception:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to get random Pokemon",
        )

