from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException

from internal.domain.aggregate.ticket_category_aggregate import TicketCategoryAggregate
from internal.domain.error import ConflictError, UnknownReferenceError
from internal.domain.model.command.category_command import (
    CreateCategoryCommand,
    UpdateCategoryCommand,
)
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.schemas import (
    CategoryResponse,
    CreateCategoryRequest,
    UpdateCategoryRequest,
)
from internal.infrastructure.support.auth import require_principal, require_staff


def create_category_router(service: SupportService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/categories",
        tags=["categories"],
        dependencies=[Depends(require_principal)],
    )

    @router.post("", status_code=201, response_model=CategoryResponse, dependencies=[Depends(require_staff)])
    async def create_category(request: CreateCategoryRequest) -> CategoryResponse:
        command = CreateCategoryCommand(
            name=request.name,
            description=request.description,
            parent_uid=request.parent_uid,
        )
        try:
            return _to_response(await service.create_category(command))
        except ConflictError as exc:
            raise HTTPException(status_code=409, detail=str(exc)) from exc
        except UnknownReferenceError as exc:
            raise HTTPException(status_code=400, detail=str(exc)) from exc

    @router.get("/{uid}", response_model=CategoryResponse)
    async def get_category(uid: UUID) -> CategoryResponse:
        category = await service.get_category(uid)
        if category is None:
            raise HTTPException(status_code=404, detail="Category not found")
        return _to_response(category)

    @router.get("", response_model=list[CategoryResponse])
    async def list_categories() -> list[CategoryResponse]:
        return [_to_response(c) for c in await service.list_categories()]

    @router.patch("/{uid}", response_model=CategoryResponse, dependencies=[Depends(require_staff)])
    async def update_category(uid: UUID, request: UpdateCategoryRequest) -> CategoryResponse:
        command = UpdateCategoryCommand(
            name=request.name,
            description=request.description,
            parent_uid=request.parent_uid,
        )
        try:
            category = await service.update_category(uid, command)
        except ConflictError as exc:
            raise HTTPException(status_code=409, detail=str(exc)) from exc
        except UnknownReferenceError as exc:
            raise HTTPException(status_code=400, detail=str(exc)) from exc
        if category is None:
            raise HTTPException(status_code=404, detail="Category not found")
        return _to_response(category)

    @router.delete("/{uid}", status_code=204, dependencies=[Depends(require_staff)])
    async def delete_category(uid: UUID) -> None:
        if not await service.delete_category(uid):
            raise HTTPException(status_code=404, detail="Category not found")

    return router


def _to_response(c: TicketCategoryAggregate) -> CategoryResponse:
    return CategoryResponse(
        uid=c.uid,
        name=c.name,
        description=c.description,
        parent_uid=c.parent_uid,
        created_at=c.created_at,
        updated_at=c.updated_at,
    )
