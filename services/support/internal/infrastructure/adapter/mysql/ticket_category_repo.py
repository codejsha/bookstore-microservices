from datetime import UTC, datetime
from uuid import UUID

from sqlalchemy import select
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker
from sqlalchemy.orm import aliased

from internal.application.port.repo.repos import TicketCategoryRepository
from internal.domain.aggregate.ticket_category_aggregate import TicketCategoryAggregate
from internal.domain.error import ConflictError
from internal.infrastructure.adapter.mysql.models import TicketCategoryEntity
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes

_ParentCategory = aliased(TicketCategoryEntity)


class MySQLTicketCategoryRepository(TicketCategoryRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, category: TicketCategoryAggregate) -> TicketCategoryAggregate:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(TicketCategoryEntity).where(TicketCategoryEntity.uid == uuid_to_bytes(category.uid))
                )
            ).scalar_one_or_none()

            if entity is None:
                entity = TicketCategoryEntity(
                    uid=uuid_to_bytes(category.uid),
                    name=category.name,
                    description=category.description,
                    parent_id=category.parent_id,
                    created_at=category.created_at,
                    updated_at=category.updated_at,
                )
                session.add(entity)
            else:
                entity.name = category.name
                entity.description = category.description
                entity.parent_id = category.parent_id
                entity.updated_at = category.updated_at
                entity.version = entity.version + 1

            try:
                await session.commit()
            except IntegrityError as exc:
                await session.rollback()
                raise ConflictError("Category name already exists") from exc
            await session.refresh(entity)
            parent_uid = await self._resolve_parent_uid(session, entity.parent_id)
            return self._to_aggregate(entity, parent_uid)

    async def find_by_uid(self, uid: UUID) -> TicketCategoryAggregate | None:
        async with self._session_factory() as session:
            row = (
                await session.execute(
                    select(TicketCategoryEntity, _ParentCategory.uid)
                    .outerjoin(_ParentCategory, TicketCategoryEntity.parent_id == _ParentCategory.id)
                    .where(
                        TicketCategoryEntity.uid == uuid_to_bytes(uid),
                        TicketCategoryEntity.deleted_at.is_(None),
                    )
                )
            ).first()
            if row is None:
                return None
            entity, parent_uid_bytes = row
            return self._to_aggregate(entity, parent_uid_bytes)

    async def find_id_by_uid(self, uid: UUID) -> int | None:
        async with self._session_factory() as session:
            row = (
                await session.execute(
                    select(TicketCategoryEntity.id).where(
                        TicketCategoryEntity.uid == uuid_to_bytes(uid),
                        TicketCategoryEntity.deleted_at.is_(None),
                    )
                )
            ).first()
            return row[0] if row else None

    async def find_by_name(self, name: str) -> TicketCategoryAggregate | None:
        async with self._session_factory() as session:
            row = (
                await session.execute(
                    select(TicketCategoryEntity, _ParentCategory.uid)
                    .outerjoin(_ParentCategory, TicketCategoryEntity.parent_id == _ParentCategory.id)
                    .where(TicketCategoryEntity.name == name)
                )
            ).first()
            if row is None:
                return None
            entity, parent_uid_bytes = row
            return self._to_aggregate(entity, parent_uid_bytes)

    async def find_all(self) -> list[TicketCategoryAggregate]:
        async with self._session_factory() as session:
            query = (
                select(TicketCategoryEntity, _ParentCategory.uid)
                .outerjoin(_ParentCategory, TicketCategoryEntity.parent_id == _ParentCategory.id)
                .where(TicketCategoryEntity.deleted_at.is_(None))
            )
            rows = (await session.execute(query)).all()
            return [self._to_aggregate(entity, parent_uid_bytes) for entity, parent_uid_bytes in rows]

    async def delete_by_uid(self, uid: UUID) -> bool:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(TicketCategoryEntity).where(
                        TicketCategoryEntity.uid == uuid_to_bytes(uid),
                        TicketCategoryEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            if entity is None:
                return False
            entity.deleted_at = datetime.now(UTC)
            await session.commit()
            return True

    @staticmethod
    async def _resolve_parent_uid(session: AsyncSession, parent_id: int | None) -> bytes | None:
        if parent_id is None:
            return None
        return (
            await session.execute(select(TicketCategoryEntity.uid).where(TicketCategoryEntity.id == parent_id))
        ).scalar_one_or_none()

    @staticmethod
    def _to_aggregate(entity: TicketCategoryEntity, parent_uid_bytes: bytes | None = None) -> TicketCategoryAggregate:
        return TicketCategoryAggregate(
            uid=bytes_to_uuid(entity.uid),
            name=entity.name,
            description=entity.description,
            parent_id=entity.parent_id,
            parent_uid=bytes_to_uuid(parent_uid_bytes) if parent_uid_bytes else None,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )
