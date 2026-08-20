from datetime import UTC, datetime
from uuid import UUID, uuid7

from internal.application.port.repo.repos import (
    CarrierRepository,
    FreightRepository,
    ShipmentRepository,
    StatsRepository,
    TrackingRepository,
)
from internal.domain.aggregate.carrier_aggregate import CarrierAggregate
from internal.domain.aggregate.freight_aggregate import FreightAggregate
from internal.domain.aggregate.shipment_aggregate import ShipmentAggregate
from internal.domain.aggregate.stats_aggregate import DeliveryDashboard
from internal.domain.aggregate.tracking_aggregate import TrackingAggregate
from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.constant.freight_status import FreightStatus
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.command.delivery_command import (
    AddTrackingCommand,
    AssignCarrierCommand,
    CreateCarrierCommand,
    CreateFreightCommand,
    CreateShipmentCommand,
    UpdateCarrierCommand,
    UpdateFreightStatusCommand,
)
from internal.domain.model.option.delivery_option import (
    CarrierFilterOption,
    FreightFilterOption,
    ShipmentFilterOption,
    StatsFilterOption,
)


class ShipmentService:
    def __init__(
        self,
        shipment_repo: ShipmentRepository,
        tracking_repo: TrackingRepository,
    ):
        self._shipment_repo = shipment_repo
        self._tracking_repo = tracking_repo

    async def create_shipment(self, command: CreateShipmentCommand) -> ShipmentAggregate:
        now = datetime.now(UTC)
        shipment = ShipmentAggregate(
            uid=uuid7(),
            order_uid=command.order_uid,
            carrier_uid=command.carrier_uid,
            origin_address=command.origin_address,
            destination_address=command.destination_address,
            destination_city=command.destination_city,
            destination_state=command.destination_state,
            destination_country_code=command.destination_country_code,
            destination_postal_code=command.destination_postal_code,
            status=ShipmentStatus.PLANNED,
            tracking_number=None,
            weight_kg=command.weight_kg,
            planned_pickup_at=command.planned_pickup_at,
            planned_delivery_at=command.planned_delivery_at,
            actual_pickup_at=None,
            actual_delivery_at=None,
            created_at=now,
            updated_at=now,
        )
        initial_tracking = TrackingAggregate(
            uid=uuid7(),
            shipment_uid=shipment.uid,
            status=ShipmentStatus.PLANNED,
            location="",
            description="Shipment planned",
            occurred_at=now,
            created_at=now,
        )
        return await self._shipment_repo.create_with_initial_tracking(shipment, initial_tracking)

    async def get_shipment(self, uid: UUID) -> ShipmentAggregate | None:
        return await self._shipment_repo.find_by_uid(uid)

    async def get_shipment_by_order_uid(self, order_uid: UUID) -> ShipmentAggregate | None:
        return await self._shipment_repo.find_by_order_uid(order_uid)

    async def list_shipments(self, option: ShipmentFilterOption) -> tuple[list[ShipmentAggregate], int]:
        return await self._shipment_repo.find_all(option)

    async def set_tracking_number(self, shipment: ShipmentAggregate) -> ShipmentAggregate:
        shipment.updated_at = datetime.now(UTC)
        return await self._shipment_repo.save(shipment)

    async def assign_carrier(self, uid: UUID, command: AssignCarrierCommand) -> ShipmentAggregate | None:
        now = datetime.now(UTC)

        def _mutate(s: ShipmentAggregate) -> ShipmentAggregate:
            s.carrier_uid = command.carrier_uid
            s.tracking_number = command.tracking_number
            s.updated_at = now
            return s

        return await self._shipment_repo.update(uid, _mutate)

    async def dispatch_shipment(self, uid: UUID) -> ShipmentAggregate | None:
        def _mutate(s: ShipmentAggregate) -> ShipmentAggregate:
            s.status = ShipmentStatus.DISPATCHED
            s.updated_at = datetime.now(UTC)
            return s

        shipment = await self._shipment_repo.transition(uid, ShipmentAggregate.can_dispatch, _mutate)
        if shipment is None:
            return None
        await self._add_tracking(shipment.uid, ShipmentStatus.DISPATCHED, "", "Shipment dispatched")
        return shipment

    async def pick_up_shipment(self, uid: UUID) -> ShipmentAggregate | None:
        def _mutate(s: ShipmentAggregate) -> ShipmentAggregate:
            now = datetime.now(UTC)
            s.status = ShipmentStatus.PICKED_UP
            s.actual_pickup_at = now
            s.updated_at = now
            return s

        shipment = await self._shipment_repo.transition(uid, ShipmentAggregate.can_pick_up, _mutate)
        if shipment is None:
            return None
        await self._add_tracking(shipment.uid, ShipmentStatus.PICKED_UP, "", "Package picked up")
        return shipment

    async def deliver_shipment(self, uid: UUID) -> ShipmentAggregate | None:
        def _mutate(s: ShipmentAggregate) -> ShipmentAggregate:
            now = datetime.now(UTC)
            s.status = ShipmentStatus.DELIVERED
            s.actual_delivery_at = now
            s.updated_at = now
            return s

        shipment = await self._shipment_repo.transition(uid, ShipmentAggregate.can_deliver, _mutate)
        if shipment is None:
            return None
        await self._add_tracking(shipment.uid, ShipmentStatus.DELIVERED, "", "Package delivered")
        return shipment

    async def cancel_shipment(self, uid: UUID) -> ShipmentAggregate | None:
        def _mutate(s: ShipmentAggregate) -> ShipmentAggregate:
            s.status = ShipmentStatus.CANCELLED
            s.updated_at = datetime.now(UTC)
            return s

        shipment = await self._shipment_repo.transition(uid, ShipmentAggregate.can_cancel, _mutate)
        if shipment is None:
            return None
        await self._add_tracking(shipment.uid, ShipmentStatus.CANCELLED, "", "Shipment cancelled")
        return shipment

    async def add_tracking(self, command: AddTrackingCommand) -> TrackingAggregate:
        shipment = await self._shipment_repo.find_by_uid(command.shipment_uid)
        if shipment is None:
            raise ValueError(f"Shipment {command.shipment_uid} not found")
        return await self._add_tracking(command.shipment_uid, command.status, command.location, command.description)

    async def get_tracking_history(self, shipment_uid: UUID) -> list[TrackingAggregate]:
        return await self._tracking_repo.find_by_shipment_uid(shipment_uid)

    async def _add_tracking(
        self, shipment_uid: UUID, status: ShipmentStatus, location: str, description: str
    ) -> TrackingAggregate:
        now = datetime.now(UTC)
        tracking = TrackingAggregate(
            uid=uuid7(),
            shipment_uid=shipment_uid,
            status=status,
            location=location,
            description=description,
            occurred_at=now,
            created_at=now,
        )
        return await self._tracking_repo.save(tracking)


