import dataclasses
from types import SimpleNamespace
from unittest.mock import AsyncMock, MagicMock

import pytest
from temporalio.exceptions import ApplicationError

from internal.domain.model.error import InsufficientStockError, InvalidQuantityError
from internal.domain.service.inventory_service import StockService
from internal.infrastructure.adapter.temporal import stock_activities as activities_module
from internal.infrastructure.adapter.temporal.stock_activities import (
    RELEASE_REASON,
    RESERVE_REASON,
    ReleaseStockInput,
    ReserveStockInput,
    StockActivities,
    StockReservationItem,
    aggregate_reservation_items,
)
from tests.conftest import make_stock


@pytest.fixture
def stock_service() -> MagicMock:
    service = MagicMock(spec=StockService)
    service.reserve_stock_for_order = AsyncMock(return_value=make_stock())
    service.release_stock_for_order = AsyncMock(return_value=make_stock())
    return service


@pytest.fixture
def activities(stock_service: MagicMock) -> StockActivities:
    return StockActivities(stock_service=stock_service)


def _payload(order_uid: str = "order-1") -> ReserveStockInput:
    return ReserveStockInput(
        orderUid=order_uid,
        items=[StockReservationItem(productId=7, quantity=2), StockReservationItem(productId=9, quantity=1)],
    )


class TestPayloadContract:
    def test_reserve_stock_input_fields_are_exact_camelcase(self) -> None:
        assert [f.name for f in dataclasses.fields(ReserveStockInput)] == ["orderUid", "items"]

    def test_stock_reservation_item_fields_are_exact_camelcase(self) -> None:
        assert [f.name for f in dataclasses.fields(StockReservationItem)] == ["productId", "quantity"]

    def test_release_stock_input_is_same_shape_as_reserve(self) -> None:
        assert ReleaseStockInput is ReserveStockInput

    def test_activity_names_match_order_workflow_contract(self) -> None:
        definition_attr = "__temporal_activity_definition"
        assert getattr(StockActivities.reserve_stock, definition_attr).name == "ReserveStock"
        assert getattr(StockActivities.release_stock, definition_attr).name == "ReleaseStock"


class TestAggregateItems:
    def test_aggregate_reservation_items_duplicate_products_merge_quantities_in_first_position(self) -> None:
        merged = aggregate_reservation_items(
            [
                StockReservationItem(productId=1, quantity=2),
                StockReservationItem(productId=2, quantity=1),
                StockReservationItem(productId=1, quantity=3),
            ]
        )
        assert [(i.productId, i.quantity) for i in merged] == [(1, 5), (2, 1)]

    def test_aggregate_reservation_items_does_not_mutate_input(self) -> None:
        items = [StockReservationItem(productId=1, quantity=2), StockReservationItem(productId=1, quantity=3)]
        aggregate_reservation_items(items)
        assert items[0].quantity == 2


class TestReserveStock:
    async def test_reserve_stock_every_item_reserved_with_saga_reason(
        self, activities: StockActivities, stock_service: MagicMock
    ) -> None:
        await activities.reserve_stock(_payload())

        commands = [call.args[0] for call in stock_service.reserve_stock_for_order.await_args_list]
        assert [(c.order_uid, c.edition_id, c.quantity, c.reason) for c in commands] == [
            ("order-1", 7, 2, RESERVE_REASON),
            ("order-1", 9, 1, RESERVE_REASON),
        ]

    async def test_reserve_stock_duplicate_products_merged_before_reserving(
        self, activities: StockActivities, stock_service: MagicMock
    ) -> None:
        payload = ReserveStockInput(
            orderUid="order-1",
            items=[StockReservationItem(productId=7, quantity=2), StockReservationItem(productId=7, quantity=4)],
        )
        await activities.reserve_stock(payload)

        assert stock_service.reserve_stock_for_order.await_count == 1
        assert stock_service.reserve_stock_for_order.await_args.args[0].quantity == 6

    async def test_reserve_stock_insufficient_raises_non_retryable_insufficient_stock(
        self, activities: StockActivities, stock_service: MagicMock
    ) -> None:
        stock_service.reserve_stock_for_order = AsyncMock(side_effect=InsufficientStockError("short"))
        with pytest.raises(ApplicationError) as excinfo:
            await activities.reserve_stock(_payload())
        assert excinfo.value.non_retryable
        assert excinfo.value.type == "InsufficientStock"

    async def test_reserve_stock_invalid_quantity_raises_non_retryable_invalid_quantity(
        self, activities: StockActivities, stock_service: MagicMock
    ) -> None:
        stock_service.reserve_stock_for_order = AsyncMock(side_effect=InvalidQuantityError("zero"))
        with pytest.raises(ApplicationError) as excinfo:
            await activities.reserve_stock(_payload())
        assert excinfo.value.non_retryable
        assert excinfo.value.type == "InvalidQuantity"

    async def test_reserve_stock_invalid_product_id_raises_non_retryable_invalid_command(
        self, activities: StockActivities, stock_service: MagicMock
    ) -> None:
        payload = ReserveStockInput(orderUid="order-1", items=[StockReservationItem(productId=0, quantity=1)])
        with pytest.raises(ApplicationError) as excinfo:
            await activities.reserve_stock(payload)
        assert excinfo.value.non_retryable
        assert excinfo.value.type == "InvalidCommand"
        stock_service.reserve_stock_for_order.assert_not_awaited()

    async def test_reserve_stock_unexpected_error_propagates_as_is(
        self, activities: StockActivities, stock_service: MagicMock
    ) -> None:
        stock_service.reserve_stock_for_order = AsyncMock(side_effect=RuntimeError("db down"))
        with pytest.raises(RuntimeError):
            await activities.reserve_stock(_payload())

    async def test_reserve_stock_blank_order_uid_falls_back_to_workflow_id(
        self, activities: StockActivities, stock_service: MagicMock, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        monkeypatch.setattr(activities_module.activity, "info", lambda: SimpleNamespace(workflow_id="wf-123"))

        await activities.reserve_stock(_payload(order_uid=""))

        assert stock_service.reserve_stock_for_order.await_args.args[0].order_uid == "wf-123"


class TestReleaseStock:
    async def test_release_stock_every_item_released_with_compensation_reason(
        self, activities: StockActivities, stock_service: MagicMock
    ) -> None:
        await activities.release_stock(_payload())

        commands = [call.args[0] for call in stock_service.release_stock_for_order.await_args_list]
        assert [(c.order_uid, c.edition_id, c.reason) for c in commands] == [
            ("order-1", 7, RELEASE_REASON),
            ("order-1", 9, RELEASE_REASON),
        ]

    async def test_release_stock_unknown_reservation_completes_silently(
        self, activities: StockActivities, stock_service: MagicMock
    ) -> None:
        stock_service.release_stock_for_order = AsyncMock(return_value=None)
        await activities.release_stock(_payload())
        assert stock_service.release_stock_for_order.await_count == 2

    async def test_release_stock_invalid_product_id_raises_non_retryable_invalid_command(
        self, activities: StockActivities, stock_service: MagicMock
    ) -> None:
        payload = ReleaseStockInput(orderUid="order-1", items=[StockReservationItem(productId=-1, quantity=1)])
        with pytest.raises(ApplicationError) as excinfo:
            await activities.release_stock(payload)
        assert excinfo.value.type == "InvalidCommand"
