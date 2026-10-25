from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import ValidationError

from generated.application.port.model.closing_create_request import ClosingCreateRequest
from generated.application.port.model.closing_find_all_response import ClosingFindAllResponse
from generated.application.port.model.closing_find_response import ClosingFindResponse
from generated.application.port.model.closing_item_response import ClosingItemResponse
from internal.domain.aggregate.closing_aggregate import MonthlyClosingAggregate
from internal.domain.constant.closing_status import ClosingStatus
from internal.domain.model.command.inventory_command import MonthlyClosingCommand
from internal.domain.model.error import InventoryError
from internal.domain.model.option.inventory_option import DEFAULT_INVENTORY_SORT, ClosingQueryOption
from internal.domain.service.inventory_service import ClosingService
from internal.infrastructure.adapter.restcontroller.error_mapping import invalid_command_error, inventory_error
from internal.infrastructure.support.auth import require_manager, require_principal


def create_closing_router(service: ClosingService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/closings",
        tags=["closings"],
        dependencies=[Depends(require_principal)],
    )

    @router.get("", response_model=ClosingFindAllResponse)
    async def list_closings(
        warehouse_uid: UUID | None = Query(None),
        year: int | None = Query(None),
        month: int | None = Query(None),
        status: ClosingStatus | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query(DEFAULT_INVENTORY_SORT),
    ) -> ClosingFindAllResponse:
        option = ClosingQueryOption(
            warehouse_uid=warehouse_uid, year=year, month=month, status=status, page=page, size=size, sort=sort
        )
        items, total = await service.list_closings(option)
        return ClosingFindAllResponse(items=[_to_closing_response(c) for c in items], total=total)

    @router.post("", response_model=ClosingFindResponse, dependencies=[Depends(require_manager)])
    async def create_closing(request: ClosingCreateRequest) -> ClosingFindResponse:
        try:
            command = MonthlyClosingCommand(**request.model_dump())
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            closing = await service.create_monthly_closing(command)
        except InventoryError as e:
            raise inventory_error(e) from e
        return _to_closing_response(closing)

    @router.get("/{uid}", response_model=ClosingFindResponse)
    async def get_closing(uid: UUID) -> ClosingFindResponse:
        closing = await service.get_closing(uid)
        if closing is None:
            raise HTTPException(status_code=404, detail="Closing not found")
        return _to_closing_response(closing)

    return router


def _to_closing_response(closing: MonthlyClosingAggregate) -> ClosingFindResponse:
    return ClosingFindResponse(
        uid=str(closing.uid),
        warehouse_uid=str(closing.warehouse_uid),
        year=closing.year,
        month=closing.month,
        status=closing.status,
        items=[
            ClosingItemResponse(
                edition_uid=str(item.edition_uid),
                opening_quantity=item.opening_quantity,
                inbound_quantity=item.inbound_quantity,
                outbound_quantity=item.outbound_quantity,
                adjust_quantity=item.adjust_quantity,
                closing_quantity=item.closing_quantity,
            )
            for item in closing.items
        ],
        closed_at=closing.closed_at,
        created_at=closing.created_at,
        updated_at=closing.updated_at,
    )
