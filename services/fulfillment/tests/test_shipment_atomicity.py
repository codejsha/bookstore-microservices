from collections.abc import Callable
from uuid import UUID, uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.aggregate.tracking_aggregate import TrackingAggregate
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.service.delivery_service import ShipmentService
from internal.infrastructure.adapter.mysql import shipment_repo as shipment_repo_module
from internal.infrastructure.adapter.mysql.shipment_repo import MySQLShipmentRepository
from internal.infrastructure.adapter.mysql.tracking_repo import MySQLTrackingRepository, to_tracking_entity
from internal.infrastructure.adapter.temporal.shipment_activities import CreateShipmentInput, ShipmentActivities
from tests.conftest import make_shipment


@pytest.fixture
def shipment_repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLShipmentRepository:
    return MySQLShipmentRepository(session_factory)


@pytest.fixture
def tracking_repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLTrackingRepository:
    return MySQLTrackingRepository(session_factory)


@pytest.fixture
def service(shipment_repo: MySQLShipmentRepository, tracking_repo: MySQLTrackingRepository) -> ShipmentService:
    return ShipmentService(shipment_repo=shipment_repo, tracking_repo=tracking_repo)


def _failing_tracking_insert(status: ShipmentStatus) -> Callable[[TrackingAggregate], object]:
    def _build(tracking: TrackingAggregate):
        if tracking.status == status:
            raise RuntimeError("tracking insert failed")
        return to_tracking_entity(tracking)

    return _build


def _input(order_uid) -> CreateShipmentInput:
    return CreateShipmentInput(
        orderUid=str(order_uid),
        originAddress="1 Origin Way, Seoul",
        destinationAddress="42 Destination Blvd",
        destinationCity="Busan",
        destinationState="Busan",
        destinationCountryCode="KR",
        destinationPostalCode="48058",
    )


class TestTransitionAtomicity:
    @pytest.mark.parametrize(
        ("initial_status", "action", "expected_status"),
        [
            (ShipmentStatus.PLANNED, "dispatch_shipment", ShipmentStatus.DISPATCHED),
            (ShipmentStatus.DISPATCHED, "pick_up_shipment", ShipmentStatus.PICKED_UP),
            (ShipmentStatus.IN_TRANSIT, "deliver_shipment", ShipmentStatus.DELIVERED),
            (ShipmentStatus.PLANNED, "cancel_shipment", ShipmentStatus.CANCELLED),
        ],
    )
    async def test_transition_commits_status_and_tracking_together(
        self,
        service: ShipmentService,
        shipment_repo: MySQLShipmentRepository,
        tracking_repo: MySQLTrackingRepository,
        initial_status: ShipmentStatus,
        action: str,
        expected_status: ShipmentStatus,
    ) -> None:
        shipment = make_shipment(status=initial_status)
        await shipment_repo.save(shipment)

        result = await getattr(service, action)(shipment.uid)

        assert result is not None
        assert result.status == expected_status
        reloaded = await shipment_repo.find_by_uid(shipment.uid)
        assert reloaded is not None
        assert reloaded.status == expected_status
        history = await tracking_repo.find_by_shipment_uid(shipment.uid)
        assert [event.status for event in history] == [expected_status]

    @pytest.mark.parametrize(
        ("initial_status", "action", "attempted_status"),
        [
            (ShipmentStatus.PLANNED, "dispatch_shipment", ShipmentStatus.DISPATCHED),
            (ShipmentStatus.DISPATCHED, "pick_up_shipment", ShipmentStatus.PICKED_UP),
            (ShipmentStatus.IN_TRANSIT, "deliver_shipment", ShipmentStatus.DELIVERED),
            (ShipmentStatus.PLANNED, "cancel_shipment", ShipmentStatus.CANCELLED),
        ],
    )
    async def test_transition_tracking_insert_failure_rolls_back_status(
        self,
        service: ShipmentService,
        shipment_repo: MySQLShipmentRepository,
        tracking_repo: MySQLTrackingRepository,
        monkeypatch: pytest.MonkeyPatch,
        initial_status: ShipmentStatus,
        action: str,
        attempted_status: ShipmentStatus,
    ) -> None:
        shipment = make_shipment(status=initial_status)
        await shipment_repo.save(shipment)

        monkeypatch.setattr(shipment_repo_module, "to_tracking_entity", _failing_tracking_insert(attempted_status))
        with pytest.raises(RuntimeError, match="tracking insert failed"):
            await getattr(service, action)(shipment.uid)

        reloaded = await shipment_repo.find_by_uid(shipment.uid)
        assert reloaded is not None
        assert reloaded.status == initial_status
        assert await tracking_repo.find_by_shipment_uid(shipment.uid) == []

    async def test_dispatch_retry_after_tracking_insert_failure_completes(
        self,
        service: ShipmentService,
        shipment_repo: MySQLShipmentRepository,
        tracking_repo: MySQLTrackingRepository,
        monkeypatch: pytest.MonkeyPatch,
    ) -> None:
        shipment = make_shipment(status=ShipmentStatus.PLANNED)
        await shipment_repo.save(shipment)

        monkeypatch.setattr(
            shipment_repo_module, "to_tracking_entity", _failing_tracking_insert(ShipmentStatus.DISPATCHED)
        )
        with pytest.raises(RuntimeError, match="tracking insert failed"):
            await service.dispatch_shipment(shipment.uid, "DLVRETRY00000001")
        monkeypatch.undo()

        result = await service.dispatch_shipment(shipment.uid, "DLVRETRY00000001")

        assert result is not None
        assert result.status == ShipmentStatus.DISPATCHED
        assert result.tracking_number == "DLVRETRY00000001"
        history = await tracking_repo.find_by_shipment_uid(shipment.uid)
        assert [event.status for event in history] == [ShipmentStatus.DISPATCHED]


