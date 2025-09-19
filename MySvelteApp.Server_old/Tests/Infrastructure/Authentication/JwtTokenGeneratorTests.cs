using System;
using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;
using FluentAssertions;
using Microsoft.Extensions.Options;
using Microsoft.IdentityModel.Tokens;
using MySvelteApp.Server.Domain.Entities;
using MySvelteApp.Server.Infrastructure.Authentication;
using MySvelteApp.Server.Tests.TestUtilities;
using Xunit;

namespace MySvelteApp.Server.Tests.Infrastructure.Authentication;

public class JwtTokenGeneratorTests
{
    private readonly IOptions<JwtOptions> _jwtOptions;
    private readonly JwtTokenGenerator _jwtTokenGenerator;

    public JwtTokenGeneratorTests()
    {
        _jwtOptions = TestHelper.CreateJwtOptions();
        _jwtTokenGenerator = new JwtTokenGenerator(_jwtOptions);
    }

    [Fact]
    public void GenerateToken_WithValidUser_ShouldReturnValidJwtToken()
    {
        // Arrange
        var user = TestData.Users.ValidUser;

        // Act
        var token = _jwtTokenGenerator.GenerateToken(user);

        // Assert
        token.Should().NotBeNullOrEmpty();

        var tokenHandler = new JwtSecurityTokenHandler();
        var jwtToken = tokenHandler.ReadJwtToken(token);

        jwtToken.Issuer.Should().Be(_jwtOptions.Value.Issuer);
        jwtToken.Audiences.Should().Contain(_jwtOptions.Value.Audience);
        jwtToken
            .ValidTo.Should()
            .BeCloseTo(
                DateTime.UtcNow.AddHours(_jwtOptions.Value.AccessTokenLifetimeHours),
                TimeSpan.FromMinutes(1)
            );

        // Build the signing key similar to JwtTokenGenerator
        var keyString = _jwtOptions.Value.Key;
        var keyBytes = keyString.StartsWith("base64:", StringComparison.OrdinalIgnoreCase)
            ? Convert.FromBase64String(keyString.Substring("base64:".Length))
            : Encoding.UTF8.GetBytes(keyString);

        var validationParams = new TokenValidationParameters
        {
            ValidateIssuerSigningKey = true,
            IssuerSigningKey = new SymmetricSecurityKey(keyBytes),
            ValidateIssuer = true,
            ValidIssuer = _jwtOptions.Value.Issuer,
            ValidateAudience = true,
            ValidAudience = _jwtOptions.Value.Audience,
            ValidateLifetime = true,
            ClockSkew = TimeSpan.FromMinutes(1),
        };

        new JwtSecurityTokenHandler().ValidateToken(token, validationParams, out _);
    }

    [Fact]
    public void GenerateToken_ShouldContainCorrectClaims()
    {
        // Arrange
        var user = TestData.Users.ValidUser;

        // Act
        var token = _jwtTokenGenerator.GenerateToken(user);

        // Assert
        var tokenHandler = new JwtSecurityTokenHandler();
        var jwtToken = tokenHandler.ReadJwtToken(token);

        var claims = jwtToken.Claims.ToList();

        claims
            .Should()
            .Contain(c => c.Type == ClaimTypes.NameIdentifier && c.Value == user.Id.ToString());
        claims
            .Should()
            .Contain(c => c.Type == JwtRegisteredClaimNames.Sub && c.Value == user.Id.ToString());
        claims.Should().Contain(c => c.Type == ClaimTypes.Name && c.Value == user.Username);
        claims.Count(c => c.Type == JwtRegisteredClaimNames.Jti).Should().BeGreaterThan(0);
        claims.First(c => c.Type == JwtRegisteredClaimNames.Jti).Value.Should().NotBeNullOrEmpty();
        claims.Count(c => c.Type == JwtRegisteredClaimNames.Iat).Should().BeGreaterThan(0);
        claims.First(c => c.Type == JwtRegisteredClaimNames.Iat).Value.Should().NotBeNullOrEmpty();
    }

    [Fact]
    public void GenerateToken_ShouldHaveUniqueJtiForEachToken()
    {
        // Arrange
        var user = TestData.Users.ValidUser;

        // Act
        var token1 = _jwtTokenGenerator.GenerateToken(user);
        var token2 = _jwtTokenGenerator.GenerateToken(user);

        // Assert
        var tokenHandler = new JwtSecurityTokenHandler();
        var jwtToken1 = tokenHandler.ReadJwtToken(token1);
        var jwtToken2 = tokenHandler.ReadJwtToken(token2);

        var jti1 = jwtToken1.Claims.First(c => c.Type == JwtRegisteredClaimNames.Jti).Value;
        var jti2 = jwtToken2.Claims.First(c => c.Type == JwtRegisteredClaimNames.Jti).Value;

        jti1.Should().NotBe(jti2);
    }

