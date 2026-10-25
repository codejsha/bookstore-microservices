from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.model.command.inventory_command import MonthlyClosingCommand
from internal.domain.model.error import DuplicateClosingError
from internal.domain.model.option.inventory_option import ClosingQueryOption
from internal.domain.service.inventory_service import ClosingService
from tests.conftest import make_closing


@pytest.fixture
def service(closing_repo: MagicMock, event_publisher: MagicMock) -> ClosingService:
    return ClosingService(closing_repo=closing_repo, event_publisher=event_publisher)


async def test_create_monthly_closing_delegates_and_publishes(
    service: ClosingService, closing_repo: MagicMock, event_publisher: MagicMock
) -> None:
    closing = make_closing()
    closing_repo.create.return_value = closing
    command = MonthlyClosingCommand(warehouse_uid=uuid4(), year=2026, month=10)

    result = await service.create_monthly_closing(command)

    assert result is closing
    params = closing_repo.create.call_args.args[0]
    assert (params.warehouse_uid, params.year, params.month) == (command.warehouse_uid, 2026, 10)
    assert event_publisher.publish.await_args.args[:3] == ("closing", str(closing.uid), "created")


async def test_create_monthly_closing_duplicate_propagates_without_publishing(
    service: ClosingService, closing_repo: MagicMock, event_publisher: MagicMock
) -> None:
    closing_repo.create.side_effect = DuplicateClosingError("dup")
    with pytest.raises(DuplicateClosingError):
        await service.create_monthly_closing(MonthlyClosingCommand(warehouse_uid=uuid4(), year=2026, month=10))
    event_publisher.publish.assert_not_awaited()


async def test_list_closings_delegates_to_repository(service: ClosingService, closing_repo: MagicMock) -> None:
    option = ClosingQueryOption(year=2026)
    closing_repo.find_all.return_value = ([], 0)
    assert await service.list_closings(option) == ([], 0)
    closing_repo.find_all.assert_awaited_once_with(option)


async def test_get_closing_missing_is_none(service: ClosingService, closing_repo: MagicMock) -> None:
    closing_repo.find_by_uid.return_value = None
    assert await service.get_closing(uuid4()) is None
