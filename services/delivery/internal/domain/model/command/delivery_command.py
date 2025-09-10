from datetime import datetime
from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.constant.freight_status import FreightStatus
from internal.domain.constant.shipment_status import ShipmentStatus

# ─── Shipment commands ──────────────────────────────────────────────────────


class CreateShipmentCommand(BaseModel):
    order_uid: UUID
    carrier_uid: UUID | None = None
    origin_address: str
    destination_address: str
    destination_city: str
    destination_state: str
    destination_country_code: str
    destination_postal_code: str
    weight_kg: float | None = None
    planned_pickup_at: datetime | None = None
    planned_delivery_at: datetime | None = None


class AssignCarrierCommand(BaseModel):
    carrier_uid: UUID
    tracking_number: str | None = None


class AddTrackingCommand(BaseModel):
    shipment_uid: UUID
    status: ShipmentStatus
    location: str
    description: str


# ─── Carrier commands ───────────────────────────────────────────────────────


class CreateCarrierCommand(BaseModel):
    name: str
    code: str
    contact_name: str | None = None
    contact_phone: str | None = None
    contact_email: str | None = None
    base_rate: float
    rate_per_kg: float


class UpdateCarrierCommand(BaseModel):
    name: str | None = None
    contact_name: str | None = None
    contact_phone: str | None = None
    contact_email: str | None = None
    base_rate: float | None = None
    rate_per_kg: float | None = None
    status: CarrierStatus | None = None


# ─── Freight commands ───────────────────────────────────────────────────────


class CreateFreightCommand(BaseModel):
    shipment_uid: UUID
    carrier_uid: UUID
    base_cost: float
    weight_surcharge: float = 0.0
    distance_surcharge: float = 0.0
    discount: float = 0.0
    currency: str = "KRW"


class UpdateFreightStatusCommand(BaseModel):
    status: FreightStatus
