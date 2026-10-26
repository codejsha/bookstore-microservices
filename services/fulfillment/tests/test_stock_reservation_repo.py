import pytest
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import ReleaseParams, ReserveParams
from internal.domain.constant.reservation_status import ReservationStatus
from internal.domain.model.error import InsufficientStockError
from internal.infrastructure.adapter.mysql.models import StockEntity, StockHistoryEntity, StockReservationEntity
from internal.infrastructure.adapter.mysql.stock_reservation_repo import MySQLStockReservationRepository
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid
from tests.conftest import seed_stock, seed_warehouse


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLStockReservationRepository:
    return MySQLStockReservationRepository(session_factory)


def _reserve(key: str, edition_id: int, quantity: int) -> ReserveParams:
    return ReserveParams(
        reservation_key=key, order_ref="order-1", edition_id=edition_id, quantity=quantity, reason="test"
    )


async def _quantities(session_factory, edition_id: int) -> dict[int, int]:
    async with session_factory() as session:
        rows = (await session.execute(select(StockEntity).where(StockEntity.edition_id == edition_id))).scalars()
        return {r.warehouse_id: r.quantity for r in rows}


async def _reservation_rows(session_factory, key: str) -> list[StockReservationEntity]:
    async with session_factory() as session:
        return list(
            (
                await session.execute(
                    select(StockReservationEntity)
                    .where(StockReservationEntity.reservation_key == key)
                    .order_by(StockReservationEntity.id)
                )
            ).scalars()
        )


async def _history_types(session_factory) -> list[tuple[str, int]]:
    async with session_factory() as session:
        rows = (await session.execute(select(StockHistoryEntity).order_by(StockHistoryEntity.id))).scalars()
        return [(r.change_type, r.change_qty) for r in rows]


async def test_reserve_single_warehouse_deducts_and_records_reservation(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)

    result = await repo.reserve(_reserve("o:1", 1, 4))

    assert result.quantity == 4
    assert result.status == ReservationStatus.RESERVED
    assert result.already_applied is False
    assert result.edition_uid == bytes_to_uuid(stock.edition_uid)
    assert await _quantities(session_factory, 1) == {warehouse.id: 6}
    assert await _history_types(session_factory) == [("RESERVATION", -4)]


async def test_reserve_spanning_warehouses_takes_largest_first(repo, session_factory) -> None:
    small = await seed_warehouse(session_factory, name="small")
    large = await seed_warehouse(session_factory, name="large")
    await seed_stock(session_factory, small, edition_id=1, quantity=3)
    await seed_stock(session_factory, large, edition_id=1, quantity=5)

    result = await repo.reserve(_reserve("o:1", 1, 7))

    assert result.quantity == 7
    assert await _quantities(session_factory, 1) == {small.id: 1, large.id: 0}
    rows = await _reservation_rows(session_factory, "o:1")
    assert {r.warehouse_id: r.quantity for r in rows} == {large.id: 5, small.id: 2}


async def test_reserve_insufficient_total_raises_and_leaves_quantities(repo, session_factory) -> None:
    first = await seed_warehouse(session_factory, name="A")
    second = await seed_warehouse(session_factory, name="B")
    await seed_stock(session_factory, first, edition_id=1, quantity=3)
    await seed_stock(session_factory, second, edition_id=1, quantity=2)

    with pytest.raises(InsufficientStockError):
        await repo.reserve(_reserve("o:1", 1, 6))

    assert await _quantities(session_factory, 1) == {first.id: 3, second.id: 2}
    assert await _reservation_rows(session_factory, "o:1") == []
    assert await _history_types(session_factory) == []


async def test_reserve_unknown_edition_raises_insufficient_stock(repo) -> None:
    with pytest.raises(InsufficientStockError):
        await repo.reserve(_reserve("o:9", 9, 1))


async def test_reserve_same_key_twice_replays_without_double_deduction(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)

    first = await repo.reserve(_reserve("o:1", 1, 4))
    second = await repo.reserve(_reserve("o:1", 1, 4))

    assert second.uid == first.uid
    assert second.already_applied is True
    assert second.quantity == 4
    assert await _quantities(session_factory, 1) == {warehouse.id: 6}
    assert len(await _reservation_rows(session_factory, "o:1")) == 1


async def test_release_restores_quantities_and_marks_released(repo, session_factory) -> None:
    small = await seed_warehouse(session_factory, name="small")
    large = await seed_warehouse(session_factory, name="large")
    await seed_stock(session_factory, small, edition_id=1, quantity=3)
    await seed_stock(session_factory, large, edition_id=1, quantity=5)
    await repo.reserve(_reserve("o:1", 1, 7))

    result = await repo.release(ReleaseParams(reservation_key="o:1", reason="compensate"))

    assert result is not None
    assert result.status == ReservationStatus.RELEASED
    assert result.already_applied is False
    assert result.quantity == 7
    assert await _quantities(session_factory, 1) == {small.id: 3, large.id: 5}
    rows = await _reservation_rows(session_factory, "o:1")
    assert all(r.status == "RELEASED" and r.released_at is not None and r.version == 2 for r in rows)
    history = await _history_types(session_factory)
    assert sorted(history) == sorted([("RESERVATION", -5), ("RESERVATION", -2), ("RELEASE", 5), ("RELEASE", 2)])


async def test_release_twice_replays_without_restoring_again(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    await repo.reserve(_reserve("o:1", 1, 4))
    await repo.release(ReleaseParams(reservation_key="o:1", reason=None))

    result = await repo.release(ReleaseParams(reservation_key="o:1", reason=None))

    assert result is not None
    assert result.already_applied is True
    assert result.status == ReservationStatus.RELEASED
    assert await _quantities(session_factory, 1) == {warehouse.id: 10}


async def test_release_unknown_key_is_none(repo) -> None:
    assert await repo.release(ReleaseParams(reservation_key="missing", reason=None)) is None
