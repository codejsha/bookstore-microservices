from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.constant.audit_status import AuditStatus
from internal.domain.model.error import InvalidTransitionError, StockNotFoundError, WarehouseNotFoundError
from internal.domain.service.inventory_service import AuditService
from internal.infrastructure.adapter.restcontroller.audit_controller import create_audit_router
from tests.conftest import make_audit

STAFF_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "STAFF,USER"}


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=AuditService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_audit_router(service))
    return TestClient(app, headers=STAFF_HEADERS)


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def _body(**overrides) -> dict:
    body = {"warehouse_uid": str(uuid4()), "items": [{"edition_uid": str(uuid4()), "actual_quantity": 3}]}
    body.update(overrides)
    return body


def test_create_audit_valid_request_ok_with_items(client: TestClient, service: MagicMock) -> None:
    audit = make_audit()
    service.create_audit.return_value = audit

    response = client.post("/api/v1/audits", json=_body(notes="count"))

    assert response.status_code == 200
    body = response.json()
    assert body["status"] == "IN_PROGRESS"
    assert body["items"][0]["difference"] == -2
    assert service.create_audit.call_args.args[0].notes == "count"


def test_create_audit_duplicate_editions_bad_request(client: TestClient, service: MagicMock) -> None:
    edition = str(uuid4())
    items = [{"edition_uid": edition, "actual_quantity": 1}, {"edition_uid": edition, "actual_quantity": 2}]
    assert client.post("/api/v1/audits", json=_body(items=items)).status_code == 400
    service.create_audit.assert_not_called()


def test_create_audit_empty_items_bad_request(client: TestClient, service: MagicMock) -> None:
    assert client.post("/api/v1/audits", json=_body(items=[])).status_code == 400


def test_create_audit_missing_warehouse_not_found(client: TestClient, service: MagicMock) -> None:
    service.create_audit.side_effect = WarehouseNotFoundError("missing")
    assert client.post("/api/v1/audits", json=_body()).status_code == 404


def test_create_audit_missing_stock_not_found(client: TestClient, service: MagicMock) -> None:
    service.create_audit.side_effect = StockNotFoundError("missing")
    assert client.post("/api/v1/audits", json=_body()).status_code == 404


def test_get_audit_missing_not_found(client: TestClient, service: MagicMock) -> None:
    service.get_audit.return_value = None
    assert client.get(f"/api/v1/audits/{uuid4()}").status_code == 404


def test_list_audits_with_status_filter_passes_option(client: TestClient, service: MagicMock) -> None:
    service.list_audits.return_value = ([make_audit()], 1)
    response = client.get("/api/v1/audits", params={"status": "COMPLETED"})
    assert response.status_code == 200
    assert service.list_audits.call_args.args[0].status == AuditStatus.COMPLETED


def test_complete_audit_already_completed_conflict(client: TestClient, service: MagicMock) -> None:
    service.complete_audit.side_effect = InvalidTransitionError("already completed")
    assert client.post(f"/api/v1/audits/{uuid4()}/complete").status_code == 409


def test_complete_audit_found_ok(client: TestClient, service: MagicMock) -> None:
    service.complete_audit.return_value = make_audit(status=AuditStatus.COMPLETED)
    response = client.post(f"/api/v1/audits/{uuid4()}/complete")
    assert response.status_code == 200
    assert response.json()["status"] == "COMPLETED"
