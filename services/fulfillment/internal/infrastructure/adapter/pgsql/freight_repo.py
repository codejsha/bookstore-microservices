from collections.abc import Callable
from uuid import UUID

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import FreightRepository
from internal.domain.aggregate.freight_aggregate import FreightAggregate
from internal.domain.constant.freight_status import FreightStatus
from internal.domain.model.option.delivery_option import FreightFilterOption
from internal.infrastructure.adapter.mysql.models import FreightEntity
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes


class MySQLFreightRepository(FreightRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, freight: FreightAggregate) -> FreightAggregate:
        async with self._session_factory() as session:
            entity = (
                await session.execute(select(FreightEntity).where(FreightEntity.uid == uuid_to_bytes(freight.uid)))
            ).scalar_one_or_none()

            if entity is None:
                entity = FreightEntity(
                    uid=uuid_to_bytes(freight.uid),
                    shipment_uid=uuid_to_bytes(freight.shipment_uid),
                    carrier_uid=uuid_to_bytes(freight.carrier_uid),
                    base_cost=freight.base_cost,
                    weight_surcharge=freight.weight_surcharge,
                    distance_surcharge=freight.distance_surcharge,
                    discount=freight.discount,
                    total_cost=freight.total_cost,
                    currency=freight.currency,
                    status=freight.status,
                    invoiced_at=freight.invoiced_at,
                    paid_at=freight.paid_at,
                    created_at=freight.created_at,
                    updated_at=freight.updated_at,
                )
                session.add(entity)
            else:
                entity.status = freight.status
                entity.invoiced_at = freight.invoiced_at
                entity.paid_at = freight.paid_at
                entity.updated_at = freight.updated_at

            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    async def find_by_uid(self, uid: UUID) -> FreightAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(FreightEntity).where(
                        FreightEntity.uid == uuid_to_bytes(uid),
                        FreightEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def find_by_shipment_uid(self, shipment_uid: UUID) -> FreightAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(FreightEntity).where(
                        FreightEntity.shipment_uid == uuid_to_bytes(shipment_uid),
                        FreightEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def find_all(self, option: FreightFilterOption) -> tuple[list[FreightAggregate], int]:
        async with self._session_factory() as session:
            query = select(FreightEntity).where(FreightEntity.deleted_at.is_(None))
            count_query = select(func.count()).select_from(FreightEntity).where(FreightEntity.deleted_at.is_(None))

            if option.shipment_uid is not None:
                uid_bytes = uuid_to_bytes(option.shipment_uid)
                query = query.where(FreightEntity.shipment_uid == uid_bytes)
                count_query = count_query.where(FreightEntity.shipment_uid == uid_bytes)
            if option.carrier_uid is not None:
                uid_bytes = uuid_to_bytes(option.carrier_uid)
                query = query.where(FreightEntity.carrier_uid == uid_bytes)
                count_query = count_query.where(FreightEntity.carrier_uid == uid_bytes)
            if option.status is not None:
                query = query.where(FreightEntity.status == option.status)
                count_query = count_query.where(FreightEntity.status == option.status)

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
        mutator: Callable[[FreightAggregate], FreightAggregate],
    ) -> FreightAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(FreightEntity)
                    .where(
                        FreightEntity.uid == uuid_to_bytes(uid),
                        FreightEntity.deleted_at.is_(None),
                    )
                    .with_for_update()
                )
            ).scalar_one_or_none()
            if entity is None:
                return None
            updated = mutator(self._to_aggregate(entity))
            entity.status = updated.status
            entity.invoiced_at = updated.invoiced_at
            entity.paid_at = updated.paid_at
            entity.updated_at = updated.updated_at
            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    @staticmethod
    def _to_aggregate(entity: FreightEntity) -> FreightAggregate:
        return FreightAggregate(
            uid=bytes_to_uuid(entity.uid),
            shipment_uid=bytes_to_uuid(entity.shipment_uid),
            carrier_uid=bytes_to_uuid(entity.carrier_uid),
            base_cost=entity.base_cost,
            weight_surcharge=entity.weight_surcharge,
            distance_surcharge=entity.distance_surcharge,
            discount=entity.discount,
            total_cost=entity.total_cost,
            currency=entity.currency,
            status=FreightStatus(entity.status),
            invoiced_at=entity.invoiced_at,
            paid_at=entity.paid_at,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )


_SORTABLE_COLUMNS = {
    "created_at": FreightEntity.created_at,
    "updated_at": FreightEntity.updated_at,
    "status": FreightEntity.status,
    "total_cost": FreightEntity.total_cost,
    "invoiced_at": FreightEntity.invoiced_at,
    "paid_at": FreightEntity.paid_at,
}
_DEFAULT_SORT_COLUMN = FreightEntity.created_at


def _parse_sort(sort: str) -> tuple[str, str]:
    parts = sort.split(":")
    field = parts[0] if len(parts) > 0 else "created_at"
    direction = parts[1] if len(parts) > 1 else "desc"
    return field, direction
