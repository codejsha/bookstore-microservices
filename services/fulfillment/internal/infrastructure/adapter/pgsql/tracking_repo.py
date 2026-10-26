from uuid import UUID

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import TrackingRepository
from internal.domain.aggregate.tracking_aggregate import TrackingAggregate
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.infrastructure.adapter.pgsql.models import TrackingEntity


def to_tracking_entity(tracking: TrackingAggregate) -> TrackingEntity:
    return TrackingEntity(
        uid=tracking.uid,
        shipment_uid=tracking.shipment_uid,
        status=tracking.status,
        location=tracking.location,
        description=tracking.description,
        occurred_at=tracking.occurred_at,
        created_at=tracking.created_at,
    )


class PgsqlTrackingRepository(TrackingRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, tracking: TrackingAggregate) -> TrackingAggregate:
        async with self._session_factory() as session:
            entity = to_tracking_entity(tracking)
            session.add(entity)
            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    async def find_by_shipment_uid(self, shipment_uid: UUID) -> list[TrackingAggregate]:
        async with self._session_factory() as session:
            query = (
                select(TrackingEntity)
                .where(TrackingEntity.shipment_uid == shipment_uid)
                .order_by(TrackingEntity.occurred_at.asc())
            )
            entities = (await session.execute(query)).scalars().all()
            return [self._to_aggregate(e) for e in entities]

    @staticmethod
    def _to_aggregate(entity: TrackingEntity) -> TrackingAggregate:
        return TrackingAggregate(
            uid=entity.uid,
            shipment_uid=entity.shipment_uid,
            status=ShipmentStatus(entity.status),
            location=entity.location,
            description=entity.description,
            occurred_at=entity.occurred_at,
            created_at=entity.created_at,
        )
