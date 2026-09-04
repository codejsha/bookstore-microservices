import pytest

from internal.domain.constant.shipment_status import ShipmentStatus
from tests.conftest import make_shipment


@pytest.mark.parametrize(
    "status,expected",
    [
        (ShipmentStatus.PLANNED, True),
        (ShipmentStatus.DISPATCHED, False),
        (ShipmentStatus.PICKED_UP, False),
        (ShipmentStatus.DELIVERED, False),
        (ShipmentStatus.CANCELLED, False),
    ],
)
def test_can_dispatch_per_status_matches_expected_flag(status: ShipmentStatus, expected: bool) -> None:
    shipment = make_shipment(status=status)
    assert shipment.can_dispatch() is expected


@pytest.mark.parametrize(
    "status,expected",
    [
        (ShipmentStatus.PLANNED, False),
        (ShipmentStatus.DISPATCHED, True),
        (ShipmentStatus.PICKED_UP, False),
    ],
)
def test_can_pick_up_per_status_matches_expected_flag(status: ShipmentStatus, expected: bool) -> None:
    shipment = make_shipment(status=status)
    assert shipment.can_pick_up() is expected


@pytest.mark.parametrize(
    "status,expected",
    [
        (ShipmentStatus.IN_TRANSIT, True),
        (ShipmentStatus.OUT_FOR_DELIVERY, True),
        (ShipmentStatus.PICKED_UP, True),
        (ShipmentStatus.DELIVERED, False),
        (ShipmentStatus.PLANNED, False),
    ],
)
def test_can_deliver_per_status_matches_expected_flag(status: ShipmentStatus, expected: bool) -> None:
    shipment = make_shipment(status=status)
    assert shipment.can_deliver() is expected


@pytest.mark.parametrize(
    "status,expected",
    [
        (ShipmentStatus.PLANNED, True),
        (ShipmentStatus.DISPATCHED, True),
        (ShipmentStatus.IN_TRANSIT, False),
        (ShipmentStatus.DELIVERED, False),
        (ShipmentStatus.CANCELLED, False),
    ],
)
def test_can_cancel_per_status_matches_expected_flag(status: ShipmentStatus, expected: bool) -> None:
    shipment = make_shipment(status=status)
    assert shipment.can_cancel() is expected
