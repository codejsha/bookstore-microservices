from datetime import datetime
from uuid import UUID

from sqlalchemy import BigInteger, Index, Integer, String, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column

from internal.infrastructure.adapter.pgsql.models import Base, Timestamp, Uid


class WarehouseEntity(Base):
    __tablename__ = "warehouse"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_warehouse_uid"),
        Index("idx_warehouse_name", "name"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    address: Mapped[str | None] = mapped_column(String(255), nullable=True)
    capacity: Mapped[int] = mapped_column(Integer, nullable=False)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False)


class StockEntity(Base):
    __tablename__ = "stock"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_stock_uid"),
        UniqueConstraint("edition_id", "warehouse_id", name="uq_stock_edition_warehouse"),
        Index("idx_stock_edition_id", "edition_id"),
        Index("idx_stock_warehouse_id", "warehouse_id"),
        Index("idx_stock_edition_uid_warehouse_uid", "edition_uid", "warehouse_uid"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    edition_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    edition_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    warehouse_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    warehouse_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False)


class StockHistoryEntity(Base):
    __tablename__ = "stock_history"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_stock_history_uid"),
        Index("idx_stock_history_stock_id_created", "stock_id", "created_at"),
        Index("idx_stock_history_stock_uid_created", "stock_uid", "created_at"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    stock_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    stock_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    change_type: Mapped[str] = mapped_column(String(50), nullable=False)
    reason: Mapped[str | None] = mapped_column(String(255), nullable=True)
    change_qty: Mapped[int] = mapped_column(Integer, nullable=False)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False)


class StockReservationEntity(Base):
    __tablename__ = "stock_reservation"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_stock_reservation_uid"),
        UniqueConstraint("reservation_key", "warehouse_id", name="uk_stock_reservation_key_warehouse"),
        Index("idx_stock_reservation_order_ref", "order_ref"),
        Index("idx_stock_reservation_edition_id", "edition_id"),
        Index("idx_stock_reservation_warehouse_id", "warehouse_id"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    reservation_key: Mapped[str] = mapped_column(String(255), nullable=False)
    order_ref: Mapped[str] = mapped_column(String(255), nullable=False)
    edition_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    edition_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    warehouse_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    warehouse_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    status: Mapped[str] = mapped_column(String(20), nullable=False)
    released_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False)


class StockTransferEntity(Base):
    __tablename__ = "stock_transfer"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_transfer_uid"),
        Index("idx_transfer_edition_id", "edition_id"),
        Index("idx_transfer_source", "source_warehouse_id"),
        Index("idx_transfer_target", "target_warehouse_id"),
        Index("idx_transfer_status", "status"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    edition_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    edition_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    source_warehouse_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    source_warehouse_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    target_warehouse_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    target_warehouse_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    status: Mapped[str] = mapped_column(String(20), nullable=False)
    reason: Mapped[str | None] = mapped_column(String(255), nullable=True)
    completed_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False)


class StockAuditEntity(Base):
    __tablename__ = "stock_audit"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_audit_uid"),
        Index("idx_audit_warehouse_id", "warehouse_id"),
        Index("idx_audit_status", "status"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    warehouse_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    warehouse_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    status: Mapped[str] = mapped_column(String(20), nullable=False)
    notes: Mapped[str | None] = mapped_column(String(500), nullable=True)
    completed_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False)


class StockAuditItemEntity(Base):
    __tablename__ = "stock_audit_item"
    __table_args__ = (Index("idx_audit_item_audit_id", "audit_id"),)

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    audit_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    edition_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    edition_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    system_quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    actual_quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    difference: Mapped[int] = mapped_column(Integer, nullable=False)


class MonthlyClosingEntity(Base):
    __tablename__ = "monthly_closing"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_closing_uid"),
        UniqueConstraint("warehouse_id", "year", "month", name="uq_closing_warehouse_period"),
        Index("idx_closing_warehouse_id", "warehouse_id"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    warehouse_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    warehouse_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    year: Mapped[int] = mapped_column(Integer, nullable=False)
    month: Mapped[int] = mapped_column(Integer, nullable=False)
    status: Mapped[str] = mapped_column(String(20), nullable=False)
    closed_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    version: Mapped[int] = mapped_column(BigInteger, nullable=False)


class MonthlyClosingItemEntity(Base):
    __tablename__ = "monthly_closing_item"
    __table_args__ = (Index("idx_closing_item_closing_id", "closing_id"),)

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    closing_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    edition_id: Mapped[int] = mapped_column(BigInteger, nullable=False)
    edition_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    opening_quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    inbound_quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    outbound_quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    adjust_quantity: Mapped[int] = mapped_column(Integer, nullable=False)
    closing_quantity: Mapped[int] = mapped_column(Integer, nullable=False)
