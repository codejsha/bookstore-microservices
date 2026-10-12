from fastapi import HTTPException
from pydantic import ValidationError

from internal.domain.error import ConflictError, UnknownReferenceError


def invalid_command_error(exc: ValidationError) -> HTTPException:
    errors = exc.errors(include_url=False, include_input=False)
    detail = "; ".join(
        ".".join(str(part) for part in err.get("loc", ())) + ": " + str(err.get("msg", "invalid value"))
        for err in errors
    )
    return HTTPException(status_code=400, detail=detail or "invalid request")


def conflict_error(exc: ConflictError) -> HTTPException:
    return HTTPException(status_code=409, detail=str(exc) or "resource already exists")


def unknown_reference_error(exc: UnknownReferenceError) -> HTTPException:
    return HTTPException(status_code=400, detail=str(exc) or "referenced resource not found")
