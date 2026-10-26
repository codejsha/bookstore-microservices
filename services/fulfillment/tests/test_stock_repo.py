from uuid import UUID, uuid4

import pytest
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import ApplyChangeParams
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.error import InsufficientStockError, StockNotFoundError, WarehouseNotFoundError
from internal.domain.model.option.inventory_option import StockQueryOption
from internal.infrastructure.adapter.pgsql.models import StockEntity, StockHistoryEntity
from internal.infrastructure.adapter.pgsql.stock_repo import PgsqlStockRepository
from tests.conftest import seed_stock, seed_warehouse


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> PgsqlStockRepository:
    return PgsqlStockRepository(session_factory)


async def _quantity(session_factory, stock_uid: UUID) -> tuple[int, int]:
    async with session_factory() as session:
        row = (await session.execute(select(StockEntity).where(StockEntity.uid == stock_uid))).scalar_one()
        return row.quantity, row.version


async def _history(session_factory, stock_uid: UUID) -> list[tuple[str, int]]:
    async with session_factory() as session:
        rows = (
            await session.execute(
                select(StockHistoryEntity)
                .where(StockHistoryEntity.stock_uid == stock_uid)
                .order_by(StockHistoryEntity.id)
            )
        ).scalars()
        return [(r.change_type, r.change_qty) for r in rows]


def _params(stock: StockEntity, delta: int, change_type: StockChangeType, reason: str | None = None):
    return ApplyChangeParams(
        edition_uid=stock.edition_uid,
        warehouse_uid=stock.warehouse_uid,
        delta=delta,
        change_type=change_type,
        reason=reason,
    )


async def test_apply_change_inbound_existing_row_increments_quantity_and_records_history(
    repo: PgsqlStockRepository, session_factory
) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)

    await repo.apply_change(_params(stock, 5, StockChangeType.INBOUND, "restock"))

    assert await _quantity(session_factory, stock.uid) == (15, 2)
    assert await _history(session_factory, stock.uid) == [("INBOUND", 5)]


async def test_apply_change_reservation_beyond_quantity_raises_insufficient_and_keeps_quantity(
    repo: PgsqlStockRepository, session_factory
) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=3)

    with pytest.raises(InsufficientStockError):
        await repo.apply_change(_params(stock, -4, StockChangeType.RESERVATION))

    assert await _quantity(session_factory, stock.uid) == (3, 1)
    assert await _history(session_factory, stock.uid) == []


async def test_apply_change_reservation_missing_row_raises_stock_not_found(
    repo: PgsqlStockRepository, session_factory
) -> None:
    warehouse = await seed_warehouse(session_factory)
    params = ApplyChangeParams(
        edition_uid=uuid4(),
        warehouse_uid=warehouse.uid,
        delta=-1,
        change_type=StockChangeType.RESERVATION,
        reason=None,
    )
    with pytest.raises(StockNotFoundError):
        await repo.apply_change(params)


async def test_apply_change_inbound_missing_warehouse_raises_warehouse_not_found(
    repo: PgsqlStockRepository, session_factory
) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=1)
    params = ApplyChangeParams(
        edition_uid=stock.edition_uid,
        warehouse_uid=uuid4(),
        delta=1,
        change_type=StockChangeType.INBOUND,
        reason=None,
    )
    with pytest.raises(WarehouseNotFoundError):
        await repo.apply_change(params)


async def test_apply_change_inbound_edition_without_any_stock_raises_stock_not_found(
    repo: PgsqlStockRepository, session_factory
) -> None:
    warehouse = await seed_warehouse(session_factory)
    params = ApplyChangeParams(
        edition_uid=uuid4(),
        warehouse_uid=warehouse.uid,
        delta=1,
        change_type=StockChangeType.INBOUND,
        reason=None,
    )
    with pytest.raises(StockNotFoundError):
        await repo.apply_change(params)


async def test_apply_change_inbound_new_warehouse_with_sibling_creates_row(
    repo: PgsqlStockRepository, session_factory
) -> None:
    first = await seed_warehouse(session_factory, name="A")
    second = await seed_warehouse(session_factory, name="B")
    sibling = await seed_stock(session_factory, first, edition_id=7, quantity=4)

    await repo.apply_change(
        ApplyChangeParams(
            edition_uid=sibling.edition_uid,
            warehouse_uid=second.uid,
            delta=6,
            change_type=StockChangeType.INBOUND,
            reason=None,
        )
    )

    stock = await repo.find_by_edition_uid(sibling.edition_uid)
    assert stock is not None
    assert stock.total_quantity == 10
    assert {w.warehouse_name: w.quantity for w in stock.warehouses} == {"A": 4, "B": 6}
    async with session_factory() as session:
        created = (await session.execute(select(StockEntity).where(StockEntity.warehouse_id == second.id))).scalar_one()
    assert created.edition_id == 7
    assert created.version == 2


async def test_find_all_groups_rows_by_edition_with_distinct_total(repo: PgsqlStockRepository, session_factory) -> None:
    first = await seed_warehouse(session_factory, name="A")
    second = await seed_warehouse(session_factory, name="B")
    edition = uuid4()
    await seed_stock(session_factory, first, edition_id=1, edition_uid=edition, quantity=3)
    await seed_stock(session_factory, second, edition_id=1, edition_uid=edition, quantity=4)
    await seed_stock(session_factory, first, edition_id=2, quantity=9)

    items, total = await repo.find_all(StockQueryOption())

    assert total == 2
    by_edition = {s.edition_uid: s for s in items}
    assert by_edition[edition].total_quantity == 7
    assert len(by_edition[edition].warehouses) == 2


async def test_find_all_by_warehouse_filters_rows(repo: PgsqlStockRepository, session_factory) -> None:
    first = await seed_warehouse(session_factory, name="A")
    second = await seed_warehouse(session_factory, name="B")
    await seed_stock(session_factory, first, edition_id=1, quantity=3)
    await seed_stock(session_factory, second, edition_id=2, quantity=4)

    items, total = await repo.find_all(StockQueryOption(warehouse_uid=second.uid))

    assert total == 1
    assert items[0].warehouses[0].warehouse_name == "B"


async def test_find_all_with_page_size_pages_editions(repo: PgsqlStockRepository, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    for edition_id in range(5):
        await seed_stock(session_factory, warehouse, edition_id=edition_id, quantity=1)

    items, total = await repo.find_all(StockQueryOption(page=1, size=2))

    assert total == 5
    assert len(items) == 2


async def test_find_by_edition_uid_missing_is_none(repo: PgsqlStockRepository) -> None:
    assert await repo.find_by_edition_uid(uuid4()) is None


async def test_find_stock_uids_by_warehouse_filters_rows(repo: PgsqlStockRepository, session_factory) -> None:
    first = await seed_warehouse(session_factory, name="A")
    second = await seed_warehouse(session_factory, name="B")
    edition = uuid4()
    a = await seed_stock(session_factory, first, edition_id=1, edition_uid=edition, quantity=3)
    await seed_stock(session_factory, second, edition_id=1, edition_uid=edition, quantity=4)

    assert len(await repo.find_stock_uids(edition, None)) == 2
    assert await repo.find_stock_uids(edition, first.uid) == [a.uid]
