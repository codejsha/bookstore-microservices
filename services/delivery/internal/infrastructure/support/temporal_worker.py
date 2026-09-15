import asyncio
import contextlib
from collections.abc import Mapping, Sequence

import structlog
from temporalio.client import Client
from temporalio.worker import Worker

from internal.config.config import TemporalConfig
from internal.infrastructure.support.temporal_auth import TemporalTokenProvider

logger = structlog.get_logger()

_INITIAL_BACKOFF_SECONDS = 1.0
_MAX_BACKOFF_SECONDS = 60.0


class TemporalWorkerRunner:
    def __init__(
        self,
        config: TemporalConfig,
        activities: Sequence,
        token_provider: TemporalTokenProvider | None = None,
    ):
        self._config = config
        self._activities = list(activities)
        self._token_provider = token_provider
        self._task: asyncio.Task | None = None
        self._refresh_task: asyncio.Task | None = None
        self._client: Client | None = None
        self._worker: Worker | None = None

    async def start(self) -> None:
        if not self._config.enabled:
            logger.info("temporal disabled — skipping worker startup")
            return
        if self._token_provider is not None:
            await self._token_provider.initialize()
            logger.info("temporal access token acquired")
            self._refresh_task = asyncio.create_task(
                self._token_provider.run(self._apply_rpc_metadata), name="temporal-token-refresh"
            )
        self._task = asyncio.create_task(self._run(), name="temporal-worker")

    async def _run(self) -> None:
        backoff = _INITIAL_BACKOFF_SECONDS
        while True:
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
                self._client = None
                await asyncio.sleep(backoff)
                backoff = min(backoff * 2, _MAX_BACKOFF_SECONDS)

    async def stop(self) -> None:
        if self._task is not None:
            logger.info("temporal worker stopping")
            self._task.cancel()
            try:
                await self._task
            except asyncio.CancelledError:
                pass
            self._task = None
        if self._refresh_task is not None:
            self._refresh_task.cancel()
            with contextlib.suppress(asyncio.CancelledError):
                await self._refresh_task
            self._refresh_task = None
        self._worker = None
        self._client = None

    def _rpc_metadata(self) -> Mapping[str, str]:
        if self._token_provider is None:
            return {}
        return self._token_provider.rpc_metadata()

    def _apply_rpc_metadata(self, metadata: Mapping[str, str]) -> None:
        client = self._client
        if client is not None:
            client.rpc_metadata = metadata
