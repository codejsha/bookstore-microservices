from datetime import UTC, datetime
from uuid import UUID

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import TemplateRepository
from internal.domain.aggregate.template_aggregate import TemplateAggregate
from internal.domain.model.option.notification_option import TemplateFilterOption
from internal.infrastructure.adapter.mysql.models import TemplateEntity
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes

_SORTABLE_COLUMNS = {
    "created_at": TemplateEntity.created_at,
    "updated_at": TemplateEntity.updated_at,
    "notification_type": TemplateEntity.notification_type,
    "channel": TemplateEntity.channel,
}
_DEFAULT_SORT_COLUMN = TemplateEntity.created_at


class MySQLTemplateRepository(TemplateRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, template: TemplateAggregate) -> TemplateAggregate:
        async with self._session_factory() as session:
            entity = (
                await session.execute(select(TemplateEntity).where(TemplateEntity.uid == uuid_to_bytes(template.uid)))
            ).scalar_one_or_none()

            if entity is None:
                entity = TemplateEntity(
                    uid=uuid_to_bytes(template.uid),
                    notification_type=template.notification_type,
                    channel=template.channel,
                    title_template=template.title_template,
                    content_template=template.content_template,
                    created_at=template.created_at,
                    updated_at=template.updated_at,
                )
                session.add(entity)
            else:
                entity.notification_type = template.notification_type
                entity.channel = template.channel
                entity.title_template = template.title_template
                entity.content_template = template.content_template
                entity.updated_at = template.updated_at

            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    async def find_by_uid(self, uid: UUID) -> TemplateAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(TemplateEntity).where(
                        TemplateEntity.uid == uuid_to_bytes(uid),
                        TemplateEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def find_all(self, option: TemplateFilterOption) -> tuple[list[TemplateAggregate], int]:
        async with self._session_factory() as session:
            query = select(TemplateEntity).where(TemplateEntity.deleted_at.is_(None))
            count_query = select(func.count()).select_from(TemplateEntity).where(TemplateEntity.deleted_at.is_(None))

            if option.notification_type is not None:
                query = query.where(TemplateEntity.notification_type == option.notification_type)
                count_query = count_query.where(TemplateEntity.notification_type == option.notification_type)
            if option.channel is not None:
                query = query.where(TemplateEntity.channel == option.channel)
                count_query = count_query.where(TemplateEntity.channel == option.channel)

            sort_field, sort_dir = self._parse_sort(option.sort)
            order_col = _SORTABLE_COLUMNS.get(sort_field, _DEFAULT_SORT_COLUMN)
            query = query.order_by(order_col.desc() if sort_dir == "desc" else order_col.asc())

            query = query.offset(option.page * option.size).limit(option.size)

            total = (await session.execute(count_query)).scalar() or 0
            entities = (await session.execute(query)).scalars().all()
            return [self._to_aggregate(e) for e in entities], total

    async def delete_by_uid(self, uid: UUID) -> bool:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(TemplateEntity).where(
                        TemplateEntity.uid == uuid_to_bytes(uid),
                        TemplateEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            if entity is None:
                return False
            entity.deleted_at = datetime.now(UTC)
            await session.commit()
            return True

    @staticmethod
    def _parse_sort(sort: str) -> tuple[str, str]:
        parts = sort.split(":")
        field = parts[0] if len(parts) > 0 else "created_at"
        direction = parts[1] if len(parts) > 1 else "desc"
        return field, direction

    @staticmethod
    def _to_aggregate(entity: TemplateEntity) -> TemplateAggregate:
        return TemplateAggregate(
            uid=bytes_to_uuid(entity.uid),
            notification_type=entity.notification_type,
            channel=entity.channel,
            title_template=entity.title_template,
            content_template=entity.content_template,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )
