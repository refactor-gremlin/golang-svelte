"""Pokemon service for external API integration."""

import asyncio
import json
import random
import time
from typing import Optional

import httpx

POKEMON_API_BASE_URL = "https://pokeapi.co/api/v2/pokemon/"
POKEMON_COUNT_URL = "https://pokeapi.co/api/v2/pokemon-species/?limit=0"
POKEMON_COUNT_CACHE_TTL = 10 * 60  # 10 minutes in seconds


class ExternalAPIError(Exception):
    """Indicates the upstream Pokemon API could not satisfy the request."""

    pass


class CountCache:
    """Cache for Pokemon count."""

    def __init__(self):
        """Initialize cache."""
        self.value: int = 0
        self.expires_at: float = 0.0


class PokemonService:
    """Service for Pokemon operations."""

    def __init__(self, http_client: Optional[httpx.AsyncClient] = None):
        """Initialize Pokemon service."""
        self.http_client = http_client or httpx.AsyncClient(timeout=30.0)
        self.cache = CountCache()
        self._cache_lock = asyncio.Lock()

    async def get_random_pokemon(self) -> dict:
        """Retrieve a random Pokemon from external API."""
        count = await self._get_pokemon_count()
        if count <= 0:
            raise ExternalAPIError("invalid pokemon count")

        random_pokemon_id = random.randint(1, count)
        pokemon_url = f"{POKEMON_API_BASE_URL}{random_pokemon_id}"

        pokemon_data = await self._fetch_pokemon(pokemon_url)

        types = [t["type"]["name"] for t in pokemon_data.get("types", [])]
        type_str = ", ".join(types)

        sprites = pokemon_data.get("sprites", {})
        front_default = sprites.get("front_default")

        return {
            "name": pokemon_data.get("name"),
            "type": type_str if types else None,
            "image": front_default,
        }

    async def _get_pokemon_count(self) -> int:
        """Get total Pokemon count with caching."""
        async with self._cache_lock:
            now = time.time()
            if self.cache.value > 0 and self.cache.expires_at > now:
                return self.cache.value

        try:
            response = await self.http_client.get(POKEMON_COUNT_URL)
            response.raise_for_status()

            data = response.json()
            count = data.get("count", 0)

            if count <= 0:
                raise ExternalAPIError("received non-positive count")

            async with self._cache_lock:
                self.cache.value = count
                self.cache.expires_at = now + POKEMON_COUNT_CACHE_TTL

            return count
        except httpx.HTTPError as e:
            raise ExternalAPIError(f"retrieve count: {e}") from e
        except Exception as e:
            raise ExternalAPIError(f"decode count response: {e}") from e

    async def _fetch_pokemon(self, url: str) -> dict:
        """Fetch Pokemon data from API."""
        try:
            response = await self.http_client.get(url)
            response.raise_for_status()
            return response.json()
        except httpx.HTTPError as e:
            raise ExternalAPIError(f"retrieve pokemon: {e}") from e
        except Exception as e:
            raise ExternalAPIError(f"decode pokemon response: {e}") from e

    async def close(self):
        """Close HTTP client."""
        await self.http_client.aclose()

