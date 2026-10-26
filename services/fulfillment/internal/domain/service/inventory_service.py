from uuid import UUID

from internal.application.port.event.event_publisher import EventPublisher
from internal.application.port.repo.inventory_repos import (
    ApplyChangeParams,
    AuditCreateParams,
    AuditItemParam,
    ClosingCreateParams,
    MonthlyClosingRepository,
    ReleaseParams,
    ReserveParams,
    StockAuditRepository,
    StockBalanceRepository,
    StockHistoryRepository,
    StockRepository,
    StockReservationRepository,
    StockTransferRepository,
    TransferCreateParams,
    WarehouseCreateParams,
    WarehouseRepository,
    WarehouseUpdateParams,
)
from internal.domain.aggregate.audit_aggregate import StockAuditAggregate
from internal.domain.aggregate.balance_aggregate import StockBalanceEntry
from internal.domain.aggregate.closing_aggregate import MonthlyClosingAggregate
from internal.domain.aggregate.stock_aggregate import StockAggregate, StockHistoryEntry
from internal.domain.aggregate.transfer_aggregate import StockTransferAggregate
from internal.domain.aggregate.warehouse_aggregate import WarehouseAggregate
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.command.inventory_command import (
    MonthlyClosingCommand,
    StockAdjustCommand,
    StockAuditCreateCommand,
    StockOrderReleaseCommand,
    StockOrderReserveCommand,
    StockReceiveCommand,
    StockReleaseCommand,
    StockReserveCommand,
    StockTransferCommand,
    WarehouseCreateCommand,
    WarehouseUpdateCommand,
)
from internal.domain.model.error import (
    InvalidQuantityError,
    InvalidTransitionError,
    SameWarehouseError,
    StockNotFoundError,
)
from internal.domain.model.option.inventory_option import (
    AuditQueryOption,
    BalanceQueryOption,
    ClosingQueryOption,
    StockHistoryQueryOption,
    StockQueryOption,
    TransferQueryOption,
    WarehouseQueryOption,
)


def reservation_key(order_ref: str, edition_id: int) -> str:
    return f"{order_ref}:{edition_id}"


