using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using MySvelteApp.Server.Application.Pokemon;
using MySvelteApp.Server.Application.Pokemon.DTOs;

namespace MySvelteApp.Server.Presentation.Controllers;

[ApiController]
[Route("[controller]")]
public class RandomPokemonController : ControllerBase
{
    private readonly IRandomPokemonService _randomPokemonService;

    public RandomPokemonController(IRandomPokemonService randomPokemonService)
    {
        _randomPokemonService = randomPokemonService;
    }

    [HttpGet]
    [AllowAnonymous]
    [ProducesResponseType(typeof(RandomPokemonDto), StatusCodes.Status200OK)]
    public async Task<ActionResult<RandomPokemonDto>> Get(CancellationToken cancellationToken)
    {
        var pokemon = await _randomPokemonService.GetRandomPokemonAsync(cancellationToken);
        return Ok(pokemon);
    }
}
