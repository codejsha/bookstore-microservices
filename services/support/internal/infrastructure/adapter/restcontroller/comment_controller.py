from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query, Response

from generated.application.port.model.comment_create_request import CommentCreateRequest
from generated.application.port.model.comment_find_all_response import CommentFindAllResponse
from generated.application.port.model.comment_item import CommentItem
from internal.domain.aggregate.ticket_comment_aggregate import TicketCommentAggregate
from internal.domain.constant.author_role import AuthorRole
from internal.domain.model.command.ticket_command import AddCommentCommand
from internal.domain.service.support_service import SupportService
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

    @router.post("", status_code=201)
    async def add_comment(
        ticket_uid: UUID,
        request: CommentCreateRequest,
        principal: Principal = Depends(require_principal),
    ) -> Response:
        staff = await _authorize(ticket_uid, principal)
        internal = bool(request.internal)
        if internal and not staff:
            raise HTTPException(status_code=403, detail="Internal comments require staff role")
        author_uid = principal_uid(principal)
        if author_uid is None:
            raise HTTPException(status_code=403, detail="Cannot resolve caller identity")
        command = AddCommentCommand(
            ticket_uid=ticket_uid,
            author_uid=author_uid,
            author_role=AuthorRole.AGENT if staff else AuthorRole.CUSTOMER,
            body=request.body,
            internal=internal,
        )
        comment = await service.add_comment(command)
        if comment is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        return Response(
            status_code=201,
            headers={"Location": f"/api/v1/tickets/{ticket_uid}/comments/{comment.uid}"},
        )

    @router.get("", response_model=CommentFindAllResponse)
    async def list_comments(
        ticket_uid: UUID,
        principal: Principal = Depends(require_principal),
        include_internal: bool = Query(False),
    ) -> CommentFindAllResponse:
        staff = await _authorize(ticket_uid, principal)
        effective_internal = include_internal and staff
        comments = await service.list_comments(ticket_uid, include_internal=effective_internal)
        if comments is None:
            raise HTTPException(status_code=404, detail="Ticket not found")
        return CommentFindAllResponse(items=[_to_comment_item(c, ticket_uid) for c in comments])

    return router


def _to_comment_item(c: TicketCommentAggregate, ticket_uid: UUID) -> CommentItem:
    return CommentItem(
        uid=str(c.uid),
        ticket_uid=str(ticket_uid),
        author_uid=str(c.author_uid),
        author_role=c.author_role,
        body=c.body,
        internal=c.internal,
        created_at=c.created_at,
        updated_at=c.updated_at,
    )
