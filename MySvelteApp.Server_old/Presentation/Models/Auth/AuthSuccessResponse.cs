namespace MySvelteApp.Server.Presentation.Models.Auth;

public class AuthSuccessResponse
{
    public string Token { get; set; } = string.Empty;
    public int UserId { get; set; }
    public string Username { get; set; } = string.Empty;
}
