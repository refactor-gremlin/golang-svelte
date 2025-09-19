using MySvelteApp.Server.Application.Authentication.DTOs;
using MySvelteApp.Server.Domain.Entities;

namespace MySvelteApp.Server.Tests.TestUtilities;

public static class TestData
{
    public static class Users
    {
        public static User ValidUser =>
            new()
            {
                Username = "testuser",
                Email = "test@example.com",
                PasswordHash = "hashed_password",
                PasswordSalt = "salt_value",
            };

        public static User AnotherValidUser =>
            new()
            {
                Username = "anotheruser",
                Email = "another@example.com",
                PasswordHash = "another_hash",
                PasswordSalt = "another_salt",
            };

        public static User NewUser =>
            new()
            {
                Username = "newuser",
                Email = "new@example.com",
                PasswordHash = "new_hash",
                PasswordSalt = "new_salt",
            };
    }

    public static class Requests
    {
        public static class Authentication
        {
            public static RegisterRequest ValidRegisterRequest =>
                new()
                {
                    Username = "testuser",
                    Email = "test@example.com",
                    Password = "ValidPassword123",
                };

            public static RegisterRequest InvalidUsernameRequest =>
                new()
                {
                    Username = "us",
                    Email = "test@example.com",
                    Password = "ValidPassword123",
                };

            public static RegisterRequest InvalidEmailRequest =>
                new()
                {
                    Username = "testuser",
                    Email = "invalid-email",
                    Password = "ValidPassword123",
                };

            public static RegisterRequest WeakPasswordRequest =>
                new()
                {
                    Username = "testuser",
                    Email = "test@example.com",
                    Password = "password",
                };

            public static LoginRequest ValidLoginRequest =>
                new() { Username = "testuser", Password = "ValidPassword123" };

            public static LoginRequest EmptyUsernameRequest =>
                new() { Username = "", Password = "ValidPassword123" };

            public static RegisterRequest InvalidDoubleDotEmailRequest =>
                new()
                {
                    Username = "testuser",
                    Email = "test..user@example.com",
                    Password = "ValidPassword123",
                };

            public static RegisterRequest LeadingDotDomainEmailRequest =>
                new()
                {
                    Username = "testuser",
                    Email = "user@.example.com",
                    Password = "ValidPassword123",
                };

            public static RegisterRequest TrailingDotDomainEmailRequest =>
                new()
                {
                    Username = "testuser",
                    Email = "user@example.com.",
                    Password = "ValidPassword123",
                };

            public static RegisterRequest EmailWithSpacesRequest =>
                new()
                {
                    Username = "testuser",
                    Email = "  test@example.com  ",
                    Password = "ValidPassword123",
                };
        }
    }

    public static class Jwt
    {
        public const string ValidKey = "base64:YWJjZGVmZ2hpamsxMjM0NTY3ODkwYWJjZGVmZ2hpams="; // 32+ bytes when decoded
        public const string ShortKey = "short";
        public const string ValidIssuer = "TestIssuer";
        public const string ValidAudience = "TestAudience";
        public const int ValidLifetimeHours = 24;
        public const string ValidPlainTextKey =
            "ThisIsAVeryLongPlainTextKeyThatIsDefinitelyLongerThan32CharactersForTesting";
    }
}
