import asyncio
from datetime import timedelta

import pytest

from internal.config.config import TemporalConfig
from internal.infrastructure.support import temporal_worker as temporal_worker_module
from internal.infrastructure.support.temporal_worker import TemporalWorkerRunner


class _FakeWorker:
    def __init__(self, client, *, task_queue, activities, **kwargs):
        self.client = client
        self.task_queue = task_queue
        self.activities = activities
        self.kwargs = kwargs
        self.shutdown_awaited = False
        self.run_cancelled = False
        self.drains = True
        self._running = asyncio.Event()
        self._stopped = asyncio.Event()

    async def run(self) -> None:
        self._running.set()
        try:
            await self._stopped.wait()
        except asyncio.CancelledError:
            self.run_cancelled = True
            raise

    async def shutdown(self) -> None:
        self.shutdown_awaited = True
        if self.drains:
            self._stopped.set()

    async def wait_until_running(self) -> None:
        await asyncio.wait_for(self._running.wait(), timeout=1.0)


class _FakeClient:
    @staticmethod
    async def connect(host, namespace="default", rpc_metadata=None):
        return _FakeClient()


@pytest.fixture
def config() -> TemporalConfig:
    return TemporalConfig(enabled=True, host="localhost:7233", namespace="default", task_queue="delivery-task-queue")


@pytest.fixture
def workers(monkeypatch: pytest.MonkeyPatch) -> list[_FakeWorker]:
    created: list[_FakeWorker] = []

    def _factory(*args, **kwargs) -> _FakeWorker:
        worker = _FakeWorker(*args, **kwargs)
        created.append(worker)
        return worker

    monkeypatch.setattr(temporal_worker_module, "Worker", _factory)
    monkeypatch.setattr(temporal_worker_module, "Client", _FakeClient)
    return created


async def _start(runner: TemporalWorkerRunner, workers: list[_FakeWorker]) -> _FakeWorker:
    await runner.start()
    for _ in range(100):
        if workers:
            break
        await asyncio.sleep(0.01)
    assert workers, "worker was never constructed"
    await workers[0].wait_until_running()
    return workers[0]


async def test_worker_is_built_with_graceful_shutdown_timeout(
    config: TemporalConfig, workers: list[_FakeWorker]
) -> None:
    runner = TemporalWorkerRunner(config, [])
    worker = await _start(runner, workers)
    try:
        assert worker.kwargs["graceful_shutdown_timeout"] == timedelta(
            seconds=temporal_worker_module._GRACEFUL_SHUTDOWN_SECONDS
        )
        assert worker.kwargs["graceful_shutdown_timeout"] > timedelta(0)
    finally:
        await runner.stop()


async def test_stop_drains_worker_without_cancelling_activities(
    config: TemporalConfig, workers: list[_FakeWorker]
) -> None:
    runner = TemporalWorkerRunner(config, [])
    worker = await _start(runner, workers)

    await runner.stop()

    assert worker.shutdown_awaited
    assert not worker.run_cancelled


async def test_stop_cancels_worker_that_does_not_drain_within_timeout(
    config: TemporalConfig, workers: list[_FakeWorker], monkeypatch: pytest.MonkeyPatch
) -> None:
    monkeypatch.setattr(temporal_worker_module, "_DRAIN_TIMEOUT_SECONDS", 0.05)
    runner = TemporalWorkerRunner(config, [])
    worker = await _start(runner, workers)
    worker.drains = False

    await runner.stop()

    assert worker.shutdown_awaited
    assert worker.run_cancelled


async def test_stop_before_start_is_noop(config: TemporalConfig, workers: list[_FakeWorker]) -> None:
    runner = TemporalWorkerRunner(config, [])
    await runner.stop()
    assert workers == []


async def test_start_disabled_does_not_build_worker(workers: list[_FakeWorker]) -> None:
    runner = TemporalWorkerRunner(TemporalConfig(enabled=False), [])
    await runner.start()
    await asyncio.sleep(0.05)
    await runner.stop()
    assert workers == []
