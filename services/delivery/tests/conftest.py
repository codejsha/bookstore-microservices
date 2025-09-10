from collections.abc import AsyncIterator
from datetime import UTC, datetime
from unittest.mock import MagicMock
from uuid import UUID, uuid4

import pytest
import pytest_asyncio
from sqlalchemy import BigInteger
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.ext.compiler import compiles
from sqlalchemy.pool import StaticPool

from internal.application.port.repo.repos import (
    CarrierRepository,
    FreightRepository,
    ShipmentRepository,
    StatsRepository,
    TrackingRepository,
)
from internal.domain.aggregate.carrier_aggregate import CarrierAggregate
from internal.domain.aggregate.freight_aggregate import FreightAggregate
from internal.domain.aggregate.shipment_aggregate import ShipmentAggregate
from internal.domain.aggregate.tracking_aggregate import TrackingAggregate
from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.constant.freight_status import FreightStatus
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.infrastructure.adapter.mysql.models import Base


@compiles(BigInteger, "sqlite")
def _compile_bigint_sqlite(_type, _compiler, **_kw):
    return "INTEGER"


@pytest.fixture
def now() -> datetime:
    return datetime.now(UTC)


@pytest.fixture
def shipment_repo() -> MagicMock:
    repo = MagicMock(spec=ShipmentRepository)
    repo.save.side_effect = lambda s: s

    async def _create_with_initial_tracking(shipment, tracking):
        return shipment

    repo.create_with_initial_tracking.side_effect = _create_with_initial_tracking

    async def _transition(uid, guard, mutator):
        shipment = repo.find_by_uid.return_value
        if shipment is None:
            return None
        if not guard(shipment):
            raise ValueError(f"Cannot transition from {shipment.status}")
        return mutator(shipment)

    repo.transition.side_effect = _transition

    async def _update(uid, mutator):
        shipment = repo.find_by_uid.return_value
        if shipment is None:
            return None
        return mutator(shipment)

    repo.update.side_effect = _update
    return repo


@pytest.fixture
def tracking_repo() -> MagicMock:
    repo = MagicMock(spec=TrackingRepository)
    repo.save.side_effect = lambda t: t
    return repo


@pytest.fixture
def carrier_repo() -> MagicMock:
    repo = MagicMock(spec=CarrierRepository)
    repo.save.side_effect = lambda c: c

    async def _update(uid, mutator):
        carrier = repo.find_by_uid.return_value
        if carrier is None:
            return None
        return mutator(carrier)

    repo.update.side_effect = _update
    return repo


@pytest.fixture
def freight_repo() -> MagicMock:
    repo = MagicMock(spec=FreightRepository)
    repo.save.side_effect = lambda f: f

    async def _update(uid, mutator):
        freight = repo.find_by_uid.return_value
        if freight is None:
            return None
        return mutator(freight)

    repo.update.side_effect = _update
    return repo


@pytest.fixture
def stats_repo() -> MagicMock:
    return MagicMock(spec=StatsRepository)


@pytest_asyncio.fixture
async def session_factory() -> AsyncIterator[async_sessionmaker[AsyncSession]]:
    engine = create_async_engine(
        "sqlite+aiosqlite:///:memory:",
        connect_args={"check_same_thread": False},
        poolclass=StaticPool,
    )
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    factory = async_sessionmaker(bind=engine, expire_on_commit=False)
    try:
        yield factory
    finally:
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.drop_all)
        await engine.dispose()


def make_shipment(
    *,
    uid: UUID | None = None,
    order_uid: UUID | None = None,
    status: ShipmentStatus = ShipmentStatus.PLANNED,
    now: datetime | None = None,
) -> ShipmentAggregate:
    ts = now or datetime.now(UTC)
    return ShipmentAggregate(
        uid=uid or uuid4(),
        order_uid=order_uid or uuid4(),
        carrier_uid=None,
        origin_address="origin",
        destination_address="dest",
        destination_city="Seoul",
        destination_state="KR",
        destination_country_code="KR",
        destination_postal_code="00000",
        status=status,
        tracking_number=None,
        weight_kg=1.0,
        planned_pickup_at=None,
        planned_delivery_at=None,
        actual_pickup_at=None,
        actual_delivery_at=None,
        created_at=ts,
        updated_at=ts,
    )


def make_carrier(
    *,
    uid: UUID | None = None,
    status: CarrierStatus = CarrierStatus.ACTIVE,
    now: datetime | None = None,
) -> CarrierAggregate:
    ts = now or datetime.now(UTC)
    return CarrierAggregate(
        uid=uid or uuid4(),
        name="ACME",
        code="ACME",
        contact_name="John",
        contact_phone="0000",
        contact_email="john@example.com",
        base_rate=1000.0,
        rate_per_kg=500.0,
        status=status,
        created_at=ts,
        updated_at=ts,
    )


def make_freight(
    *,
    uid: UUID | None = None,
    shipment_uid: UUID | None = None,
    carrier_uid: UUID | None = None,
    status: FreightStatus = FreightStatus.ESTIMATED,
    now: datetime | None = None,
) -> FreightAggregate:
    ts = now or datetime.now(UTC)
    return FreightAggregate(
        uid=uid or uuid4(),
        shipment_uid=shipment_uid or uuid4(),
        carrier_uid=carrier_uid or uuid4(),
        base_cost=1000.0,
        weight_surcharge=100.0,
        distance_surcharge=200.0,
        discount=50.0,
        total_cost=1250.0,
        currency="KRW",
        status=status,
        invoiced_at=None,
        paid_at=None,
        created_at=ts,
        updated_at=ts,
    )


def make_tracking(
    *,
    shipment_uid: UUID | None = None,
    status: ShipmentStatus = ShipmentStatus.PLANNED,
    now: datetime | None = None,
) -> TrackingAggregate:
    ts = now or datetime.now(UTC)
    return TrackingAggregate(
        uid=uuid4(),
        shipment_uid=shipment_uid or uuid4(),
        status=status,
        location="",
        description="",
        occurred_at=ts,
        created_at=ts,
    )
