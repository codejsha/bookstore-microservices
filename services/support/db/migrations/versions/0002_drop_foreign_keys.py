"""drop foreign key constraints

Revision ID: 0002_drop_foreign_keys
Revises: 0001_init_tables
Create Date: 2026-08-25

"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import op

revision: str = "0002_drop_foreign_keys"
down_revision: str | None = "0001_init_tables"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


_FOREIGN_KEYS = (
    {
        "table": "ticket",
        "constraint": "fk_ticket__category",
        "column": "category_id",
        "referred_table": "ticket_category",
        "ondelete": "SET NULL",
        "index": "ix_ticket_category_id",
        "index_added_here": True,
    },
    {
        "table": "ticket_category",
        "constraint": "fk_ticket_category__parent",
        "column": "parent_id",
        "referred_table": "ticket_category",
        "ondelete": "SET NULL",
        "index": "ix_ticket_category_parent_id",
        "index_added_here": True,
    },
    {
        "table": "ticket_comment",
        "constraint": "fk_ticket_comment__ticket",
        "column": "ticket_id",
        "referred_table": "ticket",
        "ondelete": "CASCADE",
        "index": "ix_ticket_comment_ticket_id",
        "index_added_here": False,
    },
    {
        "table": "faq",
        "constraint": "fk_faq__category",
        "column": "category_id",
        "referred_table": "ticket_category",
        "ondelete": "SET NULL",
        "index": "ix_faq_category_id",
        "index_added_here": False,
    },
)


def _foreign_key_names(table: str) -> set[str]:
    inspector = sa.inspect(op.get_bind())
    return {fk["name"] for fk in inspector.get_foreign_keys(table) if fk["name"]}


def _index_names(table: str) -> set[str]:
    inspector = sa.inspect(op.get_bind())
    return {index["name"] for index in inspector.get_indexes(table) if index["name"]}


def upgrade() -> None:
    for fk in _FOREIGN_KEYS:
        table, constraint, index = fk["table"], fk["constraint"], fk["index"]
        if constraint in _foreign_key_names(table):
            op.drop_constraint(constraint, table, type_="foreignkey")
        if index not in _index_names(table):
            op.create_index(index, table, [fk["column"]])
        if constraint in _index_names(table):
            op.drop_index(constraint, table_name=table)


def downgrade() -> None:
    for fk in reversed(_FOREIGN_KEYS):
        table, constraint, index = fk["table"], fk["constraint"], fk["index"]
        if fk["index_added_here"] and index in _index_names(table):
            op.drop_index(index, table_name=table)
        if constraint not in _foreign_key_names(table):
            op.create_foreign_key(
                constraint,
                table,
                fk["referred_table"],
                [fk["column"]],
                ["id"],
                ondelete=fk["ondelete"],
            )
