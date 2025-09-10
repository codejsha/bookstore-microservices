from datetime import datetime
from uuid import UUID


class FaqAggregate:
    def __init__(
        self,
        uid: UUID,
        question: str,
        answer: str,
        view_count: int,
        published: bool,
        created_at: datetime,
        updated_at: datetime,
        category_id: int | None = None,
        category_uid: UUID | None = None,
    ):
        self.uid = uid
        self.category_id = category_id
        self.category_uid = category_uid
        self.question = question
        self.answer = answer
        self.view_count = view_count
        self.published = published
        self.created_at = created_at
        self.updated_at = updated_at
