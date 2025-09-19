using FluentAssertions;
using MySvelteApp.Server.Domain.Entities;
using Xunit;

namespace MySvelteApp.Server.Tests.Domain.Entities;

public class UserTests
{
    [Fact]
    public void User_DefaultValues_ShouldBeCorrect()
    {
        // Act
        var user = new User
        {
            Username = "testuser",
            Email = "test@example.com",
            PasswordHash = "hash",
            PasswordSalt = "salt",
        };

        // Assert
        user.Id.Should().Be(0);
        user.Username.Should().Be("testuser");
        user.Email.Should().Be("test@example.com");
        user.PasswordHash.Should().Be("hash");
        user.PasswordSalt.Should().Be("salt");
    }

    [Fact]
    public void User_WithValidValues_ShouldSetPropertiesCorrectly()
    {
        // Arrange
        const int id = 123;
        const string username = "testuser";
        const string email = "test@example.com";
        const string passwordHash = "hashed_password";
        const string passwordSalt = "salt_value";

        // Act
        var user = new User
        {
            Id = id,
            Username = username,
            Email = email,
            PasswordHash = passwordHash,
            PasswordSalt = passwordSalt,
        };

        // Assert
        user.Id.Should().Be(id);
        user.Username.Should().Be(username);
        user.Email.Should().Be(email);
        user.PasswordHash.Should().Be(passwordHash);
        user.PasswordSalt.Should().Be(passwordSalt);
    }

    [Fact]
    public void User_Id_ShouldBeMutable()
    {
        // Arrange
        var user = new User();

        // Act
        user.Id = 456;

        // Assert
        user.Id.Should().Be(456);
    }

    [Fact]
    public void User_Properties_ShouldBeMutable()
    {
        // Arrange
        var user = new User
        {
            Username = "orig",
            Email = "orig@example.com",
            PasswordHash = "origHash",
            PasswordSalt = "origSalt",
        };

        // Act
        user.Username = "newusername";
        user.Email = "newemail@example.com";
        user.PasswordHash = "newhash";
        user.PasswordSalt = "newsalt";

        // Assert
        user.Username.Should().Be("newusername");
        user.Email.Should().Be("newemail@example.com");
        user.PasswordHash.Should().Be("newhash");
        user.PasswordSalt.Should().Be("newsalt");
    }

    [Theory]
    [InlineData(null)]
    [InlineData("")]
    [InlineData("   ")]
    public void Setting_NullOrEmpty_Username_ShouldThrow(object? value)
    {
        var user = new User
        {
            Username = "ok",
            Email = "e@example.com",
            PasswordHash = "h",
            PasswordSalt = "s",
        };

        Action act = () => user.Username = (string?)value!;
        act.Should().Throw<ArgumentException>();
    }

    [Theory]
    [InlineData(null)]
    [InlineData("")]
    [InlineData("   ")]
    public void Setting_NullOrEmpty_Email_ShouldThrow(object? value)
    {
        var user = new User
        {
            Username = "ok",
            Email = "e@example.com",
            PasswordHash = "h",
            PasswordSalt = "s",
        };

        Action act = () => user.Email = (string?)value!;
        act.Should().Throw<ArgumentException>();
    }

    [Theory]
    [InlineData(null)]
    [InlineData("")]
    [InlineData("   ")]
    public void Setting_NullOrEmpty_PasswordHash_ShouldThrow(object? value)
    {
        var user = new User
        {
            Username = "ok",
            Email = "e@example.com",
            PasswordHash = "h",
            PasswordSalt = "s",
        };

        Action act = () => user.PasswordHash = (string?)value!;
        act.Should().Throw<ArgumentException>();
    }

    [Theory]
    [InlineData(null)]
    [InlineData("")]
    [InlineData("   ")]
    public void Setting_NullOrEmpty_PasswordSalt_ShouldThrow(object? value)
    {
        var user = new User
        {
            Username = "ok",
            Email = "e@example.com",
            PasswordHash = "h",
            PasswordSalt = "s",
        };

        Action act = () => user.PasswordSalt = (string?)value!;
        act.Should().Throw<ArgumentException>();
    }

    [Fact]
    public void Users_WithSameValues_AreNotReferenceEqual()
    {
        // Arrange
        var user1 = new User
        {
            Id = 1,
            Username = "testuser",
            Email = "test@example.com",
            PasswordHash = "hash1",
            PasswordSalt = "salt1",
        };

        var user2 = new User
        {
            Id = 1,
            Username = "testuser",
            Email = "test@example.com",
            PasswordHash = "hash1",
            PasswordSalt = "salt1",
        };

        var user3 = new User
        {
            Id = 2,
            Username = "testuser",
            Email = "test@example.com",
            PasswordHash = "hash1",
            PasswordSalt = "salt1",
        };

        // Assert - Reference equality (and default Equals)
        user1.Should().NotBeSameAs(user2);
        user1.Should().NotBeSameAs(user3);

        user1.Equals(user2).Should().BeFalse();

        // Note: For value equality, you'd need to override Equals and GetHashCode
        // This test ensures the class works as expected for Entity Framework
    }
}
