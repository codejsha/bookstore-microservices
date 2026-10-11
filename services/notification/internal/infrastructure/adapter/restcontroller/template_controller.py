from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query, Response
from fastapi.exceptions import RequestValidationError
from pydantic import ValidationError
from sqlalchemy.exc import IntegrityError

from generated.application.port.model.template_create_request import TemplateCreateRequest
from generated.application.port.model.template_find_all_response import TemplateFindAllResponse
from generated.application.port.model.template_find_response import TemplateFindResponse
from generated.application.port.model.template_item import TemplateItem
from generated.application.port.model.template_update_request import TemplateUpdateRequest
from generated.application.port.model.template_update_response import TemplateUpdateResponse
from internal.domain.aggregate.template_aggregate import TemplateAggregate
from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.error import TemplateAlreadyExistsError
from internal.domain.model.command.notification_command import (
    CreateTemplateCommand,
    UpdateTemplateCommand,
)
from internal.domain.model.option.notification_option import TemplateFilterOption
from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.support.auth import require_manager, require_principal


def create_template_router(service: NotificationService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/templates",
        tags=["templates"],
        dependencies=[Depends(require_principal)],
    )

    @router.post("", status_code=201, dependencies=[Depends(require_manager)])
    async def create_template(request: TemplateCreateRequest) -> Response:
        try:
            command = CreateTemplateCommand(**request.model_dump())
        except ValidationError as e:
            raise RequestValidationError(e.errors(include_url=False, include_input=False)) from e
        try:
            aggregate = await service.create_template(command)
        except (TemplateAlreadyExistsError, IntegrityError) as exc:
            raise _conflict(exc) from exc
        return Response(
            status_code=201,
            headers={"Location": f"/api/v1/templates/{aggregate.uid}"},
        )

    @router.get("/{uid}", response_model=TemplateFindResponse)
    async def get_template(uid: UUID) -> TemplateFindResponse:
        aggregate = await service.get_template(uid)
        if aggregate is None:
            raise HTTPException(status_code=404, detail="Template not found")
        return TemplateFindResponse(**_to_template_item(aggregate).model_dump())

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
            items=[_to_template_item(a) for a in aggregates],
        )

    @router.put("/{uid}", response_model=TemplateUpdateResponse, dependencies=[Depends(require_manager)])
    async def update_template(uid: UUID, request: TemplateUpdateRequest) -> TemplateUpdateResponse:
        try:
            command = UpdateTemplateCommand(**request.model_dump(exclude_unset=True))
        except ValidationError as e:
            raise RequestValidationError(e.errors(include_url=False, include_input=False)) from e
        try:
            aggregate = await service.update_template(uid, command)
        except (TemplateAlreadyExistsError, IntegrityError) as exc:
            raise _conflict(exc) from exc
        if aggregate is None:
            raise HTTPException(status_code=404, detail="Template not found")
        return TemplateUpdateResponse(**_to_template_item(aggregate).model_dump())

    @router.delete("/{uid}", status_code=204, dependencies=[Depends(require_manager)])
    async def delete_template(uid: UUID) -> None:
        if not await service.delete_template(uid):
            raise HTTPException(status_code=404, detail="Template not found")

    return router


def _conflict(exc: Exception) -> HTTPException:
    if isinstance(exc, TemplateAlreadyExistsError):
        detail = f"Template already exists for {exc.notification_type.value} via {exc.channel.value}"
    else:
        detail = "Template already exists for this notification type and channel"
    return HTTPException(status_code=409, detail=detail)


def _to_template_item(aggregate: TemplateAggregate) -> TemplateItem:
    return TemplateItem(
        uid=str(aggregate.uid),
        notification_type=aggregate.notification_type,
        channel=aggregate.channel,
        title_template=aggregate.title_template,
        content_template=aggregate.content_template,
        created_at=aggregate.created_at,
        updated_at=aggregate.updated_at,
    )
