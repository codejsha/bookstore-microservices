from datetime import datetime
from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.author_role import AuthorRole
from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus
from internal.domain.model.command import (
    NonBlankStr100,
    NonBlankStr255,
    NonBlankStr500,
    NonBlankStr16000,
    Str500,
)

# ─── Tickets ─────────────────────────────────────────────────────────────


class CreateTicketRequest(BaseModel):
    customer_uid: UUID
    subject: NonBlankStr255
    description: NonBlankStr16000
    priority: TicketPriority = TicketPriority.MEDIUM
    category_uid: UUID | None = None


class UpdateTicketRequest(BaseModel):
    subject: NonBlankStr255 | None = None
    description: NonBlankStr16000 | None = None
    priority: TicketPriority | None = None
    category_uid: UUID | None = None
    assignee_uid: UUID | None = None


class UpdateTicketStatusRequest(BaseModel):
    status: TicketStatus


class TicketResponse(BaseModel):
    uid: UUID
    customer_uid: UUID
    category_uid: UUID | None = None
    assignee_uid: UUID | None = None
    subject: str
    description: str
    status: TicketStatus
    priority: TicketPriority
    resolved_at: datetime | None = None
    created_at: datetime
    updated_at: datetime


class PaginatedTicketResponse(BaseModel):
    content: list[TicketResponse]
    total: int
    page: int
    size: int


# ─── Comments ────────────────────────────────────────────────────────────


class AddCommentRequest(BaseModel):
    body: NonBlankStr16000
    internal: bool = False


class CommentResponse(BaseModel):
    uid: UUID
    ticket_uid: UUID
    author_uid: UUID
    author_role: AuthorRole
    body: str
    internal: bool
    created_at: datetime
    updated_at: datetime


# ─── Categories ──────────────────────────────────────────────────────────


class CreateCategoryRequest(BaseModel):
    name: NonBlankStr100
    description: Str500 | None = None
    parent_uid: UUID | None = None


class UpdateCategoryRequest(BaseModel):
    name: NonBlankStr100 | None = None
    description: Str500 | None = None
    parent_uid: UUID | None = None


class CategoryResponse(BaseModel):
    uid: UUID
    name: str
    description: str | None = None
    parent_uid: UUID | None = None
    created_at: datetime
    updated_at: datetime


# ─── FAQ ─────────────────────────────────────────────────────────────────


class CreateFaqRequest(BaseModel):
    question: NonBlankStr500
    answer: NonBlankStr16000
    category_uid: UUID | None = None
    published: bool = False


class UpdateFaqRequest(BaseModel):
    question: NonBlankStr500 | None = None
    answer: NonBlankStr16000 | None = None
    category_uid: UUID | None = None
    published: bool | None = None


class FaqResponse(BaseModel):
    uid: UUID
    category_uid: UUID | None = None
    question: str
    answer: str
    view_count: int
    published: bool
    created_at: datetime
    updated_at: datetime


class PaginatedFaqResponse(BaseModel):
    content: list[FaqResponse]
    total: int
    page: int
    size: int
