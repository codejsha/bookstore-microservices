from uuid import UUID

from fastapi import APIRouter, Depends, Query
from pydantic import ValidationError

from generated.application.port.model.balance_find_all_response import BalanceFindAllResponse
from generated.application.port.model.balance_find_response import BalanceFindResponse
from internal.domain.aggregate.balance_aggregate import StockBalanceEntry
from internal.domain.model.error import InventoryError
from internal.domain.model.option.inventory_option import DEFAULT_INVENTORY_SORT, BalanceQueryOption
from internal.domain.service.inventory_service import BalanceService
from internal.infrastructure.adapter.restcontroller.error_mapping import invalid_command_error, inventory_error
from internal.infrastructure.support.auth import require_principal


def create_balance_router(service: BalanceService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/balance",
        tags=["balance"],
        dependencies=[Depends(require_principal)],
    )

    @router.get("", response_model=BalanceFindAllResponse)
    async def get_stock_balance(
        edition_uid: UUID | None = Query(None),
        warehouse_uid: UUID | None = Query(None),
        year_month: str | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query(DEFAULT_INVENTORY_SORT),
    ) -> BalanceFindAllResponse:
        try:
            option = BalanceQueryOption(
                edition_uid=edition_uid,
                warehouse_uid=warehouse_uid,
                year_month=year_month,
                page=page,
                size=size,
                sort=sort,
            )
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            items, total = await service.get_stock_balance(option)
        except InventoryError as e:
            raise inventory_error(e) from e
        return BalanceFindAllResponse(items=[_to_balance_response(b) for b in items], total=total)

    return router


def _to_balance_response(entry: StockBalanceEntry) -> BalanceFindResponse:
    return BalanceFindResponse(
        edition_uid=str(entry.edition_uid),
        warehouse_uid=str(entry.warehouse_uid),
        warehouse_name=entry.warehouse_name,
        inbound_quantity=entry.inbound_quantity,
        outbound_quantity=entry.outbound_quantity,
        adjust_quantity=entry.adjust_quantity,
        current_quantity=entry.current_quantity,
    )
