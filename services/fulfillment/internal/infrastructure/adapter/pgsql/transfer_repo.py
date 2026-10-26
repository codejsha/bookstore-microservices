from collections.abc import Awaitable, Callable
from datetime import UTC, datetime
from uuid import UUID, uuid7

from sqlalchemy import func, select, update
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import StockTransferRepository, TransferCreateParams
from internal.domain.aggregate.transfer_aggregate import StockTransferAggregate
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.constant.transfer_status import TransferStatus
from internal.domain.model.error import InvalidTransitionError, StockNotFoundError, WarehouseNotFoundError
from internal.infrastructure.adapter.mysql.models import StockTransferEntity, WarehouseEntity
from internal.infrastructure.adapter.mysql.optimistic_retry import VersionConflictError, run_with_optimistic_retry
from internal.infrastructure.adapter.mysql.paging import page_offset, resolve_order_by
from internal.infrastructure.adapter.mysql.stock_movement import StockDelta, apply_stock_delta, lock_stock
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes

_SORTABLE = {
    "uid": StockTransferEntity.uid,
    "edition_uid": StockTransferEntity.edition_uid,
    "source_warehouse_uid": StockTransferEntity.source_warehouse_uid,
    "target_warehouse_uid": StockTransferEntity.target_warehouse_uid,
    "quantity": StockTransferEntity.quantity,
    "status": StockTransferEntity.status,
    "completed_at": StockTransferEntity.completed_at,
    "created_at": StockTransferEntity.created_at,
    "updated_at": StockTransferEntity.updated_at,
}


