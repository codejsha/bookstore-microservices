from sqlalchemy import and_, func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import StockBalanceRepository
from internal.domain.aggregate.balance_aggregate import StockBalanceEntry
from internal.domain.model.option.inventory_option import BalanceQueryOption
from internal.infrastructure.adapter.mysql.closing_repo import month_range
from internal.infrastructure.adapter.mysql.models import StockEntity, StockHistoryEntity, WarehouseEntity
from internal.infrastructure.adapter.mysql.paging import page_offset
from internal.infrastructure.adapter.mysql.stock_movement import history_buckets
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes


class MySQLStockBalanceRepository(StockBalanceRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def find_all(self, option: BalanceQueryOption) -> tuple[list[StockBalanceEntry], int]:
        stock, warehouse, history = StockEntity, WarehouseEntity, StockHistoryEntity

        join_condition = and_(history.stock_id == stock.id, history.deleted_at.is_(None))
        if option.year_month is not None:
            year, month = (int(part) for part in option.year_month.split("-"))
            start, end = month_range(year, month)
            join_condition = and_(join_condition, history.created_at >= start, history.created_at < end)

        conditions = [stock.deleted_at.is_(None)]
        if option.edition_uid is not None:
            conditions.append(stock.edition_uid == uuid_to_bytes(option.edition_uid))
        if option.warehouse_uid is not None:
            conditions.append(stock.warehouse_uid == uuid_to_bytes(option.warehouse_uid))

        inbound, outbound, adjust = history_buckets()
        rows_stmt = (
            select(stock.edition_uid, stock.warehouse_uid, warehouse.name, stock.quantity, inbound, outbound, adjust)
            .select_from(stock)
            .join(warehouse, warehouse.id == stock.warehouse_id)
            .outerjoin(history, join_condition)
            .where(*conditions)
            .group_by(stock.id, stock.edition_uid, stock.warehouse_uid, warehouse.name, stock.quantity)
            .order_by(stock.id)
            .offset(page_offset(option.page, option.size))
            .limit(option.size)
        )
        count_stmt = (
            select(func.count())
            .select_from(stock)
            .join(warehouse, warehouse.id == stock.warehouse_id)
            .where(*conditions)
        )

        async with self._session_factory() as session:
            rows = (await session.execute(rows_stmt)).all()
            total = (await session.execute(count_stmt)).scalar() or 0

        return [
            StockBalanceEntry(
                edition_uid=bytes_to_uuid(edition_uid),
                warehouse_uid=bytes_to_uuid(warehouse_uid),
                warehouse_name=name,
                inbound_quantity=int(in_qty),
                outbound_quantity=int(out_qty),
                adjust_quantity=int(adj_qty),
                current_quantity=quantity,
            )
            for edition_uid, warehouse_uid, name, quantity, in_qty, out_qty, adj_qty in rows
        ], total
