import base64
import json
from uuid import uuid4

import pytest
from fastapi import Depends, FastAPI
from fastapi.testclient import TestClient

from internal.infrastructure.support.auth import Principal, require_principal
from internal.infrastructure.support.grpc_auth import caller_user_uid


def _b64(raw: bytes) -> str:
    return base64.urlsafe_b64encode(raw).decode().rstrip("=")


def _app() -> FastAPI:
    app = FastAPI()

    @app.get("/whoami")
    def whoami(principal: Principal = Depends(require_principal)) -> dict:
        return {"sub": principal.sub, "raw": principal.raw}

    return app


# ─── REST: x-jwt-payload decode guard ────────────────────────────────────────


@pytest.mark.parametrize(
    "payload",
    [
        _b64(b"not json at all"),
        "!!!!",
        _b64(b"\xff\xfe\xff"),
    ],
)
def test_current_user_malformed_jwt_payload_unauthorized(payload: str) -> None:
    client = TestClient(_app())
    resp = client.get("/whoami", headers={"x-user-id": str(uuid4()), "x-jwt-payload": payload})
    assert resp.status_code == 401


def test_current_user_valid_jwt_payload_decodes_claims() -> None:
    sub = str(uuid4())
    payload = _b64(json.dumps({"sub": sub, "email": "a@b.com"}).encode())
    client = TestClient(_app())
    resp = client.get("/whoami", headers={"x-user-id": sub, "x-jwt-payload": payload})
    assert resp.status_code == 200
    assert resp.json()["raw"]["email"] == "a@b.com"


def test_current_user_absent_jwt_payload_has_empty_raw() -> None:
    sub = str(uuid4())
    client = TestClient(_app())
    resp = client.get("/whoami", headers={"x-user-id": sub})
    assert resp.status_code == 200
    assert resp.json()["raw"] == {}


# ─── gRPC: metadata parser fails closed on malformed payload ──────────────────


@pytest.mark.parametrize(
    "payload",
    [
        _b64(b"not json at all"),
        "!!!!",
        _b64(b"\xff\xfe\xff"),
    ],
)
def test_caller_user_uid_malformed_jwt_payload_is_none(payload: str) -> None:
    assert caller_user_uid((("x-jwt-payload", payload),)) is None


def test_caller_user_uid_valid_jwt_payload_resolves_sub() -> None:
    uid = uuid4()
    payload = _b64(json.dumps({"sub": str(uid)}).encode())
    assert caller_user_uid((("x-jwt-payload", payload),)) == uid


def test_caller_user_uid_with_x_user_id_prefers_it() -> None:
    uid = uuid4()
    assert caller_user_uid((("x-user-id", str(uid)),)) == uid
