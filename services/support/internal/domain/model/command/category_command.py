from uuid import UUID

from pydantic import BaseModel

from internal.domain.model.command import NonBlankStr100, Str500


class CreateCategoryCommand(BaseModel):
    name: NonBlankStr100
    description: Str500 | None = None
    parent_uid: UUID | None = None


class UpdateCategoryCommand(BaseModel):
    name: NonBlankStr100 | None = None
    description: Str500 | None = None
    parent_uid: UUID | None = None
