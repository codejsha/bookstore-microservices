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
    run_errors: list[Exception] = []

    def __init__(self, client, **kwargs):
        self.client = client
        self.shutdown_awaited = False
        self._stopped = asyncio.Event()
        self._run_error = FakeWorker.run_errors.pop(0) if FakeWorker.run_errors else None
        FakeWorker.instances.append(self)

    @property
    def is_running(self) -> bool:
        return False

    @property
    def is_shutdown(self) -> bool:
        return False

    async def run(self) -> None:
        if self._run_error is not None:
            raise self._run_error
        await self._stopped.wait()

    async def shutdown(self) -> None:
        self.shutdown_awaited = True
        self._stopped.set()


@pytest.fixture
def fake_temporal(monkeypatch: pytest.MonkeyPatch) -> AsyncMock:
    FakeWorker.instances = []
    FakeWorker.run_errors = []
    monkeypatch.setattr(worker_module, "_INITIAL_BACKOFF_SECONDS", 0.0)
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

    async def test_run_token_refreshed_after_worker_retry_updates_retried_worker_client(
        self, fake_temporal: AsyncMock
    ) -> None:
        FakeWorker.run_errors.append(RuntimeError("Worker validation failed"))
        runner = _runner(TemporalTokenProvider(ScriptedFetcher(AccessToken("first", 0.2), AccessToken("second", 300))))
        client = fake_temporal.return_value

        await runner.start()
        try:
            async with asyncio.timeout(2.0):
                while len(FakeWorker.instances) < 2:
                    await asyncio.sleep(0.01)
            assert fake_temporal.await_args.kwargs["rpc_metadata"] == {"authorization": "Bearer first"}
            async with asyncio.timeout(2.0):
                while client.rpc_metadata != {"authorization": "Bearer second"}:
                    await asyncio.sleep(0.01)
            retried = FakeWorker.instances[-1]
        finally:
            await runner.stop()

        assert fake_temporal.await_count == 1
        assert FakeWorker.instances[0].shutdown_awaited
        assert retried.client is client
        assert retried.client.rpc_metadata == {"authorization": "Bearer second"}


async def test_run_token_refreshed_during_connect_connected_client_gets_latest_token(
    monkeypatch: pytest.MonkeyPatch,
) -> None:
    FakeWorker.instances = []
    FakeWorker.run_errors = []
    monkeypatch.setattr(worker_module, "Worker", FakeWorker)
    provider = TemporalTokenProvider(ScriptedFetcher(AccessToken("first", 0.05), AccessToken("second", 300)))
    client = SimpleNamespace(rpc_metadata={})

    async def slow_connect(host, namespace, rpc_metadata):
        client.rpc_metadata = dict(rpc_metadata)
        async with asyncio.timeout(2.0):
            while provider.rpc_metadata() != {"authorization": "Bearer second"}:
                await asyncio.sleep(0.01)
        return client

    monkeypatch.setattr(worker_module.Client, "connect", slow_connect)
    runner = _runner(provider)
    await runner.start()
    try:
        async with asyncio.timeout(2.0):
            while not FakeWorker.instances:
                await asyncio.sleep(0.01)
    finally:
        await runner.stop()

    assert FakeWorker.instances[0].client.rpc_metadata == {"authorization": "Bearer second"}
