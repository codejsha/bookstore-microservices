from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.constant.stock_change_type import StockChangeType
from internal.domain.model.command.inventory_command import (
    StockAdjustCommand,
    StockOrderReleaseCommand,
    StockOrderReserveCommand,
    StockReceiveCommand,
    StockReserveCommand,
)
from internal.domain.model.error import InvalidQuantityError, StockNotFoundError
from internal.domain.model.option.inventory_option import StockHistoryQueryOption, StockQueryOption
from internal.domain.service.inventory_service import StockService, reservation_key
from tests.conftest import make_reservation, make_stock


@pytest.fixture
def service(
    stock_repo: MagicMock, stock_history_repo: MagicMock, reservation_repo: MagicMock, event_publisher: MagicMock
) -> StockService:
    return StockService(
        stock_repo=stock_repo,
        stock_history_repo=stock_history_repo,
        reservation_repo=reservation_repo,
        event_publisher=event_publisher,
    )


def test_reservation_key_joins_order_and_edition() -> None:
    assert reservation_key("order-1", 42) == "order-1:42"


class TestStockChanges:
    async def test_receive_stock_applies_positive_inbound_delta_and_publishes(
        self, service: StockService, stock_repo: MagicMock, event_publisher: MagicMock
    ) -> None:
        stock = make_stock()
        stock_repo.find_by_edition_uid.return_value = stock
        command = StockReceiveCommand(edition_uid=stock.edition_uid, warehouse_uid=uuid4(), quantity=5, reason="r")

        result = await service.receive_stock(command)

        assert result is stock
        params = stock_repo.apply_change.call_args.args[0]
        assert (params.delta, params.change_type, params.reason) == (5, StockChangeType.INBOUND, "r")
        event_publisher.publish.assert_awaited_once()
        assert event_publisher.publish.await_args.args[:3] == ("stock", str(stock.uid), "received")

    async def test_reserve_stock_applies_negative_reservation_delta(
        self, service: StockService, stock_repo: MagicMock
    ) -> None:
        stock_repo.find_by_edition_uid.return_value = make_stock()
        command = StockReserveCommand(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=4)

        await service.reserve_stock(command)

        params = stock_repo.apply_change.call_args.args[0]
        assert (params.delta, params.change_type) == (-4, StockChangeType.RESERVATION)

    async def test_adjust_stock_keeps_signed_delta(self, service: StockService, stock_repo: MagicMock) -> None:
        stock_repo.find_by_edition_uid.return_value = make_stock()
        command = StockAdjustCommand(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=-3, reason="damaged")

        await service.adjust_stock(command)

        params = stock_repo.apply_change.call_args.args[0]
        assert (params.delta, params.change_type, params.reason) == (-3, StockChangeType.ADJUSTMENT, "damaged")

    async def test_receive_stock_repository_error_propagates_without_publishing(
        self, service: StockService, stock_repo: MagicMock, event_publisher: MagicMock
    ) -> None:
        stock_repo.apply_change.side_effect = StockNotFoundError("missing")
        command = StockReceiveCommand(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=5)

        with pytest.raises(StockNotFoundError):
            await service.receive_stock(command)
        event_publisher.publish.assert_not_awaited()


class TestOrderReservations:
    async def test_reserve_stock_for_order_builds_key_and_returns_stock(
        self, service: StockService, reservation_repo: MagicMock, stock_repo: MagicMock
    ) -> None:
        stock = make_stock()
        reservation_repo.reserve.return_value = make_reservation(edition_uid=stock.edition_uid)
        stock_repo.find_by_edition_uid.return_value = stock

        result = await service.reserve_stock_for_order(
            StockOrderReserveCommand(order_uid="order-1", edition_id=7, quantity=2, reason="saga")
        )

        assert result is stock
        params = reservation_repo.reserve.call_args.args[0]
        assert params.reservation_key == "order-1:7"
        assert (params.order_ref, params.edition_id, params.quantity, params.reason) == ("order-1", 7, 2, "saga")
        stock_repo.find_by_edition_uid.assert_awaited_once_with(stock.edition_uid)

    async def test_reserve_stock_for_order_non_positive_quantity_raises_invalid_quantity(
        self, service: StockService, reservation_repo: MagicMock
    ) -> None:
        with pytest.raises(InvalidQuantityError):
            await service.reserve_stock_for_order(StockOrderReserveCommand(order_uid="o", edition_id=1, quantity=0))
        reservation_repo.reserve.assert_not_awaited()

    async def test_release_stock_for_order_unknown_reservation_is_none(
        self, service: StockService, reservation_repo: MagicMock, stock_repo: MagicMock
    ) -> None:
        reservation_repo.release.return_value = None

        result = await service.release_stock_for_order(StockOrderReleaseCommand(order_uid="o", edition_id=1))

        assert result is None
        assert reservation_repo.release.call_args.args[0].reservation_key == "o:1"
        stock_repo.find_by_edition_uid.assert_not_awaited()

    async def test_release_stock_for_order_released_returns_stock(
        self, service: StockService, reservation_repo: MagicMock, stock_repo: MagicMock
    ) -> None:
        stock = make_stock()
        reservation_repo.release.return_value = make_reservation(edition_uid=stock.edition_uid)
        stock_repo.find_by_edition_uid.return_value = stock

        assert await service.release_stock_for_order(StockOrderReleaseCommand(order_uid="o", edition_id=1)) is stock


class TestQueries:
    async def test_list_stocks_delegates_to_repository(self, service: StockService, stock_repo: MagicMock) -> None:
        option = StockQueryOption()
        stock_repo.find_all.return_value = ([make_stock()], 1)

        items, total = await service.list_stocks(option)

        assert total == 1
        stock_repo.find_all.assert_awaited_once_with(option)

    async def test_get_stock_history_without_stock_rows_raises_stock_not_found(
        self, service: StockService, stock_repo: MagicMock, stock_history_repo: MagicMock
    ) -> None:
        stock_repo.find_stock_uids.return_value = []

        with pytest.raises(StockNotFoundError):
            await service.get_stock_history(uuid4(), StockHistoryQueryOption())
        stock_history_repo.find_all.assert_not_awaited()

    async def test_get_stock_history_with_stock_rows_queries_their_uids(
        self, service: StockService, stock_repo: MagicMock, stock_history_repo: MagicMock
    ) -> None:
        uids = [uuid4(), uuid4()]
        stock_repo.find_stock_uids.return_value = uids
        stock_history_repo.find_all.return_value = ([], 0)
        option = StockHistoryQueryOption(warehouse_uid=uuid4())

        await service.get_stock_history(uuid4(), option)

        stock_history_repo.find_all.assert_awaited_once_with(uids, option)
