from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.infrastructure.adapter.restcontroller.health_controller import create_health_router


def _client(ping_result: bool) -> TestClient:
    async def ping() -> bool:
        return ping_result

    app = FastAPI()
    app.include_router(create_health_router(ping))
    return TestClient(app)


def test_ready_ping_succeeds_returns_200() -> None:
    response = _client(True).get("/health/ready")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_ready_ping_fails_returns_503() -> None:
    response = _client(False).get("/health/ready")
    assert response.status_code == 503
    assert response.json() == {"status": "unavailable"}


def test_health_always_returns_200() -> None:
    for ping_result in (True, False):
        response = _client(ping_result).get("/health")
        assert response.status_code == 200
        assert response.json() == {"status": "ok"}
