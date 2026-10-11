from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException
from pydantic import ValidationError

from generated.application.port.model.category_create_request import CategoryCreateRequest
from generated.application.port.model.category_find_all_response import CategoryFindAllResponse
from generated.application.port.model.category_find_response import CategoryFindResponse
from generated.application.port.model.category_item import CategoryItem
from generated.application.port.model.category_update_request import CategoryUpdateRequest
from generated.application.port.model.category_update_response import CategoryUpdateResponse
from internal.domain.aggregate.ticket_category_aggregate import TicketCategoryAggregate
from internal.domain.error import ConflictError, UnknownReferenceError
from internal.domain.model.command.category_command import (
    CreateCategoryCommand,
    UpdateCategoryCommand,
)
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.error_mapping import (
    conflict_error,
    invalid_command_error,
    unknown_reference_error,
)
from internal.infrastructure.support.auth import require_manager, require_principal


def create_category_router(service: SupportService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/categories",
        tags=["categories"],
        dependencies=[Depends(require_principal)],
    )

    @router.post("", status_code=201, response_model=CategoryItem, dependencies=[Depends(require_manager)])
    async def create_category(request: CategoryCreateRequest) -> CategoryItem:
        try:
            command = CreateCategoryCommand(**request.model_dump(exclude_none=True))
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            category = await service.create_category(command)
        except ConflictError as exc:
            raise conflict_error(exc) from exc
        except UnknownReferenceError as exc:
            raise unknown_reference_error(exc) from exc
        return _to_category_item(category)

    @router.get("/{uid}", response_model=CategoryFindResponse)
    async def get_category(uid: UUID) -> CategoryFindResponse:
        category = await service.get_category(uid)
        if category is None:
            raise HTTPException(status_code=404, detail="Category not found")
        return CategoryFindResponse(**_to_category_item(category).model_dump())

    @router.get("", response_model=CategoryFindAllResponse)
    async def list_categories() -> CategoryFindAllResponse:
        return CategoryFindAllResponse(items=[_to_category_item(c) for c in await service.list_categories()])

    @router.patch("/{uid}", response_model=CategoryUpdateResponse, dependencies=[Depends(require_manager)])
    async def update_category(uid: UUID, request: CategoryUpdateRequest) -> CategoryUpdateResponse:
        try:
            command = UpdateCategoryCommand(**request.model_dump(exclude_unset=True))
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            category = await service.update_category(uid, command)
        except ConflictError as exc:
            raise conflict_error(exc) from exc
        except UnknownReferenceError as exc:
            raise unknown_reference_error(exc) from exc
        if category is None:
            raise HTTPException(status_code=404, detail="Category not found")
        return CategoryUpdateResponse(**_to_category_item(category).model_dump())

    @router.delete("/{uid}", status_code=204, dependencies=[Depends(require_manager)])
    async def delete_category(uid: UUID) -> None:
        if not await service.delete_category(uid):
            raise HTTPException(status_code=404, detail="Category not found")

    return router


def _to_category_item(c: TicketCategoryAggregate) -> CategoryItem:
    return CategoryItem(
        uid=str(c.uid),
        name=c.name,
        description=c.description,
        parent_uid=str(c.parent_uid) if c.parent_uid else None,
        created_at=c.created_at,
        updated_at=c.updated_at,
    )
