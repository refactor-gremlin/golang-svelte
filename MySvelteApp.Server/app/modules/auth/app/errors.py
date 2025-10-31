"""Application layer error types."""


class ValidationError(Exception):
    """Indicates the payload failed validation rules."""

    def __init__(self, message: str):
        """Initialize validation error."""
        self.message = message
        super().__init__(self.message)


class ConflictError(Exception):
    """Indicates the request conflicts with existing state (e.g. duplicate username)."""

    def __init__(self, message: str):
        """Initialize conflict error."""
        self.message = message
        super().__init__(self.message)


class UnauthorizedError(Exception):
    """Indicates credentials were invalid."""

    def __init__(self, message: str):
        """Initialize unauthorized error."""
        self.message = message
        super().__init__(self.message)


def is_validation_error(err: Exception) -> bool:
    """Check if error is a ValidationError."""
    return isinstance(err, ValidationError)


def is_conflict_error(err: Exception) -> bool:
    """Check if error is a ConflictError."""
    return isinstance(err, ConflictError)


def is_unauthorized_error(err: Exception) -> bool:
    """Check if error is an UnauthorizedError."""
    return isinstance(err, UnauthorizedError)

