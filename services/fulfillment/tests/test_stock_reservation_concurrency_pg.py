from __future__ import annotations

import asyncio
import os
from collections.abc import AsyncIterator
from datetime import UTC, datetime
from uuid import uuid4, uuid7

import pytest
import pytest_asyncio
from sqlalchemy import case, func, select, update
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

import internal.infrastructure.adapter.pgsql.optimistic_retry as optimistic_retry
from internal.application.port.repo.inventory_repos import ReleaseParams, ReserveParams
from internal.infrastructure.adapter.pgsql.models import Base, StockEntity, StockReservationEntity, WarehouseEntity
from internal.infrastructure.adapter.pgsql.stock_reservation_repo import PgsqlStockReservationRepository
from tests.postgres_container import container_runtime_available, postgres_container_url

os.environ.setdefault("TESTCONTAINERS_RYUK_DISABLED", "true")

pytestmark = pytest.mark.skipif(
    not container_runtime_available(),
    reason="no reachable container runtime (Docker/Podman) for testcontainers",
)

EDITION_ID = 1
WAREHOUSE_QUANTITIES = (30, 30)
CONCURRENT_RESERVATIONS = 20
RESERVE_QUANTITY = 3


@pytest.fixture(scope="module")
def postgres_url() -> str:
    yield from postgres_container_url()


async def _seed(factory: async_sessionmaker[AsyncSession]) -> None:
    now = datetime.now(UTC)
    async with factory() as session:
        warehouses = [
            WarehouseEntity(uid=uuid7(), name=f"W{i}", address=None, capacity=1000, created_at=now, version=1)
            for i in range(len(WAREHOUSE_QUANTITIES))
        ]
        session.add_all(warehouses)
        await session.flush()
        edition_uid = uuid4()
        session.add_all(
            [
                StockEntity(
                    uid=uuid7(),
                    edition_id=EDITION_ID,
                    edition_uid=edition_uid,
                    warehouse_id=w.id,
                    warehouse_uid=w.uid,
                    quantity=q,
                    created_at=now,
                    version=1,
                )
                for w, q in zip(warehouses, WAREHOUSE_QUANTITIES, strict=True)
            ]
        )
        await session.commit()


@pytest_asyncio.fixture
async def fixture_env(postgres_url: str) -> AsyncIterator[tuple[PgsqlStockReservationRepository, async_sessionmaker]]:
    engine = create_async_engine(postgres_url, pool_size=CONCURRENT_RESERVATIONS + 2, max_overflow=5)
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.drop_all)
        await conn.run_sync(Base.metadata.create_all)
    factory = async_sessionmaker(bind=engine, expire_on_commit=False)
    await _seed(factory)
    try:
        yield PgsqlStockReservationRepository(factory), factory
    finally:
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.drop_all)
        await engine.dispose()


async def _stock_quantities(factory: async_sessionmaker) -> list[int]:
    async with factory() as session:
        return list((await session.execute(select(StockEntity.quantity).order_by(StockEntity.id))).scalars())


async def _reserved_total(factory: async_sessionmaker) -> int:
    async with factory() as session:
        return (
            await session.execute(
                select(func.coalesce(func.sum(StockReservationEntity.quantity), 0)).where(
                    StockReservationEntity.status == "RESERVED"
                )
            )
        ).scalar()


def _params(key: str, quantity: int = RESERVE_QUANTITY) -> ReserveParams:
    return ReserveParams(reservation_key=key, order_ref=key, edition_id=EDITION_ID, quantity=quantity, reason="load")


async def test_reserve_concurrent_distinct_keys_never_goes_negative_and_totals_match(fixture_env) -> None:
    repo, factory = fixture_env

    results = await asyncio.gather(*(repo.reserve(_params(f"order-{i}:1")) for i in range(CONCURRENT_RESERVATIONS)))

    assert all(r.quantity == RESERVE_QUANTITY and not r.already_applied for r in results)
    quantities = await _stock_quantities(factory)
    assert all(q >= 0 for q in quantities)
    assert sum(quantities) == sum(WAREHOUSE_QUANTITIES) - CONCURRENT_RESERVATIONS * RESERVE_QUANTITY
    assert await _reserved_total(factory) == CONCURRENT_RESERVATIONS * RESERVE_QUANTITY


async def test_reserve_concurrent_same_key_applies_once(fixture_env) -> None:
    repo, factory = fixture_env

    results = await asyncio.gather(*(repo.reserve(_params("order-same:1", 7)) for _ in range(10)))

    assert sum(1 for r in results if not r.already_applied) == 1
    assert len({r.uid for r in results}) == 1
    assert sum(await _stock_quantities(factory)) == sum(WAREHOUSE_QUANTITIES) - 7
    assert await _reserved_total(factory) == 7


async def test_release_after_concurrent_reserves_restores_every_unit(fixture_env) -> None:
    repo, factory = fixture_env
    keys = [f"order-{i}:1" for i in range(CONCURRENT_RESERVATIONS)]
    await asyncio.gather(*(repo.reserve(_params(key)) for key in keys))

    await asyncio.gather(*(repo.release(ReleaseParams(reservation_key=key, reason="undo")) for key in keys))

    assert sum(await _stock_quantities(factory)) == sum(WAREHOUSE_QUANTITIES)
    assert await _reserved_total(factory) == 0


async def _set_stock(factory: async_sessionmaker, quantity) -> None:
    async with factory() as session, session.begin():
        await session.execute(update(StockEntity).values(quantity=quantity))


async def test_release_split_reservations_racing_reserves_no_deadlock(fixture_env, monkeypatch) -> None:
    repo, factory = fixture_env
    first_stock_id = min(await _stock_ids(factory))
    split_keys = [f"split-{i}:1" for i in range(8)]
    for key in split_keys:
        await _set_stock(factory, case((StockEntity.id == first_stock_id, 1), else_=2))
        await repo.reserve(_params(key, 3))
    await _set_stock(factory, 1000)

    deadlocks: list[BaseException] = []
    original = optimistic_retry.is_retryable

    def spy(exc: BaseException) -> bool:
        if optimistic_retry.sqlstate(exc) == optimistic_retry.SQLSTATE_DEADLOCK_DETECTED:
            deadlocks.append(exc)
        return original(exc)

    monkeypatch.setattr(optimistic_retry, "is_retryable", spy)

    small_keys = [f"small-{i}:1" for i in range(12)]
    results = await asyncio.gather(
        *(repo.release(ReleaseParams(reservation_key=key, reason="undo")) for key in split_keys),
        *(repo.reserve(_params(key, 1)) for key in small_keys),
        return_exceptions=True,
    )

    assert [r for r in results if isinstance(r, BaseException)] == []
    assert deadlocks == []
    assert sum(await _stock_quantities(factory)) == 2000 + 8 * 3 - len(small_keys)
    assert await _reserved_total(factory) == len(small_keys)


async def _stock_ids(factory: async_sessionmaker) -> list[int]:
    async with factory() as session:
        return list((await session.execute(select(StockEntity.id))).scalars())
