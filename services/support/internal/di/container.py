import asyncio

from sqlalchemy import text

from internal.application.port.repo.repos import (
    FaqRepository,
    TicketCategoryRepository,
    TicketCommentRepository,
    TicketRepository,
)
from internal.config.config import Settings, vault_env_prefix
from internal.domain.service.support_service import SupportService
from internal.infrastructure.adapter.mysql.faq_repo import MySQLFaqRepository
from internal.infrastructure.adapter.mysql.ticket_category_repo import (
    MySQLTicketCategoryRepository,
)
from internal.infrastructure.adapter.mysql.ticket_comment_repo import (
    MySQLTicketCommentRepository,
)
from internal.infrastructure.adapter.mysql.ticket_repo import MySQLTicketRepository
from internal.infrastructure.support.database import (
    create_engine_and_session_factory,
    create_readiness_engine,
)

_PING_TIMEOUT = 3.0


class Container:
    def __init__(self, settings: Settings):
        self.settings = settings
        self.engine, self._session_factory = create_engine_and_session_factory(settings.database, vault_env_prefix())
        self._readiness_engine = create_readiness_engine(settings.database, vault_env_prefix())

        self.ticket_repo: TicketRepository = MySQLTicketRepository(self._session_factory)
        self.comment_repo: TicketCommentRepository = MySQLTicketCommentRepository(self._session_factory)
        self.category_repo: TicketCategoryRepository = MySQLTicketCategoryRepository(self._session_factory)
        self.faq_repo: FaqRepository = MySQLFaqRepository(self._session_factory)

        self.support_service = SupportService(
            ticket_repo=self.ticket_repo,
            comment_repo=self.comment_repo,
            category_repo=self.category_repo,
            faq_repo=self.faq_repo,
        )

    async def ping_db(self) -> bool:
        async def _run() -> None:
            async with self._readiness_engine.connect() as conn:
                await conn.execute(text("SELECT 1"))

        try:
            await asyncio.wait_for(_run(), timeout=_PING_TIMEOUT)
        except Exception:
            return False
        return True

    async def close_readiness(self) -> None:
        await self._readiness_engine.dispose()
