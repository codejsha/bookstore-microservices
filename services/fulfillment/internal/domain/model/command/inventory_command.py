from typing import Annotated
from uuid import UUID

from pydantic import AfterValidator, BaseModel, Field, StringConstraints, field_validator, model_validator


def _require_non_blank(value: str) -> str:
    if not value.strip():
        raise ValueError("must not be blank")
    return value


def _non_blank(max_length: int) -> object:
    return Annotated[str, StringConstraints(max_length=max_length), AfterValidator(_require_non_blank)]


NonBlankStr200 = _non_blank(200)
NonBlankStr255 = _non_blank(255)
NonBlankStr500 = _non_blank(500)


class StockReceiveCommand(BaseModel):
    edition_uid: UUID
    warehouse_uid: UUID
    quantity: int = Field(ge=1)
    reason: NonBlankStr255 | None = None


class StockReleaseCommand(BaseModel):
    edition_uid: UUID
    warehouse_uid: UUID
    quantity: int = Field(ge=1)
    reason: NonBlankStr255 | None = None


class StockAdjustCommand(BaseModel):
    edition_uid: UUID
    warehouse_uid: UUID
    quantity: int
    reason: NonBlankStr255

    @field_validator("quantity")
    @classmethod
    def _non_zero(cls, value: int) -> int:
        if value == 0:
            raise ValueError("must be non-zero")
        return value


class StockReserveCommand(BaseModel):
    edition_uid: UUID
    warehouse_uid: UUID
    quantity: int = Field(ge=1)
    reason: NonBlankStr255 | None = None


class StockOrderReserveCommand(BaseModel):
    order_uid: NonBlankStr200
    edition_id: int = Field(gt=0)
    quantity: int
    reason: NonBlankStr255 | None = None


class StockOrderReleaseCommand(BaseModel):
    order_uid: NonBlankStr200
    edition_id: int = Field(gt=0)
    quantity: int = 0
    reason: NonBlankStr255 | None = None


class WarehouseCreateCommand(BaseModel):
    name: NonBlankStr255
    address: NonBlankStr255 | None = None
    capacity: int = Field(ge=0)


class WarehouseUpdateCommand(BaseModel):
    name: NonBlankStr255 | None = None
    address: NonBlankStr255 | None = None
    capacity: int | None = Field(None, ge=0)


class StockTransferCommand(BaseModel):
    edition_uid: UUID
    source_warehouse_uid: UUID
    target_warehouse_uid: UUID
    quantity: int = Field(ge=1)
    reason: NonBlankStr255 | None = None

    @model_validator(mode="after")
    def _check_distinct_warehouses(self) -> StockTransferCommand:
        if self.source_warehouse_uid == self.target_warehouse_uid:
            raise ValueError("source_warehouse_uid and target_warehouse_uid must differ")
        return self


class StockAuditItemCommand(BaseModel):
    edition_uid: UUID
    actual_quantity: int = Field(ge=0)


class StockAuditCreateCommand(BaseModel):
    warehouse_uid: UUID
    items: list[StockAuditItemCommand] = Field(min_length=1)
    notes: NonBlankStr500 | None = None

    @model_validator(mode="after")
    def _check_unique_editions(self) -> StockAuditCreateCommand:
        seen: set[UUID] = set()
        for item in self.items:
            if item.edition_uid in seen:
                raise ValueError(f"items must not contain duplicate edition_uid: {item.edition_uid}")
            seen.add(item.edition_uid)
        return self


class MonthlyClosingCommand(BaseModel):
    warehouse_uid: UUID
    year: int = Field(ge=2000, le=2100)
    month: int = Field(ge=1, le=12)
