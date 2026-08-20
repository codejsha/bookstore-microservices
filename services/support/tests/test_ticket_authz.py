from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.comment_controller import create_comment_router
from internal.infrastructure.adapter.restcontroller.ticket_controller import create_ticket_router
from tests.conftest import make_comment, make_ticket

CUSTOMER_UID = uuid4()
OTHER_UID = uuid4()
CUSTOMER_HEADERS = {"x-user-id": str(CUSTOMER_UID), "x-user-roles": "PROFILE,ORDER,VIEW"}


@pytest.fixture
def service() -> MagicMock:
    return MagicMock(spec=SupportService)


@pytest.fixture
def ticket_app(service: MagicMock) -> FastAPI:
    app = FastAPI()
    app.include_router(create_ticket_router(service))
    return app


@pytest.fixture
def comment_app(service: MagicMock) -> FastAPI:
    app = FastAPI()
    app.include_router(create_comment_router(service))
    return app


def _customer(app: FastAPI) -> TestClient:
    return TestClient(app, headers=CUSTOMER_HEADERS)


# ─── Tickets ──────────────────────────────────────────────────────────────


def test_missing_auth_returns_401(ticket_app: FastAPI) -> None:
    assert TestClient(ticket_app).get("/api/v1/tickets").status_code == 401


def test_non_uuid_subject_is_rejected(ticket_app: FastAPI, service: MagicMock) -> None:
    client = TestClient(ticket_app, headers={"x-user-id": "7", "x-user-roles": "PROFILE,VIEW"})
    assert client.get("/api/v1/tickets").status_code == 403


def test_list_pinned_to_caller(ticket_app: FastAPI, service: MagicMock) -> None:
    service.list_tickets.return_value = ([], 0)
    response = _customer(ticket_app).get("/api/v1/tickets", params={"customer_uid": str(OTHER_UID)})
    assert response.status_code == 200
    option = service.list_tickets.call_args[0][0]
    assert option.customer_uid == CUSTOMER_UID


def test_read_other_customers_ticket_forbidden(ticket_app: FastAPI, service: MagicMock) -> None:
    service.get_ticket.return_value = make_ticket(customer_uid=OTHER_UID)
    response = _customer(ticket_app).get(f"/api/v1/tickets/{uuid4()}")
    assert response.status_code == 403


def test_read_own_ticket_allowed(ticket_app: FastAPI, service: MagicMock) -> None:
    service.get_ticket.return_value = make_ticket(customer_uid=CUSTOMER_UID)
    response = _customer(ticket_app).get(f"/api/v1/tickets/{uuid4()}")
    assert response.status_code == 200


def test_status_change_requires_staff(ticket_app: FastAPI, service: MagicMock) -> None:
    response = _customer(ticket_app).patch(f"/api/v1/tickets/{uuid4()}/status", json={"status": "RESOLVED"})
    assert response.status_code == 403


def test_customer_cannot_assign(ticket_app: FastAPI, service: MagicMock) -> None:
    service.get_ticket.return_value = make_ticket(customer_uid=CUSTOMER_UID)
    response = _customer(ticket_app).patch(f"/api/v1/tickets/{uuid4()}", json={"assignee_uid": str(OTHER_UID)})
    assert response.status_code == 403


# ─── Comments ─────────────────────────────────────────────────────────────


def test_internal_comments_hidden_from_customer(comment_app: FastAPI, service: MagicMock) -> None:
    service.get_ticket.return_value = make_ticket(customer_uid=CUSTOMER_UID)
    service.list_comments.return_value = [make_comment()]
    response = _customer(comment_app).get(f"/api/v1/tickets/{uuid4()}/comments", params={"include_internal": "true"})
    assert response.status_code == 200
    assert service.list_comments.call_args.kwargs["include_internal"] is False


def test_customer_cannot_post_internal_comment(comment_app: FastAPI, service: MagicMock) -> None:
    service.get_ticket.return_value = make_ticket(customer_uid=CUSTOMER_UID)
    response = _customer(comment_app).post(
        f"/api/v1/tickets/{uuid4()}/comments",
        json={"author_uid": str(CUSTOMER_UID), "author_role": "CUSTOMER", "body": "hi", "internal": True},
    )
    assert response.status_code == 403


def test_comment_on_others_ticket_forbidden(comment_app: FastAPI, service: MagicMock) -> None:
    service.get_ticket.return_value = make_ticket(customer_uid=OTHER_UID)
    response = _customer(comment_app).get(f"/api/v1/tickets/{uuid4()}/comments")
    assert response.status_code == 403
