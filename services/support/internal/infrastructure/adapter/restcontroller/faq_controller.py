from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query

from internal.domain.aggregate.faq_aggregate import FaqAggregate
from internal.domain.error import UnknownReferenceError
from internal.domain.model.command.faq_command import CreateFaqCommand, UpdateFaqCommand
from internal.domain.model.option.faq_option import FaqSearchOption
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.schemas import (
    CreateFaqRequest,
    FaqResponse,
    PaginatedFaqResponse,
    UpdateFaqRequest,
)
from internal.infrastructure.support.auth import Principal, is_staff, require_manager, require_principal


def create_faq_router(service: SupportService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/faqs",
        tags=["faqs"],
        dependencies=[Depends(require_principal)],
    )

    @router.post("", status_code=201, response_model=FaqResponse, dependencies=[Depends(require_manager)])
    async def create_faq(request: CreateFaqRequest) -> FaqResponse:
        command = CreateFaqCommand(
            question=request.question,
            answer=request.answer,
            category_uid=request.category_uid,
            published=request.published,
        )
        try:
            return _to_response(await service.create_faq(command))
        except UnknownReferenceError as exc:
            raise HTTPException(status_code=400, detail=str(exc)) from exc

    @router.get("/{uid}", response_model=FaqResponse)
    async def get_faq(
        uid: UUID,
        increment_view: bool = Query(True),
        principal: Principal = Depends(require_principal),
    ) -> FaqResponse:
        faq = await service.get_faq(uid, increment_view=increment_view)
        if faq is None:
            raise HTTPException(status_code=404, detail="FAQ not found")
        if not faq.published and not is_staff(principal):
            raise HTTPException(status_code=404, detail="FAQ not found")
        return _to_response(faq)

    @router.get("", response_model=PaginatedFaqResponse)
    async def search_faqs(
        category_uid: UUID | None = Query(None),
        q: str | None = Query(None, description="search query"),
        published_only: bool = Query(True),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query("view_count:desc"),
        principal: Principal = Depends(require_principal),
    ) -> PaginatedFaqResponse:
        effective_published_only = published_only or not is_staff(principal)
        option = FaqSearchOption(
            category_uid=category_uid,
            query=q,
            published_only=effective_published_only,
            page=page,
            size=size,
            sort=sort,
        )
        faqs, total = await service.search_faqs(option)
        return PaginatedFaqResponse(
            content=[_to_response(f) for f in faqs],
            total=total,
            page=page,
            size=size,
        )

    @router.patch("/{uid}", response_model=FaqResponse, dependencies=[Depends(require_manager)])
    async def update_faq(uid: UUID, request: UpdateFaqRequest) -> FaqResponse:
        command = UpdateFaqCommand(
            question=request.question,
            answer=request.answer,
            category_uid=request.category_uid,
            published=request.published,
        )
        try:
            faq = await service.update_faq(uid, command)
        except UnknownReferenceError as exc:
            raise HTTPException(status_code=400, detail=str(exc)) from exc
        if faq is None:
            raise HTTPException(status_code=404, detail="FAQ not found")
        return _to_response(faq)

    @router.delete("/{uid}", status_code=204, dependencies=[Depends(require_manager)])
    async def delete_faq(uid: UUID) -> None:
        if not await service.delete_faq(uid):
            raise HTTPException(status_code=404, detail="FAQ not found")

    return router


def _to_response(f: FaqAggregate) -> FaqResponse:
    return FaqResponse(
        uid=f.uid,
        category_uid=f.category_uid,
        question=f.question,
        answer=f.answer,
        view_count=f.view_count,
        published=f.published,
        created_at=f.created_at,
        updated_at=f.updated_at,
    )
