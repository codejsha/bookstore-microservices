import asyncio
from datetime import timedelta

import pytest

from internal.config.config import TemporalConfig
from internal.infrastructure.support import temporal_worker as temporal_worker_module
from internal.infrastructure.support.temporal_worker import TemporalWorkerRunner


class _FailAfterStart(RuntimeError):
    pass


class _FakeWorker:
    next_validation_gate: asyncio.Event | None = None

    def __init__(self, client, *, task_queue, activities, **kwargs):
        self.client = client
        self.task_queue = task_queue
        self.activities = activities
        self.kwargs = kwargs
        self.shutdown_awaited = False
        self.run_cancelled = False
        self.drains = True
        self.run_error: Exception | None = None
        self.validation_gate = _FakeWorker.next_validation_gate
        self.started = False
        self.finished = False
        self._running = asyncio.Event()
        self._stopped = asyncio.Event()

    @property
    def is_running(self) -> bool:
        return self.started and not self.finished

    @property
    def is_shutdown(self) -> bool:
        return self.finished

    async def run(self) -> None:
        if self.validation_gate is not None:
            await self.validation_gate.wait()
        if self.run_error is not None and not isinstance(self.run_error, _FailAfterStart):
            raise self.run_error
        self.started = True
        if self.run_error is not None:
            self.finished = True
            raise self.run_error
        self._running.set()
        try:
            await self._stopped.wait()
        except asyncio.CancelledError:
            self.run_cancelled = True
            raise
        finally:
            self.finished = True

    async def shutdown(self) -> None:
        self.shutdown_awaited = True
        if self.drains:
            self._stopped.set()

    async def wait_until_running(self) -> None:
        await asyncio.wait_for(self._running.wait(), timeout=1.0)


class _FakeClient:
    connects: list[_FakeClient] = []
    connect_errors: list[Exception] = []

    @staticmethod
    async def connect(host, namespace="default", rpc_metadata=None):
        if _FakeClient.connect_errors:
            raise _FakeClient.connect_errors.pop(0)
        client = _FakeClient()
        _FakeClient.connects.append(client)
        return client


@pytest.fixture
def config() -> TemporalConfig:
    return TemporalConfig(enabled=True, host="localhost:7233", namespace="default", task_queue="fulfillment-task-queue")


@pytest.fixture
def run_errors() -> list[Exception]:
    return []


@pytest.fixture
def workers(monkeypatch: pytest.MonkeyPatch, run_errors: list[Exception]) -> list[_FakeWorker]:
    created: list[_FakeWorker] = []
    _FakeClient.connects = []
    _FakeClient.connect_errors = []
    monkeypatch.setattr(temporal_worker_module, "_INITIAL_BACKOFF_SECONDS", 0.0)

    def _factory(*args, **kwargs) -> _FakeWorker:
        worker = _FakeWorker(*args, **kwargs)
        if run_errors:
            worker.run_error = run_errors.pop(0)
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


async def _start_until_running(runner: TemporalWorkerRunner, workers: list[_FakeWorker], count: int) -> _FakeWorker:
    await runner.start()
    async with asyncio.timeout(1.0):
        while len(workers) < count:
            await asyncio.sleep(0.01)
    await workers[count - 1].wait_until_running()
    return workers[count - 1]


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


async def test_run_worker_failure_shuts_down_failed_worker(
    config: TemporalConfig, workers: list[_FakeWorker], run_errors: list[Exception]
) -> None:
    run_errors.append(RuntimeError("Namespace default is not found"))
    runner = TemporalWorkerRunner(config, [])
    await _start_until_running(runner, workers, 2)
    try:
        assert workers[0].shutdown_awaited
        assert not workers[1].shutdown_awaited
        assert runner.is_connected
    finally:
        await runner.stop()


async def test_run_worker_failure_reuses_client(
    config: TemporalConfig, workers: list[_FakeWorker], run_errors: list[Exception]
) -> None:
    run_errors.extend([RuntimeError("Worker failed"), RuntimeError("Worker validation failed")])
    runner = TemporalWorkerRunner(config, [])
    await _start_until_running(runner, workers, 3)
    try:
        assert len(_FakeClient.connects) == 1
        assert all(worker.client is _FakeClient.connects[0] for worker in workers)
        assert runner._client is _FakeClient.connects[0]
    finally:
        await runner.stop()


