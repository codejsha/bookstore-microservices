from datetime import datetime
from uuid import UUID

from fastapi import APIRouter, Depends, Query

from generated.application.port.model.carrier_performance import CarrierPerformance
from generated.application.port.model.dashboard_response import DashboardResponse
from generated.application.port.model.shipment_status_count import ShipmentStatusCount
from internal.domain.model.option.delivery_option import StatsFilterOption
from internal.domain.service.delivery_service import StatsService
from internal.infrastructure.support.auth import require_staff


def create_stats_router(service: StatsService) -> APIRouter:
    router = APIRouter(
        prefix="/api/v1/stats",
        tags=["stats"],
        dependencies=[Depends(require_staff)],
    )

    @router.get("/dashboard", response_model=DashboardResponse)
    async def get_dashboard(
        carrier_uid: UUID | None = Query(None),
        date_from: datetime | None = Query(None),
        date_to: datetime | None = Query(None),
    ) -> DashboardResponse:
        option = StatsFilterOption(
            carrier_uid=carrier_uid,
            date_from=date_from,
            date_to=date_to,
        )
        dashboard = await service.get_dashboard(option)
        return DashboardResponse(
            total_shipments=dashboard.total_shipments,
            status_counts=[ShipmentStatusCount(status=sc.status, count=sc.count) for sc in dashboard.status_counts],
            carrier_performances=[
                CarrierPerformance(
                    carrier_uid=str(cp.carrier_uid),
                    carrier_name=cp.carrier_name,
                    total_shipments=cp.total_shipments,
                    delivered_count=cp.delivered_count,
                    failed_count=cp.failed_count,
                    avg_delivery_hours=cp.avg_delivery_hours,
                    total_freight_cost=cp.total_freight_cost,
                )
                for cp in dashboard.carrier_performances
            ],
            avg_delivery_hours=dashboard.avg_delivery_hours,
            total_freight_cost=dashboard.total_freight_cost,
            on_time_delivery_rate=dashboard.on_time_delivery_rate,
        )

    return router
