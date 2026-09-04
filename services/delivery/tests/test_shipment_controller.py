from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from sqlalchemy.exc import IntegrityError

from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.error import ConflictError
from internal.domain.service.delivery_service import ShipmentService
from internal.infrastructure.adapter.restcontroller.shipment_controller import create_shipment_router
from tests.conftest import make_shipment, make_tracking


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=ShipmentService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_shipment_router(service))
    return TestClient(app, headers={"x-user-id": "42", "x-user-roles": "MANAGE"})


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def test_create_shipment_valid_request_created(client: TestClient, service: MagicMock) -> None:
    shipment = make_shipment()
    service.create_shipment.return_value = shipment
    response = client.post(
        "/api/v1/shipments",
        json={
            "order_uid": str(uuid4()),
            "origin_address": "origin",
            "destination_address": "dest",
            "destination_city": "Seoul",
            "destination_state": "KR",
            "destination_country_code": "KR",
            "destination_postal_code": "00000",
        },
    )
    assert response.status_code == 201
    assert response.json()["uid"] == str(shipment.uid)
    service.create_shipment.assert_called_once()


def test_create_shipment_already_shipped_order_conflict(client: TestClient, service: MagicMock) -> None:
    order_uid = uuid4()
    service.create_shipment.side_effect = ConflictError(f"Shipment for order {order_uid} already exists")
    response = client.post(
        "/api/v1/shipments",
        json={
            "order_uid": str(order_uid),
            "origin_address": "origin",
            "destination_address": "dest",
            "destination_city": "Seoul",
            "destination_state": "KR",
            "destination_country_code": "KR",
            "destination_postal_code": "00000",
        },
    )
    assert response.status_code == 409
    assert str(order_uid) in response.json()["detail"]


def test_create_shipment_unique_violation_conflict(client: TestClient, service: MagicMock) -> None:
    service.create_shipment.side_effect = IntegrityError("INSERT", {}, Exception("uk_shipment_order_uid"))
    response = client.post(
        "/api/v1/shipments",
        json={
            "order_uid": str(uuid4()),
            "origin_address": "origin",
            "destination_address": "dest",
            "destination_city": "Seoul",
            "destination_state": "KR",
            "destination_country_code": "KR",
            "destination_postal_code": "00000",
        },
    )
    assert response.status_code == 409


def test_create_shipment_over_length_address_unprocessable(client: TestClient, service: MagicMock) -> None:
    response = client.post(
        "/api/v1/shipments",
        json={
            "order_uid": str(uuid4()),
            "origin_address": "x" * 501,
            "destination_address": "dest",
            "destination_city": "Seoul",
            "destination_state": "KR",
            "destination_country_code": "KR",
            "destination_postal_code": "00000",
        },
    )
    assert response.status_code == 422
    service.create_shipment.assert_not_called()


def test_assign_carrier_over_length_tracking_number_unprocessable(client: TestClient, service: MagicMock) -> None:
    response = client.patch(
        f"/api/v1/shipments/{uuid4()}/carrier",
        json={"carrier_uid": str(uuid4()), "tracking_number": "x" * 101},
    )
    assert response.status_code == 422
    service.assign_carrier.assert_not_called()


def test_create_shipment_blank_address_unprocessable(client: TestClient, service: MagicMock) -> None:
    response = client.post(
        "/api/v1/shipments",
        json={
            "order_uid": str(uuid4()),
            "origin_address": "",
            "destination_address": "dest",
            "destination_city": "Seoul",
            "destination_state": "",
            "destination_country_code": "KR",
            "destination_postal_code": "00000",
        },
    )
    assert response.status_code == 422
    service.create_shipment.assert_not_called()


def test_create_shipment_delivery_before_pickup_bad_request(client: TestClient, service: MagicMock) -> None:
    response = client.post(
        "/api/v1/shipments",
        json={
            "order_uid": str(uuid4()),
            "origin_address": "origin",
            "destination_address": "dest",
            "destination_city": "Seoul",
            "destination_state": "",
            "destination_country_code": "KR",
            "destination_postal_code": "00000",
            "planned_pickup_at": "2026-08-20T10:00:00Z",
            "planned_delivery_at": "2026-08-19T10:00:00Z",
        },
    )
    assert response.status_code == 400
    assert "planned_delivery_at" in response.json()["detail"]
    service.create_shipment.assert_not_called()


