from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import WarehouseCreateParams, WarehouseUpdateParams
from internal.domain.model.option.inventory_option import WarehouseQueryOption
from internal.infrastructure.adapter.pgsql.warehouse_repo import PgsqlWarehouseRepository


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> PgsqlWarehouseRepository:
    return PgsqlWarehouseRepository(session_factory)


async def test_create_then_find_by_uid_roundtrips(repo: PgsqlWarehouseRepository) -> None:
    created = await repo.create(WarehouseCreateParams(name="Main", address="1 Depot Rd", capacity=100))
    found = await repo.find_by_uid(created.uid)
    assert found is not None
    assert found.name == "Main"
    assert found.address == "1 Depot Rd"
    assert found.capacity == 100
    assert found.updated_at is None
    assert found.created_at.tzinfo is not None


async def test_find_by_uid_missing_is_none(repo: PgsqlWarehouseRepository) -> None:
    assert await repo.find_by_uid(uuid4()) is None


async def test_update_partial_fields_keeps_others_and_sets_updated_at(repo: PgsqlWarehouseRepository) -> None:
    created = await repo.create(WarehouseCreateParams(name="Main", address="1 Depot Rd", capacity=100))

    updated = await repo.update(WarehouseUpdateParams(uid=created.uid, name="Renamed", address=None, capacity=None))

    assert updated is not None
    assert updated.name == "Renamed"
    assert updated.address == "1 Depot Rd"
    assert updated.capacity == 100
    assert updated.updated_at is not None
    reloaded = await repo.find_by_uid(created.uid)
    assert reloaded is not None
    assert reloaded.name == "Renamed"


async def test_update_missing_is_none(repo: PgsqlWarehouseRepository) -> None:
    assert await repo.update(WarehouseUpdateParams(uid=uuid4(), name="x", address=None, capacity=None)) is None


async def test_find_all_by_name_matches_case_insensitively(repo: PgsqlWarehouseRepository) -> None:
    await repo.create(WarehouseCreateParams(name="Seoul Depot", address=None, capacity=1))
    await repo.create(WarehouseCreateParams(name="Busan Depot", address=None, capacity=1))

    items, total = await repo.find_all(WarehouseQueryOption(name="seoul"))

    assert total == 1
    assert items[0].name == "Seoul Depot"


async def test_find_all_sorted_by_name_asc_orders_items(repo: PgsqlWarehouseRepository) -> None:
    await repo.create(WarehouseCreateParams(name="B", address=None, capacity=1))
    await repo.create(WarehouseCreateParams(name="A", address=None, capacity=1))

    items, _ = await repo.find_all(WarehouseQueryOption(sort="name,asc"))

    assert [w.name for w in items] == ["A", "B"]


async def test_find_all_unknown_sort_falls_back_to_newest_first(repo: PgsqlWarehouseRepository) -> None:
    await repo.create(WarehouseCreateParams(name="older", address=None, capacity=1))
    await repo.create(WarehouseCreateParams(name="newer", address=None, capacity=1))

    items, total = await repo.find_all(WarehouseQueryOption(sort="metadata,desc"))

    assert total == 2
    assert items[0].name == "newer"


async def test_find_all_with_page_size_pages_with_total(repo: PgsqlWarehouseRepository) -> None:
    for i in range(5):
        await repo.create(WarehouseCreateParams(name=f"W{i}", address=None, capacity=1))

    items, total = await repo.find_all(WarehouseQueryOption(page=1, size=2))

    assert total == 5
    assert len(items) == 2
