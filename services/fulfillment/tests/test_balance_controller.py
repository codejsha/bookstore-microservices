from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.service.inventory_service import BalanceService
from internal.infrastructure.adapter.restcontroller.balance_controller import create_balance_router
from tests.conftest import make_balance_entry

USER_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "USER"}


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=BalanceService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_balance_router(service))
    return TestClient(app, headers=USER_HEADERS)


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def test_get_stock_balance_with_year_month_passes_option(client: TestClient, service: MagicMock) -> None:
    service.get_stock_balance.return_value = ([make_balance_entry()], 1)

    response = client.get("/api/v1/balance", params={"year_month": "2026-10", "size": 5})

    assert response.status_code == 200
    body = response.json()
    assert body["total"] == 1
    assert body["items"][0]["current_quantity"] == 12
    option = service.get_stock_balance.call_args.args[0]
    assert (option.year_month, option.size) == ("2026-10", 5)


def test_get_stock_balance_malformed_year_month_bad_request(client: TestClient, service: MagicMock) -> None:
    response = client.get("/api/v1/balance", params={"year_month": "2026/10"})
    assert response.status_code == 400
    service.get_stock_balance.assert_not_called()


def test_get_stock_balance_anonymous_unauthorized(service: MagicMock) -> None:
    app = FastAPI()
    app.include_router(create_balance_router(service))
    assert TestClient(app).get("/api/v1/balance").status_code == 401