class MySQLStockTransferRepository(StockTransferRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def find_all(self, option) -> tuple[list[StockTransferAggregate], int]:
        conditions = [StockTransferEntity.deleted_at.is_(None)]
        if option.edition_uid is not None:
            conditions.append(StockTransferEntity.edition_uid == uuid_to_bytes(option.edition_uid))
        if option.source_warehouse_uid is not None:
            conditions.append(StockTransferEntity.source_warehouse_uid == uuid_to_bytes(option.source_warehouse_uid))
        if option.target_warehouse_uid is not None:
            conditions.append(StockTransferEntity.target_warehouse_uid == uuid_to_bytes(option.target_warehouse_uid))
        if option.status is not None:
            conditions.append(StockTransferEntity.status == option.status.value)
        async with self._session_factory() as session:
            total = (
                await session.execute(select(func.count()).select_from(StockTransferEntity).where(*conditions))
            ).scalar() or 0
            entities = (
                await session.execute(
                    select(StockTransferEntity)
                    .where(*conditions)
                    .order_by(*resolve_order_by(option.sort, _SORTABLE))
                    .offset(page_offset(option.page, option.size))
                    .limit(option.size)
                )
            ).scalars()
            return [self._to_aggregate(e) for e in entities], total

    async def find_by_uid(self, uid: UUID) -> StockTransferAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(StockTransferEntity).where(
                        StockTransferEntity.uid == uuid_to_bytes(uid), StockTransferEntity.deleted_at.is_(None)
                    )
                )
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def create(self, params: TransferCreateParams) -> StockTransferAggregate:
        edition_uid = uuid_to_bytes(params.edition_uid)
        source_uid = uuid_to_bytes(params.source_warehouse_uid)
        target_uid = uuid_to_bytes(params.target_warehouse_uid)

        async def attempt() -> StockTransferAggregate:
            async with self._session_factory() as session, session.begin():
                source = await lock_stock(session, edition_uid, source_uid)
                if source is None:
                    raise StockNotFoundError(
                        f"no stock for edition {params.edition_uid} in source warehouse {params.source_warehouse_uid}"
                    )
                target = (
                    await session.execute(
                        select(WarehouseEntity).where(
                            WarehouseEntity.uid == target_uid, WarehouseEntity.deleted_at.is_(None)
                        )
                    )
                ).scalar_one_or_none()
                if target is None:
                    raise WarehouseNotFoundError(f"target warehouse {params.target_warehouse_uid} not found")
                edition_id = source.edition_id
                source_warehouse_id = source.warehouse_id
                await apply_stock_delta(
                    session,
                    StockDelta(
                        edition_id=edition_id,
                        edition_uid=edition_uid,
                        warehouse_id=source_warehouse_id,
                        warehouse_uid=source_uid,
                        delta=-params.quantity,
                        change_type=StockChangeType.OUTBOUND,
                        reason=params.reason,
                    ),
                    existing=source,
                )
                now = datetime.now(UTC)
                transfer = StockTransferEntity(
                    uid=uuid7().bytes,
                    edition_id=edition_id,
                    edition_uid=edition_uid,
                    source_warehouse_id=source_warehouse_id,
                    source_warehouse_uid=source_uid,
                    target_warehouse_id=target.id,
                    target_warehouse_uid=target.uid,
                    quantity=params.quantity,
                    status=TransferStatus.PENDING.value,
                    reason=params.reason,
                    completed_at=None,
                    created_at=now,
                    updated_at=None,
                    deleted_at=None,
                    version=1,
                )
                session.add(transfer)
                await session.flush()
                return self._to_aggregate(transfer)

        return await run_with_optimistic_retry(attempt)

    async def complete(self, uid: UUID) -> StockTransferAggregate | None:
        async def apply(session: AsyncSession, transfer: StockTransferEntity) -> tuple[TransferStatus, datetime | None]:
            if transfer.status not in (TransferStatus.PENDING.value, TransferStatus.IN_TRANSIT.value):
                raise InvalidTransitionError(f"transfer {uid} cannot be completed in status {transfer.status}")
            await apply_stock_delta(
                session,
                StockDelta(
                    edition_id=transfer.edition_id,
                    edition_uid=transfer.edition_uid,
                    warehouse_id=transfer.target_warehouse_id,
                    warehouse_uid=transfer.target_warehouse_uid,
                    delta=transfer.quantity,
                    change_type=StockChangeType.INBOUND,
                    reason=transfer.reason,
                    allow_create=True,
                ),
            )
            return TransferStatus.COMPLETED, datetime.now(UTC)

        return await self._finalize(uid, apply)

    async def cancel(self, uid: UUID) -> StockTransferAggregate | None:
        async def apply(session: AsyncSession, transfer: StockTransferEntity) -> tuple[TransferStatus, datetime | None]:
            if transfer.status != TransferStatus.PENDING.value:
                raise InvalidTransitionError(f"transfer {uid} cannot be cancelled in status {transfer.status}")
            await apply_stock_delta(
                session,
                StockDelta(
                    edition_id=transfer.edition_id,
                    edition_uid=transfer.edition_uid,
                    warehouse_id=transfer.source_warehouse_id,
                    warehouse_uid=transfer.source_warehouse_uid,
                    delta=transfer.quantity,
                    change_type=StockChangeType.INBOUND,
                    reason=transfer.reason,
                    allow_create=True,
                ),
            )
            return TransferStatus.CANCELLED, None

        return await self._finalize(uid, apply)

    async def _finalize(
        self,
        uid: UUID,
        apply: Callable[[AsyncSession, StockTransferEntity], Awaitable[tuple[TransferStatus, datetime | None]]],
    ) -> StockTransferAggregate | None:
        async def attempt() -> StockTransferAggregate | None:
            async with self._session_factory() as session, session.begin():
                transfer = (
                    await session.execute(
                        select(StockTransferEntity)
                        .where(StockTransferEntity.uid == uuid_to_bytes(uid), StockTransferEntity.deleted_at.is_(None))
                        .with_for_update()
                    )
                ).scalar_one_or_none()
                if transfer is None:
                    return None
                status, completed_at = await apply(session, transfer)
                now = datetime.now(UTC)
                result = await session.execute(
                    update(StockTransferEntity)
                    .where(StockTransferEntity.id == transfer.id, StockTransferEntity.version == transfer.version)
                    .values(
                        status=status.value, completed_at=completed_at, updated_at=now, version=transfer.version + 1
                    )
                )
                if result.rowcount == 0:
                    raise VersionConflictError(f"transfer {uid} changed concurrently")
                return StockTransferAggregate(
                    uid=uid,
                    edition_uid=bytes_to_uuid(transfer.edition_uid),
                    source_warehouse_uid=bytes_to_uuid(transfer.source_warehouse_uid),
                    target_warehouse_uid=bytes_to_uuid(transfer.target_warehouse_uid),
                    quantity=transfer.quantity,
                    status=status,
                    reason=transfer.reason,
                    completed_at=completed_at,
                    created_at=transfer.created_at,
                    updated_at=now,
                )

        return await run_with_optimistic_retry(attempt)

    @staticmethod
    def _to_aggregate(entity: StockTransferEntity) -> StockTransferAggregate:
        return StockTransferAggregate(
            uid=bytes_to_uuid(entity.uid),
            edition_uid=bytes_to_uuid(entity.edition_uid),
            source_warehouse_uid=bytes_to_uuid(entity.source_warehouse_uid),
            target_warehouse_uid=bytes_to_uuid(entity.target_warehouse_uid),
            quantity=entity.quantity,
            status=TransferStatus(entity.status),
            reason=entity.reason,
            completed_at=entity.completed_at,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )
