from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.di.container import Container


class _StubContainer:
    ping_db = Container.ping_db

    def __init__(self, session_factory) -> None:
        self._session_factory = session_factory


class _FailingSession:
    async def __aenter__(self) -> _FailingSession:
        return self

    async def __aexit__(self, *_args) -> bool:
        return False

    async def execute(self, *_args, **_kwargs):
        raise RuntimeError("connection refused")


async def test_ping_db_query_succeeds_returns_true(session_factory: async_sessionmaker[AsyncSession]) -> None:
    container = _StubContainer(session_factory)
    assert await container.ping_db() is True


async def test_ping_db_query_raises_returns_false() -> None:
    container = _StubContainer(lambda: _FailingSession())
    assert await container.ping_db() is False
