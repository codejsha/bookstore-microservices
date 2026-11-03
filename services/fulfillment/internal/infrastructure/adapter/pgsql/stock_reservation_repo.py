from collections.abc import Sequence
from datetime import UTC, datetime
from uuid import uuid7

from sqlalchemy import select, update
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.inventory_repos import (
    ReleaseParams,
    ReserveParams,
    StockReservationRepository,
)
from internal.domain.aggregate.reservation_aggregate import StockReservationAggregate
from internal.domain.constant.reservation_status import ReservationStatus
from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.error import InsufficientStockError, StockNotFoundError
from internal.infrastructure.adapter.pgsql.models import StockEntity, StockReservationEntity
from internal.infrastructure.adapter.pgsql.optimistic_retry import (
    ReservationConflictError,
    VersionConflictError,
    is_unique_violation,
    run_with_optimistic_retry,
)
from internal.infrastructure.adapter.pgsql.stock_movement import StockDelta, apply_stock_delta


def allocate_reservation(rows: Sequence[StockEntity], quantity: int) -> list[tuple[StockEntity, int]]:
    ordered = sorted(rows, key=lambda stock: (-stock.quantity, stock.id))
    remaining = quantity
    allocations: list[tuple[StockEntity, int]] = []
    for stock in ordered:
        if remaining == 0:
            break
        if stock.quantity <= 0:
            continue
        take = min(stock.quantity, remaining)
        allocations.append((stock, take))
        remaining -= take
    if remaining > 0:
        raise InsufficientStockError(f"insufficient stock: requested {quantity}")
    return allocations


