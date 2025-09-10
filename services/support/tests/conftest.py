from collections.abc import AsyncIterator
from datetime import UTC, datetime
from unittest.mock import MagicMock
from uuid import UUID, uuid4

import pytest
import pytest_asyncio
from sqlalchemy import BigInteger, event
from sqlalchemy.engine import Engine
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.ext.compiler import compiles
from sqlalchemy.pool import StaticPool

from internal.application.port.repo.repos import (
    FaqRepository,
    TicketCategoryRepository,
    TicketCommentRepository,
    TicketRepository,
)
from internal.domain.aggregate.faq_aggregate import FaqAggregate
from internal.domain.aggregate.ticket_aggregate import TicketAggregate
from internal.domain.aggregate.ticket_category_aggregate import TicketCategoryAggregate
from internal.domain.aggregate.ticket_comment_aggregate import TicketCommentAggregate
from internal.domain.constant.author_role import AuthorRole
from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus
from internal.infrastructure.adapter.mysql.models import Base


@compiles(BigInteger, "sqlite")
def _compile_bigint_sqlite(_type, _compiler, **_kw):
    return "INTEGER"


@event.listens_for(Engine, "connect")
def _enable_sqlite_fk(dbapi_connection, _record):
    if dbapi_connection.__class__.__module__.startswith("sqlite3") or dbapi_connection.__class__.__module__.startswith(
        "aiosqlite"
    ):
        cursor = dbapi_connection.cursor()
        cursor.execute("PRAGMA foreign_keys=ON")
        cursor.close()


@pytest.fixture
def ticket_repo() -> MagicMock:
    repo = MagicMock(spec=TicketRepository)
    repo.save.side_effect = lambda t: t

    async def _update(uid, mutator):
        ticket = repo.find_by_uid.return_value
        if ticket is None:
            return None
        return mutator(ticket)

    repo.update.side_effect = _update
    return repo


@pytest.fixture
def comment_repo() -> MagicMock:
    repo = MagicMock(spec=TicketCommentRepository)
    repo.save.side_effect = lambda c: c
    return repo


@pytest.fixture
def category_repo() -> MagicMock:
    repo = MagicMock(spec=TicketCategoryRepository)
    repo.save.side_effect = lambda c: c
    return repo


@pytest.fixture
def faq_repo() -> MagicMock:
    repo = MagicMock(spec=FaqRepository)
    repo.save.side_effect = lambda f: f

    async def _update(uid, mutator):
        faq = repo.find_by_uid.return_value
        if faq is None:
            return None
        return mutator(faq)

    repo.update.side_effect = _update
    return repo


@pytest_asyncio.fixture
async def session_factory() -> AsyncIterator[async_sessionmaker[AsyncSession]]:
    engine = create_async_engine(
        "sqlite+aiosqlite:///:memory:",
        connect_args={"check_same_thread": False},
        poolclass=StaticPool,
    )
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    factory = async_sessionmaker(bind=engine, expire_on_commit=False)
    try:
        yield factory
    finally:
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.drop_all)
        await engine.dispose()


def make_ticket(
    *,
    uid: UUID | None = None,
    customer_uid: UUID | None = None,
    status: TicketStatus = TicketStatus.OPEN,
    priority: TicketPriority = TicketPriority.MEDIUM,
    category_id: int | None = None,
    assignee_uid: UUID | None = None,
    resolved_at: datetime | None = None,
) -> TicketAggregate:
    now = datetime.now(UTC)
    return TicketAggregate(
        uid=uid or uuid4(),
        customer_uid=customer_uid or uuid4(),
        category_id=category_id,
        assignee_uid=assignee_uid,
        subject="subject",
        description="description",
        status=status,
        priority=priority,
        resolved_at=resolved_at,
        created_at=now,
        updated_at=now,
    )


def make_category(
    *,
    uid: UUID | None = None,
    name: str = "Billing",
    description: str | None = "About billing",
    parent_id: int | None = None,
) -> TicketCategoryAggregate:
    now = datetime.now(UTC)
    return TicketCategoryAggregate(
        uid=uid or uuid4(),
        name=name,
        description=description,
        parent_id=parent_id,
        created_at=now,
        updated_at=now,
    )


def make_faq(
    *,
    uid: UUID | None = None,
    published: bool = True,
    view_count: int = 0,
    category_id: int | None = None,
) -> FaqAggregate:
    now = datetime.now(UTC)
    return FaqAggregate(
        uid=uid or uuid4(),
        category_id=category_id,
        question="Q?",
        answer="A.",
        view_count=view_count,
        published=published,
        created_at=now,
        updated_at=now,
    )


def make_comment(
    *,
    ticket_id: int = 1,
    author_role: AuthorRole = AuthorRole.CUSTOMER,
    internal: bool = False,
) -> TicketCommentAggregate:
    now = datetime.now(UTC)
    return TicketCommentAggregate(
        uid=uuid4(),
        ticket_id=ticket_id,
        author_uid=uuid4(),
        author_role=author_role,
        body="body",
        internal=internal,
        created_at=now,
        updated_at=now,
    )
