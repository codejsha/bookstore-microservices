from uuid import UUID

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import TicketCommentRepository
from internal.domain.aggregate.ticket_comment_aggregate import TicketCommentAggregate
from internal.infrastructure.adapter.mysql.models import TicketCommentEntity
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes


class MySQLTicketCommentRepository(TicketCommentRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def save(self, comment: TicketCommentAggregate) -> TicketCommentAggregate:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(TicketCommentEntity).where(TicketCommentEntity.uid == uuid_to_bytes(comment.uid))
                )
            ).scalar_one_or_none()

            if entity is None:
                entity = TicketCommentEntity(
                    uid=uuid_to_bytes(comment.uid),
                    ticket_id=comment.ticket_id,
                    author_uid=uuid_to_bytes(comment.author_uid),
                    author_role=comment.author_role,
                    body=comment.body,
                    internal=comment.internal,
                    created_at=comment.created_at,
                    updated_at=comment.updated_at,
                )
                session.add(entity)
            else:
                entity.body = comment.body
                entity.internal = comment.internal
                entity.updated_at = comment.updated_at
                entity.version = entity.version + 1

            await session.commit()
            await session.refresh(entity)
            return self._to_aggregate(entity)

    async def find_by_uid(self, uid: UUID) -> TicketCommentAggregate | None:
        async with self._session_factory() as session:
            entity = (
                await session.execute(
                    select(TicketCommentEntity).where(
                        TicketCommentEntity.uid == uuid_to_bytes(uid),
                        TicketCommentEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            return self._to_aggregate(entity) if entity else None

    async def find_by_ticket_id(self, ticket_id: int, include_internal: bool = True) -> list[TicketCommentAggregate]:
        async with self._session_factory() as session:
            query = (
                select(TicketCommentEntity)
                .where(
                    TicketCommentEntity.ticket_id == ticket_id,
                    TicketCommentEntity.deleted_at.is_(None),
                )
                .order_by(TicketCommentEntity.created_at.asc())
            )
            if not include_internal:
                query = query.where(TicketCommentEntity.internal.is_(False))
            entities = (await session.execute(query)).scalars().all()
            return [self._to_aggregate(e) for e in entities]

    @staticmethod
    def _to_aggregate(entity: TicketCommentEntity) -> TicketCommentAggregate:
        return TicketCommentAggregate(
            uid=bytes_to_uuid(entity.uid),
            ticket_id=entity.ticket_id,
            author_uid=bytes_to_uuid(entity.author_uid),
            author_role=entity.author_role,
            body=entity.body,
            internal=entity.internal,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
        )
