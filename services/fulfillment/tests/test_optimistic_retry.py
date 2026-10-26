import asyncpg
import pytest
from sqlalchemy.dialects.postgresql.asyncpg import AsyncAdapt_asyncpg_dbapi
from sqlalchemy.exc import IntegrityError, OperationalError

from internal.domain.model.error import InsufficientStockError
from internal.infrastructure.adapter.pgsql import optimistic_retry as retry_module
from internal.infrastructure.adapter.pgsql.optimistic_retry import (
    OptimisticRetryExhaustedError,
    ReservationConflictError,
    StockCreateConflictError,
    VersionConflictError,
    is_retryable,
    is_unique_violation,
    retry_delay_seconds,
    run_with_optimistic_retry,
    sqlstate,
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


def _dbapi_error(sqlstate: str, kind: type[Exception] = AsyncAdapt_asyncpg_dbapi.Error) -> OperationalError:
    orig = kind("postgres error")
    orig.pgcode = orig.sqlstate = sqlstate
    return OperationalError("UPDATE stock", {}, orig)


def _integrity_error(cause: Exception) -> IntegrityError:
    orig = AsyncAdapt_asyncpg_dbapi.IntegrityError(str(cause))
    orig.pgcode = orig.sqlstate = getattr(cause, "sqlstate", None)
    orig.__cause__ = cause
    return IntegrityError("INSERT", {}, orig)


def test_sqlstate_translated_adapter_error_reads_sqlstate() -> None:
    assert sqlstate(_dbapi_error("40P01")) == "40P01"


def test_sqlstate_raw_asyncpg_cause_reads_class_sqlstate() -> None:
    orig = AsyncAdapt_asyncpg_dbapi.Error("deadlock")
    orig.__cause__ = asyncpg.exceptions.DeadlockDetectedError("deadlock")
    assert sqlstate(OperationalError("UPDATE stock", {}, orig)) == "40P01"


def test_sqlstate_non_dbapi_error_none() -> None:
    assert sqlstate(RuntimeError("boom")) is None


def test_is_retryable_deadlock_detected_true() -> None:
    assert is_retryable(_dbapi_error("40P01"))


def test_is_retryable_lock_not_available_true() -> None:
    assert is_retryable(_dbapi_error("55P03"))


def test_is_retryable_other_sqlstate_false() -> None:
    assert not is_retryable(_dbapi_error("28P01"))


def test_is_retryable_serialization_failure_false() -> None:
    assert not is_retryable(_dbapi_error("40001"))


def test_is_unique_violation_postgres_unique_violation_true() -> None:
    assert is_unique_violation(_integrity_error(asyncpg.exceptions.UniqueViolationError("duplicate key")))


def test_is_unique_violation_sqlite_message_true() -> None:
    assert is_unique_violation(IntegrityError("INSERT", {}, Exception("UNIQUE constraint failed: stock.uid")))


def test_is_unique_violation_not_null_violation_false() -> None:
    assert not is_unique_violation(_integrity_error(asyncpg.exceptions.NotNullViolationError("null value")))


def test_is_unique_violation_non_integrity_error_false() -> None:
    assert not is_unique_violation(_dbapi_error("23505"))


def test_retry_delay_seconds_stays_within_jitter_bounds() -> None:
    for attempt in range(6):
        delay_ms = retry_delay_seconds(attempt) * 1000
        assert 2.5 <= delay_ms <= 80.0
    assert 40.0 <= retry_delay_seconds(5) * 1000 <= 80.0
