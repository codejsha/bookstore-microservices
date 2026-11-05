from collections.abc import Sequence
from datetime import UTC, datetime
from uuid import UUID, uuid7

from sqlalchemy import func, select
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import ClosingCreateParams, MonthlyClosingRepository
from internal.domain.aggregate.closing_aggregate import MonthlyClosingAggregate, MonthlyClosingItem
from internal.domain.constant.closing_status import ClosingStatus
from internal.domain.model.error import DuplicateClosingError, WarehouseNotFoundError
from internal.domain.model.option.inventory_option import ClosingQueryOption
from internal.infrastructure.adapter.pgsql.models import (
    MonthlyClosingEntity,
    MonthlyClosingItemEntity,
    StockEntity,
    StockHistoryEntity,
    WarehouseEntity,
)
from internal.infrastructure.adapter.pgsql.optimistic_retry import is_unique_violation
from internal.infrastructure.adapter.pgsql.paging import page_offset, resolve_order_by
from internal.infrastructure.adapter.pgsql.stock_movement import history_buckets

CLOSING_ITEM_BATCH_SIZE = 500

_SORTABLE = {
    "uid": MonthlyClosingEntity.uid,
    "warehouse_uid": MonthlyClosingEntity.warehouse_uid,
    "year": MonthlyClosingEntity.year,
    "month": MonthlyClosingEntity.month,
    "status": MonthlyClosingEntity.status,
    "closed_at": MonthlyClosingEntity.closed_at,
    "created_at": MonthlyClosingEntity.created_at,
    "updated_at": MonthlyClosingEntity.updated_at,
}


def month_range(year: int, month: int) -> tuple[datetime, datetime]:
    start = datetime(year, month, 1, tzinfo=UTC)
    end = datetime(year + 1, 1, 1, tzinfo=UTC) if month == 12 else datetime(year, month + 1, 1, tzinfo=UTC)
    return start, end


