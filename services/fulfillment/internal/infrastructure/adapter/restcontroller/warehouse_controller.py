from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query, Response
from pydantic import ValidationError

from generated.application.port.model.warehouse_create_request import WarehouseCreateRequest
from generated.application.port.model.warehouse_find_all_response import WarehouseFindAllResponse
from generated.application.port.model.warehouse_find_response import WarehouseFindResponse
from generated.application.port.model.warehouse_update_request import WarehouseUpdateRequest
from generated.application.port.model.warehouse_update_response import WarehouseUpdateResponse
from internal.domain.aggregate.warehouse_aggregate import WarehouseAggregate
from internal.domain.model.command.inventory_command import WarehouseCreateCommand, WarehouseUpdateCommand
from internal.domain.model.error import InventoryError
from internal.domain.model.option.inventory_option import DEFAULT_INVENTORY_SORT, WarehouseQueryOption
from internal.domain.service.inventory_service import WarehouseService
from internal.infrastructure.adapter.restcontroller.error_mapping import invalid_command_error, inventory_error
from internal.infrastructure.support.auth import require_manager, require_principal


def create_warehouse_router(service: WarehouseService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/warehouses",
        tags=["warehouses"],
        dependencies=[Depends(require_principal)],
    )

    @router.get("", response_model=WarehouseFindAllResponse)
    async def list_warehouses(
        name: str | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query(DEFAULT_INVENTORY_SORT),
    ) -> WarehouseFindAllResponse:
        option = WarehouseQueryOption(name=name, page=page, size=size, sort=sort)
        items, total = await service.list_warehouses(option)
        return WarehouseFindAllResponse(items=[_to_warehouse_response(w) for w in items], total=total)

    @router.post("", status_code=201, dependencies=[Depends(require_manager)])
    async def register_warehouse(request: WarehouseCreateRequest) -> Response:
        try:
            command = WarehouseCreateCommand(**request.model_dump())
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            warehouse = await service.register_warehouse(command)
        except InventoryError as e:
            raise inventory_error(e) from e
        return Response(status_code=201, headers={"Location": f"/api/v1/warehouses/{warehouse.uid}"})

    @router.get("/{uid}", response_model=WarehouseFindResponse)
    async def get_warehouse(uid: UUID) -> WarehouseFindResponse:
        warehouse = await service.get_warehouse(uid)
        if warehouse is None:
            raise HTTPException(status_code=404, detail="Warehouse not found")
        return _to_warehouse_response(warehouse)

    @router.put("/{uid}", response_model=WarehouseUpdateResponse, dependencies=[Depends(require_manager)])
    async def update_warehouse(uid: UUID, request: WarehouseUpdateRequest) -> WarehouseUpdateResponse:
        try:
            command = WarehouseUpdateCommand(**request.model_dump(exclude_none=True))
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            warehouse = await service.update_warehouse(uid, command)
        except InventoryError as e:
            raise inventory_error(e) from e
        if warehouse is None:
            raise HTTPException(status_code=404, detail="Warehouse not found")
        return WarehouseUpdateResponse(**_to_warehouse_response(warehouse).model_dump())

    return router


def _to_warehouse_response(warehouse: WarehouseAggregate) -> WarehouseFindResponse:
    return WarehouseFindResponse(
        uid=str(warehouse.uid),
        name=warehouse.name,
        address=warehouse.address,
        capacity=warehouse.capacity,
        created_at=warehouse.created_at,
        updated_at=warehouse.updated_at,
    )
