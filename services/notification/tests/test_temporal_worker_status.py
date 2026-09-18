import asyncio
from collections.abc import Callable
from unittest.mock import MagicMock

import pytest

from internal.config.config import TemporalConfig
from internal.infrastructure.adapter.temporal import worker as worker_module
from internal.infrastructure.adapter.temporal.activities import NotificationActivities
from internal.infrastructure.adapter.temporal.worker import TemporalWorker

_WAIT_TIMEOUT_SECONDS = 5.0


async def _wait_until(predicate: Callable[[], bool]) -> None:
    async with asyncio.timeout(_WAIT_TIMEOUT_SECONDS):
        while not predicate():
            await asyncio.sleep(0.01)


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

    async def shutdown(self) -> None:
        self._stopped.set()


@pytest.fixture
def activities() -> MagicMock:
    return MagicMock(spec=NotificationActivities)


async def test_is_connected_temporal_disabled_false(activities: MagicMock) -> None:
    worker = TemporalWorker(TemporalConfig(enabled=False), activities)

    await worker.run()

    assert worker.is_connected is False


async def test_is_connected_connect_succeeds_true(activities: MagicMock, monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setattr(worker_module, "Client", _FakeClient)
    monkeypatch.setattr(worker_module, "Worker", _FakeWorker)
    worker = TemporalWorker(TemporalConfig(enabled=True), activities)
    task = asyncio.create_task(worker.run())

    try:
        await _wait_until(lambda: worker.is_connected)
    finally:
        await worker.stop()
        await asyncio.wait_for(task, timeout=_WAIT_TIMEOUT_SECONDS)


async def test_is_connected_connect_fails_false(activities: MagicMock, monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setattr(worker_module, "Client", _UnavailableClient)
    worker = TemporalWorker(TemporalConfig(enabled=True), activities)
    task = asyncio.create_task(worker.run())

    try:
        await asyncio.sleep(0.05)
        assert worker.is_connected is False
    finally:
        await worker.stop()
        await asyncio.wait_for(task, timeout=_WAIT_TIMEOUT_SECONDS)


async def test_is_connected_after_stop_false(activities: MagicMock, monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setattr(worker_module, "Client", _FakeClient)
    monkeypatch.setattr(worker_module, "Worker", _FakeWorker)
    worker = TemporalWorker(TemporalConfig(enabled=True), activities)
    task = asyncio.create_task(worker.run())
    await _wait_until(lambda: worker.is_connected)

    await worker.stop()
    await asyncio.wait_for(task, timeout=_WAIT_TIMEOUT_SECONDS)

    assert worker.is_connected is False
