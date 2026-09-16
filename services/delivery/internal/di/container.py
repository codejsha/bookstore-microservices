import asyncio

from sqlalchemy import text

from internal.config.config import Settings, vault_env_prefix
from internal.domain.service.delivery_service import (
    CarrierService,
    FreightService,
    ShipmentService,
    StatsService,
)
from internal.infrastructure.adapter.mysql.carrier_repo import MySQLCarrierRepository
from internal.infrastructure.adapter.mysql.freight_repo import MySQLFreightRepository
from internal.infrastructure.adapter.mysql.shipment_repo import MySQLShipmentRepository
from internal.infrastructure.adapter.mysql.stats_repo import MySQLStatsRepository
from internal.infrastructure.adapter.mysql.tracking_repo import MySQLTrackingRepository
from internal.infrastructure.adapter.temporal.shipment_activities import ShipmentActivities
from internal.infrastructure.support.database import create_session_factory
from internal.infrastructure.support.temporal_auth import TemporalTokenProvider, build_temporal_token_provider

_PING_TIMEOUT = 3.0


class Container:
    def __init__(self, settings: Settings):
        self.settings = settings
        self._session_factory = create_session_factory(settings.database, vault_env_prefix())

        self.shipment_repo = MySQLShipmentRepository(self._session_factory)
        self.tracking_repo = MySQLTrackingRepository(self._session_factory)
        self.carrier_repo = MySQLCarrierRepository(self._session_factory)
        self.freight_repo = MySQLFreightRepository(self._session_factory)
        self.stats_repo = MySQLStatsRepository(self._session_factory)

        self.shipment_service = ShipmentService(
            shipment_repo=self.shipment_repo,
            tracking_repo=self.tracking_repo,
        )
        self.carrier_service = CarrierService(carrier_repo=self.carrier_repo)
        self.freight_service = FreightService(freight_repo=self.freight_repo)
        self.stats_service = StatsService(stats_repo=self.stats_repo)

        self.shipment_activities = ShipmentActivities(shipment_service=self.shipment_service)

    async def ping_db(self) -> bool:
        async def _run() -> None:
            async with self._session_factory() as session:
                await session.execute(text("SELECT 1"))

        try:
            await asyncio.wait_for(_run(), timeout=_PING_TIMEOUT)
        except Exception:
            return False
        return True

    def temporal_token_provider(self) -> TemporalTokenProvider | None:
        return build_temporal_token_provider(self.settings.temporal.auth)

    def temporal_activities(self) -> list:
        return [self.shipment_activities.create_shipment]

    def register_grpc_servicers(self, server) -> tuple[str, ...]:
        from generated.application.port.pb.deliverypb.delivery import v1_pb2 as pb
        from generated.application.port.pb.deliverypb.delivery import v1_pb2_grpc as pb_grpc
        from internal.infrastructure.adapter.grpc.delivery_server import DeliveryServiceServicer

        pb_grpc.add_DeliveryServiceServicer_to_server(DeliveryServiceServicer(self.shipment_service), server)
        return (pb.DESCRIPTOR.services_by_name["DeliveryService"].full_name,)
