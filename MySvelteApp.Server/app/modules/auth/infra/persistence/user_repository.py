"""SQLAlchemy-based user repository implementation."""

from typing import Optional

from sqlalchemy.orm import Session

from app.modules.auth.app.errors import ConflictError
from app.modules.auth.app.ports import UserRepository
from app.modules.auth.domain.user import User
from app.modules.auth.infra.persistence.models import UserModel


class SQLAlchemyUserRepository(UserRepository):
    """User repository using SQLAlchemy."""

    def __init__(self, db: Session):
        """Initialize repository with database session."""
        self.db = db

    def add(self, user: User) -> None:
        """Persist a new user."""
        if user is None:
            raise ValueError("user cannot be nil")

        user_model = UserModel(
            username=user.username,
            email=user.email,
            password_hash=user.password_hash,
            password_salt=user.password_salt,
        )

        try:
            self.db.add(user_model)
            self.db.commit()
            self.db.refresh(user_model)
            user.id = user_model.id
        except Exception as e:
            self.db.rollback()
            # Check for unique constraint violations
            if "UNIQUE constraint failed" in str(e) or "unique constraint" in str(e).lower():
                # Detect which field caused the conflict
                if self.username_exists(user.username):
                    raise ConflictError("This username is already taken. Please choose a different one.")
                if self.email_exists(user.email):
                    raise ConflictError("This email is already registered. Please use a different email address.")
                raise ConflictError("A user with the provided credentials already exists.")
            raise

    def get_by_username(self, username: str) -> Optional[User]:
        """Retrieve a user by username."""
        username = username.strip()
        if not username:
            raise ValueError("username cannot be blank")

        user_model = self.db.query(UserModel).filter(UserModel.username == username).first()
        if user_model is None:
            return None

        return User(
            id=user_model.id,
            username=user_model.username,
            email=user_model.email,
            password_hash=user_model.password_hash,
            password_salt=user_model.password_salt,
            created_at=user_model.created_at,
            updated_at=user_model.updated_at,
        )

    def username_exists(self, username: str) -> bool:
        """Check if a username is already taken."""
        username = username.strip()
        if not username:
            raise ValueError("username cannot be blank")

        count = self.db.query(UserModel).filter(UserModel.username == username).count()
        return count > 0

    def email_exists(self, email: str) -> bool:
        """Check if an email address is already registered."""
        normalized = email.strip().lower()
        if not normalized:
            raise ValueError("email cannot be blank")

        count = self.db.query(UserModel).filter(UserModel.email == normalized).count()
        return count > 0

