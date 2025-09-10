from __future__ import annotations

import os
from collections.abc import AsyncIterator
from datetime import UTC, datetime, timedelta
from uuid import UUID, uuid4

import pytest
import pytest_asyncio
from sqlalchemy import LargeBinary
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.ext.compiler import compiles

from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.constant.freight_status import FreightStatus
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.option.delivery_option import StatsFilterOption
from internal.infrastructure.adapter.mysql.models import (
    Base,
    CarrierEntity,
    FreightEntity,
    ShipmentEntity,
)
from internal.infrastructure.adapter.mysql.stats_repo import MySQLStatsRepository
from internal.infrastructure.adapter.mysql.uuid_helper import uuid_to_bytes

os.environ.setdefault("TESTCONTAINERS_RYUK_DISABLED", "true")


@compiles(LargeBinary, "mysql")
def _mysql_binary16(_type, _compiler, **_kw) -> str:
    return "BINARY(16)"


def _container_runtime_available() -> bool:
    try:
        import docker
    except Exception:
        return False
    try:
        client = docker.from_env()
        client.ping()
        return True
    except Exception:
        return False


pytestmark = pytest.mark.skipif(
    not _container_runtime_available(),
    reason="no reachable container runtime (Docker/Podman) for testcontainers",
)


C1 = uuid4()
C2 = uuid4()
T0 = datetime(2026, 1, 1, 0, 0, 0, tzinfo=UTC)


def _shipment(order_uid: UUID, carrier: UUID, status: str, *, pickup: datetime | None, delivery: datetime | None):
    return ShipmentEntity(
        uid=uuid_to_bytes(uuid4()),
        order_uid=uuid_to_bytes(order_uid),
        carrier_uid=uuid_to_bytes(carrier),
        origin_address="origin",
        destination_address="dest",
        destination_city="Seoul",
        destination_state="Seoul",
        destination_country_code="KR",
        destination_postal_code="00000",
        status=status,
        tracking_number=None,
        weight_kg=1.0,
        planned_pickup_at=pickup,
        planned_delivery_at=delivery,
        actual_pickup_at=pickup,
        actual_delivery_at=delivery,
        created_at=T0,
        updated_at=T0,
    )


def _freight(shipment: ShipmentEntity, carrier: UUID, total_cost: float) -> FreightEntity:
    return FreightEntity(
        uid=uuid_to_bytes(uuid4()),
        shipment_uid=shipment.uid,
        carrier_uid=uuid_to_bytes(carrier),
        base_cost=total_cost,
        weight_surcharge=0.0,
        distance_surcharge=0.0,
        discount=0.0,
        total_cost=total_cost,
        currency="KRW",
        status=FreightStatus.CONFIRMED,
        invoiced_at=None,
        paid_at=None,
        created_at=T0,
        updated_at=T0,
    )


async def _seed(factory: async_sessionmaker[AsyncSession]) -> None:
    async with factory() as session:
        carriers = [
            CarrierEntity(
                uid=uuid_to_bytes(C1),
                name="Carrier One",
                code="C1",
                base_rate=1.0,
                rate_per_kg=1.0,
                status=CarrierStatus.ACTIVE,
                created_at=T0,
                updated_at=T0,
            ),
            CarrierEntity(
                uid=uuid_to_bytes(C2),
                name="Carrier Two",
                code="C2",
                base_rate=1.0,
                rate_per_kg=1.0,
                status=CarrierStatus.ACTIVE,
                created_at=T0,
                updated_at=T0,
            ),
        ]
        s1 = _shipment(uuid4(), C1, ShipmentStatus.DELIVERED, pickup=T0, delivery=T0 + timedelta(hours=20))
        s2 = _shipment(uuid4(), C1, ShipmentStatus.DELIVERED, pickup=T0, delivery=T0 + timedelta(hours=10))
        s3 = _shipment(uuid4(), C1, ShipmentStatus.IN_TRANSIT, pickup=T0, delivery=None)
        s4 = _shipment(uuid4(), C1, ShipmentStatus.FAILED, pickup=None, delivery=None)
        s5 = _shipment(uuid4(), C2, ShipmentStatus.DELIVERED, pickup=T0, delivery=T0 + timedelta(hours=5))

        session.add_all([*carriers, s1, s2, s3, s4, s5])
        await session.flush()

        session.add_all(
            [
                _freight(s1, C1, 100.0),
                _freight(s1, C1, 50.0),
                _freight(s2, C1, 30.0),
                _freight(s3, C1, 40.0),
                _freight(s5, C2, 200.0),
            ]
        )
        await session.commit()


@pytest.fixture(scope="module")
def mysql_url() -> str:
    from testcontainers.mysql import MySqlContainer

    with MySqlContainer("mysql:8.0") as mysql:
        _, _, tail = mysql.get_connection_url().partition("://")
        yield f"mysql+aiomysql://{tail}"


@pytest_asyncio.fixture
async def stats_repo(mysql_url: str) -> AsyncIterator[MySQLStatsRepository]:
    engine = create_async_engine(mysql_url)
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    factory = async_sessionmaker(bind=engine, expire_on_commit=False)
    await _seed(factory)
    try:
        yield MySQLStatsRepository(factory)
    finally:
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.drop_all)
        await engine.dispose()


async def test_dashboard_totals(stats_repo: MySQLStatsRepository) -> None:
    dash = await stats_repo.get_dashboard(StatsFilterOption())

    assert dash.total_shipments == 5
    assert dash.total_freight_cost == pytest.approx(420.0)
    assert dash.avg_delivery_hours == pytest.approx((20 + 10 + 5) / 3, abs=1e-3)

    status_counts = {sc.status: sc.count for sc in dash.status_counts}
    assert status_counts[ShipmentStatus.DELIVERED] == 3
    assert status_counts[ShipmentStatus.IN_TRANSIT] == 1
    assert status_counts[ShipmentStatus.FAILED] == 1


async def test_per_carrier_avg_is_not_biased_by_freight_fanout(stats_repo: MySQLStatsRepository) -> None:
    dash = await stats_repo.get_dashboard(StatsFilterOption())
    perf = {p.carrier_name: p for p in dash.carrier_performances}

    c1 = perf["Carrier One"]
    assert c1.total_shipments == 4
    assert c1.delivered_count == 2
    assert c1.failed_count == 1
    assert c1.avg_delivery_hours == pytest.approx(15.0)
    assert c1.total_freight_cost == pytest.approx(220.0)

    c2 = perf["Carrier Two"]
    assert c2.total_shipments == 1
    assert c2.delivered_count == 1
    assert c2.avg_delivery_hours == pytest.approx(5.0)
    assert c2.total_freight_cost == pytest.approx(200.0)


async def test_carrier_filter_scopes_dashboard(stats_repo: MySQLStatsRepository) -> None:
    dash = await stats_repo.get_dashboard(StatsFilterOption(carrier_uid=C2))

    assert dash.total_shipments == 1
    assert [p.carrier_name for p in dash.carrier_performances] == ["Carrier Two"]
    assert dash.avg_delivery_hours == pytest.approx(5.0)
