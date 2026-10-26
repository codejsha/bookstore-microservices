from abc import ABC, abstractmethod
from dataclasses import dataclass
from uuid import UUID

from internal.domain.aggregate.audit_aggregate import StockAuditAggregate
from internal.domain.aggregate.balance_aggregate import StockBalanceEntry
from internal.domain.aggregate.closing_aggregate import MonthlyClosingAggregate
from internal.domain.aggregate.reservation_aggregate import StockReservationAggregate
from internal.domain.aggregate.stock_aggregate import StockAggregate, StockHistoryEntry
from internal.domain.aggregate.transfer_aggregate import StockTransferAggregate
from internal.domain.aggregate.warehouse_aggregate import WarehouseAggregate
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.option.inventory_option import (
    AuditQueryOption,
    BalanceQueryOption,
    ClosingQueryOption,
    StockHistoryQueryOption,
    StockQueryOption,
    TransferQueryOption,
    WarehouseQueryOption,
)


@dataclass(frozen=True)
class ApplyChangeParams:
    edition_uid: UUID
    warehouse_uid: UUID
    delta: int
    change_type: StockChangeType
    reason: str | None


@dataclass(frozen=True)
class ReserveParams:
    reservation_key: str
    order_ref: str
    edition_id: int
    quantity: int
    reason: str | None


@dataclass(frozen=True)
class ReleaseParams:
    reservation_key: str
    reason: str | None


@dataclass(frozen=True)
class WarehouseCreateParams:
    name: str
    address: str | None
    capacity: int


@dataclass(frozen=True)
class WarehouseUpdateParams:
    uid: UUID
    name: str | None
    address: str | None
    capacity: int | None


@dataclass(frozen=True)
class TransferCreateParams:
    edition_uid: UUID
    source_warehouse_uid: UUID
    target_warehouse_uid: UUID
    quantity: int
    reason: str | None


@dataclass(frozen=True)
class AuditItemParam:
    edition_uid: UUID
    actual_quantity: int


@dataclass(frozen=True)
class AuditCreateParams:
    warehouse_uid: UUID
    items: tuple[AuditItemParam, ...]
    notes: str | None


@dataclass(frozen=True)
class ClosingCreateParams:
    warehouse_uid: UUID
    year: int
    month: int


class StockRepository(ABC):
    @abstractmethod
    async def find_all(self, option: StockQueryOption) -> tuple[list[StockAggregate], int]: ...

    @abstractmethod
    async def find_by_edition_uid(self, edition_uid: UUID) -> StockAggregate | None: ...

    @abstractmethod
    async def find_stock_uids(self, edition_uid: UUID, warehouse_uid: UUID | None) -> list[UUID]: ...

    @abstractmethod
    async def apply_change(self, params: ApplyChangeParams) -> None: ...


class StockHistoryRepository(ABC):
    @abstractmethod
    async def find_all(
        self, stock_uids: list[UUID], option: StockHistoryQueryOption
    ) -> tuple[list[StockHistoryEntry], int]: ...


class StockReservationRepository(ABC):
    @abstractmethod
    async def reserve(self, params: ReserveParams) -> StockReservationAggregate: ...

    @abstractmethod
    async def release(self, params: ReleaseParams) -> StockReservationAggregate | None: ...


class WarehouseRepository(ABC):
    @abstractmethod
    async def find_all(self, option: WarehouseQueryOption) -> tuple[list[WarehouseAggregate], int]: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> WarehouseAggregate | None: ...

    @abstractmethod
    async def create(self, params: WarehouseCreateParams) -> WarehouseAggregate: ...

    @abstractmethod
    async def update(self, params: WarehouseUpdateParams) -> WarehouseAggregate | None: ...


class StockTransferRepository(ABC):
    @abstractmethod
    async def find_all(self, option: TransferQueryOption) -> tuple[list[StockTransferAggregate], int]: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> StockTransferAggregate | None: ...

    @abstractmethod
    async def create(self, params: TransferCreateParams) -> StockTransferAggregate: ...

    @abstractmethod
    async def complete(self, uid: UUID) -> StockTransferAggregate | None: ...

    @abstractmethod
    async def cancel(self, uid: UUID) -> StockTransferAggregate | None: ...


class StockAuditRepository(ABC):
    @abstractmethod
    async def find_all(self, option: AuditQueryOption) -> tuple[list[StockAuditAggregate], int]: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> StockAuditAggregate | None: ...

    @abstractmethod
    async def create(self, params: AuditCreateParams) -> StockAuditAggregate: ...

    @abstractmethod
    async def complete(self, uid: UUID) -> StockAuditAggregate | None: ...


class MonthlyClosingRepository(ABC):
    @abstractmethod
    async def find_all(self, option: ClosingQueryOption) -> tuple[list[MonthlyClosingAggregate], int]: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> MonthlyClosingAggregate | None: ...

    @abstractmethod
    async def create(self, params: ClosingCreateParams) -> MonthlyClosingAggregate: ...


class StockBalanceRepository(ABC):
    @abstractmethod
    async def find_all(self, option: BalanceQueryOption) -> tuple[list[StockBalanceEntry], int]: ...
