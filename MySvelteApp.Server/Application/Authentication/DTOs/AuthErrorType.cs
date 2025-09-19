namespace MySvelteApp.Server.Application.Authentication.DTOs;

public enum AuthErrorType
{
    None = 0,
    Validation = 1,
    Conflict = 2,
    Unauthorized = 3,
}
