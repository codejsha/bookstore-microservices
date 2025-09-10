from uuid import UUID

from pydantic import BaseModel


class CreateFaqCommand(BaseModel):
    question: str
    answer: str
    category_uid: UUID | None = None
    published: bool = False


class UpdateFaqCommand(BaseModel):
    question: str | None = None
    answer: str | None = None
    category_uid: UUID | None = None
    published: bool | None = None
