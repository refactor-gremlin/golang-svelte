using MySvelteApp.Server.Application.Pokemon.DTOs;

namespace MySvelteApp.Server.Application.Pokemon;

public interface IRandomPokemonService
{
    Task<RandomPokemonDto> GetRandomPokemonAsync(CancellationToken cancellationToken = default);
}
