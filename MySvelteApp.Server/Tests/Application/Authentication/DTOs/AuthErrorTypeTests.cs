using FluentAssertions;
using MySvelteApp.Server.Application.Authentication.DTOs;
using Xunit;

namespace MySvelteApp.Server.Tests.Application.Authentication.DTOs;

public class AuthErrorTypeTests
{
    [Fact]
    public void AuthErrorType_ShouldHaveCorrectNames()
    {
        // Assert
        AuthErrorType.None.ToString().Should().Be("None");
        AuthErrorType.Validation.ToString().Should().Be("Validation");
        AuthErrorType.Conflict.ToString().Should().Be("Conflict");
        AuthErrorType.Unauthorized.ToString().Should().Be("Unauthorized");
    }

    [Fact]
    public void AuthErrorType_ShouldBeConvertibleToInt()
    {
        // Act & Assert
        int noneValue = (int)AuthErrorType.None;
        int validationValue = (int)AuthErrorType.Validation;
        int conflictValue = (int)AuthErrorType.Conflict;
        int unauthorizedValue = (int)AuthErrorType.Unauthorized;

        noneValue.Should().BeGreaterThanOrEqualTo(0);
        validationValue.Should().BeGreaterThanOrEqualTo(0);
        conflictValue.Should().BeGreaterThanOrEqualTo(0);
        unauthorizedValue.Should().BeGreaterThanOrEqualTo(0);
    }

    [Fact]
    public void AuthErrorType_ShouldBeConvertibleFromInt()
    {
        // Act & Assert
        ((AuthErrorType)0)
            .Should()
            .Be(AuthErrorType.None);
        ((AuthErrorType)1).Should().Be(AuthErrorType.Validation);
        ((AuthErrorType)2).Should().Be(AuthErrorType.Conflict);
        ((AuthErrorType)3).Should().Be(AuthErrorType.Unauthorized);
    }
}
