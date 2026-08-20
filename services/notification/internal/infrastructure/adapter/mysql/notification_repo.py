from collections.abc import Callable
from uuid import UUID

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import NotificationRepository
from internal.domain.aggregate.notification_aggregate import NotificationAggregate
from internal.domain.constant.notification_type import NotificationType
from internal.domain.model.option.notification_option import NotificationFilterOption
from internal.infrastructure.adapter.mysql.models import NotificationEntity
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes

_SORTABLE_COLUMNS = {
    "created_at": NotificationEntity.created_at,
    "updated_at": NotificationEntity.updated_at,
    "status": NotificationEntity.status,
    "sent_at": NotificationEntity.sent_at,
    "notification_type": NotificationEntity.notification_type,
    "channel": NotificationEntity.channel,
}
_DEFAULT_SORT_COLUMN = NotificationEntity.created_at


class MySQLNotificationRepository(NotificationRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, notification: NotificationAggregate) -> NotificationAggregate:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(NotificationEntity).where(NotificationEntity.uid == uuid_to_bytes(notification.uid))
                )
            ).scalar_one_or_none()

            if entity is None:
                entity = NotificationEntity(
                    uid=uuid_to_bytes(notification.uid),
                    user_uid=uuid_to_bytes(notification.user_uid),
                    order_uid=uuid_to_bytes(notification.order_uid) if notification.order_uid else None,
                    notification_type=notification.notification_type,
                    channel=notification.channel,
                    status=notification.status,
                    title=notification.title,
                    content=notification.content,
                    sent_at=notification.sent_at,
                    created_at=notification.created_at,
                    updated_at=notification.updated_at,
                )
                session.add(entity)
            else:
                entity.status = notification.status
                entity.sent_at = notification.sent_at
                entity.updated_at = notification.updated_at

            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    async def find_by_uid(self, uid: UUID) -> NotificationAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(NotificationEntity).where(
                        NotificationEntity.uid == uuid_to_bytes(uid),
                        NotificationEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def find_by_order_uid_and_type(
        self, order_uid: UUID, notification_type: NotificationType
    ) -> NotificationAggregate | None:
        async with self._session_factory() as session:
            entity = (
                (
                    await session.execute(
                        select(NotificationEntity)
                        .where(
                            NotificationEntity.order_uid == uuid_to_bytes(order_uid),
                            NotificationEntity.notification_type == notification_type,
                            NotificationEntity.deleted_at.is_(None),
                        )
                        .order_by(NotificationEntity.created_at.asc())
                    )
                )
                .scalars()
                .first()
            )
            return self._to_aggregate(entity) if entity else None

    async def find_all(self, option: NotificationFilterOption) -> tuple[list[NotificationAggregate], int]:
        async with self._session_factory() as session:
            query = select(NotificationEntity).where(NotificationEntity.deleted_at.is_(None))
            count_query = (
                select(func.count()).select_from(NotificationEntity).where(NotificationEntity.deleted_at.is_(None))
            )

            if option.user_uid is not None:
                uid_bytes = uuid_to_bytes(option.user_uid)
                query = query.where(NotificationEntity.user_uid == uid_bytes)
                count_query = count_query.where(NotificationEntity.user_uid == uid_bytes)
            if option.notification_type is not None:
                query = query.where(NotificationEntity.notification_type == option.notification_type)
                count_query = count_query.where(NotificationEntity.notification_type == option.notification_type)
            if option.channel is not None:
                query = query.where(NotificationEntity.channel == option.channel)
                count_query = count_query.where(NotificationEntity.channel == option.channel)
            if option.status is not None:
                query = query.where(NotificationEntity.status == option.status)
                count_query = count_query.where(NotificationEntity.status == option.status)

            sort_field, sort_dir = self._parse_sort(option.sort)
            order_col = _SORTABLE_COLUMNS.get(sort_field, _DEFAULT_SORT_COLUMN)
            query = query.order_by(order_col.desc() if sort_dir == "desc" else order_col.asc())

            query = query.offset(option.page * option.size).limit(option.size)

            total = (await session.execute(count_query)).scalar() or 0
            entities = (await session.execute(query)).scalars().all()
            return [self._to_aggregate(e) for e in entities], total

    async def update(
        self,
        uid: UUID,
        mutator: Callable[[NotificationAggregate], NotificationAggregate],
    ) -> NotificationAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(NotificationEntity)
                    .where(
                        NotificationEntity.uid == uuid_to_bytes(uid),
                        NotificationEntity.deleted_at.is_(None),
                    )
                    .with_for_update()
                )
            ).scalar_one_or_none()
            if entity is None:
                return None
            updated = mutator(self._to_aggregate(entity))
            entity.status = updated.status
            entity.sent_at = updated.sent_at
            entity.updated_at = updated.updated_at
            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    @staticmethod
    def _parse_sort(sort: str) -> tuple[str, str]:
        parts = sort.split(":")
        field = parts[0] if len(parts) > 0 else "created_at"
        direction = parts[1] if len(parts) > 1 else "desc"
        return field, direction

    @staticmethod
    def _to_aggregate(entity: NotificationEntity) -> NotificationAggregate:
        return NotificationAggregate(
            uid=bytes_to_uuid(entity.uid),
            user_uid=bytes_to_uuid(entity.user_uid),
            notification_type=entity.notification_type,
            channel=entity.channel,
            status=entity.status,
            title=entity.title,
            content=entity.content,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
            sent_at=entity.sent_at,
            order_uid=bytes_to_uuid(entity.order_uid) if entity.order_uid else None,
        )
