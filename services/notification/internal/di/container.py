import asyncio

from sqlalchemy import text

from internal.application.port.repo.repos import NotificationRepository, TemplateRepository
from internal.config.config import Settings, vault_env_prefix
from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.mysql.notification_repo import MySQLNotificationRepository
from internal.infrastructure.adapter.mysql.template_repo import MySQLTemplateRepository
from internal.infrastructure.adapter.temporal.activities import NotificationActivities
from internal.infrastructure.adapter.temporal.temporal_auth import build_temporal_token_provider
from internal.infrastructure.adapter.temporal.worker import TemporalWorker
from internal.infrastructure.support.database import create_session_factory

_PING_TIMEOUT = 3.0


class Container:
    def __init__(self, settings: Settings):
        self.settings = settings
        self._session_factory = create_session_factory(settings.database, vault_env_prefix())

        self.notification_repo: NotificationRepository = MySQLNotificationRepository(self._session_factory)
        self.template_repo: TemplateRepository = MySQLTemplateRepository(self._session_factory)

        self.notification_service = NotificationService(
            notification_repo=self.notification_repo,
            template_repo=self.template_repo,
        )

        self.notification_activities = NotificationActivities(self.notification_service)
        self.temporal_worker = TemporalWorker(
            settings.temporal,
            self.notification_activities,
            build_temporal_token_provider(settings.temporal.auth),
        )

    async def ping_db(self) -> bool:
        async def _run() -> None:
            async with self._session_factory() as session:
                await session.execute(text("SELECT 1"))

        try:
            await asyncio.wait_for(_run(), timeout=_PING_TIMEOUT)
        except Exception:
            return False
        return True
