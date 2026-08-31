from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus
from internal.domain.error import UnknownReferenceError
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.ticket_controller import create_ticket_router
from tests.conftest import make_ticket


@pytest.fixture
def service() -> MagicMock:
    return MagicMock(spec=SupportService)


STAFF_UID = uuid4()
TARGET_CUSTOMER_UID = uuid4()


@pytest.fixture
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_ticket_router(service))
    return TestClient(app, headers={"x-user-id": str(STAFF_UID), "x-user-roles": "MANAGE"})


def test_create_ticket_blank_subject_unprocessable(client: TestClient, service: MagicMock) -> None:
    response = client.post(
        "/api/v1/tickets",
        json={"customer_uid": str(uuid4()), "subject": "  ", "description": "Pages missing"},
    )
    assert response.status_code == 422
    service.create_ticket.assert_not_called()


def test_create_ticket_over_length_subject_unprocessable(client: TestClient, service: MagicMock) -> None:
    response = client.post(
        "/api/v1/tickets",
        json={"customer_uid": str(uuid4()), "subject": "x" * 256, "description": "Pages missing"},
    )
    assert response.status_code == 422
    service.create_ticket.assert_not_called()


def test_create_ticket_unknown_category_bad_request(client: TestClient, service: MagicMock) -> None:
    service.create_ticket.side_effect = UnknownReferenceError("Category not found")
    response = client.post(
        "/api/v1/tickets",
        json={
            "customer_uid": str(uuid4()),
            "subject": "Broken book",
            "description": "Pages missing",
            "category_uid": str(uuid4()),
        },
    )
    assert response.status_code == 400


def test_create_ticket_valid_request_created(client: TestClient, service: MagicMock) -> None:
    ticket = make_ticket()
    service.create_ticket.return_value = ticket
    response = client.post(
        "/api/v1/tickets",
        json={
            "customer_uid": str(TARGET_CUSTOMER_UID),
            "subject": "Help",
            "description": "issue",
            "priority": "HIGH",
        },
    )
    assert response.status_code == 201
    assert response.json()["uid"] == str(ticket.uid)


def test_get_ticket_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_ticket.return_value = None
    response = client.get(f"/api/v1/tickets/{uuid4()}")
    assert response.status_code == 404


def test_get_ticket_found_ok(client: TestClient, service: MagicMock) -> None:
    ticket = make_ticket()
    service.get_ticket.return_value = ticket
    response = client.get(f"/api/v1/tickets/{ticket.uid}")
    assert response.status_code == 200


def test_list_tickets_with_filters_passes_them_to_service(client: TestClient, service: MagicMock) -> None:
    service.list_tickets.return_value = ([make_ticket()], 1)
    response = client.get(
        "/api/v1/tickets",
        params={
            "customer_uid": str(TARGET_CUSTOMER_UID),
            "status": "OPEN",
            "priority": "URGENT",
            "size": 50,
        },
    )
    assert response.status_code == 200
    body = response.json()
    assert body["size"] == 50
    option = service.list_tickets.call_args[0][0]
    assert option.customer_uid == TARGET_CUSTOMER_UID
    assert option.status == TicketStatus.OPEN
    assert option.priority == TicketPriority.URGENT


def test_update_ticket_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.update_ticket.return_value = None
    response = client.patch(f"/api/v1/tickets/{uuid4()}", json={"subject": "x"})
    assert response.status_code == 404


def test_update_ticket_found_ok(client: TestClient, service: MagicMock) -> None:
    ticket = make_ticket()
    service.update_ticket.return_value = ticket
    response = client.patch(f"/api/v1/tickets/{ticket.uid}", json={"subject": "Updated"})
    assert response.status_code == 200


def test_update_ticket_status_found_ok(client: TestClient, service: MagicMock) -> None:
    ticket = make_ticket(status=TicketStatus.RESOLVED)
    service.update_ticket_status.return_value = ticket
    response = client.patch(f"/api/v1/tickets/{ticket.uid}/status", json={"status": "RESOLVED"})
    assert response.status_code == 200
    assert response.json()["status"] == "RESOLVED"


def test_update_ticket_status_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.update_ticket_status.return_value = None
    response = client.patch(f"/api/v1/tickets/{uuid4()}/status", json={"status": "OPEN"})
    assert response.status_code == 404
