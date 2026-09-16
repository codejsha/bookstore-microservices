import asyncio
from types import SimpleNamespace
from unittest.mock import AsyncMock

import pytest

from internal.config.config import TemporalConfig
from internal.infrastructure.support import temporal_worker as worker_module
from internal.infrastructure.support.temporal_auth import AccessToken, TemporalTokenProvider
from internal.infrastructure.support.temporal_worker import TemporalWorkerRunner


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
        FakeWorker.instances.append(self)

    async def run(self) -> None:
        await asyncio.Event().wait()


@pytest.fixture
def fake_temporal(monkeypatch: pytest.MonkeyPatch) -> AsyncMock:
    FakeWorker.instances = []
    client = SimpleNamespace(rpc_metadata={})
    connect = AsyncMock(return_value=client)
    monkeypatch.setattr(worker_module.Client, "connect", connect)
    monkeypatch.setattr(worker_module, "Worker", FakeWorker)
    return connect


def _runner(provider: TemporalTokenProvider | None) -> TemporalWorkerRunner:
    return TemporalWorkerRunner(TemporalConfig(), [], provider)


class TestTemporalWorkerRunnerAuth:
    async def test_start_auth_enabled_connects_with_bearer_metadata(self, fake_temporal: AsyncMock) -> None:
        runner = _runner(TemporalTokenProvider(ScriptedFetcher(AccessToken("first", 300))))

        await runner.start()
        while not FakeWorker.instances:
            await asyncio.sleep(0)
        await runner.stop()

        assert fake_temporal.await_args.kwargs["rpc_metadata"] == {"authorization": "Bearer first"}

    async def test_start_auth_disabled_connects_without_metadata(self, fake_temporal: AsyncMock) -> None:
        runner = _runner(None)

        await runner.start()
        while not FakeWorker.instances:
            await asyncio.sleep(0)
        await runner.stop()

        assert fake_temporal.await_args.kwargs["rpc_metadata"] == {}

    async def test_start_token_refreshed_updates_connected_client_metadata(self, fake_temporal: AsyncMock) -> None:
        runner = _runner(TemporalTokenProvider(ScriptedFetcher(AccessToken("first", 0.04), AccessToken("second", 300))))

        await runner.start()
        client = fake_temporal.return_value
        try:
            async with asyncio.timeout(2.0):
                while client.rpc_metadata != {"authorization": "Bearer second"}:
                    await asyncio.sleep(0.01)
        finally:
            await runner.stop()

        assert client.rpc_metadata == {"authorization": "Bearer second"}

    async def test_start_token_fetch_fails_raises_without_connecting(self, fake_temporal: AsyncMock) -> None:
        runner = _runner(TemporalTokenProvider(ScriptedFetcher(ConnectionError("keycloak down"))))

        with pytest.raises(ConnectionError):
            await runner.start()

        fake_temporal.assert_not_awaited()

    async def test_stop_auth_enabled_cancels_refresh_task(self, fake_temporal: AsyncMock) -> None:
        runner = _runner(TemporalTokenProvider(ScriptedFetcher(AccessToken("first", 300))))
        await runner.start()
        refresh_task = runner._refresh_task

        await runner.stop()

        assert refresh_task is not None and refresh_task.cancelled()
        assert runner._refresh_task is None
