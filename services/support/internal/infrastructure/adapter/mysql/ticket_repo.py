from collections.abc import Callable
from uuid import UUID

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import TicketRepository
from internal.domain.aggregate.ticket_aggregate import TicketAggregate
from internal.domain.model.option.ticket_option import TicketFilterOption
from internal.infrastructure.adapter.mysql.models import TicketCategoryEntity, TicketEntity
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes

_SORTABLE_COLUMNS = {
    "created_at": TicketEntity.created_at,
    "updated_at": TicketEntity.updated_at,
    "status": TicketEntity.status,
    "priority": TicketEntity.priority,
    "subject": TicketEntity.subject,
    "resolved_at": TicketEntity.resolved_at,
}
_DEFAULT_SORT_COLUMN = TicketEntity.created_at


class MySQLTicketRepository(TicketRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, ticket: TicketAggregate) -> TicketAggregate:
        async with self._session_factory() as session:
            entity = (
                await session.execute(select(TicketEntity).where(TicketEntity.uid == uuid_to_bytes(ticket.uid)))
            ).scalar_one_or_none()

            if entity is None:
                entity = TicketEntity(
                    uid=uuid_to_bytes(ticket.uid),
                    customer_uid=uuid_to_bytes(ticket.customer_uid),
                    category_id=ticket.category_id,
                    assignee_uid=uuid_to_bytes(ticket.assignee_uid) if ticket.assignee_uid else None,
                    subject=ticket.subject,
                    description=ticket.description,
                    status=ticket.status,
                    priority=ticket.priority,
                    resolved_at=ticket.resolved_at,
                    created_at=ticket.created_at,
                    updated_at=ticket.updated_at,
                )
                session.add(entity)
            else:
                entity.customer_uid = uuid_to_bytes(ticket.customer_uid)
                entity.category_id = ticket.category_id
                entity.assignee_uid = uuid_to_bytes(ticket.assignee_uid) if ticket.assignee_uid else None
                entity.subject = ticket.subject
                entity.description = ticket.description
                entity.status = ticket.status
                entity.priority = ticket.priority
                entity.resolved_at = ticket.resolved_at
                entity.updated_at = ticket.updated_at

            await session.commit()
            await session.refresh(entity)
            category_uid = await self._resolve_category_uid(session, entity.category_id)
            return self._to_aggregate(entity, category_uid)

    async def find_by_uid(self, uid: UUID) -> TicketAggregate | None:
        async with self._session_factory() as session:
            row = (
                await session.execute(
                    select(TicketEntity, TicketCategoryEntity.uid)
                    .outerjoin(TicketCategoryEntity, TicketEntity.category_id == TicketCategoryEntity.id)
                    .where(
                        TicketEntity.uid == uuid_to_bytes(uid),
                        TicketEntity.deleted_at.is_(None),
                    )
                )
            ).first()
            if row is None:
                return None
            entity, category_uid_bytes = row
            return self._to_aggregate(entity, category_uid_bytes)

    async def find_id_by_uid(self, uid: UUID) -> int | None:
        async with self._session_factory() as session:
            row = (
                await session.execute(
                    select(TicketEntity.id).where(
                        TicketEntity.uid == uuid_to_bytes(uid),
                        TicketEntity.deleted_at.is_(None),
                    )
                )
            ).first()
            return row[0] if row else None

    async def find_all(self, option: TicketFilterOption) -> tuple[list[TicketAggregate], int]:
        async with self._session_factory() as session:
            query = (
                select(TicketEntity, TicketCategoryEntity.uid)
                .outerjoin(TicketCategoryEntity, TicketEntity.category_id == TicketCategoryEntity.id)
                .where(TicketEntity.deleted_at.is_(None))
            )
            count_query = select(func.count()).select_from(TicketEntity).where(TicketEntity.deleted_at.is_(None))

            if option.customer_uid is not None:
                customer_uid = uuid_to_bytes(option.customer_uid)
                query = query.where(TicketEntity.customer_uid == customer_uid)
                count_query = count_query.where(TicketEntity.customer_uid == customer_uid)
            if option.assignee_uid is not None:
                assignee_uid = uuid_to_bytes(option.assignee_uid)
                query = query.where(TicketEntity.assignee_uid == assignee_uid)
                count_query = count_query.where(TicketEntity.assignee_uid == assignee_uid)
            if option.status is not None:
                query = query.where(TicketEntity.status == option.status)
                count_query = count_query.where(TicketEntity.status == option.status)
            if option.priority is not None:
                query = query.where(TicketEntity.priority == option.priority)
                count_query = count_query.where(TicketEntity.priority == option.priority)
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
                query = query.where(TicketEntity.category_id == cat_id)
                count_query = count_query.where(TicketEntity.category_id == cat_id)

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
        mutator: Callable[[TicketAggregate], TicketAggregate],
    ) -> TicketAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(TicketEntity)
                    .where(
                        TicketEntity.uid == uuid_to_bytes(uid),
                        TicketEntity.deleted_at.is_(None),
                    )
                    .with_for_update()
                )
            ).scalar_one_or_none()
            if entity is None:
                return None
            category_uid = await self._resolve_category_uid(session, entity.category_id)
            updated = mutator(self._to_aggregate(entity, category_uid))
            entity.customer_uid = uuid_to_bytes(updated.customer_uid)
            entity.category_id = updated.category_id
            entity.assignee_uid = uuid_to_bytes(updated.assignee_uid) if updated.assignee_uid else None
            entity.subject = updated.subject
            entity.description = updated.description
            entity.status = updated.status
            entity.priority = updated.priority
            entity.resolved_at = updated.resolved_at
            entity.updated_at = updated.updated_at
            entity.version = entity.version + 1
            await session.commit()
            await session.refresh(entity)
            category_uid = await self._resolve_category_uid(session, entity.category_id)
            return self._to_aggregate(entity, category_uid)

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
        field = parts[0] if len(parts) > 0 else "created_at"
        direction = parts[1] if len(parts) > 1 else "desc"
        return field, direction

    @staticmethod
    def _to_aggregate(entity: TicketEntity, category_uid_bytes: bytes | None = None) -> TicketAggregate:
        return TicketAggregate(
            uid=bytes_to_uuid(entity.uid),
            customer_uid=bytes_to_uuid(entity.customer_uid),
            category_id=entity.category_id,
            category_uid=bytes_to_uuid(category_uid_bytes) if category_uid_bytes else None,
            assignee_uid=bytes_to_uuid(entity.assignee_uid) if entity.assignee_uid else None,
            subject=entity.subject,
            description=entity.description,
            status=entity.status,
            priority=entity.priority,
            resolved_at=entity.resolved_at,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )
