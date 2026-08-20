import asyncio
from collections.abc import Sequence

import structlog
from temporalio.client import Client
from temporalio.worker import Worker

from internal.config.config import TemporalConfig

logger = structlog.get_logger()

_INITIAL_BACKOFF_SECONDS = 1.0
_MAX_BACKOFF_SECONDS = 60.0


class TemporalWorkerRunner:
    def __init__(self, config: TemporalConfig, activities: Sequence):
        self._config = config
        self._activities = list(activities)
        self._task: asyncio.Task | None = None
        self._worker: Worker | None = None

    async def start(self) -> None:
        if not self._config.enabled:
            logger.info("temporal disabled — skipping worker startup")
            return
        self._task = asyncio.create_task(self._run(), name="temporal-worker")

    async def _run(self) -> None:
        backoff = _INITIAL_BACKOFF_SECONDS
        while True:
            try:
                client = await Client.connect(self._config.host, namespace=self._config.namespace)
                self._worker = Worker(
                    client,
                    task_queue=self._config.task_queue,
                    activities=self._activities,
                )
                logger.info(
                    "temporal worker started",
                    host=self._config.host,
                    namespace=self._config.namespace,
                    task_queue=self._config.task_queue,
                )
                backoff = _INITIAL_BACKOFF_SECONDS
                await self._worker.run()
                return
            except asyncio.CancelledError:
                raise
            except Exception as e:  # noqa: BLE001
                logger.error(
                    "temporal worker unavailable — retrying",
                    host=self._config.host,
                    namespace=self._config.namespace,
                    task_queue=self._config.task_queue,
                    retry_in_seconds=backoff,
                    error=str(e),
                )
                self._worker = None
                await asyncio.sleep(backoff)
                backoff = min(backoff * 2, _MAX_BACKOFF_SECONDS)

    async def stop(self) -> None:
        if self._task is None:
            return
        logger.info("temporal worker stopping")
        self._task.cancel()
        try:
            await self._task
        except asyncio.CancelledError:
            pass
        self._task = None
        self._worker = None
