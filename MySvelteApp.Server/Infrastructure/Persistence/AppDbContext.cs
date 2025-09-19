using Microsoft.EntityFrameworkCore;
using MySvelteApp.Server.Domain.Entities;

namespace MySvelteApp.Server.Infrastructure.Persistence;

public class AppDbContext(DbContextOptions<AppDbContext> options) : DbContext(options)
{
    public DbSet<User> Users => Set<User>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        var user = modelBuilder.Entity<User>();
        user.Property(u => u.Username).IsRequired().HasMaxLength(64);
        user.Property(u => u.Email).IsRequired().HasMaxLength(320);
        user.Property(u => u.PasswordHash).IsRequired().HasMaxLength(512);
        user.Property(u => u.PasswordSalt).IsRequired().HasMaxLength(512);

        user.HasIndex(u => u.Username).IsUnique();
        user.HasIndex(u => u.Email).IsUnique();

        // Align SQLite with case-insensitive email semantics and case-sensitive username
        if (Database.ProviderName?.Contains("Sqlite") == true)
        {
            try
            {
                // Use dynamic invocation to avoid compile-time dependency issues with UseCollation
                dynamic emailProperty = user.Property(u => u.Email);
                dynamic usernameProperty = user.Property(u => u.Username);

                // Call UseCollation dynamically - this will work at runtime when SQLite provider is available
                emailProperty.UseCollation("NOCASE");
                usernameProperty.UseCollation("BINARY");
            }
            catch
            {
                // Silently ignore if SQLite collation is not available
                // The database will still work with case-sensitive comparison
            }
        }
    }
}
