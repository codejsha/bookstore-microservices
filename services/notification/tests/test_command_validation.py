from uuid import uuid4

import pytest
from pydantic import ValidationError

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.model.command.notification_command import (
    CreateTemplateCommand,
    SendNotificationCommand,
    UpdateTemplateCommand,
)


class TestSendNotificationCommand:
    def test_send_notification_command_every_field_valid_passes(self):
        SendNotificationCommand(
            user_uid=uuid4(),
            notification_type=NotificationType.SHIPMENT_DISPATCHED,
            channel=Channel.EMAIL,
            title="Your order has shipped",
            content="Order 123 is on the way.",
        )

    @pytest.mark.parametrize("field", ["title", "content"])
    def test_send_notification_command_blank_field_raises_validation_error(self, field):
        kwargs = {
            "user_uid": uuid4(),
            "notification_type": NotificationType.SHIPMENT_DISPATCHED,
            "channel": Channel.EMAIL,
            "title": "t",
            "content": "c",
        }
        kwargs[field] = "  "
        with pytest.raises(ValidationError):
            SendNotificationCommand(**kwargs)

    @pytest.mark.parametrize(("field", "limit"), [("title", 255), ("content", 10000)])
    def test_send_notification_command_over_length_field_raises_validation_error(self, field, limit):
        kwargs = {
            "user_uid": uuid4(),
            "notification_type": NotificationType.SHIPMENT_DISPATCHED,
            "channel": Channel.EMAIL,
            "title": "t",
            "content": "c",
        }
        kwargs[field] = "x" * (limit + 1)
        with pytest.raises(ValidationError):
            SendNotificationCommand(**kwargs)
        kwargs[field] = "x" * limit
        SendNotificationCommand(**kwargs)


class TestTemplateCommands:
    def test_create_template_command_every_field_valid_passes(self):
        CreateTemplateCommand(
            notification_type=NotificationType.SHIPMENT_DISPATCHED,
            channel=Channel.EMAIL,
            title_template="Order {order_uid}",
            content_template="Status: {status}",
        )

    def test_create_template_command_blank_template_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateTemplateCommand(
                notification_type=NotificationType.SHIPMENT_DISPATCHED,
                channel=Channel.EMAIL,
                title_template=" ",
                content_template="Status: {status}",
            )

    def test_update_template_command_every_field_none_passes(self):
        UpdateTemplateCommand()

    def test_update_template_command_blank_content_raises_validation_error(self):
        with pytest.raises(ValidationError):
            UpdateTemplateCommand(content_template="")

    @pytest.mark.parametrize(("field", "limit"), [("title_template", 255), ("content_template", 10000)])
    def test_create_template_command_over_length_field_raises_validation_error(self, field, limit):
        kwargs = {
            "notification_type": NotificationType.SHIPMENT_DISPATCHED,
            "channel": Channel.EMAIL,
            "title_template": "t",
            "content_template": "c",
        }
        kwargs[field] = "x" * (limit + 1)
        with pytest.raises(ValidationError):
            CreateTemplateCommand(**kwargs)
        kwargs[field] = "x" * limit
        CreateTemplateCommand(**kwargs)

    @pytest.mark.parametrize(("field", "limit"), [("title_template", 255), ("content_template", 10000)])
    def test_update_template_command_over_length_field_raises_validation_error(self, field, limit):
        with pytest.raises(ValidationError):
            UpdateTemplateCommand(**{field: "x" * (limit + 1)})
        UpdateTemplateCommand(**{field: "x" * limit})
