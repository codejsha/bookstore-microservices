from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.model.error import DuplicateClosingError, WarehouseNotFoundError
from internal.domain.service.inventory_service import ClosingService
from internal.infrastructure.adapter.restcontroller.closing_controller import create_closing_router
from tests.conftest import make_closing

MANAGER_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "MANAGER,STAFF,USER"}


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=ClosingService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_closing_router(service))
    return TestClient(app, headers=MANAGER_HEADERS)


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def test_create_closing_valid_request_ok_with_items(client: TestClient, service: MagicMock) -> None:
    closing = make_closing()
    service.create_monthly_closing.return_value = closing

    response = client.post("/api/v1/closings", json={"warehouse_uid": str(uuid4()), "year": 2026, "month": 10})

    assert response.status_code == 200
    body = response.json()
    assert body["status"] == "CLOSED"
    assert body["items"][0]["closing_quantity"] == 12


def test_create_closing_month_out_of_range_unprocessable(client: TestClient, service: MagicMock) -> None:
    response = client.post("/api/v1/closings", json={"warehouse_uid": str(uuid4()), "year": 2026, "month": 13})
    assert response.status_code == 422
    service.create_monthly_closing.assert_not_called()


def test_create_closing_duplicate_period_conflict(client: TestClient, service: MagicMock) -> None:
    service.create_monthly_closing.side_effect = DuplicateClosingError("exists")
    response = client.post("/api/v1/closings", json={"warehouse_uid": str(uuid4()), "year": 2026, "month": 10})
    assert response.status_code == 409


def test_create_closing_missing_warehouse_not_found(client: TestClient, service: MagicMock) -> None:
    service.create_monthly_closing.side_effect = WarehouseNotFoundError("missing")
    response = client.post("/api/v1/closings", json={"warehouse_uid": str(uuid4()), "year": 2026, "month": 10})
    assert response.status_code == 404


def test_get_closing_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_closing.return_value = None
    assert client.get(f"/api/v1/closings/{uuid4()}").status_code == 404


def test_list_closings_with_period_filter_passes_option(client: TestClient, service: MagicMock) -> None:
    service.list_closings.return_value = ([make_closing()], 1)
    response = client.get("/api/v1/closings", params={"year": 2026, "month": 10})
    assert response.status_code == 200
    option = service.list_closings.call_args.args[0]
    assert (option.year, option.month) == (2026, 10)
