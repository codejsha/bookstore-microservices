from datetime import datetime
from uuid import UUID

from internal.domain.constant.shipment_status import ShipmentStatus


class TrackingAggregate:
    def __init__(
        self,
        uid: UUID,
        shipment_uid: UUID,
        status: ShipmentStatus,
        location: str,
        description: str,
        occurred_at: datetime,
        created_at: datetime,
    ):
        self.uid = uid
        self.shipment_uid = shipment_uid
        self.status = status
        self.location = location
        self.description = description
        self.occurred_at = occurred_at
        self.created_at = created_at
