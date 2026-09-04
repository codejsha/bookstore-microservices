from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.author_role import AuthorRole
from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus
from internal.domain.model.command import NonBlankStr255, NonBlankStr16000


class CreateTicketCommand(BaseModel):
    customer_uid: UUID
    subject: NonBlankStr255
    description: NonBlankStr16000
    priority: TicketPriority = TicketPriority.MEDIUM
    category_uid: UUID | None = None


class UpdateTicketCommand(BaseModel):
    subject: NonBlankStr255 | None = None
    description: NonBlankStr16000 | None = None
    priority: TicketPriority | None = None
    category_uid: UUID | None = None
    assignee_uid: UUID | None = None


class UpdateTicketStatusCommand(BaseModel):
    status: TicketStatus


class AddCommentCommand(BaseModel):
    ticket_uid: UUID
    author_uid: UUID
    author_role: AuthorRole
    body: NonBlankStr16000
    internal: bool = False
