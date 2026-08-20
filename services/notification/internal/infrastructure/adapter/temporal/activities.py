from uuid import UUID

import structlog
from sqlalchemy.exc import IntegrityError
from temporalio import activity

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.model.command.notification_command import SendNotificationCommand
from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.temporal.models import ShipmentUpdateInput

_logger = structlog.get_logger()


class NotificationActivities:
    def __init__(self, notification_service: NotificationService):
        self._notification_service = notification_service

    @activity.defn(name="SendShipmentUpdate")
    async def send_shipment_update(self, payload: ShipmentUpdateInput) -> None:
        order_uid = UUID(payload.orderUid)
        try:
            notification_type = NotificationType[payload.notificationType]
        except KeyError:
            _logger.warning(
                "unknown notification_type, defaulting to SHIPMENT_DISPATCHED",
                notification_type=payload.notificationType,
            )
            notification_type = NotificationType.SHIPMENT_DISPATCHED

        existing = await self._notification_service.find_by_order_and_type(order_uid, notification_type)
        if existing is not None:
            _logger.info(
                "shipment update notification already persisted; skipping duplicate",
                notification_uid=str(existing.uid),
                order_uid=payload.orderUid,
                notification_type=notification_type.value,
                status=existing.status.value,
            )
            return

        content = f"Your order {payload.orderUid} update: {notification_type.value}."
        if payload.trackingNumber:
            content += f" Tracking number: {payload.trackingNumber}."

        command = SendNotificationCommand(
            user_uid=UUID(payload.userUid),
            order_uid=order_uid,
            notification_type=notification_type,
            channel=Channel.EMAIL,
            title="Your order has shipped",
            content=content,
        )

        try:
            notification = await self._notification_service.send_notification(command)
        except IntegrityError:
            existing = await self._notification_service.find_by_order_and_type(order_uid, notification_type)
            if existing is None:
                raise
            _logger.info(
                "shipment update notification lost insert race; duplicate skipped",
                notification_uid=str(existing.uid),
                order_uid=payload.orderUid,
                notification_type=notification_type.value,
            )
            return

        _logger.info(
            "shipment update notification queued",
            notification_uid=str(notification.uid),
            order_uid=payload.orderUid,
            notification_type=notification_type.value,
            status=notification.status.value,
        )
