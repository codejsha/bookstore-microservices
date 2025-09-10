"""init tables

Baseline for the support schema. This transcribes the legacy golang-migrate SQL
file (000001_init_tables) that this service carried before it moved to Alembic;
that file was never run by any Alembic-based pipeline, so its DDL had never
actually been applied.

Two corrections against that legacy DDL, both required for the service to work:

  - The user-identifying columns (`ticket.customer_id`, `ticket.assignee_id`,
    `ticket_comment.author_id`) were BIGINT. Platform-wide, a user is identified
    across service boundaries by their UUID -- the JWT `sub` -- and a numeric id
    is internal to whichever service owns the row. Support has no way to resolve
    a UUID subject to another service's numeric id, so ownership checks could
    never succeed. They are BINARY(16) uid columns here.
  - The `DEFAULT` clauses on `ticket_comment.internal`, `faq.view_count` and
    `faq.published` are dropped: the repo standard is that the application sets
    every value (the models carry the same defaults Python-side).

The column types are imported from the SQLAlchemy models rather than restated,
so `alembic revision --autogenerate` against a database built from this revision
produces an empty diff.

Revision ID: 0001_init_tables
Revises:
Create Date: 2026-07-13

"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import context, op

from internal.infrastructure.adapter.mysql.models import (
    AuthorRoleCol,
    TicketPriorityCol,
    TicketStatusCol,
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
        "ticket_category",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False, comment="Public UUID identifier"),
        sa.Column("name", sa.String(100), nullable=False, unique=True),
        sa.Column("description", sa.String(500), nullable=True),
        sa.Column("parent_id", sa.BigInteger(), nullable=True),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=False),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("actor", sa.BigInteger(), nullable=False),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.ForeignKeyConstraint(
            ["parent_id"], ["ticket_category.id"], name="fk_ticket_category__parent", ondelete="SET NULL"
        ),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("ix_ticket_category_uid", "ticket_category", ["uid"], unique=True)

    op.create_table(
        "ticket",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False, comment="Public UUID identifier"),
        sa.Column("customer_uid", Uid, nullable=False, comment="Owning user (JWT sub)"),
        sa.Column("category_id", sa.BigInteger(), nullable=True),
        sa.Column("assignee_uid", Uid, nullable=True, comment="Assigned agent (JWT sub)"),
        sa.Column("subject", sa.String(255), nullable=False),
        sa.Column("description", sa.Text(), nullable=False),
        sa.Column("status", TicketStatusCol, nullable=False),
        sa.Column("priority", TicketPriorityCol, nullable=False),
        sa.Column("resolved_at", Timestamp, nullable=True),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=False),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("actor", sa.BigInteger(), nullable=False),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.ForeignKeyConstraint(
            ["category_id"], ["ticket_category.id"], name="fk_ticket__category", ondelete="SET NULL"
        ),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("ix_ticket_uid", "ticket", ["uid"], unique=True)
    op.create_index("ix_ticket_customer_uid", "ticket", ["customer_uid"])
    op.create_index("ix_ticket_assignee_uid", "ticket", ["assignee_uid"])
    op.create_index("ix_ticket_status", "ticket", ["status"])
    op.create_index("ix_ticket_priority", "ticket", ["priority"])
    op.create_index("ix_ticket_created_at", "ticket", ["created_at"])

    op.create_table(
        "ticket_comment",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False, comment="Public UUID identifier"),
        sa.Column("ticket_id", sa.BigInteger(), nullable=False),
        sa.Column("author_uid", Uid, nullable=False, comment="Comment author (JWT sub)"),
        sa.Column("author_role", AuthorRoleCol, nullable=False),
        sa.Column("body", sa.Text(), nullable=False),
        sa.Column("internal", sa.Boolean(), nullable=False),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=False),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("actor", sa.BigInteger(), nullable=False),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.ForeignKeyConstraint(["ticket_id"], ["ticket.id"], name="fk_ticket_comment__ticket", ondelete="CASCADE"),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("ix_ticket_comment_uid", "ticket_comment", ["uid"], unique=True)
    op.create_index("ix_ticket_comment_ticket_id", "ticket_comment", ["ticket_id"])
    op.create_index("ix_ticket_comment_created_at", "ticket_comment", ["created_at"])

    op.create_table(
        "faq",
        sa.Column("id", sa.BigInteger(), autoincrement=True, primary_key=True),
        sa.Column("uid", Uid, nullable=False, comment="Public UUID identifier"),
        sa.Column("category_id", sa.BigInteger(), nullable=True),
        sa.Column("question", sa.String(500), nullable=False),
        sa.Column("answer", sa.Text(), nullable=False),
        sa.Column("view_count", sa.BigInteger(), nullable=False),
        sa.Column("published", sa.Boolean(), nullable=False),
        sa.Column("created_at", Timestamp, nullable=False),
        sa.Column("updated_at", Timestamp, nullable=False),
        sa.Column("deleted_at", Timestamp, nullable=True),
        sa.Column("actor", sa.BigInteger(), nullable=False),
        sa.Column("version", sa.BigInteger(), nullable=False),
        sa.ForeignKeyConstraint(["category_id"], ["ticket_category.id"], name="fk_faq__category", ondelete="SET NULL"),
        **_MYSQL_TABLE_OPTS,
    )
    op.create_index("ix_faq_uid", "faq", ["uid"], unique=True)
    op.create_index("ix_faq_category_id", "faq", ["category_id"])
    op.create_index("ix_faq_published", "faq", ["published"])
    if context.get_context().dialect.name == "mysql":
        op.create_index(
            "idx_faq_question_answer",
            "faq",
            ["question", "answer"],
            mysql_prefix="FULLTEXT",
        )


def downgrade() -> None:
    op.drop_table("faq")
    op.drop_table("ticket_comment")
    op.drop_table("ticket")
    op.drop_table("ticket_category")
