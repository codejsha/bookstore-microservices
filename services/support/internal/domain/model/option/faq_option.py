from uuid import UUID

from pydantic import BaseModel


class FaqSearchOption(BaseModel):
    category_uid: UUID | None = None
    query: str | None = None
    published_only: bool = True
    page: int = 0
    size: int = 20
    sort: str = "view_count:desc"
