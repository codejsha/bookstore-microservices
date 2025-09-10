from datetime import datetime
from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus


class SendNotificationRequest(BaseModel):
    user_uid: UUID
    notification_type: NotificationType
    channel: Channel
    title: str
    content: str


class NotificationResponse(BaseModel):
    uid: UUID
    user_uid: UUID
    notification_type: NotificationType
    channel: Channel
    status: NotificationStatus
    title: str
    content: str
    sent_at: datetime | None = None
    created_at: datetime
    updated_at: datetime


class NotificationFindAllResponse(BaseModel):
    total: int
    items: list[NotificationResponse]


class CreateTemplateRequest(BaseModel):
    notification_type: NotificationType
    channel: Channel
    title_template: str
    content_template: str


class UpdateTemplateRequest(BaseModel):
    notification_type: NotificationType | None = None
    channel: Channel | None = None
    title_template: str | None = None
    content_template: str | None = None


class TemplateResponse(BaseModel):
    uid: UUID
    notification_type: NotificationType
    channel: Channel
    title_template: str
    content_template: str
    created_at: datetime
    updated_at: datetime


class TemplateFindAllResponse(BaseModel):
    total: int
    items: list[TemplateResponse]


class ErrorResponse(BaseModel):
    detail: str
