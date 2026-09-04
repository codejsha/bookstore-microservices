from dataclasses import fields
from unittest.mock import AsyncMock
from uuid import UUID, uuid4

import pytest
from sqlalchemy.exc import IntegrityError

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus
from internal.domain.model.command.notification_command import SendNotificationCommand
from internal.infrastructure.adapter.temporal.activities import NotificationActivities
from internal.infrastructure.adapter.temporal.models import ShipmentUpdateInput
from tests.conftest import make_notification


@pytest.fixture
def notification_service() -> AsyncMock:
    service = AsyncMock()
    service.send_notification.return_value.uid = uuid4()
    service.send_notification.return_value.status.value = "PENDING"
    service.find_by_order_and_type = AsyncMock(return_value=None)
    return service


@pytest.fixture
def activities(notification_service: AsyncMock) -> NotificationActivities:
    return NotificationActivities(notification_service=notification_service)


def _sent_command(notification_service: AsyncMock) -> SendNotificationCommand:
    notification_service.send_notification.assert_awaited_once()
    (command,) = notification_service.send_notification.await_args.args
    return command


class TestSendShipmentUpdate:
    async def test_send_shipment_update_payload_camelcase_builds_send_command(
        self, activities: NotificationActivities, notification_service: AsyncMock
    ) -> None:
        user_uid = uuid4()
        order_uid = uuid4()
        payload = ShipmentUpdateInput(
            userUid=str(user_uid),
            orderUid=str(order_uid),
            notificationType="SHIPMENT_DISPATCHED",
            trackingNumber="1Z-TRACK-999",
        )

        await activities.send_shipment_update(payload)

        command = _sent_command(notification_service)
        assert isinstance(command, SendNotificationCommand)
        assert command.user_uid == UUID(payload.userUid)
        assert command.order_uid == order_uid
        assert command.notification_type == NotificationType.SHIPMENT_DISPATCHED
        assert command.channel == Channel.EMAIL
        assert str(order_uid) in command.content
        assert "1Z-TRACK-999" in command.content

    async def test_send_shipment_update_unknown_notification_type_falls_back_to_dispatched(
        self, activities: NotificationActivities, notification_service: AsyncMock
    ) -> None:
        payload = ShipmentUpdateInput(
            userUid=str(uuid4()),
            orderUid=str(uuid4()),
            notificationType="NOT_A_REAL_TYPE",
            trackingNumber="TRK-1",
        )

        await activities.send_shipment_update(payload)

        command = _sent_command(notification_service)
        assert command.notification_type == NotificationType.SHIPMENT_DISPATCHED

    async def test_send_shipment_update_tracking_number_missing_still_builds_valid_command(
        self, activities: NotificationActivities, notification_service: AsyncMock
    ) -> None:
        order_uid = uuid4()
        payload = ShipmentUpdateInput(
            userUid=str(uuid4()),
            orderUid=str(order_uid),
            notificationType="SHIPMENT_DELIVERED",
            trackingNumber=None,
        )

        await activities.send_shipment_update(payload)

        command = _sent_command(notification_service)
        assert command.notification_type == NotificationType.SHIPMENT_DELIVERED
        assert command.channel == Channel.EMAIL
        assert str(order_uid) in command.content
        assert "Tracking number" not in command.content

    async def test_send_shipment_update_without_tracking_number_leaves_it_out_of_content(
        self, activities: NotificationActivities, notification_service: AsyncMock
    ) -> None:
        payload = ShipmentUpdateInput(
            userUid=str(uuid4()),
            orderUid=str(uuid4()),
            notificationType="SHIPMENT_DISPATCHED",
        )
        assert payload.trackingNumber is None

        await activities.send_shipment_update(payload)

        command = _sent_command(notification_service)
        assert "Tracking number" not in command.content


class TestSendShipmentUpdateIdempotency:
    async def test_send_shipment_update_existing_notification_already_skips_send(
        self, activities: NotificationActivities, notification_service: AsyncMock
    ) -> None:
        order_uid = uuid4()
        existing = make_notification(
            order_uid=order_uid,
            notification_type=NotificationType.SHIPMENT_DISPATCHED,
            status=NotificationStatus.SENT,
        )
        notification_service.find_by_order_and_type = AsyncMock(return_value=existing)
        payload = ShipmentUpdateInput(
            userUid=str(uuid4()),
            orderUid=str(order_uid),
            notificationType="SHIPMENT_DISPATCHED",
            trackingNumber="TRK-1",
        )

        await activities.send_shipment_update(payload)

        notification_service.find_by_order_and_type.assert_awaited_once_with(
            order_uid, NotificationType.SHIPMENT_DISPATCHED
        )
        notification_service.send_notification.assert_not_awaited()

    async def test_send_shipment_update_insert_race_falls_back_to_existing_row(
        self, activities: NotificationActivities, notification_service: AsyncMock
    ) -> None:
        order_uid = uuid4()
        winner = make_notification(
            order_uid=order_uid,
            notification_type=NotificationType.SHIPMENT_DISPATCHED,
        )
        notification_service.find_by_order_and_type = AsyncMock(side_effect=[None, winner])
        notification_service.send_notification = AsyncMock(
            side_effect=IntegrityError("INSERT", {}, Exception("duplicate order_uid"))
        )
        payload = ShipmentUpdateInput(
            userUid=str(uuid4()),
            orderUid=str(order_uid),
            notificationType="SHIPMENT_DISPATCHED",
            trackingNumber="TRK-1",
        )

        await activities.send_shipment_update(payload)

        assert notification_service.find_by_order_and_type.await_count == 2
        notification_service.send_notification.assert_awaited_once()

    async def test_send_shipment_update_insert_races_and_no_row_found_reraises_integrity_error(
        self, activities: NotificationActivities, notification_service: AsyncMock
    ) -> None:
        order_uid = uuid4()
        notification_service.find_by_order_and_type = AsyncMock(side_effect=[None, None])
        notification_service.send_notification = AsyncMock(
            side_effect=IntegrityError("INSERT", {}, Exception("some other constraint"))
        )
        payload = ShipmentUpdateInput(
            userUid=str(uuid4()),
            orderUid=str(order_uid),
            notificationType="SHIPMENT_DISPATCHED",
        )

        with pytest.raises(IntegrityError):
            await activities.send_shipment_update(payload)


class TestShipmentUpdateInputContract:
    def test_shipment_update_input_fields_listed_are_exact_camelcase(self) -> None:
        names = {f.name for f in fields(ShipmentUpdateInput)}
        assert names == {"userUid", "orderUid", "notificationType", "trackingNumber"}

    def test_shipment_update_input_without_tracking_number_defaults_to_none(self) -> None:
        (tracking_field,) = (f for f in fields(ShipmentUpdateInput) if f.name == "trackingNumber")
        assert tracking_field.default is None
