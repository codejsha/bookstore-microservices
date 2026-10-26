from dataclasses import dataclass
from datetime import UTC, datetime
from uuid import uuid7

from sqlalchemy import Label, case, func, select, update
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession

from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.error import InsufficientStockError, StockNotFoundError
from internal.infrastructure.adapter.mysql.models import StockEntity, StockHistoryEntity
from internal.infrastructure.adapter.mysql.optimistic_retry import (
    StockCreateConflictError,
    VersionConflictError,
    is_unique_violation,
)
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid

HISTORY_INBOUND_TYPES = (StockChangeType.INBOUND.value, StockChangeType.RELEASE.value)
HISTORY_OUTBOUND_TYPES = (StockChangeType.OUTBOUND.value, StockChangeType.RESERVATION.value)


@dataclass(frozen=True)
class StockDelta:
    edition_id: int
    edition_uid: bytes
    warehouse_id: int
    warehouse_uid: bytes
    delta: int
    change_type: StockChangeType
    reason: str | None
    allow_create: bool = False


@dataclass(frozen=True)
class StockMovement:
    stock_id: int
    stock_uid: bytes
    quantity: int
    version: int


def history_buckets() -> tuple[Label, Label, Label]:
    history = StockHistoryEntity
    inbound = func.coalesce(
        func.sum(case((history.change_type.in_(HISTORY_INBOUND_TYPES), history.change_qty), else_=0)), 0
    )
    outbound = func.coalesce(
        func.sum(case((history.change_type.in_(HISTORY_OUTBOUND_TYPES), -history.change_qty), else_=0)), 0
    )
    adjust = func.coalesce(
        func.sum(case((history.change_type == StockChangeType.ADJUSTMENT.value, history.change_qty), else_=0)), 0
    )
    return inbound.label("inbound"), outbound.label("outbound"), adjust.label("adjust")


async def lock_stock(session: AsyncSession, edition_uid: bytes, warehouse_uid: bytes) -> StockEntity | None:
    stmt = (
        select(StockEntity)
        .where(
            StockEntity.edition_uid == edition_uid,
            StockEntity.warehouse_uid == warehouse_uid,
            StockEntity.deleted_at.is_(None),
        )
        .with_for_update()
    )
    return (await session.execute(stmt)).scalar_one_or_none()


async def apply_stock_delta(
    session: AsyncSession, delta: StockDelta, existing: StockEntity | None = None
) -> StockMovement:
    now = datetime.now(UTC)
    stock = existing if existing is not None else await lock_stock(session, delta.edition_uid, delta.warehouse_uid)
    edition_uid = bytes_to_uuid(delta.edition_uid)
    warehouse_uid = bytes_to_uuid(delta.warehouse_uid)

    if stock is None:
        if not delta.allow_create or delta.delta < 0:
            raise StockNotFoundError(f"stock not found for edition {edition_uid} in warehouse {warehouse_uid}")
        stock = StockEntity(
            uid=uuid7().bytes,
            edition_id=delta.edition_id,
            edition_uid=delta.edition_uid,
            warehouse_id=delta.warehouse_id,
            warehouse_uid=delta.warehouse_uid,
            quantity=0,
            created_at=now,
            updated_at=None,
            version=1,
        )
        session.add(stock)
        try:
            await session.flush()
        except IntegrityError as exc:
            if is_unique_violation(exc):
                raise StockCreateConflictError(
                    f"concurrent stock row creation for edition {edition_uid} in warehouse {warehouse_uid}"
                ) from exc
            raise

    new_quantity = stock.quantity + delta.delta
    if new_quantity < 0:
        raise InsufficientStockError(
            f"insufficient stock for edition {edition_uid} in warehouse {warehouse_uid}: "
            f"have {stock.quantity}, requested {-delta.delta}"
        )

    next_version = stock.version + 1
    result = await session.execute(
        update(StockEntity)
        .where(StockEntity.id == stock.id, StockEntity.version == stock.version)
        .values(quantity=new_quantity, version=next_version, updated_at=now)
    )
    if result.rowcount == 0:
        raise VersionConflictError(f"stock {bytes_to_uuid(stock.uid)} changed concurrently")

    session.add(
        StockHistoryEntity(
            uid=uuid7().bytes,
            stock_id=stock.id,
            stock_uid=stock.uid,
            change_type=delta.change_type.value,
            reason=delta.reason,
            change_qty=delta.delta,
            created_at=now,
            updated_at=None,
            version=1,
        )
    )
    await session.flush()
    movement = StockMovement(stock_id=stock.id, stock_uid=stock.uid, quantity=new_quantity, version=next_version)
    session.expire(stock, ["quantity", "version", "updated_at"])
    return movement
