using MySvelteApp.Server.Domain.Entities;

namespace MySvelteApp.Server.Application.Common.Interfaces;

public interface IJwtTokenGenerator
{
    string GenerateToken(User user);
}
