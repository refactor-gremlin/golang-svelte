using FluentAssertions;
using Moq;
using MySvelteApp.Server.Application.Authentication;
using MySvelteApp.Server.Application.Authentication.DTOs;
using MySvelteApp.Server.Application.Common.Interfaces;
using MySvelteApp.Server.Domain.Entities;
using MySvelteApp.Server.Tests.TestUtilities;
using Xunit;

namespace MySvelteApp.Server.Tests.Application.Authentication;

public class AuthServiceTests
{
    private readonly Mock<IUserRepository> _userRepositoryMock;
    private readonly Mock<IPasswordHasher> _passwordHasherMock;
    private readonly Mock<IJwtTokenGenerator> _jwtTokenGeneratorMock;
    private readonly AuthService _authService;

    public AuthServiceTests()
    {
        _userRepositoryMock = new Mock<IUserRepository>();
        _passwordHasherMock = new Mock<IPasswordHasher>();
        _jwtTokenGeneratorMock = new Mock<IJwtTokenGenerator>();

        _authService = new AuthService(
            _userRepositoryMock.Object,
            _passwordHasherMock.Object,
            _jwtTokenGeneratorMock.Object
        );
    }

    public class RegisterAsyncTests : AuthServiceTests
    {
        [Fact]
        public async Task RegisterAsync_WithValidRequest_ShouldReturnSuccess()
        {
            // Arrange
            var request = TestData.Requests.Authentication.ValidRegisterRequest;
            const string expectedToken = "generated.jwt.token";

            _userRepositoryMock
                .Setup(x => x.UsernameExistsAsync(request.Username, default))
                .ReturnsAsync(false);
            _userRepositoryMock
                .Setup(x => x.EmailExistsAsync(request.Email, default))
                .ReturnsAsync(false);
            _passwordHasherMock
                .Setup(x => x.HashPassword(request.Password))
                .Returns(("hashed_password", "salt"));
            _jwtTokenGeneratorMock
                .Setup(x => x.GenerateToken(It.IsAny<User>()))
                .Returns(expectedToken);

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeTrue();
            result.Token.Should().Be(expectedToken);
            result.Username.Should().Be(request.Username);
            result.UserId.Should().NotBeNull();
            result.ErrorMessage.Should().BeNull();
            result.ErrorType.Should().Be(AuthErrorType.None);

            _userRepositoryMock.Verify(
                x =>
                    x.AddAsync(
                        It.Is<User>(u =>
                            u.Username == request.Username
                            && u.Email == request.Email.ToLowerInvariant()
                            && u.PasswordHash == "hashed_password"
                            && u.PasswordSalt == "salt"
                        ),
                        default
                    ),
                Times.Once
            );
        }

        [Fact]
        public async Task RegisterAsync_ShouldTrimUsername()
        {
            var request = new RegisterRequest
            {
                Username = "  testuser  ",
                Email = "TEST@EXAMPLE.COM",
                Password = "ValidPassword123",
            };
            const string expectedToken = "generated.jwt.token";

            _userRepositoryMock
                .Setup(x => x.UsernameExistsAsync("testuser", default))
                .ReturnsAsync(false);
            _userRepositoryMock
                .Setup(x => x.EmailExistsAsync("test@example.com", default))
                .ReturnsAsync(false);
            _passwordHasherMock
                .Setup(x => x.HashPassword(request.Password))
                .Returns(("hashed_password", "salt"));
            _jwtTokenGeneratorMock
                .Setup(x => x.GenerateToken(It.IsAny<User>()))
                .Returns(expectedToken);

            var result = await _authService.RegisterAsync(request);

            result.Success.Should().BeTrue();
            result.Username.Should().Be("testuser");
            _userRepositoryMock.Verify(
                x => x.AddAsync(It.Is<User>(u => u.Username == "testuser"), default),
                Times.Once
            );
        }

        [Fact]
        public async Task RegisterAsync_WithExistingUsername_ShouldReturnConflictError()
        {
            // Arrange
            var request = TestData.Requests.Authentication.ValidRegisterRequest;

            _userRepositoryMock
                .Setup(x => x.UsernameExistsAsync(request.Username, default))
                .ReturnsAsync(true);

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("username is already taken");
            result.ErrorType.Should().Be(AuthErrorType.Conflict);
            result.Token.Should().BeNull();
            result.UserId.Should().BeNull();

            _userRepositoryMock.Verify(
                x => x.EmailExistsAsync(It.IsAny<string>(), default),
                Times.Never
            );
            _passwordHasherMock.Verify(x => x.HashPassword(It.IsAny<string>()), Times.Never);
            _userRepositoryMock.Verify(x => x.AddAsync(It.IsAny<User>(), default), Times.Never);
        }

