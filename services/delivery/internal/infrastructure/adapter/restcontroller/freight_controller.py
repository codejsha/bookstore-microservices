from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import ValidationError

from generated.application.port.model.freight_create_request import FreightCreateRequest
from generated.application.port.model.freight_find_all_response import FreightFindAllResponse
from generated.application.port.model.freight_find_response import FreightFindResponse
from generated.application.port.model.freight_item import FreightItem
from generated.application.port.model.freight_update_response import FreightUpdateResponse
from generated.application.port.model.freight_update_status_request import FreightUpdateStatusRequest
from internal.domain.constant.freight_status import FreightStatus
from internal.domain.model.command.delivery_command import CreateFreightCommand, UpdateFreightStatusCommand
from internal.domain.model.option.delivery_option import FreightFilterOption
from internal.domain.service.delivery_service import FreightService
from internal.infrastructure.adapter.restcontroller.error_mapping import invalid_command_error
from internal.infrastructure.support.auth import require_staff


def create_freight_router(service: FreightService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/freights",
        tags=["freights"],
        dependencies=[Depends(require_staff)],
    )

    @router.post("", status_code=201, response_model=FreightItem)
    async def create_freight(request: FreightCreateRequest) -> FreightItem:
        try:
            command = CreateFreightCommand(**request.model_dump(exclude_none=True))
        except ValidationError as e:
            raise invalid_command_error(e) from e
        return _to_freight_item(await service.create_freight(command))

    @router.get("/{uid}", response_model=FreightFindResponse)
    async def get_freight(uid: UUID) -> FreightFindResponse:
        agg = await service.get_freight(uid)
        if agg is None:
            raise HTTPException(status_code=404, detail="Freight not found")
        return FreightFindResponse(**_to_freight_item(agg).model_dump())

    @router.get("", response_model=FreightFindAllResponse)
    async def list_freights(
        shipment_uid: UUID | None = Query(None),
        carrier_uid: UUID | None = Query(None),
        status: FreightStatus | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query("created_at:desc"),
    ) -> FreightFindAllResponse:
        option = FreightFilterOption(
            shipment_uid=shipment_uid,
            carrier_uid=carrier_uid,
            status=status,
            page=page,
            size=size,
            sort=sort,
        )
        items, total = await service.list_freights(option)
        return FreightFindAllResponse(
            items=[_to_freight_item(a) for a in items],
            total=total,
        )

    @router.get("/shipment/{shipment_uid}", response_model=FreightFindResponse)
    async def get_freight_by_shipment(shipment_uid: UUID) -> FreightFindResponse:
        agg = await service.find_by_shipment(shipment_uid)
        if agg is None:
            raise HTTPException(status_code=404, detail="Freight not found")
        return FreightFindResponse(**_to_freight_item(agg).model_dump())

    @router.patch("/{uid}/status", response_model=FreightUpdateResponse)
    async def update_freight_status(uid: UUID, request: FreightUpdateStatusRequest) -> FreightUpdateResponse:
        command = UpdateFreightStatusCommand(status=request.status)
        agg = await service.update_status(uid, command)
        if agg is None:
            raise HTTPException(status_code=404, detail="Freight not found")
        return FreightUpdateResponse(**_to_freight_item(agg).model_dump())

    return router


def _to_freight_item(agg) -> FreightItem:
    return FreightItem(
        uid=str(agg.uid),
        shipment_uid=str(agg.shipment_uid),
        carrier_uid=str(agg.carrier_uid),
        base_cost=agg.base_cost,
        weight_surcharge=agg.weight_surcharge,
        distance_surcharge=agg.distance_surcharge,
        discount=agg.discount,
        total_cost=agg.total_cost,
        currency=agg.currency,
        status=agg.status,
        invoiced_at=agg.invoiced_at,
        paid_at=agg.paid_at,
        created_at=agg.created_at,
        updated_at=agg.updated_at,
    )
