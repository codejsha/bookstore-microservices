from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.constant.freight_status import FreightStatus
from internal.domain.service.delivery_service import FreightService
from internal.infrastructure.adapter.restcontroller.freight_controller import create_freight_router
from tests.conftest import make_freight


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=FreightService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_freight_router(service))
    return TestClient(app, headers={"x-user-id": "42", "x-user-roles": "MANAGE,STAFF,USER"})


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def test_create_freight_valid_request_created(client: TestClient, service: MagicMock) -> None:
    freight = make_freight()
    service.create_freight.return_value = freight
    response = client.post(
        "/api/v1/freights",
        json={
            "shipment_uid": str(uuid4()),
            "carrier_uid": str(uuid4()),
            "base_cost": 1000.0,
        },
    )
    assert response.status_code == 201
    assert response.json()["uid"] == str(freight.uid)


def test_get_freight_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_freight.return_value = None
    response = client.get(f"/api/v1/freights/{uuid4()}")
    assert response.status_code == 404


def test_list_freights_with_filters_passes_them_to_service(client: TestClient, service: MagicMock) -> None:
    service.list_freights.return_value = ([make_freight()], 1)
    response = client.get("/api/v1/freights", params={"status": "PAID"})
    assert response.status_code == 200
    option = service.list_freights.call_args[0][0]
    assert option.status == FreightStatus.PAID


def test_get_freight_by_shipment_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.find_by_shipment.return_value = None
    response = client.get(f"/api/v1/freights/shipment/{uuid4()}")
    assert response.status_code == 404


def test_get_freight_by_shipment_found_ok(client: TestClient, service: MagicMock) -> None:
    freight = make_freight()
    service.find_by_shipment.return_value = freight
    response = client.get(f"/api/v1/freights/shipment/{freight.shipment_uid}")
    assert response.status_code == 200


def test_update_freight_status_found_ok(client: TestClient, service: MagicMock) -> None:
    freight = make_freight(status=FreightStatus.INVOICED)
    service.update_status.return_value = freight
    response = client.patch(f"/api/v1/freights/{freight.uid}/status", json={"status": "INVOICED"})
    assert response.status_code == 200
    assert response.json()["status"] == "INVOICED"


def test_update_freight_status_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.update_status.return_value = None
    response = client.patch(f"/api/v1/freights/{uuid4()}/status", json={"status": "PAID"})
    assert response.status_code == 404
