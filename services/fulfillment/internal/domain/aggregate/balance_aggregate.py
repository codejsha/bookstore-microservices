from dataclasses import dataclass
from uuid import UUID


@dataclass
class StockBalanceEntry:
    edition_uid: UUID
    warehouse_uid: UUID
    warehouse_name: str
    inbound_quantity: int
    outbound_quantity: int
    adjust_quantity: int
    current_quantity: int
