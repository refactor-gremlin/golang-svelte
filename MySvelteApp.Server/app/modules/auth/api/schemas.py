"""Pydantic schemas for auth API."""

from pydantic import BaseModel


class RegisterRequest(BaseModel):
    """Register request schema."""

    username: str
    email: str
    password: str


class LoginRequest(BaseModel):
    """Login request schema."""

    username: str
    password: str


class AuthSuccessResponse(BaseModel):
    """Auth success response schema."""

    token: str
    userId: int
    username: str


class AuthErrorResponse(BaseModel):
    """Auth error response schema."""

    message: str

