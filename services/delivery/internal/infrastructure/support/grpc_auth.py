from __future__ import annotations

from collections.abc import Iterable
from uuid import UUID

from internal.infrastructure.support.auth import (
    HEADER_JWT_PAYLOAD,
    HEADER_USER_ID,
    _decode_payload,
)

Metadata = Iterable[tuple[str, str | bytes]]


def _as_str(value: str | bytes) -> str:
    return value.decode() if isinstance(value, bytes) else value


def caller_user_uid(metadata: Metadata | None) -> UUID | None:
    md: dict[str, str] = {}
    for key, value in metadata or ():
        md[key.lower()] = _as_str(value)

    sub = md.get(HEADER_USER_ID)
    if not sub:
        payload = md.get(HEADER_JWT_PAYLOAD)
        if payload:
            try:
                sub = _decode_payload(payload).get("sub")
            except ValueError, TypeError:
                sub = None
    if not sub:
        return None
    try:
        return UUID(sub)
    except ValueError, TypeError:
        return None
