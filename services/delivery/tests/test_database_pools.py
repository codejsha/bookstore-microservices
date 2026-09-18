from internal.config.config import DatabaseConfig
from internal.infrastructure.support.database import create_engine_and_session_factory

MESH_REQUEST_TIMEOUT_SECONDS = 15
ENV_PREFIX = "DELIVERY"


async def test_application_pool_timeout_stays_below_mesh_deadline() -> None:
    engine, _ = create_engine_and_session_factory(DatabaseConfig(), ENV_PREFIX)
    try:
        assert engine.pool._timeout < MESH_REQUEST_TIMEOUT_SECONDS
    finally:
        await engine.dispose()
