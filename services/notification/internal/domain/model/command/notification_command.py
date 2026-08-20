from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType


class SendNotificationCommand(BaseModel):
    user_uid: UUID
    notification_type: NotificationType
    channel: Channel
    title: str
    content: str
    order_uid: UUID | None = None


class CreateTemplateCommand(BaseModel):
    notification_type: NotificationType
    channel: Channel
    title_template: str
    content_template: str


class UpdateTemplateCommand(BaseModel):
    notification_type: NotificationType | None = None
    channel: Channel | None = None
    title_template: str | None = None
    content_template: str | None = None
