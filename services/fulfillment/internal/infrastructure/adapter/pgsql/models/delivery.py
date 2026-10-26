from datetime import datetime
from uuid import UUID

from sqlalchemy import BigInteger, Index, String, Text, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column

from internal.infrastructure.adapter.pgsql.models import Amount, Base, Timestamp, Uid


class ShipmentEntity(Base):
    __tablename__ = "shipment"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_shipment_uid"),
        UniqueConstraint("order_uid", name="uk_shipment_order_uid"),
        Index("idx_shipment_carrier_uid", "carrier_uid"),
        Index("idx_shipment_status", "status"),
        Index("idx_shipment_created", "created_at"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    order_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    carrier_uid: Mapped[UUID | None] = mapped_column(Uid, nullable=True)
    origin_address: Mapped[str] = mapped_column(String(500), nullable=False)
    destination_address: Mapped[str] = mapped_column(String(500), nullable=False)
    destination_city: Mapped[str] = mapped_column(String(100), nullable=False)
    destination_state: Mapped[str] = mapped_column(String(100), nullable=False)
    destination_country_code: Mapped[str] = mapped_column(String(3), nullable=False)
    destination_postal_code: Mapped[str] = mapped_column(String(20), nullable=False)
    status: Mapped[str] = mapped_column(String(20), nullable=False)
    tracking_number: Mapped[str | None] = mapped_column(String(100), nullable=True)
    weight_kg: Mapped[float | None] = mapped_column(Amount, nullable=True)
    planned_pickup_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    planned_delivery_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    actual_pickup_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    actual_delivery_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)


class TrackingEntity(Base):
    __tablename__ = "tracking"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_tracking_uid"),
        Index("idx_tracking_shipment_uid", "shipment_uid"),
        Index("idx_tracking_occurred", "occurred_at"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    shipment_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    status: Mapped[str] = mapped_column(String(20), nullable=False)
    location: Mapped[str] = mapped_column(String(255), nullable=False)
    description: Mapped[str] = mapped_column(Text, nullable=False)
    occurred_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)


class CarrierEntity(Base):
    __tablename__ = "carrier"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_carrier_uid"),
        UniqueConstraint("code", name="uk_carrier_code"),
        Index("idx_carrier_name", "name"),
        Index("idx_carrier_status", "status"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    code: Mapped[str] = mapped_column(String(50), nullable=False)
    contact_name: Mapped[str | None] = mapped_column(String(255), nullable=True)
    contact_phone: Mapped[str | None] = mapped_column(String(50), nullable=True)
    contact_email: Mapped[str | None] = mapped_column(String(255), nullable=True)
    base_rate: Mapped[float] = mapped_column(Amount, nullable=False)
    rate_per_kg: Mapped[float] = mapped_column(Amount, nullable=False)
    status: Mapped[str] = mapped_column(String(10), nullable=False)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)


class FreightEntity(Base):
    __tablename__ = "freight"
    __table_args__ = (
        UniqueConstraint("uid", name="uk_freight_uid"),
        Index("idx_freight_shipment_uid", "shipment_uid"),
        Index("idx_freight_carrier_uid", "carrier_uid"),
        Index("idx_freight_status", "status"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    shipment_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    carrier_uid: Mapped[UUID] = mapped_column(Uid, nullable=False)
    base_cost: Mapped[float] = mapped_column(Amount, nullable=False)
    weight_surcharge: Mapped[float] = mapped_column(Amount, nullable=False, default=0.0)
    distance_surcharge: Mapped[float] = mapped_column(Amount, nullable=False, default=0.0)
    discount: Mapped[float] = mapped_column(Amount, nullable=False, default=0.0)
    total_cost: Mapped[float] = mapped_column(Amount, nullable=False)
    currency: Mapped[str] = mapped_column(String(3), nullable=False, default="KRW")
    status: Mapped[str] = mapped_column(String(10), nullable=False)
    invoiced_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    paid_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
