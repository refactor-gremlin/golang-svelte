"""Structured logging setup."""

import logging
import sys
from typing import Optional

# Configure root logger
logger = logging.getLogger("mysvelteapp")
logger.setLevel(logging.INFO)

# Console handler with structured format
handler = logging.StreamHandler(sys.stdout)
handler.setLevel(logging.INFO)

# Formatter
formatter = logging.Formatter(
    "%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
handler.setFormatter(formatter)
logger.addHandler(handler)


def get_logger(name: Optional[str] = None) -> logging.Logger:
    """Get a logger instance."""
    if name:
        return logging.getLogger(f"mysvelteapp.{name}")
    return logger


def configure_logging(level: str = "INFO", format_type: str = "text"):
    """Configure logging level and format."""
    log_level = getattr(logging, level.upper(), logging.INFO)
    logger.setLevel(log_level)
    handler.setLevel(log_level)

    if format_type.lower() == "json":
        import json
        import sys
        from datetime import datetime

        class JSONFormatter(logging.Formatter):
            def format(self, record):
                log_data = {
                    "timestamp": datetime.utcnow().isoformat(),
                    "level": record.levelname,
                    "logger": record.name,
                    "message": record.getMessage(),
                }
                if record.exc_info:
                    log_data["exception"] = self.formatException(record.exc_info)
                return json.dumps(log_data)

        handler.setFormatter(JSONFormatter())
    else:
        handler.setFormatter(formatter)