def test_get_shipment_found_ok(client: TestClient, service: MagicMock) -> None:
    shipment = make_shipment()
    service.get_shipment.return_value = shipment
    response = client.get(f"/api/v1/shipments/{shipment.uid}")
    assert response.status_code == 200
    assert response.json()["uid"] == str(shipment.uid)


def test_get_shipment_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_shipment.return_value = None
    response = client.get(f"/api/v1/shipments/{uuid4()}")
    assert response.status_code == 404


def test_list_shipments_with_filters_passes_them_to_service(client: TestClient, service: MagicMock) -> None:
    service.list_shipments.return_value = ([make_shipment()], 1)
    response = client.get(
        "/api/v1/shipments",
        params={"status": "PLANNED", "page": 0, "size": 10},
    )
    assert response.status_code == 200
    body = response.json()
    assert body["total"] == 1
    option = service.list_shipments.call_args[0][0]
    assert option.status == ShipmentStatus.PLANNED
    assert option.size == 10


def test_assign_carrier_found_ok(client: TestClient, service: MagicMock) -> None:
    shipment = make_shipment()
    service.assign_carrier.return_value = shipment
    response = client.patch(
        f"/api/v1/shipments/{shipment.uid}/carrier",
        json={"carrier_uid": str(uuid4()), "tracking_number": "TRK-1"},
    )
    assert response.status_code == 200


def test_assign_carrier_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.assign_carrier.return_value = None
    response = client.patch(
        f"/api/v1/shipments/{uuid4()}/carrier",
        json={"carrier_uid": str(uuid4())},
    )
    assert response.status_code == 404


def test_dispatch_invalid_status_bad_request(client: TestClient, service: MagicMock) -> None:
    service.dispatch_shipment.side_effect = ValueError("Cannot dispatch")
    response = client.patch(f"/api/v1/shipments/{uuid4()}/dispatch")
    assert response.status_code == 400
    assert "Cannot dispatch" in response.json()["detail"]


def test_dispatch_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.dispatch_shipment.return_value = None
    response = client.patch(f"/api/v1/shipments/{uuid4()}/dispatch")
    assert response.status_code == 404


def test_pickup_dispatched_shipment_ok(client: TestClient, service: MagicMock) -> None:
    shipment = make_shipment(status=ShipmentStatus.PICKED_UP)
    service.pick_up_shipment.return_value = shipment
    response = client.patch(f"/api/v1/shipments/{shipment.uid}/pickup")
    assert response.status_code == 200
    assert response.json()["status"] == "PICKED_UP"


def test_deliver_picked_up_shipment_ok(client: TestClient, service: MagicMock) -> None:
    shipment = make_shipment(status=ShipmentStatus.DELIVERED)
    service.deliver_shipment.return_value = shipment
    response = client.patch(f"/api/v1/shipments/{shipment.uid}/deliver")
    assert response.status_code == 200


def test_cancel_cancellable_shipment_ok(client: TestClient, service: MagicMock) -> None:
    shipment = make_shipment(status=ShipmentStatus.CANCELLED)
    service.cancel_shipment.return_value = shipment
    response = client.patch(f"/api/v1/shipments/{shipment.uid}/cancel")
    assert response.status_code == 200


def test_add_tracking_found_created(client: TestClient, service: MagicMock) -> None:
    shipment_uid = uuid4()
    tracking = make_tracking(shipment_uid=shipment_uid, status=ShipmentStatus.IN_TRANSIT)
    service.add_tracking.return_value = tracking
    response = client.post(
        f"/api/v1/shipments/{shipment_uid}/tracking",
        json={"status": "IN_TRANSIT", "location": "Hub", "description": "Arrived"},
    )
    assert response.status_code == 201
    assert response.json()["status"] == "IN_TRANSIT"


def test_add_tracking_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.add_tracking.side_effect = ValueError("Shipment not found")
    response = client.post(
        f"/api/v1/shipments/{uuid4()}/tracking",
        json={"status": "IN_TRANSIT", "location": "Hub", "description": "Arrived"},
    )
    assert response.status_code == 404


def test_get_tracking_history_existing_events_ok_with_items(client: TestClient, service: MagicMock) -> None:
    shipment_uid = uuid4()
    service.get_tracking_history.return_value = [
        make_tracking(shipment_uid=shipment_uid),
        make_tracking(shipment_uid=shipment_uid, status=ShipmentStatus.DELIVERED),
    ]
    response = client.get(f"/api/v1/shipments/{shipment_uid}/tracking")
    assert response.status_code == 200
    assert len(response.json()["items"]) == 2
