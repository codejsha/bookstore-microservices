from dataclasses import dataclass, field

import structlog
from pydantic import ValidationError
from temporalio import activity
from temporalio.exceptions import ApplicationError

from internal.domain.model.command.inventory_command import StockOrderReleaseCommand, StockOrderReserveCommand
from internal.domain.model.error import InsufficientStockError, InvalidQuantityError
from internal.domain.service.inventory_service import StockService

logger = structlog.get_logger()

RESERVE_REASON = "temporal-saga-reserve"
RELEASE_REASON = "temporal-saga-compensation"

ERROR_TYPE_INSUFFICIENT_STOCK = "InsufficientStock"
ERROR_TYPE_INVALID_QUANTITY = "InvalidQuantity"
ERROR_TYPE_INVALID_COMMAND = "InvalidCommand"


@dataclass
class StockReservationItem:
    productId: int  # noqa: N815
    quantity: int


@dataclass
class ReserveStockInput:
    orderUid: str  # noqa: N815
    items: list[StockReservationItem] = field(default_factory=list)


ReleaseStockInput = ReserveStockInput


def aggregate_reservation_items(items: list[StockReservationItem]) -> list[StockReservationItem]:
    merged: list[StockReservationItem] = []
    position: dict[int, int] = {}
    for item in items:
        index = position.get(item.productId)
        if index is not None:
            merged[index].quantity += item.quantity
            continue
        position[item.productId] = len(merged)
        merged.append(StockReservationItem(productId=item.productId, quantity=item.quantity))
    return merged


def _non_retryable(exc: Exception, error_type: str) -> ApplicationError:
    return ApplicationError(str(exc), non_retryable=True, type=error_type)


class StockActivities:
    def __init__(self, stock_service: StockService):
        self._stock_service = stock_service

    @activity.defn(name="ReserveStock")
    async def reserve_stock(self, payload: ReserveStockInput) -> None:
        order_ref = _order_ref(payload)
        for item in aggregate_reservation_items(payload.items):
            try:
                command = StockOrderReserveCommand(
                    order_uid=order_ref,
                    edition_id=item.productId,
                    quantity=item.quantity,
                    reason=RESERVE_REASON,
                )
                await self._stock_service.reserve_stock_for_order(command)
            except ValidationError as exc:
                logger.error("Failed to reserve stock", order_ref=order_ref, product_id=item.productId, error=str(exc))
                raise _non_retryable(exc, ERROR_TYPE_INVALID_COMMAND) from exc
            except InsufficientStockError as exc:
                logger.error("Failed to reserve stock", order_ref=order_ref, product_id=item.productId, error=str(exc))
                raise _non_retryable(exc, ERROR_TYPE_INSUFFICIENT_STOCK) from exc
            except InvalidQuantityError as exc:
                logger.error("Failed to reserve stock", order_ref=order_ref, product_id=item.productId, error=str(exc))
                raise _non_retryable(exc, ERROR_TYPE_INVALID_QUANTITY) from exc
            logger.info("Reserved stock", order_ref=order_ref, product_id=item.productId, quantity=item.quantity)

    @activity.defn(name="ReleaseStock")
    async def release_stock(self, payload: ReleaseStockInput) -> None:
        order_ref = _order_ref(payload)
        for item in aggregate_reservation_items(payload.items):
            try:
                command = StockOrderReleaseCommand(
                    order_uid=order_ref,
                    edition_id=item.productId,
                    quantity=item.quantity,
                    reason=RELEASE_REASON,
                )
                await self._stock_service.release_stock_for_order(command)
            except ValidationError as exc:
                logger.error("Failed to release stock", order_ref=order_ref, product_id=item.productId, error=str(exc))
                raise _non_retryable(exc, ERROR_TYPE_INVALID_COMMAND) from exc
            except InsufficientStockError as exc:
                logger.error("Failed to release stock", order_ref=order_ref, product_id=item.productId, error=str(exc))
                raise _non_retryable(exc, ERROR_TYPE_INSUFFICIENT_STOCK) from exc
            except InvalidQuantityError as exc:
                logger.error("Failed to release stock", order_ref=order_ref, product_id=item.productId, error=str(exc))
                raise _non_retryable(exc, ERROR_TYPE_INVALID_QUANTITY) from exc
            logger.info("Released stock", order_ref=order_ref, product_id=item.productId, quantity=item.quantity)


def _order_ref(payload: ReserveStockInput) -> str:
    if payload.orderUid:
        return payload.orderUid
    return activity.info().workflow_id
