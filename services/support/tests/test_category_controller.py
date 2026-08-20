from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.restcontroller.category_controller import create_category_router
from tests.conftest import make_category


@pytest.fixture
def service() -> MagicMock:
    return MagicMock(spec=SupportService)


@pytest.fixture
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_category_router(service))
    return TestClient(app, headers={"x-user-id": "42", "x-user-roles": "MANAGE"})


def test_create_category_returns_201(client: TestClient, service: MagicMock) -> None:
    category = make_category()
    service.create_category.return_value = category
    response = client.post(
        "/api/v1/categories",
        json={"name": "Billing"},
    )
    assert response.status_code == 201
    assert response.json()["name"] == category.name


def test_get_category_returns_404(client: TestClient, service: MagicMock) -> None:
    service.get_category.return_value = None
    response = client.get(f"/api/v1/categories/{uuid4()}")
    assert response.status_code == 404


def test_get_category_returns_200(client: TestClient, service: MagicMock) -> None:
    category = make_category()
    service.get_category.return_value = category
    response = client.get(f"/api/v1/categories/{category.uid}")
    assert response.status_code == 200


def test_list_categories(client: TestClient, service: MagicMock) -> None:
    service.list_categories.return_value = [make_category(), make_category(name="Other")]
    response = client.get("/api/v1/categories")
    assert response.status_code == 200
    assert len(response.json()) == 2


def test_update_category_returns_404(client: TestClient, service: MagicMock) -> None:
    service.update_category.return_value = None
    response = client.patch(f"/api/v1/categories/{uuid4()}", json={"name": "X"})
    assert response.status_code == 404


def test_update_category_returns_200(client: TestClient, service: MagicMock) -> None:
    category = make_category()
    service.update_category.return_value = category
    response = client.patch(f"/api/v1/categories/{category.uid}", json={"name": "New"})
    assert response.status_code == 200


def test_delete_category_returns_204(client: TestClient, service: MagicMock) -> None:
    service.delete_category.return_value = True
    response = client.delete(f"/api/v1/categories/{uuid4()}")
    assert response.status_code == 204


def test_delete_category_returns_404(client: TestClient, service: MagicMock) -> None:
    service.delete_category.return_value = False
    response = client.delete(f"/api/v1/categories/{uuid4()}")
    assert response.status_code == 404
