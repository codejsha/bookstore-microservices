from datetime import datetime
from uuid import UUID

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType


class TemplateAggregate:
    def __init__(
        self,
        uid: UUID,
        notification_type: NotificationType,
        channel: Channel,
        title_template: str,
        content_template: str,
        created_at: datetime,
        updated_at: datetime,
    ):
        self.uid = uid
        self.notification_type = notification_type
        self.channel = channel
        self.title_template = title_template
        self.content_template = content_template
        self.created_at = created_at
        self.updated_at = updated_at
