from datetime import datetime
from uuid import UUID


class TicketCategoryAggregate:
    def __init__(
        self,
        uid: UUID,
        name: str,
        created_at: datetime,
        updated_at: datetime,
        description: str | None = None,
        parent_id: int | None = None,
        parent_uid: UUID | None = None,
    ):
        self.uid = uid
        self.name = name
        self.description = description
        self.parent_id = parent_id
        self.parent_uid = parent_uid
        self.created_at = created_at
        self.updated_at = updated_at
