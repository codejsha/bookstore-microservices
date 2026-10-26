from datetime import UTC, datetime, timedelta
from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.constant.shipment_status import ShipmentStatus
from internal.infrastructure.adapter.mysql.tracking_repo import MySQLTrackingRepository
from tests.conftest import make_tracking


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLTrackingRepository:
    return MySQLTrackingRepository(session_factory)


async def test_save_new_inserts_row(repo: MySQLTrackingRepository) -> None:
    tracking = make_tracking()
    saved = await repo.save(tracking)
    assert saved.uid == tracking.uid
    assert saved.shipment_uid == tracking.shipment_uid


async def test_find_by_shipment_uid_no_events_is_empty(repo: MySQLTrackingRepository) -> None:
    assert await repo.find_by_shipment_uid(uuid4()) == []


async def test_find_by_shipment_uid_orders_events(repo: MySQLTrackingRepository) -> None:
    shipment_uid = uuid4()
    base = datetime.now(UTC)
    later = make_tracking(shipment_uid=shipment_uid, status=ShipmentStatus.DISPATCHED, now=base + timedelta(hours=1))
    earlier = make_tracking(shipment_uid=shipment_uid, status=ShipmentStatus.PLANNED, now=base)
    await repo.save(later)
    await repo.save(earlier)

    events = await repo.find_by_shipment_uid(shipment_uid)
    assert [e.status for e in events] == [ShipmentStatus.PLANNED, ShipmentStatus.DISPATCHED]


async def test_find_by_shipment_uid_excludes_other_shipments_events(repo: MySQLTrackingRepository) -> None:
    shipment_a = uuid4()
    shipment_b = uuid4()
    await repo.save(make_tracking(shipment_uid=shipment_a))
    await repo.save(make_tracking(shipment_uid=shipment_b))

    events_a = await repo.find_by_shipment_uid(shipment_a)
    assert len(events_a) == 1
    assert events_a[0].shipment_uid == shipment_a
