from dataclasses import dataclass, field
from datetime import datetime
from uuid import UUID

from internal.domain.constant.closing_status import ClosingStatus


@dataclass
class MonthlyClosingItem:
    edition_uid: UUID
    opening_quantity: int
    inbound_quantity: int
    outbound_quantity: int
    adjust_quantity: int
    closing_quantity: int


@dataclass
class MonthlyClosingAggregate:
    uid: UUID
    warehouse_uid: UUID
    year: int
    month: int
    status: ClosingStatus
    closed_at: datetime | None
    created_at: datetime
    updated_at: datetime | None
    items: list[MonthlyClosingItem] = field(default_factory=list)
