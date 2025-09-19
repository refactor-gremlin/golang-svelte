using MySvelteApp.Server.Application.Weather;
using MySvelteApp.Server.Application.Weather.DTOs;

namespace MySvelteApp.Server.Infrastructure.Weather;

public class WeatherForecastService : IWeatherForecastService
{
    private static readonly string[] Summaries =
    {
        "Freezing",
        "Bracing",
        "Chilly",
        "Cool",
        "Mild",
        "Warm",
        "Balmy",
        "Hot",
        "Sweltering",
        "Scorching",
    };

    public IEnumerable<WeatherForecastDto> GetForecasts()
    {
        return Enumerable
            .Range(1, 5)
            .Select(index => new WeatherForecastDto
            {
                Date = DateOnly.FromDateTime(DateTime.UtcNow.AddDays(index)),
                TemperatureC = Random.Shared.Next(-20, 55),
                Summary = Summaries[Random.Shared.Next(Summaries.Length)],
            });
    }
}
