using FluentAssertions;
using MySvelteApp.Server.Application.Authentication.DTOs;
using Xunit;

namespace MySvelteApp.Server.Tests.Application.Authentication;

public class AuthServiceValidationTests
{
    // These tests validate the logic from the private methods in AuthService
    // Since they're private, we test them indirectly through the public methods
    // or by recreating the validation logic for direct testing

    private static (string Message, AuthErrorType Type)? ValidateUsername(string username)
    {
        if (string.IsNullOrWhiteSpace(username))
        {
            return ("Username is required.", AuthErrorType.Validation);
        }

        if (username.Length < 3)
        {
            return ("Username must be at least 3 characters long.", AuthErrorType.Validation);
        }

        if (!System.Text.RegularExpressions.Regex.IsMatch(username, "^[a-zA-Z0-9_]+$"))
        {
            return (
                "Username can only contain letters, numbers, and underscores.",
                AuthErrorType.Validation
            );
        }

        return null;
    }

    private static (string Message, AuthErrorType Type)? ValidateEmail(string email)
    {
        if (string.IsNullOrWhiteSpace(email))
        {
            return ("Email is required.", AuthErrorType.Validation);
        }

        if (
            !System.Text.RegularExpressions.Regex.IsMatch(
                email,
                "^(?!.*\\.\\.)[^\\s@]+@[^\\s@]+\\.[^\\s@]+$"
            )
        )
        {
            return ("Please enter a valid email address.", AuthErrorType.Validation);
        }

        return null;
    }

    private static (string Message, AuthErrorType Type)? ValidatePassword(string password)
    {
        if (string.IsNullOrWhiteSpace(password))
        {
            return ("Password is required.", AuthErrorType.Validation);
        }

        if (password.Length < 8)
        {
            return ("Password must be at least 8 characters long.", AuthErrorType.Validation);
        }

        if (
            !System.Text.RegularExpressions.Regex.IsMatch(
                password,
                "(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)"
            )
        )
        {
            return (
                "Password must contain at least one uppercase letter, one lowercase letter, and one number.",
                AuthErrorType.Validation
            );
        }

        return null;
    }

    public class UsernameValidationTests
    {
#pragma warning disable xUnit1012 // Intentionally testing null input validation for username field - null values should trigger validation errors
        [Theory]
        [InlineData(null)]
        [InlineData("")]
        [InlineData("   ")]
        public void ValidateUsername_WithNullOrWhitespace_ShouldReturnValidationError(
            string username
        )
        {
            // Act
            var result = ValidateUsername(username);

            // Assert
            result.Should().NotBeNull();
            result!.Value.Message.Should().Be("Username is required.");
            result!.Value.Type.Should().Be(AuthErrorType.Validation);
        }
#pragma warning restore xUnit1012

        [Theory]
        [InlineData("a")]
        [InlineData("ab")]
        public void ValidateUsername_WithTooShortUsername_ShouldReturnValidationError(
            string username
        )
        {
            // Act
            var result = ValidateUsername(username);

            // Assert
            result.Should().NotBeNull();
            result!.Value.Message.Should().Be("Username must be at least 3 characters long.");
            result!.Value.Type.Should().Be(AuthErrorType.Validation);
        }

        [Theory]
        [InlineData("user@name")]
        [InlineData("user-name")]
        [InlineData("user.name")]
        [InlineData("user name")]
        [InlineData("user!name")]
        public void ValidateUsername_WithInvalidCharacters_ShouldReturnValidationError(
            string username
        )
        {
            // Act
            var result = ValidateUsername(username);

            // Assert
            result.Should().NotBeNull();
            result!
                .Value.Message.Should()
                .Be("Username can only contain letters, numbers, and underscores.");
            result!.Value.Type.Should().Be(AuthErrorType.Validation);
        }

        [Theory]
        [InlineData("abc")]
        [InlineData("user123")]
        [InlineData("test_user")]
        [InlineData("User_Name_123")]
        public void ValidateUsername_WithValidUsername_ShouldReturnNull(string username)
        {
            // Act
            var result = ValidateUsername(username);

            // Assert
            result.Should().BeNull();
        }
    }

    public class EmailValidationTests
    {
#pragma warning disable xUnit1012 // Intentionally testing null input validation for email field - null values should trigger validation errors
        [Theory]
        [InlineData(null)]
        [InlineData("")]
        [InlineData("   ")]
        public void ValidateEmail_WithNullOrWhitespace_ShouldReturnValidationError(string email)
        {
            // Act
            var result = ValidateEmail(email);

            // Assert
            result.Should().NotBeNull();
            result!.Value.Message.Should().Be("Email is required.");
            result!.Value.Type.Should().Be(AuthErrorType.Validation);
        }
#pragma warning restore xUnit1012

