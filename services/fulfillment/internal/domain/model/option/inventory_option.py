from typing import Annotated
from uuid import UUID

from pydantic import BaseModel, StringConstraints

from internal.domain.constant.audit_status import AuditStatus
from internal.domain.constant.closing_status import ClosingStatus
from internal.domain.constant.transfer_status import TransferStatus

DEFAULT_INVENTORY_SORT = "created_at,desc;uid,asc"

YearMonth = Annotated[str, StringConstraints(pattern=r"^\d{4}-(0[1-9]|1[0-2])$")]


class StockQueryOption(BaseModel):
    edition_uid: UUID | None = None
    warehouse_uid: UUID | None = None
    page: int = 0
    size: int = 20
    sort: str = DEFAULT_INVENTORY_SORT


class StockHistoryQueryOption(BaseModel):
    warehouse_uid: UUID | None = None
    page: int = 0
    size: int = 20
    sort: str = DEFAULT_INVENTORY_SORT


class WarehouseQueryOption(BaseModel):
    name: str | None = None
    page: int = 0
    size: int = 20
    sort: str = DEFAULT_INVENTORY_SORT


class TransferQueryOption(BaseModel):
    edition_uid: UUID | None = None
    source_warehouse_uid: UUID | None = None
    target_warehouse_uid: UUID | None = None
    status: TransferStatus | None = None
    page: int = 0
    size: int = 20
    sort: str = DEFAULT_INVENTORY_SORT


class AuditQueryOption(BaseModel):
    warehouse_uid: UUID | None = None
    status: AuditStatus | None = None
    page: int = 0
    size: int = 20
    sort: str = DEFAULT_INVENTORY_SORT


class ClosingQueryOption(BaseModel):
    warehouse_uid: UUID | None = None
    year: int | None = None
    month: int | None = None
    status: ClosingStatus | None = None
    page: int = 0
    size: int = 20
    sort: str = DEFAULT_INVENTORY_SORT


class BalanceQueryOption(BaseModel):
    edition_uid: UUID | None = None
    warehouse_uid: UUID | None = None
    year_month: YearMonth | None = None
    page: int = 0
    size: int = 20
    sort: str = DEFAULT_INVENTORY_SORT
