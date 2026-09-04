from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus
from internal.domain.error import TemplateAlreadyExistsError
from internal.domain.model.command.notification_command import (
    CreateTemplateCommand,
    SendNotificationCommand,
    UpdateTemplateCommand,
)
from internal.domain.model.option.notification_option import (
    NotificationFilterOption,
    TemplateFilterOption,
)
from internal.domain.service.notification_service import NotificationService
from tests.conftest import make_notification, make_template


@pytest.fixture
def service(notification_repo: MagicMock, template_repo: MagicMock) -> NotificationService:
    return NotificationService(notification_repo=notification_repo, template_repo=template_repo)


class TestSendNotification:
    async def test_send_notification_valid_command_persists_pending_notification(
        self, service: NotificationService, notification_repo: MagicMock
    ) -> None:
        command = SendNotificationCommand(
            user_uid=uuid4(),
            notification_type=NotificationType.ORDER_PLACED,
            channel=Channel.EMAIL,
            title="Welcome",
            content="body",
        )
        result = await service.send_notification(command)
        assert result.status == NotificationStatus.PENDING
        assert result.title == "Welcome"
        assert result.sent_at is None
        notification_repo.save.assert_called_once()


class TestGetNotification:
    async def test_get_notification_roundtrips(
        self, service: NotificationService, notification_repo: MagicMock
    ) -> None:
        notification = make_notification()
        notification_repo.find_by_uid.return_value = notification
        assert await service.get_notification(notification.uid) is notification

    async def test_get_notification_missing_is_none(
        self, service: NotificationService, notification_repo: MagicMock
    ) -> None:
        notification_repo.find_by_uid.return_value = None
        assert await service.get_notification(uuid4()) is None


class TestListNotifications:
    async def test_list_notifications_with_option_delegates_to_repository(
        self, service: NotificationService, notification_repo: MagicMock
    ) -> None:
        option = NotificationFilterOption()
        notification_repo.find_all.return_value = ([make_notification()], 1)
        result, total = await service.list_notifications(option)
        assert total == 1
        notification_repo.find_all.assert_called_once_with(option)


class TestMarkAsSent:
    async def test_mark_as_sent_missing_is_none(
        self, service: NotificationService, notification_repo: MagicMock
    ) -> None:
        notification_repo.find_by_uid.return_value = None
        assert await service.mark_as_sent(uuid4()) is None
        notification_repo.save.assert_not_called()

    async def test_mark_as_sent_found_sets_status_and_sent_at(
        self, service: NotificationService, notification_repo: MagicMock
    ) -> None:
        notification = make_notification()
        notification_repo.find_by_uid.return_value = notification
        result = await service.mark_as_sent(notification.uid)
        assert result is not None
        assert result.status == NotificationStatus.SENT
        assert result.sent_at is not None


class TestMarkAsFailed:
    async def test_mark_as_failed_missing_is_none(
        self, service: NotificationService, notification_repo: MagicMock
    ) -> None:
        notification_repo.find_by_uid.return_value = None
        assert await service.mark_as_failed(uuid4()) is None

    async def test_mark_as_failed_found_sets_failed_status_without_sent_at(
        self, service: NotificationService, notification_repo: MagicMock
    ) -> None:
        notification = make_notification()
        notification_repo.find_by_uid.return_value = notification
        result = await service.mark_as_failed(notification.uid)
        assert result is not None
        assert result.status == NotificationStatus.FAILED
        assert result.sent_at is None


