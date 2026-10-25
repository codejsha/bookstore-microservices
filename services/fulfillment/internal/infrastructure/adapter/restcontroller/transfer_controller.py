from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import ValidationError

from generated.application.port.model.transfer_create_request import TransferCreateRequest
from generated.application.port.model.transfer_find_all_response import TransferFindAllResponse
from generated.application.port.model.transfer_find_response import TransferFindResponse
from internal.domain.aggregate.transfer_aggregate import StockTransferAggregate
from internal.domain.constant.transfer_status import TransferStatus
from internal.domain.model.command.inventory_command import StockTransferCommand
from internal.domain.model.error import InventoryError
from internal.domain.model.option.inventory_option import DEFAULT_INVENTORY_SORT, TransferQueryOption
from internal.domain.service.inventory_service import TransferService
from internal.infrastructure.adapter.restcontroller.error_mapping import invalid_command_error, inventory_error
from internal.infrastructure.support.auth import require_principal, require_staff


def create_transfer_router(service: TransferService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/transfers",
        tags=["transfers"],
        dependencies=[Depends(require_principal)],
    )

    @router.get("", response_model=TransferFindAllResponse)
    async def list_transfers(
        edition_uid: UUID | None = Query(None),
        source_warehouse_uid: UUID | None = Query(None),
        target_warehouse_uid: UUID | None = Query(None),
        status: TransferStatus | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query(DEFAULT_INVENTORY_SORT),
    ) -> TransferFindAllResponse:
        option = TransferQueryOption(
            edition_uid=edition_uid,
            source_warehouse_uid=source_warehouse_uid,
            target_warehouse_uid=target_warehouse_uid,
            status=status,
            page=page,
            size=size,
            sort=sort,
        )
        items, total = await service.list_transfers(option)
        return TransferFindAllResponse(items=[_to_transfer_response(t) for t in items], total=total)

    @router.post("", response_model=TransferFindResponse, dependencies=[Depends(require_staff)])
    async def create_transfer(request: TransferCreateRequest) -> TransferFindResponse:
        try:
            command = StockTransferCommand(**request.model_dump())
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            transfer = await service.transfer_stock(command)
        except InventoryError as e:
            raise inventory_error(e) from e
        return _to_transfer_response(transfer)

    @router.get("/{uid}", response_model=TransferFindResponse)
    async def get_transfer(uid: UUID) -> TransferFindResponse:
        transfer = await service.get_transfer(uid)
        if transfer is None:
            raise HTTPException(status_code=404, detail="Transfer not found")
        return _to_transfer_response(transfer)

    @router.post("/{uid}/complete", response_model=TransferFindResponse, dependencies=[Depends(require_staff)])
    async def complete_transfer(uid: UUID) -> TransferFindResponse:
        try:
            transfer = await service.complete_transfer(uid)
        except InventoryError as e:
            raise inventory_error(e) from e
        if transfer is None:
            raise HTTPException(status_code=404, detail="Transfer not found")
        return _to_transfer_response(transfer)

    @router.post("/{uid}/cancel", response_model=TransferFindResponse, dependencies=[Depends(require_staff)])
    async def cancel_transfer(uid: UUID) -> TransferFindResponse:
        try:
            transfer = await service.cancel_transfer(uid)
        except InventoryError as e:
            raise inventory_error(e) from e
        if transfer is None:
            raise HTTPException(status_code=404, detail="Transfer not found")
        return _to_transfer_response(transfer)

    return router


def _to_transfer_response(transfer: StockTransferAggregate) -> TransferFindResponse:
    return TransferFindResponse(
        uid=str(transfer.uid),
        edition_uid=str(transfer.edition_uid),
        source_warehouse_uid=str(transfer.source_warehouse_uid),
        target_warehouse_uid=str(transfer.target_warehouse_uid),
        quantity=transfer.quantity,
        status=transfer.status,
        reason=transfer.reason,
        completed_at=transfer.completed_at,
        created_at=transfer.created_at,
        updated_at=transfer.updated_at,
    )
