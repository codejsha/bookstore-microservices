from datetime import UTC, datetime
from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import ApplyChangeParams, ClosingCreateParams
from internal.domain.constant.closing_status import ClosingStatus
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.error import DuplicateClosingError, WarehouseNotFoundError
from internal.domain.model.option.inventory_option import ClosingQueryOption
from internal.infrastructure.adapter.pgsql.closing_repo import PgsqlMonthlyClosingRepository, month_range
from internal.infrastructure.adapter.pgsql.stock_repo import PgsqlStockRepository
from tests.conftest import seed_stock, seed_warehouse


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> PgsqlMonthlyClosingRepository:
    return PgsqlMonthlyClosingRepository(session_factory)


@pytest.fixture
def stock_repo(session_factory: async_sessionmaker[AsyncSession]) -> PgsqlStockRepository:
    return PgsqlStockRepository(session_factory)


def _this_month() -> tuple[int, int]:
    now = datetime.now(UTC)
    return now.year, now.month


def _previous_month() -> tuple[int, int]:
    year, month = _this_month()
    return (year - 1, 12) if month == 1 else (year, month - 1)


def test_month_range_december_rolls_into_next_year() -> None:
    start, end = month_range(2026, 12)
    assert start == datetime(2026, 12, 1, tzinfo=UTC)
    assert end == datetime(2027, 1, 1, tzinfo=UTC)


def test_month_range_is_utc_aware() -> None:
    start, end = month_range(2026, 3)
    assert start.tzinfo is UTC and end.tzinfo is UTC
    assert end == datetime(2026, 4, 1, tzinfo=UTC)


async def test_create_computes_period_movements_and_opening_quantity(repo, stock_repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    edition_uid = stock.edition_uid
    warehouse_uid = warehouse.uid
    for delta, change_type in (
        (5, StockChangeType.INBOUND),
        (-3, StockChangeType.RESERVATION),
        (2, StockChangeType.ADJUSTMENT),
    ):
        await stock_repo.apply_change(
            ApplyChangeParams(
                edition_uid=edition_uid, warehouse_uid=warehouse_uid, delta=delta, change_type=change_type, reason=None
            )
        )
    year, month = _this_month()

    closing = await repo.create(ClosingCreateParams(warehouse_uid=warehouse_uid, year=year, month=month))

    assert closing.status == ClosingStatus.CLOSED
    assert closing.closed_at is not None
    assert len(closing.items) == 1
    item = closing.items[0]
    assert item.edition_uid == edition_uid
    assert (item.inbound_quantity, item.outbound_quantity, item.adjust_quantity) == (5, 3, 2)
    assert item.closing_quantity == 14
    assert item.opening_quantity == 10


async def test_create_other_period_has_no_movements(repo, stock_repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    await stock_repo.apply_change(
        ApplyChangeParams(
            edition_uid=stock.edition_uid,
            warehouse_uid=warehouse.uid,
            delta=5,
            change_type=StockChangeType.INBOUND,
            reason=None,
        )
    )
    year, month = _previous_month()

    closing = await repo.create(ClosingCreateParams(warehouse_uid=warehouse.uid, year=year, month=month))

    item = closing.items[0]
    assert (item.inbound_quantity, item.outbound_quantity, item.adjust_quantity) == (0, 0, 0)
    assert item.opening_quantity == 15
    assert item.closing_quantity == 15


async def test_create_duplicate_period_raises_duplicate_closing(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    params = ClosingCreateParams(warehouse_uid=warehouse.uid, year=2026, month=9)
    await repo.create(params)

    with pytest.raises(DuplicateClosingError):
        await repo.create(params)


async def test_create_missing_warehouse_raises_warehouse_not_found(repo) -> None:
    with pytest.raises(WarehouseNotFoundError):
        await repo.create(ClosingCreateParams(warehouse_uid=uuid4(), year=2026, month=9))


async def test_find_all_by_period_filters_closings(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    warehouse_uid = warehouse.uid
    await repo.create(ClosingCreateParams(warehouse_uid=warehouse_uid, year=2026, month=8))
    september = await repo.create(ClosingCreateParams(warehouse_uid=warehouse_uid, year=2026, month=9))

    items, total = await repo.find_all(ClosingQueryOption(year=2026, month=9))

    assert total == 1
    assert items[0].uid == september.uid
    found = await repo.find_by_uid(september.uid)
    assert found is not None and found.month == 9
