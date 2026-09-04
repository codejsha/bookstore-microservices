from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.constant.carrier_status import CarrierStatus
from internal.domain.model.command.delivery_command import (
    CreateCarrierCommand,
    UpdateCarrierCommand,
)
from internal.domain.model.error import ConflictError
from internal.domain.model.option.delivery_option import CarrierFilterOption
from internal.domain.service.delivery_service import CarrierService
from tests.conftest import make_carrier


@pytest.fixture
def service(carrier_repo: MagicMock) -> CarrierService:
    return CarrierService(carrier_repo=carrier_repo)


class TestCreateCarrier:
    async def test_create_carrier_unused_code_persists_active_carrier(
        self, service: CarrierService, carrier_repo: MagicMock
    ) -> None:
        command = CreateCarrierCommand(
            name="ACME",
            code="ACME",
            base_rate=1000.0,
            rate_per_kg=500.0,
        )
        result = await service.create_carrier(command)
        assert result.name == "ACME"
        assert result.status == CarrierStatus.ACTIVE
        carrier_repo.save.assert_called_once()

    async def test_create_carrier_duplicate_code_raises_conflict_error(
        self, service: CarrierService, carrier_repo: MagicMock
    ) -> None:
        carrier_repo.find_by_code.return_value = make_carrier()
        command = CreateCarrierCommand(name="ACME", code="ACME", base_rate=1000.0, rate_per_kg=500.0)
        with pytest.raises(ConflictError):
            await service.create_carrier(command)
        carrier_repo.save.assert_not_called()


class TestGetCarrier:
    async def test_get_carrier_roundtrips(self, service: CarrierService, carrier_repo: MagicMock) -> None:
        carrier = make_carrier()
        carrier_repo.find_by_uid.return_value = carrier
        assert await service.get_carrier(carrier.uid) is carrier

    async def test_get_carrier_missing_is_none(self, service: CarrierService, carrier_repo: MagicMock) -> None:
        carrier_repo.find_by_uid.return_value = None
        assert await service.get_carrier(uuid4()) is None


class TestListCarriers:
    async def test_list_carriers_with_option_delegates_to_repository(
        self, service: CarrierService, carrier_repo: MagicMock
    ) -> None:
        option = CarrierFilterOption()
        carrier_repo.find_all.return_value = ([make_carrier()], 1)
        result, total = await service.list_carriers(option)
        assert total == 1
        carrier_repo.find_all.assert_called_once_with(option)


class TestUpdateCarrier:
    async def test_update_carrier_missing_is_none(self, service: CarrierService, carrier_repo: MagicMock) -> None:
        carrier_repo.find_by_uid.return_value = None
        result = await service.update_carrier(uuid4(), UpdateCarrierCommand(name="X"))
        assert result is None
        carrier_repo.save.assert_not_called()

    async def test_update_carrier_partial_fields_updates_only_those(
        self, service: CarrierService, carrier_repo: MagicMock
    ) -> None:
        carrier = make_carrier()
        original_name = carrier.name
        original_code = carrier.code
        carrier_repo.find_by_uid.return_value = carrier
        command = UpdateCarrierCommand(contact_name="Jane", base_rate=2000.0)
        result = await service.update_carrier(carrier.uid, command)
        assert result is not None
        assert result.name == original_name
        assert result.code == original_code
        assert result.contact_name == "Jane"
        assert result.base_rate == 2000.0

    async def test_update_carrier_with_status_changes_status(
        self, service: CarrierService, carrier_repo: MagicMock
    ) -> None:
        carrier = make_carrier(status=CarrierStatus.ACTIVE)
        carrier_repo.find_by_uid.return_value = carrier
        result = await service.update_carrier(carrier.uid, UpdateCarrierCommand(status=CarrierStatus.INACTIVE))
        assert result is not None
        assert result.status == CarrierStatus.INACTIVE
