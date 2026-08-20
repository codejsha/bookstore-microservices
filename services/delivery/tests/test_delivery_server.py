from datetime import UTC, datetime
from unittest.mock import AsyncMock, MagicMock
from uuid import uuid4

import grpc
import pytest

from generated.application.port.pb.deliverypb.delivery import v1_pb2 as pb
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.option.delivery_option import ShipmentFilterOption
from internal.domain.service.delivery_service import ShipmentService
from internal.infrastructure.adapter.grpc.delivery_server import (
    _PB_TO_STATUS,
    _STATUS_TO_PB,
    DeliveryServiceServicer,
)
from tests.conftest import make_shipment

CALLER_UID = str(uuid4())
CALLER_METADATA = (("x-user-id", CALLER_UID),)


class _Aborted(Exception):
    def __init__(self, code: grpc.StatusCode, detail: str):
        self.code = code
        self.detail = detail


class FakeContext:
    def __init__(self, metadata: tuple[tuple[str, str], ...] = CALLER_METADATA) -> None:
        self.code: grpc.StatusCode | None = None
        self.detail: str | None = None
        self._metadata = metadata

    def invocation_metadata(self) -> tuple[tuple[str, str], ...]:
        return self._metadata

    async def abort(self, code: grpc.StatusCode, detail: str) -> None:
        self.code = code
        self.detail = detail
        raise _Aborted(code, detail)


@pytest.fixture
def service() -> MagicMock:
    return MagicMock(spec=ShipmentService)


@pytest.fixture
def servicer(service: MagicMock) -> DeliveryServiceServicer:
    return DeliveryServiceServicer(service=service)


class TestTrackShipment:
    async def test_maps_aggregate_to_proto(self, servicer: DeliveryServiceServicer, service: MagicMock) -> None:
        now = datetime(2026, 6, 28, 12, 0, 0, tzinfo=UTC)
        shipment = make_shipment(status=ShipmentStatus.IN_TRANSIT, now=now)
        shipment.tracking_number = "DLVABCDEF0123456789"
        shipment.actual_delivery_at = now
        service.get_shipment = AsyncMock(return_value=shipment)

        response = await servicer.TrackShipment(pb.TrackShipmentRequest(uid=str(shipment.uid)), FakeContext())

        s = response.shipment
        assert s.uid == str(shipment.uid)
        assert s.order_uid == str(shipment.order_uid)
        assert s.status == pb.SHIPMENT_STATUS_IN_TRANSIT
        assert s.tracking_number == "DLVABCDEF0123456789"
        assert s.destination_city == shipment.destination_city
        assert s.destination_state == shipment.destination_state
        assert s.destination_country_code == shipment.destination_country_code
        assert s.destination_postal_code == shipment.destination_postal_code
        assert s.created_at == now.isoformat()
        assert s.updated_at == now.isoformat()
        assert s.actual_delivery_at == now.isoformat()
        assert s.planned_delivery_at == ""
        assert s.carrier_uid == ""
        service.get_shipment.assert_awaited_once_with(shipment.uid)

    async def test_aborts_not_found(self, servicer: DeliveryServiceServicer, service: MagicMock) -> None:
        service.get_shipment = AsyncMock(return_value=None)
        context = FakeContext()
        with pytest.raises(_Aborted):
            await servicer.TrackShipment(pb.TrackShipmentRequest(uid=str(uuid4())), context)
        assert context.code == grpc.StatusCode.NOT_FOUND

    async def test_aborts_invalid_uid(self, servicer: DeliveryServiceServicer, service: MagicMock) -> None:
        context = FakeContext()
        with pytest.raises(_Aborted):
            await servicer.TrackShipment(pb.TrackShipmentRequest(uid="not-a-uuid"), context)
        assert context.code == grpc.StatusCode.INVALID_ARGUMENT

    async def test_unauthenticated_without_caller_metadata(
        self, servicer: DeliveryServiceServicer, service: MagicMock
    ) -> None:
        service.get_shipment = AsyncMock(return_value=make_shipment())
        context = FakeContext(metadata=())
        with pytest.raises(_Aborted):
            await servicer.TrackShipment(pb.TrackShipmentRequest(uid=str(uuid4())), context)
        assert context.code == grpc.StatusCode.UNAUTHENTICATED
        service.get_shipment.assert_not_awaited()

    async def test_unauthenticated_with_non_uuid_caller(
        self, servicer: DeliveryServiceServicer, service: MagicMock
    ) -> None:
        service.get_shipment = AsyncMock(return_value=make_shipment())
        context = FakeContext(metadata=(("x-user-id", "not-a-uuid"),))
        with pytest.raises(_Aborted):
            await servicer.TrackShipment(pb.TrackShipmentRequest(uid=str(uuid4())), context)
        assert context.code == grpc.StatusCode.UNAUTHENTICATED


