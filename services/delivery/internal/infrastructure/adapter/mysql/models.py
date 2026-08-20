from datetime import UTC, datetime

from sqlalchemy import BigInteger, DateTime, Float, LargeBinary, String, Text, TypeDecorator
from sqlalchemy.dialects import mysql
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


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


class Base(DeclarativeBase):
    pass


class ShipmentEntity(Base):
    __tablename__ = "shipment"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(LargeBinary(16), unique=True, nullable=False, index=True)
    order_uid: Mapped[bytes] = mapped_column(LargeBinary(16), unique=True, nullable=False, index=True)
    carrier_uid: Mapped[bytes | None] = mapped_column(LargeBinary(16), nullable=True, index=True)
    origin_address: Mapped[str] = mapped_column(String(500), nullable=False)
    destination_address: Mapped[str] = mapped_column(String(500), nullable=False)
    destination_city: Mapped[str] = mapped_column(String(100), nullable=False)
    destination_state: Mapped[str] = mapped_column(String(100), nullable=False)
    destination_country_code: Mapped[str] = mapped_column(String(3), nullable=False)
    destination_postal_code: Mapped[str] = mapped_column(String(20), nullable=False)
    status: Mapped[str] = mapped_column(String(20), nullable=False)
    tracking_number: Mapped[str | None] = mapped_column(String(100), nullable=True)
    weight_kg: Mapped[float | None] = mapped_column(Float, nullable=True)
    planned_pickup_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    planned_delivery_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    actual_pickup_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    actual_delivery_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)


class TrackingEntity(Base):
    __tablename__ = "tracking"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(LargeBinary(16), unique=True, nullable=False, index=True)
    shipment_uid: Mapped[bytes] = mapped_column(LargeBinary(16), nullable=False, index=True)
    status: Mapped[str] = mapped_column(String(20), nullable=False)
    location: Mapped[str] = mapped_column(String(255), nullable=False)
    description: Mapped[str] = mapped_column(Text, nullable=False)
    occurred_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)


class CarrierEntity(Base):
    __tablename__ = "carrier"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(LargeBinary(16), unique=True, nullable=False, index=True)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    code: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    contact_name: Mapped[str | None] = mapped_column(String(255), nullable=True)
    contact_phone: Mapped[str | None] = mapped_column(String(50), nullable=True)
    contact_email: Mapped[str | None] = mapped_column(String(255), nullable=True)
    base_rate: Mapped[float] = mapped_column(Float, nullable=False)
    rate_per_kg: Mapped[float] = mapped_column(Float, nullable=False)
    status: Mapped[str] = mapped_column(String(10), nullable=False)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)


class FreightEntity(Base):
    __tablename__ = "freight"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    uid: Mapped[bytes] = mapped_column(LargeBinary(16), unique=True, nullable=False, index=True)
    shipment_uid: Mapped[bytes] = mapped_column(LargeBinary(16), nullable=False, index=True)
    carrier_uid: Mapped[bytes] = mapped_column(LargeBinary(16), nullable=False, index=True)
    base_cost: Mapped[float] = mapped_column(Float, nullable=False)
    weight_surcharge: Mapped[float] = mapped_column(Float, nullable=False, default=0.0)
    distance_surcharge: Mapped[float] = mapped_column(Float, nullable=False, default=0.0)
    discount: Mapped[float] = mapped_column(Float, nullable=False, default=0.0)
    total_cost: Mapped[float] = mapped_column(Float, nullable=False)
    currency: Mapped[str] = mapped_column(String(3), nullable=False, default="KRW")
    status: Mapped[str] = mapped_column(String(10), nullable=False)
    invoiced_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    paid_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
    created_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    updated_at: Mapped[datetime] = mapped_column(Timestamp, nullable=False)
    deleted_at: Mapped[datetime | None] = mapped_column(Timestamp, nullable=True)