        [Fact]
        public async Task RegisterAsync_WithExistingEmail_ShouldReturnConflictError()
        {
            // Arrange
            var request = TestData.Requests.Authentication.ValidRegisterRequest;

            _userRepositoryMock
                .Setup(x => x.UsernameExistsAsync(request.Username, default))
                .ReturnsAsync(false);
            _userRepositoryMock
                .Setup(x => x.EmailExistsAsync(request.Email, default))
                .ReturnsAsync(true);

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("email is already registered");
            result.ErrorType.Should().Be(AuthErrorType.Conflict);

            _passwordHasherMock.Verify(x => x.HashPassword(It.IsAny<string>()), Times.Never);
            _userRepositoryMock.Verify(x => x.AddAsync(It.IsAny<User>(), default), Times.Never);
        }

        [Fact]
        public async Task RegisterAsync_WithInvalidUsername_ShouldReturnValidationError()
        {
            // Arrange
            var request = TestData.Requests.Authentication.InvalidUsernameRequest;

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("at least 3 characters");
            result.ErrorType.Should().Be(AuthErrorType.Validation);

            _userRepositoryMock.Verify(
                x => x.UsernameExistsAsync(It.IsAny<string>(), default),
                Times.Never
            );
        }

        [Fact]
        public async Task RegisterAsync_WithInvalidEmail_ShouldReturnValidationError()
        {
            // Arrange
            var request = TestData.Requests.Authentication.InvalidEmailRequest;

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("valid email");
            result.ErrorType.Should().Be(AuthErrorType.Validation);
        }

        [Fact]
        public async Task RegisterAsync_WithWeakPassword_ShouldReturnValidationError()
        {
            // Arrange
            var request = TestData.Requests.Authentication.WeakPasswordRequest;

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("uppercase letter");
            result.ErrorType.Should().Be(AuthErrorType.Validation);
        }

        [Fact]
        public async Task RegisterAsync_WithEmptyUsername_ShouldReturnValidationError()
        {
            // Arrange
            var request = new RegisterRequest
            {
                Username = "",
                Email = "test@example.com",
                Password = "ValidPassword123",
            };

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("Username is required");
            result.ErrorType.Should().Be(AuthErrorType.Validation);
        }

        [Fact]
        public async Task RegisterAsync_WithNullEmail_ShouldReturnValidationError()
        {
            // Arrange
            var request = new RegisterRequest
            {
                Username = "testuser",
                Email = null!,
                Password = "ValidPassword123",
            };

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("Email is required");
            result.ErrorType.Should().Be(AuthErrorType.Validation);
        }

        [Fact]
        public async Task RegisterAsync_WithNullPassword_ShouldReturnValidationError()
        {
            // Arrange
            var request = new RegisterRequest
            {
                Username = "testuser",
                Email = "test@example.com",
                Password = null!,
            };

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("Password is required");
            result.ErrorType.Should().Be(AuthErrorType.Validation);
        }

        [Fact]
        public async Task RegisterAsync_WithUsernameContainingInvalidChars_ShouldReturnValidationError()
        {
            // Arrange
            var request = new RegisterRequest
            {
                Username = "test@user",
                Email = "test@example.com",
                Password = "ValidPassword123",
            };

            // Act
            var result = await _authService.RegisterAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("letters, numbers, and underscores");
            result.ErrorType.Should().Be(AuthErrorType.Validation);
        }

        [Fact]
        public async Task RegisterAsync_ShouldNormalizeEmailAndPersistUser()
        {
            // Arrange
            var request = new RegisterRequest
            {
                Username = "testuser",
                Email = "Test@Example.COM",
                Password = "ValidPassword123",
            };

            _userRepositoryMock
                .Setup(x => x.UsernameExistsAsync("testuser", default))
                .ReturnsAsync(false);
            _userRepositoryMock
                .Setup(x => x.EmailExistsAsync("test@example.com", default))
                .ReturnsAsync(false);
            _passwordHasherMock
                .Setup(x => x.HashPassword(request.Password))
                .Returns(("hashed_password", "salt"));
            _jwtTokenGeneratorMock.Setup(x => x.GenerateToken(It.IsAny<User>())).Returns("token");

            // Act
            await _authService.RegisterAsync(request);

            // Assert
            _userRepositoryMock.Verify(
                x =>
                    x.AddAsync(
                        It.Is<User>(u => u.Username == "testuser" && u.Email == "test@example.com"),
                        default
                    ),
                Times.Once
            );
        }
    }