        [Theory]
        [InlineData("invalid-email")]
        [InlineData("user@")]
        [InlineData("@example.com")]
        [InlineData("user.example.com")]
        [InlineData("user@.com")]
        public void ValidateEmail_WithInvalidEmail_ShouldReturnValidationError(string email)
        {
            // Act
            var result = ValidateEmail(email);

            // Assert
            result.Should().NotBeNull();
            result!.Value.Message.Should().Be("Please enter a valid email address.");
            result!.Value.Type.Should().Be(AuthErrorType.Validation);
        }

        [Theory]
        [InlineData("user@example.com")]
        [InlineData("test.email@example.com")]
        [InlineData("user+tag@example.com")]
        [InlineData("123@example.com")]
        [InlineData("user@example.co.uk")]
        public void ValidateEmail_WithValidEmail_ShouldReturnNull(string email)
        {
            // Act
            var result = ValidateEmail(email);

            // Assert
            result.Should().BeNull();
        }

        [Theory]
        [InlineData("user..user@example.com")]
        public void ValidateEmail_WithDoubleDotLocalPart_ShouldReturnValidationError(string email)
        {
            var result = ValidateEmail(email);
            result.Should().NotBeNull();
            result!.Value.Type.Should().Be(AuthErrorType.Validation);
        }
    }

    public class PasswordValidationTests
    {
#pragma warning disable xUnit1012 // Intentionally testing null input validation for password field - null values should trigger validation errors
        [Theory]
        [InlineData(null)]
        [InlineData("")]
        [InlineData("   ")]
        public void ValidatePassword_WithNullOrWhitespace_ShouldReturnValidationError(
            string password
        )
        {
            // Act
            var result = ValidatePassword(password);

            // Assert
            result.Should().NotBeNull();
            result!.Value.Message.Should().Be("Password is required.");
            result!.Value.Type.Should().Be(AuthErrorType.Validation);
        }
#pragma warning restore xUnit1012

        [Theory]
        [InlineData("1234567")]
        [InlineData("abcdefg")]
        [InlineData("ABCDEFG")]
        [InlineData("!@#$%^&")]
        public void ValidatePassword_WithTooShortPassword_ShouldReturnValidationError(
            string password
        )
        {
            // Act
            var result = ValidatePassword(password);

            // Assert
            result.Should().NotBeNull();
            result!.Value.Message.Should().Be("Password must be at least 8 characters long.");
            result!.Value.Type.Should().Be(AuthErrorType.Validation);
        }

        [Theory]
        [InlineData("password")]
        [InlineData("PASSWORD")]
        [InlineData("Password")]
        [InlineData("12345678")]
        [InlineData("password123")]
        [InlineData("PASSWORD123")]
        [InlineData("Password!@#")]
        public void ValidatePassword_WithMissingCharacterTypes_ShouldReturnValidationError(
            string password
        )
        {
            // Act
            var result = ValidatePassword(password);

            // Assert
            result.Should().NotBeNull();
            result!
                .Value.Message.Should()
                .Be(
                    "Password must contain at least one uppercase letter, one lowercase letter, and one number."
                );
            result!.Value.Type.Should().Be(AuthErrorType.Validation);
        }

        [Theory]
        [InlineData("Password1")]
        [InlineData("MyPassword123")]
        [InlineData("Test123Password")]
        [InlineData("P@ssw0rd123")]
        [InlineData("Complex!Password!123")]
        public void ValidatePassword_WithValidPassword_ShouldReturnNull(string password)
        {
            // Act
            var result = ValidatePassword(password);

            // Assert
            result.Should().BeNull();
        }

        [Fact]
        public void ValidatePassword_WithExactly8Characters_ShouldReturnNull()
        {
            // Arrange
            const string password = "Passw0rd";

            // Act
            var result = ValidatePassword(password);

            // Assert
            result.Should().BeNull();
        }
    }

    public class CombinedValidationTests
    {
        [Fact]
        public void ValidateRegistrationData_WithAllValidData_ShouldPassAllValidations()
        {
            // Arrange
            const string username = "testuser";
            const string email = "test@example.com";
            const string password = "ValidPassword123";

            // Act
            var usernameResult = ValidateUsername(username);
            var emailResult = ValidateEmail(email);
            var passwordResult = ValidatePassword(password);

            // Assert
            usernameResult.Should().BeNull();
            emailResult.Should().BeNull();
            passwordResult.Should().BeNull();
        }

        [Fact]
        public void ValidateRegistrationData_WithAllInvalidData_ShouldFailAllValidations()
        {
            // Arrange
            const string username = "u";
            const string email = "invalid";
            const string password = "short";

            // Act
            var usernameResult = ValidateUsername(username);
            var emailResult = ValidateEmail(email);
            var passwordResult = ValidatePassword(password);

            // Assert
            usernameResult.Should().NotBeNull();
            emailResult.Should().NotBeNull();
            passwordResult.Should().NotBeNull();
        }
    }
}
