from uuid import UUID

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import StockHistoryRepository
from internal.domain.aggregate.stock_aggregate import StockHistoryEntry
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.option.inventory_option import StockHistoryQueryOption
from internal.infrastructure.adapter.mysql.models import StockHistoryEntity
from internal.infrastructure.adapter.mysql.paging import page_offset
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes


class MySQLStockHistoryRepository(StockHistoryRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def find_all(
        self, stock_uids: list[UUID], option: StockHistoryQueryOption
    ) -> tuple[list[StockHistoryEntry], int]:
        if not stock_uids:
            return [], 0
        conditions = [
            StockHistoryEntity.stock_uid.in_([uuid_to_bytes(uid) for uid in stock_uids]),
            StockHistoryEntity.deleted_at.is_(None),
        ]
        async with self._session_factory() as session:
            total = (
                await session.execute(select(func.count()).select_from(StockHistoryEntity).where(*conditions))
            ).scalar() or 0
            entities = (
                await session.execute(
                    select(StockHistoryEntity)
                    .where(*conditions)
                    .order_by(StockHistoryEntity.created_at.desc(), StockHistoryEntity.uid.desc())
                    .offset(page_offset(option.page, option.size))
                    .limit(option.size)
                )
            ).scalars()
            return [self._to_entry(e) for e in entities], total

    @staticmethod
    def _to_entry(entity: StockHistoryEntity) -> StockHistoryEntry:
        return StockHistoryEntry(
            uid=bytes_to_uuid(entity.uid),
            change_type=StockChangeType(entity.change_type),
            reason=entity.reason,
            change_qty=entity.change_qty,
            created_at=entity.created_at,
        )
