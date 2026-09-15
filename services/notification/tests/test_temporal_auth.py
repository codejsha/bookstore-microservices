import asyncio
import json
import threading
import time
import urllib.error
import urllib.parse
from collections.abc import Iterator
from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path
from types import SimpleNamespace

import pytest

from internal.config.config import Settings, TemporalAuthConfig
from internal.infrastructure.adapter.temporal.temporal_auth import (
    AccessToken,
    KeycloakTokenFetcher,
    TemporalTokenError,
    TemporalTokenProvider,
    build_temporal_token_provider,
)

ENV_PREFIX = Settings.model_config["env_prefix"]


class FakeClock:
    def __init__(self, now: float = 1000.0):
        self.now = now

    def __call__(self) -> float:
        return self.now


class ScriptedFetcher:
    def __init__(self, *steps: AccessToken | Exception):
        self._steps = list(steps)
        self.calls = 0

    async def __call__(self) -> AccessToken:
        self.calls += 1
        step = self._steps.pop(0) if len(self._steps) > 1 else self._steps[0]
        if isinstance(step, Exception):
            raise step
        return step


def _provider(fetcher: ScriptedFetcher, clock: FakeClock) -> TemporalTokenProvider:
    return TemporalTokenProvider(fetcher, refresh_ratio=0.75, min_backoff=1.0, max_backoff=8.0, clock=clock)


def _unavailable() -> Exception:
    return ConnectionError("token endpoint unavailable")


class TestTemporalTokenProvider:
    async def test_initialize_300s_lifetime_schedules_refresh_at_75_percent(self) -> None:
        provider = _provider(ScriptedFetcher(AccessToken("first", 300)), FakeClock())

        await provider.initialize()

        assert provider.next_refresh_delay() == pytest.approx(225.0)

    async def test_initialize_fetch_fails_propagates_error(self) -> None:
        provider = _provider(ScriptedFetcher(_unavailable()), FakeClock())

        with pytest.raises(ConnectionError):
            await provider.initialize()

    def test_rpc_metadata_before_initialize_raises(self) -> None:
        provider = _provider(ScriptedFetcher(AccessToken("first", 300)), FakeClock())

        with pytest.raises(TemporalTokenError):
            provider.rpc_metadata()

    async def test_rpc_metadata_repeated_calls_reuses_cached_token(self) -> None:
        fetcher = ScriptedFetcher(AccessToken("first", 300))
        provider = _provider(fetcher, FakeClock())
        await provider.initialize()

        for _ in range(5):
            assert provider.rpc_metadata() == {"authorization": "Bearer first"}

        assert fetcher.calls == 1

    async def test_refresh_fetch_succeeds_replaces_token_and_reschedules(self) -> None:
        clock = FakeClock()
        provider = _provider(ScriptedFetcher(AccessToken("first", 300), AccessToken("second", 120)), clock)
        await provider.initialize()
        clock.now += 225

        delay = await provider.refresh()

        assert provider.rpc_metadata() == {"authorization": "Bearer second"}
        assert delay == pytest.approx(90.0)

    async def test_refresh_fetch_fails_keeps_previous_token(self) -> None:
        clock = FakeClock()
        provider = _provider(ScriptedFetcher(AccessToken("first", 300), _unavailable()), clock)
        await provider.initialize()
        clock.now += 225

        delay = await provider.refresh()

        assert provider.rpc_metadata() == {"authorization": "Bearer first"}
        assert delay == 1.0

    async def test_refresh_consecutive_failures_backs_off_up_to_max(self) -> None:
        provider = _provider(ScriptedFetcher(AccessToken("first", 300), _unavailable()), FakeClock())
        await provider.initialize()

        delays = [await provider.refresh() for _ in range(5)]

        assert delays == [1.0, 2.0, 4.0, 8.0, 8.0]

    async def test_refresh_success_after_failures_resets_backoff(self) -> None:
        provider = _provider(
            ScriptedFetcher(
                AccessToken("first", 300), _unavailable(), _unavailable(), AccessToken("second", 300), _unavailable()
            ),
            FakeClock(),
        )
        await provider.initialize()
        for _ in range(3):
            await provider.refresh()

        delay = await provider.refresh()

        assert delay == 1.0
        assert provider.rpc_metadata() == {"authorization": "Bearer second"}

    async def test_run_refresh_due_notifies_listener_with_new_metadata(self) -> None:
        provider = TemporalTokenProvider(
            ScriptedFetcher(AccessToken("first", 0.04), AccessToken("second", 300)), clock=time.monotonic
        )
        await provider.initialize()
        received: list[dict[str, str]] = []
        notified = asyncio.Event()

        def on_refresh(metadata) -> None:
            received.append(dict(metadata))
            notified.set()

        task = asyncio.create_task(provider.run(on_refresh))
        try:
            await asyncio.wait_for(notified.wait(), timeout=2.0)
        finally:
            task.cancel()
            with pytest.raises(asyncio.CancelledError):
                await task

        assert received == [{"authorization": "Bearer second"}]


