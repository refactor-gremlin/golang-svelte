namespace MySvelteApp.Server.Application.Authentication.DTOs;

public class AuthResult
{
    public bool Success { get; init; }
    public string? ErrorMessage { get; init; }
    public int? UserId { get; init; }
    public string? Username { get; init; }
    public string? Token { get; init; }
    public AuthErrorType ErrorType { get; init; } = AuthErrorType.None;
}
