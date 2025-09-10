from datetime import UTC, datetime
from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.model.option.delivery_option import CarrierFilterOption
from internal.infrastructure.adapter.mysql.carrier_repo import MySQLCarrierRepository
from tests.conftest import make_carrier


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLCarrierRepository:
    return MySQLCarrierRepository(session_factory)


async def test_save_inserts_new_carrier(repo: MySQLCarrierRepository) -> None:
    carrier = make_carrier()
    saved = await repo.save(carrier)
    assert saved.uid == carrier.uid
    assert saved.status == CarrierStatus.ACTIVE


async def test_find_by_uid_returns_carrier(repo: MySQLCarrierRepository) -> None:
    carrier = make_carrier()
    await repo.save(carrier)
    found = await repo.find_by_uid(carrier.uid)
    assert found is not None
    assert found.code == carrier.code


async def test_find_by_uid_returns_none_when_missing(repo: MySQLCarrierRepository) -> None:
    assert await repo.find_by_uid(uuid4()) is None


async def test_save_updates_existing_carrier(repo: MySQLCarrierRepository) -> None:
    carrier = make_carrier()
    await repo.save(carrier)
    carrier.name = "Renamed"
    carrier.status = CarrierStatus.INACTIVE
    carrier.updated_at = datetime.now(UTC)
    updated = await repo.save(carrier)
    assert updated.name == "Renamed"
    assert updated.status == CarrierStatus.INACTIVE


async def test_find_all_filters_by_status(repo: MySQLCarrierRepository) -> None:
    a = make_carrier(status=CarrierStatus.ACTIVE)
    b = make_carrier(status=CarrierStatus.INACTIVE)
    a.code = "A"
    b.code = "B"
    await repo.save(a)
    await repo.save(b)
    items, total = await repo.find_all(CarrierFilterOption(status=CarrierStatus.ACTIVE))
    assert total == 1
    assert items[0].status == CarrierStatus.ACTIVE


async def test_find_all_filters_by_name_substring(repo: MySQLCarrierRepository) -> None:
    a = make_carrier()
    a.name = "FedEx Express"
    a.code = "FEDEX"
    b = make_carrier()
    b.name = "UPS"
    b.code = "UPS"
    await repo.save(a)
    await repo.save(b)
    items, total = await repo.find_all(CarrierFilterOption(name="Express"))
    assert total == 1
    assert items[0].name == "FedEx Express"


async def test_find_all_orders_by_name_asc(repo: MySQLCarrierRepository) -> None:
    a = make_carrier()
    a.name = "Beta"
    a.code = "B"
    b = make_carrier()
    b.name = "Alpha"
    b.code = "A"
    await repo.save(a)
    await repo.save(b)
    items, _ = await repo.find_all(CarrierFilterOption(sort="name:asc"))
    assert [c.name for c in items] == ["Alpha", "Beta"]
