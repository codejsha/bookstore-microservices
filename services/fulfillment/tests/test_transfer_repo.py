from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import TransferCreateParams
from internal.domain.constant.transfer_status import TransferStatus
from internal.domain.model.error import InvalidTransitionError, StockNotFoundError, WarehouseNotFoundError
from internal.domain.model.option.inventory_option import TransferQueryOption
from internal.infrastructure.adapter.mysql.stock_repo import MySQLStockRepository
from internal.infrastructure.adapter.mysql.transfer_repo import MySQLStockTransferRepository
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid
from tests.conftest import seed_stock, seed_warehouse


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLStockTransferRepository:
    return MySQLStockTransferRepository(session_factory)


@pytest.fixture
def stock_repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLStockRepository:
    return MySQLStockRepository(session_factory)


async def _setup(session_factory, source_qty: int = 10):
    source = await seed_warehouse(session_factory, name="source")
    target = await seed_warehouse(session_factory, name="target")
    stock = await seed_stock(session_factory, source, edition_id=1, quantity=source_qty)
    return source, target, stock


def _params(stock, source, target, quantity: int = 4) -> TransferCreateParams:
    return TransferCreateParams(
        edition_uid=bytes_to_uuid(stock.edition_uid),
        source_warehouse_uid=bytes_to_uuid(source.uid),
        target_warehouse_uid=bytes_to_uuid(target.uid),
        quantity=quantity,
        reason="rebalance",
    )


async def _warehouse_quantities(stock_repo, edition_uid) -> dict[str, int]:
    stock = await stock_repo.find_by_edition_uid(edition_uid)
    assert stock is not None
    return {w.warehouse_name: w.quantity for w in stock.warehouses}


async def test_create_deducts_source_and_records_pending_transfer(repo, stock_repo, session_factory) -> None:
    source, target, stock = await _setup(session_factory)

    transfer = await repo.create(_params(stock, source, target))

    assert transfer.status == TransferStatus.PENDING
    assert transfer.quantity == 4
    assert transfer.reason == "rebalance"
    assert await _warehouse_quantities(stock_repo, transfer.edition_uid) == {"source": 6}


async def test_create_missing_source_stock_raises_stock_not_found(repo, session_factory) -> None:
    source, target, stock = await _setup(session_factory)
    params = TransferCreateParams(
        edition_uid=uuid4(),
        source_warehouse_uid=bytes_to_uuid(source.uid),
        target_warehouse_uid=bytes_to_uuid(target.uid),
        quantity=1,
        reason=None,
    )
    with pytest.raises(StockNotFoundError):
        await repo.create(params)


async def test_create_missing_target_warehouse_raises_warehouse_not_found(repo, session_factory) -> None:
    source, _, stock = await _setup(session_factory)
    params = TransferCreateParams(
        edition_uid=bytes_to_uuid(stock.edition_uid),
        source_warehouse_uid=bytes_to_uuid(source.uid),
        target_warehouse_uid=uuid4(),
        quantity=1,
        reason=None,
    )
    with pytest.raises(WarehouseNotFoundError):
        await repo.create(params)


async def test_complete_pending_creates_target_row_and_marks_completed(repo, stock_repo, session_factory) -> None:
    source, target, stock = await _setup(session_factory)
    transfer = await repo.create(_params(stock, source, target))

    done = await repo.complete(transfer.uid)

    assert done is not None
    assert done.status == TransferStatus.COMPLETED
    assert done.completed_at is not None
    assert done.updated_at is not None
    assert await _warehouse_quantities(stock_repo, transfer.edition_uid) == {"source": 6, "target": 4}
    reloaded = await repo.find_by_uid(transfer.uid)
    assert reloaded is not None
    assert reloaded.status == TransferStatus.COMPLETED


async def test_complete_completed_transfer_raises_invalid_transition(repo, session_factory) -> None:
    source, target, stock = await _setup(session_factory)
    transfer = await repo.create(_params(stock, source, target))
    await repo.complete(transfer.uid)

    with pytest.raises(InvalidTransitionError):
        await repo.complete(transfer.uid)


async def test_cancel_pending_restores_source_and_marks_cancelled(repo, stock_repo, session_factory) -> None:
    source, target, stock = await _setup(session_factory)
    transfer = await repo.create(_params(stock, source, target))

    done = await repo.cancel(transfer.uid)

    assert done is not None
    assert done.status == TransferStatus.CANCELLED
    assert done.completed_at is None
    assert await _warehouse_quantities(stock_repo, transfer.edition_uid) == {"source": 10}


async def test_cancel_completed_transfer_raises_invalid_transition(repo, session_factory) -> None:
    source, target, stock = await _setup(session_factory)
    transfer = await repo.create(_params(stock, source, target))
    await repo.complete(transfer.uid)

    with pytest.raises(InvalidTransitionError):
        await repo.cancel(transfer.uid)


async def test_complete_missing_is_none(repo) -> None:
    assert await repo.complete(uuid4()) is None


async def test_find_all_by_status_filters_transfers(repo, session_factory) -> None:
    source, target, stock = await _setup(session_factory, source_qty=20)
    first = await repo.create(_params(stock, source, target, quantity=2))
    await repo.create(_params(stock, source, target, quantity=3))
    await repo.complete(first.uid)

    items, total = await repo.find_all(TransferQueryOption(status=TransferStatus.PENDING))

    assert total == 1
    assert items[0].quantity == 3
