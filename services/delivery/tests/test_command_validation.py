from datetime import datetime
from uuid import uuid4

import pytest
from pydantic import ValidationError

from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.command.delivery_command import (
    AddTrackingCommand,
    AssignCarrierCommand,
    CreateCarrierCommand,
    CreateFreightCommand,
    CreateShipmentCommand,
    UpdateCarrierCommand,
)


def _shipment_kwargs(**overrides):
    kwargs = {
        "order_uid": uuid4(),
        "origin_address": "1 Warehouse Rd",
        "destination_address": "2 Home St",
        "destination_city": "Seoul",
        "destination_state": "",
        "destination_country_code": "KR",
        "destination_postal_code": "12345",
    }
    kwargs.update(overrides)
    return kwargs


class TestCreateShipmentCommand:
    def test_create_shipment_command_every_field_valid_passes(self):
        CreateShipmentCommand(**_shipment_kwargs())

    def test_create_shipment_command_blank_state_passes(self):
        CreateShipmentCommand(**_shipment_kwargs(destination_state=""))

    @pytest.mark.parametrize(
        "field", ["origin_address", "destination_address", "destination_city", "destination_postal_code"]
    )
    def test_create_shipment_command_blank_required_field_raises_validation_error(self, field):
        with pytest.raises(ValidationError):
            CreateShipmentCommand(**_shipment_kwargs(**{field: "  "}))

    def test_create_shipment_command_non_positive_weight_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateShipmentCommand(**_shipment_kwargs(weight_kg=0))

    @pytest.mark.parametrize(
        ("field", "size"),
        [
            ("origin_address", 500),
            ("destination_address", 500),
            ("destination_city", 100),
            ("destination_state", 100),
            ("destination_country_code", 3),
            ("destination_postal_code", 20),
        ],
    )
    def test_create_shipment_command_over_length_field_raises_validation_error(self, field, size):
        CreateShipmentCommand(**_shipment_kwargs(**{field: "x" * size}))
        with pytest.raises(ValidationError):
            CreateShipmentCommand(**_shipment_kwargs(**{field: "x" * (size + 1)}))

    def test_create_shipment_command_delivery_before_pickup_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateShipmentCommand(
                **_shipment_kwargs(
                    planned_pickup_at=datetime(2026, 8, 20, 10, 0),
                    planned_delivery_at=datetime(2026, 8, 19, 10, 0),
                )
            )


class TestCarrierCommands:
    def test_create_carrier_command_every_field_valid_passes(self):
        CreateCarrierCommand(name="CJ", code="CJK", base_rate=1000.0, rate_per_kg=100.0)

    def test_create_carrier_command_blank_code_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateCarrierCommand(name="CJ", code=" ", base_rate=1000.0, rate_per_kg=100.0)

    def test_create_carrier_command_negative_rate_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateCarrierCommand(name="CJ", code="CJK", base_rate=-1.0, rate_per_kg=100.0)

    def test_create_carrier_command_malformed_email_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateCarrierCommand(name="CJ", code="CJK", contact_email="nope", base_rate=1000.0, rate_per_kg=100.0)

    @pytest.mark.parametrize(
        ("field", "size"),
        [("name", 255), ("code", 50), ("contact_name", 255), ("contact_phone", 50)],
    )
    def test_create_carrier_command_over_length_field_raises_validation_error(self, field, size):
        kwargs = {"name": "CJ", "code": "CJK", "base_rate": 1000.0, "rate_per_kg": 100.0}
        CreateCarrierCommand(**{**kwargs, field: "x" * size})
        with pytest.raises(ValidationError):
            CreateCarrierCommand(**{**kwargs, field: "x" * (size + 1)})

    def test_create_carrier_command_over_length_email_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateCarrierCommand(
                name="CJ",
                code="CJK",
                contact_email="x" * 250 + "@example.com",
                base_rate=1000.0,
                rate_per_kg=100.0,
            )

    def test_update_carrier_command_over_length_name_raises_validation_error(self):
        with pytest.raises(ValidationError):
            UpdateCarrierCommand(name="x" * 256)

    def test_update_carrier_command_every_field_none_passes(self):
        UpdateCarrierCommand()

    def test_update_carrier_command_blank_name_raises_validation_error(self):
        with pytest.raises(ValidationError):
            UpdateCarrierCommand(name=" ")


class TestFreightCommand:
    def test_create_freight_command_every_field_valid_passes(self):
        CreateFreightCommand(shipment_uid=uuid4(), carrier_uid=uuid4(), base_cost=5000.0)

    def test_create_freight_command_lowercase_currency_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateFreightCommand(shipment_uid=uuid4(), carrier_uid=uuid4(), base_cost=5000.0, currency="krw")

    def test_create_freight_command_discount_exceeds_cost_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateFreightCommand(shipment_uid=uuid4(), carrier_uid=uuid4(), base_cost=1000.0, discount=2000.0)

    def test_create_freight_command_negative_surcharge_raises_validation_error(self):
        with pytest.raises(ValidationError):
            CreateFreightCommand(shipment_uid=uuid4(), carrier_uid=uuid4(), base_cost=1000.0, weight_surcharge=-1.0)


class TestTrackingAndAssign:
    def test_tracking_command_blank_location_raises_validation_error(self):
        with pytest.raises(ValidationError):
            AddTrackingCommand(
                shipment_uid=uuid4(), status=ShipmentStatus.IN_TRANSIT, location=" ", description="moving"
            )

    def test_tracking_command_over_length_location_raises_validation_error(self):
        with pytest.raises(ValidationError):
            AddTrackingCommand(
                shipment_uid=uuid4(),
                status=ShipmentStatus.IN_TRANSIT,
                location="x" * 256,
                description="moving",
            )

    def test_assign_command_blank_tracking_number_raises_validation_error(self):
        with pytest.raises(ValidationError):
            AssignCarrierCommand(carrier_uid=uuid4(), tracking_number="")

    def test_assign_command_over_length_tracking_number_raises_validation_error(self):
        AssignCarrierCommand(carrier_uid=uuid4(), tracking_number="x" * 100)
        with pytest.raises(ValidationError):
            AssignCarrierCommand(carrier_uid=uuid4(), tracking_number="x" * 101)
