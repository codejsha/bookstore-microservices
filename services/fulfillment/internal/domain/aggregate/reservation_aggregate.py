from dataclasses import dataclass
from uuid import UUID

from internal.domain.constant.reservation_status import ReservationStatus


@dataclass
class StockReservationAggregate:
    uid: UUID
    reservation_key: str
    edition_uid: UUID
    warehouse_uid: UUID
    quantity: int
    status: ReservationStatus
    already_applied: bool
