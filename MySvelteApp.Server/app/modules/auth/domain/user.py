"""User domain entity with business rules and validation."""

import re
from dataclasses import dataclass
from datetime import datetime
from typing import Optional

# Domain constants
MIN_USERNAME_LENGTH = 3
MAX_USERNAME_LENGTH = 64
MAX_EMAIL_LENGTH = 320

# Validation patterns
USERNAME_REGEX = re.compile(r"^[a-zA-Z0-9_]+$")
EMAIL_REGEX = re.compile(r"^[^\s@]+@[^\s@.]+\.[^\s@.]+$")


def valid_username(username: str) -> bool:
    """Check if username satisfies domain rules."""
    username = username.strip()
    if len(username) < MIN_USERNAME_LENGTH or len(username) > MAX_USERNAME_LENGTH:
        return False
    return bool(USERNAME_REGEX.match(username))


def valid_email(email: str) -> bool:
    """Check if email satisfies domain rules."""
    email = email.strip()
    if len(email) == 0 or len(email) > MAX_EMAIL_LENGTH:
        return False
    if ".." in email:
        return False
    return bool(EMAIL_REGEX.match(email))


@dataclass
class User:
    """User domain entity."""

    id: Optional[int] = None
    username: str = ""
    email: str = ""
    password_hash: str = ""
    password_salt: str = ""
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None

    @classmethod
    def create(
        cls, username: str, email: str, password_hash: str, password_salt: str
    ) -> "User":
        """Create a new User with validation."""
        username = username.strip()
        if len(username) == 0:
            raise ValueError("username cannot be empty")
        if len(username) < MIN_USERNAME_LENGTH:
            raise ValueError(f"username must be at least {MIN_USERNAME_LENGTH} characters")
        if len(username) > MAX_USERNAME_LENGTH:
            raise ValueError(f"username must not exceed {MAX_USERNAME_LENGTH} characters")
        if not USERNAME_REGEX.match(username):
            raise ValueError("username may only contain letters, numbers, and underscores")

        if len(password_hash) == 0:
            raise ValueError("password hash cannot be empty")
        if len(password_salt) == 0:
            raise ValueError("password salt cannot be empty")

        trimmed_email = email.strip()
        if len(trimmed_email) == 0:
            raise ValueError("email cannot be empty")
        normalized_email = trimmed_email.lower()
        if len(normalized_email) > MAX_EMAIL_LENGTH:
            raise ValueError(f"email must not exceed {MAX_EMAIL_LENGTH} characters")
        if not valid_email(normalized_email):
            raise ValueError("email format is invalid")

        return cls(
            username=username,
            email=normalized_email,
            password_hash=password_hash,
            password_salt=password_salt,
        )