class TestListShipments:
    async def test_maps_list_and_total(self, servicer: DeliveryServiceServicer, service: MagicMock) -> None:
        shipments = [make_shipment(), make_shipment()]
        service.list_shipments = AsyncMock(return_value=(shipments, 2))

        response = await servicer.ListShipments(
            pb.ListShipmentsRequest(order_uid=str(uuid4()), page_size=20), FakeContext()
        )

        assert response.total_size == 2
        assert [s.uid for s in response.shipments] == [str(s.uid) for s in shipments]

    async def test_passes_order_uid_and_status_filters(
        self, servicer: DeliveryServiceServicer, service: MagicMock
    ) -> None:
        service.list_shipments = AsyncMock(return_value=([], 0))
        order_uid = uuid4()
        request = pb.ListShipmentsRequest(
            order_uid=str(order_uid),
            status=pb.SHIPMENT_STATUS_DISPATCHED,
            page_size=20,
        )

        await servicer.ListShipments(request, FakeContext())

        option = service.list_shipments.await_args.args[0]
        assert isinstance(option, ShipmentFilterOption)
        assert option.order_uid == order_uid
        assert option.status == ShipmentStatus.DISPATCHED

    async def test_status_unspecified_means_no_status_filter(
        self, servicer: DeliveryServiceServicer, service: MagicMock
    ) -> None:
        service.list_shipments = AsyncMock(return_value=([], 0))
        order_uid = uuid4()
        await servicer.ListShipments(pb.ListShipmentsRequest(order_uid=str(order_uid), page_size=20), FakeContext())

        option = service.list_shipments.await_args.args[0]
        assert option.status is None
        assert option.order_uid == order_uid

    async def test_next_page_token_set_when_more_results(
        self, servicer: DeliveryServiceServicer, service: MagicMock
    ) -> None:
        service.list_shipments = AsyncMock(return_value=([make_shipment()], 50))
        response = await servicer.ListShipments(
            pb.ListShipmentsRequest(order_uid=str(uuid4()), page_size=20), FakeContext()
        )
        assert response.next_page_token == "1"

    async def test_next_page_token_empty_on_last_page(
        self, servicer: DeliveryServiceServicer, service: MagicMock
    ) -> None:
        service.list_shipments = AsyncMock(return_value=([make_shipment()], 10))
        response = await servicer.ListShipments(
            pb.ListShipmentsRequest(order_uid=str(uuid4()), page_size=20), FakeContext()
        )
        assert response.next_page_token == ""

    async def test_unauthenticated_without_caller_metadata(
        self, servicer: DeliveryServiceServicer, service: MagicMock
    ) -> None:
        service.list_shipments = AsyncMock(return_value=([make_shipment()], 1))
        context = FakeContext(metadata=())
        with pytest.raises(_Aborted):
            await servicer.ListShipments(pb.ListShipmentsRequest(order_uid=str(uuid4())), context)
        assert context.code == grpc.StatusCode.UNAUTHENTICATED
        service.list_shipments.assert_not_awaited()

    async def test_permission_denied_without_order_scope(
        self, servicer: DeliveryServiceServicer, service: MagicMock
    ) -> None:
        service.list_shipments = AsyncMock(return_value=([make_shipment()], 1))
        context = FakeContext()
        with pytest.raises(_Aborted):
            await servicer.ListShipments(pb.ListShipmentsRequest(status=pb.SHIPMENT_STATUS_DELIVERED), context)
        assert context.code == grpc.StatusCode.PERMISSION_DENIED
        service.list_shipments.assert_not_awaited()

    async def test_list_is_scoped_to_requested_order(
        self, servicer: DeliveryServiceServicer, service: MagicMock
    ) -> None:
        own_order = uuid4()
        own_shipment = make_shipment(order_uid=own_order)
        service.list_shipments = AsyncMock(return_value=([own_shipment], 1))

        response = await servicer.ListShipments(
            pb.ListShipmentsRequest(order_uid=str(own_order), page_size=20), FakeContext()
        )

        option = service.list_shipments.await_args.args[0]
        assert option.order_uid == own_order
        assert [s.order_uid for s in response.shipments] == [str(own_order)]


class TestStatusEnumMapping:
    @pytest.mark.parametrize(
        ("domain_status", "pb_status"),
        [
            (ShipmentStatus.PLANNED, pb.SHIPMENT_STATUS_PLANNED),
            (ShipmentStatus.DISPATCHED, pb.SHIPMENT_STATUS_DISPATCHED),
            (ShipmentStatus.PICKED_UP, pb.SHIPMENT_STATUS_PICKED_UP),
            (ShipmentStatus.IN_TRANSIT, pb.SHIPMENT_STATUS_IN_TRANSIT),
            (ShipmentStatus.OUT_FOR_DELIVERY, pb.SHIPMENT_STATUS_OUT_FOR_DELIVERY),
            (ShipmentStatus.DELIVERED, pb.SHIPMENT_STATUS_DELIVERED),
            (ShipmentStatus.FAILED, pb.SHIPMENT_STATUS_FAILED),
            (ShipmentStatus.CANCELLED, pb.SHIPMENT_STATUS_CANCELLED),
        ],
    )
    def test_domain_to_pb_and_back(self, domain_status: ShipmentStatus, pb_status: int) -> None:
        assert _STATUS_TO_PB[domain_status] == pb_status
        assert _PB_TO_STATUS[pb_status] == domain_status

    def test_every_domain_status_is_mapped(self) -> None:
        assert set(_STATUS_TO_PB.keys()) == set(ShipmentStatus)

    def test_unspecified_has_no_domain_mapping(self) -> None:
        assert pb.SHIPMENT_STATUS_UNSPECIFIED not in _PB_TO_STATUS
