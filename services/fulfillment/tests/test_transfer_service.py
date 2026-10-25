from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.constant.transfer_status import TransferStatus
from internal.domain.model.command.inventory_command import StockTransferCommand
from internal.domain.model.error import InvalidTransitionError, SameWarehouseError
from internal.domain.service.inventory_service import TransferService
from tests.conftest import make_transfer


@pytest.fixture
def service(transfer_repo: MagicMock, event_publisher: MagicMock) -> TransferService:
    return TransferService(transfer_repo=transfer_repo, event_publisher=event_publisher)


def _command() -> StockTransferCommand:
    return StockTransferCommand.model_construct(
        edition_uid=uuid4(), source_warehouse_uid=uuid4(), target_warehouse_uid=uuid4(), quantity=2, reason=None
    )


async def test_transfer_stock_creates_pending_transfer_and_publishes(
    service: TransferService, transfer_repo: MagicMock, event_publisher: MagicMock
) -> None:
    transfer = make_transfer()
    transfer_repo.create.return_value = transfer

    result = await service.transfer_stock(_command())

    assert result is transfer
    assert event_publisher.publish.await_args.args[:3] == ("transfer", str(transfer.uid), "created")


async def test_transfer_stock_same_warehouse_raises_same_warehouse_error(
    service: TransferService, transfer_repo: MagicMock
) -> None:
    warehouse = uuid4()
    command = StockTransferCommand.model_construct(
        edition_uid=uuid4(), source_warehouse_uid=warehouse, target_warehouse_uid=warehouse, quantity=2, reason=None
    )
    with pytest.raises(SameWarehouseError):
        await service.transfer_stock(command)
    transfer_repo.create.assert_not_awaited()


async def test_complete_transfer_missing_is_none(service: TransferService, transfer_repo: MagicMock) -> None:
    transfer_repo.find_by_uid.return_value = None
    assert await service.complete_transfer(uuid4()) is None
    transfer_repo.complete.assert_not_awaited()


@pytest.mark.parametrize("status", [TransferStatus.COMPLETED, TransferStatus.CANCELLED])
async def test_complete_transfer_finished_status_raises_invalid_transition(
    service: TransferService, transfer_repo: MagicMock, status: TransferStatus
) -> None:
    transfer_repo.find_by_uid.return_value = make_transfer(status=status)
    with pytest.raises(InvalidTransitionError):
        await service.complete_transfer(uuid4())
    transfer_repo.complete.assert_not_awaited()


@pytest.mark.parametrize("status", [TransferStatus.PENDING, TransferStatus.IN_TRANSIT])
async def test_complete_transfer_open_status_delegates_and_publishes(
    service: TransferService, transfer_repo: MagicMock, event_publisher: MagicMock, status: TransferStatus
) -> None:
    transfer = make_transfer(status=status)
    transfer_repo.find_by_uid.return_value = transfer
    transfer_repo.complete.return_value = make_transfer(uid=transfer.uid, status=TransferStatus.COMPLETED)

    result = await service.complete_transfer(transfer.uid)

    assert result is not None and result.status == TransferStatus.COMPLETED
    assert event_publisher.publish.await_args.args[:3] == ("transfer", str(transfer.uid), "completed")


async def test_cancel_transfer_in_transit_raises_invalid_transition(
    service: TransferService, transfer_repo: MagicMock
) -> None:
    transfer_repo.find_by_uid.return_value = make_transfer(status=TransferStatus.IN_TRANSIT)
    with pytest.raises(InvalidTransitionError):
        await service.cancel_transfer(uuid4())


async def test_cancel_transfer_pending_delegates_and_publishes(
    service: TransferService, transfer_repo: MagicMock, event_publisher: MagicMock
) -> None:
    transfer = make_transfer(status=TransferStatus.PENDING)
    transfer_repo.find_by_uid.return_value = transfer
    transfer_repo.cancel.return_value = make_transfer(uid=transfer.uid, status=TransferStatus.CANCELLED)

    result = await service.cancel_transfer(transfer.uid)

    assert result is not None and result.status == TransferStatus.CANCELLED
    assert event_publisher.publish.await_args.args[:3] == ("transfer", str(transfer.uid), "cancelled")
