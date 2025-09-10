from dataclasses import dataclass


@dataclass
class ShipmentUpdateInput:
    userUid: str  # noqa: N815
    orderUid: str  # noqa: N815
    notificationType: str  # noqa: N815
    trackingNumber: str | None = None  # noqa: N815
