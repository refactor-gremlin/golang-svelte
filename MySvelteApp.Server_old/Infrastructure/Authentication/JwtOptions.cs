using System.ComponentModel.DataAnnotations;
using System.Text;

namespace MySvelteApp.Server.Infrastructure.Authentication;

/// <summary>
/// Public validator class for JwtOptions to support DataAnnotations validation.
/// </summary>
public static class JwtOptionsValidator
{
    /// <summary>
    /// Validates that a string value is not null, empty, or whitespace-only.
    /// </summary>
    public static ValidationResult ValidateNotWhitespace(object? value, ValidationContext context)
    {
        var str = value as string;
        if (string.IsNullOrWhiteSpace(str))
        {
            var displayName = context.DisplayName ?? context.MemberName ?? "Value";
            return new ValidationResult($"{displayName} cannot be blank or whitespace-only.");
        }
        return ValidationResult.Success!;
    }

    public static ValidationResult ValidateKeyStrength(object? value, ValidationContext context)
    {
        var s = value as string ?? string.Empty;
        var bytes = DeriveKeyBytes(s);
        return bytes.Length >= 32
            ? ValidationResult.Success!
            : new ValidationResult("Jwt:Key must be at least 32 bytes (256-bit) after decoding.");
    }

    private static byte[] DeriveKeyBytes(string key)
    {
        return key.StartsWith("base64:", StringComparison.Ordinal)
            ? Convert.FromBase64String(key["base64:".Length..])
            : Encoding.UTF8.GetBytes(key);
    }
}

public sealed class JwtOptions
{
    [Required]
    [MinLength(32, ErrorMessage = "Jwt:Key must be at least 32 characters long.")]
    [CustomValidation(typeof(JwtOptionsValidator), nameof(JwtOptionsValidator.ValidateKeyStrength))]
    public string Key { get; set; } = string.Empty;

    [Required]
    [CustomValidation(
        typeof(JwtOptionsValidator),
        nameof(JwtOptionsValidator.ValidateNotWhitespace)
    )]
    public string Issuer { get; set; } = string.Empty;

    [Required]
    [CustomValidation(
        typeof(JwtOptionsValidator),
        nameof(JwtOptionsValidator.ValidateNotWhitespace)
    )]
    public string Audience { get; set; } = string.Empty;

    /// <summary>
    /// Lifetime of the access token in hours.
    /// </summary>
    [Range(1, 168, ErrorMessage = "Jwt:AccessTokenLifetimeHours must be between 1 and 168 hours.")]
    public int AccessTokenLifetimeHours { get; set; } = 24;
}
