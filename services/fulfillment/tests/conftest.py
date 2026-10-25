from collections.abc import AsyncIterator
from datetime import UTC, datetime
from unittest.mock import MagicMock
from uuid import UUID, uuid4, uuid7

import pytest
import pytest_asyncio
from sqlalchemy import BigInteger
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.ext.compiler import compiles
from sqlalchemy.pool import StaticPool

from internal.application.port.event.event_publisher import EventPublisher
from internal.application.port.repo.inventory_repos import (
    MonthlyClosingRepository,
    StockAuditRepository,
    StockBalanceRepository,
    StockHistoryRepository,
    StockRepository,
    StockReservationRepository,
    StockTransferRepository,
    WarehouseRepository,
)
from internal.application.port.repo.repos import (
    CarrierRepository,
    FreightRepository,
    ShipmentRepository,
    StatsRepository,
    TrackingRepository,
)
from internal.domain.aggregate.audit_aggregate import StockAuditAggregate, StockAuditItem
from internal.domain.aggregate.balance_aggregate import StockBalanceEntry
from internal.domain.aggregate.carrier_aggregate import CarrierAggregate
from internal.domain.aggregate.closing_aggregate import MonthlyClosingAggregate, MonthlyClosingItem
from internal.domain.aggregate.freight_aggregate import FreightAggregate
from internal.domain.aggregate.reservation_aggregate import StockReservationAggregate
from internal.domain.aggregate.shipment_aggregate import ShipmentAggregate
from internal.domain.aggregate.stock_aggregate import StockAggregate, StockWarehouse
from internal.domain.aggregate.tracking_aggregate import TrackingAggregate
from internal.domain.aggregate.transfer_aggregate import StockTransferAggregate
from internal.domain.aggregate.warehouse_aggregate import WarehouseAggregate
from internal.domain.constant.audit_status import AuditStatus
from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.constant.closing_status import ClosingStatus
from internal.domain.constant.freight_status import FreightStatus
from internal.domain.constant.reservation_status import ReservationStatus
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.constant.transfer_status import TransferStatus
from internal.infrastructure.adapter.mysql.models import Base, StockEntity, WarehouseEntity


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
    repo.find_by_order_uid.return_value = None

    async def _create_with_initial_tracking(shipment, tracking):
        return shipment

    repo.create_with_initial_tracking.side_effect = _create_with_initial_tracking

    async def _transition_with_tracking(uid, guard, mutator, tracking_factory):
        shipment = repo.find_by_uid.return_value
        if shipment is None:
            return None
        if not guard(shipment):
            raise ValueError(f"Cannot transition from {shipment.status}")
        updated = mutator(shipment)
        tracking_factory(updated)
        return updated

    repo.transition_with_tracking.side_effect = _transition_with_tracking

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
    repo.find_by_code.return_value = None

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


@pytest.fixture
def event_publisher() -> MagicMock:
    return MagicMock(spec=EventPublisher)


@pytest.fixture
def stock_repo() -> MagicMock:
    return MagicMock(spec=StockRepository)


@pytest.fixture
def stock_history_repo() -> MagicMock:
    return MagicMock(spec=StockHistoryRepository)


@pytest.fixture
def reservation_repo() -> MagicMock:
    return MagicMock(spec=StockReservationRepository)


@pytest.fixture
def warehouse_repo() -> MagicMock:
    return MagicMock(spec=WarehouseRepository)


@pytest.fixture
def transfer_repo() -> MagicMock:
    return MagicMock(spec=StockTransferRepository)


@pytest.fixture
def audit_repo() -> MagicMock:
    return MagicMock(spec=StockAuditRepository)


@pytest.fixture
def closing_repo() -> MagicMock:
    return MagicMock(spec=MonthlyClosingRepository)


@pytest.fixture
def balance_repo() -> MagicMock:
    return MagicMock(spec=StockBalanceRepository)


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


def make_warehouse(
    *,
    uid: UUID | None = None,
    name: str = "Main",
    capacity: int = 100,
    now: datetime | None = None,
) -> WarehouseAggregate:
    ts = now or datetime.now(UTC)
    return WarehouseAggregate(
        uid=uid or uuid4(),
        name=name,
        address="1 Depot Rd",
        capacity=capacity,
        created_at=ts,
        updated_at=None,
    )


