from uuid import UUID

from pydantic import BaseModel

from internal.domain.model.command import NonBlankStr500, NonBlankStr16000


class CreateFaqCommand(BaseModel):
    question: NonBlankStr500
    answer: NonBlankStr16000
    category_uid: UUID | None = None
    published: bool = False


class UpdateFaqCommand(BaseModel):
    question: NonBlankStr500 | None = None
    answer: NonBlankStr16000 | None = None
    category_uid: UUID | None = None
    published: bool | None = None
