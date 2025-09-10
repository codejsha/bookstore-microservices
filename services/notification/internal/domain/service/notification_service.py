from datetime import UTC, datetime
from uuid import UUID, uuid7

from internal.application.port.repo.repos import NotificationRepository, TemplateRepository
from internal.domain.aggregate.notification_aggregate import NotificationAggregate
from internal.domain.aggregate.template_aggregate import TemplateAggregate
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus
from internal.domain.model.command.notification_command import (
    CreateTemplateCommand,
    SendNotificationCommand,
    UpdateTemplateCommand,
)
from internal.domain.model.option.notification_option import (
    NotificationFilterOption,
    TemplateFilterOption,
)


class NotificationService:
    def __init__(
        self,
        notification_repo: NotificationRepository,
        template_repo: TemplateRepository,
    ):
        self._notification_repo = notification_repo
        self._template_repo = template_repo

    async def send_notification(self, command: SendNotificationCommand) -> NotificationAggregate:
        now = datetime.now(UTC)
        notification = NotificationAggregate(
            uid=uuid7(),
            user_uid=command.user_uid,
            notification_type=command.notification_type,
            channel=command.channel,
            status=NotificationStatus.PENDING,
            title=command.title,
            content=command.content,
            created_at=now,
            updated_at=now,
            order_uid=command.order_uid,
        )
        return await self._notification_repo.save(notification)

    async def get_notification(self, uid: UUID) -> NotificationAggregate | None:
        return await self._notification_repo.find_by_uid(uid)

    async def find_by_order_and_type(
        self, order_uid: UUID, notification_type: NotificationType
    ) -> NotificationAggregate | None:
        return await self._notification_repo.find_by_order_uid_and_type(order_uid, notification_type)

    async def list_notifications(self, option: NotificationFilterOption) -> tuple[list[NotificationAggregate], int]:
        return await self._notification_repo.find_all(option)

    async def mark_as_sent(self, uid: UUID) -> NotificationAggregate | None:
        now = datetime.now(UTC)

        def _mutate(notification: NotificationAggregate) -> NotificationAggregate:
            notification.status = NotificationStatus.SENT
            notification.sent_at = now
            notification.updated_at = now
            return notification

        return await self._notification_repo.update(uid, _mutate)

    async def mark_as_failed(self, uid: UUID) -> NotificationAggregate | None:
        now = datetime.now(UTC)

        def _mutate(notification: NotificationAggregate) -> NotificationAggregate:
            notification.status = NotificationStatus.FAILED
            notification.updated_at = now
            return notification

        return await self._notification_repo.update(uid, _mutate)

    async def create_template(self, command: CreateTemplateCommand) -> TemplateAggregate:
        now = datetime.now(UTC)
        template = TemplateAggregate(
            uid=uuid7(),
            notification_type=command.notification_type,
            channel=command.channel,
            title_template=command.title_template,
            content_template=command.content_template,
            created_at=now,
            updated_at=now,
        )
        return await self._template_repo.save(template)

    async def get_template(self, uid: UUID) -> TemplateAggregate | None:
        return await self._template_repo.find_by_uid(uid)

    async def list_templates(self, option: TemplateFilterOption) -> tuple[list[TemplateAggregate], int]:
        return await self._template_repo.find_all(option)

    async def update_template(self, uid: UUID, command: UpdateTemplateCommand) -> TemplateAggregate | None:
        template = await self._template_repo.find_by_uid(uid)
        if template is None:
            return None
        now = datetime.now(UTC)
        if command.notification_type is not None:
            template.notification_type = command.notification_type
        if command.channel is not None:
            template.channel = command.channel
        if command.title_template is not None:
            template.title_template = command.title_template
        if command.content_template is not None:
            template.content_template = command.content_template
        template.updated_at = now
        return await self._template_repo.save(template)

    async def delete_template(self, uid: UUID) -> bool:
        return await self._template_repo.delete_by_uid(uid)
