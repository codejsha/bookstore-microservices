"""drop check constraints

Revision ID: 0003_drop_check_constraints
Revises: 0002_shipment_order_uid_unique
Create Date: 2026-08-25

"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import op

revision: str = "0003_drop_check_constraints"
down_revision: str | None = "0002_shipment_order_uid_unique"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.drop_constraint("chk_carrier_status", "carrier", type_="check")
    op.drop_constraint("chk_shipment_status", "shipment", type_="check")
    op.drop_constraint("chk_freight_status", "freight", type_="check")


def downgrade() -> None:
    op.create_check_constraint(
        "chk_freight_status",
        "freight",
        sa.text("status IN ('ESTIMATED', 'CONFIRMED', 'INVOICED', 'PAID')"),
    )
    op.create_check_constraint(
        "chk_shipment_status",
        "shipment",
        sa.text(
            "status IN ('PLANNED', 'DISPATCHED', 'PICKED_UP', 'IN_TRANSIT',"
            " 'OUT_FOR_DELIVERY', 'DELIVERED', 'FAILED', 'CANCELLED')"
        ),
    )
    op.create_check_constraint(
        "chk_carrier_status",
        "carrier",
        sa.text("status IN ('ACTIVE', 'INACTIVE')"),
    )
