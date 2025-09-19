using FluentAssertions;
using Microsoft.EntityFrameworkCore;
using MySvelteApp.Server.Domain.Entities;
using MySvelteApp.Server.Infrastructure.Persistence;
using MySvelteApp.Server.Infrastructure.Persistence.Repositories;
using MySvelteApp.Server.Tests.TestUtilities;
using Xunit;

namespace MySvelteApp.Server.Tests.Infrastructure.Persistence;

public class UserRepositoryTests : IDisposable
{
    private readonly string _dbName;
    private readonly AppDbContext _dbContext;
    private readonly UserRepository _userRepository;

    public UserRepositoryTests()
    {
        _dbName = $"UserRepositoryTests_{Guid.NewGuid():N}";
        _dbContext = TestHelper.CreateInMemoryDbContext(_dbName);
        _userRepository = new UserRepository(_dbContext);
    }

    public void Dispose()
    {
        _dbContext.Database.EnsureDeleted();
        _dbContext.Dispose();
    }

    [Fact]
    public async Task GetByUsernameAsync_WithExistingUsername_ShouldReturnUser()
    {
        // Arrange
        var user = TestData.Users.ValidUser;
        await TestHelper.SeedUsersAsync(_dbContext, user);

        // Act
        var result = await _userRepository.GetByUsernameAsync(user.Username);

        // Assert
        result.Should().NotBeNull();
        result!.Id.Should().Be(user.Id);
        result.Username.Should().Be(user.Username);
        result.Email.Should().Be(user.Email);
        result.PasswordHash.Should().Be(user.PasswordHash);
        result.PasswordSalt.Should().Be(user.PasswordSalt);
    }

    [Fact]
    public async Task GetByUsernameAsync_WithNonexistentUsername_ShouldReturnNull()
    {
        // Arrange
        var nonexistentUsername = "nonexistentuser";

        // Act
        var result = await _userRepository.GetByUsernameAsync(nonexistentUsername);

        // Assert
        result.Should().BeNull();
    }

    [Fact]
    public async Task GetByUsernameAsync_WithCaseSensitiveUsername_ShouldReturnNull()
    {
        // Arrange
        var user = TestData.Users.ValidUser;
        await TestHelper.SeedUsersAsync(_dbContext, user);

        var wrongCaseUsername = user.Username.ToUpper();

        // Act
        var result = await _userRepository.GetByUsernameAsync(wrongCaseUsername);

        // Assert
        result.Should().BeNull();
    }

    [Fact]
    public async Task UsernameExistsAsync_WithExistingUsername_ShouldReturnTrue()
    {
        // Arrange
        var user = TestData.Users.ValidUser;
        await TestHelper.SeedUsersAsync(_dbContext, user);

        // Act
        var result = await _userRepository.UsernameExistsAsync(user.Username);

        // Assert
        result.Should().BeTrue();
    }

    [Fact]
    public async Task UsernameExistsAsync_WithNonexistentUsername_ShouldReturnFalse()
    {
        // Arrange
        var nonexistentUsername = "nonexistentuser";

        // Act
        var result = await _userRepository.UsernameExistsAsync(nonexistentUsername);

        // Assert
        result.Should().BeFalse();
    }

    [Fact]
    public async Task UsernameExistsAsync_WithDifferentCase_ShouldReturnFalse()
    {
        // Arrange
        var user = TestData.Users.ValidUser;
        await TestHelper.SeedUsersAsync(_dbContext, user);

        // Act
        var result = await _userRepository.UsernameExistsAsync(user.Username.ToUpper());

        // Assert
        result.Should().BeFalse();
    }

    [Fact]
    public async Task EmailExistsAsync_WithExistingEmail_ShouldReturnTrue()
    {
        // Arrange
        var user = TestData.Users.ValidUser;
        await TestHelper.SeedUsersAsync(_dbContext, user);

        // Act
        var result = await _userRepository.EmailExistsAsync(user.Email);

        // Assert
        result.Should().BeTrue();
    }

