from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.constant.author_role import AuthorRole
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.comment_controller import create_comment_router
from tests.conftest import make_comment, make_ticket

STAFF_UID = uuid4()
CUSTOMER_UID = uuid4()


@pytest.fixture
def service() -> MagicMock:
    return MagicMock(spec=SupportService)


@pytest.fixture
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_comment_router(service))
    return TestClient(app, headers={"x-user-id": str(STAFF_UID), "x-user-roles": "MANAGE"})


@pytest.fixture
def customer_client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_comment_router(service))
    return TestClient(app, headers={"x-user-id": str(CUSTOMER_UID), "x-user-roles": "VIEW"})


def test_add_comment_returns_201(client: TestClient, service: MagicMock) -> None:
    comment = make_comment()
    service.add_comment.return_value = comment
    response = client.post(
        f"/api/v1/tickets/{uuid4()}/comments",
        json={"author_uid": str(STAFF_UID), "author_role": "CUSTOMER", "body": "Hi"},
    )
    assert response.status_code == 201


def test_add_comment_returns_404_when_ticket_missing(client: TestClient, service: MagicMock) -> None:
    service.add_comment.return_value = None
    response = client.post(
        f"/api/v1/tickets/{uuid4()}/comments",
        json={"author_uid": str(STAFF_UID), "author_role": "CUSTOMER", "body": "Hi"},
    )
    assert response.status_code == 404


def test_list_comments_returns_404_when_ticket_missing(client: TestClient, service: MagicMock) -> None:
    service.list_comments.return_value = None
    response = client.get(f"/api/v1/tickets/{uuid4()}/comments")
    assert response.status_code == 404


def test_list_comments_returns_200(client: TestClient, service: MagicMock) -> None:
    service.list_comments.return_value = [make_comment(), make_comment(internal=True)]
    response = client.get(f"/api/v1/tickets/{uuid4()}/comments")
    assert response.status_code == 200
    assert len(response.json()) == 2


def test_non_staff_cannot_forge_author_uid_or_role(customer_client: TestClient, service: MagicMock) -> None:
    service.get_ticket.return_value = make_ticket(customer_uid=CUSTOMER_UID)
    service.add_comment.return_value = make_comment()
    response = customer_client.post(
        f"/api/v1/tickets/{uuid4()}/comments",
        json={"author_uid": str(uuid4()), "author_role": "AGENT", "body": "Hi"},
    )
    assert response.status_code == 201
    command = service.add_comment.call_args[0][0]
    assert command.author_uid == CUSTOMER_UID
    assert command.author_role == AuthorRole.CUSTOMER


def test_staff_comment_uses_own_uid_and_agent_role(client: TestClient, service: MagicMock) -> None:
    service.get_ticket.return_value = make_ticket()
    service.add_comment.return_value = make_comment(author_role=AuthorRole.AGENT)
    response = client.post(
        f"/api/v1/tickets/{uuid4()}/comments",
        json={"body": "Internal note"},
    )
    assert response.status_code == 201
    command = service.add_comment.call_args[0][0]
    assert command.author_uid == STAFF_UID
    assert command.author_role == AuthorRole.AGENT


def test_list_comments_passes_include_internal_flag(client: TestClient, service: MagicMock) -> None:
    service.list_comments.return_value = [make_comment()]
    response = client.get(f"/api/v1/tickets/{uuid4()}/comments", params={"include_internal": "false"})
    assert response.status_code == 200
    _, kwargs = service.list_comments.call_args
    assert kwargs["include_internal"] is False
