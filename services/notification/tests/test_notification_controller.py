from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus
from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.restcontroller.notification_controller import (
    create_notification_router,
)
from tests.conftest import make_notification


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=NotificationService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_notification_router(service))
    return TestClient(app, headers={"x-user-id": "42", "x-user-roles": "MANAGE"})


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def test_send_notification_returns_201(client: TestClient, service: MagicMock) -> None:
    notification = make_notification()
    service.send_notification.return_value = notification
    response = client.post(
        "/api/v1/notifications",
        json={
            "user_uid": str(uuid4()),
            "notification_type": "ORDER_PLACED",
            "channel": "EMAIL",
            "title": "Hello",
            "content": "Body",
        },
    )
    assert response.status_code == 201
    assert response.content == b""
    assert response.headers["location"] == f"/api/v1/notifications/{notification.uid}"


def test_get_notification_returns_404_when_missing(client: TestClient, service: MagicMock) -> None:
    service.get_notification.return_value = None
    response = client.get(f"/api/v1/notifications/{uuid4()}")
    assert response.status_code == 404


def test_get_notification_returns_200(client: TestClient, service: MagicMock) -> None:
    notification = make_notification()
    service.get_notification.return_value = notification
    response = client.get(f"/api/v1/notifications/{notification.uid}")
    assert response.status_code == 200
    assert response.json()["title"] == notification.title


def test_list_passes_filters(client: TestClient, service: MagicMock) -> None:
    service.list_notifications.return_value = ([make_notification()], 1)
    response = client.get(
        "/api/v1/notifications",
        params={"channel": "SMS", "status": "PENDING", "page": 1, "size": 5},
    )
    assert response.status_code == 200
    body = response.json()
    assert body["total"] == 1
    assert len(body["items"]) == 1
    option = service.list_notifications.call_args[0][0]
    assert option.channel == Channel.SMS
    assert option.status == NotificationStatus.PENDING
    assert option.page == 1
    assert option.size == 5


def test_list_passes_notification_type_filter(client: TestClient, service: MagicMock) -> None:
    service.list_notifications.return_value = ([], 0)
    response = client.get("/api/v1/notifications", params={"notification_type": "ORDER_PLACED"})
    assert response.status_code == 200
    option = service.list_notifications.call_args[0][0]
    assert option.notification_type == NotificationType.ORDER_PLACED


def test_mark_as_sent_returns_404(client: TestClient, service: MagicMock) -> None:
    service.mark_as_sent.return_value = None
    response = client.patch(f"/api/v1/notifications/{uuid4()}/sent")
    assert response.status_code == 404


def test_mark_as_sent_returns_200(client: TestClient, service: MagicMock) -> None:
    notification = make_notification(status=NotificationStatus.SENT)
    service.mark_as_sent.return_value = notification
    response = client.patch(f"/api/v1/notifications/{notification.uid}/sent")
    assert response.status_code == 200
    assert response.json()["status"] == "SENT"


def test_mark_as_failed_returns_200(client: TestClient, service: MagicMock) -> None:
    notification = make_notification(status=NotificationStatus.FAILED)
    service.mark_as_failed.return_value = notification
    response = client.patch(f"/api/v1/notifications/{notification.uid}/failed")
    assert response.status_code == 200
    assert response.json()["status"] == "FAILED"


def test_mark_as_failed_returns_404(client: TestClient, service: MagicMock) -> None:
    service.mark_as_failed.return_value = None
    response = client.patch(f"/api/v1/notifications/{uuid4()}/failed")
    assert response.status_code == 404
