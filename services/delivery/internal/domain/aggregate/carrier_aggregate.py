from datetime import datetime
from uuid import UUID

from internal.domain.constant.carrier_status import CarrierStatus


class CarrierAggregate:
    def __init__(
        self,
        uid: UUID,
        name: str,
        code: str,
        contact_name: str | None,
        contact_phone: str | None,
        contact_email: str | None,
        base_rate: float,
        rate_per_kg: float,
        status: CarrierStatus,
        created_at: datetime,
        updated_at: datetime,
    ):
        self.uid = uid
        self.name = name
        self.code = code
        self.contact_name = contact_name
        self.contact_phone = contact_phone
        self.contact_email = contact_email
        self.base_rate = base_rate
        self.rate_per_kg = rate_per_kg
        self.status = status
        self.created_at = created_at
        self.updated_at = updated_at
