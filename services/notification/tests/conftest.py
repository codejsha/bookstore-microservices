from collections.abc import AsyncIterator
from datetime import UTC, datetime
from unittest.mock import MagicMock
from uuid import UUID, uuid4

import pytest
import pytest_asyncio
from sqlalchemy import BigInteger
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.ext.compiler import compiles
from sqlalchemy.pool import StaticPool

from internal.application.port.repo.repos import NotificationRepository, TemplateRepository
from internal.domain.aggregate.notification_aggregate import NotificationAggregate
from internal.domain.aggregate.template_aggregate import TemplateAggregate
from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus
from internal.infrastructure.adapter.mysql.models import Base


@compiles(BigInteger, "sqlite")
def _compile_bigint_sqlite(_type, _compiler, **_kw):
    return "INTEGER"


@pytest.fixture
def notification_repo() -> MagicMock:
    repo = MagicMock(spec=NotificationRepository)
    repo.save.side_effect = lambda n: n

    async def _update(uid, mutator):
        notification = repo.find_by_uid.return_value
        if notification is None:
            return None
        return mutator(notification)

    repo.update.side_effect = _update
    return repo


@pytest.fixture
def template_repo() -> MagicMock:
    repo = MagicMock(spec=TemplateRepository)
    repo.save.side_effect = lambda t: t
    return repo


@pytest_asyncio.fixture
async def session_factory() -> AsyncIterator[async_sessionmaker[AsyncSession]]:
    engine = create_async_engine(
        "sqlite+aiosqlite:///:memory:",
        connect_args={"check_same_thread": False},
        poolclass=StaticPool,
    )
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    factory = async_sessionmaker(bind=engine, expire_on_commit=False)
    try:
        yield factory
    finally:
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.drop_all)
        await engine.dispose()


def make_notification(
    *,
    uid: UUID | None = None,
    user_uid: UUID | None = None,
    order_uid: UUID | None = None,
    status: NotificationStatus = NotificationStatus.PENDING,
    notification_type: NotificationType = NotificationType.ORDER_PLACED,
    channel: Channel = Channel.EMAIL,
) -> NotificationAggregate:
    now = datetime.now(UTC)
    return NotificationAggregate(
        uid=uid or uuid4(),
        user_uid=user_uid or uuid4(),
        notification_type=notification_type,
        channel=channel,
        status=status,
        title="title",
        content="content",
        created_at=now,
        updated_at=now,
        order_uid=order_uid,
    )


def make_template(
    *,
    uid: UUID | None = None,
    notification_type: NotificationType = NotificationType.ORDER_PLACED,
    channel: Channel = Channel.EMAIL,
) -> TemplateAggregate:
    now = datetime.now(UTC)
    return TemplateAggregate(
        uid=uid or uuid4(),
        notification_type=notification_type,
        channel=channel,
        title_template="Hello {name}",
        content_template="Order {order_id} placed",
        created_at=now,
        updated_at=now,
    )
