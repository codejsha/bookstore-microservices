from uuid import UUID

from pydantic import BaseModel

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus


class NotificationFilterOption(BaseModel):
    user_uid: UUID | None = None
    notification_type: NotificationType | None = None
    channel: Channel | None = None
    status: NotificationStatus | None = None
    page: int = 0
    size: int = 20
    sort: str = "created_at:desc"


class TemplateFilterOption(BaseModel):
    notification_type: NotificationType | None = None
    channel: Channel | None = None
    page: int = 0
    size: int = 20
    sort: str = "created_at:desc"
