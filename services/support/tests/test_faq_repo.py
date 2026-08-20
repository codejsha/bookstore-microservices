from datetime import UTC, datetime
from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.model.option.faq_option import FaqSearchOption
from internal.infrastructure.adapter.mysql.faq_repo import MySQLFaqRepository
from internal.infrastructure.adapter.mysql.ticket_category_repo import (
    MySQLTicketCategoryRepository,
)
from tests.conftest import make_category, make_faq


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLFaqRepository:
    return MySQLFaqRepository(session_factory)


@pytest.fixture
def category_repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLTicketCategoryRepository:
    return MySQLTicketCategoryRepository(session_factory)


async def test_save_inserts_faq(repo: MySQLFaqRepository) -> None:
    faq = make_faq()
    saved = await repo.save(faq)
    assert saved.uid == faq.uid


async def test_find_by_uid_returns_saved(repo: MySQLFaqRepository) -> None:
    faq = make_faq()
    await repo.save(faq)
    found = await repo.find_by_uid(faq.uid)
    assert found is not None
    assert found.question == faq.question


async def test_find_by_uid_returns_none(repo: MySQLFaqRepository) -> None:
    assert await repo.find_by_uid(uuid4()) is None


async def test_datetime_round_trips_as_tz_aware_utc(repo: MySQLFaqRepository) -> None:
    faq = make_faq()
    await repo.save(faq)

    found = await repo.find_by_uid(faq.uid)
    assert found is not None
    assert found.created_at.tzinfo is not None
    assert found.created_at.utcoffset() == UTC.utcoffset(None)
    assert found.created_at.isoformat().endswith("+00:00")


async def test_save_updates_existing(repo: MySQLFaqRepository) -> None:
    faq = make_faq(published=False)
    await repo.save(faq)
    faq.question = "Updated?"
    faq.published = True
    faq.updated_at = datetime.now(UTC)
    updated = await repo.save(faq)
    assert updated.question == "Updated?"
    assert updated.published is True


async def test_search_published_only_excludes_unpublished(repo: MySQLFaqRepository) -> None:
    await repo.save(make_faq(published=True))
    await repo.save(make_faq(published=False))
    items, total = await repo.search(FaqSearchOption(published_only=True))
    assert total == 1
    assert items[0].published is True


async def test_search_published_only_false_includes_all(repo: MySQLFaqRepository) -> None:
    await repo.save(make_faq(published=True))
    await repo.save(make_faq(published=False))
    items, total = await repo.search(FaqSearchOption(published_only=False))
    assert total == 2


async def test_search_query_matches_question_and_answer(repo: MySQLFaqRepository) -> None:
    a = make_faq(published=True)
    a.question = "How do I reset my password?"
    b = make_faq(published=True)
    b.question = "Refund?"
    b.answer = "We refund within 7 days."
    c = make_faq(published=True)
    c.question = "Shipping policy?"
    c.answer = "We ship in 3 days."
    await repo.save(a)
    await repo.save(b)
    await repo.save(c)
    items, total = await repo.search(FaqSearchOption(query="refund"))
    assert total == 1


async def test_search_filters_by_category(
    repo: MySQLFaqRepository, category_repo: MySQLTicketCategoryRepository
) -> None:
    category = make_category()
    await category_repo.save(category)
    cat_id = await category_repo.find_id_by_uid(category.uid)
    await repo.save(make_faq(category_id=cat_id, published=True))
    await repo.save(make_faq(published=True))
    items, total = await repo.search(FaqSearchOption(category_uid=category.uid))
    assert total == 1


async def test_search_unknown_category_returns_empty(repo: MySQLFaqRepository) -> None:
    await repo.save(make_faq(published=True))
    items, total = await repo.search(FaqSearchOption(category_uid=uuid4()))
    assert total == 0
    assert items == []


async def test_resolves_category_uid(repo: MySQLFaqRepository, category_repo: MySQLTicketCategoryRepository) -> None:
    category = make_category()
    await category_repo.save(category)
    cat_id = await category_repo.find_id_by_uid(category.uid)
    faq = make_faq(category_id=cat_id, published=True)
    await repo.save(faq)

    found = await repo.find_by_uid(faq.uid)
    assert found is not None
    assert found.category_uid == category.uid

    items, _ = await repo.search(FaqSearchOption(category_uid=category.uid))
    assert items[0].category_uid == category.uid


async def test_update_applies_mutator_atomically(repo: MySQLFaqRepository) -> None:
    faq = make_faq(published=False)
    await repo.save(faq)

    def _mutate(f):
        f.question = "Locked question?"
        f.published = True
        f.updated_at = datetime.now(UTC)
        return f

    updated = await repo.update(faq.uid, _mutate)
    assert updated is not None
    assert updated.question == "Locked question?"
    assert updated.published is True
    reloaded = await repo.find_by_uid(faq.uid)
    assert reloaded is not None
    assert reloaded.published is True


async def test_update_returns_none_when_missing(repo: MySQLFaqRepository) -> None:
    assert await repo.update(uuid4(), lambda f: f) is None


async def test_search_unknown_sort_field_falls_back_to_default(repo: MySQLFaqRepository) -> None:
    await repo.save(make_faq(published=True, view_count=1))
    await repo.save(make_faq(published=True, view_count=9))
    items, total = await repo.search(FaqSearchOption(sort="metadata:desc"))
    assert total == 2
    assert items[0].view_count == 9


async def test_increment_view_count(repo: MySQLFaqRepository) -> None:
    faq = make_faq(view_count=5)
    await repo.save(faq)
    await repo.increment_view_count(faq.uid)
    found = await repo.find_by_uid(faq.uid)
    assert found is not None
    assert found.view_count == 6


async def test_delete_returns_true_when_exists(repo: MySQLFaqRepository) -> None:
    faq = make_faq()
    await repo.save(faq)
    assert await repo.delete_by_uid(faq.uid) is True
    assert await repo.find_by_uid(faq.uid) is None


async def test_delete_returns_false_when_missing(repo: MySQLFaqRepository) -> None:
    assert await repo.delete_by_uid(uuid4()) is False
