"""HMAC-SHA512 password hasher implementation."""

import base64
import hmac
import os
import secrets
from hashlib import sha512

from app.modules.auth.app.ports import PasswordHasher

DEFAULT_SALT_SIZE = 64


class HMACPasswordHasher(PasswordHasher):
    """HMAC-SHA512 password hasher matching Go implementation."""

    def __init__(self, salt_size: int = DEFAULT_SALT_SIZE):
        """Initialize hasher with salt size."""
        self.salt_size = salt_size

    def hash_password(self, password: str) -> tuple[str, str]:
        """Generate a base64-encoded hash and salt."""
        if not password:
            raise ValueError("password cannot be empty")

        # Generate random salt
        salt = secrets.token_bytes(self.salt_size)

        # Compute HMAC-SHA512
        mac = hmac.new(salt, password.encode("utf-8"), sha512)
        hash_bytes = mac.digest()

        # Encode as base64
        hash_str = base64.b64encode(hash_bytes).decode("utf-8")
        salt_str = base64.b64encode(salt).decode("utf-8")

        return hash_str, salt_str

    def verify_password(self, password: str, stored_hash: str, stored_salt: str) -> bool:
        """Verify password against stored hash and salt."""
        if not password:
            raise ValueError("password cannot be empty")
        if not stored_hash or not stored_salt:
            raise ValueError("stored hash and salt must be provided")

        try:
            # Decode salt and hash
            decoded_salt = base64.b64decode(stored_salt)
            decoded_hash = base64.b64decode(stored_hash)

            # Recompute hash
            mac = hmac.new(decoded_salt, password.encode("utf-8"), sha512)
            computed_hash = mac.digest()

            # Constant-time comparison
            return hmac.compare_digest(computed_hash, decoded_hash)
        except Exception:
            return False

