from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import ApplyChangeParams
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.option.inventory_option import StockHistoryQueryOption
from internal.infrastructure.adapter.pgsql.stock_history_repo import PgsqlStockHistoryRepository
from internal.infrastructure.adapter.pgsql.stock_repo import PgsqlStockRepository
from tests.conftest import seed_stock, seed_warehouse


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> PgsqlStockHistoryRepository:
    return PgsqlStockHistoryRepository(session_factory)


@pytest.fixture
def stock_repo(session_factory: async_sessionmaker[AsyncSession]) -> PgsqlStockRepository:
    return PgsqlStockRepository(session_factory)


async def test_find_all_empty_uids_is_empty(repo: PgsqlStockHistoryRepository) -> None:
    assert await repo.find_all([], StockHistoryQueryOption()) == ([], 0)


async def test_find_all_lists_newest_first_with_total(
    repo: PgsqlStockHistoryRepository, stock_repo: PgsqlStockRepository, session_factory
) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    edition_uid = stock.edition_uid
    warehouse_uid = stock.warehouse_uid
    for delta, change_type in (
        (5, StockChangeType.INBOUND),
        (-2, StockChangeType.RESERVATION),
        (3, StockChangeType.ADJUSTMENT),
    ):
        await stock_repo.apply_change(
            ApplyChangeParams(
                edition_uid=edition_uid, warehouse_uid=warehouse_uid, delta=delta, change_type=change_type, reason=None
            )
        )

    items, total = await repo.find_all([stock.uid], StockHistoryQueryOption(size=2))

    assert total == 3
    assert len(items) == 2
    assert items[0].change_type == StockChangeType.ADJUSTMENT
    assert items[0].change_qty == 3
    assert items[1].change_type == StockChangeType.RESERVATION


async def test_find_all_unknown_stock_uid_is_empty(repo: PgsqlStockHistoryRepository) -> None:
    assert await repo.find_all([uuid4()], StockHistoryQueryOption()) == ([], 0)