    [Fact]
    public async Task EmailExistsAsync_WithNonexistentEmail_ShouldReturnFalse()
    {
        // Arrange
        var nonexistentEmail = "nonexistent@example.com";

        // Act
        var result = await _userRepository.EmailExistsAsync(nonexistentEmail);

        // Assert
        result.Should().BeFalse();
    }

    [Fact]
    public async Task EmailExistsAsync_WithDifferentCaseEmail_ShouldReturnTrue()
    {
        // Arrange
        var user = TestData.Users.ValidUser;
        await TestHelper.SeedUsersAsync(_dbContext, user);

        var differentCaseEmail = user.Email.ToUpper();

        // Act
        var result = await _userRepository.EmailExistsAsync(differentCaseEmail);

        // Assert
        result.Should().BeTrue();
    }

    [Fact]
    public async Task AddAsync_WithValidUser_ShouldPersistUser()
    {
        // Arrange
        var newUser = TestData.Users.NewUser;

        // Act
        await _userRepository.AddAsync(newUser);

        // Assert
        var persistedUser = await _dbContext.Users.FirstOrDefaultAsync(u =>
            u.Username == newUser.Username
        );
        persistedUser.Should().NotBeNull();
        persistedUser!.Username.Should().Be(newUser.Username);
        persistedUser.Email.Should().Be(newUser.Email);
        persistedUser.PasswordHash.Should().Be(newUser.PasswordHash);
        persistedUser.PasswordSalt.Should().Be(newUser.PasswordSalt);
        persistedUser.Id.Should().NotBe(0); // Should have been assigned by EF
    }

    [Fact]
    public async Task AddAsync_ShouldSaveChangesToDatabase()
    {
        // Arrange
        var testDbName = $"PersistenceTestDb_{Guid.NewGuid():N}";
        var newUser = TestData.Users.NewUser;

        // Create a fresh context with a known database name for this test
        using var testContext = TestHelper.CreateInMemoryDbContext(testDbName);
        var testRepository = new UserRepository(testContext);

        // Act
        await testRepository.AddAsync(newUser);

        // Assert
        // First verify the user doesn't exist in a different database (isolation test)
        using var differentContext = TestHelper.CreateInMemoryDbContext("DifferentDb");
        var userInDifferentDb = await differentContext.Users.FirstOrDefaultAsync(u =>
            u.Username == newUser.Username
        );
        userInDifferentDb.Should().BeNull(); // Should not exist in different database

        // Then verify the user exists in the same database (persistence test)
        using var sameContext = TestHelper.CreateInMemoryDbContext(testDbName);
        var persistedUser = await sameContext.Users.FirstOrDefaultAsync(u =>
            u.Username == newUser.Username
        );
        persistedUser.Should().NotBeNull(); // Should exist in the same database
        persistedUser!.Username.Should().Be(newUser.Username);
        persistedUser.Email.Should().Be(newUser.Email);
    }

    [Fact]
    public async Task GetByUsernameAsync_AfterAddAsync_ShouldReturnPersistedUser()
    {
        // Arrange
        var newUser = TestData.Users.NewUser;
        await _userRepository.AddAsync(newUser);

        // Act
        var result = await _userRepository.GetByUsernameAsync(newUser.Username);

        // Assert
        result.Should().NotBeNull();
        result!.Username.Should().Be(newUser.Username);
        result.Email.Should().Be(newUser.Email);
    }

    [Fact]
    public async Task UsernameExistsAsync_AfterAddAsync_ShouldReturnTrue()
    {
        // Arrange
        var newUser = TestData.Users.NewUser;
        await _userRepository.AddAsync(newUser);

        // Act
        var result = await _userRepository.UsernameExistsAsync(newUser.Username);

        // Assert
        result.Should().BeTrue();
    }