class PgsqlStockReservationRepository(StockReservationRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def reserve(self, params: ReserveParams) -> StockReservationAggregate:
        async def attempt() -> StockReservationAggregate:
            async with self._session_factory() as session, session.begin():
                rows = await self._lock_edition_stock(session, {params.edition_id})
                existing = await self._find_by_key(session, params.reservation_key, for_update=True)
                if existing:
                    return self._to_aggregate(existing, already_applied=True)
                if not rows:
                    raise InsufficientStockError(f"no stock found for edition_id {params.edition_id}")
                allocations = allocate_reservation(rows, params.quantity)

                now = datetime.now(UTC)
                created: list[StockReservationEntity] = []
                for stock, take in allocations:
                    stock_edition_uid = stock.edition_uid
                    stock_warehouse_id = stock.warehouse_id
                    stock_warehouse_uid = stock.warehouse_uid
                    await apply_stock_delta(
                        session,
                        StockDelta(
                            edition_id=stock.edition_id,
                            edition_uid=stock_edition_uid,
                            warehouse_id=stock_warehouse_id,
                            warehouse_uid=stock_warehouse_uid,
                            delta=-take,
                            change_type=StockChangeType.RESERVATION,
                            reason=params.reason,
                        ),
                        existing=stock,
                    )
                    reservation = StockReservationEntity(
                        uid=uuid7(),
                        reservation_key=params.reservation_key,
                        order_ref=params.order_ref,
                        edition_id=params.edition_id,
                        edition_uid=stock_edition_uid,
                        warehouse_id=stock_warehouse_id,
                        warehouse_uid=stock_warehouse_uid,
                        quantity=take,
                        status=ReservationStatus.RESERVED.value,
                        released_at=None,
                        created_at=now,
                        updated_at=None,
                        version=1,
                    )
                    session.add(reservation)
                    try:
                        await session.flush()
                    except IntegrityError as exc:
                        if is_unique_violation(exc):
                            raise ReservationConflictError(params.reservation_key) from exc
                        raise
                    created.append(reservation)
                return self._to_aggregate(created, already_applied=False)

        try:
            return await run_with_optimistic_retry(attempt)
        except ReservationConflictError:
            async with self._session_factory() as session:
                existing = await self._find_by_key(session, params.reservation_key, for_update=False)
            if existing:
                return self._to_aggregate(existing, already_applied=True)
            raise

    async def release(self, params: ReleaseParams) -> StockReservationAggregate | None:
        async def attempt() -> StockReservationAggregate | None:
            async with self._session_factory() as session, session.begin():
                pending = await self._find_by_key(session, params.reservation_key, for_update=False)
                if not pending:
                    return None
                if all(row.status == ReservationStatus.RELEASED.value for row in pending):
                    return self._to_aggregate(pending, already_applied=True)
                stocks = {
                    (stock.edition_id, stock.warehouse_id): stock
                    for stock in await self._lock_edition_stock(session, {row.edition_id for row in pending})
                }
                rows = await self._find_by_key(session, params.reservation_key, for_update=True)
                if not rows:
                    return None
                if all(row.status == ReservationStatus.RELEASED.value for row in rows):
                    return self._to_aggregate(rows, already_applied=True)

                now = datetime.now(UTC)
                first = rows[0]
                summary = (first.uid, first.reservation_key, first.edition_uid, first.warehouse_uid)
                total = sum(row.quantity for row in rows)
                for reservation in rows:
                    if reservation.status == ReservationStatus.RELEASED.value:
                        continue
                    target = stocks.get((reservation.edition_id, reservation.warehouse_id))
                    if target is None:
                        raise StockNotFoundError(
                            f"stock not found for edition {reservation.edition_uid} "
                            f"in warehouse {reservation.warehouse_uid}"
                        )
                    await apply_stock_delta(
                        session,
                        StockDelta(
                            edition_id=target.edition_id,
                            edition_uid=target.edition_uid,
                            warehouse_id=target.warehouse_id,
                            warehouse_uid=target.warehouse_uid,
                            delta=reservation.quantity,
                            change_type=StockChangeType.RELEASE,
                            reason=params.reason,
                        ),
                        existing=target,
                    )
                    marked = await session.execute(
                        update(StockReservationEntity)
                        .where(
                            StockReservationEntity.id == reservation.id,
                            StockReservationEntity.version == reservation.version,
                        )
                        .values(
                            status=ReservationStatus.RELEASED.value,
                            released_at=now,
                            version=reservation.version + 1,
                            updated_at=now,
                        )
                    )
                    if marked.rowcount == 0:
                        raise VersionConflictError(f"reservation {reservation.uid} changed concurrently")
                    session.expire(reservation, ["status", "released_at", "version", "updated_at"])

                uid, key, edition_uid, warehouse_uid = summary
                return StockReservationAggregate(
                    uid=uid,
                    reservation_key=key,
                    edition_uid=edition_uid,
                    warehouse_uid=warehouse_uid,
                    quantity=total,
                    status=ReservationStatus.RELEASED,
                    already_applied=False,
                )

        return await run_with_optimistic_retry(attempt)

    @staticmethod
    async def _find_by_key(session: AsyncSession, key: str, for_update: bool) -> list[StockReservationEntity]:
        stmt = (
            select(StockReservationEntity)
            .where(StockReservationEntity.reservation_key == key, StockReservationEntity.deleted_at.is_(None))
            .order_by(StockReservationEntity.id)
        )
        if for_update:
            stmt = stmt.with_for_update().execution_options(populate_existing=True)
        return list((await session.execute(stmt)).scalars().all())

    @staticmethod
    async def _lock_edition_stock(session: AsyncSession, edition_ids: set[int]) -> list[StockEntity]:
        stmt = (
            select(StockEntity)
            .where(StockEntity.edition_id.in_(edition_ids), StockEntity.deleted_at.is_(None))
            .order_by(StockEntity.id)
            .with_for_update()
            .execution_options(populate_existing=True)
        )
        return list((await session.execute(stmt)).scalars().all())

    @staticmethod
    def _to_aggregate(rows: Sequence[StockReservationEntity], already_applied: bool) -> StockReservationAggregate:
        first = rows[0]
        return StockReservationAggregate(
            uid=first.uid,
            reservation_key=first.reservation_key,
            edition_uid=first.edition_uid,
            warehouse_uid=first.warehouse_uid,
            quantity=sum(row.quantity for row in rows),
            status=ReservationStatus(first.status),
            already_applied=already_applied,
        )
