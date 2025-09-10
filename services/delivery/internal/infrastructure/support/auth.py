from __future__ import annotations

import base64
import binascii
import json
from dataclasses import dataclass, field
from typing import Any

from fastapi import Depends, HTTPException, Request

HEADER_USER_ID = "x-user-id"
HEADER_USER_EMAIL = "x-user-email"
HEADER_USER_NAME = "x-user-name"
HEADER_USER_ROLES = "x-user-roles"
HEADER_USER_SCOPES = "x-user-scopes"
HEADER_JWT_PAYLOAD = "x-jwt-payload"


@dataclass
class Principal:
    sub: str
    email: str | None = None
    name: str | None = None
    roles: list[str] = field(default_factory=list)
    scopes: list[str] = field(default_factory=list)
    raw: dict[str, Any] = field(default_factory=dict)

    def has_role(self, role: str) -> bool:
        return role in self.roles

    def has_scope(self, scope: str) -> bool:
        return scope in self.scopes


def _decode_payload(value: str) -> dict[str, Any]:
    padding = "=" * (-len(value) % 4)
    decoded = base64.urlsafe_b64decode(value + padding)
    return json.loads(decoded)


def _parse_claim_list(value: str | None, separator: str | None) -> list[str]:
    if not value:
        return []
    decoded = _decode_claim_list(value)
    if decoded is not None:
        return decoded
    parts = value.split(separator) if separator else value.split()
    return [part.strip() for part in parts if part.strip()]


def _decode_claim_list(value: str) -> list[str] | None:
    padding = "=" * (-len(value) % 4)
    try:
        decoded = json.loads(base64.urlsafe_b64decode(value + padding))
    except binascii.Error, ValueError, UnicodeDecodeError:
        return None
    if not isinstance(decoded, list):
        return None
    return [str(item).strip() for item in decoded if str(item).strip()]


def get_principal(request: Request) -> Principal | None:
    sub = request.headers.get(HEADER_USER_ID)
    if not sub:
        return None
    roles = _parse_claim_list(request.headers.get(HEADER_USER_ROLES), ",")
    scopes = _parse_claim_list(request.headers.get(HEADER_USER_SCOPES), None)
    payload = request.headers.get(HEADER_JWT_PAYLOAD)
    if payload:
        try:
            raw = _decode_payload(payload)
        except (binascii.Error, ValueError, json.JSONDecodeError, UnicodeDecodeError) as exc:
            raise HTTPException(status_code=401, detail="Invalid authentication payload") from exc
    else:
        raw = {}
    return Principal(
        sub=sub,
        email=request.headers.get(HEADER_USER_EMAIL),
        name=request.headers.get(HEADER_USER_NAME),
        roles=roles,
        scopes=scopes,
        raw=raw,
    )


STAFF_ROLES: tuple[str, ...] = ("MANAGE", "SYSTEM")


def require_principal(principal: Principal | None = Depends(get_principal)) -> Principal:
    if principal is None:
        raise HTTPException(status_code=401, detail="Authentication required")
    return principal


def is_staff(principal: Principal) -> bool:
    return any(principal.has_role(role) for role in STAFF_ROLES)


def require_staff(principal: Principal = Depends(require_principal)) -> Principal:
    if not is_staff(principal):
        raise HTTPException(status_code=403, detail="Staff role required")
    return principal
