from datetime import datetime
from uuid import UUID

from internal.domain.constant.freight_status import FreightStatus


class FreightAggregate:
    def __init__(
        self,
        uid: UUID,
        shipment_uid: UUID,
        carrier_uid: UUID,
        base_cost: float,
        weight_surcharge: float,
        distance_surcharge: float,
        discount: float,
        total_cost: float,
        currency: str,
        status: FreightStatus,
        invoiced_at: datetime | None,
        paid_at: datetime | None,
        created_at: datetime,
        updated_at: datetime,
    ):
        self.uid = uid
        self.shipment_uid = shipment_uid
        self.carrier_uid = carrier_uid
        self.base_cost = base_cost
        self.weight_surcharge = weight_surcharge
        self.distance_surcharge = distance_surcharge
        self.discount = discount
        self.total_cost = total_cost
        self.currency = currency
        self.status = status
        self.invoiced_at = invoiced_at
        self.paid_at = paid_at
        self.created_at = created_at
        self.updated_at = updated_at