class StockService:
    def __init__(
        self,
        stock_repo: StockRepository,
        stock_history_repo: StockHistoryRepository,
        reservation_repo: StockReservationRepository,
        event_publisher: EventPublisher,
    ):
        self._stock_repo = stock_repo
        self._stock_history_repo = stock_history_repo
        self._reservation_repo = reservation_repo
        self._events = event_publisher

    async def list_stocks(self, option: StockQueryOption) -> tuple[list[StockAggregate], int]:
        return await self._stock_repo.find_all(option)

    async def get_stock(self, edition_uid: UUID) -> StockAggregate | None:
        return await self._stock_repo.find_by_edition_uid(edition_uid)

    async def get_stock_history(
        self, edition_uid: UUID, option: StockHistoryQueryOption
    ) -> tuple[list[StockHistoryEntry], int]:
        stock_uids = await self._stock_repo.find_stock_uids(edition_uid, option.warehouse_uid)
        if not stock_uids:
            raise StockNotFoundError(f"stock not found for edition {edition_uid}")
        return await self._stock_history_repo.find_all(stock_uids, option)

    async def receive_stock(self, command: StockReceiveCommand) -> StockAggregate | None:
        stock = await self._apply_change(
            command.edition_uid, command.warehouse_uid, command.quantity, StockChangeType.INBOUND, command.reason
        )
        await self._publish_stock(stock, "received", command.edition_uid, command.warehouse_uid, command.quantity)
        return stock

    async def release_stock(self, command: StockReleaseCommand) -> StockAggregate | None:
        stock = await self._apply_change(
            command.edition_uid, command.warehouse_uid, command.quantity, StockChangeType.RELEASE, command.reason
        )
        await self._publish_stock(stock, "released", command.edition_uid, command.warehouse_uid, command.quantity)
        return stock

    async def adjust_stock(self, command: StockAdjustCommand) -> StockAggregate | None:
        stock = await self._apply_change(
            command.edition_uid, command.warehouse_uid, command.quantity, StockChangeType.ADJUSTMENT, command.reason
        )
        await self._publish_stock(stock, "adjusted", command.edition_uid, command.warehouse_uid, command.quantity)
        return stock

    async def reserve_stock(self, command: StockReserveCommand) -> StockAggregate | None:
        stock = await self._apply_change(
            command.edition_uid, command.warehouse_uid, -command.quantity, StockChangeType.RESERVATION, command.reason
        )
        await self._publish_stock(stock, "reserved", command.edition_uid, command.warehouse_uid, command.quantity)
        return stock

    async def reserve_stock_for_order(self, command: StockOrderReserveCommand) -> StockAggregate | None:
        if command.quantity <= 0:
            raise InvalidQuantityError(f"reserve quantity must be positive: got {command.quantity}")
        reservation = await self._reservation_repo.reserve(
            ReserveParams(
                reservation_key=reservation_key(command.order_uid, command.edition_id),
                order_ref=command.order_uid,
                edition_id=command.edition_id,
                quantity=command.quantity,
                reason=command.reason,
            )
        )
        return await self.get_stock(reservation.edition_uid)

    async def release_stock_for_order(self, command: StockOrderReleaseCommand) -> StockAggregate | None:
        reservation = await self._reservation_repo.release(
            ReleaseParams(
                reservation_key=reservation_key(command.order_uid, command.edition_id),
                reason=command.reason,
            )
        )
        if reservation is None:
            return None
        return await self.get_stock(reservation.edition_uid)

    async def _apply_change(
        self,
        edition_uid: UUID,
        warehouse_uid: UUID,
        delta: int,
        change_type: StockChangeType,
        reason: str | None,
    ) -> StockAggregate | None:
        await self._stock_repo.apply_change(
            ApplyChangeParams(
                edition_uid=edition_uid,
                warehouse_uid=warehouse_uid,
                delta=delta,
                change_type=change_type,
                reason=reason,
            )
        )
        return await self.get_stock(edition_uid)

    async def _publish_stock(
        self, stock: StockAggregate | None, action: str, edition_uid: UUID, warehouse_uid: UUID, quantity: int
    ) -> None:
        if stock is None:
            return
        await self._events.publish(
            "stock",
            str(stock.uid),
            action,
            {"edition_uid": str(edition_uid), "warehouse_uid": str(warehouse_uid), "quantity": quantity},
        )


class WarehouseService:
    def __init__(self, warehouse_repo: WarehouseRepository, event_publisher: EventPublisher):
        self._warehouse_repo = warehouse_repo
        self._events = event_publisher

    async def list_warehouses(self, option: WarehouseQueryOption) -> tuple[list[WarehouseAggregate], int]:
        return await self._warehouse_repo.find_all(option)

    async def get_warehouse(self, uid: UUID) -> WarehouseAggregate | None:
        return await self._warehouse_repo.find_by_uid(uid)

    async def register_warehouse(self, command: WarehouseCreateCommand) -> WarehouseAggregate:
        warehouse = await self._warehouse_repo.create(
            WarehouseCreateParams(name=command.name, address=command.address, capacity=command.capacity)
        )
        await self._events.publish("warehouse", str(warehouse.uid), "created", {"name": command.name})
        return warehouse

    async def update_warehouse(self, uid: UUID, command: WarehouseUpdateCommand) -> WarehouseAggregate | None:
        warehouse = await self._warehouse_repo.update(
            WarehouseUpdateParams(uid=uid, name=command.name, address=command.address, capacity=command.capacity)
        )
        if warehouse is None:
            return None
        await self._events.publish("warehouse", str(uid), "updated", {})
        return warehouse


