from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query

from internal.domain.aggregate.ticket_aggregate import TicketAggregate
from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus
from internal.domain.model.command.ticket_command import (
    CreateTicketCommand,
    UpdateTicketCommand,
    UpdateTicketStatusCommand,
)
from internal.domain.model.option.ticket_option import TicketFilterOption
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.schemas import (
    CreateTicketRequest,
    PaginatedTicketResponse,
    TicketResponse,
    UpdateTicketRequest,
    UpdateTicketStatusRequest,
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

    @router.post("", status_code=201, response_model=TicketResponse)
    async def create_ticket(
        request: CreateTicketRequest,
        principal: Principal = Depends(require_principal),
    ) -> TicketResponse:
        customer_uid = request.customer_uid if is_staff(principal) else _caller_customer_uid(principal)
        command = CreateTicketCommand(
            customer_uid=customer_uid,
            subject=request.subject,
            description=request.description,
            priority=request.priority,
            category_uid=request.category_uid,
        )
        return _to_response(await service.create_ticket(command))

    @router.get("/{uid}", response_model=TicketResponse)
    async def get_ticket(
        uid: UUID,
        principal: Principal = Depends(require_principal),
    ) -> TicketResponse:
        ticket = await service.get_ticket(uid)
        if ticket is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        if not is_staff(principal) and not owns(principal, ticket.customer_uid):
            raise HTTPException(status_code=403, detail="Forbidden")
        return _to_response(ticket)

    @router.get("", response_model=PaginatedTicketResponse)
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
    ) -> PaginatedTicketResponse:
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
        return PaginatedTicketResponse(
            content=[_to_response(t) for t in tickets],
            total=total,
            page=page,
            size=size,
        )

    @router.patch("/{uid}", response_model=TicketResponse)
    async def update_ticket(
        uid: UUID,
        request: UpdateTicketRequest,
        principal: Principal = Depends(require_principal),
    ) -> TicketResponse:
        ticket = await service.get_ticket(uid)
        if ticket is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        staff = is_staff(principal)
        if not staff and not owns(principal, ticket.customer_uid):
            raise HTTPException(status_code=403, detail="Forbidden")
        if request.assignee_uid is not None and not staff:
            raise HTTPException(status_code=403, detail="Assignment requires staff role")
        command = UpdateTicketCommand(
            subject=request.subject,
            description=request.description,
            priority=request.priority,
            category_uid=request.category_uid,
            assignee_uid=request.assignee_uid,
        )
        ticket = await service.update_ticket(uid, command)
        if ticket is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        return _to_response(ticket)

    @router.patch("/{uid}/status", response_model=TicketResponse)
    async def update_status(
        uid: UUID,
        request: UpdateTicketStatusRequest,
        _staff: Principal = Depends(require_staff),
    ) -> TicketResponse:
        ticket = await service.update_ticket_status(uid, UpdateTicketStatusCommand(status=request.status))
        if ticket is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        return _to_response(ticket)

    return router


def _to_response(t: TicketAggregate) -> TicketResponse:
    return TicketResponse(
        uid=t.uid,
        customer_uid=t.customer_uid,
        category_uid=t.category_uid,
        assignee_uid=t.assignee_uid,
        subject=t.subject,
        description=t.description,
        status=t.status,
        priority=t.priority,
        resolved_at=t.resolved_at,
        created_at=t.created_at,
        updated_at=t.updated_at,
    )
