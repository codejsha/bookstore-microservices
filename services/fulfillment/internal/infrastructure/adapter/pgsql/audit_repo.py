from collections.abc import Sequence
from datetime import UTC, datetime
from uuid import UUID, uuid7

from sqlalchemy import func, select, update
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import AuditCreateParams, StockAuditRepository
from internal.domain.aggregate.audit_aggregate import StockAuditAggregate, StockAuditItem
from internal.domain.constant.audit_status import AuditStatus
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.error import InvalidTransitionError, StockNotFoundError, WarehouseNotFoundError
from internal.domain.model.option.inventory_option import AuditQueryOption
from internal.infrastructure.adapter.pgsql.models import (
    StockAuditEntity,
    StockAuditItemEntity,
    StockEntity,
    WarehouseEntity,
)
from internal.infrastructure.adapter.pgsql.optimistic_retry import VersionConflictError, run_with_optimistic_retry
from internal.infrastructure.adapter.pgsql.paging import page_offset, resolve_order_by
from internal.infrastructure.adapter.pgsql.stock_movement import StockDelta, apply_stock_delta

_SORTABLE = {
    "uid": StockAuditEntity.uid,
    "warehouse_uid": StockAuditEntity.warehouse_uid,
    "status": StockAuditEntity.status,
    "notes": StockAuditEntity.notes,
    "completed_at": StockAuditEntity.completed_at,
    "created_at": StockAuditEntity.created_at,
    "updated_at": StockAuditEntity.updated_at,
}


