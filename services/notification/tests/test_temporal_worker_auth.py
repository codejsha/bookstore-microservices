import asyncio
from types import SimpleNamespace
from unittest.mock import AsyncMock, MagicMock

import pytest

from internal.config.config import TemporalConfig
from internal.infrastructure.adapter.temporal import worker as worker_module
from internal.infrastructure.adapter.temporal.temporal_auth import AccessToken, TemporalTokenProvider
from internal.infrastructure.adapter.temporal.worker import TemporalWorker


class ScriptedFetcher:
    def __init__(self, *steps: AccessToken | Exception):
        self._steps = list(steps)

    async def __call__(self) -> AccessToken:
        step = self._steps.pop(0) if len(self._steps) > 1 else self._steps[0]
        if isinstance(step, Exception):
            raise step
        return step


class FakeWorker:
    instances: list[FakeWorker] = []

    def __init__(self, client, **kwargs):
        self.client = client
        self.release = asyncio.Event()
        FakeWorker.instances.append(self)

    async def run(self) -> None:
        await self.release.wait()

    async def shutdown(self) -> None:
        self.release.set()


@pytest.fixture
def fake_temporal(monkeypatch: pytest.MonkeyPatch) -> AsyncMock:
    FakeWorker.instances = []
    client = SimpleNamespace(rpc_metadata={})
    connect = AsyncMock(return_value=client)
    monkeypatch.setattr(worker_module.Client, "connect", connect)
    monkeypatch.setattr(worker_module, "Worker", FakeWorker)
    return connect


def _worker(provider: TemporalTokenProvider | None) -> TemporalWorker:
    return TemporalWorker(TemporalConfig(), MagicMock(), provider)


class TestTemporalWorkerAuth:
    async def test_run_auth_enabled_connects_with_bearer_metadata(self, fake_temporal: AsyncMock) -> None:
        worker = _worker(TemporalTokenProvider(ScriptedFetcher(AccessToken("first", 300))))
        await worker.prepare()

        task = asyncio.create_task(worker.run())
        while not FakeWorker.instances:
            await asyncio.sleep(0)
        await worker.stop()
        await asyncio.wait_for(task, timeout=2.0)

        assert fake_temporal.await_args.kwargs["rpc_metadata"] == {"authorization": "Bearer first"}

    async def test_run_auth_disabled_connects_without_metadata(self, fake_temporal: AsyncMock) -> None:
        worker = _worker(None)
        await worker.prepare()

        task = asyncio.create_task(worker.run())
        while not FakeWorker.instances:
            await asyncio.sleep(0)
        await worker.stop()
        await asyncio.wait_for(task, timeout=2.0)

        assert fake_temporal.await_args.kwargs["rpc_metadata"] == {}

    async def test_run_token_refreshed_updates_connected_client_metadata(self, fake_temporal: AsyncMock) -> None:
        worker = _worker(TemporalTokenProvider(ScriptedFetcher(AccessToken("first", 0.04), AccessToken("second", 300))))
        await worker.prepare()

        task = asyncio.create_task(worker.run())
        client = fake_temporal.return_value
        async with asyncio.timeout(2.0):
            while client.rpc_metadata != {"authorization": "Bearer second"}:
                await asyncio.sleep(0.01)
        await worker.stop()
        await asyncio.wait_for(task, timeout=2.0)

        assert client.rpc_metadata == {"authorization": "Bearer second"}

    async def test_prepare_token_fetch_fails_raises(self, fake_temporal: AsyncMock) -> None:
        worker = _worker(TemporalTokenProvider(ScriptedFetcher(ConnectionError("keycloak down"))))

        with pytest.raises(ConnectionError):
            await worker.prepare()

        fake_temporal.assert_not_awaited()
