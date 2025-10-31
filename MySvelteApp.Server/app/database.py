"""Database setup and session management."""

from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

from app.config import Settings

Base = declarative_base()


class Database:
    """Database connection manager."""

    def __init__(self, settings: Settings):
        """Initialize database with settings."""
        self.settings = settings
        self.engine = create_engine(
            settings.database_dsn,
            connect_args={"check_same_thread": False} if "sqlite" in settings.database_dsn else {},
            echo=False,
        )
        self.SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=self.engine)

    def create_tables(self):
        """Create all database tables."""
        Base.metadata.create_all(bind=self.engine)

    def get_session(self):
        """Get a database session."""
        return self.SessionLocal()


def get_db(database: Database):
    """Dependency for FastAPI to get database session."""
    db = database.get_session()
    try:
        yield db
    finally:
        db.close()

