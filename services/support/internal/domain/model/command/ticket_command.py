from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.author_role import AuthorRole
from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus


class CreateTicketCommand(BaseModel):
    customer_uid: UUID
    subject: str
    description: str
    priority: TicketPriority = TicketPriority.MEDIUM
    category_uid: UUID | None = None


class UpdateTicketCommand(BaseModel):
    subject: str | None = None
    description: str | None = None
    priority: TicketPriority | None = None
    category_uid: UUID | None = None
    assignee_uid: UUID | None = None


class UpdateTicketStatusCommand(BaseModel):
    status: TicketStatus


class AddCommentCommand(BaseModel):
    ticket_uid: UUID
    author_uid: UUID
    author_role: AuthorRole
    body: str
    internal: bool = False
