import asyncio

from sqlalchemy import text

from internal.config.config import Settings, vault_env_prefix
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
from internal.infrastructure.adapter.event.noop_publisher import NoopEventPublisher
from internal.infrastructure.adapter.mysql.audit_repo import MySQLStockAuditRepository
from internal.infrastructure.adapter.mysql.balance_repo import MySQLStockBalanceRepository
from internal.infrastructure.adapter.mysql.carrier_repo import MySQLCarrierRepository
from internal.infrastructure.adapter.mysql.closing_repo import MySQLMonthlyClosingRepository
from internal.infrastructure.adapter.mysql.freight_repo import MySQLFreightRepository
from internal.infrastructure.adapter.mysql.shipment_repo import MySQLShipmentRepository
from internal.infrastructure.adapter.mysql.stats_repo import MySQLStatsRepository
from internal.infrastructure.adapter.mysql.stock_history_repo import MySQLStockHistoryRepository
from internal.infrastructure.adapter.mysql.stock_repo import MySQLStockRepository
from internal.infrastructure.adapter.mysql.stock_reservation_repo import MySQLStockReservationRepository
from internal.infrastructure.adapter.mysql.tracking_repo import MySQLTrackingRepository
from internal.infrastructure.adapter.mysql.transfer_repo import MySQLStockTransferRepository
from internal.infrastructure.adapter.mysql.warehouse_repo import MySQLWarehouseRepository
from internal.infrastructure.adapter.temporal.shipment_activities import ShipmentActivities
from internal.infrastructure.adapter.temporal.stock_activities import StockActivities
from internal.infrastructure.support.database import (
    create_engine_and_session_factory,
    create_readiness_engine,
)
from internal.infrastructure.support.temporal_auth import TemporalTokenProvider, build_temporal_token_provider

_PING_TIMEOUT = 3.0


class Container:
    def __init__(self, settings: Settings):
        self.settings = settings
        self.engine, self._session_factory = create_engine_and_session_factory(settings.database, vault_env_prefix())
        self._readiness_engine = create_readiness_engine(settings.database, vault_env_prefix())

        self.event_publisher = NoopEventPublisher()

        self.shipment_repo = MySQLShipmentRepository(self._session_factory)
        self.tracking_repo = MySQLTrackingRepository(self._session_factory)
        self.carrier_repo = MySQLCarrierRepository(self._session_factory)
        self.freight_repo = MySQLFreightRepository(self._session_factory)
        self.stats_repo = MySQLStatsRepository(self._session_factory)

        self.stock_repo = MySQLStockRepository(self._session_factory)
        self.stock_history_repo = MySQLStockHistoryRepository(self._session_factory)
        self.stock_reservation_repo = MySQLStockReservationRepository(self._session_factory)
        self.warehouse_repo = MySQLWarehouseRepository(self._session_factory)
        self.transfer_repo = MySQLStockTransferRepository(self._session_factory)
        self.audit_repo = MySQLStockAuditRepository(self._session_factory)
        self.closing_repo = MySQLMonthlyClosingRepository(self._session_factory)
        self.balance_repo = MySQLStockBalanceRepository(self._session_factory)

        self.shipment_service = ShipmentService(
            shipment_repo=self.shipment_repo,
            tracking_repo=self.tracking_repo,
        )
        self.carrier_service = CarrierService(carrier_repo=self.carrier_repo)
        self.freight_service = FreightService(freight_repo=self.freight_repo)
        self.stats_service = StatsService(stats_repo=self.stats_repo)

        self.stock_service = StockService(
            stock_repo=self.stock_repo,
            stock_history_repo=self.stock_history_repo,
            reservation_repo=self.stock_reservation_repo,
            event_publisher=self.event_publisher,
        )
        self.warehouse_service = WarehouseService(
            warehouse_repo=self.warehouse_repo, event_publisher=self.event_publisher
        )
        self.transfer_service = TransferService(transfer_repo=self.transfer_repo, event_publisher=self.event_publisher)
        self.audit_service = AuditService(audit_repo=self.audit_repo, event_publisher=self.event_publisher)
        self.closing_service = ClosingService(closing_repo=self.closing_repo, event_publisher=self.event_publisher)
        self.balance_service = BalanceService(balance_repo=self.balance_repo)

        self.shipment_activities = ShipmentActivities(shipment_service=self.shipment_service)
        self.stock_activities = StockActivities(stock_service=self.stock_service)

    async def ping_db(self) -> bool:
        async def _run() -> None:
            async with self._readiness_engine.connect() as conn:
                await conn.execute(text("SELECT 1"))

        try:
            await asyncio.wait_for(_run(), timeout=_PING_TIMEOUT)
        except Exception:
            return False
        return True

    async def close_readiness(self) -> None:
        await self._readiness_engine.dispose()

    async def dispose(self) -> None:
        await self.engine.dispose()
        await self.close_readiness()

    def temporal_token_provider(self) -> TemporalTokenProvider | None:
        return build_temporal_token_provider(self.settings.temporal.auth)

    def temporal_activities(self) -> list:
        return [
            self.shipment_activities.create_shipment,
            self.stock_activities.reserve_stock,
            self.stock_activities.release_stock,
        ]

    def register_grpc_servicers(self, server) -> tuple[str, ...]:
        from generated.application.port.pb.fulfillmentpb.fulfillment import v1_pb2 as pb
        from generated.application.port.pb.fulfillmentpb.fulfillment import v1_pb2_grpc as pb_grpc
        from internal.infrastructure.adapter.grpc.fulfillment_server import FulfillmentServiceServicer

        pb_grpc.add_FulfillmentServiceServicer_to_server(FulfillmentServiceServicer(self.shipment_service), server)
        return (pb.DESCRIPTOR.services_by_name["FulfillmentService"].full_name,)
