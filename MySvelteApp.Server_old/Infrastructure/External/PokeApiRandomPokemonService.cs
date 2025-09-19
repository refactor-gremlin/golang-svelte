using System.Text.Json;
using MySvelteApp.Server.Application.Pokemon;
using MySvelteApp.Server.Application.Pokemon.DTOs;

namespace MySvelteApp.Server.Infrastructure.External;

public class PokeApiRandomPokemonService : IRandomPokemonService
{
    private const string PokemonApiBaseUrl = "https://pokeapi.co/api/v2/pokemon/";
    private readonly HttpClient _httpClient;

    public PokeApiRandomPokemonService(HttpClient httpClient)
    {
        _httpClient = httpClient;
    }

    public async Task<RandomPokemonDto> GetRandomPokemonAsync(
        CancellationToken cancellationToken = default
    )
    {
        var count = await GetPokemonCountAsync(cancellationToken);
        var randomPokemon = Random.Shared.Next(1, count + 1);
        var pokemonUrl = $"{PokemonApiBaseUrl}{randomPokemon}";

        using var pokemonResponse = await _httpClient.GetAsync(pokemonUrl, cancellationToken);
        pokemonResponse.EnsureSuccessStatusCode();
        await using var pokemonContent = await pokemonResponse.Content.ReadAsStreamAsync(
            cancellationToken
        );
        var pokeApi = await JsonSerializer.DeserializeAsync<PokeApiResponse>(
            pokemonContent,
            cancellationToken: cancellationToken
        );
        return pokeApi is null
            ? throw new InvalidOperationException("Failed to deserialize Pokemon data.")
            : new RandomPokemonDto
            {
                Name = pokeApi.name,
                Type = string.Join(", ", pokeApi.types.Select(t => t.type.name)),
                Image = pokeApi.sprites.front_default,
            };
    }

    private async Task<int> GetPokemonCountAsync(CancellationToken cancellationToken)
    {
        using var response = await _httpClient.GetAsync(
            "https://pokeapi.co/api/v2/pokemon-species/?limit=0",
            cancellationToken
        );
        response.EnsureSuccessStatusCode();
        await using var content = await response.Content.ReadAsStreamAsync(cancellationToken);
        using var document = await JsonDocument.ParseAsync(
            content,
            cancellationToken: cancellationToken
        );
        return document.RootElement.GetProperty("count").GetInt32();
    }

    private sealed class PokeApiResponse
    {
        public string name { get; set; } = string.Empty;
        public List<PokeApiType> types { get; set; } = new();
        public PokeApiSprites sprites { get; set; } = new();
    }

    private sealed class PokeApiType
    {
        public TypeInfo type { get; set; } = new();
    }

    private sealed class TypeInfo
    {
        public string name { get; set; } = string.Empty;
    }

    private sealed class PokeApiSprites
    {
        public string? front_default { get; set; }
    }
}
