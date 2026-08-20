from dataclasses import dataclass
from uuid import UUID, uuid4

import structlog
from sqlalchemy.exc import IntegrityError
from temporalio import activity

from internal.domain.aggregate.shipment_aggregate import ShipmentAggregate
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.command.delivery_command import CreateShipmentCommand
from internal.domain.service.delivery_service import ShipmentService

logger = structlog.get_logger()


@dataclass
class CreateShipmentInput:
    orderUid: str  # noqa: N815
    originAddress: str  # noqa: N815
    destinationAddress: str  # noqa: N815
    destinationCity: str  # noqa: N815
    destinationState: str  # noqa: N815
    destinationCountryCode: str  # noqa: N815
    destinationPostalCode: str  # noqa: N815


@dataclass
class CreateShipmentOutput:
    shipmentUid: str  # noqa: N815
    trackingNumber: str | None  # noqa: N815
    status: str  # noqa: N815


def _generate_tracking_number() -> str:
    return f"DLV{uuid4().hex[:16].upper()}"


class ShipmentActivities:
    def __init__(self, shipment_service: ShipmentService):
        self._shipment_service = shipment_service

    @activity.defn(name="CreateShipment")
    async def create_shipment(self, payload: CreateShipmentInput) -> CreateShipmentOutput:
        order_uid = UUID(payload.orderUid)
        logger.info("CreateShipment activity invoked", order_uid=payload.orderUid)

        command = CreateShipmentCommand(
            order_uid=order_uid,
            origin_address=payload.originAddress,
            destination_address=payload.destinationAddress,
            destination_city=payload.destinationCity,
            destination_state=payload.destinationState,
            destination_country_code=payload.destinationCountryCode,
            destination_postal_code=payload.destinationPostalCode,
        )

        shipment = await self._get_or_create_shipment(order_uid, command)

        if shipment.status == ShipmentStatus.PLANNED:
            dispatched = await self._shipment_service.dispatch_shipment(shipment.uid)
            if dispatched is None:
                raise RuntimeError(f"failed to dispatch shipment {shipment.uid}")
            shipment = dispatched

        if not shipment.tracking_number:
            shipment.tracking_number = _generate_tracking_number()
            shipment = await self._shipment_service.set_tracking_number(shipment)

        logger.info(
            "CreateShipment activity completed",
            shipment_uid=str(shipment.uid),
            tracking_number=shipment.tracking_number,
            status=shipment.status.value,
        )
        return CreateShipmentOutput(
            shipmentUid=str(shipment.uid),
            trackingNumber=shipment.tracking_number,
            status=shipment.status.value,
        )

    async def _get_or_create_shipment(self, order_uid: UUID, command: CreateShipmentCommand) -> ShipmentAggregate:
        existing = await self._shipment_service.get_shipment_by_order_uid(order_uid)
        if existing is not None:
            logger.info(
                "CreateShipment idempotent hit; reusing existing shipment",
                order_uid=str(order_uid),
                shipment_uid=str(existing.uid),
                status=existing.status.value,
            )
            return existing
        try:
            return await self._shipment_service.create_shipment(command)
        except IntegrityError:
            existing = await self._shipment_service.get_shipment_by_order_uid(order_uid)
            if existing is None:
                raise
            logger.info(
                "CreateShipment lost insert race; reusing existing shipment",
                order_uid=str(order_uid),
                shipment_uid=str(existing.uid),
            )
            return existing
