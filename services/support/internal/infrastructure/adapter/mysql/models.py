from datetime import UTC, datetime

from sqlalchemy import (
    BigInteger,
    Boolean,
    DateTime,
    Enum,
    LargeBinary,
    String,
    Text,
    TypeDecorator,
)
from sqlalchemy.dialects import mysql
from sqlalchemy.engine import Dialect
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column

from internal.domain.constant.author_role import AuthorRole
from internal.domain.constant.ticket_priority import TicketPriority
from internal.domain.constant.ticket_status import TicketStatus

Uid = LargeBinary(16).with_variant(mysql.BINARY(16), "mysql")


class UTCDateTime(TypeDecorator):
    impl = DateTime
    cache_ok = True

    def load_dialect_impl(self, dialect: Dialect):
        if dialect.name == "mysql":
            return dialect.type_descriptor(mysql.DATETIME(fsp=6))
        return dialect.type_descriptor(DateTime(timezone=True))

    def process_bind_param(self, value: datetime | None, dialect: Dialect) -> datetime | None:
        if value is None:
            return None
        if value.tzinfo is not None:
            value = value.astimezone(UTC)
        return value.replace(tzinfo=None)

    def process_result_value(self, value: datetime | None, dialect: Dialect) -> datetime | None:
        if value is None:
            return None
        if value.tzinfo is None:
            return value.replace(tzinfo=UTC)
        return value.astimezone(UTC)


Timestamp = UTCDateTime()
TicketStatusCol = Enum(TicketStatus, native_enum=False, length=20)
TicketPriorityCol = Enum(TicketPriority, native_enum=False, length=10)
AuthorRoleCol = Enum(AuthorRole, native_enum=False, length=20)


class Base(DeclarativeBase):
    pass


class TicketCategoryEntity(Base):
    __tablename__ = "ticket_category"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(Uid, unique=True, nullable=False, index=True)
    name: Mapped[str] = mapped_column(String(100), unique=True, nullable=False)
    description: Mapped[str | None] = mapped_column(String(500), nullable=True)
    parent_id: Mapped[int | None] = mapped_column(BigInteger, nullable=True, index=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    actor: Mapped[int] = mapped_column(BigInteger, nullable=False, default=0)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False, default=0)


class TicketEntity(Base):
    __tablename__ = "ticket"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(Uid, unique=True, nullable=False, index=True)
    customer_uid: Mapped[bytes] = mapped_column(Uid, nullable=False, index=True)
    category_id: Mapped[int | None] = mapped_column(BigInteger, nullable=True, index=True)
    assignee_uid: Mapped[bytes | None] = mapped_column(Uid, nullable=True, index=True)
    subject: Mapped[str] = mapped_column(String(255), nullable=False)
    description: Mapped[str] = mapped_column(Text, nullable=False)
    status: Mapped[TicketStatus] = mapped_column(TicketStatusCol, nullable=False, index=True)
    priority: Mapped[TicketPriority] = mapped_column(TicketPriorityCol, nullable=False, index=True)
    resolved_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False, index=True)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    actor: Mapped[int] = mapped_column(BigInteger, nullable=False, default=0)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False, default=0)


class TicketCommentEntity(Base):
    __tablename__ = "ticket_comment"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(Uid, unique=True, nullable=False, index=True)
    ticket_id: Mapped[int] = mapped_column(BigInteger, nullable=False, index=True)
    author_uid: Mapped[bytes] = mapped_column(Uid, nullable=False)
    author_role: Mapped[AuthorRole] = mapped_column(AuthorRoleCol, nullable=False)
    body: Mapped[str] = mapped_column(Text, nullable=False)
    internal: Mapped[bool] = mapped_column(Boolean, nullable=False, default=False)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False, index=True)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    actor: Mapped[int] = mapped_column(BigInteger, nullable=False, default=0)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False, default=0)


class FaqEntity(Base):
    __tablename__ = "faq"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(Uid, unique=True, nullable=False, index=True)
    category_id: Mapped[int | None] = mapped_column(BigInteger, nullable=True, index=True)
    question: Mapped[str] = mapped_column(String(500), nullable=False)
    answer: Mapped[str] = mapped_column(Text, nullable=False)
    view_count: Mapped[int] = mapped_column(BigInteger, nullable=False, default=0)
    published: Mapped[bool] = mapped_column(Boolean, nullable=False, default=False, index=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    actor: Mapped[int] = mapped_column(BigInteger, nullable=False, default=0)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False, default=0)
