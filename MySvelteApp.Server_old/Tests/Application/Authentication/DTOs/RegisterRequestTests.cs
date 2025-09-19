using FluentAssertions;
using MySvelteApp.Server.Application.Authentication.DTOs;
using Xunit;

namespace MySvelteApp.Server.Tests.Application.Authentication.DTOs;

public class RegisterRequestTests
{
    [Fact]
    public void RegisterRequest_DefaultValues_ShouldBeEmptyStrings()
    {
        // Act
        var request = new RegisterRequest();

        // Assert
        request.Username.Should().BeEmpty();
        request.Email.Should().BeEmpty();
        request.Password.Should().BeEmpty();
    }

    [Theory]
    [InlineData("testuser", "test@example.com", "testpassword")]
    [InlineData("user.name", "user+repo@example.co", "P@ssw0rd!")]
    public void RegisterRequest_WithValues_ShouldSetProperties(
        string username,
        string email,
        string password
    )
    {
        var request = new RegisterRequest
        {
            Username = username,
            Email = email,
            Password = password,
        };

        request.Username.Should().Be(username);
        request.Email.Should().Be(email);
        request.Password.Should().Be(password);
    }
}
