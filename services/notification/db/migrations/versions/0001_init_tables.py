"""init tables

Baseline for the notification schema. This transcribes the three legacy
golang-migrate SQL files (000001_init_tables, 000002_alter_datetime_precision,
000003_notification_order_dedup) that this service carried before it moved to
Alembic; those files were never run by any Alembic-based pipeline, so their DDL
had never actually been applied and is collapsed into this single baseline.

The column types are imported from the SQLAlchemy models rather than restated,
so `alembic revision --autogenerate` against a database built from this revision
produces an empty diff.

Revision ID: 0001_init_tables
Revises:
Create Date: 2026-07-13

"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import op

from internal.infrastructure.adapter.mysql.models import (
    ChannelCol,
    NotificationStatusCol,
    NotificationTypeCol,
    Timestamp,
    Uid,
)

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
        "notification",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False, comment="Public UUID identifier"),
        sa.Column("user_uid", Uid, nullable=False),
        sa.Column("order_uid", Uid, nullable=True),
        sa.Column("notification_type", NotificationTypeCol, nullable=False),
        sa.Column("channel", ChannelCol, nullable=False),
        sa.Column("status", NotificationStatusCol, nullable=False),
        sa.Column("title", sa.String(255), nullable=False),
        sa.Column("content", sa.Text(), nullable=False),
        sa.Column("sent_at", Timestamp, nullable=True),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=False),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.UniqueConstraint("order_uid", "notification_type", name="uq_notification_order_type"),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("ix_notification_uid", "notification", ["uid"], unique=True)
    op.create_index("ix_notification_user_uid", "notification", ["user_uid"])
    op.create_index("ix_notification_order_uid", "notification", ["order_uid"])

    op.create_table(
        "template",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False, comment="Public UUID identifier"),
        sa.Column("notification_type", NotificationTypeCol, nullable=False),
        sa.Column("channel", ChannelCol, nullable=False),
        sa.Column("title_template", sa.String(255), nullable=False),
        sa.Column("content_template", sa.Text(), nullable=False),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=False),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.UniqueConstraint("notification_type", "channel", name="uq_template_type_channel"),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("ix_template_uid", "template", ["uid"], unique=True)


def downgrade() -> None:
    op.drop_index("ix_template_uid", table_name="template")
    op.drop_table("template")
    op.drop_index("ix_notification_order_uid", table_name="notification")
    op.drop_index("ix_notification_user_uid", table_name="notification")
    op.drop_index("ix_notification_uid", table_name="notification")
    op.drop_table("notification")