class TestCreateTemplate:
    async def test_create_template_slot_free_persists_template(
        self, service: NotificationService, template_repo: MagicMock
    ) -> None:
        command = CreateTemplateCommand(
            notification_type=NotificationType.ORDER_PLACED,
            channel=Channel.EMAIL,
            title_template="Hi {name}",
            content_template="Order {order_id}",
        )
        result = await service.create_template(command)
        assert result.notification_type == NotificationType.ORDER_PLACED
        assert result.channel == Channel.EMAIL
        assert result.title_template == "Hi {name}"
        template_repo.save.assert_called_once()

    async def test_create_template_type_and_occupied_channel_raises_template_already_exists_error(
        self, service: NotificationService, template_repo: MagicMock
    ) -> None:
        template_repo.find_by_type_and_channel.return_value = make_template(
            notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL
        )
        command = CreateTemplateCommand(
            notification_type=NotificationType.ORDER_PLACED,
            channel=Channel.EMAIL,
            title_template="Hi {name}",
            content_template="Order {order_id}",
        )
        with pytest.raises(TemplateAlreadyExistsError):
            await service.create_template(command)
        template_repo.save.assert_not_called()

    async def test_create_template_soft_deleted_slot_purges_it_and_saves(
        self, service: NotificationService, template_repo: MagicMock
    ) -> None:
        template_repo.find_by_type_and_channel.return_value = None
        template_repo.purge_deleted_by_type_and_channel.return_value = True
        command = CreateTemplateCommand(
            notification_type=NotificationType.ORDER_PLACED,
            channel=Channel.EMAIL,
            title_template="Hi {name}",
            content_template="Order {order_id}",
        )
        await service.create_template(command)
        template_repo.purge_deleted_by_type_and_channel.assert_called_once_with(
            NotificationType.ORDER_PLACED, Channel.EMAIL
        )
        template_repo.save.assert_called_once()


class TestGetTemplate:
    async def test_get_template_roundtrips(self, service: NotificationService, template_repo: MagicMock) -> None:
        template = make_template()
        template_repo.find_by_uid.return_value = template
        assert await service.get_template(template.uid) is template


class TestListTemplates:
    async def test_list_templates_with_option_delegates_to_repository(
        self, service: NotificationService, template_repo: MagicMock
    ) -> None:
        templates = [make_template()]
        template_repo.find_all.return_value = (templates, 1)
        option = TemplateFilterOption()
        assert await service.list_templates(option) == (templates, 1)
        template_repo.find_all.assert_called_once_with(option)


class TestUpdateTemplate:
    async def test_update_template_missing_is_none(
        self, service: NotificationService, template_repo: MagicMock
    ) -> None:
        template_repo.find_by_uid.return_value = None
        result = await service.update_template(uuid4(), UpdateTemplateCommand(title_template="X"))
        assert result is None
        template_repo.save.assert_not_called()

    async def test_update_template_partial_fields_updates_only_those(
        self, service: NotificationService, template_repo: MagicMock
    ) -> None:
        template = make_template()
        original_content = template.content_template
        original_type = template.notification_type
        template_repo.find_by_uid.return_value = template
        command = UpdateTemplateCommand(channel=Channel.SMS, title_template="New title")
        result = await service.update_template(template.uid, command)
        assert result is not None
        assert result.channel == Channel.SMS
        assert result.title_template == "New title"
        assert result.content_template == original_content
        assert result.notification_type == original_type

    async def test_update_template_moving_onto_occupied_slot_raises_template_already_exists_error(
        self, service: NotificationService, template_repo: MagicMock
    ) -> None:
        template = make_template(notification_type=NotificationType.ORDER_PLACED, channel=Channel.EMAIL)
        template_repo.find_by_uid.return_value = template
        template_repo.find_by_type_and_channel.return_value = make_template(
            notification_type=NotificationType.ORDER_PLACED, channel=Channel.SMS
        )
        with pytest.raises(TemplateAlreadyExistsError):
            await service.update_template(template.uid, UpdateTemplateCommand(channel=Channel.SMS))
        template_repo.save.assert_not_called()

    async def test_update_template_type_and_channel_unchanged_skips_slot_check_and_saves(
        self, service: NotificationService, template_repo: MagicMock
    ) -> None:
        template = make_template()
        template_repo.find_by_uid.return_value = template
        await service.update_template(template.uid, UpdateTemplateCommand(title_template="New title"))
        template_repo.find_by_type_and_channel.assert_not_called()
        template_repo.save.assert_called_once()


class TestDeleteTemplate:
    async def test_delete_template_delegates_to_repository(
        self, service: NotificationService, template_repo: MagicMock
    ) -> None:
        template_repo.delete_by_uid.return_value = True
        uid = uuid4()
        assert await service.delete_template(uid) is True
        template_repo.delete_by_uid.assert_called_once_with(uid)
