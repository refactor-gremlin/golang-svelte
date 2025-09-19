using System.Security.Cryptography;
using System.Text;
using MySvelteApp.Server.Application.Common.Interfaces;

namespace MySvelteApp.Server.Infrastructure.Security;

public class PasswordHasher : IPasswordHasher
{
    private const int SaltSize = 16;
    private const int KeySize = 64; // bytes
    private const int Iterations = 100_000; // consider higher or configurable

    public (string Hash, string Salt) HashPassword(string password)
    {
        var salt = RandomNumberGenerator.GetBytes(SaltSize);
        var hash = Rfc2898DeriveBytes.Pbkdf2(
            password,
            salt,
            Iterations,
            HashAlgorithmName.SHA512,
            KeySize
        );
        return (Convert.ToBase64String(hash), Convert.ToBase64String(salt));
    }

    public bool VerifyPassword(string password, string hash, string salt)
    {
        var saltBytes = Convert.FromBase64String(salt);
        var storedHash = Convert.FromBase64String(hash);
        var computed = Rfc2898DeriveBytes.Pbkdf2(
            password,
            saltBytes,
            Iterations,
            HashAlgorithmName.SHA512,
            KeySize
        );
        return CryptographicOperations.FixedTimeEquals(computed, storedHash);
    }
}
