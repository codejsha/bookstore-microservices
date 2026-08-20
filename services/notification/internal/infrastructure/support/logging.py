import logging
import sys
import time
from collections.abc import Awaitable, Callable
from typing import Any

import structlog
from fastapi import Request, Response
from opentelemetry import trace

from internal.infrastructure.support.auth import HEADER_USER_ID


def _add_otel_context(_logger: Any, _method: str, event_dict: dict) -> dict:
    span = trace.get_current_span()
    if span is None:
        return event_dict
    ctx = span.get_span_context()
    if ctx.is_valid:
        event_dict["trace_id"] = format(ctx.trace_id, "032x")
        event_dict["span_id"] = format(ctx.span_id, "016x")
    return event_dict


class _DropUvicornAccessFilter(logging.Filter):
    def filter(self, record: logging.LogRecord) -> bool:
        return False


def _log_access(request: Request, status: int, elapsed_ns: int) -> None:
    fields: dict[str, Any] = {
        "method": request.method,
        "path": request.url.path,
        "status": status,
        "latency_us": elapsed_ns // 1000,
        "client_ip": request.client.host if request.client else "",
    }
    forwarded = request.headers.get("x-forwarded-for")
    if forwarded:
        fields["client_ip"] = forwarded.split(",")[0].strip()
    if request.url.query:
        fields["query"] = request.url.query
    user_uid = request.headers.get(HEADER_USER_ID)
    if user_uid:
        fields["user_uid"] = user_uid

    logger = structlog.get_logger()
    if status >= 500:
        logger.error("", **fields)
    elif status >= 400:
        logger.warning("", **fields)
    else:
        logger.info("", **fields)


async def access_log_middleware(
    request: Request,
    call_next: Callable[[Request], Awaitable[Response]],
) -> Response:
    start = time.perf_counter_ns()
    user_uid = request.headers.get(HEADER_USER_ID)
    if user_uid:
        trace.get_current_span().set_attribute("enduser.id", user_uid)
    try:
        response = await call_next(request)
    except Exception:
        _log_access(request, 500, time.perf_counter_ns() - start)
        raise

    status = response.status_code
    if request.url.path == "/health" and status < 400:
        return response

    _log_access(request, status, time.perf_counter_ns() - start)
    return response


def configure_logging(level: str = "info", debug: bool = False) -> None:
    log_level = getattr(logging, level.upper(), logging.INFO)

    logging.basicConfig(
        format="%(message)s",
        stream=sys.stdout,
        level=log_level,
    )

    logging.getLogger("uvicorn.access").addFilter(_DropUvicornAccessFilter())

    renderer: Any = structlog.dev.ConsoleRenderer() if debug else structlog.processors.JSONRenderer()

    structlog.configure(
        processors=[
            structlog.contextvars.merge_contextvars,
            structlog.processors.add_log_level,
            structlog.processors.TimeStamper(fmt="iso"),
            _add_otel_context,
            renderer,
        ],
        wrapper_class=structlog.make_filtering_bound_logger(log_level),
        logger_factory=structlog.stdlib.LoggerFactory(),
        cache_logger_on_first_use=True,
    )
