from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.service.inventory_service import WarehouseService
from internal.infrastructure.adapter.restcontroller.warehouse_controller import create_warehouse_router
from tests.conftest import make_warehouse

MANAGER_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "MANAGER,STAFF,USER"}


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=WarehouseService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_warehouse_router(service))
    return TestClient(app, headers=MANAGER_HEADERS)


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def test_register_warehouse_valid_request_created_with_location(client: TestClient, service: MagicMock) -> None:
    warehouse = make_warehouse()
    service.register_warehouse.return_value = warehouse

    response = client.post("/api/v1/warehouses", json={"name": "Main", "address": "1 Depot Rd", "capacity": 100})

    assert response.status_code == 201
    assert response.headers["location"] == f"/api/v1/warehouses/{warehouse.uid}"
    assert response.content == b""


def test_register_warehouse_blank_name_bad_request(client: TestClient, service: MagicMock) -> None:
    response = client.post("/api/v1/warehouses", json={"name": " ", "capacity": 100})
    assert response.status_code == 400
    service.register_warehouse.assert_not_called()


def test_register_warehouse_negative_capacity_unprocessable(client: TestClient, service: MagicMock) -> None:
    assert client.post("/api/v1/warehouses", json={"name": "Main", "capacity": -1}).status_code == 422
    service.register_warehouse.assert_not_called()


def test_get_warehouse_found_ok(client: TestClient, service: MagicMock) -> None:
    warehouse = make_warehouse()
    service.get_warehouse.return_value = warehouse

    response = client.get(f"/api/v1/warehouses/{warehouse.uid}")

    assert response.status_code == 200
    body = response.json()
    assert body["uid"] == str(warehouse.uid)
    assert body["updated_at"] is None


def test_get_warehouse_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_warehouse.return_value = None
    assert client.get(f"/api/v1/warehouses/{uuid4()}").status_code == 404


def test_list_warehouses_with_name_filter_passes_option(client: TestClient, service: MagicMock) -> None:
    service.list_warehouses.return_value = ([make_warehouse()], 1)

    response = client.get("/api/v1/warehouses", params={"name": "Main", "sort": "name,asc", "page": 1, "size": 3})

    assert response.status_code == 200
    option = service.list_warehouses.call_args.args[0]
    assert (option.name, option.sort, option.page, option.size) == ("Main", "name,asc", 1, 3)


def test_update_warehouse_found_ok_with_partial_body(client: TestClient, service: MagicMock) -> None:
    warehouse = make_warehouse(name="Renamed")
    service.update_warehouse.return_value = warehouse

    response = client.put(f"/api/v1/warehouses/{warehouse.uid}", json={"name": "Renamed"})

    assert response.status_code == 200
    assert response.json()["name"] == "Renamed"
    command = service.update_warehouse.call_args.args[1]
    assert command.name == "Renamed" and command.capacity is None


def test_update_warehouse_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.update_warehouse.return_value = None
    assert client.put(f"/api/v1/warehouses/{uuid4()}", json={"name": "x"}).status_code == 404
