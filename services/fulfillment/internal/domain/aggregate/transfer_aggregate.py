from dataclasses import dataclass
from datetime import datetime
from uuid import UUID

from internal.domain.constant.transfer_status import TransferStatus


@dataclass
class StockTransferAggregate:
    uid: UUID
    edition_uid: UUID
    source_warehouse_uid: UUID
    target_warehouse_uid: UUID
    quantity: int
    status: TransferStatus
    reason: str | None
    completed_at: datetime | None
    created_at: datetime
    updated_at: datetime | None

    def can_complete(self) -> bool:
        return self.status in (TransferStatus.PENDING, TransferStatus.IN_TRANSIT)

    def can_cancel(self) -> bool:
        return self.status == TransferStatus.PENDING
