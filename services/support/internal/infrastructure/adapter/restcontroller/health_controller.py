from collections.abc import Awaitable, Callable

from fastapi import APIRouter
from fastapi.responses import JSONResponse


def create_health_router(ping: Callable[[], Awaitable[bool]]) -> APIRouter:
    router = APIRouter(tags=["health"])

    @router.get("/health")
    async def health() -> dict[str, str]:
        return {"status": "ok"}

    @router.get("/health/ready")
    async def ready():
        if await ping():
            return {"status": "ok"}
        return JSONResponse(status_code=503, content={"status": "unavailable"})

    return router
