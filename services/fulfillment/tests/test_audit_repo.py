from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import AuditCreateParams, AuditItemParam
from internal.domain.constant.audit_status import AuditStatus
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.error import InvalidTransitionError, StockNotFoundError, WarehouseNotFoundError
from internal.domain.model.option.inventory_option import AuditQueryOption, StockHistoryQueryOption
from internal.infrastructure.adapter.mysql.audit_repo import MySQLStockAuditRepository
from internal.infrastructure.adapter.mysql.stock_history_repo import MySQLStockHistoryRepository
from internal.infrastructure.adapter.mysql.stock_repo import MySQLStockRepository
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid
from tests.conftest import seed_stock, seed_warehouse


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLStockAuditRepository:
    return MySQLStockAuditRepository(session_factory)


@pytest.fixture
def stock_repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLStockRepository:
    return MySQLStockRepository(session_factory)


async def test_create_snapshots_system_quantity_and_difference(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)

    audit = await repo.create(
        AuditCreateParams(
            warehouse_uid=bytes_to_uuid(warehouse.uid),
            items=(AuditItemParam(edition_uid=bytes_to_uuid(stock.edition_uid), actual_quantity=8),),
            notes="cycle count",
        )
    )

    assert audit.status == AuditStatus.IN_PROGRESS
    assert audit.notes == "cycle count"
    assert len(audit.items) == 1
    assert audit.items[0].system_quantity == 10
    assert audit.items[0].actual_quantity == 8
    assert audit.items[0].difference == -2


async def test_create_missing_warehouse_raises_warehouse_not_found(repo) -> None:
    with pytest.raises(WarehouseNotFoundError):
        await repo.create(
            AuditCreateParams(
                warehouse_uid=uuid4(), items=(AuditItemParam(edition_uid=uuid4(), actual_quantity=1),), notes=None
            )
        )


async def test_create_missing_stock_raises_stock_not_found_and_persists_nothing(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)

    with pytest.raises(StockNotFoundError):
        await repo.create(
            AuditCreateParams(
                warehouse_uid=bytes_to_uuid(warehouse.uid),
                items=(AuditItemParam(edition_uid=uuid4(), actual_quantity=1),),
                notes=None,
            )
        )

    assert await repo.find_all(AuditQueryOption()) == ([], 0)


async def test_complete_applies_adjustments_and_marks_completed(repo, stock_repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    short = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    exact = await seed_stock(session_factory, warehouse, edition_id=2, quantity=5)
    audit = await repo.create(
        AuditCreateParams(
            warehouse_uid=bytes_to_uuid(warehouse.uid),
            items=(
                AuditItemParam(edition_uid=bytes_to_uuid(short.edition_uid), actual_quantity=7),
                AuditItemParam(edition_uid=bytes_to_uuid(exact.edition_uid), actual_quantity=5),
            ),
            notes=None,
        )
    )

    done = await repo.complete(audit.uid)

    assert done is not None
    assert done.status == AuditStatus.COMPLETED
    assert done.completed_at is not None
    adjusted = await stock_repo.find_by_edition_uid(bytes_to_uuid(short.edition_uid))
    assert adjusted is not None and adjusted.total_quantity == 7
    untouched = await stock_repo.find_by_edition_uid(bytes_to_uuid(exact.edition_uid))
    assert untouched is not None and untouched.total_quantity == 5
    history, _ = await MySQLStockHistoryRepository(session_factory).find_all(
        [bytes_to_uuid(short.uid)], StockHistoryQueryOption()
    )
    assert [(h.change_type, h.change_qty, h.reason) for h in history] == [
        (StockChangeType.ADJUSTMENT, -3, f"audit {audit.uid}")
    ]


async def test_complete_completed_audit_raises_invalid_transition(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    audit = await repo.create(
        AuditCreateParams(
            warehouse_uid=bytes_to_uuid(warehouse.uid),
            items=(AuditItemParam(edition_uid=bytes_to_uuid(stock.edition_uid), actual_quantity=10),),
            notes=None,
        )
    )
    await repo.complete(audit.uid)

    with pytest.raises(InvalidTransitionError):
        await repo.complete(audit.uid)


async def test_complete_missing_is_none(repo) -> None:
    assert await repo.complete(uuid4()) is None


async def test_find_by_uid_returns_items(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    created = await repo.create(
        AuditCreateParams(
            warehouse_uid=bytes_to_uuid(warehouse.uid),
            items=(AuditItemParam(edition_uid=bytes_to_uuid(stock.edition_uid), actual_quantity=9),),
            notes=None,
        )
    )

    found = await repo.find_by_uid(created.uid)

    assert found is not None
    assert [i.difference for i in found.items] == [-1]


async def test_find_all_by_status_filters_audits(repo, session_factory) -> None:
    warehouse = await seed_warehouse(session_factory)
    stock = await seed_stock(session_factory, warehouse, edition_id=1, quantity=10)
    params = AuditCreateParams(
        warehouse_uid=bytes_to_uuid(warehouse.uid),
        items=(AuditItemParam(edition_uid=bytes_to_uuid(stock.edition_uid), actual_quantity=10),),
        notes=None,
    )
    first = await repo.create(params)
    await repo.create(params)
    await repo.complete(first.uid)

    items, total = await repo.find_all(AuditQueryOption(status=AuditStatus.COMPLETED))

    assert total == 1
    assert items[0].uid == first.uid
