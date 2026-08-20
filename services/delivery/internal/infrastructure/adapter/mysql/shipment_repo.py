from collections.abc import Callable
from uuid import UUID

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import ShipmentRepository
from internal.domain.aggregate.shipment_aggregate import ShipmentAggregate
from internal.domain.aggregate.tracking_aggregate import TrackingAggregate
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.option.delivery_option import ShipmentFilterOption
from internal.infrastructure.adapter.mysql.models import ShipmentEntity, TrackingEntity
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes


class MySQLShipmentRepository(ShipmentRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, shipment: ShipmentAggregate) -> ShipmentAggregate:
        async with self._session_factory() as session:
            entity = (
                await session.execute(select(ShipmentEntity).where(ShipmentEntity.uid == uuid_to_bytes(shipment.uid)))
            ).scalar_one_or_none()

            if entity is None:
                entity = ShipmentEntity(
                    uid=uuid_to_bytes(shipment.uid),
                    order_uid=uuid_to_bytes(shipment.order_uid),
                    carrier_uid=uuid_to_bytes(shipment.carrier_uid) if shipment.carrier_uid else None,
                    origin_address=shipment.origin_address,
                    destination_address=shipment.destination_address,
                    destination_city=shipment.destination_city,
                    destination_state=shipment.destination_state,
                    destination_country_code=shipment.destination_country_code,
                    destination_postal_code=shipment.destination_postal_code,
                    status=shipment.status,
                    tracking_number=shipment.tracking_number,
                    weight_kg=shipment.weight_kg,
                    planned_pickup_at=shipment.planned_pickup_at,
                    planned_delivery_at=shipment.planned_delivery_at,
                    actual_pickup_at=shipment.actual_pickup_at,
                    actual_delivery_at=shipment.actual_delivery_at,
                    created_at=shipment.created_at,
                    updated_at=shipment.updated_at,
                )
                session.add(entity)
            else:
                entity.carrier_uid = uuid_to_bytes(shipment.carrier_uid) if shipment.carrier_uid else None
                entity.status = shipment.status
                entity.tracking_number = shipment.tracking_number
                entity.actual_pickup_at = shipment.actual_pickup_at
                entity.actual_delivery_at = shipment.actual_delivery_at
                entity.updated_at = shipment.updated_at

            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    async def create_with_initial_tracking(
        self, shipment: ShipmentAggregate, tracking: TrackingAggregate
    ) -> ShipmentAggregate:
        async with self._session_factory() as session:
            shipment_entity = ShipmentEntity(
                uid=uuid_to_bytes(shipment.uid),
                order_uid=uuid_to_bytes(shipment.order_uid),
                carrier_uid=uuid_to_bytes(shipment.carrier_uid) if shipment.carrier_uid else None,
                origin_address=shipment.origin_address,
                destination_address=shipment.destination_address,
                destination_city=shipment.destination_city,
                destination_state=shipment.destination_state,
                destination_country_code=shipment.destination_country_code,
                destination_postal_code=shipment.destination_postal_code,
                status=shipment.status,
                tracking_number=shipment.tracking_number,
                weight_kg=shipment.weight_kg,
                planned_pickup_at=shipment.planned_pickup_at,
                planned_delivery_at=shipment.planned_delivery_at,
                actual_pickup_at=shipment.actual_pickup_at,
                actual_delivery_at=shipment.actual_delivery_at,
                created_at=shipment.created_at,
                updated_at=shipment.updated_at,
            )
            tracking_entity = TrackingEntity(
                uid=uuid_to_bytes(tracking.uid),
                shipment_uid=uuid_to_bytes(tracking.shipment_uid),
                status=tracking.status,
                location=tracking.location,
                description=tracking.description,
                occurred_at=tracking.occurred_at,
                created_at=tracking.created_at,
            )
            session.add(shipment_entity)
            session.add(tracking_entity)
            await session.commit()
            await session.refresh(shipment_entity)
            return self._to_aggregate(shipment_entity)

    async def find_by_uid(self, uid: UUID) -> ShipmentAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(ShipmentEntity).where(
                        ShipmentEntity.uid == uuid_to_bytes(uid),
                        ShipmentEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def find_by_order_uid(self, order_uid: UUID) -> ShipmentAggregate | None:
        async with self._session_factory() as session:
            entity = (
                (
                    await session.execute(
                        select(ShipmentEntity)
                        .where(
                            ShipmentEntity.order_uid == uuid_to_bytes(order_uid),
                            ShipmentEntity.deleted_at.is_(None),
                        )
                        .order_by(ShipmentEntity.created_at.asc())
                    )
                )
                .scalars()
                .first()
            )
            return self._to_aggregate(entity) if entity else None

    async def find_all(self, option: ShipmentFilterOption) -> tuple[list[ShipmentAggregate], int]:
        async with self._session_factory() as session:
            query = select(ShipmentEntity).where(ShipmentEntity.deleted_at.is_(None))
            count_query = select(func.count()).select_from(ShipmentEntity).where(ShipmentEntity.deleted_at.is_(None))

            if option.order_uid is not None:
                uid_bytes = uuid_to_bytes(option.order_uid)
                query = query.where(ShipmentEntity.order_uid == uid_bytes)
                count_query = count_query.where(ShipmentEntity.order_uid == uid_bytes)
            if option.carrier_uid is not None:
                uid_bytes = uuid_to_bytes(option.carrier_uid)
                query = query.where(ShipmentEntity.carrier_uid == uid_bytes)
                count_query = count_query.where(ShipmentEntity.carrier_uid == uid_bytes)
            if option.status is not None:
                query = query.where(ShipmentEntity.status == option.status)
                count_query = count_query.where(ShipmentEntity.status == option.status)

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
        mutator: Callable[[ShipmentAggregate], ShipmentAggregate],
    ) -> ShipmentAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(ShipmentEntity)
                    .where(
                        ShipmentEntity.uid == uuid_to_bytes(uid),
                        ShipmentEntity.deleted_at.is_(None),
                    )
                    .with_for_update()
                )
            ).scalar_one_or_none()
            if entity is None:
                return None
            updated = mutator(self._to_aggregate(entity))
            entity.carrier_uid = uuid_to_bytes(updated.carrier_uid) if updated.carrier_uid else None
            entity.status = updated.status
            entity.tracking_number = updated.tracking_number
            entity.actual_pickup_at = updated.actual_pickup_at
            entity.actual_delivery_at = updated.actual_delivery_at
            entity.updated_at = updated.updated_at
            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    async def transition(
        self,
        uid: UUID,
        guard: Callable[[ShipmentAggregate], bool],
        mutator: Callable[[ShipmentAggregate], ShipmentAggregate],
    ) -> ShipmentAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(ShipmentEntity)
                    .where(
                        ShipmentEntity.uid == uuid_to_bytes(uid),
                        ShipmentEntity.deleted_at.is_(None),
                    )
                    .with_for_update()
                )
            ).scalar_one_or_none()
            if entity is None:
                return None
            aggregate = self._to_aggregate(entity)
            if not guard(aggregate):
                raise ValueError(f"Cannot transition shipment from status {aggregate.status}")
            updated = mutator(aggregate)
            entity.carrier_uid = uuid_to_bytes(updated.carrier_uid) if updated.carrier_uid else None
            entity.status = updated.status
            entity.tracking_number = updated.tracking_number
            entity.actual_pickup_at = updated.actual_pickup_at
            entity.actual_delivery_at = updated.actual_delivery_at
            entity.updated_at = updated.updated_at
            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    @staticmethod
    def _to_aggregate(entity: ShipmentEntity) -> ShipmentAggregate:
        return ShipmentAggregate(
            uid=bytes_to_uuid(entity.uid),
            order_uid=bytes_to_uuid(entity.order_uid),
            carrier_uid=bytes_to_uuid(entity.carrier_uid) if entity.carrier_uid else None,
            origin_address=entity.origin_address,
            destination_address=entity.destination_address,
            destination_city=entity.destination_city,
            destination_state=entity.destination_state,
            destination_country_code=entity.destination_country_code,
            destination_postal_code=entity.destination_postal_code,
            status=ShipmentStatus(entity.status),
            tracking_number=entity.tracking_number,
            weight_kg=entity.weight_kg,
            planned_pickup_at=entity.planned_pickup_at,
            planned_delivery_at=entity.planned_delivery_at,
            actual_pickup_at=entity.actual_pickup_at,
            actual_delivery_at=entity.actual_delivery_at,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )


_SORTABLE_COLUMNS = {
    "created_at": ShipmentEntity.created_at,
    "updated_at": ShipmentEntity.updated_at,
    "status": ShipmentEntity.status,
    "planned_pickup_at": ShipmentEntity.planned_pickup_at,
    "planned_delivery_at": ShipmentEntity.planned_delivery_at,
    "actual_pickup_at": ShipmentEntity.actual_pickup_at,
    "actual_delivery_at": ShipmentEntity.actual_delivery_at,
}
_DEFAULT_SORT_COLUMN = ShipmentEntity.created_at


def _parse_sort(sort: str) -> tuple[str, str]:
    parts = sort.split(":")
    field = parts[0] if len(parts) > 0 else "created_at"
    direction = parts[1] if len(parts) > 1 else "desc"
    return field, direction
