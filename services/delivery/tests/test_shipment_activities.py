import dataclasses
from unittest.mock import AsyncMock, MagicMock
from uuid import UUID, uuid4

import pytest
from sqlalchemy.exc import DataError, IntegrityError
from temporalio.exceptions import ApplicationError

from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.command.delivery_command import CreateShipmentCommand
from internal.domain.model.error import ConflictError
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
    def test_create_shipment_input_fields_are_exact_camelcase(self) -> None:
        assert [f.name for f in dataclasses.fields(CreateShipmentInput)] == [
            "orderUid",
            "originAddress",
            "destinationAddress",
            "destinationCity",
            "destinationState",
            "destinationCountryCode",
            "destinationPostalCode",
        ]

    def test_create_shipment_output_fields_are_exact_camelcase(self) -> None:
        assert [f.name for f in dataclasses.fields(CreateShipmentOutput)] == [
            "shipmentUid",
            "trackingNumber",
            "status",
        ]


class TestCreateShipmentActivity:
    async def test_create_shipment_output_is_camelcase_with_dispatched_status(
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

    async def test_create_shipment_maps_input_to_create_command(
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

    async def test_create_shipment_planned_dispatches_then_sets_tracking_number(
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

    async def test_create_shipment_failed_dispatch_raises_runtime_error(
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
    async def test_create_shipment_already_shipped_order_reuses_existing_without_creating(
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

    async def test_create_shipment_still_planned_resumes_dispatch_and_tracking(
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

    async def test_create_shipment_insert_race_falls_back_to_winner_row(
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

    async def test_create_shipment_insert_races_and_no_row_found_reraises_integrity_error(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        shipment_service.get_shipment_by_order_uid = AsyncMock(side_effect=[None, None])
        shipment_service.create_shipment = AsyncMock(
            side_effect=IntegrityError("INSERT", {}, Exception("some other constraint"))
        )

        with pytest.raises(IntegrityError):
            await activities.create_shipment(_input(order_uid))


class TestCreateShipmentFailFast:
    async def test_create_shipment_malformed_order_uid_raises_non_retryable_error(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        payload = dataclasses.replace(_input(uuid4()), orderUid="not-a-uuid")
        with pytest.raises(ApplicationError) as excinfo:
            await activities.create_shipment(payload)
        assert excinfo.value.non_retryable
        assert excinfo.value.type == "InvalidCommand"
        shipment_service.get_shipment_by_order_uid.assert_not_awaited()

    async def test_create_shipment_over_length_address_raises_non_retryable_error(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        payload = dataclasses.replace(_input(uuid4()), originAddress="x" * 501)
        with pytest.raises(ApplicationError) as excinfo:
            await activities.create_shipment(payload)
        assert excinfo.value.non_retryable
        shipment_service.create_shipment.assert_not_called()

    async def test_create_shipment_db_data_error_raises_non_retryable_error(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        shipment_service.create_shipment = AsyncMock(
            side_effect=DataError("INSERT", {}, Exception("Data too long for column"))
        )
        with pytest.raises(ApplicationError) as excinfo:
            await activities.create_shipment(_input(uuid4()))
        assert excinfo.value.non_retryable
        assert excinfo.value.type == "InvalidCommand"

    async def test_create_shipment_conflict_raised_falls_back_to_existing_row(
        self, activities: ShipmentActivities, shipment_service: MagicMock
    ) -> None:
        order_uid = uuid4()
        winner = make_shipment(order_uid=order_uid, status=ShipmentStatus.DISPATCHED)
        winner.tracking_number = "DLVRACEWINNER0002"
        shipment_service.get_shipment_by_order_uid = AsyncMock(side_effect=[None, winner])
        shipment_service.create_shipment = AsyncMock(
            side_effect=ConflictError(f"Shipment for order {order_uid} already exists")
        )
        shipment_service.dispatch_shipment = AsyncMock()

        result = await activities.create_shipment(_input(order_uid))

        assert result.shipmentUid == str(winner.uid)
        assert result.trackingNumber == "DLVRACEWINNER0002"
        shipment_service.dispatch_shipment.assert_not_awaited()
