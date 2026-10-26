from dataclasses import dataclass, field
from datetime import datetime
from uuid import UUID

from internal.domain.constant.audit_status import AuditStatus


@dataclass
class StockAuditItem:
    edition_uid: UUID
    system_quantity: int
    actual_quantity: int
    difference: int


@dataclass
class StockAuditAggregate:
    uid: UUID
    warehouse_uid: UUID
    status: AuditStatus
    notes: str | None
    completed_at: datetime | None
    created_at: datetime
    updated_at: datetime | None
    items: list[StockAuditItem] = field(default_factory=list)

    def can_complete(self) -> bool:
        return self.status != AuditStatus.COMPLETED