def make_stock(
    *,
    edition_uid: UUID | None = None,
    warehouse_uid: UUID | None = None,
    quantity: int = 10,
) -> StockAggregate:
    return StockAggregate(
        uid=uuid4(),
        edition_uid=edition_uid or uuid4(),
        total_quantity=quantity,
        warehouses=[
            StockWarehouse(warehouse_uid=warehouse_uid or uuid4(), warehouse_name="Main", quantity=quantity),
        ],
    )


def make_transfer(
    *,
    uid: UUID | None = None,
    status: TransferStatus = TransferStatus.PENDING,
    quantity: int = 5,
    now: datetime | None = None,
) -> StockTransferAggregate:
    ts = now or datetime.now(UTC)
    return StockTransferAggregate(
        uid=uid or uuid4(),
        edition_uid=uuid4(),
        source_warehouse_uid=uuid4(),
        target_warehouse_uid=uuid4(),
        quantity=quantity,
        status=status,
        reason=None,
        completed_at=None,
        created_at=ts,
        updated_at=None,
    )


def make_audit(
    *,
    uid: UUID | None = None,
    status: AuditStatus = AuditStatus.IN_PROGRESS,
    now: datetime | None = None,
) -> StockAuditAggregate:
    ts = now or datetime.now(UTC)
    return StockAuditAggregate(
        uid=uid or uuid4(),
        warehouse_uid=uuid4(),
        status=status,
        notes=None,
        completed_at=None,
        created_at=ts,
        updated_at=None,
        items=[StockAuditItem(edition_uid=uuid4(), system_quantity=10, actual_quantity=8, difference=-2)],
    )


def make_closing(*, uid: UUID | None = None, now: datetime | None = None) -> MonthlyClosingAggregate:
    ts = now or datetime.now(UTC)
    return MonthlyClosingAggregate(
        uid=uid or uuid4(),
        warehouse_uid=uuid4(),
        year=2026,
        month=10,
        status=ClosingStatus.CLOSED,
        closed_at=ts,
        created_at=ts,
        updated_at=None,
        items=[
            MonthlyClosingItem(
                edition_uid=uuid4(),
                opening_quantity=10,
                inbound_quantity=5,
                outbound_quantity=3,
                adjust_quantity=0,
                closing_quantity=12,
            )
        ],
    )


def make_balance_entry() -> StockBalanceEntry:
    return StockBalanceEntry(
        edition_uid=uuid4(),
        warehouse_uid=uuid4(),
        warehouse_name="Main",
        inbound_quantity=5,
        outbound_quantity=3,
        adjust_quantity=0,
        current_quantity=12,
    )


def make_reservation(
    *,
    edition_uid: UUID | None = None,
    quantity: int = 3,
    status: ReservationStatus = ReservationStatus.RESERVED,
    already_applied: bool = False,
) -> StockReservationAggregate:
    return StockReservationAggregate(
        uid=uuid4(),
        reservation_key="order-1:1",
        edition_uid=edition_uid or uuid4(),
        warehouse_uid=uuid4(),
        quantity=quantity,
        status=status,
        already_applied=already_applied,
    )


async def seed_warehouse(session_factory: async_sessionmaker[AsyncSession], name: str = "Main") -> WarehouseEntity:
    entity = WarehouseEntity(
        uid=uuid7().bytes,
        name=name,
        address=None,
        capacity=1000,
        created_at=datetime.now(UTC),
        updated_at=None,
        deleted_at=None,
        version=1,
    )
    async with session_factory() as session:
        session.add(entity)
        await session.commit()
    return entity


async def seed_stock(
    session_factory: async_sessionmaker[AsyncSession],
    warehouse: WarehouseEntity,
    edition_id: int,
    edition_uid: UUID | None = None,
    quantity: int = 10,
) -> StockEntity:
    entity = StockEntity(
        uid=uuid7().bytes,
        edition_id=edition_id,
        edition_uid=(edition_uid or uuid4()).bytes,
        warehouse_id=warehouse.id,
        warehouse_uid=warehouse.uid,
        quantity=quantity,
        created_at=datetime.now(UTC),
        updated_at=None,
        deleted_at=None,
        version=1,
    )
    async with session_factory() as session:
        session.add(entity)
        await session.commit()
    return entity
