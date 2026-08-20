from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query, Response

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus
from internal.domain.model.command.notification_command import SendNotificationCommand
from internal.domain.model.option.notification_option import NotificationFilterOption
from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.restcontroller.schemas import (
    NotificationFindAllResponse,
    NotificationResponse,
    SendNotificationRequest,
)
from internal.infrastructure.support.auth import (
    Principal,
    is_staff,
    require_principal,
    require_staff,
)


def create_notification_router(service: NotificationService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/notifications",
        tags=["notifications"],
        dependencies=[Depends(require_principal)],
    )

    @router.post("", status_code=201)
    async def send_notification(
        request: SendNotificationRequest,
        _staff: Principal = Depends(require_staff),
    ) -> Response:
        command = SendNotificationCommand(
            user_uid=request.user_uid,
            notification_type=request.notification_type,
            channel=request.channel,
            title=request.title,
            content=request.content,
        )
        aggregate = await service.send_notification(command)
        return Response(
            status_code=201,
            headers={"Location": f"/api/v1/notifications/{aggregate.uid}"},
        )

    @router.get("/{uid}", response_model=NotificationResponse)
    async def get_notification(
        uid: UUID,
        principal: Principal = Depends(require_principal),
    ) -> NotificationResponse:
        aggregate = await service.get_notification(uid)
        if aggregate is None:
            raise HTTPException(status_code=404, detail="Notification not found")
        if not is_staff(principal) and str(aggregate.user_uid) != principal.sub:
            raise HTTPException(status_code=403, detail="Forbidden")
        return _to_response(aggregate)

    @router.get("", response_model=NotificationFindAllResponse)
    async def list_notifications(
        principal: Principal = Depends(require_principal),
        user_uid: UUID | None = Query(None),
        notification_type: NotificationType | None = Query(None),
        channel: Channel | None = Query(None),
        status: NotificationStatus | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query("created_at:desc"),
    ) -> NotificationFindAllResponse:
        if not is_staff(principal):
            try:
                user_uid = UUID(principal.sub)
            except ValueError:
                raise HTTPException(status_code=403, detail="Cannot resolve caller identity")
        option = NotificationFilterOption(
            user_uid=user_uid,
            notification_type=notification_type,
            channel=channel,
            status=status,
            page=page,
            size=size,
            sort=sort,
        )
        aggregates, total = await service.list_notifications(option)
        return NotificationFindAllResponse(
            total=total,
            items=[_to_response(a) for a in aggregates],
        )

    @router.patch("/{uid}/sent", response_model=NotificationResponse)
    async def mark_as_sent(
        uid: UUID,
        _staff: Principal = Depends(require_staff),
    ) -> NotificationResponse:
        aggregate = await service.mark_as_sent(uid)
        if aggregate is None:
            raise HTTPException(status_code=404, detail="Notification not found")
        return _to_response(aggregate)

    @router.patch("/{uid}/failed", response_model=NotificationResponse)
    async def mark_as_failed(
        uid: UUID,
        _staff: Principal = Depends(require_staff),
    ) -> NotificationResponse:
        aggregate = await service.mark_as_failed(uid)
        if aggregate is None:
            raise HTTPException(status_code=404, detail="Notification not found")
        return _to_response(aggregate)

    return router


def _to_response(aggregate) -> NotificationResponse:
    return NotificationResponse(
        uid=aggregate.uid,
        user_uid=aggregate.user_uid,
        notification_type=aggregate.notification_type,
        channel=aggregate.channel,
        status=aggregate.status,
        title=aggregate.title,
        content=aggregate.content,
        sent_at=aggregate.sent_at,
        created_at=aggregate.created_at,
        updated_at=aggregate.updated_at,
    )
