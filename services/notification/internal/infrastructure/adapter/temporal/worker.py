import asyncio
import contextlib
from collections.abc import Mapping
from datetime import timedelta

import structlog
from temporalio.client import Client
from temporalio.worker import Worker

from internal.config.config import TemporalConfig
from internal.infrastructure.adapter.temporal.activities import NotificationActivities
from internal.infrastructure.adapter.temporal.temporal_auth import TemporalTokenProvider

_logger = structlog.get_logger()

_INITIAL_BACKOFF_SECONDS = 1.0
_MAX_BACKOFF_SECONDS = 60.0
_GRACEFUL_SHUTDOWN_SECONDS = 5.0


class TemporalWorker:
    def __init__(
        self,
        config: TemporalConfig,
        activities: NotificationActivities,
        token_provider: TemporalTokenProvider | None = None,
    ):
        self._config = config
        self._activities = activities
        self._token_provider = token_provider
        self._client: Client | None = None
        self._worker: Worker | None = None
        self._stopped = asyncio.Event()

    async def prepare(self) -> None:
        if not self._config.enabled or self._token_provider is None:
            return
        await self._token_provider.initialize()
        _logger.info("temporal access token acquired")

    async def run(self) -> None:
        if not self._config.enabled:
            _logger.info("temporal worker disabled, skipping startup")
            return

        refresh_task: asyncio.Task | None = None
        if self._token_provider is not None:
            refresh_task = asyncio.create_task(
                self._token_provider.run(self._apply_rpc_metadata), name="temporal-token-refresh"
            )
        try:
            await self._run_worker()
        finally:
            if refresh_task is not None:
                refresh_task.cancel()
                with contextlib.suppress(asyncio.CancelledError):
                    await refresh_task

    async def _run_worker(self) -> None:
        backoff = _INITIAL_BACKOFF_SECONDS
        while not self._stopped.is_set():
            try:
                client = await Client.connect(
                    self._config.host,
                    namespace=self._config.namespace,
                    rpc_metadata=self._rpc_metadata(),
                )
                self._client = client
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
                self._client = None

    async def stop(self) -> None:
        self._stopped.set()
        worker = self._worker
        if worker is not None:
            await worker.shutdown()

    def _rpc_metadata(self) -> Mapping[str, str]:
        if self._token_provider is None:
            return {}
        return self._token_provider.rpc_metadata()

    def _apply_rpc_metadata(self, metadata: Mapping[str, str]) -> None:
        client = self._client
        if client is not None:
            client.rpc_metadata = metadata
