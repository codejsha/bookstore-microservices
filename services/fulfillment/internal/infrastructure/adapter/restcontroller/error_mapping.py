from fastapi import HTTPException
from pydantic import ValidationError

from internal.domain.model.error import (
    ConflictError,
    DuplicateClosingError,
    InsufficientStockError,
    InvalidQuantityError,
    InvalidTransitionError,
    SameWarehouseError,
    StockNotFoundError,
    WarehouseNotFoundError,
)


def invalid_command_error(exc: ValidationError) -> HTTPException:
    errors = exc.errors(include_url=False, include_input=False)
    detail = "; ".join(
        ".".join(str(part) for part in err.get("loc", ())) + ": " + str(err.get("msg", "invalid value"))
        for err in errors
    )
    return HTTPException(status_code=400, detail=detail or "invalid request")


def conflict_error(exc: ConflictError) -> HTTPException:
    return HTTPException(status_code=409, detail=str(exc) or "resource already exists")


def duplicate_key_error(detail: str) -> HTTPException:
    return HTTPException(status_code=409, detail=detail)


_INVENTORY_STATUS: tuple[tuple[type[Exception], int], ...] = (
    (InvalidQuantityError, 400),
    (InsufficientStockError, 400),
    (SameWarehouseError, 400),
    (StockNotFoundError, 404),
    (WarehouseNotFoundError, 404),
    (DuplicateClosingError, 409),
    (InvalidTransitionError, 409),
)


def inventory_error(exc: Exception) -> HTTPException:
    if isinstance(exc, ValidationError):
        return invalid_command_error(exc)
    for error_type, status in _INVENTORY_STATUS:
        if isinstance(exc, error_type):
            return HTTPException(status_code=status, detail=str(exc) or error_type.__name__)
    return HTTPException(status_code=500, detail="internal server error")
