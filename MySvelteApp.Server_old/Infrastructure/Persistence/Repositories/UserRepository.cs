using Microsoft.EntityFrameworkCore;
using MySvelteApp.Server.Application.Common.Interfaces;
using MySvelteApp.Server.Domain.Entities;

namespace MySvelteApp.Server.Infrastructure.Persistence.Repositories;

public class UserRepository : IUserRepository
{
    private readonly AppDbContext _dbContext;

    public UserRepository(AppDbContext dbContext)
    {
        _dbContext = dbContext;
    }

    public async Task<User?> GetByUsernameAsync(
        string username,
        CancellationToken cancellationToken = default
    )
    {
        return await _dbContext
            .Users.AsNoTracking()
            .FirstOrDefaultAsync(u => u.Username == username, cancellationToken);
    }

    public async Task<bool> UsernameExistsAsync(
        string username,
        CancellationToken cancellationToken = default
    )
    {
        return await _dbContext
            .Users.AsNoTracking()
            .AnyAsync(u => u.Username == username, cancellationToken);
    }

    public async Task<bool> EmailExistsAsync(
        string email,
        CancellationToken cancellationToken = default
    )
    {
        var normalizedEmail = email.Trim().ToLowerInvariant();
        return await _dbContext
            .Users.AsNoTracking()
            .AnyAsync(u => u.Email == normalizedEmail, cancellationToken);
    }

    public async Task AddAsync(User user, CancellationToken cancellationToken = default)
    {
        await _dbContext.Users.AddAsync(user, cancellationToken);
        await _dbContext.SaveChangesAsync(cancellationToken);
    }
}
