using FluentAssertions;
using MySvelteApp.Server.Application.Authentication.DTOs;
using Xunit;

namespace MySvelteApp.Server.Tests.Application.Authentication.DTOs;

public class LoginRequestTests
{
    [Fact]
    public void LoginRequest_DefaultValues_ShouldBeEmptyStrings()
    {
        // Act
        var request = new LoginRequest();

        // Assert
        request.Username.Should().BeEmpty();
        request.Password.Should().BeEmpty();
    }

    [Theory]
    [InlineData("testuser", "testpassword")]
    [InlineData("user2", "P@ssw0rd!")]
    public void LoginRequest_WithValidValues_ShouldSetPropertiesCorrectly(
        string username,
        string password
    )
    {
        // Act
        var request = new LoginRequest { Username = username, Password = password };

        // Assert
        request.Username.Should().Be(username);
        request.Password.Should().Be(password);
    }

    // Null handling is validated via service validation tests; DTO remains non-nullable
}
