from datetime import UTC, datetime
from uuid import uuid4

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.model.option.notification_option import TemplateFilterOption
from internal.infrastructure.adapter.mysql.template_repo import MySQLTemplateRepository
from tests.conftest import make_template


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLTemplateRepository:
    return MySQLTemplateRepository(session_factory)


async def test_save_inserts_template(repo: MySQLTemplateRepository) -> None:
    template = make_template()
    saved = await repo.save(template)
    assert saved.uid == template.uid


async def test_find_by_uid_returns_saved(repo: MySQLTemplateRepository) -> None:
    template = make_template()
    await repo.save(template)
    found = await repo.find_by_uid(template.uid)
    assert found is not None
    assert found.title_template == template.title_template


async def test_find_by_uid_returns_none(repo: MySQLTemplateRepository) -> None:
    assert await repo.find_by_uid(uuid4()) is None


async def test_save_updates_existing_template(repo: MySQLTemplateRepository) -> None:
    template = make_template(channel=Channel.EMAIL)
    await repo.save(template)
    template.channel = Channel.SMS
    template.title_template = "Updated"
    template.updated_at = datetime.now(UTC)
    updated = await repo.save(template)
    assert updated.channel == Channel.SMS
    assert updated.title_template == "Updated"


async def test_find_all_excludes_deleted(repo: MySQLTemplateRepository) -> None:
    a = make_template(channel=Channel.EMAIL)
    b = make_template(channel=Channel.SMS)
    await repo.save(a)
    await repo.save(b)
    await repo.delete_by_uid(a.uid)
    templates, total = await repo.find_all(TemplateFilterOption())
    uids = {t.uid for t in templates}
    assert a.uid not in uids
    assert b.uid in uids
    assert total == 1


async def test_find_all_filters_by_channel(repo: MySQLTemplateRepository) -> None:
    await repo.save(make_template(channel=Channel.EMAIL))
    await repo.save(make_template(channel=Channel.SMS))
    templates, total = await repo.find_all(TemplateFilterOption(channel=Channel.SMS))
    assert total == 1
    assert all(t.channel == Channel.SMS for t in templates)


async def test_find_all_filters_by_notification_type(repo: MySQLTemplateRepository) -> None:
    await repo.save(make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL))
    await repo.save(make_template(notification_type=NotificationType.PAYMENT_SUCCESS, channel=Channel.SMS))
    templates, total = await repo.find_all(TemplateFilterOption(notification_type=NotificationType.PAYMENT_SUCCESS))
    assert total == 1
    assert all(t.notification_type == NotificationType.PAYMENT_SUCCESS for t in templates)


async def test_find_all_paginates(repo: MySQLTemplateRepository) -> None:
    await repo.save(make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL))
    await repo.save(make_template(notification_type=NotificationType.ORDER_CONFIRMED, channel=Channel.SMS))
    await repo.save(make_template(notification_type=NotificationType.PAYMENT_SUCCESS, channel=Channel.PUSH))
    page0, total = await repo.find_all(TemplateFilterOption(page=0, size=2))
    assert total == 3
    assert len(page0) == 2
    page1, _ = await repo.find_all(TemplateFilterOption(page=1, size=2))
    assert len(page1) == 1


async def test_find_all_unknown_sort_field_falls_back_to_default(repo: MySQLTemplateRepository) -> None:
    await repo.save(make_template(channel=Channel.EMAIL))
    templates, total = await repo.find_all(TemplateFilterOption(sort="secret_column:desc"))
    assert total == 1
    assert len(templates) == 1


async def test_delete_by_uid_returns_true_when_exists(repo: MySQLTemplateRepository) -> None:
    template = make_template()
    await repo.save(template)
    assert await repo.delete_by_uid(template.uid) is True
    assert await repo.find_by_uid(template.uid) is None


async def test_delete_by_uid_returns_false_when_missing(repo: MySQLTemplateRepository) -> None:
    assert await repo.delete_by_uid(uuid4()) is False
