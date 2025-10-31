"""JWT token generator implementation."""

import base64
import uuid
from datetime import datetime, timedelta
from typing import Optional

import jwt

from app.modules.auth.app.ports import TokenGenerator
from app.modules.auth.domain.user import User


class JWTTokenGenerator(TokenGenerator):
    """JWT token generator using PyJWT."""

    def __init__(
        self,
        key: str,
        issuer: str,
        audience: str,
        access_token_lifetime_hours: int,
    ):
        """Initialize JWT generator."""
        self.signing_key = self._decode_key(key)
        self.issuer = issuer
        self.audience = audience
        self.access_token_lifetime_hours = access_token_lifetime_hours

        # Validate key length
        if len(self.signing_key) < 32:
            raise ValueError("jwt: key must be at least 32 bytes after decoding")

        # Validate options
        if not issuer or not issuer.strip():
            raise ValueError("jwt: issuer must be provided")
        if not audience or not audience.strip():
            raise ValueError("jwt: audience must be provided")
        if access_token_lifetime_hours < 1 or access_token_lifetime_hours > 168:
            raise ValueError("jwt: access token lifetime must be between 1 and 168 hours")

    def generate_token(self, user: User) -> str:
        """Generate a signed JWT for the user."""
        if user is None:
            raise ValueError("user must not be None")

        now = datetime.utcnow()
        expires_at = now + timedelta(hours=self.access_token_lifetime_hours)

        claims = {
            "name": user.username,
            "nameid": str(user.id) if user.id else "",
            "sub": str(user.id) if user.id else "",
            "iss": self.issuer,
            "aud": self.audience,
            "iat": int(now.timestamp()),
            "exp": int(expires_at.timestamp()),
            "jti": str(uuid.uuid4()),
        }

        token = jwt.encode(claims, self.signing_key, algorithm="HS256")
        return token

    def _decode_key(self, key: str) -> bytes:
        """Decode key from base64 or plain text format."""
        if key.startswith("base64:"):
            decoded = base64.b64decode(key[7:])
            return decoded
        return key.encode("utf-8")

