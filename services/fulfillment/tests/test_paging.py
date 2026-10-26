from internal.infrastructure.adapter.mysql.models import WarehouseEntity
from internal.infrastructure.adapter.mysql.paging import page_offset, parse_sort, resolve_order_by

_SORTABLE = {"name": WarehouseEntity.name, "created_at": WarehouseEntity.created_at, "uid": WarehouseEntity.uid}


def test_parse_sort_comma_format_splits_field_and_direction() -> None:
    assert parse_sort("name,desc") == [("name", "desc")]


def test_parse_sort_colon_format_splits_field_and_direction() -> None:
    assert parse_sort("name:desc") == [("name", "desc")]


def test_parse_sort_multiple_items_keeps_order() -> None:
    assert parse_sort("created_at,desc;uid,asc") == [("created_at", "desc"), ("uid", "asc")]


def test_parse_sort_missing_direction_defaults_to_asc() -> None:
    assert parse_sort("name") == [("name", "asc")]


def test_parse_sort_unknown_direction_falls_back_to_asc() -> None:
    assert parse_sort("name,sideways") == [("name", "asc")]


def test_parse_sort_blank_is_empty() -> None:
    assert parse_sort("") == []
    assert parse_sort(None) == []


def test_resolve_order_by_unknown_field_falls_back_to_default() -> None:
    clauses = resolve_order_by("metadata,desc", _SORTABLE)
    assert [str(c) for c in clauses] == ["warehouse.created_at DESC", "warehouse.uid ASC"]


def test_resolve_order_by_known_field_without_underscore_matches() -> None:
    clauses = resolve_order_by("createdAt,desc", _SORTABLE)
    assert [str(c) for c in clauses] == ["warehouse.created_at DESC"]


def test_resolve_order_by_mixed_known_and_unknown_keeps_known_only() -> None:
    clauses = resolve_order_by("bogus,asc;name,asc", _SORTABLE)
    assert [str(c) for c in clauses] == ["warehouse.name ASC"]


def test_page_offset_multiplies_page_by_size() -> None:
    assert page_offset(2, 20) == 40


def test_page_offset_negative_page_clamps_to_zero() -> None:
    assert page_offset(-1, 20) == 0
