from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import ValidationError
from sqlalchemy.exc import IntegrityError

from generated.application.port.model.carrier_create_request import CarrierCreateRequest
from generated.application.port.model.carrier_find_all_response import CarrierFindAllResponse
from generated.application.port.model.carrier_find_response import CarrierFindResponse
from generated.application.port.model.carrier_item import CarrierItem
from generated.application.port.model.carrier_update_request import CarrierUpdateRequest
from generated.application.port.model.carrier_update_response import CarrierUpdateResponse
from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.model.command.delivery_command import CreateCarrierCommand, UpdateCarrierCommand
from internal.domain.model.error import ConflictError
from internal.domain.model.option.delivery_option import CarrierFilterOption
from internal.domain.service.delivery_service import CarrierService
from internal.infrastructure.adapter.restcontroller.error_mapping import (
    conflict_error,
    duplicate_key_error,
    invalid_command_error,
)
from internal.infrastructure.support.auth import require_staff


def create_carrier_router(service: CarrierService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/carriers",
        tags=["carriers"],
        dependencies=[Depends(require_staff)],
    )

    @router.post("", status_code=201, response_model=CarrierItem)
    async def create_carrier(request: CarrierCreateRequest) -> CarrierItem:
        try:
            command = CreateCarrierCommand(**request.model_dump())
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            carrier = await service.create_carrier(command)
        except ConflictError as e:
            raise conflict_error(e) from e
        except IntegrityError as e:
            raise duplicate_key_error(f"Carrier code {command.code} is already in use") from e
        return _to_carrier_item(carrier)

    @router.get("/{uid}", response_model=CarrierFindResponse)
    async def get_carrier(uid: UUID) -> CarrierFindResponse:
        agg = await service.get_carrier(uid)
        if agg is None:
            raise HTTPException(status_code=404, detail="Carrier not found")
        return CarrierFindResponse(**_to_carrier_item(agg).model_dump())

    @router.get("", response_model=CarrierFindAllResponse)
    async def list_carriers(
        name: str | None = Query(None),
        status: CarrierStatus | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query("name:asc"),
    ) -> CarrierFindAllResponse:
        option = CarrierFilterOption(name=name, status=status, page=page, size=size, sort=sort)
        items, total = await service.list_carriers(option)
        return CarrierFindAllResponse(
            items=[_to_carrier_item(a) for a in items],
            total=total,
        )

    @router.put("/{uid}", response_model=CarrierUpdateResponse)
    async def update_carrier(uid: UUID, request: CarrierUpdateRequest) -> CarrierUpdateResponse:
        try:
            command = UpdateCarrierCommand(**request.model_dump(exclude_unset=True))
        except ValidationError as e:
            raise invalid_command_error(e) from e
        agg = await service.update_carrier(uid, command)
        if agg is None:
            raise HTTPException(status_code=404, detail="Carrier not found")
        return CarrierUpdateResponse(**_to_carrier_item(agg).model_dump())

    return router


def _to_carrier_item(agg) -> CarrierItem:
    return CarrierItem(
        uid=str(agg.uid),
        name=agg.name,
        code=agg.code,
        contact_name=agg.contact_name,
        contact_phone=agg.contact_phone,
        contact_email=agg.contact_email,
        base_rate=agg.base_rate,
        rate_per_kg=agg.rate_per_kg,
        status=agg.status,
        created_at=agg.created_at,
        updated_at=agg.updated_at,
    )
