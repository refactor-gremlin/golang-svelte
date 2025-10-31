"""SQLAlchemy models for auth module."""

from datetime import datetime

from sqlalchemy import Column, DateTime, Integer, String
from sqlalchemy.sql import func

from app.database import Base


class UserModel(Base):
    """User database model."""

    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String(64), unique=True, nullable=False, index=True)
    email = Column(String(320), unique=True, nullable=False, index=True)
    password_hash = Column(String(512), nullable=False)
    password_salt = Column(String(256), nullable=False)
    created_at = Column(DateTime, nullable=False, default=func.now())
    updated_at = Column(DateTime, nullable=False, default=func.now(), onupdate=func.now())

