import asyncio
from collections.abc import Callable

import pytest

from internal.config.config import GrpcConfig, TemporalConfig
from internal.infrastructure.support import temporal_worker as temporal_worker_module
from internal.infrastructure.support.grpc_server import GrpcServerRunner
from internal.infrastructure.support.temporal_worker import TemporalWorkerRunner

_WAIT_TIMEOUT_SECONDS = 5.0


async def _wait_until(predicate: Callable[[], bool]) -> None:
    async with asyncio.timeout(_WAIT_TIMEOUT_SECONDS):
        while not predicate():
            await asyncio.sleep(0.01)


def _no_servicers(_server) -> tuple[str, ...]:
    return ()


class _FakeClient:
    @staticmethod
    async def connect(*_args, **_kwargs) -> object:
        return object()


class _UnavailableClient:
    @staticmethod
    async def connect(*_args, **_kwargs) -> object:
        raise RuntimeError("temporal unavailable")


class _FakeWorker:
    def __init__(self, *_args, **_kwargs) -> None:
        self._stopped = asyncio.Event()

    async def run(self) -> None:
        await self._stopped.wait()


class TestGrpcServerRunner:
    async def test_is_running_before_start_false(self) -> None:
        runner = GrpcServerRunner(GrpcConfig(enabled=True, port=0), _no_servicers)

        assert runner.is_running is False

    async def test_is_running_after_start_true(self) -> None:
        runner = GrpcServerRunner(GrpcConfig(enabled=True, port=0), _no_servicers)

        await runner.start()
        try:
            assert runner.is_running is True
        finally:
            await runner.stop()

    async def test_is_running_after_stop_false(self) -> None:
        runner = GrpcServerRunner(GrpcConfig(enabled=True, port=0), _no_servicers)

        await runner.start()
        await runner.stop()

        assert runner.is_running is False

    async def test_is_running_start_fails_false(self) -> None:
        def _failing_registration(_server) -> tuple[str, ...]:
            raise RuntimeError("servicer registration failed")

        runner = GrpcServerRunner(GrpcConfig(enabled=True, port=0), _failing_registration)

        with pytest.raises(RuntimeError):
            await runner.start()

        assert runner.is_running is False


class TestTemporalWorkerRunner:
    async def test_is_connected_temporal_disabled_false(self) -> None:
        runner = TemporalWorkerRunner(TemporalConfig(enabled=False), [])

        await runner.start()
        try:
            assert runner.is_connected is False
        finally:
            await runner.stop()

    async def test_is_connected_connect_succeeds_true(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setattr(temporal_worker_module, "Client", _FakeClient)
        monkeypatch.setattr(temporal_worker_module, "Worker", _FakeWorker)
        runner = TemporalWorkerRunner(TemporalConfig(enabled=True), [])

        await runner.start()
        try:
            await _wait_until(lambda: runner.is_connected)
        finally:
            await runner.stop()

    async def test_is_connected_connect_fails_false(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setattr(temporal_worker_module, "Client", _UnavailableClient)
        runner = TemporalWorkerRunner(TemporalConfig(enabled=True), [])

        await runner.start()
        try:
            await asyncio.sleep(0.05)
            assert runner.is_connected is False
        finally:
            await runner.stop()

    async def test_is_connected_after_stop_false(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setattr(temporal_worker_module, "Client", _FakeClient)
        monkeypatch.setattr(temporal_worker_module, "Worker", _FakeWorker)
        runner = TemporalWorkerRunner(TemporalConfig(enabled=True), [])

        await runner.start()
        await _wait_until(lambda: runner.is_connected)
        await runner.stop()

        assert runner.is_connected is False
