from datetime import UTC, datetime
from uuid import UUID, uuid7

import structlog

from internal.application.port.repo.repos import (
    FaqRepository,
    TicketCategoryRepository,
    TicketCommentRepository,
    TicketRepository,
)
from internal.domain.aggregate.faq_aggregate import FaqAggregate
from internal.domain.aggregate.ticket_aggregate import TicketAggregate
from internal.domain.aggregate.ticket_category_aggregate import TicketCategoryAggregate
from internal.domain.aggregate.ticket_comment_aggregate import TicketCommentAggregate
from internal.domain.constant.ticket_status import TicketStatus
from internal.domain.error import ConflictError, UnknownReferenceError
from internal.domain.model.command.category_command import (
    CreateCategoryCommand,
    UpdateCategoryCommand,
)
from internal.domain.model.command.faq_command import CreateFaqCommand, UpdateFaqCommand
from internal.domain.model.command.ticket_command import (
    AddCommentCommand,
    CreateTicketCommand,
    UpdateTicketCommand,
    UpdateTicketStatusCommand,
)
from internal.domain.model.option.faq_option import FaqSearchOption
from internal.domain.model.option.ticket_option import TicketFilterOption

logger = structlog.get_logger()


class SupportService:
    def __init__(
        self,
        ticket_repo: TicketRepository,
        comment_repo: TicketCommentRepository,
        category_repo: TicketCategoryRepository,
        faq_repo: FaqRepository,
    ):
        self._ticket_repo = ticket_repo
        self._comment_repo = comment_repo
        self._category_repo = category_repo
        self._faq_repo = faq_repo

    async def _resolve_category_id(self, category_uid: UUID | None) -> int | None:
        if category_uid is None:
            return None
        category_id = await self._category_repo.find_id_by_uid(category_uid)
        if category_id is None:
            raise UnknownReferenceError("Category not found")
        return category_id

    # ─── Tickets ──────────────────────────────────────────────────────

    async def create_ticket(self, command: CreateTicketCommand) -> TicketAggregate:
        now = datetime.now(UTC)
        category_id = await self._resolve_category_id(command.category_uid)
        ticket = TicketAggregate(
            uid=uuid7(),
            customer_uid=command.customer_uid,
            category_id=category_id,
            subject=command.subject,
            description=command.description,
            status=TicketStatus.OPEN,
            priority=command.priority,
            created_at=now,
            updated_at=now,
        )
        return await self._ticket_repo.save(ticket)

    async def get_ticket(self, uid: UUID) -> TicketAggregate | None:
        return await self._ticket_repo.find_by_uid(uid)

    async def list_tickets(self, option: TicketFilterOption) -> tuple[list[TicketAggregate], int]:
        return await self._ticket_repo.find_all(option)

    async def update_ticket(self, uid: UUID, command: UpdateTicketCommand) -> TicketAggregate | None:
        category_provided = command.category_uid is not None
        category_id = await self._resolve_category_id(command.category_uid)
        now = datetime.now(UTC)

        def _mutate(ticket: TicketAggregate) -> TicketAggregate:
            if command.subject is not None:
                ticket.subject = command.subject
            if command.description is not None:
                ticket.description = command.description
            if command.priority is not None:
                ticket.priority = command.priority
            if command.assignee_uid is not None:
                ticket.assignee_uid = command.assignee_uid
            if category_provided:
                ticket.category_id = category_id
            ticket.updated_at = now
            return ticket

        return await self._ticket_repo.update(uid, _mutate)

    async def update_ticket_status(self, uid: UUID, command: UpdateTicketStatusCommand) -> TicketAggregate | None:
        now = datetime.now(UTC)

        def _mutate(ticket: TicketAggregate) -> TicketAggregate:
            ticket.status = command.status
            ticket.updated_at = now
            if command.status == TicketStatus.RESOLVED and ticket.resolved_at is None:
                ticket.resolved_at = now
            return ticket

        return await self._ticket_repo.update(uid, _mutate)

    # ─── Comments ─────────────────────────────────────────────────────

    async def add_comment(self, command: AddCommentCommand) -> TicketCommentAggregate | None:
        ticket_id = await self._ticket_repo.find_id_by_uid(command.ticket_uid)
        if ticket_id is None:
            return None
        now = datetime.now(UTC)
        comment = TicketCommentAggregate(
            uid=uuid7(),
            ticket_id=ticket_id,
            author_uid=command.author_uid,
            author_role=command.author_role,
            body=command.body,
            internal=command.internal,
            created_at=now,
            updated_at=now,
        )
        return await self._comment_repo.save(comment)

    async def list_comments(
        self, ticket_uid: UUID, include_internal: bool = True
    ) -> list[TicketCommentAggregate] | None:
        ticket_id = await self._ticket_repo.find_id_by_uid(ticket_uid)
        if ticket_id is None:
            return None
        return await self._comment_repo.find_by_ticket_id(ticket_id, include_internal=include_internal)

    # ─── Categories ───────────────────────────────────────────────────

    async def create_category(self, command: CreateCategoryCommand) -> TicketCategoryAggregate:
        now = datetime.now(UTC)
        if await self._category_repo.find_by_name(command.name) is not None:
            raise ConflictError("Category name already exists")
        parent_id = await self._resolve_category_id(command.parent_uid)
        category = TicketCategoryAggregate(
            uid=uuid7(),
            name=command.name,
            description=command.description,
            parent_id=parent_id,
            created_at=now,
            updated_at=now,
        )
        return await self._category_repo.save(category)

    async def get_category(self, uid: UUID) -> TicketCategoryAggregate | None:
        return await self._category_repo.find_by_uid(uid)

    async def list_categories(self) -> list[TicketCategoryAggregate]:
        return await self._category_repo.find_all()

    async def update_category(self, uid: UUID, command: UpdateCategoryCommand) -> TicketCategoryAggregate | None:
        category = await self._category_repo.find_by_uid(uid)
        if category is None:
            return None
        now = datetime.now(UTC)
        if command.name is not None and command.name != category.name:
            existing = await self._category_repo.find_by_name(command.name)
            if existing is not None and existing.uid != category.uid:
                raise ConflictError("Category name already exists")
        if command.name is not None:
            category.name = command.name
        if command.description is not None:
            category.description = command.description
        if command.parent_uid is not None:
            category.parent_id = await self._resolve_category_id(command.parent_uid)
        category.updated_at = now
        return await self._category_repo.save(category)

    async def delete_category(self, uid: UUID) -> bool:
        return await self._category_repo.delete_by_uid(uid)

    # ─── FAQ ──────────────────────────────────────────────────────────

    async def create_faq(self, command: CreateFaqCommand) -> FaqAggregate:
        now = datetime.now(UTC)
        category_id = await self._resolve_category_id(command.category_uid)
        faq = FaqAggregate(
            uid=uuid7(),
            category_id=category_id,
            question=command.question,
            answer=command.answer,
            view_count=0,
            published=command.published,
            created_at=now,
            updated_at=now,
        )
        return await self._faq_repo.save(faq)

    async def get_faq(self, uid: UUID, increment_view: bool = True) -> FaqAggregate | None:
        faq = await self._faq_repo.find_by_uid(uid)
        if faq is None:
            return None
        if increment_view and faq.published:
            try:
                await self._faq_repo.increment_view_count(uid)
                faq.view_count += 1
            except Exception:
                logger.warning("faq view_count increment failed", faq_uid=str(uid), exc_info=True)
        return faq

    async def search_faqs(self, option: FaqSearchOption) -> tuple[list[FaqAggregate], int]:
        return await self._faq_repo.search(option)

    async def update_faq(self, uid: UUID, command: UpdateFaqCommand) -> FaqAggregate | None:
        category_provided = command.category_uid is not None
        category_id = await self._resolve_category_id(command.category_uid)
        now = datetime.now(UTC)

        def _mutate(faq: FaqAggregate) -> FaqAggregate:
            if command.question is not None:
                faq.question = command.question
            if command.answer is not None:
                faq.answer = command.answer
            if command.published is not None:
                faq.published = command.published
            if category_provided:
                faq.category_id = category_id
            faq.updated_at = now
            return faq

        return await self._faq_repo.update(uid, _mutate)

    async def delete_faq(self, uid: UUID) -> bool:
        return await self._faq_repo.delete_by_uid(uid)
