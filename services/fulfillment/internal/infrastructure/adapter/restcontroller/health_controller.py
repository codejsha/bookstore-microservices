from collections.abc import Awaitable, Callable, Mapping

from fastapi import APIRouter
from fastapi.responses import JSONResponse

ComponentCheck = Callable[[], bool]


def create_health_router(
    ping: Callable[[], Awaitable[bool]],
    components: Mapping[str, ComponentCheck] | None = None,
) -> APIRouter:
    router = APIRouter(tags=["health"])
    checks: dict[str, ComponentCheck] = dict(components or {})

    @router.get("/health")
    async def health() -> dict[str, str]:
        return {"status": "ok"}

    @router.get("/health/ready")
    async def ready():
        unavailable = sorted(name for name, check in checks.items() if not check())
        if not unavailable and not await ping():
            unavailable = ["database"]
        if unavailable:
            return JSONResponse(
                status_code=503,
                content={"status": "unavailable", "components": unavailable},
            )
        return {"status": "ok"}

    return router
