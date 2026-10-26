from datetime import datetime
from uuid import UUID

from internal.domain.constant.shipment_status import ShipmentStatus


class ShipmentAggregate:
    def __init__(
        self,
        uid: UUID,
        order_uid: UUID,
        carrier_uid: UUID | None,
        origin_address: str,
        destination_address: str,
        destination_city: str,
        destination_state: str,
        destination_country_code: str,
        destination_postal_code: str,
        status: ShipmentStatus,
        tracking_number: str | None,
        weight_kg: float | None,
        planned_pickup_at: datetime | None,
        planned_delivery_at: datetime | None,
        actual_pickup_at: datetime | None,
        actual_delivery_at: datetime | None,
        created_at: datetime,
        updated_at: datetime,
    ):
        self.uid = uid
        self.order_uid = order_uid
        self.carrier_uid = carrier_uid
        self.origin_address = origin_address
        self.destination_address = destination_address
        self.destination_city = destination_city
        self.destination_state = destination_state
        self.destination_country_code = destination_country_code
        self.destination_postal_code = destination_postal_code
        self.status = status
        self.tracking_number = tracking_number
        self.weight_kg = weight_kg
        self.planned_pickup_at = planned_pickup_at
        self.planned_delivery_at = planned_delivery_at
        self.actual_pickup_at = actual_pickup_at
        self.actual_delivery_at = actual_delivery_at
        self.created_at = created_at
        self.updated_at = updated_at

    def can_dispatch(self) -> bool:
        return self.status == ShipmentStatus.PLANNED

    def can_pick_up(self) -> bool:
        return self.status == ShipmentStatus.DISPATCHED

    def can_deliver(self) -> bool:
        return self.status in (
            ShipmentStatus.PICKED_UP,
            ShipmentStatus.IN_TRANSIT,
            ShipmentStatus.OUT_FOR_DELIVERY,
        )

    def can_cancel(self) -> bool:
        return self.status in (ShipmentStatus.PLANNED, ShipmentStatus.DISPATCHED)
