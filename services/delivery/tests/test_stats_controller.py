from datetime import datetime
from unittest.mock import MagicMock

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.domain.aggregate.stats_aggregate import (
    CarrierPerformance,
    DeliveryDashboard,
    ShipmentStatusCount,
)
from internal.domain.service.delivery_service import StatsService
from internal.infrastructure.adapter.restcontroller.stats_controller import create_stats_router


@pytest.fixture(scope="module")
def service() -> MagicMock:
    return MagicMock(spec=StatsService)


@pytest.fixture(scope="module")
def client(service: MagicMock) -> TestClient:
    app = FastAPI()
    app.include_router(create_stats_router(service))
    return TestClient(app, headers={"x-user-id": "42", "x-user-roles": "MANAGE"})


@pytest.fixture(autouse=True)
def _reset(service: MagicMock) -> None:
    service.reset_mock(return_value=True, side_effect=True)


def test_dashboard_ok_with_service_metrics(client: TestClient, service: MagicMock) -> None:
    service.get_dashboard.return_value = DeliveryDashboard(
        total_shipments=5,
        status_counts=[ShipmentStatusCount(status="PLANNED", count=3)],
        carrier_performances=[
            CarrierPerformance(
                carrier_uid="00000000-0000-0000-0000-000000000001",
                carrier_name="ACME",
                total_shipments=5,
                delivered_count=3,
                failed_count=1,
                avg_delivery_hours=24.0,
                total_freight_cost=10000.0,
            )
        ],
        avg_delivery_hours=24.0,
        total_freight_cost=10000.0,
        on_time_delivery_rate=80.0,
    )
    response = client.get("/api/v1/stats/dashboard")
    assert response.status_code == 200
    body = response.json()
    assert body["total_shipments"] == 5
    assert body["status_counts"][0]["status"] == "PLANNED"
    assert body["carrier_performances"][0]["carrier_name"] == "ACME"
    assert body["on_time_delivery_rate"] == 80.0


def test_dashboard_with_filters_passes_them_to_service(client: TestClient, service: MagicMock) -> None:
    service.get_dashboard.return_value = DeliveryDashboard(
        total_shipments=0,
        status_counts=[],
        carrier_performances=[],
        avg_delivery_hours=None,
        total_freight_cost=0.0,
        on_time_delivery_rate=None,
    )
    response = client.get(
        "/api/v1/stats/dashboard",
        params={"date_from": "2026-01-01", "date_to": "2026-12-31"},
    )
    assert response.status_code == 200
    option = service.get_dashboard.call_args[0][0]
    assert option.date_from == datetime(2026, 1, 1)
    assert option.date_to == datetime(2026, 12, 31)
