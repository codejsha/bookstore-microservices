from collections.abc import Awaitable, Callable
from typing import Any
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import BaseModel, ValidationError

from generated.application.port.model.stock_adjust_request import StockAdjustRequest
from generated.application.port.model.stock_find_all_response import StockFindAllResponse
from generated.application.port.model.stock_find_response import StockFindResponse
from generated.application.port.model.stock_history_find_all_response import StockHistoryFindAllResponse
from generated.application.port.model.stock_history_find_response import StockHistoryFindResponse
from generated.application.port.model.stock_receive_request import StockReceiveRequest
from generated.application.port.model.stock_release_request import StockReleaseRequest
from generated.application.port.model.stock_reserve_request import StockReserveRequest
from generated.application.port.model.stock_warehouse_item import StockWarehouseItem
from internal.domain.aggregate.stock_aggregate import StockAggregate, StockHistoryEntry
from internal.domain.model.command.inventory_command import (
    StockAdjustCommand,
    StockReceiveCommand,
    StockReleaseCommand,
    StockReserveCommand,
)
from internal.domain.model.error import InventoryError
from internal.domain.model.option.inventory_option import (
    DEFAULT_INVENTORY_SORT,
    StockHistoryQueryOption,
    StockQueryOption,
)
from internal.domain.service.inventory_service import StockService
from internal.infrastructure.adapter.restcontroller.error_mapping import invalid_command_error, inventory_error
from internal.infrastructure.support.auth import require_principal, require_staff


def create_stock_router(service: StockService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/stocks",
        tags=["stocks"],
        dependencies=[Depends(require_principal)],
    )

    @router.get("", response_model=StockFindAllResponse)
    async def list_stocks(
        edition_uid: UUID | None = Query(None),
        warehouse_uid: UUID | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query(DEFAULT_INVENTORY_SORT),
    ) -> StockFindAllResponse:
        option = StockQueryOption(edition_uid=edition_uid, warehouse_uid=warehouse_uid, page=page, size=size, sort=sort)
        items, total = await service.list_stocks(option)
        return StockFindAllResponse(items=[_to_stock_response(s) for s in items], total=total)

    @router.post("/receive", response_model=StockFindResponse, dependencies=[Depends(require_staff)])
    async def receive_stock(request: StockReceiveRequest) -> StockFindResponse:
        return await _mutate(service.receive_stock, StockReceiveCommand, request.model_dump())

    @router.post("/release", response_model=StockFindResponse, dependencies=[Depends(require_staff)])
    async def release_stock(request: StockReleaseRequest) -> StockFindResponse:
        return await _mutate(service.release_stock, StockReleaseCommand, request.model_dump())

    @router.post("/adjust", response_model=StockFindResponse, dependencies=[Depends(require_staff)])
    async def adjust_stock(request: StockAdjustRequest) -> StockFindResponse:
        return await _mutate(service.adjust_stock, StockAdjustCommand, request.model_dump())

    @router.post("/reserve", response_model=StockFindResponse, dependencies=[Depends(require_staff)])
    async def reserve_stock(request: StockReserveRequest) -> StockFindResponse:
        return await _mutate(service.reserve_stock, StockReserveCommand, request.model_dump())

    @router.get("/{edition_uid}", response_model=StockFindResponse)
    async def get_stock(edition_uid: UUID) -> StockFindResponse:
        stock = await service.get_stock(edition_uid)
        if stock is None:
            raise HTTPException(status_code=404, detail="Stock not found")
        return _to_stock_response(stock)

    @router.get("/{edition_uid}/history", response_model=StockHistoryFindAllResponse)
    async def get_stock_history(
        edition_uid: UUID,
        warehouse_uid: UUID | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query(DEFAULT_INVENTORY_SORT),
    ) -> StockHistoryFindAllResponse:
        option = StockHistoryQueryOption(warehouse_uid=warehouse_uid, page=page, size=size, sort=sort)
        try:
            items, total = await service.get_stock_history(edition_uid, option)
        except InventoryError as e:
            raise inventory_error(e) from e
        return StockHistoryFindAllResponse(items=[_to_history_response(h) for h in items], total=total)

    return router


async def _mutate[C: BaseModel](
    operation: Callable[[C], Awaitable[StockAggregate | None]],
    command_type: type[C],
    payload: dict[str, Any],
) -> StockFindResponse:
    try:
        command = command_type(**payload)
    except ValidationError as e:
        raise invalid_command_error(e) from e
    try:
        stock = await operation(command)
    except InventoryError as e:
        raise inventory_error(e) from e
    if stock is None:
        raise HTTPException(status_code=404, detail="Stock not found")
    return _to_stock_response(stock)


def _to_stock_response(stock: StockAggregate) -> StockFindResponse:
    return StockFindResponse(
        uid=str(stock.uid),
        edition_uid=str(stock.edition_uid),
        total_quantity=stock.total_quantity,
        warehouses=[
            StockWarehouseItem(
                warehouse_uid=str(w.warehouse_uid),
                warehouse_name=w.warehouse_name,
                quantity=w.quantity,
            )
            for w in stock.warehouses
        ],
    )


def _to_history_response(entry: StockHistoryEntry) -> StockHistoryFindResponse:
    return StockHistoryFindResponse(
        uid=str(entry.uid),
        change_type=entry.change_type,
        reason=entry.reason,
        change_qty=entry.change_qty,
        created_at=entry.created_at,
    )
