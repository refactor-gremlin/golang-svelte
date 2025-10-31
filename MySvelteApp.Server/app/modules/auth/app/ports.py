"""Interface definitions for authentication dependencies."""

from abc import ABC, abstractmethod
from typing import Optional

from app.modules.auth.domain.user import User


class UserRepository(ABC):
    """User persistence interface."""

    @abstractmethod
    def add(self, user: User) -> None:
        """Persist a new user."""
        pass

    @abstractmethod
    def get_by_username(self, username: str) -> Optional[User]:
        """Retrieve a user by username. Returns None if not found."""
        pass

    @abstractmethod
    def username_exists(self, username: str) -> bool:
        """Check if a username is already taken."""
        pass

    @abstractmethod
    def email_exists(self, email: str) -> bool:
        """Check if an email address is already registered."""
        pass


class PasswordHasher(ABC):
    """Password hashing interface."""

    @abstractmethod
    def hash_password(self, password: str) -> tuple[str, str]:
        """Create a secure hash and salt for the given password. Returns (hash, salt)."""
        pass

    @abstractmethod
    def verify_password(self, password: str, stored_hash: str, stored_salt: str) -> bool:
        """Verify if the provided password matches the stored hash and salt."""
        pass


class TokenGenerator(ABC):
    """Token generation interface."""

    @abstractmethod
    def generate_token(self, user: User) -> str:
        """Create a new authentication token for the given user."""
        pass