class PgsqlStockAuditRepository(StockAuditRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def find_all(self, option: AuditQueryOption) -> tuple[list[StockAuditAggregate], int]:
        conditions = [StockAuditEntity.deleted_at.is_(None)]
        if option.warehouse_uid is not None:
            conditions.append(StockAuditEntity.warehouse_uid == option.warehouse_uid)
        if option.status is not None:
            conditions.append(StockAuditEntity.status == option.status.value)
        async with self._session_factory() as session:
            total = (
                await session.execute(select(func.count()).select_from(StockAuditEntity).where(*conditions))
            ).scalar() or 0
            audits = list(
                (
                    await session.execute(
                        select(StockAuditEntity)
                        .where(*conditions)
                        .order_by(*resolve_order_by(option.sort, _SORTABLE))
                        .offset(page_offset(option.page, option.size))
                        .limit(option.size)
                    )
                ).scalars()
            )
            items_by_audit = await self._load_items(session, [a.id for a in audits])
            return [self._to_aggregate(a, items_by_audit.get(a.id, [])) for a in audits], total

    async def find_by_uid(self, uid: UUID) -> StockAuditAggregate | None:
        async with self._session_factory() as session:
            audit = (
                await session.execute(
                    select(StockAuditEntity).where(StockAuditEntity.uid == uid, StockAuditEntity.deleted_at.is_(None))
                )
            ).scalar_one_or_none()
            if audit is None:
                return None
            items = await self._load_items(session, [audit.id])
            return self._to_aggregate(audit, items.get(audit.id, []))

    async def create(self, params: AuditCreateParams) -> StockAuditAggregate:
        warehouse_uid = params.warehouse_uid
        async with self._session_factory() as session, session.begin():
            warehouse = (
                await session.execute(
                    select(WarehouseEntity).where(
                        WarehouseEntity.uid == warehouse_uid, WarehouseEntity.deleted_at.is_(None)
                    )
                )
            ).scalar_one_or_none()
            if warehouse is None:
                raise WarehouseNotFoundError(f"warehouse {params.warehouse_uid} not found")

            now = datetime.now(UTC)
            audit = StockAuditEntity(
                uid=uuid7(),
                warehouse_id=warehouse.id,
                warehouse_uid=warehouse.uid,
                status=AuditStatus.IN_PROGRESS.value,
                notes=params.notes,
                completed_at=None,
                created_at=now,
                updated_at=None,
                deleted_at=None,
                version=1,
            )
            session.add(audit)
            await session.flush()

            items: list[StockAuditItemEntity] = []
            for item in params.items:
                stock = (
                    await session.execute(
                        select(StockEntity).where(
                            StockEntity.edition_uid == item.edition_uid,
                            StockEntity.warehouse_uid == warehouse.uid,
                            StockEntity.deleted_at.is_(None),
                        )
                    )
                ).scalar_one_or_none()
                if stock is None:
                    raise StockNotFoundError(
                        f"no stock for edition {item.edition_uid} in warehouse {params.warehouse_uid}"
                    )
                entity = StockAuditItemEntity(
                    audit_id=audit.id,
                    edition_id=stock.edition_id,
                    edition_uid=stock.edition_uid,
                    system_quantity=stock.quantity,
                    actual_quantity=item.actual_quantity,
                    difference=item.actual_quantity - stock.quantity,
                )
                session.add(entity)
                items.append(entity)
            await session.flush()
            return self._to_aggregate(audit, items)

    async def complete(self, uid: UUID) -> StockAuditAggregate | None:
        async def attempt() -> StockAuditAggregate | None:
            async with self._session_factory() as session, session.begin():
                audit = (
                    await session.execute(
                        select(StockAuditEntity)
                        .where(StockAuditEntity.uid == uid, StockAuditEntity.deleted_at.is_(None))
                        .with_for_update()
                    )
                ).scalar_one_or_none()
                if audit is None:
                    return None
                if audit.status == AuditStatus.COMPLETED.value:
                    raise InvalidTransitionError(f"audit {uid} is already completed")

                items = (await self._load_items(session, [audit.id])).get(audit.id, [])
                reason = f"audit {uid}"
                for item in items:
                    if item.difference == 0:
                        continue
                    await apply_stock_delta(
                        session,
                        StockDelta(
                            edition_id=item.edition_id,
                            edition_uid=item.edition_uid,
                            warehouse_id=audit.warehouse_id,
                            warehouse_uid=audit.warehouse_uid,
                            delta=item.difference,
                            change_type=StockChangeType.ADJUSTMENT,
                            reason=reason,
                            allow_create=True,
                        ),
                    )

                now = datetime.now(UTC)
                result = await session.execute(
                    update(StockAuditEntity)
                    .where(StockAuditEntity.id == audit.id, StockAuditEntity.version == audit.version)
                    .values(
                        status=AuditStatus.COMPLETED.value,
                        completed_at=now,
                        updated_at=now,
                        version=audit.version + 1,
                    )
                )
                if result.rowcount == 0:
                    raise VersionConflictError(f"audit {uid} changed concurrently")
                return StockAuditAggregate(
                    uid=uid,
                    warehouse_uid=audit.warehouse_uid,
                    status=AuditStatus.COMPLETED,
                    notes=audit.notes,
                    completed_at=now,
                    created_at=audit.created_at,
                    updated_at=now,
                    items=[self._to_item(i) for i in items],
                )

        return await run_with_optimistic_retry(attempt)

    @staticmethod
    async def _load_items(session: AsyncSession, audit_ids: Sequence[int]) -> dict[int, list[StockAuditItemEntity]]:
        grouped: dict[int, list[StockAuditItemEntity]] = {}
        if not audit_ids:
            return grouped
        items = (
            await session.execute(
                select(StockAuditItemEntity)
                .where(StockAuditItemEntity.audit_id.in_(list(audit_ids)))
                .order_by(StockAuditItemEntity.audit_id, StockAuditItemEntity.id)
            )
        ).scalars()
        for item in items:
            grouped.setdefault(item.audit_id, []).append(item)
        return grouped

    @staticmethod
    def _to_item(entity: StockAuditItemEntity) -> StockAuditItem:
        return StockAuditItem(
            edition_uid=entity.edition_uid,
            system_quantity=entity.system_quantity,
            actual_quantity=entity.actual_quantity,
            difference=entity.difference,
        )

    @classmethod
    def _to_aggregate(cls, audit: StockAuditEntity, items: Sequence[StockAuditItemEntity]) -> StockAuditAggregate:
        return StockAuditAggregate(
            uid=audit.uid,
            warehouse_uid=audit.warehouse_uid,
            status=AuditStatus(audit.status),
            notes=audit.notes,
            completed_at=audit.completed_at,
            created_at=audit.created_at,
            updated_at=audit.updated_at,
            items=[cls._to_item(i) for i in items],
        )
