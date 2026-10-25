import pytest
from sqlalchemy.exc import IntegrityError, OperationalError

from internal.domain.model.error import InsufficientStockError
from internal.infrastructure.adapter.mysql import optimistic_retry as retry_module
from internal.infrastructure.adapter.mysql.optimistic_retry import (
    OptimisticRetryExhaustedError,
    ReservationConflictError,
    StockCreateConflictError,
    VersionConflictError,
    is_retryable,
    is_unique_violation,
    retry_delay_seconds,
    run_with_optimistic_retry,
)


@pytest.fixture(autouse=True)
def _no_sleep(monkeypatch: pytest.MonkeyPatch) -> None:
    async def _sleep(_seconds: float) -> None:
        return None

    monkeypatch.setattr(retry_module.asyncio, "sleep", _sleep)


async def test_run_with_optimistic_retry_version_conflict_then_success_returns_result() -> None:
    calls = {"n": 0}

    async def attempt() -> str:
        calls["n"] += 1
        if calls["n"] < 3:
            raise VersionConflictError()
        return "ok"

    assert await run_with_optimistic_retry(attempt) == "ok"
    assert calls["n"] == 3


async def test_run_with_optimistic_retry_stock_create_conflict_retries() -> None:
    calls = {"n": 0}

    async def attempt() -> str:
        calls["n"] += 1
        if calls["n"] == 1:
            raise StockCreateConflictError()
        return "ok"

    assert await run_with_optimistic_retry(attempt) == "ok"
    assert calls["n"] == 2


async def test_run_with_optimistic_retry_non_retryable_raises_immediately() -> None:
    calls = {"n": 0}

    async def attempt() -> None:
        calls["n"] += 1
        raise InsufficientStockError("nope")

    with pytest.raises(InsufficientStockError):
        await run_with_optimistic_retry(attempt)
    assert calls["n"] == 1


async def test_run_with_optimistic_retry_reservation_conflict_is_not_retried() -> None:
    calls = {"n": 0}

    async def attempt() -> None:
        calls["n"] += 1
        raise ReservationConflictError("key")

    with pytest.raises(ReservationConflictError):
        await run_with_optimistic_retry(attempt)
    assert calls["n"] == 1


async def test_run_with_optimistic_retry_every_attempt_conflicts_raises_exhausted() -> None:
    calls = {"n": 0}

    async def attempt() -> None:
        calls["n"] += 1
        raise VersionConflictError()

    with pytest.raises(OptimisticRetryExhaustedError):
        await run_with_optimistic_retry(attempt, attempts=3)
    assert calls["n"] == 3


def _operational_error(code: int) -> OperationalError:
    return OperationalError("UPDATE stock", {}, Exception(code, "mysql error"))


def test_is_retryable_mysql_deadlock_true() -> None:
    assert is_retryable(_operational_error(1213))


def test_is_retryable_mysql_lock_wait_timeout_true() -> None:
    assert is_retryable(_operational_error(1205))


def test_is_retryable_other_mysql_error_false() -> None:
    assert not is_retryable(_operational_error(1045))


def test_is_unique_violation_mysql_duplicate_entry_true() -> None:
    assert is_unique_violation(IntegrityError("INSERT", {}, Exception(1062, "Duplicate entry")))


def test_is_unique_violation_sqlite_message_true() -> None:
    assert is_unique_violation(IntegrityError("INSERT", {}, Exception("UNIQUE constraint failed: stock.uid")))


def test_is_unique_violation_other_integrity_error_false() -> None:
    assert not is_unique_violation(IntegrityError("INSERT", {}, Exception(1048, "Column cannot be null")))


def test_retry_delay_seconds_stays_within_jitter_bounds() -> None:
    for attempt in range(6):
        delay_ms = retry_delay_seconds(attempt) * 1000
        assert 2.5 <= delay_ms <= 80.0
    assert 40.0 <= retry_delay_seconds(5) * 1000 <= 80.0
