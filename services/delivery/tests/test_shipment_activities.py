import dataclasses
from unittest.mock import AsyncMock, MagicMock
from uuid import UUID, uuid4

import pytest
from sqlalchemy.exc import IntegrityError

from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.command.delivery_command import CreateShipmentCommand
from internal.domain.service.delivery_service import ShipmentService
from internal.infrastructure.adapter.temporal.shipment_activities import (
    CreateShipmentInput,
    CreateShipmentOutput,
    ShipmentActivities,
)
from tests.conftest import make_shipment


def _input(order_uid: UUID) -> CreateShipmentInput:
    return CreateShipmentInput(
        orderUid=str(order_uid),
        originAddress="1 Origin Way, Seoul",
        destinationAddress="42 Destination Blvd",
        destinationCity="Busan",
        destinationState="Busan",
        destinationCountryCode="KR",
        destinationPostalCode="48058",
    )


@pytest.fixture
def shipment_service() -> MagicMock:
    service = MagicMock(spec=ShipmentService)
    service.set_tracking_number = AsyncMock(side_effect=lambda s: s)
    service.get_shipment_by_order_uid = AsyncMock(return_value=None)
    return service


@pytest.fixture
def activities(shipment_service: MagicMock) -> ShipmentActivities:
    return ShipmentActivities(shipment_service=shipment_service)


class TestCreateShipmentContract:
    def test_input_dataclass_fields_are_exact_camelcase(self) -> None:
        assert [f.name for f in dataclasses.fields(CreateShipmentInput)] == [
            "orderUid",
            "originAddress",
            "destinationAddress",
            "destinationCity",
            "destinationState",
            "destinationCountryCode",
            "destinationPostalCode",
        ]

    def test_output_dataclass_fields_are_exact_camelcase(self) -> None:
        assert [f.name for f in dataclasses.fields(CreateShipmentOutput)] == [
            "shipmentUid",
            "trackingNumber",
            "status",
        ]


