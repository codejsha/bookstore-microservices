from datetime import datetime
from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.constant.freight_status import FreightStatus
from internal.domain.constant.shipment_status import ShipmentStatus


class ShipmentFilterOption(BaseModel):
    order_uid: UUID | None = None
    carrier_uid: UUID | None = None
    status: ShipmentStatus | None = None
    page: int = 0
    size: int = 20
    sort: str = "created_at:desc"


class CarrierFilterOption(BaseModel):
    name: str | None = None
    status: CarrierStatus | None = None
    page: int = 0
    size: int = 20
    sort: str = "name:asc"


class FreightFilterOption(BaseModel):
    shipment_uid: UUID | None = None
    carrier_uid: UUID | None = None
    status: FreightStatus | None = None
    page: int = 0
    size: int = 20
    sort: str = "created_at:desc"


class StatsFilterOption(BaseModel):
    carrier_uid: UUID | None = None
    date_from: datetime | None = None
    date_to: datetime | None = None
