from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus


class TicketFilterOption(BaseModel):
    customer_uid: UUID | None = None
    assignee_uid: UUID | None = None
    category_uid: UUID | None = None
    status: TicketStatus | None = None
    priority: TicketPriority | None = None
    page: int = 0
    size: int = 20
    sort: str = "created_at:desc"
