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


async def test_save_new_inserts_row(repo: MySQLTemplateRepository) -> None:
    template = make_template()
    saved = await repo.save(template)
    assert saved.uid == template.uid


async def test_find_by_uid_roundtrips(repo: MySQLTemplateRepository) -> None:
    template = make_template()
    await repo.save(template)
    found = await repo.find_by_uid(template.uid)
    assert found is not None
    assert found.title_template == template.title_template


async def test_find_by_uid_missing_is_none(repo: MySQLTemplateRepository) -> None:
    assert await repo.find_by_uid(uuid4()) is None


async def test_save_found_updates_row(repo: MySQLTemplateRepository) -> None:
    template = make_template(channel=Channel.EMAIL)
    await repo.save(template)
    template.channel = Channel.SMS
    template.title_template = "Updated"
    template.updated_at = datetime.now(UTC)
    updated = await repo.save(template)
    assert updated.channel == Channel.SMS
    assert updated.title_template == "Updated"


async def test_find_all_soft_deleted_excludes_it(repo: MySQLTemplateRepository) -> None:
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


async def test_find_all_by_channel_matches_templates(repo: MySQLTemplateRepository) -> None:
    await repo.save(make_template(channel=Channel.EMAIL))
    await repo.save(make_template(channel=Channel.SMS))
    templates, total = await repo.find_all(TemplateFilterOption(channel=Channel.SMS))
    assert total == 1
    assert all(t.channel == Channel.SMS for t in templates)


async def test_find_all_by_type_matches_templates(repo: MySQLTemplateRepository) -> None:
    await repo.save(make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL))
    await repo.save(make_template(notification_type=NotificationType.PAYMENT_SUCCESS, channel=Channel.SMS))
    templates, total = await repo.find_all(TemplateFilterOption(notification_type=NotificationType.PAYMENT_SUCCESS))
    assert total == 1
    assert all(t.notification_type == NotificationType.PAYMENT_SUCCESS for t in templates)


async def test_find_all_with_page_size_pages_with_total(repo: MySQLTemplateRepository) -> None:
    await repo.save(make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL))
    await repo.save(make_template(notification_type=NotificationType.ORDER_CONFIRMED, channel=Channel.SMS))
    await repo.save(make_template(notification_type=NotificationType.PAYMENT_SUCCESS, channel=Channel.PUSH))
    page0, total = await repo.find_all(TemplateFilterOption(page=0, size=2))
    assert total == 3
    assert len(page0) == 2
    page1, _ = await repo.find_all(TemplateFilterOption(page=1, size=2))
    assert len(page1) == 1


async def test_find_all_unknown_sort_field_falls_back_to_default_order(repo: MySQLTemplateRepository) -> None:
    await repo.save(make_template(channel=Channel.EMAIL))
    templates, total = await repo.find_all(TemplateFilterOption(sort="secret_column:desc"))
    assert total == 1
    assert len(templates) == 1


async def test_delete_by_uid_found_reports_true_and_hides_it(repo: MySQLTemplateRepository) -> None:
    template = make_template()
    await repo.save(template)
    assert await repo.delete_by_uid(template.uid) is True
    assert await repo.find_by_uid(template.uid) is None


async def test_delete_by_uid_missing_reports_false(repo: MySQLTemplateRepository) -> None:
    assert await repo.delete_by_uid(uuid4()) is False


async def test_find_by_type_and_channel_active_is_found(repo: MySQLTemplateRepository) -> None:
    template = make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL)
    await repo.save(template)
    found = await repo.find_by_type_and_channel(NotificationType.ORDER_PLACED, Channel.EMAIL)
    assert found is not None
    assert found.uid == template.uid


async def test_find_by_type_and_channel_soft_deleted_is_none(repo: MySQLTemplateRepository) -> None:
    template = make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL)
    await repo.save(template)
    await repo.delete_by_uid(template.uid)
    assert await repo.find_by_type_and_channel(NotificationType.ORDER_PLACED, Channel.EMAIL) is None


async def test_purge_deleted_by_type_and_channel_soft_deleted_slot_frees_it(repo: MySQLTemplateRepository) -> None:
    template = make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL)
    await repo.save(template)
    await repo.delete_by_uid(template.uid)
    assert await repo.purge_deleted_by_type_and_channel(NotificationType.ORDER_PLACED, Channel.EMAIL) is True
    replacement = make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL)
    saved = await repo.save(replacement)
    assert saved.uid == replacement.uid


async def test_purge_deleted_by_type_and_channel_active_slot_keeps_row_and_reports_false(
    repo: MySQLTemplateRepository,
) -> None:
    template = make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL)
    await repo.save(template)
    assert await repo.purge_deleted_by_type_and_channel(NotificationType.ORDER_PLACED, Channel.EMAIL) is False
    assert await repo.find_by_uid(template.uid) is not None
