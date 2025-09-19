using MySvelteApp.Server.Application.Weather.DTOs;

namespace MySvelteApp.Server.Application.Weather;

public interface IWeatherForecastService
{
    IEnumerable<WeatherForecastDto> GetForecasts();
}
