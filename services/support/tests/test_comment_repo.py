import pytest
import pytest_asyncio
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.constant.author_role import AuthorRole
from internal.infrastructure.adapter.mysql.ticket_comment_repo import (
    MySQLTicketCommentRepository,
)
from internal.infrastructure.adapter.mysql.ticket_repo import MySQLTicketRepository
from tests.conftest import make_comment, make_ticket


@pytest.fixture
def comment_repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLTicketCommentRepository:
    return MySQLTicketCommentRepository(session_factory)


@pytest_asyncio.fixture
async def ticket_id(session_factory: async_sessionmaker[AsyncSession]) -> int:
    repo = MySQLTicketRepository(session_factory)
    ticket = make_ticket()
    await repo.save(ticket)
    return await repo.find_id_by_uid(ticket.uid)


async def test_save_new_inserts_row(comment_repo: MySQLTicketCommentRepository, ticket_id: int) -> None:
    comment = make_comment(ticket_id=ticket_id)
    saved = await comment_repo.save(comment)
    assert saved.uid == comment.uid


async def test_find_by_uid_roundtrips(comment_repo: MySQLTicketCommentRepository, ticket_id: int) -> None:
    comment = make_comment(ticket_id=ticket_id)
    await comment_repo.save(comment)
    found = await comment_repo.find_by_uid(comment.uid)
    assert found is not None
    assert found.body == comment.body


async def test_find_by_ticket_id_include_internal_omitted_keeps_internal_comments(
    comment_repo: MySQLTicketCommentRepository, ticket_id: int
) -> None:
    await comment_repo.save(make_comment(ticket_id=ticket_id, internal=False))
    await comment_repo.save(make_comment(ticket_id=ticket_id, internal=True))
    items = await comment_repo.find_by_ticket_id(ticket_id)
    assert len(items) == 2


async def test_find_by_ticket_id_include_internal_false_excludes_internal_comments(
    comment_repo: MySQLTicketCommentRepository, ticket_id: int
) -> None:
    await comment_repo.save(make_comment(ticket_id=ticket_id, internal=False))
    await comment_repo.save(make_comment(ticket_id=ticket_id, internal=True))
    items = await comment_repo.find_by_ticket_id(ticket_id, include_internal=False)
    assert len(items) == 1
    assert items[0].internal is False


async def test_find_by_ticket_id_orders_comments_oldest_first(
    comment_repo: MySQLTicketCommentRepository, ticket_id: int
) -> None:
    from datetime import UTC, datetime, timedelta

    base = datetime.now(UTC)
    earlier = make_comment(ticket_id=ticket_id, author_role=AuthorRole.CUSTOMER)
    earlier.created_at = base
    earlier.updated_at = base
    later = make_comment(ticket_id=ticket_id, author_role=AuthorRole.AGENT)
    later.created_at = base + timedelta(minutes=5)
    later.updated_at = later.created_at
    await comment_repo.save(later)
    await comment_repo.save(earlier)

    items = await comment_repo.find_by_ticket_id(ticket_id)
    assert items[0].author_role == AuthorRole.CUSTOMER
    assert items[1].author_role == AuthorRole.AGENT
