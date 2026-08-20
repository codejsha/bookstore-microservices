from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query, Response

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.model.command.notification_command import (
    CreateTemplateCommand,
    UpdateTemplateCommand,
)
from internal.domain.model.option.notification_option import TemplateFilterOption
from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.restcontroller.schemas import (
    CreateTemplateRequest,
    TemplateFindAllResponse,
    TemplateResponse,
    UpdateTemplateRequest,
)
from internal.infrastructure.support.auth import require_principal, require_staff


def create_template_router(service: NotificationService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/templates",
        tags=["templates"],
        dependencies=[Depends(require_principal)],
    )

    @router.post("", status_code=201, dependencies=[Depends(require_staff)])
    async def create_template(request: CreateTemplateRequest) -> Response:
        command = CreateTemplateCommand(
            notification_type=request.notification_type,
            channel=request.channel,
            title_template=request.title_template,
            content_template=request.content_template,
        )
        aggregate = await service.create_template(command)
        return Response(
            status_code=201,
            headers={"Location": f"/api/v1/templates/{aggregate.uid}"},
        )

    @router.get("/{uid}", response_model=TemplateResponse)
    async def get_template(uid: UUID) -> TemplateResponse:
        aggregate = await service.get_template(uid)
        if aggregate is None:
            raise HTTPException(status_code=404, detail="Template not found")
        return _to_response(aggregate)

    @router.get("", response_model=TemplateFindAllResponse)
    async def list_templates(
        notification_type: NotificationType | None = Query(None),
        channel: Channel | None = Query(None),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query("created_at:desc"),
    ) -> TemplateFindAllResponse:
        option = TemplateFilterOption(
            notification_type=notification_type,
            channel=channel,
            page=page,
            size=size,
            sort=sort,
        )
        aggregates, total = await service.list_templates(option)
        return TemplateFindAllResponse(
            total=total,
            items=[_to_response(a) for a in aggregates],
        )

    @router.put("/{uid}", response_model=TemplateResponse, dependencies=[Depends(require_staff)])
    async def update_template(uid: UUID, request: UpdateTemplateRequest) -> TemplateResponse:
        command = UpdateTemplateCommand(
            notification_type=request.notification_type,
            channel=request.channel,
            title_template=request.title_template,
            content_template=request.content_template,
        )
        aggregate = await service.update_template(uid, command)
        if aggregate is None:
            raise HTTPException(status_code=404, detail="Template not found")
        return _to_response(aggregate)

    @router.delete("/{uid}", status_code=204, dependencies=[Depends(require_staff)])
    async def delete_template(uid: UUID) -> None:
        if not await service.delete_template(uid):
            raise HTTPException(status_code=404, detail="Template not found")

    return router


def _to_response(aggregate) -> TemplateResponse:
    return TemplateResponse(
        uid=aggregate.uid,
        notification_type=aggregate.notification_type,
        channel=aggregate.channel,
        title_template=aggregate.title_template,
        content_template=aggregate.content_template,
        created_at=aggregate.created_at,
        updated_at=aggregate.updated_at,
    )
