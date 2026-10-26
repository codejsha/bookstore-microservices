from dataclasses import dataclass, field
from datetime import datetime
from uuid import UUID

from internal.domain.constant.stock_change_type import StockChangeType


@dataclass
class StockWarehouse:
    warehouse_uid: UUID
    warehouse_name: str
    quantity: int


@dataclass
class StockAggregate:
    uid: UUID
    edition_uid: UUID
    total_quantity: int
    warehouses: list[StockWarehouse] = field(default_factory=list)


@dataclass
class StockHistoryEntry:
    uid: UUID
    change_type: StockChangeType
    reason: str | None
    change_qty: int
    created_at: datetime
