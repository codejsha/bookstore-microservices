from datetime import UTC, datetime

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import ApplyChangeParams
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.option.inventory_option import BalanceQueryOption
from internal.infrastructure.adapter.pgsql.balance_repo import PgsqlStockBalanceRepository
from internal.infrastructure.adapter.pgsql.stock_repo import PgsqlStockRepository
from tests.conftest import seed_stock, seed_warehouse


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> PgsqlStockBalanceRepository:
    return PgsqlStockBalanceRepository(session_factory)


@pytest.fixture
def stock_repo(session_factory: async_sessionmaker[AsyncSession]) -> PgsqlStockRepository:
    return PgsqlStockRepository(session_factory)


async def _move(stock_repo, stock, warehouse, delta, change_type) -> None:
    await stock_repo.apply_change(
        ApplyChangeParams(
            edition_uid=stock.edition_uid,
            warehouse_uid=warehouse.uid,
            delta=delta,
            change_type=change_type,
            reason=None,
        )
    )


async def test_find_all_sums_history_buckets_per_stock_row(repo, stock_repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    await _move(stock_repo, stock, warehouse, 5, StockChangeType.INBOUND)
    await _move(stock_repo, stock, warehouse, -3, StockChangeType.RESERVATION)
    await _move(stock_repo, stock, warehouse, 1, StockChangeType.RELEASE)
    await _move(stock_repo, stock, warehouse, -2, StockChangeType.ADJUSTMENT)

    items, total = await repo.find_all(BalanceQueryOption())

    assert total == 1
    entry = items[0]
    assert entry.warehouse_name == "Main"
    assert (entry.inbound_quantity, entry.outbound_quantity, entry.adjust_quantity) == (6, 3, -2)
    assert entry.current_quantity == 11


async def test_find_all_without_history_reports_zero_movements(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    await seed_stock(session_factory, warehouse, edition_id=1, quantity=4)

    items, total = await repo.find_all(BalanceQueryOption())

    assert total == 1
    assert (items[0].inbound_quantity, items[0].outbound_quantity, items[0].adjust_quantity) == (0, 0, 0)
    assert items[0].current_quantity == 4


async def test_find_all_by_year_month_outside_period_excludes_movements(repo, stock_repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    await _move(stock_repo, stock, warehouse, 5, StockChangeType.INBOUND)

    now = datetime.now(UTC)
    other = f"{now.year - 1}-{now.month:02d}"
    items, _ = await repo.find_all(BalanceQueryOption(year_month=other))

    assert items[0].inbound_quantity == 0
    assert items[0].current_quantity == 15


async def test_find_all_by_year_month_inside_period_includes_movements(repo, stock_repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    await _move(stock_repo, stock, warehouse, 5, StockChangeType.INBOUND)

    now = datetime.now(UTC)
    items, _ = await repo.find_all(BalanceQueryOption(year_month=f"{now.year}-{now.month:02d}"))

    assert items[0].inbound_quantity == 5


async def test_find_all_by_edition_filters_rows_and_total(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    first = await seed_stock(session_factory, warehouse, edition_id=1, quantity=1)
    await seed_stock(session_factory, warehouse, edition_id=2, quantity=2)

    items, total = await repo.find_all(BalanceQueryOption(edition_uid=first.edition_uid))

    assert total == 1
    assert items[0].edition_uid == first.edition_uid
