from unittest.mock import MagicMock
from uuid import uuid4

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.aggregate.stats_aggregate import DeliveryDashboard
from internal.domain.service.delivery_service import (
    CarrierService,
    FreightService,
    ShipmentService,
    StatsService,
)
from internal.domain.service.inventory_service import (
    AuditService,
    BalanceService,
    ClosingService,
    StockService,
    TransferService,
    WarehouseService,
)
from internal.infrastructure.adapter.restcontroller.audit_controller import create_audit_router
from internal.infrastructure.adapter.restcontroller.balance_controller import create_balance_router
from internal.infrastructure.adapter.restcontroller.carrier_controller import create_carrier_router
from internal.infrastructure.adapter.restcontroller.closing_controller import create_closing_router
from internal.infrastructure.adapter.restcontroller.freight_controller import create_freight_router
from internal.infrastructure.adapter.restcontroller.shipment_controller import create_shipment_router
from internal.infrastructure.adapter.restcontroller.stats_controller import create_stats_router
from internal.infrastructure.adapter.restcontroller.stock_controller import create_stock_router
from internal.infrastructure.adapter.restcontroller.transfer_controller import create_transfer_router
from internal.infrastructure.adapter.restcontroller.warehouse_controller import create_warehouse_router
from tests.conftest import make_audit, make_closing, make_stock, make_transfer, make_warehouse

CUSTOMER_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "USER"}
STAFF_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "STAFF,USER"}
MANAGER_HEADERS = {"x-user-id": str(uuid4()), "x-user-roles": "MANAGER,USER"}
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


def _stock_body() -> dict:
    return {"edition_uid": str(uuid4()), "warehouse_uid": str(uuid4()), "quantity": 1, "reason": "x"}


def _stock_app() -> FastAPI:
    service = MagicMock(spec=StockService)
    service.list_stocks.return_value = ([], 0)
    service.get_stock.return_value = make_stock()
    service.receive_stock.return_value = make_stock()
    service.release_stock.return_value = make_stock()
    service.adjust_stock.return_value = make_stock()
    service.reserve_stock.return_value = make_stock()
    return _app(create_stock_router(service))


def test_stocks_anonymous_caller_unauthorized() -> None:
    assert TestClient(_stock_app()).get("/api/v1/stocks").status_code == 401


def test_stocks_read_customer_caller_ok() -> None:
    client = TestClient(_stock_app(), headers=CUSTOMER_HEADERS)
    assert client.get("/api/v1/stocks").status_code == 200
    assert client.get(f"/api/v1/stocks/{uuid4()}").status_code == 200


@pytest.mark.parametrize("path", ["receive", "release", "adjust", "reserve"])
def test_stocks_write_customer_caller_forbidden(path: str) -> None:
    client = TestClient(_stock_app(), headers=CUSTOMER_HEADERS)
    assert client.post(f"/api/v1/stocks/{path}", json=_stock_body()).status_code == 403


@pytest.mark.parametrize("path", ["receive", "release", "adjust", "reserve"])
def test_stocks_write_staff_caller_ok(path: str) -> None:
    client = TestClient(_stock_app(), headers=STAFF_HEADERS)
    assert client.post(f"/api/v1/stocks/{path}", json=_stock_body()).status_code == 200


def test_stocks_write_system_caller_ok() -> None:
    client = TestClient(_stock_app(), headers=SYSTEM_HEADERS)
    assert client.post("/api/v1/stocks/receive", json=_stock_body()).status_code == 200


def _warehouse_app() -> FastAPI:
    service = MagicMock(spec=WarehouseService)
    service.list_warehouses.return_value = ([], 0)
    service.register_warehouse.return_value = make_warehouse()
    service.update_warehouse.return_value = make_warehouse()
    return _app(create_warehouse_router(service))


def test_warehouses_read_customer_caller_ok() -> None:
    assert TestClient(_warehouse_app(), headers=CUSTOMER_HEADERS).get("/api/v1/warehouses").status_code == 200


def test_warehouses_create_staff_caller_forbidden() -> None:
    client = TestClient(_warehouse_app(), headers=STAFF_HEADERS)
    assert client.post("/api/v1/warehouses", json={"name": "Main", "capacity": 1}).status_code == 403


def test_warehouses_update_staff_caller_forbidden() -> None:
    client = TestClient(_warehouse_app(), headers=STAFF_HEADERS)
    assert client.put(f"/api/v1/warehouses/{uuid4()}", json={"name": "Main"}).status_code == 403


def test_warehouses_create_manager_caller_created() -> None:
    client = TestClient(_warehouse_app(), headers=MANAGER_HEADERS)
    assert client.post("/api/v1/warehouses", json={"name": "Main", "capacity": 1}).status_code == 201


