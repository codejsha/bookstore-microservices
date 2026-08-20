from uuid import UUID

from pydantic import BaseModel


class CreateCategoryCommand(BaseModel):
    name: str
    description: str | None = None
    parent_uid: UUID | None = None


class UpdateCategoryCommand(BaseModel):
    name: str | None = None
    description: str | None = None
    parent_uid: UUID | None = None
