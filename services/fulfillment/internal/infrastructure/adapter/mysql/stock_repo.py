from collections.abc import Sequence
from uuid import UUID

from sqlalchemy import Row, distinct, func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import ApplyChangeParams, StockRepository
from internal.domain.aggregate.stock_aggregate import StockAggregate, StockWarehouse
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.error import StockNotFoundError, WarehouseNotFoundError
from internal.domain.model.option.inventory_option import StockQueryOption
from internal.infrastructure.adapter.mysql.models import StockEntity, WarehouseEntity
from internal.infrastructure.adapter.mysql.optimistic_retry import run_with_optimistic_retry
from internal.infrastructure.adapter.mysql.paging import page_offset
from internal.infrastructure.adapter.mysql.stock_movement import StockDelta, apply_stock_delta, lock_stock
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes


def _stock_rows():
    return (
        select(
            StockEntity.uid,
            StockEntity.edition_uid,
            StockEntity.warehouse_uid,
            WarehouseEntity.name,
            StockEntity.quantity,
        )
        .select_from(StockEntity)
        .outerjoin(WarehouseEntity, WarehouseEntity.uid == StockEntity.warehouse_uid)
    )


def _group_rows(rows: Sequence[Row]) -> list[StockAggregate]:
    aggregates: list[StockAggregate] = []
    by_edition: dict[bytes, StockAggregate] = {}
    for uid, edition_uid, warehouse_uid, warehouse_name, quantity in rows:
        aggregate = by_edition.get(edition_uid)
        if aggregate is None:
            aggregate = StockAggregate(uid=bytes_to_uuid(uid), edition_uid=bytes_to_uuid(edition_uid), total_quantity=0)
            by_edition[edition_uid] = aggregate
            aggregates.append(aggregate)
        aggregate.warehouses.append(
            StockWarehouse(
                warehouse_uid=bytes_to_uuid(warehouse_uid),
                warehouse_name=warehouse_name or "",
                quantity=quantity,
            )
        )
        aggregate.total_quantity += quantity
    return aggregates


class MySQLStockRepository(StockRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def find_all(self, option: StockQueryOption) -> tuple[list[StockAggregate], int]:
        conditions = self._conditions(option.edition_uid, option.warehouse_uid)
        async with self._session_factory() as session:
            total = (
                await session.execute(select(func.count(distinct(StockEntity.edition_uid))).where(*conditions))
            ).scalar() or 0
            edition_stmt = (
                select(StockEntity.edition_uid)
                .where(*conditions)
                .group_by(StockEntity.edition_uid)
                .order_by(StockEntity.edition_uid)
                .offset(page_offset(option.page, option.size))
                .limit(option.size)
            )
            editions = list((await session.execute(edition_stmt)).scalars().all())
            if not editions:
                return [], total
            rows = (
                await session.execute(
                    _stock_rows()
                    .where(*conditions, StockEntity.edition_uid.in_(editions))
                    .order_by(StockEntity.edition_uid, StockEntity.id)
                )
            ).all()
            return _group_rows(rows), total

    async def find_by_edition_uid(self, edition_uid: UUID) -> StockAggregate | None:
        async with self._session_factory() as session:
            rows = (
                await session.execute(
                    _stock_rows().where(*self._conditions(edition_uid, None)).order_by(StockEntity.id)
                )
            ).all()
        if not rows:
            return None
        return _group_rows(rows)[0]

    async def find_stock_uids(self, edition_uid: UUID, warehouse_uid: UUID | None) -> list[UUID]:
        async with self._session_factory() as session:
            uids = (
                await session.execute(
                    select(StockEntity.uid)
                    .where(*self._conditions(edition_uid, warehouse_uid))
                    .order_by(StockEntity.id)
                )
            ).scalars()
            return [bytes_to_uuid(uid) for uid in uids]

    async def apply_change(self, params: ApplyChangeParams) -> None:
        edition_uid = uuid_to_bytes(params.edition_uid)
        warehouse_uid = uuid_to_bytes(params.warehouse_uid)

        async def attempt() -> None:
            async with self._session_factory() as session, session.begin():
                target = await lock_stock(session, edition_uid, warehouse_uid)
                if target is not None:
                    await apply_stock_delta(
                        session,
                        StockDelta(
                            edition_id=target.edition_id,
                            edition_uid=edition_uid,
                            warehouse_id=target.warehouse_id,
                            warehouse_uid=warehouse_uid,
                            delta=params.delta,
                            change_type=params.change_type,
                            reason=params.reason,
                        ),
                        existing=target,
                    )
                    return
                if params.change_type != StockChangeType.INBOUND or params.delta <= 0:
                    raise StockNotFoundError(
                        f"stock not found for edition {params.edition_uid} in warehouse {params.warehouse_uid}"
                    )
                warehouse = (
                    await session.execute(
                        select(WarehouseEntity).where(
                            WarehouseEntity.uid == warehouse_uid, WarehouseEntity.deleted_at.is_(None)
                        )
                    )
                ).scalar_one_or_none()
                if warehouse is None:
                    raise WarehouseNotFoundError(f"warehouse {params.warehouse_uid} not found")
                sibling = (
                    await session.execute(
                        select(StockEntity)
                        .where(StockEntity.edition_uid == edition_uid, StockEntity.deleted_at.is_(None))
                        .order_by(StockEntity.id)
                        .limit(1)
                    )
                ).scalar_one_or_none()
                if sibling is None:
                    raise StockNotFoundError(
                        f"edition {params.edition_uid} has no stock in any warehouse to resolve its id"
                    )
                await apply_stock_delta(
                    session,
                    StockDelta(
                        edition_id=sibling.edition_id,
                        edition_uid=edition_uid,
                        warehouse_id=warehouse.id,
                        warehouse_uid=warehouse.uid,
                        delta=params.delta,
                        change_type=params.change_type,
                        reason=params.reason,
                        allow_create=True,
                    ),
                )

        await run_with_optimistic_retry(attempt)

    @staticmethod
    def _conditions(edition_uid: UUID | None, warehouse_uid: UUID | None) -> list:
        conditions = [StockEntity.deleted_at.is_(None)]
        if edition_uid is not None:
            conditions.append(StockEntity.edition_uid == uuid_to_bytes(edition_uid))
        if warehouse_uid is not None:
            conditions.append(StockEntity.warehouse_uid == uuid_to_bytes(warehouse_uid))
        return conditions
