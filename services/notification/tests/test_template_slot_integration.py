from unittest.mock import MagicMock

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.error import TemplateAlreadyExistsError
from internal.domain.model.command.notification_command import CreateTemplateCommand, UpdateTemplateCommand
from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.mysql.template_repo import MySQLTemplateRepository


@pytest.fixture
def service(session_factory: async_sessionmaker[AsyncSession]) -> NotificationService:
    return NotificationService(MagicMock(), MySQLTemplateRepository(session_factory))


def _create_command(
    notification_type: NotificationType = NotificationType.ORDER_PLACED,
    channel: Channel = Channel.EMAIL,
) -> CreateTemplateCommand:
    return CreateTemplateCommand(
        notification_type=notification_type,
        channel=channel,
        title_template="Hi {name}",
        content_template="Order {order_id}",
    )


async def test_create_template_occupied_slot_raises_template_already_exists_error(service: NotificationService) -> None:
    await service.create_template(_create_command())
    with pytest.raises(TemplateAlreadyExistsError):
        await service.create_template(_create_command())


async def test_create_template_soft_deleted_succeeds(service: NotificationService) -> None:
    first = await service.create_template(_create_command())
    assert await service.delete_template(first.uid) is True
    second = await service.create_template(_create_command())
    assert second.uid != first.uid


async def test_update_template_moving_onto_soft_deleted_slot_succeeds(service: NotificationService) -> None:
    holder = await service.create_template(_create_command(channel=Channel.EMAIL))
    mover = await service.create_template(_create_command(channel=Channel.SMS))
    assert await service.delete_template(holder.uid) is True
    moved = await service.update_template(mover.uid, UpdateTemplateCommand(channel=Channel.EMAIL))
    assert moved is not None
    assert moved.channel == Channel.EMAIL


async def test_update_template_moving_onto_active_slot_raises_template_already_exists_error(
    service: NotificationService,
) -> None:
    await service.create_template(_create_command(channel=Channel.EMAIL))
    mover = await service.create_template(_create_command(channel=Channel.SMS))
    with pytest.raises(TemplateAlreadyExistsError):
        await service.update_template(mover.uid, UpdateTemplateCommand(channel=Channel.EMAIL))
