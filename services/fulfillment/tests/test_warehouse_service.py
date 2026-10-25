from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.model.command.inventory_command import WarehouseCreateCommand, WarehouseUpdateCommand
from internal.domain.model.option.inventory_option import WarehouseQueryOption
from internal.domain.service.inventory_service import WarehouseService
from tests.conftest import make_warehouse


@pytest.fixture
def service(warehouse_repo: MagicMock, event_publisher: MagicMock) -> WarehouseService:
    return WarehouseService(warehouse_repo=warehouse_repo, event_publisher=event_publisher)


async def test_register_warehouse_creates_and_publishes(
    service: WarehouseService, warehouse_repo: MagicMock, event_publisher: MagicMock
) -> None:
    warehouse = make_warehouse()
    warehouse_repo.create.return_value = warehouse

    result = await service.register_warehouse(WarehouseCreateCommand(name="Main", capacity=100))

    assert result is warehouse
    params = warehouse_repo.create.call_args.args[0]
    assert (params.name, params.address, params.capacity) == ("Main", None, 100)
    assert event_publisher.publish.await_args.args[:3] == ("warehouse", str(warehouse.uid), "created")


async def test_update_warehouse_missing_is_none_without_publishing(
    service: WarehouseService, warehouse_repo: MagicMock, event_publisher: MagicMock
) -> None:
    warehouse_repo.update.return_value = None

    assert await service.update_warehouse(uuid4(), WarehouseUpdateCommand(name="x")) is None
    event_publisher.publish.assert_not_awaited()


async def test_update_warehouse_found_passes_fields_and_publishes(
    service: WarehouseService, warehouse_repo: MagicMock, event_publisher: MagicMock
) -> None:
    warehouse = make_warehouse()
    warehouse_repo.update.return_value = warehouse

    result = await service.update_warehouse(warehouse.uid, WarehouseUpdateCommand(capacity=5))

    assert result is warehouse
    params = warehouse_repo.update.call_args.args[0]
    assert (params.uid, params.name, params.capacity) == (warehouse.uid, None, 5)
    assert event_publisher.publish.await_args.args[:3] == ("warehouse", str(warehouse.uid), "updated")


async def test_list_warehouses_delegates_to_repository(service: WarehouseService, warehouse_repo: MagicMock) -> None:
    option = WarehouseQueryOption(name="x")
    warehouse_repo.find_all.return_value = ([], 0)

    assert await service.list_warehouses(option) == ([], 0)
    warehouse_repo.find_all.assert_awaited_once_with(option)


async def test_get_warehouse_missing_is_none(service: WarehouseService, warehouse_repo: MagicMock) -> None:
    warehouse_repo.find_by_uid.return_value = None
    assert await service.get_warehouse(uuid4()) is None
