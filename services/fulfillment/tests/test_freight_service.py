from unittest.mock import MagicMock
from uuid import uuid4

import pytest

from internal.domain.constant.freight_status import FreightStatus
from internal.domain.model.command.delivery_command import (
    CreateFreightCommand,
    UpdateFreightStatusCommand,
)
from internal.domain.model.option.delivery_option import FreightFilterOption
from internal.domain.service.delivery_service import FreightService
from tests.conftest import make_freight


@pytest.fixture
def service(freight_repo: MagicMock) -> FreightService:
    return FreightService(freight_repo=freight_repo)


class TestCreateFreight:
    async def test_create_freight_with_surcharges_computes_total_cost(
        self, service: FreightService, freight_repo: MagicMock
    ) -> None:
        command = CreateFreightCommand(
            shipment_uid=uuid4(),
            carrier_uid=uuid4(),
            base_cost=1000.0,
            weight_surcharge=200.0,
            distance_surcharge=300.0,
            discount=100.0,
        )
        result = await service.create_freight(command)
        assert result.total_cost == pytest.approx(1400.0)
        assert result.currency == "KRW"
        assert result.status == FreightStatus.ESTIMATED
        freight_repo.save.assert_called_once()

    async def test_create_freight_without_surcharges_defaults_them_to_zero(
        self, service: FreightService, freight_repo: MagicMock
    ) -> None:
        command = CreateFreightCommand(shipment_uid=uuid4(), carrier_uid=uuid4(), base_cost=500.0)
        result = await service.create_freight(command)
        assert result.total_cost == pytest.approx(500.0)


class TestGetFreight:
    async def test_get_freight_roundtrips(self, service: FreightService, freight_repo: MagicMock) -> None:
        freight = make_freight()
        freight_repo.find_by_uid.return_value = freight
        assert await service.get_freight(freight.uid) is freight


class TestFindByShipment:
    async def test_find_by_shipment_delegates_to_repository(
        self, service: FreightService, freight_repo: MagicMock
    ) -> None:
        shipment_uid = uuid4()
        freight = make_freight(shipment_uid=shipment_uid)
        freight_repo.find_by_shipment_uid.return_value = freight
        assert await service.find_by_shipment(shipment_uid) is freight
        freight_repo.find_by_shipment_uid.assert_called_once_with(shipment_uid)


class TestListFreights:
    async def test_list_freights_with_option_delegates_to_repository(
        self, service: FreightService, freight_repo: MagicMock
    ) -> None:
        option = FreightFilterOption()
        freight_repo.find_all.return_value = ([make_freight()], 1)
        result, total = await service.list_freights(option)
        assert total == 1


class TestUpdateStatus:
    async def test_update_status_missing_is_none(self, service: FreightService, freight_repo: MagicMock) -> None:
        freight_repo.find_by_uid.return_value = None
        result = await service.update_status(uuid4(), UpdateFreightStatusCommand(status=FreightStatus.PAID))
        assert result is None
        freight_repo.save.assert_not_called()

    async def test_update_status_invoiced_sets_invoiced_at(
        self, service: FreightService, freight_repo: MagicMock
    ) -> None:
        freight = make_freight(status=FreightStatus.CONFIRMED)
        freight_repo.find_by_uid.return_value = freight
        result = await service.update_status(freight.uid, UpdateFreightStatusCommand(status=FreightStatus.INVOICED))
        assert result is not None
        assert result.status == FreightStatus.INVOICED
        assert result.invoiced_at is not None
        assert result.paid_at is None

    async def test_update_status_paid_sets_paid_at(self, service: FreightService, freight_repo: MagicMock) -> None:
        freight = make_freight(status=FreightStatus.INVOICED)
        freight_repo.find_by_uid.return_value = freight
        result = await service.update_status(freight.uid, UpdateFreightStatusCommand(status=FreightStatus.PAID))
        assert result is not None
        assert result.status == FreightStatus.PAID
        assert result.paid_at is not None

    async def test_update_status_neither_invoiced_nor_paid_leaves_timestamps(
        self, service: FreightService, freight_repo: MagicMock
    ) -> None:
        freight = make_freight(status=FreightStatus.ESTIMATED)
        freight_repo.find_by_uid.return_value = freight
        result = await service.update_status(freight.uid, UpdateFreightStatusCommand(status=FreightStatus.CONFIRMED))
        assert result is not None
        assert result.status == FreightStatus.CONFIRMED
        assert result.invoiced_at is None
        assert result.paid_at is None
