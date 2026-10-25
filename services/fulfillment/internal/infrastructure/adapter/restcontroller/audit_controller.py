from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import ValidationError

from generated.application.port.model.audit_create_request import AuditCreateRequest
from generated.application.port.model.audit_find_all_response import AuditFindAllResponse
from generated.application.port.model.audit_find_response import AuditFindResponse
from generated.application.port.model.audit_item_response import AuditItemResponse
from internal.domain.aggregate.audit_aggregate import StockAuditAggregate
from internal.domain.constant.audit_status import AuditStatus
from internal.domain.model.command.inventory_command import StockAuditCreateCommand
from internal.domain.model.error import InventoryError
from internal.domain.model.option.inventory_option import DEFAULT_INVENTORY_SORT, AuditQueryOption
from internal.domain.service.inventory_service import AuditService
from internal.infrastructure.adapter.restcontroller.error_mapping import invalid_command_error, inventory_error
from internal.infrastructure.support.auth import require_principal, require_staff


def create_audit_router(service: AuditService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/audits",
        tags=["audits"],
        dependencies=[Depends(require_principal)],
    )

    @router.get("", response_model=AuditFindAllResponse)
    async def list_audits(
        warehouse_uid: UUID | None = Query(None),
        status: AuditStatus | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query(DEFAULT_INVENTORY_SORT),
    ) -> AuditFindAllResponse:
        option = AuditQueryOption(warehouse_uid=warehouse_uid, status=status, page=page, size=size, sort=sort)
        items, total = await service.list_audits(option)
        return AuditFindAllResponse(items=[_to_audit_response(a) for a in items], total=total)

    @router.post("", response_model=AuditFindResponse, dependencies=[Depends(require_staff)])
    async def create_audit(request: AuditCreateRequest) -> AuditFindResponse:
        try:
            command = StockAuditCreateCommand(**request.model_dump())
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            audit = await service.create_audit(command)
        except InventoryError as e:
            raise inventory_error(e) from e
        return _to_audit_response(audit)

    @router.get("/{uid}", response_model=AuditFindResponse)
    async def get_audit(uid: UUID) -> AuditFindResponse:
        audit = await service.get_audit(uid)
        if audit is None:
            raise HTTPException(status_code=404, detail="Audit not found")
        return _to_audit_response(audit)

    @router.post("/{uid}/complete", response_model=AuditFindResponse, dependencies=[Depends(require_staff)])
    async def complete_audit(uid: UUID) -> AuditFindResponse:
        try:
            audit = await service.complete_audit(uid)
        except InventoryError as e:
            raise inventory_error(e) from e
        if audit is None:
            raise HTTPException(status_code=404, detail="Audit not found")
        return _to_audit_response(audit)

    return router


def _to_audit_response(audit: StockAuditAggregate) -> AuditFindResponse:
    return AuditFindResponse(
        uid=str(audit.uid),
        warehouse_uid=str(audit.warehouse_uid),
        status=audit.status,
        items=[
            AuditItemResponse(
                edition_uid=str(item.edition_uid),
                system_quantity=item.system_quantity,
                actual_quantity=item.actual_quantity,
                difference=item.difference,
            )
            for item in audit.items
        ],
        notes=audit.notes,
        completed_at=audit.completed_at,
        created_at=audit.created_at,
        updated_at=audit.updated_at,
    )
