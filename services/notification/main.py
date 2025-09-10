import asyncio
import contextlib
from contextlib import asynccontextmanager

import structlog
import uvicorn
from fastapi import FastAPI

from internal.config.config import Settings
from internal.di.container import Container
from internal.infrastructure.adapter.restcontroller.notification_controller import create_notification_router
from internal.infrastructure.adapter.restcontroller.template_controller import create_template_router
from internal.infrastructure.support.logging import access_log_middleware, configure_logging
from internal.infrastructure.support.telemetry import instrument_fastapi, setup_telemetry

_WORKER_SHUTDOWN_TIMEOUT = 10.0


def create_app() -> FastAPI:
    settings = Settings()

    configure_logging(
        level=settings.app.logging_level,
        debug=settings.server.mode == "debug",
    )
    setup_telemetry(settings.telemetry)

    container = Container(settings)

    @asynccontextmanager
    async def lifespan(_app: FastAPI):
        worker_task = asyncio.create_task(container.temporal_worker.run())
        try:
            yield
        finally:
            await container.temporal_worker.stop()
            try:
                await asyncio.wait_for(worker_task, timeout=_WORKER_SHUTDOWN_TIMEOUT)
            except asyncio.CancelledError:
                pass
            except TimeoutError:
                worker_task.cancel()
                with contextlib.suppress(asyncio.CancelledError):
                    await worker_task

    app = FastAPI(
        title="Notification Service",
        version="1.0.0",
        lifespan=lifespan,
    )

    app.include_router(create_notification_router(container.notification_service))
    app.include_router(create_template_router(container.notification_service))

    @app.get("/health")
    def health():
        return {"status": "ok"}

    app.middleware("http")(access_log_middleware)
    instrument_fastapi(app)

    structlog.get_logger().info("notification service started", port=settings.server.port)
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
