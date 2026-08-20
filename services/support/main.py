from collections.abc import AsyncIterator
from contextlib import asynccontextmanager

import structlog
import uvicorn
from fastapi import FastAPI

from internal.config.config import Settings
from internal.di.container import Container
from internal.infrastructure.adapter.restcontroller.category_controller import create_category_router
from internal.infrastructure.adapter.restcontroller.comment_controller import create_comment_router
from internal.infrastructure.adapter.restcontroller.faq_controller import create_faq_router
from internal.infrastructure.adapter.restcontroller.ticket_controller import create_ticket_router
from internal.infrastructure.support.logging import access_log_middleware, configure_logging
from internal.infrastructure.support.telemetry import instrument_fastapi, setup_telemetry

logger = structlog.get_logger()


def create_app() -> FastAPI:
    settings = Settings()

    configure_logging(
        level=settings.app.logging_level,
        debug=settings.server.mode == "debug",
    )
    setup_telemetry(settings.telemetry)

    container = Container(settings)

    @asynccontextmanager
    async def lifespan(_app: FastAPI) -> AsyncIterator[None]:
        try:
            yield
        finally:
            await container.engine.dispose()

    app = FastAPI(title="Support Service", version="1.0.0", lifespan=lifespan)
    app.include_router(create_ticket_router(container.support_service))
    app.include_router(create_comment_router(container.support_service))
    app.include_router(create_category_router(container.support_service))
    app.include_router(create_faq_router(container.support_service))

    @app.get("/health")
    def health():
        return {"status": "ok"}

    app.middleware("http")(access_log_middleware)
    instrument_fastapi(app)

    logger.info(
        "support service initialized",
        rest_port=settings.server.port,
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
