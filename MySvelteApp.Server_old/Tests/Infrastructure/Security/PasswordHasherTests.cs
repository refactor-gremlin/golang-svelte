using FluentAssertions;
using MySvelteApp.Server.Infrastructure.Security;
using Xunit;

namespace MySvelteApp.Server.Tests.Infrastructure.Security;

public class PasswordHasherTests
{
    private readonly PasswordHasher _passwordHasher;

    public PasswordHasherTests()
    {
        _passwordHasher = new PasswordHasher();
    }

    [Fact]
    public void HashPassword_ShouldReturnNonEmptyHashAndSalt()
    {
        // Arrange
        const string password = "TestPassword123";

        // Act
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Assert
        hash.Should().NotBeNullOrEmpty();
        salt.Should().NotBeNullOrEmpty();
        hash.Should().NotBe(password);
        salt.Should().NotBe(password);
    }

    [Fact]
    public void HashPassword_ShouldReturnBase64EncodedStrings()
    {
        // Arrange
        const string password = "TestPassword123";

        // Act
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Assert
        // Base64 strings should be valid and decodable
        Convert.FromBase64String(hash).Should().NotBeNull();
        Convert.FromBase64String(salt).Should().NotBeNull();
    }

    [Fact]
    public void HashPassword_ShouldGenerateDifferentHashesForSamePassword()
    {
        // Arrange
        const string password = "TestPassword123";

        // Act
        var (hash1, salt1) = _passwordHasher.HashPassword(password);
        var (hash2, salt2) = _passwordHasher.HashPassword(password);

        // Assert
        hash1.Should().NotBe(hash2);
        salt1.Should().NotBe(salt2);
    }

    [Fact]
    public void VerifyPassword_WithCorrectPassword_ShouldReturnTrue()
    {
        // Arrange
        const string password = "TestPassword123";
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Act
        var result = _passwordHasher.VerifyPassword(password, hash, salt);

        // Assert
        result.Should().BeTrue();
    }

    [Theory]
    [InlineData("WrongPassword456")]
    [InlineData("testpassword123")]
    [InlineData("TestPassword124")]
    [InlineData("")]
    public void VerifyPassword_WithWrongInputs_ShouldReturnFalse(string candidate)
    {
        const string password = "TestPassword123";
        var (hash, salt) = _passwordHasher.HashPassword(password);

        _passwordHasher.VerifyPassword(candidate, hash, salt).Should().BeFalse();
    }

    [Fact]
    public void VerifyPassword_WithNullPassword_ShouldThrow()
    {
        // Arrange
        const string password = "TestPassword123";
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Act
        Action act = () => _passwordHasher.VerifyPassword(null!, hash, salt);

        // Assert
        act.Should().Throw<ArgumentNullException>();
    }

    [Fact]
    public void HashPassword_WithNullPassword_ShouldThrow()
    {
        Action act = () => _passwordHasher.HashPassword(null!);
        act.Should().Throw<ArgumentNullException>();
    }

    [Fact]
    public void HashPassword_WithEmptyPassword_ShouldStillWork()
    {
        // Arrange
        const string password = "";

        // Act
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Assert
        hash.Should().NotBeNullOrEmpty();
        salt.Should().NotBeNullOrEmpty();

        // Verify it can be verified
        var result = _passwordHasher.VerifyPassword(password, hash, salt);
        result.Should().BeTrue();
    }

    [Fact]
    public void HashPassword_WithLongPassword_ShouldWork()
    {
        // Arrange
        string password = new string('A', 1000); // Very long password

        // Act
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Assert
        hash.Should().NotBeNullOrEmpty();
        salt.Should().NotBeNullOrEmpty();

        // Verify it can be verified
        var result = _passwordHasher.VerifyPassword(password, hash, salt);
        result.Should().BeTrue();
    }

    [Fact]
    public void HashPassword_WithSpecialCharacters_ShouldWork()
    {
        // Arrange
        const string password = "P@ssw0rd!#$%^&*()_+-=[]{}|;:,.<>?";

        // Act
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Assert
        hash.Should().NotBeNullOrEmpty();
        salt.Should().NotBeNullOrEmpty();

        // Verify it can be verified
        var result = _passwordHasher.VerifyPassword(password, hash, salt);
        result.Should().BeTrue();
    }

    [Fact]
    public void HashPassword_OutputSizes_ShouldBePositive()
    {
        // Arrange
        const string password = "TestPassword123";

        // Act
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Assert (algorithm-agnostic)
        var hashBytes = Convert.FromBase64String(hash);
        hashBytes.Length.Should().BeGreaterThan(0);

        var saltBytes = Convert.FromBase64String(salt);
        saltBytes.Length.Should().BeGreaterThan(0);
    }

    [Fact]
    public void VerifyPassword_WithTamperedHash_ShouldReturnFalse()
    {
        // Arrange
        const string password = "TestPassword123";
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Tamper with the hash (guaranteed change and valid Base64)
        var hashBytes = Convert.FromBase64String(hash);
        hashBytes[0] ^= 0x01;
        var tamperedHash = Convert.ToBase64String(hashBytes);

        // Act
        var result = _passwordHasher.VerifyPassword(password, tamperedHash, salt);

        // Assert
        result.Should().BeFalse();
    }

    [Fact]
    public void VerifyPassword_WithTamperedSalt_ShouldReturnFalse()
    {
        // Arrange
        const string password = "TestPassword123";
        var (hash, salt) = _passwordHasher.HashPassword(password);

        // Tamper with the salt (guaranteed change and valid Base64)
        var saltBytes = Convert.FromBase64String(salt);
        saltBytes[0] ^= 0x01;
        var tamperedSalt = Convert.ToBase64String(saltBytes);

        // Act
        var result = _passwordHasher.VerifyPassword(password, hash, tamperedSalt);

        // Assert
        result.Should().BeFalse();
    }

    [Fact]
    public void VerifyPassword_WithInvalidBase64Hash_ShouldThrow()
    {
        // Arrange
        var (hash, salt) = _passwordHasher.HashPassword("pw");

        // Act
        Action act = () => _passwordHasher.VerifyPassword("pw", "###", salt);

        // Assert
        act.Should().Throw<FormatException>();
    }

    [Fact]
    public void VerifyPassword_WithInvalidBase64Salt_ShouldThrow()
    {
        // Arrange
        var (hash, salt) = _passwordHasher.HashPassword("pw");

        // Act
        Action act = () => _passwordHasher.VerifyPassword("pw", hash, "###");

        // Assert
        act.Should().Throw<FormatException>();
    }
}
