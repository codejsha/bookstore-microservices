"""shipment order_uid unique

Enforce one shipment per order at the database level. This backstops the
idempotency of the at-least-once ``CreateShipment`` Temporal activity: even if
two retries race past the application-level get-or-create lookup, the unique
constraint rejects the duplicate insert and the activity falls back to the
existing shipment.

Revision ID: 0002_shipment_order_uid_unique
Revises: 0001_init_tables
Create Date: 2026-07-02

"""

from collections.abc import Sequence

from alembic import op

revision: str = "0002_shipment_order_uid_unique"
down_revision: str | None = "0001_init_tables"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.drop_index("idx_shipment_order_uid", table_name="shipment")
    op.create_unique_constraint("uk_shipment_order_uid", "shipment", ["order_uid"])


def downgrade() -> None:
    op.drop_constraint("uk_shipment_order_uid", "shipment", type_="unique")
    op.create_index("idx_shipment_order_uid", "shipment", ["order_uid"])
