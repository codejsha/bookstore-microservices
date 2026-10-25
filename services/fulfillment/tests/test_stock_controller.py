from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.error import InsufficientStockError, StockNotFoundError, WarehouseNotFoundError
from internal.domain.service.inventory_service import StockService
from internal.infrastructure.adapter.restcontroller.stock_controller import create_stock_router
from tests.conftest import make_stock

STAFF_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "STAFF,USER"}


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=StockService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_stock_router(service))
    return TestClient(app, headers=STAFF_HEADERS)


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def _receive_body(**overrides) -> dict:
    body = {"edition_uid": str(uuid4()), "warehouse_uid": str(uuid4()), "quantity": 5, "reason": "restock"}
    body.update(overrides)
    return body


def test_list_stocks_with_filters_passes_option_to_service(client: TestClient, service: MagicMock) -> None:
    service.list_stocks.return_value = ([make_stock()], 1)
    warehouse_uid = uuid4()

    response = client.get("/api/v1/stocks", params={"warehouse_uid": str(warehouse_uid), "size": 5, "sort": "uid,asc"})

    assert response.status_code == 200
    assert response.json()["total"] == 1
    option = service.list_stocks.call_args.args[0]
    assert option.warehouse_uid == warehouse_uid
    assert option.size == 5
    assert option.sort == "uid,asc"


def test_list_stocks_malformed_uuid_filter_unprocessable(client: TestClient, service: MagicMock) -> None:
    assert client.get("/api/v1/stocks", params={"edition_uid": "nope"}).status_code == 422
    service.list_stocks.assert_not_called()


def test_get_stock_found_ok_with_warehouses(client: TestClient, service: MagicMock) -> None:
    stock = make_stock(quantity=7)
    service.get_stock.return_value = stock

    response = client.get(f"/api/v1/stocks/{stock.edition_uid}")

    assert response.status_code == 200
    body = response.json()
    assert body["edition_uid"] == str(stock.edition_uid)
    assert body["total_quantity"] == 7
    assert body["warehouses"][0]["warehouse_name"] == "Main"
    assert "id" not in body


def test_get_stock_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_stock.return_value = None
    assert client.get(f"/api/v1/stocks/{uuid4()}").status_code == 404


def test_receive_stock_valid_request_ok_with_body(client: TestClient, service: MagicMock) -> None:
    stock = make_stock()
    service.receive_stock.return_value = stock

    response = client.post("/api/v1/stocks/receive", json=_receive_body())

    assert response.status_code == 200
    assert response.json()["uid"] == str(stock.uid)
    command = service.receive_stock.call_args.args[0]
    assert command.quantity == 5 and command.reason == "restock"


def test_receive_stock_zero_quantity_unprocessable(client: TestClient, service: MagicMock) -> None:
    assert client.post("/api/v1/stocks/receive", json=_receive_body(quantity=0)).status_code == 422
    service.receive_stock.assert_not_called()


def test_receive_stock_malformed_uuid_bad_request(client: TestClient, service: MagicMock) -> None:
    response = client.post("/api/v1/stocks/receive", json=_receive_body(edition_uid="ed-1"))
    assert response.status_code == 400
    assert "edition_uid" in response.json()["detail"]
    service.receive_stock.assert_not_called()


def test_receive_stock_missing_stock_not_found(client: TestClient, service: MagicMock) -> None:
    service.receive_stock.side_effect = StockNotFoundError("stock not found")
    assert client.post("/api/v1/stocks/receive", json=_receive_body()).status_code == 404


def test_receive_stock_missing_warehouse_not_found(client: TestClient, service: MagicMock) -> None:
    service.receive_stock.side_effect = WarehouseNotFoundError("warehouse not found")
    assert client.post("/api/v1/stocks/receive", json=_receive_body()).status_code == 404


def test_reserve_stock_insufficient_bad_request(client: TestClient, service: MagicMock) -> None:
    service.reserve_stock.side_effect = InsufficientStockError("insufficient stock")
    response = client.post("/api/v1/stocks/reserve", json=_receive_body())
    assert response.status_code == 400
    assert response.json()["detail"] == "insufficient stock"


def test_adjust_stock_zero_quantity_bad_request(client: TestClient, service: MagicMock) -> None:
    response = client.post("/api/v1/stocks/adjust", json=_receive_body(quantity=0))
    assert response.status_code == 400
    service.adjust_stock.assert_not_called()


def test_adjust_stock_negative_quantity_ok(client: TestClient, service: MagicMock) -> None:
    service.adjust_stock.return_value = make_stock()
    assert client.post("/api/v1/stocks/adjust", json=_receive_body(quantity=-2)).status_code == 200
    assert service.adjust_stock.call_args.args[0].quantity == -2


def test_release_stock_valid_request_ok(client: TestClient, service: MagicMock) -> None:
    service.release_stock.return_value = make_stock()
    assert client.post("/api/v1/stocks/release", json=_receive_body()).status_code == 200


def test_get_stock_history_found_ok_with_items(client: TestClient, service: MagicMock) -> None:
    from datetime import UTC, datetime

    from internal.domain.aggregate.stock_aggregate import StockHistoryEntry

    entry = StockHistoryEntry(
        uid=uuid4(), change_type=StockChangeType.INBOUND, reason=None, change_qty=3, created_at=datetime.now(UTC)
    )
    service.get_stock_history.return_value = ([entry], 1)

    response = client.get(f"/api/v1/stocks/{uuid4()}/history", params={"warehouse_uid": str(uuid4()), "size": 10})

    assert response.status_code == 200
    body = response.json()
    assert body["total"] == 1
    assert body["items"][0]["change_type"] == "INBOUND"
    option = service.get_stock_history.call_args.args[1]
    assert option.size == 10 and option.warehouse_uid is not None


def test_get_stock_history_missing_stock_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_stock_history.side_effect = StockNotFoundError("stock not found")
    assert client.get(f"/api/v1/stocks/{uuid4()}/history").status_code == 404
