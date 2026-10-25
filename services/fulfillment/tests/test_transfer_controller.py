from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.constant.transfer_status import TransferStatus
from internal.domain.model.error import InvalidTransitionError, StockNotFoundError
from internal.domain.service.inventory_service import TransferService
from internal.infrastructure.adapter.restcontroller.transfer_controller import create_transfer_router
from tests.conftest import make_transfer

STAFF_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "STAFF,USER"}


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=TransferService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_transfer_router(service))
    return TestClient(app, headers=STAFF_HEADERS)


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def _body(**overrides) -> dict:
    body = {
        "edition_uid": str(uuid4()),
        "source_warehouse_uid": str(uuid4()),
        "target_warehouse_uid": str(uuid4()),
        "quantity": 3,
    }
    body.update(overrides)
    return body


def test_create_transfer_valid_request_ok_with_body(client: TestClient, service: MagicMock) -> None:
    transfer = make_transfer()
    service.transfer_stock.return_value = transfer

    response = client.post("/api/v1/transfers", json=_body())

    assert response.status_code == 200
    assert response.json()["status"] == "PENDING"
    assert response.json()["uid"] == str(transfer.uid)


def test_create_transfer_same_warehouse_bad_request(client: TestClient, service: MagicMock) -> None:
    warehouse = str(uuid4())
    response = client.post(
        "/api/v1/transfers", json=_body(source_warehouse_uid=warehouse, target_warehouse_uid=warehouse)
    )
    assert response.status_code == 400
    service.transfer_stock.assert_not_called()


def test_create_transfer_missing_source_stock_not_found(client: TestClient, service: MagicMock) -> None:
    service.transfer_stock.side_effect = StockNotFoundError("no stock")
    assert client.post("/api/v1/transfers", json=_body()).status_code == 404


def test_get_transfer_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_transfer.return_value = None
    assert client.get(f"/api/v1/transfers/{uuid4()}").status_code == 404


def test_list_transfers_with_status_filter_passes_option(client: TestClient, service: MagicMock) -> None:
    service.list_transfers.return_value = ([make_transfer()], 1)

    response = client.get("/api/v1/transfers", params={"status": "PENDING"})

    assert response.status_code == 200
    assert service.list_transfers.call_args.args[0].status == TransferStatus.PENDING


def test_complete_transfer_invalid_transition_conflict(client: TestClient, service: MagicMock) -> None:
    service.complete_transfer.side_effect = InvalidTransitionError("cannot complete")
    assert client.post(f"/api/v1/transfers/{uuid4()}/complete").status_code == 409


def test_complete_transfer_found_ok(client: TestClient, service: MagicMock) -> None:
    service.complete_transfer.return_value = make_transfer(status=TransferStatus.COMPLETED)
    response = client.post(f"/api/v1/transfers/{uuid4()}/complete")
    assert response.status_code == 200
    assert response.json()["status"] == "COMPLETED"


def test_cancel_transfer_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.cancel_transfer.return_value = None
    assert client.post(f"/api/v1/transfers/{uuid4()}/cancel").status_code == 404
