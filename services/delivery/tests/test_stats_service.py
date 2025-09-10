from unittest.mock import MagicMock

from internal.domain.aggregate.stats_aggregate import DeliveryDashboard
from internal.domain.model.option.delivery_option import StatsFilterOption
from internal.domain.service.delivery_service import StatsService


async def test_get_dashboard_delegates_to_repository(stats_repo: MagicMock) -> None:
    dashboard = DeliveryDashboard(
        total_shipments=10,
        status_counts=[],
        carrier_performances=[],
        avg_delivery_hours=24.0,
        total_freight_cost=10000.0,
        on_time_delivery_rate=0.95,
    )
    stats_repo.get_dashboard.return_value = dashboard
    service = StatsService(stats_repo=stats_repo)
    option = StatsFilterOption()
    assert await service.get_dashboard(option) is dashboard
    stats_repo.get_dashboard.assert_called_once_with(option)
