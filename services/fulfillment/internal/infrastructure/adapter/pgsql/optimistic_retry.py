import asyncio
import random
from collections.abc import Awaitable, Callable

from sqlalchemy.exc import DBAPIError, IntegrityError

MAX_ATTEMPTS = 3
BASE_DELAY_MS = 5
CAP_DELAY_MS = 80

SQLSTATE_UNIQUE_VIOLATION = "23505"
SQLSTATE_DEADLOCK_DETECTED = "40P01"
SQLSTATE_LOCK_NOT_AVAILABLE = "55P03"


class VersionConflictError(Exception):
    pass


class StockCreateConflictError(Exception):
    pass


class ReservationConflictError(Exception):
    pass


class OptimisticRetryExhaustedError(Exception):
    pass


def sqlstate(exc: BaseException) -> str | None:
    orig = getattr(exc, "orig", None)
    for candidate in (orig, getattr(orig, "__cause__", None)):
        state = getattr(candidate, "sqlstate", None) or getattr(candidate, "pgcode", None)
        if isinstance(state, str) and state:
            return state
    return None


def is_unique_violation(exc: BaseException) -> bool:
    if not isinstance(exc, IntegrityError):
        return False
    if sqlstate(exc) == SQLSTATE_UNIQUE_VIOLATION:
        return True
    return "UNIQUE constraint failed" in str(exc)


def is_retryable(exc: BaseException) -> bool:
    if isinstance(exc, VersionConflictError | StockCreateConflictError):
        return True
    return isinstance(exc, DBAPIError) and sqlstate(exc) in (
        SQLSTATE_DEADLOCK_DETECTED,
        SQLSTATE_LOCK_NOT_AVAILABLE,
    )


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
