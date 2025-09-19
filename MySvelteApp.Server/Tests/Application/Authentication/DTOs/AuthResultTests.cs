using FluentAssertions;
using MySvelteApp.Server.Application.Authentication.DTOs;
using Xunit;

namespace MySvelteApp.Server.Tests.Application.Authentication.DTOs;

public class AuthResultTests
{
    [Fact]
    public void AuthResult_DefaultValues_ShouldBeCorrect()
    {
        // Act
        var result = new AuthResult();

        // Assert
        result.Success.Should().BeFalse();
        result.ErrorMessage.Should().BeNull();
        result.UserId.Should().BeNull();
        result.Username.Should().BeNull();
        result.Token.Should().BeNull();
        result.ErrorType.Should().Be(AuthErrorType.None);
    }

    [Fact]
    public void AuthResult_WithSuccessValues_ShouldHaveCorrectProperties()
    {
        // Arrange
        const int userId = 123;
        const string username = "testuser";
        const string token = "jwt.token.here";

        // Act
        var result = new AuthResult
        {
            Success = true,
            UserId = userId,
            Username = username,
            Token = token,
        };

        // Assert
        result.Success.Should().BeTrue();
        result.UserId.Should().Be(userId);
        result.Username.Should().Be(username);
        result.Token.Should().Be(token);
        result.ErrorMessage.Should().BeNull();
        result.ErrorType.Should().Be(AuthErrorType.None);
    }

    [Fact]
    public void AuthResult_WithErrorValues_ShouldHaveCorrectProperties()
    {
        // Arrange
        const string errorMessage = "Invalid credentials";
        const AuthErrorType errorType = AuthErrorType.Unauthorized;

        // Act
        var result = new AuthResult
        {
            Success = false,
            ErrorMessage = errorMessage,
            ErrorType = errorType,
        };

        // Assert
        result.Success.Should().BeFalse();
        result.ErrorMessage.Should().Be(errorMessage);
        result.ErrorType.Should().Be(errorType);
        result.UserId.Should().BeNull();
        result.Username.Should().BeNull();
        result.Token.Should().BeNull();
    }

    [Fact]
    public void AuthResult_InitAssignments_ShouldPersistValues()
    {
        // Arrange
        var result = new AuthResult
        {
            Success = false,
            UserId = 123,
            Username = "testuser",
            Token = "jwt.token",
            ErrorMessage = "error",
            ErrorType = AuthErrorType.Validation,
        };

        // Act & Assert - Values set via object initializer should persist
        result.Success.Should().BeFalse();
        result.UserId.Should().Be(123);
        result.Username.Should().Be("testuser");
        result.Token.Should().Be("jwt.token");
        result.ErrorMessage.Should().Be("error");
        result.ErrorType.Should().Be(AuthErrorType.Validation);
    }
}
