"""init tables

Revision ID: 0001_init_tables
Revises:
Create Date: 2026-04-28

"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import mysql

revision: str = "0001_init_tables"
down_revision: str | None = None
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


_MYSQL_TABLE_OPTS = {
    "mysql_engine": "InnoDB",
    "mysql_charset": "utf8mb4",
    "mysql_collate": "utf8mb4_unicode_ci",
}


def upgrade() -> None:
    op.create_table(
        "carrier",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", mysql.BINARY(16), nullable=False, comment="Public UUID identifier"),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("code", sa.String(50), nullable=False),
        sa.Column("contact_name", sa.String(255), nullable=True),
        sa.Column("contact_phone", sa.String(50), nullable=True),
        sa.Column("contact_email", sa.String(255), nullable=True),
        sa.Column("base_rate", mysql.DOUBLE(), nullable=False),
        sa.Column("rate_per_kg", mysql.DOUBLE(), nullable=False),
        sa.Column("status", sa.String(10), nullable=False, comment="ACTIVE, INACTIVE"),
        sa.Column("created_at", mysql.DATETIME(fsp=6), nullable=False),
        sa.Column("updated_at", mysql.DATETIME(fsp=6), nullable=False),
        sa.Column("deleted_at", mysql.DATETIME(fsp=6), nullable=True),
        sa.CheckConstraint("status IN ('ACTIVE', 'INACTIVE')", name="chk_carrier_status"),
        sa.UniqueConstraint("uid", name="uk_carrier_uid"),
        sa.UniqueConstraint("code", name="uk_carrier_code"),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("idx_carrier_name", "carrier", ["name"])
    op.create_index("idx_carrier_status", "carrier", ["status"])

    op.create_table(
        "shipment",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", mysql.BINARY(16), nullable=False, comment="Public UUID identifier"),
        sa.Column("order_uid", mysql.BINARY(16), nullable=False),
        sa.Column("carrier_uid", mysql.BINARY(16), nullable=True),
        sa.Column("origin_address", sa.String(500), nullable=False),
        sa.Column("destination_address", sa.String(500), nullable=False),
        sa.Column("destination_city", sa.String(100), nullable=False),
        sa.Column("destination_state", sa.String(100), nullable=False),
        sa.Column("destination_country_code", sa.String(3), nullable=False),
        sa.Column("destination_postal_code", sa.String(20), nullable=False),
        sa.Column(
            "status",
            sa.String(20),
            nullable=False,
            comment="PLANNED, DISPATCHED, PICKED_UP, IN_TRANSIT, OUT_FOR_DELIVERY, DELIVERED, FAILED, CANCELLED",
        ),
        sa.Column("tracking_number", sa.String(100), nullable=True),
        sa.Column("weight_kg", mysql.DOUBLE(), nullable=True),
        sa.Column("planned_pickup_at", mysql.DATETIME(fsp=6), nullable=True),
        sa.Column("planned_delivery_at", mysql.DATETIME(fsp=6), nullable=True),
        sa.Column("actual_pickup_at", mysql.DATETIME(fsp=6), nullable=True),
        sa.Column("actual_delivery_at", mysql.DATETIME(fsp=6), nullable=True),
        sa.Column("created_at", mysql.DATETIME(fsp=6), nullable=False),
        sa.Column("updated_at", mysql.DATETIME(fsp=6), nullable=False),
        sa.Column("deleted_at", mysql.DATETIME(fsp=6), nullable=True),
        sa.CheckConstraint(
            "status IN ('PLANNED', 'DISPATCHED', 'PICKED_UP', 'IN_TRANSIT',"
            " 'OUT_FOR_DELIVERY', 'DELIVERED', 'FAILED', 'CANCELLED')",
            name="chk_shipment_status",
        ),
        sa.UniqueConstraint("uid", name="uk_shipment_uid"),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("idx_shipment_order_uid", "shipment", ["order_uid"])
    op.create_index("idx_shipment_carrier_uid", "shipment", ["carrier_uid"])
    op.create_index("idx_shipment_status", "shipment", ["status"])
    op.create_index("idx_shipment_created", "shipment", ["created_at"])

    op.create_table(
        "tracking",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", mysql.BINARY(16), nullable=False, comment="Public UUID identifier"),
        sa.Column("shipment_uid", mysql.BINARY(16), nullable=False),
        sa.Column("status", sa.String(20), nullable=False),
        sa.Column("location", sa.String(255), nullable=False),
        sa.Column("description", sa.Text(), nullable=False),
        sa.Column("occurred_at", mysql.DATETIME(fsp=6), nullable=False),
        sa.Column("created_at", mysql.DATETIME(fsp=6), nullable=False),
        sa.UniqueConstraint("uid", name="uk_tracking_uid"),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("idx_tracking_shipment_uid", "tracking", ["shipment_uid"])
    op.create_index("idx_tracking_occurred", "tracking", ["occurred_at"])

    op.create_table(
        "freight",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", mysql.BINARY(16), nullable=False, comment="Public UUID identifier"),
        sa.Column("shipment_uid", mysql.BINARY(16), nullable=False),
        sa.Column("carrier_uid", mysql.BINARY(16), nullable=False),
        sa.Column("base_cost", mysql.DOUBLE(), nullable=False),
        sa.Column("weight_surcharge", mysql.DOUBLE(), nullable=False),
        sa.Column("distance_surcharge", mysql.DOUBLE(), nullable=False),
        sa.Column("discount", mysql.DOUBLE(), nullable=False),
        sa.Column("total_cost", mysql.DOUBLE(), nullable=False),
        sa.Column("currency", sa.String(3), nullable=False),
        sa.Column("status", sa.String(10), nullable=False, comment="ESTIMATED, CONFIRMED, INVOICED, PAID"),
        sa.Column("invoiced_at", mysql.DATETIME(fsp=6), nullable=True),
        sa.Column("paid_at", mysql.DATETIME(fsp=6), nullable=True),
        sa.Column("created_at", mysql.DATETIME(fsp=6), nullable=False),
        sa.Column("updated_at", mysql.DATETIME(fsp=6), nullable=False),
        sa.Column("deleted_at", mysql.DATETIME(fsp=6), nullable=True),
        sa.CheckConstraint(
            "status IN ('ESTIMATED', 'CONFIRMED', 'INVOICED', 'PAID')",
            name="chk_freight_status",
        ),
        sa.UniqueConstraint("uid", name="uk_freight_uid"),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("idx_freight_shipment_uid", "freight", ["shipment_uid"])
    op.create_index("idx_freight_carrier_uid", "freight", ["carrier_uid"])
    op.create_index("idx_freight_status", "freight", ["status"])


def downgrade() -> None:
    op.drop_table("freight")
    op.drop_table("tracking")
    op.drop_table("shipment")
    op.drop_table("carrier")
