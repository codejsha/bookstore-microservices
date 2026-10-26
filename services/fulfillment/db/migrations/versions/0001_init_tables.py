"""init tables

Revision ID: 0001_init_tables
Revises:
Create Date: 2026-10-26

"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import op

from internal.infrastructure.adapter.pgsql.models import Amount, Timestamp, Uid

revision: str = "0001_init_tables"
down_revision: str | None = None
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.create_table(
        "carrier",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("code", sa.String(50), nullable=False),
        sa.Column("contact_name", sa.String(255), nullable=True),
        sa.Column("contact_phone", sa.String(50), nullable=True),
        sa.Column("contact_email", sa.String(255), nullable=True),
        sa.Column("base_rate", Amount, nullable=False),
        sa.Column("rate_per_kg", Amount, nullable=False),
        sa.Column("status", sa.String(10), nullable=False),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=False),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.UniqueConstraint("uid", name="uk_carrier_uid"),
        sa.UniqueConstraint("code", name="uk_carrier_code"),
    )
    op.create_index("idx_carrier_name", "carrier", ["name"])
    op.create_index("idx_carrier_status", "carrier", ["status"])

    op.create_table(
        "shipment",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("order_uid", Uid, nullable=False),
        sa.Column("carrier_uid", Uid, nullable=True),
        sa.Column("origin_address", sa.String(500), nullable=False),
        sa.Column("destination_address", sa.String(500), nullable=False),
        sa.Column("destination_city", sa.String(100), nullable=False),
        sa.Column("destination_state", sa.String(100), nullable=False),
        sa.Column("destination_country_code", sa.String(3), nullable=False),
        sa.Column("destination_postal_code", sa.String(20), nullable=False),
        sa.Column("status", sa.String(20), nullable=False),
        sa.Column("tracking_number", sa.String(100), nullable=True),
        sa.Column("weight_kg", Amount, nullable=True),
        sa.Column("planned_pickup_at", Timestamp, nullable=True),
        sa.Column("planned_delivery_at", Timestamp, nullable=True),
        sa.Column("actual_pickup_at", Timestamp, nullable=True),
        sa.Column("actual_delivery_at", Timestamp, nullable=True),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=False),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.UniqueConstraint("uid", name="uk_shipment_uid"),
        sa.UniqueConstraint("order_uid", name="uk_shipment_order_uid"),
    )
    op.create_index("idx_shipment_carrier_uid", "shipment", ["carrier_uid"])
    op.create_index("idx_shipment_status", "shipment", ["status"])
    op.create_index("idx_shipment_created", "shipment", ["created_at"])

    op.create_table(
        "tracking",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("shipment_uid", Uid, nullable=False),
        sa.Column("status", sa.String(20), nullable=False),
        sa.Column("location", sa.String(255), nullable=False),
        sa.Column("description", sa.Text(), nullable=False),
        sa.Column("occurred_at", Timestamp, nullable=False),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.UniqueConstraint("uid", name="uk_tracking_uid"),
    )
    op.create_index("idx_tracking_shipment_uid", "tracking", ["shipment_uid"])
    op.create_index("idx_tracking_occurred", "tracking", ["occurred_at"])

    op.create_table(
        "freight",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("shipment_uid", Uid, nullable=False),
        sa.Column("carrier_uid", Uid, nullable=False),
        sa.Column("base_cost", Amount, nullable=False),
        sa.Column("weight_surcharge", Amount, nullable=False),
        sa.Column("distance_surcharge", Amount, nullable=False),
        sa.Column("discount", Amount, nullable=False),
        sa.Column("total_cost", Amount, nullable=False),
        sa.Column("currency", sa.String(3), nullable=False),
        sa.Column("status", sa.String(10), nullable=False),
        sa.Column("invoiced_at", Timestamp, nullable=True),
        sa.Column("paid_at", Timestamp, nullable=True),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=False),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.UniqueConstraint("uid", name="uk_freight_uid"),
    )
    op.create_index("idx_freight_shipment_uid", "freight", ["shipment_uid"])
    op.create_index("idx_freight_carrier_uid", "freight", ["carrier_uid"])
    op.create_index("idx_freight_status", "freight", ["status"])

    op.create_table(
        "warehouse",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("name", sa.String(255), nullable=False),
        sa.Column("address", sa.String(255), nullable=True),
        sa.Column("capacity", sa.Integer(), nullable=False),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=True),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("uid", name="uk_warehouse_uid"),
    )
    op.create_index("idx_warehouse_name", "warehouse", ["name"])

    op.create_table(
        "stock",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("edition_id", sa.BigInteger(), nullable=False),
        sa.Column("edition_uid", Uid, nullable=False),
        sa.Column("warehouse_id", sa.BigInteger(), nullable=False),
        sa.Column("warehouse_uid", Uid, nullable=False),
        sa.Column("quantity", sa.Integer(), nullable=False),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=True),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("uid", name="uk_stock_uid"),
        sa.UniqueConstraint("edition_id", "warehouse_id", name="uq_stock_edition_warehouse"),
    )
    op.create_index("idx_stock_edition_id", "stock", ["edition_id"])
    op.create_index("idx_stock_warehouse_id", "stock", ["warehouse_id"])
    op.create_index("idx_stock_edition_uid_warehouse_uid", "stock", ["edition_uid", "warehouse_uid"])

    op.create_table(
        "stock_history",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("stock_id", sa.BigInteger(), nullable=False),
        sa.Column("stock_uid", Uid, nullable=False),
        sa.Column("change_type", sa.String(50), nullable=False),
        sa.Column("reason", sa.String(255), nullable=True),
        sa.Column("change_qty", sa.Integer(), nullable=False),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=True),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("uid", name="uk_stock_history_uid"),
    )
    op.create_index("idx_stock_history_stock_id_created", "stock_history", ["stock_id", "created_at"])
    op.create_index("idx_stock_history_stock_uid_created", "stock_history", ["stock_uid", "created_at"])

    op.create_table(
        "stock_reservation",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("reservation_key", sa.String(255), nullable=False),
        sa.Column("order_ref", sa.String(255), nullable=False),
        sa.Column("edition_id", sa.BigInteger(), nullable=False),
        sa.Column("edition_uid", Uid, nullable=False),
        sa.Column("warehouse_id", sa.BigInteger(), nullable=False),
        sa.Column("warehouse_uid", Uid, nullable=False),
        sa.Column("quantity", sa.Integer(), nullable=False),
        sa.Column("status", sa.String(20), nullable=False),
        sa.Column("released_at", Timestamp, nullable=True),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=True),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("uid", name="uk_stock_reservation_uid"),
        sa.UniqueConstraint("reservation_key", "warehouse_id", name="uk_stock_reservation_key_warehouse"),
    )
    op.create_index("idx_stock_reservation_order_ref", "stock_reservation", ["order_ref"])
    op.create_index("idx_stock_reservation_edition_id", "stock_reservation", ["edition_id"])
    op.create_index("idx_stock_reservation_warehouse_id", "stock_reservation", ["warehouse_id"])

    op.create_table(
        "stock_transfer",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("edition_id", sa.BigInteger(), nullable=False),
        sa.Column("edition_uid", Uid, nullable=False),
        sa.Column("source_warehouse_id", sa.BigInteger(), nullable=False),
        sa.Column("source_warehouse_uid", Uid, nullable=False),
        sa.Column("target_warehouse_id", sa.BigInteger(), nullable=False),
        sa.Column("target_warehouse_uid", Uid, nullable=False),
        sa.Column("quantity", sa.Integer(), nullable=False),
        sa.Column("status", sa.String(20), nullable=False),
        sa.Column("reason", sa.String(255), nullable=True),
        sa.Column("completed_at", Timestamp, nullable=True),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=True),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("uid", name="uk_transfer_uid"),
    )
    op.create_index("idx_transfer_edition_id", "stock_transfer", ["edition_id"])
    op.create_index("idx_transfer_source", "stock_transfer", ["source_warehouse_id"])
    op.create_index("idx_transfer_target", "stock_transfer", ["target_warehouse_id"])
    op.create_index("idx_transfer_status", "stock_transfer", ["status"])

    op.create_table(
        "stock_audit",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("warehouse_id", sa.BigInteger(), nullable=False),
        sa.Column("warehouse_uid", Uid, nullable=False),
        sa.Column("status", sa.String(20), nullable=False),
        sa.Column("notes", sa.String(500), nullable=True),
        sa.Column("completed_at", Timestamp, nullable=True),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=True),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("uid", name="uk_audit_uid"),
    )
    op.create_index("idx_audit_warehouse_id", "stock_audit", ["warehouse_id"])
    op.create_index("idx_audit_status", "stock_audit", ["status"])

    op.create_table(
        "stock_audit_item",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("audit_id", sa.BigInteger(), nullable=False),
        sa.Column("edition_id", sa.BigInteger(), nullable=False),
        sa.Column("edition_uid", Uid, nullable=False),
        sa.Column("system_quantity", sa.Integer(), nullable=False),
        sa.Column("actual_quantity", sa.Integer(), nullable=False),
        sa.Column("difference", sa.Integer(), nullable=False),
    )
    op.create_index("idx_audit_item_audit_id", "stock_audit_item", ["audit_id"])

    op.create_table(
        "monthly_closing",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False),
        sa.Column("warehouse_id", sa.BigInteger(), nullable=False),
        sa.Column("warehouse_uid", Uid, nullable=False),
        sa.Column("year", sa.Integer(), nullable=False),
        sa.Column("month", sa.Integer(), nullable=False),
        sa.Column("status", sa.String(20), nullable=False),
        sa.Column("closed_at", Timestamp, nullable=True),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=True),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("uid", name="uk_closing_uid"),
        sa.UniqueConstraint("warehouse_id", "year", "month", name="uq_closing_warehouse_period"),
    )
    op.create_index("idx_closing_warehouse_id", "monthly_closing", ["warehouse_id"])

    op.create_table(
        "monthly_closing_item",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("closing_id", sa.BigInteger(), nullable=False),
        sa.Column("edition_id", sa.BigInteger(), nullable=False),
        sa.Column("edition_uid", Uid, nullable=False),
        sa.Column("opening_quantity", sa.Integer(), nullable=False),
        sa.Column("inbound_quantity", sa.Integer(), nullable=False),
        sa.Column("outbound_quantity", sa.Integer(), nullable=False),
        sa.Column("adjust_quantity", sa.Integer(), nullable=False),
        sa.Column("closing_quantity", sa.Integer(), nullable=False),
    )
    op.create_index("idx_closing_item_closing_id", "monthly_closing_item", ["closing_id"])


def downgrade() -> None:
    op.drop_table("monthly_closing_item")
    op.drop_table("monthly_closing")
    op.drop_table("stock_audit_item")
    op.drop_table("stock_audit")
    op.drop_table("stock_transfer")
    op.drop_table("stock_reservation")
    op.drop_table("stock_history")
    op.drop_table("stock")
    op.drop_table("warehouse")
    op.drop_table("freight")
    op.drop_table("tracking")
    op.drop_table("shipment")
    op.drop_table("carrier")