    [Fact]
    public async Task EmailExistsAsync_AfterAddAsync_ShouldReturnTrue()
    {
        // Arrange
        var newUser = TestData.Users.NewUser;
        await _userRepository.AddAsync(newUser);

        // Act
        var result = await _userRepository.EmailExistsAsync(newUser.Email);

        // Assert
        result.Should().BeTrue();
    }

    [Fact]
    public async Task Operations_WithMultipleUsers_ShouldWorkCorrectly()
    {
        // Arrange
        var user1 = TestData.Users.ValidUser;
        var user2 = TestData.Users.AnotherValidUser;
        await TestHelper.SeedUsersAsync(_dbContext, user1, user2);

        // Act & Assert
        var result1 = await _userRepository.GetByUsernameAsync(user1.Username);
        var result2 = await _userRepository.GetByUsernameAsync(user2.Username);

        result1.Should().NotBeNull();
        result2.Should().NotBeNull();
        result1!.Id.Should().NotBe(result2!.Id);

        (await _userRepository.UsernameExistsAsync(user1.Username)).Should().BeTrue();
        (await _userRepository.UsernameExistsAsync(user2.Username)).Should().BeTrue();
        (await _userRepository.EmailExistsAsync(user1.Email)).Should().BeTrue();
        (await _userRepository.EmailExistsAsync(user2.Email)).Should().BeTrue();
    }

    [Fact]
    public async Task Operations_WithCancellationToken_ShouldRespectCancellation()
    {
        // Arrange
        var user = TestData.Users.ValidUser;
        await TestHelper.SeedUsersAsync(_dbContext, user);

        using var cts = new CancellationTokenSource();
        cts.Cancel();

        // Act & Assert
        await Assert.ThrowsAnyAsync<OperationCanceledException>(() =>
            _userRepository.GetByUsernameAsync(user.Username, cts.Token)
        );

        await Assert.ThrowsAnyAsync<OperationCanceledException>(() =>
            _userRepository.UsernameExistsAsync(user.Username, cts.Token)
        );

        await Assert.ThrowsAnyAsync<OperationCanceledException>(() =>
            _userRepository.EmailExistsAsync(user.Email, cts.Token)
        );

        await Assert.ThrowsAnyAsync<OperationCanceledException>(() =>
            _userRepository.AddAsync(TestData.Users.NewUser, cts.Token)
        );
    }

    [Fact]
    public async Task AddAsync_WithDuplicateUsername_ShouldThrowDbUpdateException()
    {
        // Arrange - Use SQLite in-memory which enforces unique constraints
        using var sqliteContext = TestHelper.CreateSqliteInMemoryDbContext();
        var sqliteRepository = new UserRepository(sqliteContext);

        var user1 = TestData.Users.ValidUser;
        var user2 = new User
        {
            Username = user1.Username, // Same username
            Email = "different@example.com",
            PasswordHash = "different_hash",
            PasswordSalt = "different_salt",
        };

        await TestHelper.SeedUsersAsync(sqliteContext, user1);

        // Act & Assert - Should throw due to unique constraint violation
        await Assert.ThrowsAsync<DbUpdateException>(() => sqliteRepository.AddAsync(user2));
    }

    [Fact]
    public async Task AddAsync_WithDuplicateEmail_ShouldThrowDbUpdateException()
    {
        // Arrange - Use SQLite in-memory which enforces unique constraints
        using var sqliteContext = TestHelper.CreateSqliteInMemoryDbContext();
        var sqliteRepository = new UserRepository(sqliteContext);

        var user1 = TestData.Users.ValidUser;
        var user2 = new User
        {
            Username = "differentuser",
            Email = user1.Email, // Same email
            PasswordHash = "different_hash",
            PasswordSalt = "different_salt",
        };

        await TestHelper.SeedUsersAsync(sqliteContext, user1);

        // Act & Assert - Should throw due to unique constraint violation
        await Assert.ThrowsAsync<DbUpdateException>(() => sqliteRepository.AddAsync(user2));
    }
}
