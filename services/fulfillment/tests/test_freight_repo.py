from datetime import UTC, datetime
from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.constant.freight_status import FreightStatus
from internal.domain.model.option.delivery_option import FreightFilterOption
from internal.infrastructure.adapter.mysql.freight_repo import MySQLFreightRepository
from tests.conftest import make_freight


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLFreightRepository:
    return MySQLFreightRepository(session_factory)


async def test_save_new_inserts_row(repo: MySQLFreightRepository) -> None:
    freight = make_freight()
    saved = await repo.save(freight)
    assert saved.uid == freight.uid
    assert saved.status == FreightStatus.ESTIMATED


async def test_find_by_uid_roundtrips(repo: MySQLFreightRepository) -> None:
    freight = make_freight()
    await repo.save(freight)
    found = await repo.find_by_uid(freight.uid)
    assert found is not None
    assert found.total_cost == freight.total_cost


async def test_find_by_shipment_uid_roundtrips(repo: MySQLFreightRepository) -> None:
    shipment_uid = uuid4()
    freight = make_freight(shipment_uid=shipment_uid)
    await repo.save(freight)
    found = await repo.find_by_shipment_uid(shipment_uid)
    assert found is not None
    assert found.shipment_uid == shipment_uid


async def test_save_found_updates_status_and_timestamps(repo: MySQLFreightRepository) -> None:
    freight = make_freight(status=FreightStatus.ESTIMATED)
    await repo.save(freight)
    freight.status = FreightStatus.INVOICED
    freight.invoiced_at = datetime.now(UTC)
    freight.updated_at = datetime.now(UTC)
    updated = await repo.save(freight)
    assert updated.status == FreightStatus.INVOICED
    assert updated.invoiced_at is not None


async def test_find_all_by_status_and_carrier_matches_freights(repo: MySQLFreightRepository) -> None:
    carrier_a = uuid4()
    carrier_b = uuid4()
    await repo.save(make_freight(carrier_uid=carrier_a, status=FreightStatus.ESTIMATED))
    await repo.save(make_freight(carrier_uid=carrier_a, status=FreightStatus.PAID))
    await repo.save(make_freight(carrier_uid=carrier_b, status=FreightStatus.ESTIMATED))

    items, total = await repo.find_all(FreightFilterOption(carrier_uid=carrier_a, status=FreightStatus.ESTIMATED))
    assert total == 1
    assert items[0].carrier_uid == carrier_a
    assert items[0].status == FreightStatus.ESTIMATED
