import asyncio
from datetime import timedelta

import structlog
from temporalio.client import Client
from temporalio.worker import Worker

from internal.config.config import TemporalConfig
from internal.infrastructure.adapter.temporal.activities import NotificationActivities

_logger = structlog.get_logger()

_INITIAL_BACKOFF_SECONDS = 1.0
_MAX_BACKOFF_SECONDS = 60.0
_GRACEFUL_SHUTDOWN_SECONDS = 5.0


class TemporalWorker:
    def __init__(self, config: TemporalConfig, activities: NotificationActivities):
        self._config = config
        self._activities = activities
        self._worker: Worker | None = None
        self._stopped = asyncio.Event()

    async def run(self) -> None:
        if not self._config.enabled:
            _logger.info("temporal worker disabled, skipping startup")
            return

        backoff = _INITIAL_BACKOFF_SECONDS
        while not self._stopped.is_set():
            try:
                client = await Client.connect(self._config.host, namespace=self._config.namespace)
                self._worker = Worker(
                    client,
                    task_queue=self._config.task_queue,
                    activities=[self._activities.send_shipment_update],
                    graceful_shutdown_timeout=timedelta(seconds=_GRACEFUL_SHUTDOWN_SECONDS),
                )
                _logger.info(
                    "temporal worker started",
                    host=self._config.host,
                    namespace=self._config.namespace,
                    task_queue=self._config.task_queue,
                )
                backoff = _INITIAL_BACKOFF_SECONDS
                await self._worker.run()
                return
            except asyncio.CancelledError:
                _logger.info("temporal worker shutting down")
                raise
            except Exception as exc:  # noqa: BLE001
                _logger.error(
                    "temporal worker unavailable — retrying",
                    host=self._config.host,
                    namespace=self._config.namespace,
                    task_queue=self._config.task_queue,
                    retry_in_seconds=backoff,
                    error=str(exc),
                )
                try:
                    await asyncio.wait_for(self._stopped.wait(), timeout=backoff)
                except TimeoutError:
                    pass
                backoff = min(backoff * 2, _MAX_BACKOFF_SECONDS)
            finally:
                self._worker = None

    async def stop(self) -> None:
        self._stopped.set()
        worker = self._worker
        if worker is not None:
            await worker.shutdown()
