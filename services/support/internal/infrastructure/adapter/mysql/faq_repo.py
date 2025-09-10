from collections.abc import Callable
from datetime import UTC, datetime
from uuid import UUID

from sqlalchemy import func, or_, select, update
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import FaqRepository
from internal.domain.aggregate.faq_aggregate import FaqAggregate
from internal.domain.model.option.faq_option import FaqSearchOption
from internal.infrastructure.adapter.mysql.models import FaqEntity, TicketCategoryEntity
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes

_SORTABLE_COLUMNS = {
    "view_count": FaqEntity.view_count,
    "created_at": FaqEntity.created_at,
    "updated_at": FaqEntity.updated_at,
    "question": FaqEntity.question,
    "published": FaqEntity.published,
}
_DEFAULT_SORT_COLUMN = FaqEntity.view_count


class MySQLFaqRepository(FaqRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, faq: FaqAggregate) -> FaqAggregate:
        async with self._session_factory() as session:
            entity = (
                await session.execute(select(FaqEntity).where(FaqEntity.uid == uuid_to_bytes(faq.uid)))
            ).scalar_one_or_none()

            if entity is None:
                entity = FaqEntity(
                    uid=uuid_to_bytes(faq.uid),
                    category_id=faq.category_id,
                    question=faq.question,
                    answer=faq.answer,
                    view_count=faq.view_count,
                    published=faq.published,
                    created_at=faq.created_at,
                    updated_at=faq.updated_at,
                )
                session.add(entity)
            else:
                entity.category_id = faq.category_id
                entity.question = faq.question
                entity.answer = faq.answer
                entity.published = faq.published
                entity.updated_at = faq.updated_at

            await session.commit()
            await session.refresh(entity)
            category_uid = await self._resolve_category_uid(session, entity.category_id)
            return self._to_aggregate(entity, category_uid)

    async def find_by_uid(self, uid: UUID) -> FaqAggregate | None:
        async with self._session_factory() as session:
            row = (
                await session.execute(
                    select(FaqEntity, TicketCategoryEntity.uid)
                    .outerjoin(TicketCategoryEntity, FaqEntity.category_id == TicketCategoryEntity.id)
                    .where(
                        FaqEntity.uid == uuid_to_bytes(uid),
                        FaqEntity.deleted_at.is_(None),
                    )
                )
            ).first()
            if row is None:
                return None
            entity, category_uid_bytes = row
            return self._to_aggregate(entity, category_uid_bytes)

    async def search(self, option: FaqSearchOption) -> tuple[list[FaqAggregate], int]:
        async with self._session_factory() as session:
            query = (
                select(FaqEntity, TicketCategoryEntity.uid)
                .outerjoin(TicketCategoryEntity, FaqEntity.category_id == TicketCategoryEntity.id)
                .where(FaqEntity.deleted_at.is_(None))
            )
            count_query = select(func.count()).select_from(FaqEntity).where(FaqEntity.deleted_at.is_(None))

            if option.published_only:
                query = query.where(FaqEntity.published.is_(True))
                count_query = count_query.where(FaqEntity.published.is_(True))

            if option.category_uid is not None:
                cat_id = (
                    await session.execute(
                        select(TicketCategoryEntity.id).where(
                            TicketCategoryEntity.uid == uuid_to_bytes(option.category_uid)
                        )
                    )
                ).scalar()
                if cat_id is None:
                    return [], 0
                query = query.where(FaqEntity.category_id == cat_id)
                count_query = count_query.where(FaqEntity.category_id == cat_id)

            if option.query:
                pattern = f"%{option.query}%"
                query = query.where(or_(FaqEntity.question.like(pattern), FaqEntity.answer.like(pattern)))
                count_query = count_query.where(or_(FaqEntity.question.like(pattern), FaqEntity.answer.like(pattern)))

            sort_field, sort_dir = self._parse_sort(option.sort)
            order_col = _SORTABLE_COLUMNS.get(sort_field, _DEFAULT_SORT_COLUMN)
            query = query.order_by(order_col.desc() if sort_dir == "desc" else order_col.asc())

            query = query.offset(option.page * option.size).limit(option.size)

            total = (await session.execute(count_query)).scalar() or 0
            rows = (await session.execute(query)).all()
            return [self._to_aggregate(entity, category_uid_bytes) for entity, category_uid_bytes in rows], total

    async def update(
        self,
        uid: UUID,
        mutator: Callable[[FaqAggregate], FaqAggregate],
    ) -> FaqAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(FaqEntity)
                    .where(
                        FaqEntity.uid == uuid_to_bytes(uid),
                        FaqEntity.deleted_at.is_(None),
                    )
                    .with_for_update()
                )
            ).scalar_one_or_none()
            if entity is None:
                return None
            category_uid = await self._resolve_category_uid(session, entity.category_id)
            updated = mutator(self._to_aggregate(entity, category_uid))
            entity.category_id = updated.category_id
            entity.question = updated.question
            entity.answer = updated.answer
            entity.published = updated.published
            entity.updated_at = updated.updated_at
            entity.version = entity.version + 1
            await session.commit()
            await session.refresh(entity)
            category_uid = await self._resolve_category_uid(session, entity.category_id)
            return self._to_aggregate(entity, category_uid)

    async def increment_view_count(self, uid: UUID) -> None:
        async with self._session_factory() as session:
            await session.execute(
                update(FaqEntity).where(FaqEntity.uid == uuid_to_bytes(uid)).values(view_count=FaqEntity.view_count + 1)
            )
            await session.commit()

    async def delete_by_uid(self, uid: UUID) -> bool:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(FaqEntity).where(
                        FaqEntity.uid == uuid_to_bytes(uid),
                        FaqEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            if entity is None:
                return False
            entity.deleted_at = datetime.now(UTC)
            await session.commit()
            return True

    @staticmethod
    async def _resolve_category_uid(session: AsyncSession, category_id: int | None) -> bytes | None:
        if category_id is None:
            return None
        return (
            await session.execute(select(TicketCategoryEntity.uid).where(TicketCategoryEntity.id == category_id))
        ).scalar_one_or_none()

    @staticmethod
    def _parse_sort(sort: str) -> tuple[str, str]:
        parts = sort.split(":")
        field = parts[0] if len(parts) > 0 else "view_count"
        direction = parts[1] if len(parts) > 1 else "desc"
        return field, direction

    @staticmethod
    def _to_aggregate(entity: FaqEntity, category_uid_bytes: bytes | None = None) -> FaqAggregate:
        return FaqAggregate(
            uid=bytes_to_uuid(entity.uid),
            category_id=entity.category_id,
            category_uid=bytes_to_uuid(category_uid_bytes) if category_uid_bytes else None,
            question=entity.question,
            answer=entity.answer,
            view_count=entity.view_count,
            published=entity.published,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )
