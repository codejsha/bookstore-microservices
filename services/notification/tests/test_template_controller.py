from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from sqlalchemy.exc import IntegrityError

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.error import TemplateAlreadyExistsError
from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.restcontroller.template_controller import create_template_router
from tests.conftest import make_template


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=NotificationService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_template_router(service))
    return TestClient(app, headers={"x-user-id": "42", "x-user-roles": "MANAGE"})


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def test_create_template_valid_request_created(client: TestClient, service: MagicMock) -> None:
    template = make_template()
    service.create_template.return_value = template
    response = client.post(
        "/api/v1/templates",
        json={
            "notification_type": "ORDER_PLACED",
            "channel": "EMAIL",
            "title_template": "Hi {name}",
            "content_template": "Body",
        },
    )
    assert response.status_code == 201
    assert response.content == b""
    assert response.headers["location"] == f"/api/v1/templates/{template.uid}"


def test_create_template_occupied_slot_conflict(client: TestClient, service: MagicMock) -> None:
    service.create_template.side_effect = TemplateAlreadyExistsError(NotificationType.ORDER_PLACED, Channel.EMAIL)
    response = client.post(
        "/api/v1/templates",
        json={
            "notification_type": "ORDER_PLACED",
            "channel": "EMAIL",
            "title_template": "Hi {name}",
            "content_template": "Body",
        },
    )
    assert response.status_code == 409


def test_create_template_unique_violation_conflict(client: TestClient, service: MagicMock) -> None:
    service.create_template.side_effect = IntegrityError("insert", {}, Exception("duplicate"))
    response = client.post(
        "/api/v1/templates",
        json={
            "notification_type": "ORDER_PLACED",
            "channel": "EMAIL",
            "title_template": "Hi {name}",
            "content_template": "Body",
        },
    )
    assert response.status_code == 409


def test_update_template_occupied_slot_conflict(client: TestClient, service: MagicMock) -> None:
    service.update_template.side_effect = TemplateAlreadyExistsError(NotificationType.ORDER_PLACED, Channel.SMS)
    response = client.put(
        f"/api/v1/templates/{uuid4()}",
        json={"channel": "SMS"},
    )
    assert response.status_code == 409


def test_get_template_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_template.return_value = None
    response = client.get(f"/api/v1/templates/{uuid4()}")
    assert response.status_code == 404


def test_get_template_found_ok(client: TestClient, service: MagicMock) -> None:
    template = make_template()
    service.get_template.return_value = template
    response = client.get(f"/api/v1/templates/{template.uid}")
    assert response.status_code == 200


def test_list_templates_found_ok_with_total(client: TestClient, service: MagicMock) -> None:
    service.list_templates.return_value = ([make_template(), make_template()], 2)
    response = client.get("/api/v1/templates")
    assert response.status_code == 200
    body = response.json()
    assert body["total"] == 2
    assert len(body["items"]) == 2


def test_list_templates_with_filters_passes_them_to_service(client: TestClient, service: MagicMock) -> None:
    service.list_templates.return_value = ([make_template()], 1)
    response = client.get(
        "/api/v1/templates",
        params={
            "notification_type": "ORDER_PLACED",
            "channel": "SMS",
            "page": 2,
            "size": 5,
            "sort": "updated_at:asc",
        },
    )
    assert response.status_code == 200
    option = service.list_templates.call_args[0][0]
    assert option.notification_type == NotificationType.ORDER_PLACED
    assert option.channel == Channel.SMS
    assert option.page == 2
    assert option.size == 5
    assert option.sort == "updated_at:asc"


def test_update_template_blank_title_unprocessable(client: TestClient, service: MagicMock) -> None:
    response = client.put(
        f"/api/v1/templates/{uuid4()}",
        json={"title_template": "  "},
    )
    assert response.status_code == 422
    service.update_template.assert_not_called()


def test_create_template_over_length_title_unprocessable(client: TestClient, service: MagicMock) -> None:
    response = client.post(
        "/api/v1/templates",
        json={
            "notification_type": "ORDER_PLACED",
            "channel": "EMAIL",
            "title_template": "x" * 256,
            "content_template": "Body",
        },
    )
    assert response.status_code == 422
    service.create_template.assert_not_called()


def test_update_template_over_length_content_unprocessable(client: TestClient, service: MagicMock) -> None:
    response = client.put(
        f"/api/v1/templates/{uuid4()}",
        json={"content_template": "x" * 10001},
    )
    assert response.status_code == 422
    service.update_template.assert_not_called()


def test_update_template_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.update_template.return_value = None
    response = client.put(
        f"/api/v1/templates/{uuid4()}",
        json={"title_template": "New"},
    )
    assert response.status_code == 404


def test_update_template_found_ok(client: TestClient, service: MagicMock) -> None:
    template = make_template()
    service.update_template.return_value = template
    response = client.put(
        f"/api/v1/templates/{template.uid}",
        json={"title_template": "Updated"},
    )
    assert response.status_code == 200


def test_delete_template_found_no_content(client: TestClient, service: MagicMock) -> None:
    service.delete_template.return_value = True
    response = client.delete(f"/api/v1/templates/{uuid4()}")
    assert response.status_code == 204


def test_delete_template_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.delete_template.return_value = False
    response = client.delete(f"/api/v1/templates/{uuid4()}")
    assert response.status_code == 404
