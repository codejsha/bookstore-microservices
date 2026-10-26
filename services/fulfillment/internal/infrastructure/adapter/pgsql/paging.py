from collections.abc import Mapping, Sequence

from sqlalchemy import ColumnElement
from sqlalchemy.sql.expression import UnaryExpression

DEFAULT_SORT: tuple[tuple[str, str], ...] = (("created_at", "desc"), ("uid", "asc"))


def parse_sort(sort: str | None) -> list[tuple[str, str]]:
    items: list[tuple[str, str]] = []
    for raw in (sort or "").split(";"):
        item = raw.strip()
        if not item:
            continue
        separator = "," if "," in item else ":"
        field, _, direction = item.partition(separator)
        field = field.strip().lower()
        if not field:
            continue
        direction = direction.strip().lower()
        items.append((field, "desc" if direction == "desc" else "asc"))
    return items


def sortable_by_wire_name(columns: Mapping[str, ColumnElement]) -> dict[str, ColumnElement]:
    out: dict[str, ColumnElement] = {}
    for name, column in columns.items():
        key = name.lower()
        out[key] = column
        out[key.replace("_", "")] = column
    return out


def resolve_order_by(
    sort: str | None,
    sortable: Mapping[str, ColumnElement],
    default: Sequence[tuple[str, str]] = DEFAULT_SORT,
) -> list[UnaryExpression]:
    normalized = sortable_by_wire_name(sortable)
    kept = [(field, direction) for field, direction in parse_sort(sort) if field in normalized]
    if not kept:
        kept = [(field, direction) for field, direction in default if field in normalized]
    return [normalized[field].desc() if direction == "desc" else normalized[field].asc() for field, direction in kept]


def page_offset(page: int, size: int) -> int:
    return max(page, 0) * size
