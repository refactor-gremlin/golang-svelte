"""Authentication service with business logic."""

import re
import string
from typing import Optional

from app.modules.auth.app.errors import (
    ConflictError,
    UnauthorizedError,
    ValidationError,
)
from app.modules.auth.app.ports import PasswordHasher, TokenGenerator, UserRepository
from app.modules.auth.domain.user import (
    MAX_EMAIL_LENGTH,
    MAX_USERNAME_LENGTH,
    MIN_USERNAME_LENGTH,
    User,
    valid_email,
    valid_username,
)

MIN_PASSWORD_LENGTH = 8
MAX_PASSWORD_LENGTH = 512


class AuthService:
    """Authentication service."""

    def __init__(
        self,
        user_repository: UserRepository,
        password_hasher: PasswordHasher,
        token_generator: TokenGenerator,
    ):
        """Initialize auth service with dependencies."""
        self.user_repository = user_repository
        self.password_hasher = password_hasher
        self.token_generator = token_generator

    def register(self, username: str, email: str, password: str) -> dict:
        """Register a new user account."""
        # Validate input
        self._validate_register(username, email, password)

        trimmed_username = username.strip()
        normalized_email = email.strip().lower()

        # Check uniqueness
        if self.user_repository.username_exists(trimmed_username):
            raise ConflictError("This username is already taken. Please choose a different one.")

        if self.user_repository.email_exists(normalized_email):
            raise ConflictError("This email is already registered. Please use a different email address.")

        # Hash password
        password_hash, password_salt = self.password_hasher.hash_password(password)

        # Create user entity
        user = User.create(trimmed_username, normalized_email, password_hash, password_salt)

        # Persist user
        self.user_repository.add(user)

        # Generate token
        token = self.token_generator.generate_token(user)

        return {
            "token": token,
            "userId": user.id or 0,
            "username": user.username,
        }

    def login(self, username: str, password: str) -> dict:
        """Authenticate a user with username and password."""
        # Validate input
        if not username or not username.strip():
            raise ValidationError("Username is required.")
        if not password or not password.strip():
            raise ValidationError("Password is required.")

        trimmed_username = username.strip()

        # Get user
        user = self.user_repository.get_by_username(trimmed_username)
        if user is None:
            raise UnauthorizedError("Invalid username or password. Please check your credentials and try again.")

        # Verify password
        is_valid = self.password_hasher.verify_password(password, user.password_hash, user.password_salt)
        if not is_valid:
            raise UnauthorizedError("Invalid username or password. Please check your credentials and try again.")

        # Generate token
        token = self.token_generator.generate_token(user)

        return {
            "token": token,
            "userId": user.id or 0,
            "username": user.username,
        }

    def _validate_register(self, username: str, email: str, password: str) -> None:
        """Validate registration input."""
        username = username.strip()
        if not username:
            raise ValidationError("Username is required.")
        if len(username) < MIN_USERNAME_LENGTH:
            raise ValidationError(f"Username must be at least {MIN_USERNAME_LENGTH} characters long.")
        if len(username) > MAX_USERNAME_LENGTH:
            raise ValidationError(f"Username must not exceed {MAX_USERNAME_LENGTH} characters.")
        if not valid_username(username):
            raise ValidationError("Username can only contain letters, numbers, and underscores.")

        email = email.strip()
        if not email:
            raise ValidationError("Email is required.")
        if len(email) > MAX_EMAIL_LENGTH:
            raise ValidationError(f"Email must not exceed {MAX_EMAIL_LENGTH} characters.")
        if not valid_email(email):
            raise ValidationError("Please enter a valid email address.")

        password = password.strip()
        if not password:
            raise ValidationError("Password is required.")
        if len(password) < MIN_PASSWORD_LENGTH:
            raise ValidationError("Password must be at least 8 characters long.")
        if len(password) > MAX_PASSWORD_LENGTH:
            raise ValidationError("Password must not exceed 512 characters.")
        if not self._password_meets_requirements(password):
            raise ValidationError("Password must contain at least one uppercase letter, one lowercase letter, and one number.")

    def _password_meets_requirements(self, password: str) -> bool:
        """Check if password meets complexity requirements."""
        has_upper = any(c.isupper() for c in password)
        has_lower = any(c.islower() for c in password)
        has_digit = any(c.isdigit() for c in password)
        return has_upper and has_lower and has_digit