def test_warehouses_update_system_caller_ok() -> None:
    client = TestClient(_warehouse_app(), headers=SYSTEM_HEADERS)
    assert client.put(f"/api/v1/warehouses/{uuid4()}", json={"name": "Main"}).status_code == 200


def _transfer_app() -> FastAPI:
    service = MagicMock(spec=TransferService)
    service.list_transfers.return_value = ([], 0)
    service.transfer_stock.return_value = make_transfer()
    service.complete_transfer.return_value = make_transfer()
    service.cancel_transfer.return_value = make_transfer()
    return _app(create_transfer_router(service))


def _transfer_body() -> dict:
    return {
        "edition_uid": str(uuid4()),
        "source_warehouse_uid": str(uuid4()),
        "target_warehouse_uid": str(uuid4()),
        "quantity": 1,
    }


def test_transfers_read_customer_caller_ok() -> None:
    assert TestClient(_transfer_app(), headers=CUSTOMER_HEADERS).get("/api/v1/transfers").status_code == 200


def test_transfers_create_customer_caller_forbidden() -> None:
    client = TestClient(_transfer_app(), headers=CUSTOMER_HEADERS)
    assert client.post("/api/v1/transfers", json=_transfer_body()).status_code == 403


@pytest.mark.parametrize("action", ["complete", "cancel"])
def test_transfers_finalize_customer_caller_forbidden(action: str) -> None:
    client = TestClient(_transfer_app(), headers=CUSTOMER_HEADERS)
    assert client.post(f"/api/v1/transfers/{uuid4()}/{action}").status_code == 403


@pytest.mark.parametrize("action", ["complete", "cancel"])
def test_transfers_finalize_staff_caller_ok(action: str) -> None:
    client = TestClient(_transfer_app(), headers=STAFF_HEADERS)
    assert client.post(f"/api/v1/transfers/{uuid4()}/{action}").status_code == 200


def _audit_app() -> FastAPI:
    service = MagicMock(spec=AuditService)
    service.list_audits.return_value = ([], 0)
    service.create_audit.return_value = make_audit()
    service.complete_audit.return_value = make_audit()
    return _app(create_audit_router(service))


def _audit_body() -> dict:
    return {"warehouse_uid": str(uuid4()), "items": [{"edition_uid": str(uuid4()), "actual_quantity": 1}]}


def test_audits_create_customer_caller_forbidden() -> None:
    client = TestClient(_audit_app(), headers=CUSTOMER_HEADERS)
    assert client.post("/api/v1/audits", json=_audit_body()).status_code == 403


def test_audits_complete_customer_caller_forbidden() -> None:
    client = TestClient(_audit_app(), headers=CUSTOMER_HEADERS)
    assert client.post(f"/api/v1/audits/{uuid4()}/complete").status_code == 403


def test_audits_create_staff_caller_ok() -> None:
    client = TestClient(_audit_app(), headers=STAFF_HEADERS)
    assert client.post("/api/v1/audits", json=_audit_body()).status_code == 200


def test_audits_read_customer_caller_ok() -> None:
    assert TestClient(_audit_app(), headers=CUSTOMER_HEADERS).get("/api/v1/audits").status_code == 200


def _closing_app() -> FastAPI:
    service = MagicMock(spec=ClosingService)
    service.list_closings.return_value = ([], 0)
    service.create_monthly_closing.return_value = make_closing()
    return _app(create_closing_router(service))


def _closing_body() -> dict:
    return {"warehouse_uid": str(uuid4()), "year": 2026, "month": 10}


def test_closings_create_staff_caller_forbidden() -> None:
    client = TestClient(_closing_app(), headers=STAFF_HEADERS)
    assert client.post("/api/v1/closings", json=_closing_body()).status_code == 403


def test_closings_create_manager_caller_ok() -> None:
    client = TestClient(_closing_app(), headers=MANAGER_HEADERS)
    assert client.post("/api/v1/closings", json=_closing_body()).status_code == 200


def test_closings_create_system_caller_ok() -> None:
    client = TestClient(_closing_app(), headers=SYSTEM_HEADERS)
    assert client.post("/api/v1/closings", json=_closing_body()).status_code == 200


def test_closings_read_customer_caller_ok() -> None:
    assert TestClient(_closing_app(), headers=CUSTOMER_HEADERS).get("/api/v1/closings").status_code == 200


def test_balance_anonymous_caller_unauthorized() -> None:
    app = _app(create_balance_router(MagicMock(spec=BalanceService)))
    assert TestClient(app).get("/api/v1/balance").status_code == 401


def test_balance_customer_caller_ok() -> None:
    service = MagicMock(spec=BalanceService)
    service.get_stock_balance.return_value = ([], 0)
    app = _app(create_balance_router(service))
    assert TestClient(app, headers=CUSTOMER_HEADERS).get("/api/v1/balance").status_code == 200
