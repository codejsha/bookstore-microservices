from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, Query
from pydantic import ValidationError

from generated.application.port.model.faq_create_request import FaqCreateRequest
from generated.application.port.model.faq_find_all_response import FaqFindAllResponse
from generated.application.port.model.faq_find_response import FaqFindResponse
from generated.application.port.model.faq_item import FaqItem
from generated.application.port.model.faq_update_request import FaqUpdateRequest
from generated.application.port.model.faq_update_response import FaqUpdateResponse
from internal.domain.aggregate.faq_aggregate import FaqAggregate
from internal.domain.error import UnknownReferenceError
from internal.domain.model.command.faq_command import CreateFaqCommand, UpdateFaqCommand
from internal.domain.model.option.faq_option import FaqSearchOption
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.error_mapping import (
    invalid_command_error,
    unknown_reference_error,
)
from internal.infrastructure.support.auth import Principal, is_staff, require_manager, require_principal


def create_faq_router(service: SupportService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/faqs",
        tags=["faqs"],
        dependencies=[Depends(require_principal)],
    )

    @router.post("", status_code=201, response_model=FaqItem, dependencies=[Depends(require_manager)])
    async def create_faq(request: FaqCreateRequest) -> FaqItem:
        try:
            command = CreateFaqCommand(**request.model_dump(exclude_none=True))
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            faq = await service.create_faq(command)
        except UnknownReferenceError as exc:
            raise unknown_reference_error(exc) from exc
        return _to_faq_item(faq)

    @router.get("/{uid}", response_model=FaqFindResponse)
    async def get_faq(
        uid: UUID,
        increment_view: bool = Query(True),
        principal: Principal = Depends(require_principal),
    ) -> FaqFindResponse:
        faq = await service.get_faq(uid, increment_view=increment_view)
        if faq is None:
            raise HTTPException(status_code=404, detail="FAQ not found")
        if not faq.published and not is_staff(principal):
            raise HTTPException(status_code=404, detail="FAQ not found")
        return FaqFindResponse(**_to_faq_item(faq).model_dump())

    @router.get("", response_model=FaqFindAllResponse)
    async def search_faqs(
        category_uid: UUID | None = Query(None),
        q: str | None = Query(None, description="search query"),
        published_only: bool = Query(True),
        page: int = Query(0, ge=0),
        size: int = Query(20, ge=1, le=100),
        sort: str = Query("view_count:desc"),
        principal: Principal = Depends(require_principal),
    ) -> FaqFindAllResponse:
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
        return FaqFindAllResponse(
            items=[_to_faq_item(f) for f in faqs],
            total=total,
        )

    @router.patch("/{uid}", response_model=FaqUpdateResponse, dependencies=[Depends(require_manager)])
    async def update_faq(uid: UUID, request: FaqUpdateRequest) -> FaqUpdateResponse:
        try:
            command = UpdateFaqCommand(**request.model_dump(exclude_unset=True))
        except ValidationError as e:
            raise invalid_command_error(e) from e
        try:
            faq = await service.update_faq(uid, command)
        except UnknownReferenceError as exc:
            raise unknown_reference_error(exc) from exc
        if faq is None:
            raise HTTPException(status_code=404, detail="FAQ not found")
        return FaqUpdateResponse(**_to_faq_item(faq).model_dump())

    @router.delete("/{uid}", status_code=204, dependencies=[Depends(require_manager)])
    async def delete_faq(uid: UUID) -> None:
        if not await service.delete_faq(uid):
            raise HTTPException(status_code=404, detail="FAQ not found")

    return router


def _to_faq_item(f: FaqAggregate) -> FaqItem:
    return FaqItem(
        uid=str(f.uid),
        category_uid=str(f.category_uid) if f.category_uid else None,
        question=f.question,
        answer=f.answer,
        view_count=f.view_count,
        published=f.published,
        created_at=f.created_at,
        updated_at=f.updated_at,
    )
