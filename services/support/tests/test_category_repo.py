from datetime import UTC, datetime
from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.infrastructure.adapter.mysql.ticket_category_repo import (
    MySQLTicketCategoryRepository,
)
from tests.conftest import make_category


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLTicketCategoryRepository:
    return MySQLTicketCategoryRepository(session_factory)


async def test_save_inserts_category(repo: MySQLTicketCategoryRepository) -> None:
    category = make_category()
    saved = await repo.save(category)
    assert saved.uid == category.uid


async def test_find_by_uid_returns_saved(repo: MySQLTicketCategoryRepository) -> None:
    category = make_category()
    await repo.save(category)
    found = await repo.find_by_uid(category.uid)
    assert found is not None
    assert found.name == category.name


async def test_find_by_uid_returns_none(repo: MySQLTicketCategoryRepository) -> None:
    assert await repo.find_by_uid(uuid4()) is None


async def test_find_id_by_uid_returns_id(repo: MySQLTicketCategoryRepository) -> None:
    category = make_category()
    await repo.save(category)
    cat_id = await repo.find_id_by_uid(category.uid)
    assert cat_id is not None and cat_id > 0


async def test_save_updates_existing(repo: MySQLTicketCategoryRepository) -> None:
    category = make_category(name="Original")
    await repo.save(category)
    category.name = "Updated"
    category.description = "new"
    category.updated_at = datetime.now(UTC)
    updated = await repo.save(category)
    assert updated.name == "Updated"
    assert updated.description == "new"


async def test_find_all_excludes_deleted(repo: MySQLTicketCategoryRepository) -> None:
    a = make_category(name="A")
    b = make_category(name="B")
    await repo.save(a)
    await repo.save(b)
    await repo.delete_by_uid(a.uid)
    found = await repo.find_all()
    names = {c.name for c in found}
    assert "A" not in names
    assert "B" in names


async def test_delete_returns_true_when_exists(repo: MySQLTicketCategoryRepository) -> None:
    category = make_category()
    await repo.save(category)
    assert await repo.delete_by_uid(category.uid) is True
    assert await repo.find_by_uid(category.uid) is None


async def test_delete_returns_false_when_missing(repo: MySQLTicketCategoryRepository) -> None:
    assert await repo.delete_by_uid(uuid4()) is False


async def test_save_with_parent_id(repo: MySQLTicketCategoryRepository) -> None:
    parent = make_category(name="Parent")
    await repo.save(parent)
    parent_id = await repo.find_id_by_uid(parent.uid)
    child = make_category(name="Child", parent_id=parent_id)
    saved = await repo.save(child)
    assert saved.parent_id == parent_id


async def test_resolves_parent_uid(repo: MySQLTicketCategoryRepository) -> None:
    parent = make_category(name="Parent")
    await repo.save(parent)
    parent_id = await repo.find_id_by_uid(parent.uid)
    child = make_category(name="Child", parent_id=parent_id)
    saved = await repo.save(child)
    assert saved.parent_uid == parent.uid

    found = await repo.find_by_uid(child.uid)
    assert found is not None
    assert found.parent_uid == parent.uid

    items = await repo.find_all()
    child_row = next(c for c in items if c.uid == child.uid)
    assert child_row.parent_uid == parent.uid
