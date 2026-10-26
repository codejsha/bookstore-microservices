from __future__ import annotations

import asyncio
import os
from collections.abc import AsyncIterator, Iterator
from datetime import UTC, datetime
from uuid import uuid7

import pytest
import pytest_asyncio
from sqlalchemy import select
from sqlalchemy.exc import DBAPIError, IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from internal.infrastructure.adapter.pgsql.models import Base, WarehouseEntity
from internal.infrastructure.adapter.pgsql.optimistic_retry import (
    SQLSTATE_DEADLOCK_DETECTED,
    SQLSTATE_LOCK_NOT_AVAILABLE,
    SQLSTATE_UNIQUE_VIOLATION,
    is_retryable,
    is_unique_violation,
    sqlstate,
)
from tests.postgres_container import container_runtime_available, postgres_container_url

os.environ.setdefault("TESTCONTAINERS_RYUK_DISABLED", "true")

pytestmark = pytest.mark.skipif(
    not container_runtime_available(),
    reason="no reachable container runtime (Docker/Podman) for testcontainers",
)


@pytest.fixture(scope="module")
def postgres_url() -> Iterator[str]:
    yield from postgres_container_url()


def _warehouse(name: str) -> WarehouseEntity:
    return WarehouseEntity(
        uid=uuid7(),
        name=name,
        address=None,
        capacity=10,
        created_at=datetime.now(UTC),
        updated_at=None,
        deleted_at=None,
        version=1,
    )


@pytest_asyncio.fixture
async def factory(postgres_url: str) -> AsyncIterator[async_sessionmaker[AsyncSession]]:
    engine = create_async_engine(postgres_url, pool_size=4, max_overflow=2)
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    session_factory = async_sessionmaker(bind=engine, expire_on_commit=False)
    async with session_factory() as session:
        session.add_all([_warehouse("A"), _warehouse("B")])
        await session.commit()
    try:
        yield session_factory
    finally:
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.drop_all)
        await engine.dispose()


async def _lock(session: AsyncSession, name: str, nowait: bool = False) -> None:
    await session.execute(select(WarehouseEntity).where(WarehouseEntity.name == name).with_for_update(nowait=nowait))


async def test_sqlstate_duplicate_uid_insert_unique_violation(factory: async_sessionmaker[AsyncSession]) -> None:
    async with factory() as session:
        existing = (await session.execute(select(WarehouseEntity).where(WarehouseEntity.name == "A"))).scalar_one()
        duplicate = _warehouse("A-dup")
        duplicate.uid = existing.uid
        session.add(duplicate)
        with pytest.raises(IntegrityError) as raised:
            await session.flush()

    assert sqlstate(raised.value) == SQLSTATE_UNIQUE_VIOLATION
    assert is_unique_violation(raised.value)
    assert not is_retryable(raised.value)


async def test_sqlstate_for_update_nowait_on_locked_row_lock_not_available(
    factory: async_sessionmaker[AsyncSession],
) -> None:
    async with factory() as holder, holder.begin(), factory() as contender, contender.begin():
        await _lock(holder, "A")
        with pytest.raises(DBAPIError) as raised:
            await _lock(contender, "A", nowait=True)

    assert sqlstate(raised.value) == SQLSTATE_LOCK_NOT_AVAILABLE
    assert is_retryable(raised.value)
    assert not is_unique_violation(raised.value)


async def test_sqlstate_crossed_row_locks_deadlock_detected(factory: async_sessionmaker[AsyncSession]) -> None:
    first_locked = asyncio.Event()
    second_locked = asyncio.Event()

    async def worker(own: str, other: str, mine: asyncio.Event, theirs: asyncio.Event) -> BaseException | None:
        async with factory() as session, session.begin():
            await _lock(session, own)
            mine.set()
            await theirs.wait()
            try:
                await _lock(session, other)
            except DBAPIError as exc:
                return exc
        return None

    outcomes = await asyncio.gather(
        worker("A", "B", first_locked, second_locked),
        worker("B", "A", second_locked, first_locked),
    )
    failures = [exc for exc in outcomes if exc is not None]

    assert len(failures) == 1
    assert sqlstate(failures[0]) == SQLSTATE_DEADLOCK_DETECTED
    assert is_retryable(failures[0])
