from collections.abc import Callable
from uuid import UUID

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import CarrierRepository
from internal.domain.aggregate.carrier_aggregate import CarrierAggregate
from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.model.option.delivery_option import CarrierFilterOption
from internal.infrastructure.adapter.mysql.models import CarrierEntity
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes


class MySQLCarrierRepository(CarrierRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, carrier: CarrierAggregate) -> CarrierAggregate:
        async with self._session_factory() as session:
            entity = (
                await session.execute(select(CarrierEntity).where(CarrierEntity.uid == uuid_to_bytes(carrier.uid)))
            ).scalar_one_or_none()

            if entity is None:
                entity = CarrierEntity(
                    uid=uuid_to_bytes(carrier.uid),
                    name=carrier.name,
                    code=carrier.code,
                    contact_name=carrier.contact_name,
                    contact_phone=carrier.contact_phone,
                    contact_email=carrier.contact_email,
                    base_rate=carrier.base_rate,
                    rate_per_kg=carrier.rate_per_kg,
                    status=carrier.status,
                    created_at=carrier.created_at,
                    updated_at=carrier.updated_at,
                )
                session.add(entity)
            else:
                entity.name = carrier.name
                entity.contact_name = carrier.contact_name
                entity.contact_phone = carrier.contact_phone
                entity.contact_email = carrier.contact_email
                entity.base_rate = carrier.base_rate
                entity.rate_per_kg = carrier.rate_per_kg
                entity.status = carrier.status
                entity.updated_at = carrier.updated_at

            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    async def find_by_uid(self, uid: UUID) -> CarrierAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(CarrierEntity).where(
                        CarrierEntity.uid == uuid_to_bytes(uid),
                        CarrierEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def find_by_code(self, code: str) -> CarrierAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(select(CarrierEntity).where(CarrierEntity.code == code))
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def find_all(self, option: CarrierFilterOption) -> tuple[list[CarrierAggregate], int]:
        async with self._session_factory() as session:
            query = select(CarrierEntity).where(CarrierEntity.deleted_at.is_(None))
            count_query = select(func.count()).select_from(CarrierEntity).where(CarrierEntity.deleted_at.is_(None))

            if option.name is not None:
                query = query.where(CarrierEntity.name.like(f"%{option.name}%"))
                count_query = count_query.where(CarrierEntity.name.like(f"%{option.name}%"))
            if option.status is not None:
                query = query.where(CarrierEntity.status == option.status)
                count_query = count_query.where(CarrierEntity.status == option.status)

            sort_field, sort_dir = _parse_sort(option.sort)
            order_col = _SORTABLE_COLUMNS.get(sort_field, _DEFAULT_SORT_COLUMN)
            query = query.order_by(order_col.desc() if sort_dir == "desc" else order_col.asc())
            query = query.offset(option.page * option.size).limit(option.size)

            total = (await session.execute(count_query)).scalar() or 0
            entities = (await session.execute(query)).scalars().all()
            return [self._to_aggregate(e) for e in entities], total

    async def update(
        self,
        uid: UUID,
        mutator: Callable[[CarrierAggregate], CarrierAggregate],
    ) -> CarrierAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(CarrierEntity)
                    .where(
                        CarrierEntity.uid == uuid_to_bytes(uid),
                        CarrierEntity.deleted_at.is_(None),
                    )
                    .with_for_update()
                )
            ).scalar_one_or_none()
            if entity is None:
                return None
            updated = mutator(self._to_aggregate(entity))
            entity.name = updated.name
            entity.contact_name = updated.contact_name
            entity.contact_phone = updated.contact_phone
            entity.contact_email = updated.contact_email
            entity.base_rate = updated.base_rate
            entity.rate_per_kg = updated.rate_per_kg
            entity.status = updated.status
            entity.updated_at = updated.updated_at
            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    @staticmethod
    def _to_aggregate(entity: CarrierEntity) -> CarrierAggregate:
        return CarrierAggregate(
            uid=bytes_to_uuid(entity.uid),
            name=entity.name,
            code=entity.code,
            contact_name=entity.contact_name,
            contact_phone=entity.contact_phone,
            contact_email=entity.contact_email,
            base_rate=entity.base_rate,
            rate_per_kg=entity.rate_per_kg,
            status=CarrierStatus(entity.status),
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )


_SORTABLE_COLUMNS = {
    "name": CarrierEntity.name,
    "code": CarrierEntity.code,
    "status": CarrierEntity.status,
    "base_rate": CarrierEntity.base_rate,
    "rate_per_kg": CarrierEntity.rate_per_kg,
    "created_at": CarrierEntity.created_at,
    "updated_at": CarrierEntity.updated_at,
}
_DEFAULT_SORT_COLUMN = CarrierEntity.name


def _parse_sort(sort: str) -> tuple[str, str]:
    parts = sort.split(":")
    field = parts[0] if len(parts) > 0 else "name"
    direction = parts[1] if len(parts) > 1 else "asc"
    return field, direction
