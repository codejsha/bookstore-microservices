from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.command.delivery_command import (
    AddTrackingCommand,
    AssignCarrierCommand,
    CreateShipmentCommand,
)
from internal.domain.model.error import ConflictError
from internal.domain.model.option.delivery_option import ShipmentFilterOption
from internal.domain.service.delivery_service import ShipmentService
from tests.conftest import make_shipment, make_tracking


@pytest.fixture
def service(shipment_repo: MagicMock, tracking_repo: MagicMock) -> ShipmentService:
    return ShipmentService(shipment_repo=shipment_repo, tracking_repo=tracking_repo)


def _create_command() -> CreateShipmentCommand:
    return CreateShipmentCommand(
        order_uid=uuid4(),
        carrier_uid=None,
        origin_address="origin",
        destination_address="dest",
        destination_city="Seoul",
        destination_state="KR",
        destination_country_code="KR",
        destination_postal_code="00000",
        weight_kg=2.5,
    )


class TestCreateShipment:
    async def test_create_shipment_unshipped_order_persists_planned_shipment(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        result = await service.create_shipment(_create_command())
        assert result.status == ShipmentStatus.PLANNED
        assert result.tracking_number is None
        assert result.actual_pickup_at is None
        assert result.actual_delivery_at is None
        shipment_repo.create_with_initial_tracking.assert_called_once()
        shipment_repo.save.assert_not_called()

    async def test_create_shipment_already_shipped_order_raises_conflict_error(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        command = _create_command()
        shipment_repo.find_by_order_uid.return_value = make_shipment(order_uid=command.order_uid)
        with pytest.raises(ConflictError):
            await service.create_shipment(command)
        shipment_repo.create_with_initial_tracking.assert_not_called()

    async def test_create_shipment_persisted_creates_initial_tracking_event_with_it(
        self, service: ShipmentService, shipment_repo: MagicMock, tracking_repo: MagicMock
    ) -> None:
        await service.create_shipment(_create_command())
        shipment_repo.create_with_initial_tracking.assert_called_once()
        shipment_arg, tracking_arg = shipment_repo.create_with_initial_tracking.call_args[0]
        assert tracking_arg.status == ShipmentStatus.PLANNED
        assert tracking_arg.shipment_uid == shipment_arg.uid
        tracking_repo.save.assert_not_called()


class TestGetShipment:
    async def test_get_shipment_roundtrips(self, service: ShipmentService, shipment_repo: MagicMock) -> None:
        shipment = make_shipment()
        shipment_repo.find_by_uid.return_value = shipment
        assert await service.get_shipment(shipment.uid) is shipment

    async def test_get_shipment_missing_is_none(self, service: ShipmentService, shipment_repo: MagicMock) -> None:
        shipment_repo.find_by_uid.return_value = None
        assert await service.get_shipment(uuid4()) is None


class TestListShipments:
    async def test_list_shipments_with_option_delegates_to_repository(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        option = ShipmentFilterOption()
        shipments = [make_shipment()]
        shipment_repo.find_all.return_value = (shipments, 1)
        result, total = await service.list_shipments(option)
        assert result == shipments
        assert total == 1
        shipment_repo.find_all.assert_called_once_with(option)


class TestAssignCarrier:
    async def test_assign_carrier_missing_is_none(self, service: ShipmentService, shipment_repo: MagicMock) -> None:
        shipment_repo.find_by_uid.return_value = None
        result = await service.assign_carrier(uuid4(), AssignCarrierCommand(carrier_uid=uuid4()))
        assert result is None
        shipment_repo.save.assert_not_called()

    async def test_assign_carrier_with_tracking_number_updates_carrier_and_tracking_number(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        shipment = make_shipment()
        shipment_repo.find_by_uid.return_value = shipment
        carrier_uid = uuid4()
        result = await service.assign_carrier(
            shipment.uid, AssignCarrierCommand(carrier_uid=carrier_uid, tracking_number="TRK-1")
        )
        assert result is not None
        assert result.carrier_uid == carrier_uid
        assert result.tracking_number == "TRK-1"

    async def test_assign_carrier_without_tracking_number_keeps_existing_one(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        shipment = make_shipment()
        shipment.tracking_number = "DLVEXISTING01"
        shipment_repo.find_by_uid.return_value = shipment
        result = await service.assign_carrier(shipment.uid, AssignCarrierCommand(carrier_uid=uuid4()))
        assert result is not None
        assert result.tracking_number == "DLVEXISTING01"


class TestDispatchShipment:
    async def test_dispatch_shipment_planned_status_transitions_to_dispatched(
        self, service: ShipmentService, shipment_repo: MagicMock, tracking_repo: MagicMock
    ) -> None:
        shipment = make_shipment(status=ShipmentStatus.PLANNED)
        shipment_repo.find_by_uid.return_value = shipment
        result = await service.dispatch_shipment(shipment.uid)
        assert result is not None
        assert result.status == ShipmentStatus.DISPATCHED
        tracking_repo.save.assert_called_once()

    async def test_dispatch_shipment_missing_is_none(self, service: ShipmentService, shipment_repo: MagicMock) -> None:
        shipment_repo.find_by_uid.return_value = None
        assert await service.dispatch_shipment(uuid4()) is None

    @pytest.mark.parametrize(
        "status",
        [
            ShipmentStatus.DISPATCHED,
            ShipmentStatus.PICKED_UP,
            ShipmentStatus.DELIVERED,
            ShipmentStatus.CANCELLED,
        ],
    )
    async def test_dispatch_shipment_invalid_status_raises_value_error(
        self, service: ShipmentService, shipment_repo: MagicMock, status: ShipmentStatus
    ) -> None:
        shipment_repo.find_by_uid.return_value = make_shipment(status=status)
        with pytest.raises(ValueError):
            await service.dispatch_shipment(uuid4())


class TestPickUpShipment:
    async def test_pick_up_shipment_dispatched_status_transitions_to_picked_up(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        shipment = make_shipment(status=ShipmentStatus.DISPATCHED)
        shipment_repo.find_by_uid.return_value = shipment
        result = await service.pick_up_shipment(shipment.uid)
        assert result is not None
        assert result.status == ShipmentStatus.PICKED_UP
        assert result.actual_pickup_at is not None

    async def test_pick_up_shipment_invalid_status_raises_value_error(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        shipment_repo.find_by_uid.return_value = make_shipment(status=ShipmentStatus.PLANNED)
        with pytest.raises(ValueError):
            await service.pick_up_shipment(uuid4())


class TestDeliverShipment:
    @pytest.mark.parametrize("status", [ShipmentStatus.IN_TRANSIT, ShipmentStatus.OUT_FOR_DELIVERY])
    async def test_deliver_shipment_picked_up_status_transitions_to_delivered(
        self, service: ShipmentService, shipment_repo: MagicMock, status: ShipmentStatus
    ) -> None:
        shipment = make_shipment(status=status)
        shipment_repo.find_by_uid.return_value = shipment
        result = await service.deliver_shipment(shipment.uid)
        assert result is not None
        assert result.status == ShipmentStatus.DELIVERED
        assert result.actual_delivery_at is not None

    async def test_deliver_shipment_invalid_status_raises_value_error(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        shipment_repo.find_by_uid.return_value = make_shipment(status=ShipmentStatus.PLANNED)
        with pytest.raises(ValueError):
            await service.deliver_shipment(uuid4())


class TestCancelShipment:
    @pytest.mark.parametrize("status", [ShipmentStatus.PLANNED, ShipmentStatus.DISPATCHED])
    async def test_cancel_shipment_cancellable_status_transitions_to_cancelled(
        self, service: ShipmentService, shipment_repo: MagicMock, status: ShipmentStatus
    ) -> None:
        shipment_repo.find_by_uid.return_value = make_shipment(status=status)
        result = await service.cancel_shipment(uuid4())
        assert result is not None
        assert result.status == ShipmentStatus.CANCELLED

    async def test_cancel_shipment_invalid_status_raises_value_error(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        shipment_repo.find_by_uid.return_value = make_shipment(status=ShipmentStatus.DELIVERED)
        with pytest.raises(ValueError):
            await service.cancel_shipment(uuid4())


class TestAddTracking:
    async def test_add_tracking_found_creates_tracking_event(
        self, service: ShipmentService, shipment_repo: MagicMock, tracking_repo: MagicMock
    ) -> None:
        shipment = make_shipment()
        shipment_repo.find_by_uid.return_value = shipment
        command = AddTrackingCommand(
            shipment_uid=shipment.uid,
            status=ShipmentStatus.IN_TRANSIT,
            location="Hub",
            description="Arrived at hub",
        )
        result = await service.add_tracking(command)
        assert result.status == ShipmentStatus.IN_TRANSIT
        assert result.location == "Hub"
        assert result.description == "Arrived at hub"
        tracking_repo.save.assert_called_once()

    async def test_add_tracking_missing_raises_value_error(
        self, service: ShipmentService, shipment_repo: MagicMock
    ) -> None:
        shipment_repo.find_by_uid.return_value = None
        command = AddTrackingCommand(
            shipment_uid=uuid4(),
            status=ShipmentStatus.IN_TRANSIT,
            location="Hub",
            description="Arrived at hub",
        )
        with pytest.raises(ValueError):
            await service.add_tracking(command)


class TestGetTrackingHistory:
    async def test_get_tracking_history_passes_through_repository_events(
        self, service: ShipmentService, tracking_repo: MagicMock
    ) -> None:
        shipment_uid = uuid4()
        events = [make_tracking(shipment_uid=shipment_uid)]
        tracking_repo.find_by_shipment_uid.return_value = events
        assert await service.get_tracking_history(shipment_uid) == events
        tracking_repo.find_by_shipment_uid.assert_called_once_with(shipment_uid)
