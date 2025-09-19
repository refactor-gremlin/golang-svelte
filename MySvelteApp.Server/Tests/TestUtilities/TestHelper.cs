using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Options;
using MySvelteApp.Server.Domain.Entities;
using MySvelteApp.Server.Infrastructure.Authentication;
using MySvelteApp.Server.Infrastructure.Persistence;

namespace MySvelteApp.Server.Tests.TestUtilities;

public static class TestHelper
{
    public static AppDbContext CreateInMemoryDbContext(string? dbName = null)
    {
        var databaseName = dbName ?? Guid.NewGuid().ToString();
        var options = new DbContextOptionsBuilder<AppDbContext>()
            .UseInMemoryDatabase(databaseName)
            .Options;

        return new AppDbContext(options);
    }

    public static AppDbContext CreateSqliteInMemoryDbContext()
    {
        var connection = new SqliteConnection("Filename=:memory:");
        connection.Open();

        var options = new DbContextOptionsBuilder<AppDbContext>().UseSqlite(connection).Options;

        var context = new AppDbContext(options);
        context.Database.EnsureCreated();

        return context;
    }

    public static IOptions<JwtOptions> CreateJwtOptions() =>
        Options.Create(
            new JwtOptions
            {
                Key = TestData.Jwt.ValidKey,
                Issuer = TestData.Jwt.ValidIssuer,
                Audience = TestData.Jwt.ValidAudience,
                AccessTokenLifetimeHours = TestData.Jwt.ValidLifetimeHours,
            }
        );

    public static async Task SeedUsersAsync(AppDbContext context, params User[] users)
    {
        await context.Users.AddRangeAsync(users);
        await context.SaveChangesAsync();
    }

    public static void ClearDatabase(AppDbContext context)
    {
        context.Database.EnsureDeleted();
    }
}