    public class LoginAsyncTests : AuthServiceTests
    {
        [Fact]
        public async Task LoginAsync_WithValidCredentials_ShouldReturnSuccess()
        {
            // Arrange
            var request = TestData.Requests.Authentication.ValidLoginRequest;
            var user = TestData.Users.ValidUser;
            const string expectedToken = "generated.jwt.token";

            _userRepositoryMock
                .Setup(x => x.GetByUsernameAsync(request.Username, default))
                .ReturnsAsync(user);
            _passwordHasherMock
                .Setup(x =>
                    x.VerifyPassword(request.Password, user.PasswordHash, user.PasswordSalt)
                )
                .Returns(true);
            _jwtTokenGeneratorMock.Setup(x => x.GenerateToken(user)).Returns(expectedToken);

            // Act
            var result = await _authService.LoginAsync(request);

            // Assert
            result.Success.Should().BeTrue();
            result.Token.Should().Be(expectedToken);
            result.UserId.Should().Be(user.Id);
            result.Username.Should().Be(user.Username);
        }

        [Fact]
        public async Task LoginAsync_WithNonexistentUsername_ShouldReturnUnauthorized()
        {
            // Arrange
            var request = TestData.Requests.Authentication.ValidLoginRequest;

            _userRepositoryMock
                .Setup(x => x.GetByUsernameAsync(request.Username, default))
                .ReturnsAsync((User?)null);

            // Act
            var result = await _authService.LoginAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("Invalid username or password");
            result.ErrorType.Should().Be(AuthErrorType.Unauthorized);

            _passwordHasherMock.Verify(
                x => x.VerifyPassword(It.IsAny<string>(), It.IsAny<string>(), It.IsAny<string>()),
                Times.Never
            );
        }

        [Fact]
        public async Task LoginAsync_WithWrongPassword_ShouldReturnUnauthorized()
        {
            // Arrange
            var request = TestData.Requests.Authentication.ValidLoginRequest;
            var user = TestData.Users.ValidUser;

            _userRepositoryMock
                .Setup(x => x.GetByUsernameAsync(request.Username, default))
                .ReturnsAsync(user);
            _passwordHasherMock
                .Setup(x =>
                    x.VerifyPassword(request.Password, user.PasswordHash, user.PasswordSalt)
                )
                .Returns(false);

            // Act
            var result = await _authService.LoginAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("Invalid username or password");
            result.ErrorType.Should().Be(AuthErrorType.Unauthorized);

            _jwtTokenGeneratorMock.Verify(x => x.GenerateToken(It.IsAny<User>()), Times.Never);
        }

        [Fact]
        public async Task LoginAsync_WithEmptyUsername_ShouldReturnValidationError()
        {
            // Arrange
            var request = TestData.Requests.Authentication.EmptyUsernameRequest;

            // Act
            var result = await _authService.LoginAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("Username is required");
            result.ErrorType.Should().Be(AuthErrorType.Validation);

            _userRepositoryMock.Verify(
                x => x.GetByUsernameAsync(It.IsAny<string>(), default),
                Times.Never
            );
        }

        [Fact]
        public async Task LoginAsync_WithEmptyPassword_ShouldReturnValidationError()
        {
            // Arrange
            var request = new LoginRequest { Username = "testuser", Password = "" };

            // Act
            var result = await _authService.LoginAsync(request);

            // Assert
            result.Success.Should().BeFalse();
            result.ErrorMessage.Should().Contain("Password is required");
            result.ErrorType.Should().Be(AuthErrorType.Validation);
        }

        [Fact]
        public async Task LoginAsync_ShouldTrimUsername()
        {
            // Arrange
            var request = new LoginRequest { Username = "  testuser  ", Password = "password" };
            var user = TestData.Users.ValidUser;

            _userRepositoryMock
                .Setup(x => x.GetByUsernameAsync("testuser", default))
                .ReturnsAsync(user);
            _passwordHasherMock
                .Setup(x =>
                    x.VerifyPassword(request.Password, user.PasswordHash, user.PasswordSalt)
                )
                .Returns(true);
            _jwtTokenGeneratorMock.Setup(x => x.GenerateToken(user)).Returns("token");

            // Act
            await _authService.LoginAsync(request);

            // Assert
            _userRepositoryMock.Verify(x => x.GetByUsernameAsync("testuser", default), Times.Once);
        }
    }
}
