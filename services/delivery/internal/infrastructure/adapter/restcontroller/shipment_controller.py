from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query

from generated.application.port.model.shipment_assign_carrier_request import ShipmentAssignCarrierRequest
from generated.application.port.model.shipment_create_request import ShipmentCreateRequest
from generated.application.port.model.shipment_find_all_response import ShipmentFindAllResponse
from generated.application.port.model.shipment_find_response import ShipmentFindResponse
from generated.application.port.model.shipment_item import ShipmentItem
from generated.application.port.model.shipment_update_response import ShipmentUpdateResponse
from generated.application.port.model.tracking_create_request import TrackingCreateRequest
from generated.application.port.model.tracking_find_all_response import TrackingFindAllResponse
from generated.application.port.model.tracking_item import TrackingItem
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.command.delivery_command import (
    AddTrackingCommand,
    AssignCarrierCommand,
    CreateShipmentCommand,
)
from internal.domain.model.option.delivery_option import ShipmentFilterOption
from internal.domain.service.delivery_service import ShipmentService
from internal.infrastructure.support.auth import require_principal, require_staff


def create_shipment_router(service: ShipmentService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/shipments",
        tags=["shipments"],
        dependencies=[Depends(require_principal)],
    )

    @router.post("", status_code=201, response_model=ShipmentItem, dependencies=[Depends(require_staff)])
    async def create_shipment(request: ShipmentCreateRequest) -> ShipmentItem:
        command = CreateShipmentCommand(**request.model_dump())
        return _to_shipment_item(await service.create_shipment(command))

    @router.get("/{uid}", response_model=ShipmentFindResponse)
    async def get_shipment(uid: UUID) -> ShipmentFindResponse:
        agg = await service.get_shipment(uid)
        if agg is None:
            raise HTTPException(status_code=404, detail="Shipment not found")
        return ShipmentFindResponse(**_to_shipment_item(agg).model_dump())

    @router.get("", response_model=ShipmentFindAllResponse, dependencies=[Depends(require_staff)])
    async def list_shipments(
        order_uid: UUID | None = Query(None),
        carrier_uid: UUID | None = Query(None),
        status: ShipmentStatus | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query("created_at:desc"),
    ) -> ShipmentFindAllResponse:
        option = ShipmentFilterOption(
            order_uid=order_uid,
            carrier_uid=carrier_uid,
            status=status,
            page=page,
            size=size,
            sort=sort,
        )
        items, total = await service.list_shipments(option)
        return ShipmentFindAllResponse(
            items=[_to_shipment_item(a) for a in items],
            total=total,
        )

    @router.patch("/{uid}/carrier", response_model=ShipmentUpdateResponse, dependencies=[Depends(require_staff)])
    async def assign_carrier(uid: UUID, request: ShipmentAssignCarrierRequest) -> ShipmentUpdateResponse:
        command = AssignCarrierCommand(carrier_uid=request.carrier_uid, tracking_number=request.tracking_number)
        agg = await service.assign_carrier(uid, command)
        if agg is None:
            raise HTTPException(status_code=404, detail="Shipment not found")
        return ShipmentUpdateResponse(**_to_shipment_item(agg).model_dump())

    @router.patch("/{uid}/dispatch", response_model=ShipmentUpdateResponse, dependencies=[Depends(require_staff)])
    async def dispatch_shipment(uid: UUID) -> ShipmentUpdateResponse:
        try:
            agg = await service.dispatch_shipment(uid)
        except ValueError as e:
            raise HTTPException(status_code=400, detail=str(e))
        if agg is None:
            raise HTTPException(status_code=404, detail="Shipment not found")
        return ShipmentUpdateResponse(**_to_shipment_item(agg).model_dump())

    @router.patch("/{uid}/pickup", response_model=ShipmentUpdateResponse, dependencies=[Depends(require_staff)])
    async def pick_up_shipment(uid: UUID) -> ShipmentUpdateResponse:
        try:
            agg = await service.pick_up_shipment(uid)
        except ValueError as e:
            raise HTTPException(status_code=400, detail=str(e))
        if agg is None:
            raise HTTPException(status_code=404, detail="Shipment not found")
        return ShipmentUpdateResponse(**_to_shipment_item(agg).model_dump())

    @router.patch("/{uid}/deliver", response_model=ShipmentUpdateResponse, dependencies=[Depends(require_staff)])
    async def deliver_shipment(uid: UUID) -> ShipmentUpdateResponse:
        try:
            agg = await service.deliver_shipment(uid)
        except ValueError as e:
            raise HTTPException(status_code=400, detail=str(e))
        if agg is None:
            raise HTTPException(status_code=404, detail="Shipment not found")
        return ShipmentUpdateResponse(**_to_shipment_item(agg).model_dump())

    @router.patch("/{uid}/cancel", response_model=ShipmentUpdateResponse, dependencies=[Depends(require_staff)])
    async def cancel_shipment(uid: UUID) -> ShipmentUpdateResponse:
        try:
            agg = await service.cancel_shipment(uid)
        except ValueError as e:
            raise HTTPException(status_code=400, detail=str(e))
        if agg is None:
            raise HTTPException(status_code=404, detail="Shipment not found")
        return ShipmentUpdateResponse(**_to_shipment_item(agg).model_dump())

    @router.post(
        "/{shipment_uid}/tracking",
        status_code=201,
        response_model=TrackingItem,
        dependencies=[Depends(require_staff)],
    )
    async def add_tracking(shipment_uid: UUID, request: TrackingCreateRequest) -> TrackingItem:
        command = AddTrackingCommand(
            shipment_uid=shipment_uid,
            status=request.status,
            location=request.location,
            description=request.description,
        )
        try:
            agg = await service.add_tracking(command)
        except ValueError as e:
            raise HTTPException(status_code=404, detail=str(e))
        return _to_tracking_item(agg)

    @router.get("/{shipment_uid}/tracking", response_model=TrackingFindAllResponse)
    async def get_tracking_history(shipment_uid: UUID) -> TrackingFindAllResponse:
        items = await service.get_tracking_history(shipment_uid)
        return TrackingFindAllResponse(
            items=[_to_tracking_item(a) for a in items],
        )

    return router


def _to_shipment_item(agg) -> ShipmentItem:
    return ShipmentItem(
        uid=str(agg.uid),
        order_uid=str(agg.order_uid),
        carrier_uid=str(agg.carrier_uid) if agg.carrier_uid else None,
        origin_address=agg.origin_address,
        destination_address=agg.destination_address,
        destination_city=agg.destination_city,
        destination_state=agg.destination_state,
        destination_country_code=agg.destination_country_code,
        destination_postal_code=agg.destination_postal_code,
        status=agg.status,
        tracking_number=agg.tracking_number,
        weight_kg=agg.weight_kg,
        planned_pickup_at=agg.planned_pickup_at,
        planned_delivery_at=agg.planned_delivery_at,
        actual_pickup_at=agg.actual_pickup_at,
        actual_delivery_at=agg.actual_delivery_at,
        created_at=agg.created_at,
        updated_at=agg.updated_at,
    )


def _to_tracking_item(agg) -> TrackingItem:
    return TrackingItem(
        uid=str(agg.uid),
        shipment_uid=str(agg.shipment_uid),
        status=agg.status,
        location=agg.location,
        description=agg.description,
        occurred_at=agg.occurred_at,
        created_at=agg.created_at,
    )
