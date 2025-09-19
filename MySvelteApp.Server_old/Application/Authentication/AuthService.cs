using System.Text.RegularExpressions;
using MySvelteApp.Server.Application.Authentication.DTOs;
using MySvelteApp.Server.Application.Common.Interfaces;
using MySvelteApp.Server.Domain.Entities;

namespace MySvelteApp.Server.Application.Authentication;

public class AuthService : IAuthService
{
    private static readonly Regex UsernameRegex = new("^[a-zA-Z0-9_]+$", RegexOptions.Compiled);
    private static readonly Regex EmailRegex = new(
        "^(?!.*\\.\\.)[^\\s@]+@[^\\s@.]+\\.[^\\s@.]+$",
        RegexOptions.Compiled | RegexOptions.CultureInvariant
    );
    private static readonly Regex PasswordRegex = new(
        "(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)",
        RegexOptions.Compiled
    );

    // Max length constants that mirror EF persistence constraints
    private const int MaxUsernameLength = 64;
    private const int MaxEmailLength = 320;
    private const int MaxPasswordLength = 512;

    private readonly IUserRepository _userRepository;
    private readonly IPasswordHasher _passwordHasher;
    private readonly IJwtTokenGenerator _jwtTokenGenerator;

    public AuthService(
        IUserRepository userRepository,
        IPasswordHasher passwordHasher,
        IJwtTokenGenerator jwtTokenGenerator
    )
    {
        _userRepository = userRepository;
        _passwordHasher = passwordHasher;
        _jwtTokenGenerator = jwtTokenGenerator;
    }

    public async Task<AuthResult> RegisterAsync(
        RegisterRequest request,
        CancellationToken cancellationToken = default
    )
    {
        var validationError = ValidateRegisterRequest(request);
        if (validationError is not null)
        {
            return CreateError(validationError.Value.Message, validationError.Value.Type);
        }

        var trimmedUsername = request.Username.Trim();
        var normalizedEmail = request.Email.Trim().ToLowerInvariant();

        if (await _userRepository.UsernameExistsAsync(trimmedUsername, cancellationToken))
        {
            return CreateError(
                "This username is already taken. Please choose a different one.",
                AuthErrorType.Conflict
            );
        }

        if (await _userRepository.EmailExistsAsync(normalizedEmail, cancellationToken))
        {
            return CreateError(
                "This email is already registered. Please use a different email address.",
                AuthErrorType.Conflict
            );
        }

        var (hash, salt) = _passwordHasher.HashPassword(request.Password);

        var user = new User
        {
            Username = trimmedUsername,
            Email = normalizedEmail,
            PasswordHash = hash,
            PasswordSalt = salt,
        };

        await _userRepository.AddAsync(user, cancellationToken);

        var token = _jwtTokenGenerator.GenerateToken(user);

        return new AuthResult
        {
            Success = true,
            Token = token,
            UserId = user.Id,
            Username = user.Username,
        };
    }

    public async Task<AuthResult> LoginAsync(
        LoginRequest request,
        CancellationToken cancellationToken = default
    )
    {
        var validationError = ValidateLoginRequest(request);
        if (validationError is not null)
        {
            return CreateError(validationError.Value.Message, validationError.Value.Type);
        }

        var trimmedUsername = request.Username.Trim();
        var user = await _userRepository.GetByUsernameAsync(trimmedUsername, cancellationToken);

        if (user is null)
        {
            return CreateUnauthorized();
        }

        var passwordValid = _passwordHasher.VerifyPassword(
            request.Password,
            user.PasswordHash,
            user.PasswordSalt
        );
        if (!passwordValid)
        {
            return CreateUnauthorized();
        }

        var token = _jwtTokenGenerator.GenerateToken(user);

        return new AuthResult
        {
            Success = true,
            Token = token,
            UserId = user.Id,
            Username = user.Username,
        };
    }

    private static (string Message, AuthErrorType Type)? ValidateRegisterRequest(
        RegisterRequest request
    )
    {
        var username = request.Username?.Trim();
        var email = request.Email?.Trim();

        return string.IsNullOrWhiteSpace(username)
                ? ((string Message, AuthErrorType Type)?)
                    ("Username is required.", AuthErrorType.Validation)
            : username.Length < 3
                ? ((string Message, AuthErrorType Type)?)
                    ("Username must be at least 3 characters long.", AuthErrorType.Validation)
            : username.Length > MaxUsernameLength
                ? ((string Message, AuthErrorType Type)?)
                    (
                        $"Username must not exceed {MaxUsernameLength} characters.",
                        AuthErrorType.Validation
                    )
            : !UsernameRegex.IsMatch(username)
                ? ((string Message, AuthErrorType Type)?)
                    (
                        "Username can only contain letters, numbers, and underscores.",
                        AuthErrorType.Validation
                    )
            : string.IsNullOrWhiteSpace(email)
                ? ((string Message, AuthErrorType Type)?)
                    ("Email is required.", AuthErrorType.Validation)
            : email.Length > MaxEmailLength
                ? ((string Message, AuthErrorType Type)?)
                    (
                        $"Email must not exceed {MaxEmailLength} characters.",
                        AuthErrorType.Validation
                    )
            : !EmailRegex.IsMatch(email)
                ? ((string Message, AuthErrorType Type)?)
                    ("Please enter a valid email address.", AuthErrorType.Validation)
            : string.IsNullOrWhiteSpace(request.Password)
                ? ((string Message, AuthErrorType Type)?)
                    ("Password is required.", AuthErrorType.Validation)
            : request.Password.Length < 8
                ? ((string Message, AuthErrorType Type)?)
                    ("Password must be at least 8 characters long.", AuthErrorType.Validation)
            : request.Password.Length > MaxPasswordLength
                ? ((string Message, AuthErrorType Type)?)
                    (
                        $"Password must not exceed {MaxPasswordLength} characters.",
                        AuthErrorType.Validation
                    )
            : !PasswordRegex.IsMatch(request.Password)
                ? (
                    "Password must contain at least one uppercase letter, one lowercase letter, and one number.",
                    AuthErrorType.Validation
                )
            : null;
    }

    private static (string Message, AuthErrorType Type)? ValidateLoginRequest(LoginRequest request)
    {
        var username = request.Username?.Trim();
        return string.IsNullOrWhiteSpace(username)
                ? ((string Message, AuthErrorType Type)?)
                    ("Username is required.", AuthErrorType.Validation)
            : string.IsNullOrWhiteSpace(request.Password)
                ? ("Password is required.", AuthErrorType.Validation)
            : null;
    }

    private static AuthResult CreateError(string message, AuthErrorType errorType)
    {
        return new AuthResult
        {
            Success = false,
            ErrorMessage = message,
            ErrorType = errorType,
        };
    }

    private static AuthResult CreateUnauthorized()
    {
        return new AuthResult
        {
            Success = false,
            ErrorMessage =
                "Invalid username or password. Please check your credentials and try again.",
            ErrorType = AuthErrorType.Unauthorized,
        };
    }
}
