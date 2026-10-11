from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import ValidationError

from generated.application.port.model.ticket_create_request import TicketCreateRequest
from generated.application.port.model.ticket_find_all_response import TicketFindAllResponse
from generated.application.port.model.ticket_find_response import TicketFindResponse
from generated.application.port.model.ticket_item import TicketItem
from generated.application.port.model.ticket_update_request import TicketUpdateRequest
from generated.application.port.model.ticket_update_response import TicketUpdateResponse
from generated.application.port.model.ticket_update_status_request import TicketUpdateStatusRequest
from internal.domain.aggregate.ticket_aggregate import TicketAggregate
from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus
from internal.domain.error import UnknownReferenceError
from internal.domain.model.command.ticket_command import (
    CreateTicketCommand,
    UpdateTicketCommand,
    UpdateTicketStatusCommand,
)
from internal.domain.model.option.ticket_option import TicketFilterOption
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.error_mapping import (
    invalid_command_error,
    unknown_reference_error,
)
from internal.infrastructure.support.auth import (
    Principal,
    is_staff,
    owns,
    principal_uid,
    require_principal,
    require_staff,
)


def _caller_customer_uid(principal: Principal) -> UUID:
    uid = principal_uid(principal)
    if uid is None:
        raise HTTPException(status_code=403, detail="Cannot resolve caller identity")
    return uid


def create_ticket_router(service: SupportService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/tickets",
        tags=["tickets"],
        dependencies=[Depends(require_principal)],
    )

    @router.post("", status_code=201, response_model=TicketItem)
    async def create_ticket(
        request: TicketCreateRequest,
        principal: Principal = Depends(require_principal),
    ) -> TicketItem:
        fields = request.model_dump(exclude_none=True)
        if not is_staff(principal):
            fields["customer_uid"] = _caller_customer_uid(principal)
        try:
            command = CreateTicketCommand(**fields)
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            ticket = await service.create_ticket(command)
        except UnknownReferenceError as exc:
            raise unknown_reference_error(exc) from exc
        return _to_ticket_item(ticket)

    @router.get("/{uid}", response_model=TicketFindResponse)
    async def get_ticket(
        uid: UUID,
        principal: Principal = Depends(require_principal),
    ) -> TicketFindResponse:
        ticket = await service.get_ticket(uid)
        if ticket is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        if not is_staff(principal) and not owns(principal, ticket.customer_uid):
            raise HTTPException(status_code=403, detail="Forbidden")
        return TicketFindResponse(**_to_ticket_item(ticket).model_dump())

    @router.get("", response_model=TicketFindAllResponse)
    async def list_tickets(
        principal: Principal = Depends(require_principal),
        customer_uid: UUID | None = Query(None),
        assignee_uid: UUID | None = Query(None),
        category_uid: UUID | None = Query(None),
        status: TicketStatus | None = Query(None),
        priority: TicketPriority | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query("created_at:desc"),
    ) -> TicketFindAllResponse:
        if not is_staff(principal):
            customer_uid = _caller_customer_uid(principal)
        option = TicketFilterOption(
            customer_uid=customer_uid,
            assignee_uid=assignee_uid,
            category_uid=category_uid,
            status=status,
            priority=priority,
            page=page,
            size=size,
            sort=sort,
        )
        tickets, total = await service.list_tickets(option)
        return TicketFindAllResponse(
            items=[_to_ticket_item(t) for t in tickets],
            total=total,
        )

    @router.patch("/{uid}", response_model=TicketUpdateResponse)
    async def update_ticket(
        uid: UUID,
        request: TicketUpdateRequest,
        principal: Principal = Depends(require_principal),
    ) -> TicketUpdateResponse:
        ticket = await service.get_ticket(uid)
        if ticket is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        staff = is_staff(principal)
        if not staff and not owns(principal, ticket.customer_uid):
            raise HTTPException(status_code=403, detail="Forbidden")
        if request.assignee_uid is not None and not staff:
            raise HTTPException(status_code=403, detail="Assignment requires staff role")
        try:
            command = UpdateTicketCommand(**request.model_dump(exclude_unset=True))
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            ticket = await service.update_ticket(uid, command)
        except UnknownReferenceError as exc:
            raise unknown_reference_error(exc) from exc
        if ticket is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        return TicketUpdateResponse(**_to_ticket_item(ticket).model_dump())

    @router.patch("/{uid}/status", response_model=TicketUpdateResponse)
    async def update_status(
        uid: UUID,
        request: TicketUpdateStatusRequest,
        _staff: Principal = Depends(require_staff),
    ) -> TicketUpdateResponse:
        ticket = await service.update_ticket_status(uid, UpdateTicketStatusCommand(status=request.status))
        if ticket is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        return TicketUpdateResponse(**_to_ticket_item(ticket).model_dump())

    return router


def _to_ticket_item(t: TicketAggregate) -> TicketItem:
    return TicketItem(
        uid=str(t.uid),
        customer_uid=str(t.customer_uid),
        category_uid=str(t.category_uid) if t.category_uid else None,
        assignee_uid=str(t.assignee_uid) if t.assignee_uid else None,
        subject=t.subject,
        description=t.description,
        status=t.status,
        priority=t.priority,
        resolved_at=t.resolved_at,
        created_at=t.created_at,
        updated_at=t.updated_at,
    )
