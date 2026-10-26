from datetime import UTC, datetime

from sqlalchemy import DateTime, Double, TypeDecorator, Uuid
from sqlalchemy.dialects import postgresql
from sqlalchemy.ext.compiler import compiles
from sqlalchemy.orm import DeclarativeBase

TIMESTAMP_PRECISION = 6


class UTCDateTime(TypeDecorator):
    impl = DateTime
    cache_ok = True

    def load_dialect_impl(self, dialect):
        if dialect.name == "postgresql":
            return dialect.type_descriptor(postgresql.TIMESTAMP(precision=TIMESTAMP_PRECISION))
        return dialect.type_descriptor(DateTime())

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


@compiles(UTCDateTime, "postgresql")
def _compile_utc_datetime_postgresql(_type, _compiler, **_kw) -> str:
    return f"TIMESTAMP({TIMESTAMP_PRECISION}) WITHOUT TIME ZONE"


Timestamp = UTCDateTime()
Uid = Uuid(as_uuid=True)
Amount = Double()


class Base(DeclarativeBase):
    pass


from internal.infrastructure.adapter.pgsql.models.delivery import (  # noqa: E402
    CarrierEntity,
    FreightEntity,
    ShipmentEntity,
    TrackingEntity,
)
from internal.infrastructure.adapter.pgsql.models.inventory import (  # noqa: E402
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
