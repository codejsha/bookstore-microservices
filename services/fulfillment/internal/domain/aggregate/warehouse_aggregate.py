from dataclasses import dataclass
from datetime import datetime
from uuid import UUID


@dataclass
class WarehouseAggregate:
    uid: UUID
    name: str
    address: str | None
    capacity: int
    created_at: datetime
    updated_at: datetime | None
