from datetime import UTC, datetime, timedelta
from uuid import uuid4

import pytest
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.option.delivery_option import ShipmentFilterOption
from internal.infrastructure.adapter.mysql.shipment_repo import MySQLShipmentRepository
from internal.infrastructure.adapter.mysql.tracking_repo import MySQLTrackingRepository
from tests.conftest import make_shipment, make_tracking


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLShipmentRepository:
    return MySQLShipmentRepository(session_factory)


@pytest.fixture
def tracking_reader(session_factory: async_sessionmaker[AsyncSession]) -> MySQLTrackingRepository:
    return MySQLTrackingRepository(session_factory)


async def test_find_by_uid_has_timezone_aware_utc_datetimes(repo: MySQLShipmentRepository) -> None:
    shipment = make_shipment(now=datetime.now(UTC))
    await repo.save(shipment)
    found = await repo.find_by_uid(shipment.uid)
    assert found is not None
    assert found.created_at.tzinfo is not None
    assert found.created_at.utcoffset() == timedelta(0)
    assert found.created_at.isoformat().endswith("+00:00")
    assert found.updated_at.isoformat().endswith("+00:00")


async def test_create_with_initial_tracking_both_valid_persists_both(
    repo: MySQLShipmentRepository, tracking_reader: MySQLTrackingRepository
) -> None:
    shipment = make_shipment()
    tracking = make_tracking(shipment_uid=shipment.uid, status=ShipmentStatus.PLANNED)
    saved = await repo.create_with_initial_tracking(shipment, tracking)
    assert saved.uid == shipment.uid

    found = await repo.find_by_uid(shipment.uid)
    assert found is not None
    history = await tracking_reader.find_by_shipment_uid(shipment.uid)
    assert len(history) == 1
    assert history[0].status == ShipmentStatus.PLANNED


async def test_create_with_initial_tracking_duplicate_order_rolls_back_both(
    repo: MySQLShipmentRepository, tracking_reader: MySQLTrackingRepository
) -> None:
    order_uid = uuid4()
    first = make_shipment(order_uid=order_uid)
    await repo.create_with_initial_tracking(first, make_tracking(shipment_uid=first.uid))

    dup = make_shipment(order_uid=order_uid)
    with pytest.raises(IntegrityError):
        await repo.create_with_initial_tracking(dup, make_tracking(shipment_uid=dup.uid))

    assert await tracking_reader.find_by_shipment_uid(dup.uid) == []


async def test_save_new_inserts_row(repo: MySQLShipmentRepository) -> None:
    shipment = make_shipment()
    saved = await repo.save(shipment)
    assert saved.uid == shipment.uid
    assert saved.status == ShipmentStatus.PLANNED


async def test_find_by_uid_roundtrips(repo: MySQLShipmentRepository) -> None:
    shipment = make_shipment()
    await repo.save(shipment)
    found = await repo.find_by_uid(shipment.uid)
    assert found is not None
    assert found.uid == shipment.uid
    assert found.order_uid == shipment.order_uid


async def test_find_by_uid_missing_is_none(repo: MySQLShipmentRepository) -> None:
    assert await repo.find_by_uid(uuid4()) is None


async def test_save_found_updates_row(repo: MySQLShipmentRepository) -> None:
    shipment = make_shipment()
    await repo.save(shipment)
    carrier_uid = uuid4()
    shipment.carrier_uid = carrier_uid
    shipment.status = ShipmentStatus.DISPATCHED
    shipment.tracking_number = "TRK-1"
    shipment.updated_at = datetime.now(UTC)
    updated = await repo.save(shipment)
    assert updated.carrier_uid == carrier_uid
    assert updated.status == ShipmentStatus.DISPATCHED
    assert updated.tracking_number == "TRK-1"


async def test_find_by_order_uid_roundtrips(repo: MySQLShipmentRepository) -> None:
    order_uid = uuid4()
    shipment = make_shipment(order_uid=order_uid)
    await repo.save(shipment)
    found = await repo.find_by_order_uid(order_uid)
    assert found is not None
    assert found.uid == shipment.uid
    assert found.order_uid == order_uid


async def test_find_by_order_uid_missing_is_none(repo: MySQLShipmentRepository) -> None:
    assert await repo.find_by_order_uid(uuid4()) is None


async def test_save_duplicate_order_uid_raises_integrity_error(repo: MySQLShipmentRepository) -> None:
    order_uid = uuid4()
    await repo.save(make_shipment(order_uid=order_uid))
    with pytest.raises(IntegrityError):
        await repo.save(make_shipment(order_uid=order_uid))


async def test_find_all_by_status_matches_shipments(repo: MySQLShipmentRepository) -> None:
    await repo.save(make_shipment(status=ShipmentStatus.PLANNED))
    await repo.save(make_shipment(status=ShipmentStatus.DELIVERED))
    await repo.save(make_shipment(status=ShipmentStatus.PLANNED))

    items, total = await repo.find_all(ShipmentFilterOption(status=ShipmentStatus.PLANNED))
    assert total == 2
    assert len(items) == 2
    assert all(item.status == ShipmentStatus.PLANNED for item in items)


async def test_find_all_by_order_uid_matches_shipments(repo: MySQLShipmentRepository) -> None:
    target_order = uuid4()
    await repo.save(make_shipment(order_uid=target_order))
    await repo.save(make_shipment())

    items, total = await repo.find_all(ShipmentFilterOption(order_uid=target_order))
    assert total == 1
    assert items[0].order_uid == target_order


async def test_find_all_with_page_size_pages_with_total(repo: MySQLShipmentRepository) -> None:
    base = datetime.now(UTC)
    for i in range(5):
        await repo.save(make_shipment(now=base + timedelta(seconds=i)))
    items, total = await repo.find_all(ShipmentFilterOption(page=0, size=2))
    assert total == 5
    assert len(items) == 2


async def test_find_all_sorted_by_created_at_desc_is_newest_first(repo: MySQLShipmentRepository) -> None:
    base = datetime.now(UTC)
    older = make_shipment(now=base - timedelta(days=1))
    newer = make_shipment(now=base)
    await repo.save(older)
    await repo.save(newer)

    items, _ = await repo.find_all(ShipmentFilterOption(sort="created_at:desc"))
    assert items[0].uid == newer.uid
    assert items[1].uid == older.uid


async def test_find_all_unknown_sort_field_falls_back_to_default_order(repo: MySQLShipmentRepository) -> None:
    base = datetime.now(UTC)
    older = make_shipment(now=base - timedelta(days=1))
    newer = make_shipment(now=base)
    await repo.save(older)
    await repo.save(newer)

    items, total = await repo.find_all(ShipmentFilterOption(sort="metadata:desc"))
    assert total == 2
    assert items[0].uid == newer.uid
    assert items[1].uid == older.uid


async def test_update_found_applies_mutator(repo: MySQLShipmentRepository) -> None:
    shipment = make_shipment()
    await repo.save(shipment)
    carrier_uid = uuid4()

    def _mutate(s):
        s.carrier_uid = carrier_uid
        s.tracking_number = "TRK-9"
        s.updated_at = datetime.now(UTC)
        return s

    updated = await repo.update(shipment.uid, _mutate)
    assert updated is not None
    assert updated.carrier_uid == carrier_uid
    assert updated.tracking_number == "TRK-9"
    reloaded = await repo.find_by_uid(shipment.uid)
    assert reloaded is not None
    assert reloaded.tracking_number == "TRK-9"


async def test_update_missing_is_none(repo: MySQLShipmentRepository) -> None:
    assert await repo.update(uuid4(), lambda s: s) is None
