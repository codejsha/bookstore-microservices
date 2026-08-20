from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query

from internal.domain.aggregate.ticket_comment_aggregate import TicketCommentAggregate
from internal.domain.constant.author_role import AuthorRole
from internal.domain.model.command.ticket_command import AddCommentCommand
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.schemas import AddCommentRequest, CommentResponse
from internal.infrastructure.support.auth import (
    Principal,
    is_staff,
    owns,
    principal_uid,
    require_principal,
)


def create_comment_router(service: SupportService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/tickets/{ticket_uid}/comments",
        tags=["comments"],
        dependencies=[Depends(require_principal)],
    )

    async def _authorize(ticket_uid: UUID, principal: Principal) -> bool:
        ticket = await service.get_ticket(ticket_uid)
        if ticket is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        staff = is_staff(principal)
        if not staff and not owns(principal, ticket.customer_uid):
            raise HTTPException(status_code=403, detail="Forbidden")
        return staff

    @router.post("", status_code=201, response_model=CommentResponse)
    async def add_comment(
        ticket_uid: UUID,
        request: AddCommentRequest,
        principal: Principal = Depends(require_principal),
    ) -> CommentResponse:
        staff = await _authorize(ticket_uid, principal)
        if request.internal and not staff:
            raise HTTPException(status_code=403, detail="Internal comments require staff role")
        author_uid = principal_uid(principal)
        if author_uid is None:
            raise HTTPException(status_code=403, detail="Cannot resolve caller identity")
        command = AddCommentCommand(
            ticket_uid=ticket_uid,
            author_uid=author_uid,
            author_role=AuthorRole.AGENT if staff else AuthorRole.CUSTOMER,
            body=request.body,
            internal=request.internal,
        )
        comment = await service.add_comment(command)
        if comment is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        return _to_response(comment, ticket_uid)

    @router.get("", response_model=list[CommentResponse])
    async def list_comments(
        ticket_uid: UUID,
        principal: Principal = Depends(require_principal),
        include_internal: bool = Query(False),
    ) -> list[CommentResponse]:
        staff = await _authorize(ticket_uid, principal)
        effective_internal = include_internal and staff
        comments = await service.list_comments(ticket_uid, include_internal=effective_internal)
        if comments is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        return [_to_response(c, ticket_uid) for c in comments]

    return router


def _to_response(c: TicketCommentAggregate, ticket_uid: UUID) -> CommentResponse:
    return CommentResponse(
        uid=c.uid,
        ticket_uid=ticket_uid,
        author_uid=c.author_uid,
        author_role=c.author_role,
        body=c.body,
        internal=c.internal,
        created_at=c.created_at,
        updated_at=c.updated_at,
    )