class TestCreateShipmentActivity:
    async def test_returns_camelcase_output_with_dispatched_status(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        shipment_uid = uuid4()
        planned = make_shipment(uid=shipment_uid, order_uid=order_uid, status=ShipmentStatus.PLANNED)
        dispatched = make_shipment(uid=shipment_uid, order_uid=order_uid, status=ShipmentStatus.DISPATCHED)
        shipment_service.create_shipment = AsyncMock(return_value=planned)
        shipment_service.dispatch_shipment = AsyncMock(return_value=dispatched)

        result = await activities.create_shipment(_input(order_uid))

        assert isinstance(result, CreateShipmentOutput)
        serialized = dataclasses.asdict(result)
        assert set(serialized.keys()) == {"shipmentUid", "trackingNumber", "status"}
        assert serialized["shipmentUid"] == str(shipment_uid)
        assert serialized["status"] == "DISPATCHED"
        assert serialized["trackingNumber"]
        assert serialized["trackingNumber"].startswith("DLV")

    async def test_creates_shipment_with_mapped_command(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        planned = make_shipment(order_uid=order_uid, status=ShipmentStatus.PLANNED)
        dispatched = make_shipment(uid=planned.uid, order_uid=order_uid, status=ShipmentStatus.DISPATCHED)
        shipment_service.create_shipment = AsyncMock(return_value=planned)
        shipment_service.dispatch_shipment = AsyncMock(return_value=dispatched)

        payload = _input(order_uid)
        await activities.create_shipment(payload)

        shipment_service.create_shipment.assert_awaited_once()
        command = shipment_service.create_shipment.await_args.args[0]
        assert isinstance(command, CreateShipmentCommand)
        assert command.order_uid == UUID(payload.orderUid)
        assert command.origin_address == payload.originAddress
        assert command.destination_address == payload.destinationAddress
        assert command.destination_city == payload.destinationCity
        assert command.destination_state == payload.destinationState
        assert command.destination_country_code == payload.destinationCountryCode
        assert command.destination_postal_code == payload.destinationPostalCode

    async def test_dispatches_then_sets_tracking_number(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        planned = make_shipment(order_uid=order_uid, status=ShipmentStatus.PLANNED)
        dispatched = make_shipment(uid=planned.uid, order_uid=order_uid, status=ShipmentStatus.DISPATCHED)
        shipment_service.create_shipment = AsyncMock(return_value=planned)
        shipment_service.dispatch_shipment = AsyncMock(return_value=dispatched)

        result = await activities.create_shipment(_input(order_uid))

        shipment_service.dispatch_shipment.assert_awaited_once_with(planned.uid)
        shipment_service.set_tracking_number.assert_awaited_once()
        saved_arg = shipment_service.set_tracking_number.await_args.args[0]
        assert saved_arg is dispatched
        assert saved_arg.tracking_number is not None
        assert result.trackingNumber == saved_arg.tracking_number

    async def test_raises_when_dispatch_returns_none(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        planned = make_shipment(order_uid=order_uid, status=ShipmentStatus.PLANNED)
        shipment_service.create_shipment = AsyncMock(return_value=planned)
        shipment_service.dispatch_shipment = AsyncMock(return_value=None)

        with pytest.raises(RuntimeError, match="failed to dispatch"):
            await activities.create_shipment(_input(order_uid))
        shipment_service.set_tracking_number.assert_not_awaited()


class TestCreateShipmentIdempotency:
    async def test_returns_existing_shipment_without_creating(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        existing = make_shipment(order_uid=order_uid, status=ShipmentStatus.DISPATCHED)
        existing.tracking_number = "DLVEXISTINGTRACK01"
        shipment_service.get_shipment_by_order_uid = AsyncMock(return_value=existing)
        shipment_service.create_shipment = AsyncMock()
        shipment_service.dispatch_shipment = AsyncMock()

        result = await activities.create_shipment(_input(order_uid))

        assert result.shipmentUid == str(existing.uid)
        assert result.trackingNumber == "DLVEXISTINGTRACK01"
        assert result.status == "DISPATCHED"
        shipment_service.create_shipment.assert_not_awaited()
        shipment_service.dispatch_shipment.assert_not_awaited()
        shipment_service.set_tracking_number.assert_not_awaited()

    async def test_resumes_partially_processed_shipment(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        planned = make_shipment(order_uid=order_uid, status=ShipmentStatus.PLANNED)
        dispatched = make_shipment(uid=planned.uid, order_uid=order_uid, status=ShipmentStatus.DISPATCHED)
        shipment_service.get_shipment_by_order_uid = AsyncMock(return_value=planned)
        shipment_service.create_shipment = AsyncMock()
        shipment_service.dispatch_shipment = AsyncMock(return_value=dispatched)

        result = await activities.create_shipment(_input(order_uid))

        shipment_service.create_shipment.assert_not_awaited()
        shipment_service.dispatch_shipment.assert_awaited_once_with(planned.uid)
        shipment_service.set_tracking_number.assert_awaited_once()
        assert result.status == "DISPATCHED"
        assert result.trackingNumber and result.trackingNumber.startswith("DLV")

    async def test_insert_race_falls_back_to_existing_row(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        winner = make_shipment(order_uid=order_uid, status=ShipmentStatus.DISPATCHED)
        winner.tracking_number = "DLVRACEWINNER0001"
        shipment_service.get_shipment_by_order_uid = AsyncMock(side_effect=[None, winner])
        shipment_service.create_shipment = AsyncMock(
            side_effect=IntegrityError("INSERT", {}, Exception("duplicate order_uid"))
        )
        shipment_service.dispatch_shipment = AsyncMock()

        result = await activities.create_shipment(_input(order_uid))

        assert result.shipmentUid == str(winner.uid)
        assert result.trackingNumber == "DLVRACEWINNER0001"
        assert shipment_service.get_shipment_by_order_uid.await_count == 2
        shipment_service.dispatch_shipment.assert_not_awaited()
        shipment_service.set_tracking_number.assert_not_awaited()

    async def test_insert_race_reraises_when_no_row_found(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        shipment_service.get_shipment_by_order_uid = AsyncMock(side_effect=[None, None])
        shipment_service.create_shipment = AsyncMock(
            side_effect=IntegrityError("INSERT", {}, Exception("some other constraint"))
        )

        with pytest.raises(IntegrityError):
            await activities.create_shipment(_input(order_uid))