class TestCreateShipmentActivityAtomicity:
    async def test_create_shipment_dispatch_failure_leaves_shipment_planned_without_tracking(
        self,
        service: ShipmentService,
        shipment_repo: MySQLShipmentRepository,
        tracking_repo: MySQLTrackingRepository,
        monkeypatch: pytest.MonkeyPatch,
    ) -> None:
        activities = ShipmentActivities(shipment_service=service)
        order_uid = uuid4()

        monkeypatch.setattr(
            shipment_repo_module, "to_tracking_entity", _failing_tracking_insert(ShipmentStatus.DISPATCHED)
        )
        with pytest.raises(RuntimeError, match="tracking insert failed"):
            await activities.create_shipment(_input(order_uid))

        persisted = await shipment_repo.find_by_order_uid(order_uid)
        assert persisted is not None
        assert persisted.status == ShipmentStatus.PLANNED
        assert persisted.tracking_number is None
        history = await tracking_repo.find_by_shipment_uid(persisted.uid)
        assert [event.status for event in history] == [ShipmentStatus.PLANNED]

    async def test_create_shipment_retry_after_dispatch_failure_completes(
        self,
        service: ShipmentService,
        shipment_repo: MySQLShipmentRepository,
        tracking_repo: MySQLTrackingRepository,
        monkeypatch: pytest.MonkeyPatch,
    ) -> None:
        activities = ShipmentActivities(shipment_service=service)
        order_uid = uuid4()

        monkeypatch.setattr(
            shipment_repo_module, "to_tracking_entity", _failing_tracking_insert(ShipmentStatus.DISPATCHED)
        )
        with pytest.raises(RuntimeError, match="tracking insert failed"):
            await activities.create_shipment(_input(order_uid))
        monkeypatch.undo()

        result = await activities.create_shipment(_input(order_uid))

        assert result.status == ShipmentStatus.DISPATCHED.value
        assert result.trackingNumber is not None
        assert result.trackingNumber.startswith("DLV")

        persisted = await shipment_repo.find_by_order_uid(order_uid)
        assert persisted is not None
        assert persisted.uid == UUID(result.shipmentUid)
        assert persisted.status == ShipmentStatus.DISPATCHED
        assert persisted.tracking_number == result.trackingNumber
        history = await tracking_repo.find_by_shipment_uid(persisted.uid)
        assert [event.status for event in history] == [ShipmentStatus.PLANNED, ShipmentStatus.DISPATCHED]

    async def test_create_shipment_creates_single_shipment_for_repeated_activity_runs(
        self,
        service: ShipmentService,
        shipment_repo: MySQLShipmentRepository,
        tracking_repo: MySQLTrackingRepository,
    ) -> None:
        activities = ShipmentActivities(shipment_service=service)
        order_uid = uuid4()

        first = await activities.create_shipment(_input(order_uid))
        second = await activities.create_shipment(_input(order_uid))

        assert first.shipmentUid == second.shipmentUid
        assert first.trackingNumber == second.trackingNumber
        persisted = await shipment_repo.find_by_order_uid(order_uid)
        assert persisted is not None
        history = await tracking_repo.find_by_shipment_uid(persisted.uid)
        assert [event.status for event in history] == [ShipmentStatus.PLANNED, ShipmentStatus.DISPATCHED]
