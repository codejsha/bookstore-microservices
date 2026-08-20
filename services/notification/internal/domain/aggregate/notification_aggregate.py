from datetime import datetime
from uuid import UUID

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus


class NotificationAggregate:
    def __init__(
        self,
        uid: UUID,
        user_uid: UUID,
        notification_type: NotificationType,
        channel: Channel,
        status: NotificationStatus,
        title: str,
        content: str,
        created_at: datetime,
        updated_at: datetime,
        sent_at: datetime | None = None,
        order_uid: UUID | None = None,
    ):
        self.uid = uid
        self.user_uid = user_uid
        self.notification_type = notification_type
        self.channel = channel
        self.status = status
        self.title = title
        self.content = content
        self.created_at = created_at
        self.updated_at = updated_at
        self.sent_at = sent_at
        self.order_uid = order_uid
