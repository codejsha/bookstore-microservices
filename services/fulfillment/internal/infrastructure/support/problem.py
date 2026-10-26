from http import HTTPStatus

from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
from starlette.exceptions import HTTPException as StarletteHTTPException

PROBLEM_MEDIA_TYPE = "application/problem+json"


def problem(status: int, detail: str | None = None, errors: list[str] | None = None) -> dict:
    body: dict = {"title": HTTPStatus(status).phrase, "status": status}
    if detail is not None:
        body["detail"] = detail
    if errors is not None:
        body["errors"] = errors
    return body


def register_problem_handlers(app: FastAPI) -> None:
    @app.exception_handler(StarletteHTTPException)
    async def handle_http_exception(_request: Request, exc: StarletteHTTPException) -> JSONResponse:
        return JSONResponse(
            status_code=exc.status_code,
            content=problem(exc.status_code, str(exc.detail) if exc.detail is not None else None),
            headers=exc.headers,
            media_type=PROBLEM_MEDIA_TYPE,
        )

    @app.exception_handler(RequestValidationError)
    async def handle_validation_error(_request: Request, exc: RequestValidationError) -> JSONResponse:
        errors = sorted(
            f"{'.'.join(str(part) for part in error['loc'] if part != 'body')}: {error['msg']}"
            for error in exc.errors()
        )
        return JSONResponse(
            status_code=422,
            content=problem(422, "request validation failed", errors),
            media_type=PROBLEM_MEDIA_TYPE,
        )
