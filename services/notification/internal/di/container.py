from internal.application.port.repo.repos import NotificationRepository, TemplateRepository
from internal.config.config import Settings
from internal.domain.service.notification_service import NotificationService
from internal.infrastructure.adapter.mysql.notification_repo import MySQLNotificationRepository
from internal.infrastructure.adapter.mysql.template_repo import MySQLTemplateRepository
from internal.infrastructure.adapter.temporal.activities import NotificationActivities
from internal.infrastructure.adapter.temporal.worker import TemporalWorker
from internal.infrastructure.support.database import create_session_factory


class Container:
    def __init__(self, settings: Settings):
        self.settings = settings
        self._session_factory = create_session_factory(settings.database)

        self.notification_repo: NotificationRepository = MySQLNotificationRepository(self._session_factory)
        self.template_repo: TemplateRepository = MySQLTemplateRepository(self._session_factory)

        self.notification_service = NotificationService(
            notification_repo=self.notification_repo,
            template_repo=self.template_repo,
        )

        self.notification_activities = NotificationActivities(self.notification_service)
        self.temporal_worker = TemporalWorker(settings.temporal, self.notification_activities)
