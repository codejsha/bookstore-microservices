from abc import ABC, abstractmethod
from collections.abc import Callable
from uuid import UUID

from internal.domain.aggregate.carrier_aggregate import CarrierAggregate
from internal.domain.aggregate.freight_aggregate import FreightAggregate
from internal.domain.aggregate.shipment_aggregate import ShipmentAggregate
from internal.domain.aggregate.stats_aggregate import DeliveryDashboard
from internal.domain.aggregate.tracking_aggregate import TrackingAggregate
from internal.domain.model.option.delivery_option import (
    CarrierFilterOption,
    FreightFilterOption,
    ShipmentFilterOption,
    StatsFilterOption,
)


class ShipmentRepository(ABC):
    @abstractmethod
    async def save(self, shipment: ShipmentAggregate) -> ShipmentAggregate: ...

    @abstractmethod
    async def create_with_initial_tracking(
        self, shipment: ShipmentAggregate, tracking: TrackingAggregate
    ) -> ShipmentAggregate: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> ShipmentAggregate | None: ...

    @abstractmethod
    async def find_by_order_uid(self, order_uid: UUID) -> ShipmentAggregate | None: ...

    @abstractmethod
    async def find_all(self, option: ShipmentFilterOption) -> tuple[list[ShipmentAggregate], int]: ...

    @abstractmethod
    async def update(
        self,
        uid: UUID,
        mutator: Callable[[ShipmentAggregate], ShipmentAggregate],
    ) -> ShipmentAggregate | None: ...

    @abstractmethod
    async def transition(
        self,
        uid: UUID,
        guard: Callable[[ShipmentAggregate], bool],
        mutator: Callable[[ShipmentAggregate], ShipmentAggregate],
    ) -> ShipmentAggregate | None: ...


class TrackingRepository(ABC):
    @abstractmethod
    async def save(self, tracking: TrackingAggregate) -> TrackingAggregate: ...

    @abstractmethod
    async def find_by_shipment_uid(self, shipment_uid: UUID) -> list[TrackingAggregate]: ...


class CarrierRepository(ABC):
    @abstractmethod
    async def save(self, carrier: CarrierAggregate) -> CarrierAggregate: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> CarrierAggregate | None: ...

    @abstractmethod
    async def find_by_code(self, code: str) -> CarrierAggregate | None: ...

    @abstractmethod
    async def find_all(self, option: CarrierFilterOption) -> tuple[list[CarrierAggregate], int]: ...

    @abstractmethod
    async def update(
        self,
        uid: UUID,
        mutator: Callable[[CarrierAggregate], CarrierAggregate],
    ) -> CarrierAggregate | None: ...


class FreightRepository(ABC):
    @abstractmethod
    async def save(self, freight: FreightAggregate) -> FreightAggregate: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> FreightAggregate | None: ...

    @abstractmethod
    async def find_by_shipment_uid(self, shipment_uid: UUID) -> FreightAggregate | None: ...

    @abstractmethod
    async def find_all(self, option: FreightFilterOption) -> tuple[list[FreightAggregate], int]: ...

    @abstractmethod
    async def update(
        self,
        uid: UUID,
        mutator: Callable[[FreightAggregate], FreightAggregate],
    ) -> FreightAggregate | None: ...


class StatsRepository(ABC):
    @abstractmethod
    async def get_dashboard(self, option: StatsFilterOption) -> DeliveryDashboard: ...
