from unittest.mock import MagicMock

import pytest

from internal.domain.model.option.inventory_option import BalanceQueryOption
from internal.domain.service.inventory_service import BalanceService
from tests.conftest import make_balance_entry


@pytest.fixture
def service(balance_repo: MagicMock) -> BalanceService:
    return BalanceService(balance_repo=balance_repo)


async def test_get_stock_balance_delegates_to_repository(service: BalanceService, balance_repo: MagicMock) -> None:
    option = BalanceQueryOption(year_month="2026-10")
    entry = make_balance_entry()
    balance_repo.find_all.return_value = ([entry], 1)

    items, total = await service.get_stock_balance(option)

    assert items == [entry]
    assert total == 1
    balance_repo.find_all.assert_awaited_once_with(option)
