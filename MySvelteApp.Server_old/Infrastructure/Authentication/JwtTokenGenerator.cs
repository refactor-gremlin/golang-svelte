using System.Globalization;
using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using Microsoft.Extensions.Options;
using Microsoft.IdentityModel.Tokens;
using MySvelteApp.Server.Application.Common.Interfaces;
using MySvelteApp.Server.Domain.Entities;

namespace MySvelteApp.Server.Infrastructure.Authentication;

public class JwtTokenGenerator(IOptions<JwtOptions> jwtOptions) : IJwtTokenGenerator
{
    private readonly JwtOptions _jwtOptions = jwtOptions.Value;

    public string GenerateToken(User user)
    {
        var keyBytes = DeriveKeyBytes(_jwtOptions.Key);
        var securityKey = new SymmetricSecurityKey(keyBytes);
        var credentials = new SigningCredentials(securityKey, SecurityAlgorithms.HmacSha256);

        var claims = new[]
        {
            new Claim(ClaimTypes.NameIdentifier, user.Id.ToString(CultureInfo.InvariantCulture)),
            new Claim(JwtRegisteredClaimNames.Sub, user.Id.ToString(CultureInfo.InvariantCulture)),
            new Claim(ClaimTypes.Name, user.Username),
            new Claim(JwtRegisteredClaimNames.Jti, Guid.NewGuid().ToString()),
            new Claim(
                JwtRegisteredClaimNames.Iat,
                EpochTime.GetIntDate(DateTime.UtcNow).ToString(CultureInfo.InvariantCulture),
                ClaimValueTypes.Integer64
            ),
        };

        var token = new JwtSecurityToken(
            issuer: _jwtOptions.Issuer,
            audience: _jwtOptions.Audience,
            claims: claims,
            expires: DateTime.UtcNow.AddHours(_jwtOptions.AccessTokenLifetimeHours),
            signingCredentials: credentials
        );

        return new JwtSecurityTokenHandler().WriteToken(token);
    }

    private static byte[] DeriveKeyBytes(string key)
    {
        return key.StartsWith("base64:", StringComparison.Ordinal)
            ? Convert.FromBase64String(key["base64:".Length..])
            : Encoding.UTF8.GetBytes(key);
    }
}
