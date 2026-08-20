from typing import NoReturn
from uuid import UUID

import grpc

from generated.application.port.pb.deliverypb.delivery import v1_pb2 as pb
from generated.application.port.pb.deliverypb.delivery import v1_pb2_grpc as pb_grpc
from internal.domain.aggregate.shipment_aggregate import ShipmentAggregate
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.option.delivery_option import ShipmentFilterOption
from internal.domain.service.delivery_service import ShipmentService
from internal.infrastructure.support.grpc_auth import caller_user_uid

_STATUS_TO_PB = {
    ShipmentStatus.PLANNED: pb.SHIPMENT_STATUS_PLANNED,
    ShipmentStatus.DISPATCHED: pb.SHIPMENT_STATUS_DISPATCHED,
    ShipmentStatus.PICKED_UP: pb.SHIPMENT_STATUS_PICKED_UP,
    ShipmentStatus.IN_TRANSIT: pb.SHIPMENT_STATUS_IN_TRANSIT,
    ShipmentStatus.OUT_FOR_DELIVERY: pb.SHIPMENT_STATUS_OUT_FOR_DELIVERY,
    ShipmentStatus.DELIVERED: pb.SHIPMENT_STATUS_DELIVERED,
    ShipmentStatus.FAILED: pb.SHIPMENT_STATUS_FAILED,
    ShipmentStatus.CANCELLED: pb.SHIPMENT_STATUS_CANCELLED,
}
_PB_TO_STATUS = {v: k for k, v in _STATUS_TO_PB.items()}


async def _abort(context: grpc.aio.ServicerContext, code: grpc.StatusCode, detail: str) -> NoReturn:
    await context.abort(code, detail)
    raise RuntimeError("unreachable")  # pragma: no cover


async def _parse_uid(value: str, context: grpc.aio.ServicerContext) -> UUID:
    try:
        return UUID(value)
    except ValueError, TypeError:
        await _abort(context, grpc.StatusCode.INVALID_ARGUMENT, f"Invalid UUID: {value}")


async def _require_caller(context: grpc.aio.ServicerContext) -> UUID:
    caller = caller_user_uid(context.invocation_metadata())
    if caller is None:
        await _abort(context, grpc.StatusCode.UNAUTHENTICATED, "Missing or invalid caller identity")
    return caller


def _optional_uid(value: str) -> UUID | None:
    if not value:
        return None
    try:
        return UUID(value)
    except ValueError, TypeError:
        return None


def _shipment_to_pb(s: ShipmentAggregate) -> pb.Shipment:
    return pb.Shipment(
        uid=str(s.uid),
        order_uid=str(s.order_uid),
        carrier_uid=str(s.carrier_uid) if s.carrier_uid else "",
        status=_STATUS_TO_PB[s.status],
        tracking_number=s.tracking_number or "",
        destination_city=s.destination_city,
        destination_state=s.destination_state,
        destination_country_code=s.destination_country_code,
        destination_postal_code=s.destination_postal_code,
        planned_delivery_at=s.planned_delivery_at.isoformat() if s.planned_delivery_at else "",
        actual_delivery_at=s.actual_delivery_at.isoformat() if s.actual_delivery_at else "",
        created_at=s.created_at.isoformat(),
        updated_at=s.updated_at.isoformat() if s.updated_at else "",
    )


class DeliveryServiceServicer(pb_grpc.DeliveryServiceServicer):
    def __init__(self, service: ShipmentService):
        self._service = service

    async def TrackShipment(
        self, request: pb.TrackShipmentRequest, context: grpc.aio.ServicerContext
    ) -> pb.TrackShipmentResponse:
        await _require_caller(context)
        shipment = await self._service.get_shipment(await _parse_uid(request.uid, context))
        if shipment is None:
            await _abort(context, grpc.StatusCode.NOT_FOUND, "Shipment not found")
        return pb.TrackShipmentResponse(shipment=_shipment_to_pb(shipment))

    async def ListShipments(
        self, request: pb.ListShipmentsRequest, context: grpc.aio.ServicerContext
    ) -> pb.ListShipmentsResponse:
        await _require_caller(context)
        order_uid = _optional_uid(request.order_uid)
        if order_uid is None:
            await _abort(
                context,
                grpc.StatusCode.PERMISSION_DENIED,
                "order_uid is required; unscoped shipment listing is not permitted",
            )
        page = _parse_page_token(request.page_token)
        size = request.page_size or 20
        option = ShipmentFilterOption(
            order_uid=order_uid,
            status=_PB_TO_STATUS.get(request.status) if request.status else None,
            page=page,
            size=size,
            sort=request.order_by or "created_at:desc",
        )
        shipments, total = await self._service.list_shipments(option)
        next_page_token = str(page + 1) if (page + 1) * size < total else ""
        return pb.ListShipmentsResponse(
            shipments=[_shipment_to_pb(s) for s in shipments],
            next_page_token=next_page_token,
            total_size=total,
        )


def _parse_page_token(token: str) -> int:
    if not token:
        return 0
    try:
        page = int(token)
        return page if page >= 0 else 0
    except ValueError:
        return 0
