from datetime import datetime
from uuid import UUID

from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus


class TicketAggregate:
    def __init__(
        self,
        uid: UUID,
        customer_uid: UUID,
        subject: str,
        description: str,
        status: TicketStatus,
        priority: TicketPriority,
        created_at: datetime,
        updated_at: datetime,
        category_id: int | None = None,
        assignee_uid: UUID | None = None,
        resolved_at: datetime | None = None,
        category_uid: UUID | None = None,
    ):
        self.uid = uid
        self.customer_uid = customer_uid
        self.category_id = category_id
        self.category_uid = category_uid
        self.assignee_uid = assignee_uid
        self.subject = subject
        self.description = description
        self.status = status
        self.priority = priority
        self.resolved_at = resolved_at
        self.created_at = created_at
        self.updated_at = updated_at
