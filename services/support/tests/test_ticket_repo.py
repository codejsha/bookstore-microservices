from datetime import UTC, datetime, timedelta
from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus
from internal.domain.model.option.ticket_option import TicketFilterOption
from internal.infrastructure.adapter.mysql.ticket_category_repo import (
    MySQLTicketCategoryRepository,
)
from internal.infrastructure.adapter.mysql.ticket_repo import MySQLTicketRepository
from tests.conftest import make_category, make_ticket


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLTicketRepository:
    return MySQLTicketRepository(session_factory)


@pytest.fixture
def category_repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLTicketCategoryRepository:
    return MySQLTicketCategoryRepository(session_factory)


async def test_save_new_inserts_row(repo: MySQLTicketRepository) -> None:
    ticket = make_ticket()
    saved = await repo.save(ticket)
    assert saved.uid == ticket.uid
    assert saved.status == TicketStatus.OPEN


async def test_find_by_uid_roundtrips(repo: MySQLTicketRepository) -> None:
    ticket = make_ticket()
    await repo.save(ticket)
    found = await repo.find_by_uid(ticket.uid)
    assert found is not None
    assert found.subject == ticket.subject


async def test_find_by_uid_missing_is_none(repo: MySQLTicketRepository) -> None:
    assert await repo.find_by_uid(uuid4()) is None


async def test_find_id_by_uid_found_resolves_internal_id(repo: MySQLTicketRepository) -> None:
    ticket = make_ticket()
    await repo.save(ticket)
    ticket_id = await repo.find_id_by_uid(ticket.uid)
    assert ticket_id is not None
    assert ticket_id > 0


async def test_find_id_by_uid_missing_is_none(repo: MySQLTicketRepository) -> None:
    assert await repo.find_id_by_uid(uuid4()) is None


async def test_save_found_updates_row(repo: MySQLTicketRepository) -> None:
    ticket = make_ticket()
    await repo.save(ticket)
    ticket.subject = "Updated"
    ticket.priority = TicketPriority.URGENT
    ticket.status = TicketStatus.IN_PROGRESS
    ticket.updated_at = datetime.now(UTC)
    updated = await repo.save(ticket)
    assert updated.subject == "Updated"
    assert updated.priority == TicketPriority.URGENT
    assert updated.status == TicketStatus.IN_PROGRESS


async def test_find_all_by_customer_uid_matches_tickets(repo: MySQLTicketRepository) -> None:
    mine, theirs = uuid4(), uuid4()
    await repo.save(make_ticket(customer_uid=mine))
    await repo.save(make_ticket(customer_uid=theirs))
    items, total = await repo.find_all(TicketFilterOption(customer_uid=mine))
    assert total == 1
    assert items[0].customer_uid == mine


async def test_find_all_by_status_and_priority_matches_tickets(repo: MySQLTicketRepository) -> None:
    await repo.save(make_ticket(status=TicketStatus.OPEN, priority=TicketPriority.HIGH))
    await repo.save(make_ticket(status=TicketStatus.OPEN, priority=TicketPriority.LOW))
    await repo.save(make_ticket(status=TicketStatus.RESOLVED, priority=TicketPriority.HIGH))

    items, total = await repo.find_all(TicketFilterOption(status=TicketStatus.OPEN, priority=TicketPriority.HIGH))
    assert total == 1


async def test_find_all_by_category_uid_matches_tickets(
    repo: MySQLTicketRepository, category_repo: MySQLTicketCategoryRepository
) -> None:
    category = make_category()
    await category_repo.save(category)
    cat_id = await category_repo.find_id_by_uid(category.uid)
    await repo.save(make_ticket(category_id=cat_id))
    await repo.save(make_ticket())
    items, total = await repo.find_all(TicketFilterOption(category_uid=category.uid))
    assert total == 1
    assert items[0].category_id == cat_id


async def test_find_all_unknown_category_uid_is_empty(
    repo: MySQLTicketRepository,
) -> None:
    await repo.save(make_ticket())
    items, total = await repo.find_all(TicketFilterOption(category_uid=uuid4()))
    assert total == 0
    assert items == []


async def test_find_all_with_page_size_pages_newest_first_with_total(repo: MySQLTicketRepository) -> None:
    base = datetime.now(UTC)
    tickets = []
    for i in range(3):
        t = make_ticket()
        t.created_at = base + timedelta(seconds=i)
        t.updated_at = t.created_at
        await repo.save(t)
        tickets.append(t)
    items, total = await repo.find_all(TicketFilterOption(page=0, size=2))
    assert total == 3
    assert items[0].uid == tickets[2].uid


async def test_find_all_unknown_sort_field_falls_back_to_default_order(repo: MySQLTicketRepository) -> None:
    base = datetime.now(UTC)
    tickets = []
    for i in range(2):
        t = make_ticket()
        t.created_at = base + timedelta(seconds=i)
        t.updated_at = t.created_at
        await repo.save(t)
        tickets.append(t)
    items, total = await repo.find_all(TicketFilterOption(sort="metadata:desc"))
    assert total == 2
    assert items[0].uid == tickets[1].uid


async def test_find_by_uid_existing_category_resolves_category_uid(
    repo: MySQLTicketRepository, category_repo: MySQLTicketCategoryRepository
) -> None:
    category = make_category()
    await category_repo.save(category)
    cat_id = await category_repo.find_id_by_uid(category.uid)
    ticket = make_ticket(category_id=cat_id)
    await repo.save(ticket)

    found = await repo.find_by_uid(ticket.uid)
    assert found is not None
    assert found.category_uid == category.uid

    items, _ = await repo.find_all(TicketFilterOption(category_uid=category.uid))
    assert items[0].category_uid == category.uid


async def test_update_found_applies_mutator(repo: MySQLTicketRepository) -> None:
    ticket = make_ticket()
    await repo.save(ticket)

    def _mutate(t):
        t.subject = "Locked-update subject"
        t.status = TicketStatus.IN_PROGRESS
        t.updated_at = datetime.now(UTC)
        return t

    updated = await repo.update(ticket.uid, _mutate)
    assert updated is not None
    assert updated.subject == "Locked-update subject"
    assert updated.status == TicketStatus.IN_PROGRESS
    reloaded = await repo.find_by_uid(ticket.uid)
    assert reloaded is not None
    assert reloaded.subject == "Locked-update subject"


async def test_update_missing_is_none(repo: MySQLTicketRepository) -> None:
    assert await repo.update(uuid4(), lambda t: t) is None
