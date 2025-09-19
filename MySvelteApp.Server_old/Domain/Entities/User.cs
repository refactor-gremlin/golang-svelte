namespace MySvelteApp.Server.Domain.Entities;

public class User
{
    public int Id { get; set; }
    private string _username = string.Empty;
    private string _email = string.Empty;
    private string _passwordHash = string.Empty;
    private string _passwordSalt = string.Empty;

    public string Username
    {
        get => _username;
        set
        {
            if (string.IsNullOrWhiteSpace(value))
            {
                throw new ArgumentException("Username cannot be null or whitespace", nameof(value));
            }
            if (value.Length > 64)
            {
                throw new ArgumentException("Username exceeds 64 characters.", nameof(value));
            }
            _username = value;
        }
    }

    public string Email
    {
        get => _email;
        set
        {
            if (string.IsNullOrWhiteSpace(value))
            {
                throw new ArgumentException("Email cannot be null or whitespace", nameof(value));
            }
            var v = value.Trim().ToLowerInvariant();
            if (v.Length > 320)
            {
                throw new ArgumentException("Email exceeds 320 characters.", nameof(value));
            }
            _email = v;
        }
    }

    public string PasswordHash
    {
        get => _passwordHash;
        set
        {
            if (string.IsNullOrWhiteSpace(value))
            {
                throw new ArgumentException(
                    "PasswordHash cannot be null or whitespace",
                    nameof(value)
                );
            }
            _passwordHash = value;
        }
    }

    public string PasswordSalt
    {
        get => _passwordSalt;
        set
        {
            if (string.IsNullOrWhiteSpace(value))
            {
                throw new ArgumentException(
                    "PasswordSalt cannot be null or whitespace",
                    nameof(value)
                );
            }
            _passwordSalt = value;
        }
    }
}
