"""Auth API routes."""

from fastapi import APIRouter, Depends, HTTPException, Request, status
from sqlalchemy.orm import Session

from app.database import Database
from app.modules.auth.api.schemas import (
    AuthErrorResponse,
    AuthSuccessResponse,
    LoginRequest,
    RegisterRequest,
)
from app.modules.auth.app.errors import (
    ConflictError,
    UnauthorizedError,
    ValidationError,
)
from app.modules.auth.app.service import AuthService

router = APIRouter(prefix="/auth", tags=["auth"])


def get_db(request: Request):
    """Get database session from app state."""
    database: Database = request.app.state.database
    db = database.get_session()
    try:
        yield db
    finally:
        db.close()


def get_auth_service(db: Session = Depends(get_db)) -> AuthService:
    """Dependency to get auth service."""
    from app.modules.auth.infra.persistence.user_repository import SQLAlchemyUserRepository
    from app.modules.auth.infra.security.password_hasher import HMACPasswordHasher
    from app.modules.auth.infra.token.jwt_generator import JWTTokenGenerator
    from app.main import get_settings

    settings = get_settings()
    user_repo = SQLAlchemyUserRepository(db)
    password_hasher = HMACPasswordHasher()
    token_generator = JWTTokenGenerator(
        key=settings.jwt_key,
        issuer=settings.jwt_issuer,
        audience=settings.jwt_audience,
        access_token_lifetime_hours=settings.jwt_access_token_lifetime_hours,
    )
    return AuthService(user_repo, password_hasher, token_generator)


@router.post(
    "/register",
    response_model=AuthSuccessResponse,
    status_code=status.HTTP_200_OK,
    responses={
        400: {"model": AuthErrorResponse},
        409: {"model": AuthErrorResponse},
    },
)
def register(
    request: RegisterRequest,
    service: AuthService = Depends(get_auth_service),
):
    """Register a new user account."""
    try:
        result = service.register(request.username, request.email, request.password)
        return AuthSuccessResponse(**result)
    except ValidationError as e:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=e.message)
    except ConflictError as e:
        raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail=e.message)
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to process request.",
        )


@router.post(
    "/login",
    response_model=AuthSuccessResponse,
    status_code=status.HTTP_200_OK,
    responses={
        400: {"model": AuthErrorResponse},
        401: {"model": AuthErrorResponse},
    },
)
def login(
    request: LoginRequest,
    service: AuthService = Depends(get_auth_service),
):
    """Authenticate a user with username and password."""
    try:
        result = service.login(request.username, request.password)
        return AuthSuccessResponse(**result)
    except ValidationError as e:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=e.message)
    except UnauthorizedError as e:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail=e.message)
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to process request.",
        )

