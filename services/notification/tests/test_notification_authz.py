from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.restcontroller.notification_controller import create_notification_router
from tests.conftest import make_notification

CUSTOMER_UID = str(uuid4())
OTHER_UID = str(uuid4())


@pytest.fixture
def service() -> MagicMock:
    return MagicMock(spec=NotificationService)


@pytest.fixture
def app(service: MagicMock) -> FastAPI:
    application = FastAPI()
    application.include_router(create_notification_router(service))
    return application


def _customer_client(app: FastAPI) -> TestClient:
    return TestClient(app, headers={"x-user-id": CUSTOMER_UID, "x-user-roles": "PROFILE,ORDER,VIEW"})


def test_missing_auth_returns_401(app: FastAPI) -> None:
    response = TestClient(app).get("/api/v1/notifications")
    assert response.status_code == 401


def test_corrupt_jwt_payload_returns_401(app: FastAPI) -> None:
    response = TestClient(app).get(
        "/api/v1/notifications",
        headers={"x-user-id": CUSTOMER_UID, "x-jwt-payload": "bm90IGpzb24="},
    )
    assert response.status_code == 401


def test_list_is_pinned_to_caller(app: FastAPI, service: MagicMock) -> None:
    service.list_notifications.return_value = ([], 0)
    response = _customer_client(app).get("/api/v1/notifications", params={"user_uid": OTHER_UID})
    assert response.status_code == 200
    option = service.list_notifications.call_args[0][0]
    assert str(option.user_uid) == CUSTOMER_UID


def test_read_others_notification_forbidden(app: FastAPI, service: MagicMock) -> None:
    service.get_notification.return_value = make_notification(user_uid=uuid4())
    response = _customer_client(app).get(f"/api/v1/notifications/{uuid4()}")
    assert response.status_code == 403


def test_read_own_notification_allowed(app: FastAPI, service: MagicMock) -> None:
    from uuid import UUID

    service.get_notification.return_value = make_notification(user_uid=UUID(CUSTOMER_UID))
    response = _customer_client(app).get(f"/api/v1/notifications/{uuid4()}")
    assert response.status_code == 200


def test_customer_cannot_send_notification(app: FastAPI, service: MagicMock) -> None:
    response = _customer_client(app).post(
        "/api/v1/notifications",
        json={
            "user_uid": OTHER_UID,
            "notification_type": "ORDER_PLACED",
            "channel": "EMAIL",
            "title": "x",
            "content": "y",
        },
    )
    assert response.status_code == 403
