from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.service.delivery_service import CarrierService
from internal.infrastructure.adapter.restcontroller.carrier_controller import create_carrier_router
from tests.conftest import make_carrier


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=CarrierService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_carrier_router(service))
    return TestClient(app, headers={"x-user-id": "42", "x-user-roles": "MANAGE"})


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def test_create_carrier_returns_201(client: TestClient, service: MagicMock) -> None:
    carrier = make_carrier()
    service.create_carrier.return_value = carrier
    response = client.post(
        "/api/v1/carriers",
        json={"name": "ACME", "code": "ACME", "base_rate": 1000.0, "rate_per_kg": 500.0},
    )
    assert response.status_code == 201
    assert response.json()["uid"] == str(carrier.uid)


def test_get_carrier_returns_200(client: TestClient, service: MagicMock) -> None:
    carrier = make_carrier()
    service.get_carrier.return_value = carrier
    response = client.get(f"/api/v1/carriers/{carrier.uid}")
    assert response.status_code == 200


def test_get_carrier_returns_404(client: TestClient, service: MagicMock) -> None:
    service.get_carrier.return_value = None
    response = client.get(f"/api/v1/carriers/{uuid4()}")
    assert response.status_code == 404


def test_list_carriers_passes_filters(client: TestClient, service: MagicMock) -> None:
    service.list_carriers.return_value = ([make_carrier()], 1)
    response = client.get(
        "/api/v1/carriers",
        params={"name": "ACME", "status": "ACTIVE", "page": 0, "size": 5},
    )
    assert response.status_code == 200
    option = service.list_carriers.call_args[0][0]
    assert option.name == "ACME"
    assert option.status == CarrierStatus.ACTIVE
    assert option.size == 5


def test_update_carrier_returns_200(client: TestClient, service: MagicMock) -> None:
    carrier = make_carrier()
    service.update_carrier.return_value = carrier
    response = client.put(
        f"/api/v1/carriers/{carrier.uid}",
        json={"name": "Renamed"},
    )
    assert response.status_code == 200


def test_update_carrier_returns_404(client: TestClient, service: MagicMock) -> None:
    service.update_carrier.return_value = None
    response = client.put(f"/api/v1/carriers/{uuid4()}", json={"name": "X"})
    assert response.status_code == 404
