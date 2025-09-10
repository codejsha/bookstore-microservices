from datetime import UTC, datetime, timedelta
from uuid import uuid4

import pytest
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus
from internal.domain.model.option.notification_option import NotificationFilterOption
from internal.infrastructure.adapter.mysql.notification_repo import (
    MySQLNotificationRepository,
)
from tests.conftest import make_notification


@pytest.fixture
def repo(session_factory: async_sessionmaker[AsyncSession]) -> MySQLNotificationRepository:
    return MySQLNotificationRepository(session_factory)


async def test_save_inserts_new(repo: MySQLNotificationRepository) -> None:
    notification = make_notification()
    saved = await repo.save(notification)
    assert saved.uid == notification.uid
    assert saved.status == NotificationStatus.PENDING


async def test_find_by_uid_returns_saved(repo: MySQLNotificationRepository) -> None:
    notification = make_notification()
    await repo.save(notification)
    found = await repo.find_by_uid(notification.uid)
    assert found is not None
    assert found.title == notification.title


async def test_find_by_uid_returns_none(repo: MySQLNotificationRepository) -> None:
    assert await repo.find_by_uid(uuid4()) is None


async def test_fetched_datetime_is_tz_aware_utc(repo: MySQLNotificationRepository) -> None:
    notification = make_notification()
    await repo.save(notification)

    found = await repo.find_by_uid(notification.uid)
    assert found is not None
    assert found.created_at.tzinfo is not None
    assert found.created_at.utcoffset() == timedelta(0)
    assert found.created_at.isoformat().endswith("+00:00")


async def test_save_updates_status_and_sent_at(repo: MySQLNotificationRepository) -> None:
    notification = make_notification()
    await repo.save(notification)
    notification.status = NotificationStatus.SENT
    notification.sent_at = datetime.now(UTC)
    notification.updated_at = datetime.now(UTC)
    updated = await repo.save(notification)
    assert updated.status == NotificationStatus.SENT
    assert updated.sent_at is not None


async def test_find_by_order_uid_and_type_returns_saved(repo: MySQLNotificationRepository) -> None:
    order_uid = uuid4()
    notification = make_notification(
        order_uid=order_uid,
        notification_type=NotificationType.SHIPMENT_DISPATCHED,
    )
    await repo.save(notification)
    found = await repo.find_by_order_uid_and_type(order_uid, NotificationType.SHIPMENT_DISPATCHED)
    assert found is not None
    assert found.uid == notification.uid
    assert found.order_uid == order_uid


async def test_find_by_order_uid_and_type_returns_none_for_other_type(repo: MySQLNotificationRepository) -> None:
    order_uid = uuid4()
    await repo.save(make_notification(order_uid=order_uid, notification_type=NotificationType.SHIPMENT_DISPATCHED))
    assert await repo.find_by_order_uid_and_type(order_uid, NotificationType.SHIPMENT_DELIVERED) is None


async def test_order_type_unique_constraint_rejects_duplicate(repo: MySQLNotificationRepository) -> None:
    order_uid = uuid4()
    await repo.save(make_notification(order_uid=order_uid, notification_type=NotificationType.SHIPMENT_DISPATCHED))
    with pytest.raises(IntegrityError):
        await repo.save(make_notification(order_uid=order_uid, notification_type=NotificationType.SHIPMENT_DISPATCHED))


async def test_null_order_uid_notifications_do_not_collide(repo: MySQLNotificationRepository) -> None:
    await repo.save(make_notification(notification_type=NotificationType.ORDER_PLACED))
    await repo.save(make_notification(notification_type=NotificationType.ORDER_PLACED))
    items, total = await repo.find_all(NotificationFilterOption(notification_type=NotificationType.ORDER_PLACED))
    assert total == 2


async def test_find_all_filters_by_user_uid(repo: MySQLNotificationRepository) -> None:
    target_user = uuid4()
    await repo.save(make_notification(user_uid=target_user))
    await repo.save(make_notification())
    items, total = await repo.find_all(NotificationFilterOption(user_uid=target_user))
    assert total == 1


async def test_find_all_filters_by_channel(repo: MySQLNotificationRepository) -> None:
    await repo.save(make_notification(channel=Channel.EMAIL))
    await repo.save(make_notification(channel=Channel.SMS))
    items, total = await repo.find_all(NotificationFilterOption(channel=Channel.SMS))
    assert total == 1
    assert items[0].channel == Channel.SMS


async def test_find_all_filters_by_status_and_type(repo: MySQLNotificationRepository) -> None:
    await repo.save(
        make_notification(
            notification_type=NotificationType.ORDER_PLACED,
            status=NotificationStatus.PENDING,
        )
    )
    await repo.save(
        make_notification(
            notification_type=NotificationType.ORDER_PLACED,
            status=NotificationStatus.SENT,
        )
    )
    await repo.save(
        make_notification(
            notification_type=NotificationType.PAYMENT_SUCCESS,
            status=NotificationStatus.SENT,
        )
    )
    items, total = await repo.find_all(
        NotificationFilterOption(
            notification_type=NotificationType.ORDER_PLACED,
            status=NotificationStatus.SENT,
        )
    )
    assert total == 1


async def test_find_all_paginates_in_descending_order(repo: MySQLNotificationRepository) -> None:
    base = datetime.now(UTC)
    notifications = []
    for i in range(3):
        n = make_notification()
        n.created_at = base + timedelta(seconds=i)
        n.updated_at = n.created_at
        await repo.save(n)
        notifications.append(n)
    items, total = await repo.find_all(NotificationFilterOption(page=0, size=2))
    assert total == 3
    assert len(items) == 2
    assert items[0].uid == notifications[2].uid


async def test_find_all_unknown_sort_field_falls_back_to_default(repo: MySQLNotificationRepository) -> None:
    base = datetime.now(UTC)
    notifications = []
    for i in range(2):
        n = make_notification()
        n.created_at = base + timedelta(seconds=i)
        n.updated_at = n.created_at
        await repo.save(n)
        notifications.append(n)
    items, total = await repo.find_all(NotificationFilterOption(sort="metadata:desc"))
    assert total == 2
    assert items[0].uid == notifications[1].uid


async def test_update_applies_mutator_atomically(repo: MySQLNotificationRepository) -> None:
    notification = make_notification()
    await repo.save(notification)
    sent_at = datetime.now(UTC)

    def _mutate(n):
        n.status = NotificationStatus.SENT
        n.sent_at = sent_at
        n.updated_at = sent_at
        return n

    updated = await repo.update(notification.uid, _mutate)
    assert updated is not None
    assert updated.status == NotificationStatus.SENT
    assert updated.sent_at is not None
    reloaded = await repo.find_by_uid(notification.uid)
    assert reloaded is not None
    assert reloaded.status == NotificationStatus.SENT


async def test_update_returns_none_when_missing(repo: MySQLNotificationRepository) -> None:
    assert await repo.update(uuid4(), lambda n: n) is None
