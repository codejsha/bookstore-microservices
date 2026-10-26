from datetime import datetime
from typing import Annotated
from uuid import UUID

from pydantic import AfterValidator, BaseModel, Field, StringConstraints, model_validator

from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.constant.freight_status import FreightStatus
from internal.domain.constant.shipment_status import ShipmentStatus


def _require_non_blank(value: str) -> str:
    if not value.strip():
        raise ValueError("must not be blank")
    return value


def _non_blank(max_length: int) -> object:
    return Annotated[str, StringConstraints(max_length=max_length), AfterValidator(_require_non_blank)]


NonBlankStr = Annotated[str, AfterValidator(_require_non_blank)]
NonBlankStr3 = _non_blank(3)
NonBlankStr20 = _non_blank(20)
NonBlankStr50 = _non_blank(50)
NonBlankStr100 = _non_blank(100)
NonBlankStr255 = _non_blank(255)
NonBlankStr500 = _non_blank(500)
Str100 = Annotated[str, StringConstraints(max_length=100)]
CurrencyCode = Annotated[str, StringConstraints(pattern=r"^[A-Z]{3}$")]
EmailStr = Annotated[str, StringConstraints(max_length=255, pattern=r"^[^@\s]+@[^@\s]+\.[^@\s]+$")]

# ─── Shipment commands ──────────────────────────────────────────────────────


class CreateShipmentCommand(BaseModel):
    order_uid: UUID
    carrier_uid: UUID | None = None
    origin_address: NonBlankStr500
    destination_address: NonBlankStr500
    destination_city: NonBlankStr100
    destination_state: Str100
    destination_country_code: NonBlankStr3
    destination_postal_code: NonBlankStr20
    weight_kg: float | None = Field(None, gt=0)
    planned_pickup_at: datetime | None = None
    planned_delivery_at: datetime | None = None

    @model_validator(mode="after")
    def _check_planned_window(self) -> CreateShipmentCommand:
        if (
            self.planned_pickup_at is not None
            and self.planned_delivery_at is not None
            and self.planned_delivery_at < self.planned_pickup_at
        ):
            raise ValueError("planned_delivery_at must not precede planned_pickup_at")
        return self


class AssignCarrierCommand(BaseModel):
    carrier_uid: UUID
    tracking_number: NonBlankStr100 | None = None


class AddTrackingCommand(BaseModel):
    shipment_uid: UUID
    status: ShipmentStatus
    location: NonBlankStr255
    description: NonBlankStr


# ─── Carrier commands ───────────────────────────────────────────────────────


class CreateCarrierCommand(BaseModel):
    name: NonBlankStr255
    code: NonBlankStr50
    contact_name: NonBlankStr255 | None = None
    contact_phone: NonBlankStr50 | None = None
    contact_email: EmailStr | None = None
    base_rate: float = Field(ge=0)
    rate_per_kg: float = Field(ge=0)


class UpdateCarrierCommand(BaseModel):
    name: NonBlankStr255 | None = None
    contact_name: NonBlankStr255 | None = None
    contact_phone: NonBlankStr50 | None = None
    contact_email: EmailStr | None = None
    base_rate: float | None = Field(None, ge=0)
    rate_per_kg: float | None = Field(None, ge=0)
    status: CarrierStatus | None = None


# ─── Freight commands ───────────────────────────────────────────────────────


class CreateFreightCommand(BaseModel):
    shipment_uid: UUID
    carrier_uid: UUID
    base_cost: float = Field(ge=0)
    weight_surcharge: float = Field(0.0, ge=0)
    distance_surcharge: float = Field(0.0, ge=0)
    discount: float = Field(0.0, ge=0)
    currency: CurrencyCode = "KRW"

    @model_validator(mode="after")
    def _check_discount(self) -> CreateFreightCommand:
        if self.discount > self.base_cost + self.weight_surcharge + self.distance_surcharge:
            raise ValueError("discount must not exceed base cost plus surcharges")
        return self


class UpdateFreightStatusCommand(BaseModel):
    status: FreightStatus
