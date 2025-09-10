from abc import ABC, abstractmethod
from uuid import UUID

from internal.domain.aggregate.notification_aggregate import NotificationAggregate
from internal.domain.aggregate.template_aggregate import TemplateAggregate
from internal.domain.model.command.notification_command import (
    CreateTemplateCommand,
    SendNotificationCommand,
    UpdateTemplateCommand,
)
from internal.domain.model.option.notification_option import (
    NotificationFilterOption,
    TemplateFilterOption,
)


class NotificationUseCase(ABC):
    @abstractmethod
    def send_notification(self, command: SendNotificationCommand) -> NotificationAggregate: ...

    @abstractmethod
    def get_notification(self, uid: UUID) -> NotificationAggregate | None: ...

    @abstractmethod
    def list_notifications(self, option: NotificationFilterOption) -> tuple[list[NotificationAggregate], int]: ...

    @abstractmethod
    def mark_as_sent(self, uid: UUID) -> NotificationAggregate | None: ...

    @abstractmethod
    def mark_as_failed(self, uid: UUID) -> NotificationAggregate | None: ...


class TemplateUseCase(ABC):
    @abstractmethod
    def create_template(self, command: CreateTemplateCommand) -> TemplateAggregate: ...

    @abstractmethod
    def get_template(self, uid: UUID) -> TemplateAggregate | None: ...

    @abstractmethod
    def list_templates(self, option: TemplateFilterOption) -> tuple[list[TemplateAggregate], int]: ...

    @abstractmethod
    def update_template(self, uid: UUID, command: UpdateTemplateCommand) -> TemplateAggregate | None: ...

    @abstractmethod
    def delete_template(self, uid: UUID) -> bool: ...
