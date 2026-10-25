import shutil
from collections.abc import AsyncIterator
from pathlib import Path

import pytest_asyncio
from sqlalchemy.ext.asyncio import AsyncEngine, create_async_engine

from internal.di.container import Container


class _StubContainer:
    ping_db = Container.ping_db
    close_readiness = Container.close_readiness
    dispose = Container.dispose

    def __init__(self, readiness_engine: AsyncEngine, engine: object | None = None) -> None:
        self._readiness_engine = readiness_engine
        self.engine = engine


class _CountingEngine:
    def __init__(self) -> None:
        self.dispose_calls = 0

    async def dispose(self) -> None:
        self.dispose_calls += 1


@pytest_asyncio.fixture
async def readiness_engine(tmp_path: Path) -> AsyncIterator[AsyncEngine]:
    db_dir = tmp_path / "readiness"
    db_dir.mkdir()
    engine = create_async_engine(f"sqlite+aiosqlite:///{db_dir / 'readiness.db'}")
    try:
        yield engine
    finally:
        await engine.dispose()


async def test_ping_db_readiness_engine_reachable_true(readiness_engine: AsyncEngine) -> None:
    container = _StubContainer(readiness_engine)
    assert await container.ping_db() is True


async def test_ping_db_readiness_engine_unreachable_false(tmp_path: Path) -> None:
    engine = create_async_engine(f"sqlite+aiosqlite:///{tmp_path / 'absent' / 'readiness.db'}")
    container = _StubContainer(engine)
    try:
        assert await container.ping_db() is False
    finally:
        await engine.dispose()


async def test_ping_db_readiness_engine_disposed_false(readiness_engine: AsyncEngine, tmp_path: Path) -> None:
    container = _StubContainer(readiness_engine)
    assert await container.ping_db() is True

    await container.close_readiness()
    shutil.rmtree(tmp_path / "readiness")

    assert await container.ping_db() is False


async def test_dispose_application_and_readiness_engines_disposed_once() -> None:
    application = _CountingEngine()
    readiness = _CountingEngine()
    container = _StubContainer(readiness, application)

    await container.dispose()

    assert application.dispose_calls == 1
    assert readiness.dispose_calls == 1
