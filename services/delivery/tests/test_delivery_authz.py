from unittest.mock import MagicMock
from uuid import uuid4

from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.aggregate.stats_aggregate import DeliveryDashboard
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

CUSTOMER_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "USER"}
STAFF_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "STAFF,USER"}
SYSTEM_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "SYSTEM"}


def _app(router) -> FastAPI:
    app = FastAPI()
    app.include_router(router)
    return app


def _empty_dashboard() -> DeliveryDashboard:
    return DeliveryDashboard(
        total_shipments=0,
        status_counts=[],
        carrier_performances=[],
        avg_delivery_hours=None,
        total_freight_cost=0.0,
        on_time_delivery_rate=None,
    )


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


def test_carriers_staff_caller_passes_authz() -> None:
    service = MagicMock(spec=CarrierService)
    service.get_carrier.return_value = None
    app = _app(create_carrier_router(service))
    resp = TestClient(app, headers=STAFF_HEADERS).get(f"/api/v1/carriers/{uuid4()}")
    assert resp.status_code == 404


def test_stats_staff_caller_forbidden() -> None:
    app = _app(create_stats_router(MagicMock(spec=StatsService)))
    assert TestClient(app, headers=STAFF_HEADERS).get("/api/v1/stats/dashboard").status_code == 403


def test_carriers_system_caller_passes_authz() -> None:
    service = MagicMock(spec=CarrierService)
    service.get_carrier.return_value = None
    app = _app(create_carrier_router(service))
    resp = TestClient(app, headers=SYSTEM_HEADERS).get(f"/api/v1/carriers/{uuid4()}")
    assert resp.status_code == 404


def test_stats_system_caller_ok() -> None:
    service = MagicMock(spec=StatsService)
    service.get_dashboard.return_value = _empty_dashboard()
    app = _app(create_stats_router(service))
    assert TestClient(app, headers=SYSTEM_HEADERS).get("/api/v1/stats/dashboard").status_code == 200
