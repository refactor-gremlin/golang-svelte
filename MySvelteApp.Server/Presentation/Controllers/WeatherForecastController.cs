using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using MySvelteApp.Server.Application.Weather;
using MySvelteApp.Server.Application.Weather.DTOs;

namespace MySvelteApp.Server.Presentation.Controllers;

[ApiController]
[Route("[controller]")]
public class WeatherForecastController : ControllerBase
{
    private readonly IWeatherForecastService _weatherForecastService;

    public WeatherForecastController(IWeatherForecastService weatherForecastService)
    {
        _weatherForecastService = weatherForecastService;
    }

    [HttpGet]
    [AllowAnonymous]
    [ProducesResponseType(typeof(IEnumerable<WeatherForecastDto>), StatusCodes.Status200OK)]
    public ActionResult<IEnumerable<WeatherForecastDto>> Get()
    {
        var forecasts = _weatherForecastService.GetForecasts();
        return Ok(forecasts);
    }
}
