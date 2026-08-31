from unittest.mock import MagicMock
from uuid import uuid4

from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.service.delivery_service import (
    CarrierService,
    FreightService,
    ShipmentService,
    StatsService,
)
from internal.infrastructure.adapter.restcontroller.carrier_controller import create_carrier_router
from internal.infrastructure.adapter.restcontroller.freight_controller import create_freight_router
from internal.infrastructure.adapter.restcontroller.shipment_controller import create_shipment_router
from internal.infrastructure.adapter.restcontroller.stats_controller import create_stats_router

CUSTOMER_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "PROFILE,ORDER,VIEW"}


def _app(router) -> FastAPI:
    app = FastAPI()
    app.include_router(router)
    return app


def test_carriers_anonymous_caller_unauthorized() -> None:
    app = _app(create_carrier_router(MagicMock(spec=CarrierService)))
    assert TestClient(app).get("/api/v1/carriers").status_code == 401


def test_carriers_customer_caller_forbidden() -> None:
    app = _app(create_carrier_router(MagicMock(spec=CarrierService)))
    assert TestClient(app, headers=CUSTOMER_HEADERS).get("/api/v1/carriers").status_code == 403


def test_stats_customer_caller_forbidden() -> None:
    app = _app(create_stats_router(MagicMock(spec=StatsService)))
    assert TestClient(app, headers=CUSTOMER_HEADERS).get("/api/v1/stats/dashboard").status_code == 403


def test_freights_customer_caller_forbidden() -> None:
    app = _app(create_freight_router(MagicMock(spec=FreightService)))
    assert TestClient(app, headers=CUSTOMER_HEADERS).get("/api/v1/freights").status_code == 403


def test_shipment_list_customer_caller_forbidden() -> None:
    app = _app(create_shipment_router(MagicMock(spec=ShipmentService)))
    assert TestClient(app, headers=CUSTOMER_HEADERS).get("/api/v1/shipments").status_code == 403


def test_shipment_dispatch_customer_caller_forbidden() -> None:
    app = _app(create_shipment_router(MagicMock(spec=ShipmentService)))
    resp = TestClient(app, headers=CUSTOMER_HEADERS).patch(f"/api/v1/shipments/{uuid4()}/dispatch")
    assert resp.status_code == 403


def test_shipment_read_customer_caller_passes_authz() -> None:
    service = MagicMock(spec=ShipmentService)
    service.get_shipment.return_value = None
    app = _app(create_shipment_router(service))
    resp = TestClient(app, headers=CUSTOMER_HEADERS).get(f"/api/v1/shipments/{uuid4()}")
    assert resp.status_code == 404
