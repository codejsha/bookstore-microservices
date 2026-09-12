from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.error import UnknownReferenceError
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.faq_controller import create_faq_router
from tests.conftest import make_faq


@pytest.fixture
def service() -> MagicMock:
    return MagicMock(spec=SupportService)


@pytest.fixture
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_faq_router(service))
    return TestClient(app, headers={"x-user-id": "42", "x-user-roles": "MANAGE,STAFF,USER"})


@pytest.fixture
def customer_client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_faq_router(service))
    return TestClient(app, headers={"x-user-id": "7", "x-user-roles": "USER"})


def test_create_faq_unknown_category_bad_request(client: TestClient, service: MagicMock) -> None:
    service.create_faq.side_effect = UnknownReferenceError("Category not found")
    response = client.post(
        "/api/v1/faqs",
        json={"question": "How to refund?", "answer": "Go to orders.", "category_uid": str(uuid4())},
    )
    assert response.status_code == 400


def test_create_faq_valid_request_created(client: TestClient, service: MagicMock) -> None:
    faq = make_faq()
    service.create_faq.return_value = faq
    response = client.post(
        "/api/v1/faqs",
        json={"question": "Q?", "answer": "A.", "published": True},
    )
    assert response.status_code == 201


def test_get_faq_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_faq.return_value = None
    response = client.get(f"/api/v1/faqs/{uuid4()}")
    assert response.status_code == 404


def test_get_faq_found_ok(client: TestClient, service: MagicMock) -> None:
    faq = make_faq()
    service.get_faq.return_value = faq
    response = client.get(f"/api/v1/faqs/{faq.uid}")
    assert response.status_code == 200


def test_get_faq_with_increment_view_passes_it_to_service(client: TestClient, service: MagicMock) -> None:
    service.get_faq.return_value = make_faq()
    client.get(f"/api/v1/faqs/{uuid4()}", params={"increment_view": "false"})
    _, kwargs = service.get_faq.call_args
    assert kwargs["increment_view"] is False


def test_search_faqs_staff_caller_honors_published_only_flag(client: TestClient, service: MagicMock) -> None:
    service.search_faqs.return_value = ([make_faq()], 1)
    response = client.get(
        "/api/v1/faqs",
        params={"q": "refund", "published_only": "false", "size": 5},
    )
    assert response.status_code == 200
    option = service.search_faqs.call_args[0][0]
    assert option.query == "refund"
    assert option.published_only is False
    assert option.size == 5


def test_search_faqs_customer_caller_forces_published_only(customer_client: TestClient, service: MagicMock) -> None:
    service.search_faqs.return_value = ([make_faq()], 1)
    response = customer_client.get(
        "/api/v1/faqs",
        params={"q": "refund", "published_only": "false", "size": 5},
    )
    assert response.status_code == 200
    option = service.search_faqs.call_args[0][0]
    assert option.published_only is True
    assert option.query == "refund"
    assert option.size == 5


def test_get_faq_unpublished_and_customer_caller_not_found(customer_client: TestClient, service: MagicMock) -> None:
    service.get_faq.return_value = make_faq(published=False)
    response = customer_client.get(f"/api/v1/faqs/{uuid4()}")
    assert response.status_code == 404


def test_get_faq_unpublished_and_staff_caller_ok(client: TestClient, service: MagicMock) -> None:
    service.get_faq.return_value = make_faq(published=False)
    response = client.get(f"/api/v1/faqs/{uuid4()}")
    assert response.status_code == 200


def test_update_faq_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.update_faq.return_value = None
    response = client.patch(f"/api/v1/faqs/{uuid4()}", json={"question": "x"})
    assert response.status_code == 404


def test_update_faq_found_ok(client: TestClient, service: MagicMock) -> None:
    faq = make_faq()
    service.update_faq.return_value = faq
    response = client.patch(f"/api/v1/faqs/{faq.uid}", json={"question": "Updated"})
    assert response.status_code == 200


def test_delete_faq_found_no_content(client: TestClient, service: MagicMock) -> None:
    service.delete_faq.return_value = True
    response = client.delete(f"/api/v1/faqs/{uuid4()}")
    assert response.status_code == 204


def test_delete_faq_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.delete_faq.return_value = False
    response = client.delete(f"/api/v1/faqs/{uuid4()}")
    assert response.status_code == 404
