import asyncio
import random
from collections.abc import Awaitable, Callable

from sqlalchemy.exc import DBAPIError, IntegrityError

MAX_ATTEMPTS = 3
BASE_DELAY_MS = 5
CAP_DELAY_MS = 80

MYSQL_DUPLICATE_ENTRY = 1062
MYSQL_LOCK_WAIT_TIMEOUT = 1205
MYSQL_DEADLOCK = 1213


class VersionConflictError(Exception):
    pass


class StockCreateConflictError(Exception):
    pass


class ReservationConflictError(Exception):
    pass


class OptimisticRetryExhaustedError(Exception):
    pass


def mysql_error_code(exc: BaseException) -> int | None:
    orig = getattr(exc, "orig", None)
    args = getattr(orig, "args", ())
    if args and isinstance(args[0], int):
        return args[0]
    return None


def is_unique_violation(exc: BaseException) -> bool:
    if not isinstance(exc, IntegrityError):
        return False
    if mysql_error_code(exc) == MYSQL_DUPLICATE_ENTRY:
        return True
    return "UNIQUE constraint failed" in str(exc)


def is_retryable(exc: BaseException) -> bool:
    if isinstance(exc, VersionConflictError | StockCreateConflictError):
        return True
    return isinstance(exc, DBAPIError) and mysql_error_code(exc) in (MYSQL_LOCK_WAIT_TIMEOUT, MYSQL_DEADLOCK)


def retry_delay_seconds(attempt: int, base_ms: int = BASE_DELAY_MS, cap_ms: int = CAP_DELAY_MS) -> float:
    delay = base_ms
    for _ in range(max(attempt, 0)):
        if delay >= cap_ms:
            break
        delay *= 2
    delay = min(delay, cap_ms)
    half = delay / 2
    return (half + random.uniform(0, half)) / 1000


async def run_with_optimistic_retry[T](
    fn: Callable[[], Awaitable[T]],
    attempts: int = MAX_ATTEMPTS,
    base_ms: int = BASE_DELAY_MS,
    cap_ms: int = CAP_DELAY_MS,
) -> T:
    last: BaseException | None = None
    for attempt in range(attempts):
        try:
            return await fn()
        except Exception as exc:
            if not is_retryable(exc):
                raise
            last = exc
            if attempt < attempts - 1:
                await asyncio.sleep(retry_delay_seconds(attempt, base_ms, cap_ms))
    raise OptimisticRetryExhaustedError(f"exhausted {attempts} optimistic-lock retries") from last
