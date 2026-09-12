from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.restcontroller.notification_controller import create_notification_router
from internal.infrastructure.adapter.restcontroller.template_controller import create_template_router
from tests.conftest import make_notification, make_template

CUSTOMER_UID = str(uuid4())
OTHER_UID = str(uuid4())
STAFF_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "STAFF,USER"}
SYSTEM_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "SYSTEM"}
SEND_PAYLOAD = {
    "user_uid": OTHER_UID,
    "notification_type": "ORDER_PLACED",
    "channel": "EMAIL",
    "title": "x",
    "content": "y",
}
TEMPLATE_PAYLOAD = {
    "notification_type": "ORDER_PLACED",
    "channel": "EMAIL",
    "title_template": "Hi {name}",
    "content_template": "Body",
}


@pytest.fixture
def service() -> MagicMock:
    return MagicMock(spec=NotificationService)


@pytest.fixture
def app(service: MagicMock) -> FastAPI:
    application = FastAPI()
    application.include_router(create_notification_router(service))
    return application


@pytest.fixture
def template_app(service: MagicMock) -> FastAPI:
    application = FastAPI()
    application.include_router(create_template_router(service))
    return application


def _customer_client(app: FastAPI) -> TestClient:
    return TestClient(app, headers={"x-user-id": CUSTOMER_UID, "x-user-roles": "USER"})


def test_notifications_anonymous_caller_unauthorized(app: FastAPI) -> None:
    response = TestClient(app).get("/api/v1/notifications")
    assert response.status_code == 401


def test_notifications_corrupt_jwt_payload_unauthorized(app: FastAPI) -> None:
    response = TestClient(app).get(
        "/api/v1/notifications",
        headers={"x-user-id": CUSTOMER_UID, "x-jwt-payload": "bm90IGpzb24="},
    )
    assert response.status_code == 401


def test_list_notifications_customer_caller_pins_filter_to_caller(app: FastAPI, service: MagicMock) -> None:
    service.list_notifications.return_value = ([], 0)
    response = _customer_client(app).get("/api/v1/notifications", params={"user_uid": OTHER_UID})
    assert response.status_code == 200
    option = service.list_notifications.call_args[0][0]
    assert str(option.user_uid) == CUSTOMER_UID


def test_get_notification_owned_by_another_user_forbidden(app: FastAPI, service: MagicMock) -> None:
    service.get_notification.return_value = make_notification(user_uid=uuid4())
    response = _customer_client(app).get(f"/api/v1/notifications/{uuid4()}")
    assert response.status_code == 403


def test_get_notification_owned_by_caller_ok(app: FastAPI, service: MagicMock) -> None:
    from uuid import UUID

    service.get_notification.return_value = make_notification(user_uid=UUID(CUSTOMER_UID))
    response = _customer_client(app).get(f"/api/v1/notifications/{uuid4()}")
    assert response.status_code == 200


def test_send_notification_customer_caller_forbidden(app: FastAPI, service: MagicMock) -> None:
    response = _customer_client(app).post("/api/v1/notifications", json=SEND_PAYLOAD)
    assert response.status_code == 403


def test_send_notification_staff_caller_created(app: FastAPI, service: MagicMock) -> None:
    service.send_notification.return_value = make_notification()
    response = TestClient(app, headers=STAFF_HEADERS).post("/api/v1/notifications", json=SEND_PAYLOAD)
    assert response.status_code == 201


def test_send_notification_system_caller_created(app: FastAPI, service: MagicMock) -> None:
    service.send_notification.return_value = make_notification()
    response = TestClient(app, headers=SYSTEM_HEADERS).post("/api/v1/notifications", json=SEND_PAYLOAD)
    assert response.status_code == 201


def test_create_template_staff_caller_forbidden(template_app: FastAPI, service: MagicMock) -> None:
    response = TestClient(template_app, headers=STAFF_HEADERS).post("/api/v1/templates", json=TEMPLATE_PAYLOAD)
    assert response.status_code == 403


def test_create_template_system_caller_created(template_app: FastAPI, service: MagicMock) -> None:
    service.create_template.return_value = make_template()
    response = TestClient(template_app, headers=SYSTEM_HEADERS).post("/api/v1/templates", json=TEMPLATE_PAYLOAD)
    assert response.status_code == 201