async def test_run_failed_worker_shutdown_hangs_retries_after_timeout(
    config: TemporalConfig, workers: list[_FakeWorker], run_errors: list[Exception], monkeypatch: pytest.MonkeyPatch
) -> None:
    monkeypatch.setattr(temporal_worker_module, "_FAILED_SHUTDOWN_TIMEOUT_SECONDS", 0.05)
    run_errors.append(RuntimeError("Worker validation failed"))
    shutdown_cancelled: list[bool] = []

    async def _never_completes(self: _FakeWorker) -> None:
        self.shutdown_awaited = True
        try:
            await asyncio.Event().wait()
        except asyncio.CancelledError:
            shutdown_cancelled.append(True)
            raise

    async def _completes(self: _FakeWorker) -> None:
        self.shutdown_awaited = True
        self._stopped.set()

    async def _shutdown(self: _FakeWorker) -> None:
        if self.run_error is not None:
            await _never_completes(self)
        else:
            await _completes(self)

    monkeypatch.setattr(_FakeWorker, "shutdown", _shutdown)
    runner = TemporalWorkerRunner(config, [])
    await _start_until_running(runner, workers, 2)
    try:
        assert workers[0].shutdown_awaited
        assert shutdown_cancelled == [True]
        assert len(_FakeClient.connects) == 1
    finally:
        await runner.stop()


async def test_run_connect_failure_retries_until_connected(config: TemporalConfig, workers: list[_FakeWorker]) -> None:
    _FakeClient.connect_errors.extend([ConnectionError("unavailable"), ConnectionError("unavailable")])
    runner = TemporalWorkerRunner(config, [])
    worker = await _start_until_running(runner, workers, 1)
    try:
        assert len(_FakeClient.connects) == 1
        assert worker.client is _FakeClient.connects[0]
    finally:
        await runner.stop()


async def test_is_connected_worker_not_yet_validated_false(
    config: TemporalConfig, workers: list[_FakeWorker], monkeypatch: pytest.MonkeyPatch
) -> None:
    gate = asyncio.Event()
    monkeypatch.setattr(_FakeWorker, "next_validation_gate", gate)
    runner = TemporalWorkerRunner(config, [])
    await runner.start()
    try:
        async with asyncio.timeout(1.0):
            while not workers:
                await asyncio.sleep(0.01)
        await asyncio.sleep(0.02)
        assert not runner.is_connected
        gate.set()
        await workers[0].wait_until_running()
        assert runner.is_connected
    finally:
        await runner.stop()


async def test_run_validation_keeps_failing_backoff_keeps_growing(
    config: TemporalConfig, workers: list[_FakeWorker], run_errors: list[Exception], monkeypatch: pytest.MonkeyPatch
) -> None:
    monkeypatch.setattr(temporal_worker_module, "_INITIAL_BACKOFF_SECONDS", 0.001)
    run_errors.extend([RuntimeError("Worker validation failed")] * 3)
    backoffs: list[float] = []
    runner = TemporalWorkerRunner(config, [])
    monkeypatch.setattr(runner, "_log_unavailable", lambda backoff, error: backoffs.append(backoff))
    await _start_until_running(runner, workers, 4)
    try:
        assert backoffs == [0.001, 0.002, 0.004]
    finally:
        await runner.stop()


async def test_run_started_worker_fails_backoff_resets(
    config: TemporalConfig, workers: list[_FakeWorker], run_errors: list[Exception], monkeypatch: pytest.MonkeyPatch
) -> None:
    monkeypatch.setattr(temporal_worker_module, "_INITIAL_BACKOFF_SECONDS", 0.001)
    run_errors.extend(
        [
            RuntimeError("Worker validation failed"),
            RuntimeError("Worker validation failed"),
            _FailAfterStart("Poll failed"),
        ]
    )
    backoffs: list[float] = []
    runner = TemporalWorkerRunner(config, [])
    monkeypatch.setattr(runner, "_log_unavailable", lambda backoff, error: backoffs.append(backoff))
    await _start_until_running(runner, workers, 4)
    try:
        assert backoffs == [0.001, 0.002, 0.001]
    finally:
        await runner.stop()