class PgsqlMonthlyClosingRepository(MonthlyClosingRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def find_all(self, option: ClosingQueryOption) -> tuple[list[MonthlyClosingAggregate], int]:
        conditions = [MonthlyClosingEntity.deleted_at.is_(None)]
        if option.warehouse_uid is not None:
            conditions.append(MonthlyClosingEntity.warehouse_uid == option.warehouse_uid)
        if option.year is not None:
            conditions.append(MonthlyClosingEntity.year == option.year)
        if option.month is not None:
            conditions.append(MonthlyClosingEntity.month == option.month)
        if option.status is not None:
            conditions.append(MonthlyClosingEntity.status == option.status.value)
        async with self._session_factory() as session:
            total = (
                await session.execute(select(func.count()).select_from(MonthlyClosingEntity).where(*conditions))
            ).scalar() or 0
            closings = list(
                (
                    await session.execute(
                        select(MonthlyClosingEntity)
                        .where(*conditions)
                        .order_by(*resolve_order_by(option.sort, _SORTABLE))
                        .offset(page_offset(option.page, option.size))
                        .limit(option.size)
                    )
                ).scalars()
            )
            items_by_closing = await self._load_items(session, [c.id for c in closings])
            return [self._to_aggregate(c, items_by_closing.get(c.id, [])) for c in closings], total

    async def find_by_uid(self, uid: UUID) -> MonthlyClosingAggregate | None:
        async with self._session_factory() as session:
            closing = (
                await session.execute(
                    select(MonthlyClosingEntity).where(
                        MonthlyClosingEntity.uid == uid, MonthlyClosingEntity.deleted_at.is_(None)
                    )
                )
            ).scalar_one_or_none()
            if closing is None:
                return None
            items = await self._load_items(session, [closing.id])
            return self._to_aggregate(closing, items.get(closing.id, []))

    async def create(self, params: ClosingCreateParams) -> MonthlyClosingAggregate:
        start, end = month_range(params.year, params.month)
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

            stocks = list(
                (
                    await session.execute(
                        select(StockEntity)
                        .where(StockEntity.warehouse_id == warehouse.id, StockEntity.deleted_at.is_(None))
                        .order_by(StockEntity.id)
                    )
                ).scalars()
            )

            now = datetime.now(UTC)
            closing = MonthlyClosingEntity(
                uid=uuid7(),
                warehouse_id=warehouse.id,
                warehouse_uid=warehouse.uid,
                year=params.year,
                month=params.month,
                status=ClosingStatus.CLOSED.value,
                closed_at=now,
                created_at=now,
                updated_at=None,
                deleted_at=None,
                version=1,
            )
            session.add(closing)
            try:
                await session.flush()
            except IntegrityError as exc:
                if is_unique_violation(exc):
                    raise DuplicateClosingError(
                        f"closing already exists for warehouse {params.warehouse_uid} {params.year}-{params.month:02d}"
                    ) from exc
                raise

            movements = await self._period_movements(session, warehouse.id, start, end)
            later = await self._net_change_since(session, warehouse.id, end)
            items = [
                MonthlyClosingItemEntity(
                    closing_id=closing.id,
                    edition_id=stock.edition_id,
                    edition_uid=stock.edition_uid,
                    opening_quantity=closing_quantity - (inbound - outbound + adjust),
                    inbound_quantity=inbound,
                    outbound_quantity=outbound,
                    adjust_quantity=adjust,
                    closing_quantity=closing_quantity,
                )
                for stock in stocks
                for inbound, outbound, adjust in (movements.get(stock.id, (0, 0, 0)),)
                for closing_quantity in (stock.quantity - later.get(stock.id, 0),)
            ]
            for offset in range(0, len(items), CLOSING_ITEM_BATCH_SIZE):
                session.add_all(items[offset : offset + CLOSING_ITEM_BATCH_SIZE])
                await session.flush()
            return self._to_aggregate(closing, items)

    @staticmethod
    async def _period_movements(
        session: AsyncSession, warehouse_id: int, start: datetime, end: datetime
    ) -> dict[int, tuple[int, int, int]]:
        inbound, outbound, adjust = history_buckets()
        warehouse_stocks = select(StockEntity.id).where(
            StockEntity.warehouse_id == warehouse_id, StockEntity.deleted_at.is_(None)
        )
        rows = (
            await session.execute(
                select(StockHistoryEntity.stock_id, inbound, outbound, adjust)
                .where(
                    StockHistoryEntity.stock_id.in_(warehouse_stocks),
                    StockHistoryEntity.deleted_at.is_(None),
                    StockHistoryEntity.created_at >= start,
                    StockHistoryEntity.created_at < end,
                )
                .group_by(StockHistoryEntity.stock_id)
            )
        ).all()
        return {stock_id: (int(i), int(o), int(a)) for stock_id, i, o, a in rows}

    @staticmethod
    async def _net_change_since(session: AsyncSession, warehouse_id: int, since: datetime) -> dict[int, int]:
        warehouse_stocks = select(StockEntity.id).where(
            StockEntity.warehouse_id == warehouse_id, StockEntity.deleted_at.is_(None)
        )
        rows = (
            await session.execute(
                select(StockHistoryEntity.stock_id, func.coalesce(func.sum(StockHistoryEntity.change_qty), 0))
                .where(
                    StockHistoryEntity.stock_id.in_(warehouse_stocks),
                    StockHistoryEntity.deleted_at.is_(None),
                    StockHistoryEntity.created_at >= since,
                )
                .group_by(StockHistoryEntity.stock_id)
            )
        ).all()
        return {stock_id: int(net) for stock_id, net in rows}

    @staticmethod
    async def _load_items(
        session: AsyncSession, closing_ids: Sequence[int]
    ) -> dict[int, list[MonthlyClosingItemEntity]]:
        grouped: dict[int, list[MonthlyClosingItemEntity]] = {}
        if not closing_ids:
            return grouped
        items = (
            await session.execute(
                select(MonthlyClosingItemEntity)
                .where(MonthlyClosingItemEntity.closing_id.in_(list(closing_ids)))
                .order_by(MonthlyClosingItemEntity.closing_id, MonthlyClosingItemEntity.id)
            )
        ).scalars()
        for item in items:
            grouped.setdefault(item.closing_id, []).append(item)
        return grouped

    @staticmethod
    def _to_aggregate(
        closing: MonthlyClosingEntity, items: Sequence[MonthlyClosingItemEntity]
    ) -> MonthlyClosingAggregate:
        return MonthlyClosingAggregate(
            uid=closing.uid,
            warehouse_uid=closing.warehouse_uid,
            year=closing.year,
            month=closing.month,
            status=ClosingStatus(closing.status),
            closed_at=closing.closed_at,
            created_at=closing.created_at,
            updated_at=closing.updated_at,
            items=[
                MonthlyClosingItem(
                    edition_uid=i.edition_uid,
                    opening_quantity=i.opening_quantity,
                    inbound_quantity=i.inbound_quantity,
                    outbound_quantity=i.outbound_quantity,
                    adjust_quantity=i.adjust_quantity,
                    closing_quantity=i.closing_quantity,
                )
                for i in items
            ],
        )
