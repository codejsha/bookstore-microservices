from collections.abc import AsyncIterator
from contextlib import asynccontextmanager

import structlog
import uvicorn
from fastapi import FastAPI

from internal.config.config import Settings
from internal.di.container import Container
from internal.infrastructure.adapter.restcontroller.carrier_controller import create_carrier_router
from internal.infrastructure.adapter.restcontroller.freight_controller import create_freight_router
from internal.infrastructure.adapter.restcontroller.shipment_controller import create_shipment_router
from internal.infrastructure.adapter.restcontroller.stats_controller import create_stats_router
from internal.infrastructure.support.grpc_server import GrpcServerRunner
from internal.infrastructure.support.logging import access_log_middleware, configure_logging
from internal.infrastructure.support.telemetry import instrument_fastapi, setup_telemetry
from internal.infrastructure.support.temporal_worker import TemporalWorkerRunner

logger = structlog.get_logger()


def _build_grpc_runner(settings: Settings, container: Container) -> GrpcServerRunner | None:
    if not settings.grpc.enabled:
        logger.info("gRPC server disabled (grpc.enabled=false) — REST-only mode")
        return None
    try:
        from generated.application.port.pb.deliverypb.delivery import v1_pb2_grpc  # noqa: F401
    except ImportError as e:
        logger.warning("gRPC stubs not available — run `make grpc`", error=str(e))
        return None
    return GrpcServerRunner(settings.grpc, container.register_grpc_servicers)


def create_app() -> FastAPI:
    settings = Settings()

    configure_logging(
        level=settings.app.logging_level,
        debug=settings.server.mode == "debug",
    )
    setup_telemetry(settings.telemetry)

    container = Container(settings)

    temporal_runner = TemporalWorkerRunner(settings.temporal, container.temporal_activities())
    grpc_runner = _build_grpc_runner(settings, container)

    @asynccontextmanager
    async def lifespan(_app: FastAPI) -> AsyncIterator[None]:
        await temporal_runner.start()
        if grpc_runner is not None:
            try:
                await grpc_runner.start()
            except Exception as e:  # noqa: BLE001
                logger.warning("gRPC server failed to start — continuing without it", error=str(e))
        try:
            yield
        finally:
            if grpc_runner is not None:
                await grpc_runner.stop()
            await temporal_runner.stop()

    app = FastAPI(
        title="Delivery Service",
        version="1.0.0",
        lifespan=lifespan,
    )

    app.include_router(create_shipment_router(container.shipment_service))
    app.include_router(create_carrier_router(container.carrier_service))
    app.include_router(create_freight_router(container.freight_service))
    app.include_router(create_stats_router(container.stats_service))

    @app.get("/health")
    def health():
        return {"status": "ok"}

    app.middleware("http")(access_log_middleware)
    instrument_fastapi(app)

    logger.info(
        "delivery service initialized",
        rest_port=settings.server.port,
        grpc_port=settings.grpc.port,
        temporal_task_queue=settings.temporal.task_queue,
    )
    return app


app = create_app()

if __name__ == "__main__":
    settings = Settings()
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=settings.server.port,
        reload=settings.server.mode == "debug",
    )
