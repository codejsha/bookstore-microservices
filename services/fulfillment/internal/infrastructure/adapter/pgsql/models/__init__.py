from datetime import UTC, datetime

from sqlalchemy import DateTime, Float, LargeBinary, TypeDecorator
from sqlalchemy.dialects import mysql
from sqlalchemy.orm import DeclarativeBase


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
Uid = LargeBinary(16).with_variant(mysql.BINARY(16), "mysql")
Amount = Float().with_variant(mysql.DOUBLE(), "mysql")

MYSQL_TABLE_OPTS = {
    "mysql_engine": "InnoDB",
    "mysql_charset": "utf8mb4",
    "mysql_collate": "utf8mb4_unicode_ci",
}


class Base(DeclarativeBase):
    pass


from internal.infrastructure.adapter.mysql.models.delivery import (  # noqa: E402
    CarrierEntity,
    FreightEntity,
    ShipmentEntity,
    TrackingEntity,
)
from internal.infrastructure.adapter.mysql.models.inventory import (  # noqa: E402
    MonthlyClosingEntity,
    MonthlyClosingItemEntity,
    StockAuditEntity,
    StockAuditItemEntity,
    StockEntity,
    StockHistoryEntity,
    StockReservationEntity,
    StockTransferEntity,
    WarehouseEntity,
)

__all__ = [
    "MYSQL_TABLE_OPTS",
    "Amount",
    "Base",
    "CarrierEntity",
    "FreightEntity",
    "MonthlyClosingEntity",
    "MonthlyClosingItemEntity",
    "ShipmentEntity",
    "StockAuditEntity",
    "StockAuditItemEntity",
    "StockEntity",
    "StockHistoryEntity",
    "StockReservationEntity",
    "StockTransferEntity",
    "Timestamp",
    "TrackingEntity",
    "UTCDateTime",
    "Uid",
    "WarehouseEntity",
]