@pytest.fixture
def token_server() -> Iterator[SimpleNamespace]:
    state = SimpleNamespace(status=200, body={"access_token": "abc.def.ghi", "expires_in": 300}, form=None)

    class Handler(BaseHTTPRequestHandler):
        def do_POST(self) -> None:
            length = int(self.headers["Content-Length"])
            state.form = urllib.parse.parse_qs(self.rfile.read(length).decode("utf-8"))
            payload = json.dumps(state.body).encode("utf-8")
            self.send_response(state.status)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(payload)))
            self.end_headers()
            self.wfile.write(payload)

        def log_message(self, *args) -> None:
            pass

    server = HTTPServer(("127.0.0.1", 0), Handler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    state.url = f"http://127.0.0.1:{server.server_port}/realms/bookstore/protocol/openid-connect/token"
    try:
        yield state
    finally:
        server.shutdown()
        server.server_close()


class TestKeycloakTokenFetcher:
    async def test_fetch_client_credentials_grant_parses_token(self, token_server: SimpleNamespace) -> None:
        fetcher = KeycloakTokenFetcher(token_server.url, "temporal-worker", "s3cret", timeout_seconds=2.0)

        token = await fetcher()

        assert token == AccessToken("abc.def.ghi", 300.0)
        assert token_server.form == {
            "grant_type": ["client_credentials"],
            "client_id": ["temporal-worker"],
            "client_secret": ["s3cret"],
        }

    async def test_fetch_error_status_raises_http_error(self, token_server: SimpleNamespace) -> None:
        token_server.status = 401
        token_server.body = {"error": "unauthorized_client"}
        fetcher = KeycloakTokenFetcher(token_server.url, "temporal-worker", "wrong", timeout_seconds=2.0)

        with pytest.raises(urllib.error.HTTPError):
            await fetcher()

    async def test_fetch_response_without_access_token_raises(self, token_server: SimpleNamespace) -> None:
        token_server.body = {"expires_in": 300}
        fetcher = KeycloakTokenFetcher(token_server.url, "temporal-worker", "s3cret", timeout_seconds=2.0)

        with pytest.raises(TemporalTokenError):
            await fetcher()


class TestBuildTemporalTokenProvider:
    def test_build_auth_disabled_is_none(self) -> None:
        assert build_temporal_token_provider(TemporalAuthConfig()) is None

    def test_build_missing_client_secret_raises(self) -> None:
        config = TemporalAuthConfig(enabled=True, token_url="http://keycloak", client_id="temporal-worker")

        with pytest.raises(ValueError, match="client_secret"):
            build_temporal_token_provider(config)

    def test_build_complete_config_creates_provider(self) -> None:
        config = TemporalAuthConfig(
            enabled=True, token_url="http://keycloak", client_id="temporal-worker", client_secret="s3cret"
        )

        assert isinstance(build_temporal_token_provider(config), TemporalTokenProvider)


class TestTemporalAuthSettings:
    def test_settings_default_env_files_include_vault_temporal_env(self) -> None:
        assert "/vault/secrets/temporal.env" in Settings.model_config["env_file"]

    def test_settings_temporal_env_file_sets_auth_credentials(
        self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        monkeypatch.delenv("APP_CONFIG_SERVER", raising=False)
        temporal_env = tmp_path / "temporal.env"
        temporal_env.write_text(
            f"{ENV_PREFIX}TEMPORAL__AUTH__CLIENT_ID=temporal-worker\n{ENV_PREFIX}TEMPORAL__AUTH__CLIENT_SECRET=s3cret\n",
            encoding="utf-8",
        )

        settings = Settings(_env_file=(str(tmp_path / "db.env"), str(temporal_env)))

        assert settings.temporal.auth.client_id == "temporal-worker"
        assert settings.temporal.auth.client_secret == "s3cret"
        assert settings.temporal.auth.enabled is False
