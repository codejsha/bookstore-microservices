from datetime import datetime
from uuid import UUID

from internal.domain.constant.author_role import AuthorRole


class TicketCommentAggregate:
    def __init__(
        self,
        uid: UUID,
        ticket_id: int,
        author_uid: UUID,
        author_role: AuthorRole,
        body: str,
        internal: bool,
        created_at: datetime,
        updated_at: datetime,
    ):
        self.uid = uid
        self.ticket_id = ticket_id
        self.author_uid = author_uid
        self.author_role = author_role
        self.body = body
        self.internal = internal
        self.created_at = created_at
        self.updated_at = updated_at
