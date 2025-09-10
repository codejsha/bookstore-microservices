from datetime import UTC, datetime

from sqlalchemy import BigInteger, DateTime, Enum, LargeBinary, String, Text, TypeDecorator, UniqueConstraint
from sqlalchemy.dialects import mysql
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column

from internal.domain.constant.channel import Channel
from internal.domain.constant.notification_type import NotificationType
from internal.domain.constant.status import NotificationStatus

Uid = LargeBinary(16).with_variant(mysql.BINARY(16), "mysql")


class UTCDateTime(TypeDecorator):
    impl = DateTime
    cache_ok = True

    def load_dialect_impl(self, dialect):
        if dialect.name == "mysql":
            return dialect.type_descriptor(mysql.DATETIME(fsp=6))
        return dialect.type_descriptor(DateTime(timezone=True))

    def process_bind_param(self, value: datetime | None, dialect) -> datetime | None:
        if value is None:
            return None
        if value.tzinfo is not None:
            value = value.astimezone(UTC)
        return value.replace(tzinfo=None)

    def process_result_value(self, value: datetime | None, dialect) -> datetime | None:
        if value is None:
            return None
        if value.tzinfo is None:
            return value.replace(tzinfo=UTC)
        return value.astimezone(UTC)


Timestamp = UTCDateTime()
NotificationTypeCol = Enum(NotificationType, native_enum=False, length=30)
ChannelCol = Enum(Channel, native_enum=False, length=10)
NotificationStatusCol = Enum(NotificationStatus, native_enum=False, length=10)


class Base(DeclarativeBase):
    pass


class NotificationEntity(Base):
    __tablename__ = "notification"

    __table_args__ = (UniqueConstraint("order_uid", "notification_type", name="uq_notification_order_type"),)

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(Uid, unique=True, nullable=False, index=True)
    user_uid: Mapped[bytes] = mapped_column(Uid, nullable=False, index=True)
    order_uid: Mapped[bytes | None] = mapped_column(Uid, nullable=True, index=True)
    notification_type: Mapped[NotificationType] = mapped_column(NotificationTypeCol, nullable=False)
    channel: Mapped[Channel] = mapped_column(ChannelCol, nullable=False)
    status: Mapped[NotificationStatus] = mapped_column(
        NotificationStatusCol, nullable=False, default=NotificationStatus.PENDING
    )
    title: Mapped[str] = mapped_column(String(255), nullable=False)
    content: Mapped[str] = mapped_column(Text, nullable=False)
    sent_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)


class TemplateEntity(Base):
    __tablename__ = "template"

    __table_args__ = (UniqueConstraint("notification_type", "channel", name="uq_template_type_channel"),)

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(Uid, unique=True, nullable=False, index=True)
    notification_type: Mapped[NotificationType] = mapped_column(NotificationTypeCol, nullable=False)
    channel: Mapped[Channel] = mapped_column(ChannelCol, nullable=False)
    title_template: Mapped[str] = mapped_column(String(255), nullable=False)
    content_template: Mapped[str] = mapped_column(Text, nullable=False)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
