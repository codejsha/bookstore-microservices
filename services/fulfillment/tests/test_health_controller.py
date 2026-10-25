from fastapi import FastAPI
from fastapi.testclient import TestClient

from internal.infrastructure.adapter.restcontroller.health_controller import ComponentCheck, create_health_router


def _client(ping_result: bool, components: dict[str, ComponentCheck] | None = None) -> TestClient:
    async def ping() -> bool:
        return ping_result

    app = FastAPI()
    app.include_router(create_health_router(ping, components))
    return TestClient(app)


def test_ready_ping_succeeds_returns_200() -> None:
    response = _client(True).get("/health/ready")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_ready_ping_fails_returns_503() -> None:
    response = _client(False).get("/health/ready")
    assert response.status_code == 503
    assert response.json() == {"status": "unavailable", "components": ["database"]}


def test_health_always_returns_200() -> None:
    for ping_result in (True, False):
        response = _client(ping_result).get("/health")
        assert response.status_code == 200
        assert response.json() == {"status": "ok"}


def test_health_component_down_returns_200() -> None:
    response = _client(True, {"worker": lambda: False}).get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_ready_all_components_up_returns_200() -> None:
    components: dict[str, ComponentCheck] = {"worker": lambda: True, "server": lambda: True}
    response = _client(True, components).get("/health/ready")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_ready_no_components_configured_returns_200() -> None:
    response = _client(True, {}).get("/health/ready")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_ready_single_component_down_returns_503() -> None:
    components: dict[str, ComponentCheck] = {"worker": lambda: False, "server": lambda: True}
    response = _client(True, components).get("/health/ready")
    assert response.status_code == 503
    assert response.json() == {"status": "unavailable", "components": ["worker"]}


def test_ready_every_component_down_reports_all() -> None:
    components: dict[str, ComponentCheck] = {"worker": lambda: False, "server": lambda: False}
    response = _client(True, components).get("/health/ready")
    assert response.status_code == 503
    assert response.json() == {"status": "unavailable", "components": ["server", "worker"]}


def test_ready_component_down_skips_database_ping() -> None:
    pinged = False

    async def ping() -> bool:
        nonlocal pinged
        pinged = True
        return True

    app = FastAPI()
    app.include_router(create_health_router(ping, {"worker": lambda: False}))

    response = TestClient(app).get("/health/ready")

    assert response.status_code == 503
    assert pinged is False
