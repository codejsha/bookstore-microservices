from abc import ABC, abstractmethod
from collections.abc import Callable
from uuid import UUID

from internal.domain.aggregate.notification_aggregate import NotificationAggregate
from internal.domain.aggregate.template_aggregate import TemplateAggregate
from internal.domain.constant.notification_type import NotificationType
from internal.domain.model.option.notification_option import (
    NotificationFilterOption,
    TemplateFilterOption,
)


class NotificationRepository(ABC):
    @abstractmethod
    async def save(self, notification: NotificationAggregate) -> NotificationAggregate: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> NotificationAggregate | None: ...

    @abstractmethod
    async def find_by_order_uid_and_type(
        self, order_uid: UUID, notification_type: NotificationType
    ) -> NotificationAggregate | None: ...

    @abstractmethod
    async def find_all(self, option: NotificationFilterOption) -> tuple[list[NotificationAggregate], int]: ...

    @abstractmethod
    async def update(
        self,
        uid: UUID,
        mutator: Callable[[NotificationAggregate], NotificationAggregate],
    ) -> NotificationAggregate | None: ...


class TemplateRepository(ABC):
    @abstractmethod
    async def save(self, template: TemplateAggregate) -> TemplateAggregate: ...

    @abstractmethod
    async def find_by_uid(self, uid: UUID) -> TemplateAggregate | None: ...

    @abstractmethod
    async def find_all(self, option: TemplateFilterOption) -> tuple[list[TemplateAggregate], int]: ...

    @abstractmethod
    async def delete_by_uid(self, uid: UUID) -> bool: ...