    [Fact]
    public void GenerateToken_WithDifferentUsers_ShouldHaveDifferentClaims()
    {
        // Arrange
        var user1 = new User
        {
            Id = 1,
            Username = TestData.Users.ValidUser.Username,
            Email = TestData.Users.ValidUser.Email,
            PasswordHash = TestData.Users.ValidUser.PasswordHash,
            PasswordSalt = TestData.Users.ValidUser.PasswordSalt,
        };
        var user2 = new User
        {
            Id = 2,
            Username = TestData.Users.AnotherValidUser.Username,
            Email = TestData.Users.AnotherValidUser.Email,
            PasswordHash = TestData.Users.AnotherValidUser.PasswordHash,
            PasswordSalt = TestData.Users.AnotherValidUser.PasswordSalt,
        };

        // Act
        var token1 = _jwtTokenGenerator.GenerateToken(user1);
        var token2 = _jwtTokenGenerator.GenerateToken(user2);

        // Assert
        var tokenHandler = new JwtSecurityTokenHandler();
        var jwtToken1 = tokenHandler.ReadJwtToken(token1);
        var jwtToken2 = tokenHandler.ReadJwtToken(token2);

        var claims1 = jwtToken1.Claims.ToList();
        var claims2 = jwtToken2.Claims.ToList();

        claims1
            .First(c => c.Type == ClaimTypes.NameIdentifier)
            .Value.Should()
            .NotBe(claims2.First(c => c.Type == ClaimTypes.NameIdentifier).Value);
        claims1
            .First(c => c.Type == ClaimTypes.Name)
            .Value.Should()
            .NotBe(claims2.First(c => c.Type == ClaimTypes.Name).Value);
    }

    [Fact]
    public void GenerateToken_WithBase64Key_ShouldWorkCorrectly()
    {
        // Arrange
        var jwtOptions = Options.Create(
            new JwtOptions
            {
                Key = TestData.Jwt.ValidKey, // base64 encoded key
                Issuer = TestData.Jwt.ValidIssuer,
                Audience = TestData.Jwt.ValidAudience,
                AccessTokenLifetimeHours = TestData.Jwt.ValidLifetimeHours,
            }
        );

        var generator = new JwtTokenGenerator(jwtOptions);
        var user = TestData.Users.ValidUser;

        // Act
        var token = generator.GenerateToken(user);

        // Assert
        token.Should().NotBeNullOrEmpty();

        var tokenHandler = new JwtSecurityTokenHandler();
        var jwtToken = tokenHandler.ReadJwtToken(token);

        jwtToken.Issuer.Should().Be(jwtOptions.Value.Issuer);
        jwtToken.Audiences.Should().Contain(jwtOptions.Value.Audience);

        // Validate signature using the same key material
        var keyString = jwtOptions.Value.Key;
        var keyBytes = keyString.StartsWith("base64:", StringComparison.OrdinalIgnoreCase)
            ? Convert.FromBase64String(keyString.Substring("base64:".Length))
            : Encoding.UTF8.GetBytes(keyString);
        var validationParams = new TokenValidationParameters
        {
            ValidateIssuerSigningKey = true,
            IssuerSigningKey = new SymmetricSecurityKey(keyBytes),
            ValidateIssuer = true,
            ValidIssuer = jwtOptions.Value.Issuer,
            ValidateAudience = true,
            ValidAudience = jwtOptions.Value.Audience,
            ValidateLifetime = true,
            ClockSkew = TimeSpan.FromMinutes(1),
        };
        new JwtSecurityTokenHandler().ValidateToken(token, validationParams, out _);
    }

    [Fact]
    public void GenerateToken_WithPlainTextKey_ShouldWorkCorrectly()
    {
        // Arrange
        var longPlainKey = TestData.Jwt.ValidPlainTextKey;
        var jwtOptions = Options.Create(
            new JwtOptions
            {
                Key = longPlainKey,
                Issuer = TestData.Jwt.ValidIssuer,
                Audience = TestData.Jwt.ValidAudience,
                AccessTokenLifetimeHours = TestData.Jwt.ValidLifetimeHours,
            }
        );

        var generator = new JwtTokenGenerator(jwtOptions);
        var user = TestData.Users.ValidUser;

        // Act
        var token = generator.GenerateToken(user);

        // Assert
        token.Should().NotBeNullOrEmpty();

        var tokenHandler = new JwtSecurityTokenHandler();
        var jwtToken = tokenHandler.ReadJwtToken(token);

        jwtToken.Issuer.Should().Be(jwtOptions.Value.Issuer);
        jwtToken.Audiences.Should().Contain(jwtOptions.Value.Audience);

        // Validate signature using the plain-text key
        var keyBytes = Encoding.UTF8.GetBytes(jwtOptions.Value.Key);
        var validationParams = new TokenValidationParameters
        {
            ValidateIssuerSigningKey = true,
            IssuerSigningKey = new SymmetricSecurityKey(keyBytes),
            ValidateIssuer = true,
            ValidIssuer = jwtOptions.Value.Issuer,
            ValidateAudience = true,
            ValidAudience = jwtOptions.Value.Audience,
            ValidateLifetime = true,
            ClockSkew = TimeSpan.FromMinutes(1),
        };
        new JwtSecurityTokenHandler().ValidateToken(token, validationParams, out _);
    }

    [Fact]
    public void GenerateToken_ShouldHaveCorrectExpirationTime()
    {
        // Arrange
        var user = TestData.Users.ValidUser;
        var beforeGeneration = DateTime.UtcNow;

        // Act
        var token = _jwtTokenGenerator.GenerateToken(user);

        // Assert
        var tokenHandler = new JwtSecurityTokenHandler();
        var jwtToken = tokenHandler.ReadJwtToken(token);

        var expectedExpiration = beforeGeneration.AddHours(
            _jwtOptions.Value.AccessTokenLifetimeHours
        );
        jwtToken.ValidTo.Should().BeCloseTo(expectedExpiration, TimeSpan.FromMinutes(1));
    }

    [Fact]
    public void GenerateToken_ShouldUseHmacSha256Algorithm()
    {
        // Arrange
        var user = TestData.Users.ValidUser;

        // Act
        var token = _jwtTokenGenerator.GenerateToken(user);

        // Assert
        var tokenHandler = new JwtSecurityTokenHandler();
        var jwtToken = tokenHandler.ReadJwtToken(token);

        jwtToken.SignatureAlgorithm.Should().Be("HS256");
    }
}
