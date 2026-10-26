from uuid import uuid4

import pytest
from pydantic import ValidationError

from internal.domain.model.command.inventory_command import (
    MonthlyClosingCommand,
    StockAdjustCommand,
    StockAuditCreateCommand,
    StockAuditItemCommand,
    StockOrderReleaseCommand,
    StockOrderReserveCommand,
    StockReceiveCommand,
    StockReleaseCommand,
    StockReserveCommand,
    StockTransferCommand,
    WarehouseCreateCommand,
    WarehouseUpdateCommand,
)


class TestStockCommands:
    @pytest.mark.parametrize("command_type", [StockReceiveCommand, StockReleaseCommand, StockReserveCommand])
    def test_stock_command_every_field_valid_passes(self, command_type):
        command_type(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=5, reason="restock")

    @pytest.mark.parametrize("command_type", [StockReceiveCommand, StockReleaseCommand, StockReserveCommand])
    @pytest.mark.parametrize("quantity", [0, -1])
    def test_stock_command_non_positive_quantity_raises_validation_error(self, command_type, quantity):
        with pytest.raises(ValidationError):
            command_type(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=quantity)

    def test_stock_command_non_uuid_edition_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockReceiveCommand(edition_uid="ed-1", warehouse_uid=uuid4(), quantity=1)

    def test_stock_command_blank_reason_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockReceiveCommand(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=1, reason=" ")

    def test_stock_command_over_length_reason_raises_validation_error(self):
        StockReserveCommand(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=1, reason="r" * 255)
        with pytest.raises(ValidationError):
            StockReserveCommand(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=1, reason="r" * 256)

    def test_adjust_command_zero_quantity_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockAdjustCommand(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=0, reason="noop")

    def test_adjust_command_negative_quantity_passes(self):
        StockAdjustCommand(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=-3, reason="damaged")

    def test_adjust_command_blank_reason_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockAdjustCommand(edition_uid=uuid4(), warehouse_uid=uuid4(), quantity=1, reason=" ")


class TestOrderCommands:
    def test_order_reserve_command_every_field_valid_passes(self):
        StockOrderReserveCommand(order_uid="o-1", edition_id=42, quantity=1)

    def test_order_reserve_command_blank_order_uid_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockOrderReserveCommand(order_uid=" ", edition_id=42, quantity=1)

    def test_order_reserve_command_over_length_order_uid_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockOrderReserveCommand(order_uid="o" * 201, edition_id=42, quantity=1)

    def test_order_release_command_zero_edition_id_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockOrderReleaseCommand(order_uid="o-1", edition_id=0)

    def test_order_reserve_command_zero_quantity_passes_validation(self):
        StockOrderReserveCommand(order_uid="o-1", edition_id=1, quantity=0)


class TestTransferCommand:
    def test_transfer_command_every_field_valid_passes(self):
        StockTransferCommand(
            edition_uid=uuid4(), source_warehouse_uid=uuid4(), target_warehouse_uid=uuid4(), quantity=1
        )

    def test_transfer_command_same_warehouse_raises_validation_error(self):
        warehouse = uuid4()
        with pytest.raises(ValidationError):
            StockTransferCommand(
                edition_uid=uuid4(), source_warehouse_uid=warehouse, target_warehouse_uid=warehouse, quantity=1
            )

    def test_transfer_command_zero_quantity_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockTransferCommand(
                edition_uid=uuid4(), source_warehouse_uid=uuid4(), target_warehouse_uid=uuid4(), quantity=0
            )


class TestWarehouseCommands:
    def test_warehouse_create_command_every_field_valid_passes(self):
        WarehouseCreateCommand(name="Main", capacity=100)

    def test_warehouse_create_command_blank_name_raises_validation_error(self):
        with pytest.raises(ValidationError):
            WarehouseCreateCommand(name=" ", capacity=100)

    def test_warehouse_create_command_over_length_name_raises_validation_error(self):
        with pytest.raises(ValidationError):
            WarehouseCreateCommand(name="n" * 256, capacity=100)

    def test_warehouse_create_command_negative_capacity_raises_validation_error(self):
        with pytest.raises(ValidationError):
            WarehouseCreateCommand(name="Main", capacity=-1)

    def test_warehouse_update_command_every_field_none_passes(self):
        WarehouseUpdateCommand()

    def test_warehouse_update_command_negative_capacity_raises_validation_error(self):
        with pytest.raises(ValidationError):
            WarehouseUpdateCommand(capacity=-5)


class TestAuditCommand:
    def test_audit_create_command_every_field_valid_passes(self):
        StockAuditCreateCommand(
            warehouse_uid=uuid4(), items=[StockAuditItemCommand(edition_uid=uuid4(), actual_quantity=3)]
        )

    def test_audit_create_command_empty_items_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockAuditCreateCommand(warehouse_uid=uuid4(), items=[])

    def test_audit_create_command_negative_actual_quantity_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockAuditCreateCommand(
                warehouse_uid=uuid4(), items=[StockAuditItemCommand(edition_uid=uuid4(), actual_quantity=-1)]
            )

    def test_audit_create_command_duplicate_editions_raises_validation_error(self):
        edition = uuid4()
        with pytest.raises(ValidationError):
            StockAuditCreateCommand(
                warehouse_uid=uuid4(),
                items=[
                    StockAuditItemCommand(edition_uid=edition, actual_quantity=1),
                    StockAuditItemCommand(edition_uid=edition, actual_quantity=2),
                ],
            )

    def test_audit_create_command_over_length_notes_raises_validation_error(self):
        with pytest.raises(ValidationError):
            StockAuditCreateCommand(
                warehouse_uid=uuid4(),
                items=[StockAuditItemCommand(edition_uid=uuid4(), actual_quantity=1)],
                notes="n" * 501,
            )


class TestClosingCommand:
    def test_closing_command_every_field_valid_passes(self):
        MonthlyClosingCommand(warehouse_uid=uuid4(), year=2026, month=8)

    @pytest.mark.parametrize(("year", "month"), [(1999, 8), (2101, 8), (2026, 0), (2026, 13)])
    def test_closing_command_period_out_of_range_raises_validation_error(self, year, month):
        with pytest.raises(ValidationError):
            MonthlyClosingCommand(warehouse_uid=uuid4(), year=year, month=month)
