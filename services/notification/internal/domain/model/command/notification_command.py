from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.model.command import NonBlankStr255, NonBlankStr10000


class SendNotificationCommand(BaseModel):
    user_uid: UUID
    notification_type: NotificationType
    channel: Channel
    title: NonBlankStr255
    content: NonBlankStr10000
    order_uid: UUID | None = None


class CreateTemplateCommand(BaseModel):
    notification_type: NotificationType
    channel: Channel
    title_template: NonBlankStr255
    content_template: NonBlankStr10000


class UpdateTemplateCommand(BaseModel):
    notification_type: NotificationType | None = None
    channel: Channel | None = None
    title_template: NonBlankStr255 | None = None
    content_template: NonBlankStr10000 | None = None
