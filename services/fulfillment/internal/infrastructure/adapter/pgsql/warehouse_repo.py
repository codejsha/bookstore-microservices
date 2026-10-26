from datetime import UTC, datetime
from uuid import UUID, uuid7

from sqlalchemy import func, select, update
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import (
    WarehouseCreateParams,
    WarehouseRepository,
    WarehouseUpdateParams,
)
from internal.domain.aggregate.warehouse_aggregate import WarehouseAggregate
from internal.domain.model.option.inventory_option import WarehouseQueryOption
from internal.infrastructure.adapter.pgsql.models import WarehouseEntity
from internal.infrastructure.adapter.pgsql.optimistic_retry import VersionConflictError, run_with_optimistic_retry
from internal.infrastructure.adapter.pgsql.paging import page_offset, resolve_order_by

_SORTABLE = {
    "uid": WarehouseEntity.uid,
    "name": WarehouseEntity.name,
    "address": WarehouseEntity.address,
    "capacity": WarehouseEntity.capacity,
    "created_at": WarehouseEntity.created_at,
    "updated_at": WarehouseEntity.updated_at,
}


class PgsqlWarehouseRepository(WarehouseRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def find_all(self, option: WarehouseQueryOption) -> tuple[list[WarehouseAggregate], int]:
        conditions = [WarehouseEntity.deleted_at.is_(None)]
        if option.name:
            conditions.append(func.lower(WarehouseEntity.name).like(f"%{option.name.lower()}%"))
        async with self._session_factory() as session:
            total = (
                await session.execute(select(func.count()).select_from(WarehouseEntity).where(*conditions))
            ).scalar() or 0
            entities = (
                await session.execute(
                    select(WarehouseEntity)
                    .where(*conditions)
                    .order_by(*resolve_order_by(option.sort, _SORTABLE))
                    .offset(page_offset(option.page, option.size))
                    .limit(option.size)
                )
            ).scalars()
            return [self._to_aggregate(e) for e in entities], total

    async def find_by_uid(self, uid: UUID) -> WarehouseAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(WarehouseEntity).where(WarehouseEntity.uid == uid, WarehouseEntity.deleted_at.is_(None))
                )
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def create(self, params: WarehouseCreateParams) -> WarehouseAggregate:
        now = datetime.now(UTC)
        entity = WarehouseEntity(
            uid=uuid7(),
            name=params.name,
            address=params.address,
            capacity=params.capacity,
            created_at=now,
            updated_at=None,
            deleted_at=None,
            version=1,
        )
        async with self._session_factory() as session:
            session.add(entity)
            await session.commit()
            return self._to_aggregate(entity)

    async def update(self, params: WarehouseUpdateParams) -> WarehouseAggregate | None:
        async def attempt() -> WarehouseAggregate | None:
            async with self._session_factory() as session, session.begin():
                entity = (
                    await session.execute(
                        select(WarehouseEntity)
                        .where(WarehouseEntity.uid == params.uid, WarehouseEntity.deleted_at.is_(None))
                        .with_for_update()
                    )
                ).scalar_one_or_none()
                if entity is None:
                    return None
                now = datetime.now(UTC)
                name = params.name if params.name is not None else entity.name
                address = params.address if params.address is not None else entity.address
                capacity = params.capacity if params.capacity is not None else entity.capacity
                result = await session.execute(
                    update(WarehouseEntity)
                    .where(WarehouseEntity.id == entity.id, WarehouseEntity.version == entity.version)
                    .values(name=name, address=address, capacity=capacity, updated_at=now, version=entity.version + 1)
                )
                if result.rowcount == 0:
                    raise VersionConflictError(f"warehouse {params.uid} changed concurrently")
                return WarehouseAggregate(
                    uid=params.uid,
                    name=name,
                    address=address,
                    capacity=capacity,
                    created_at=entity.created_at,
                    updated_at=now,
                )

        return await run_with_optimistic_retry(attempt)

    @staticmethod
    def _to_aggregate(entity: WarehouseEntity) -> WarehouseAggregate:
        return WarehouseAggregate(
            uid=entity.uid,
            name=entity.name,
            address=entity.address,
            capacity=entity.capacity,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )
