from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType


class TemplateAlreadyExistsError(Exception):
    def __init__(self, notification_type: NotificationType, channel: Channel):
        self.notification_type = notification_type
        self.channel = channel
        super().__init__(f"template already exists for {notification_type.value} via {channel.value}")