class CarrierService:
    def __init__(self, carrier_repo: CarrierRepository):
        self._carrier_repo = carrier_repo

    async def create_carrier(self, command: CreateCarrierCommand) -> CarrierAggregate:
        now = datetime.now(UTC)
        carrier = CarrierAggregate(
            uid=uuid7(),
            name=command.name,
            code=command.code,
            contact_name=command.contact_name,
            contact_phone=command.contact_phone,
            contact_email=command.contact_email,
            base_rate=command.base_rate,
            rate_per_kg=command.rate_per_kg,
            status=CarrierStatus.ACTIVE,
            created_at=now,
            updated_at=now,
        )
        return await self._carrier_repo.save(carrier)

    async def get_carrier(self, uid: UUID) -> CarrierAggregate | None:
        return await self._carrier_repo.find_by_uid(uid)

    async def list_carriers(self, option: CarrierFilterOption) -> tuple[list[CarrierAggregate], int]:
        return await self._carrier_repo.find_all(option)

    async def update_carrier(self, uid: UUID, command: UpdateCarrierCommand) -> CarrierAggregate | None:
        now = datetime.now(UTC)

        def _mutate(carrier: CarrierAggregate) -> CarrierAggregate:
            if command.name is not None:
                carrier.name = command.name
            if command.contact_name is not None:
                carrier.contact_name = command.contact_name
            if command.contact_phone is not None:
                carrier.contact_phone = command.contact_phone
            if command.contact_email is not None:
                carrier.contact_email = command.contact_email
            if command.base_rate is not None:
                carrier.base_rate = command.base_rate
            if command.rate_per_kg is not None:
                carrier.rate_per_kg = command.rate_per_kg
            if command.status is not None:
                carrier.status = command.status
            carrier.updated_at = now
            return carrier

        return await self._carrier_repo.update(uid, _mutate)


class FreightService:
    def __init__(self, freight_repo: FreightRepository):
        self._freight_repo = freight_repo

    async def create_freight(self, command: CreateFreightCommand) -> FreightAggregate:
        now = datetime.now(UTC)
        total = command.base_cost + command.weight_surcharge + command.distance_surcharge - command.discount
        freight = FreightAggregate(
            uid=uuid7(),
            shipment_uid=command.shipment_uid,
            carrier_uid=command.carrier_uid,
            base_cost=command.base_cost,
            weight_surcharge=command.weight_surcharge,
            distance_surcharge=command.distance_surcharge,
            discount=command.discount,
            total_cost=total,
            currency=command.currency,
            status=FreightStatus.ESTIMATED,
            invoiced_at=None,
            paid_at=None,
            created_at=now,
            updated_at=now,
        )
        return await self._freight_repo.save(freight)

    async def get_freight(self, uid: UUID) -> FreightAggregate | None:
        return await self._freight_repo.find_by_uid(uid)

    async def list_freights(self, option: FreightFilterOption) -> tuple[list[FreightAggregate], int]:
        return await self._freight_repo.find_all(option)

    async def find_by_shipment(self, shipment_uid: UUID) -> FreightAggregate | None:
        return await self._freight_repo.find_by_shipment_uid(shipment_uid)

    async def update_status(self, uid: UUID, command: UpdateFreightStatusCommand) -> FreightAggregate | None:
        now = datetime.now(UTC)

        def _mutate(freight: FreightAggregate) -> FreightAggregate:
            freight.status = command.status
            if command.status == FreightStatus.INVOICED:
                freight.invoiced_at = now
            elif command.status == FreightStatus.PAID:
                freight.paid_at = now
            freight.updated_at = now
            return freight

        return await self._freight_repo.update(uid, _mutate)


class StatsService:
    def __init__(self, stats_repo: StatsRepository):
        self._stats_repo = stats_repo

    async def get_dashboard(self, option: StatsFilterOption) -> DeliveryDashboard:
        return await self._stats_repo.get_dashboard(option)