class TransferService:
    def __init__(self, transfer_repo: StockTransferRepository, event_publisher: EventPublisher):
        self._transfer_repo = transfer_repo
        self._events = event_publisher

    async def transfer_stock(self, command: StockTransferCommand) -> StockTransferAggregate:
        if command.source_warehouse_uid == command.target_warehouse_uid:
            raise SameWarehouseError("transfer source and target warehouses must differ")
        transfer = await self._transfer_repo.create(
            TransferCreateParams(
                edition_uid=command.edition_uid,
                source_warehouse_uid=command.source_warehouse_uid,
                target_warehouse_uid=command.target_warehouse_uid,
                quantity=command.quantity,
                reason=command.reason,
            )
        )
        await self._events.publish(
            "transfer",
            str(transfer.uid),
            "created",
            {
                "edition_uid": str(command.edition_uid),
                "source_warehouse_uid": str(command.source_warehouse_uid),
                "target_warehouse_uid": str(command.target_warehouse_uid),
                "quantity": command.quantity,
            },
        )
        return transfer

    async def list_transfers(self, option: TransferQueryOption) -> tuple[list[StockTransferAggregate], int]:
        return await self._transfer_repo.find_all(option)

    async def get_transfer(self, uid: UUID) -> StockTransferAggregate | None:
        return await self._transfer_repo.find_by_uid(uid)

    async def complete_transfer(self, uid: UUID) -> StockTransferAggregate | None:
        transfer = await self._transfer_repo.find_by_uid(uid)
        if transfer is None:
            return None
        if not transfer.can_complete():
            raise InvalidTransitionError(f"transfer {uid} cannot be completed in status {transfer.status}")
        done = await self._transfer_repo.complete(uid)
        if done is None:
            return None
        await self._events.publish("transfer", str(uid), "completed", {})
        return done

    async def cancel_transfer(self, uid: UUID) -> StockTransferAggregate | None:
        transfer = await self._transfer_repo.find_by_uid(uid)
        if transfer is None:
            return None
        if not transfer.can_cancel():
            raise InvalidTransitionError(f"transfer {uid} cannot be cancelled in status {transfer.status}")
        done = await self._transfer_repo.cancel(uid)
        if done is None:
            return None
        await self._events.publish("transfer", str(uid), "cancelled", {})
        return done


class AuditService:
    def __init__(self, audit_repo: StockAuditRepository, event_publisher: EventPublisher):
        self._audit_repo = audit_repo
        self._events = event_publisher

    async def create_audit(self, command: StockAuditCreateCommand) -> StockAuditAggregate:
        audit = await self._audit_repo.create(
            AuditCreateParams(
                warehouse_uid=command.warehouse_uid,
                items=tuple(
                    AuditItemParam(edition_uid=item.edition_uid, actual_quantity=item.actual_quantity)
                    for item in command.items
                ),
                notes=command.notes,
            )
        )
        await self._events.publish(
            "audit",
            str(audit.uid),
            "created",
            {"warehouse_uid": str(command.warehouse_uid), "item_count": len(command.items)},
        )
        return audit

    async def list_audits(self, option: AuditQueryOption) -> tuple[list[StockAuditAggregate], int]:
        return await self._audit_repo.find_all(option)

    async def get_audit(self, uid: UUID) -> StockAuditAggregate | None:
        return await self._audit_repo.find_by_uid(uid)

    async def complete_audit(self, uid: UUID) -> StockAuditAggregate | None:
        audit = await self._audit_repo.find_by_uid(uid)
        if audit is None:
            return None
        if not audit.can_complete():
            raise InvalidTransitionError(f"audit {uid} is already completed")
        done = await self._audit_repo.complete(uid)
        if done is None:
            return None
        await self._events.publish("audit", str(uid), "completed", {})
        return done


class ClosingService:
    def __init__(self, closing_repo: MonthlyClosingRepository, event_publisher: EventPublisher):
        self._closing_repo = closing_repo
        self._events = event_publisher

    async def create_monthly_closing(self, command: MonthlyClosingCommand) -> MonthlyClosingAggregate:
        closing = await self._closing_repo.create(
            ClosingCreateParams(warehouse_uid=command.warehouse_uid, year=command.year, month=command.month)
        )
        await self._events.publish(
            "closing",
            str(closing.uid),
            "created",
            {"warehouse_uid": str(command.warehouse_uid), "year": command.year, "month": command.month},
        )
        return closing

    async def list_closings(self, option: ClosingQueryOption) -> tuple[list[MonthlyClosingAggregate], int]:
        return await self._closing_repo.find_all(option)

    async def get_closing(self, uid: UUID) -> MonthlyClosingAggregate | None:
        return await self._closing_repo.find_by_uid(uid)


class BalanceService:
    def __init__(self, balance_repo: StockBalanceRepository):
        self._balance_repo = balance_repo

    async def get_stock_balance(self, option: BalanceQueryOption) -> tuple[list[StockBalanceEntry], int]:
        return await self._balance_repo.find_all(option)
