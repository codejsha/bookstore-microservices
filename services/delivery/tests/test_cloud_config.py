import pytest

from internal.config import cloud_config as cc
from internal.config.config import (
    CLOUD_CONFIG_APP_NAME,
    DatabaseConfig,
    ServerConfig,
    Settings,
)

ENV_PREFIX = "DELIVERY"


class TestExpandSourceMap:
    def test_nested_dotted_keys(self) -> None:
        flat = {
            "server.port": 8080,
            "server.mode": "prod",
            "database.host": "db",
            "database.port": 3306,
            "telemetry.enabled": True,
        }
        assert cc._expand_source_map(flat) == {
            "server": {"port": 8080, "mode": "prod"},
            "database": {"host": "db", "port": 3306},
            "telemetry": {"enabled": True},
        }

    def test_array_index_notation(self) -> None:
        flat = {
            "servers[0].host": "a",
            "servers[1].host": "b",
            "tags[0]": "x",
            "tags[1]": "y",
        }
        result = cc._expand_source_map(flat)
        assert result == {
            "servers": [{"host": "a"}, {"host": "b"}],
            "tags": ["x", "y"],
        }

    def test_malformed_bracket_treated_as_literal(self) -> None:
        assert cc._expand_source_map({"weird[abc": 1}) == {"weird[abc": 1}


class TestMergePropertySources:
    def test_earlier_source_wins(self) -> None:
        sources = [
            {"name": "prod", "source": {"server.mode": "prod", "database.host": "prod-db"}},
            {"name": "base", "source": {"server.mode": "base", "server.port": 8080}},
        ]
        merged = cc._merge_property_sources(sources)
        assert merged == {
            "server.mode": "prod",
            "database.host": "prod-db",
            "server.port": 8080,
        }


class TestUrlBuilding:
    def test_strip_configserver_prefix(self) -> None:
        assert cc._strip_configserver_prefix("configserver:http://host:8888") == "http://host:8888"
        assert cc._strip_configserver_prefix("http://host:8888") == "http://host:8888"

    def test_build_url_without_label_lowercases(self) -> None:
        url = cc._build_url("http://host:8888/", "Support", "Dev", "")
        assert url == "http://host:8888/support/dev"

    def test_build_url_with_label(self) -> None:
        url = cc._build_url("http://host:8888", "support", "prod", "MAIN")
        assert url == "http://host:8888/support/prod/main"


class TestFetchCloudConfig:
    def test_no_server_returns_empty(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.delenv("APP_CONFIG_SERVER", raising=False)
        assert cc.fetch_cloud_config(CLOUD_CONFIG_APP_NAME) == {}

    def test_empty_server_returns_empty(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv("APP_CONFIG_SERVER", "   ")
        assert cc.fetch_cloud_config(CLOUD_CONFIG_APP_NAME) == {}

    def test_failure_raises_after_retries(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv("APP_CONFIG_SERVER", "configserver:http://unreachable:8888")
        monkeypatch.setenv("APP_CONFIG_PROFILE", "dev")
        monkeypatch.delenv("APP_CONFIG_LABEL", raising=False)

        calls = {"n": 0}

        def boom(*_args: object, **_kwargs: object) -> None:
            calls["n"] += 1
            raise OSError("connection refused")

        monkeypatch.setattr(cc.urllib.request, "urlopen", boom)
        monkeypatch.setattr(cc.time, "sleep", lambda _s: None)

        with pytest.raises(cc.CloudConfigUnavailableError):
            cc.fetch_cloud_config(CLOUD_CONFIG_APP_NAME)
        assert calls["n"] == cc._MAX_ATTEMPTS

    def test_success_parses_and_expands(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv("APP_CONFIG_SERVER", "configserver:http://host:8888")
        monkeypatch.setenv("APP_CONFIG_PROFILE", "prod")
        monkeypatch.delenv("APP_CONFIG_LABEL", raising=False)

        captured = {}
        payload = {
            "propertySources": [
                {"name": "prod", "source": {"server.mode": "prod", "database.host": "prod-db"}},
                {"name": "base", "source": {"server.mode": "base", "server.port": 8080}},
            ]
        }

        class _Resp:
            def __enter__(self) -> _Resp:
                return self

            def __exit__(self, *_a: object) -> None:
                return None

            def read(self) -> bytes:
                import json

                return json.dumps(payload).encode()

        def fake_urlopen(url: str, timeout: float = 0) -> _Resp:
            captured["url"] = url
            return _Resp()

        monkeypatch.setattr(cc.urllib.request, "urlopen", fake_urlopen)

        result = cc.fetch_cloud_config(CLOUD_CONFIG_APP_NAME)
        assert captured["url"] == f"http://host:8888/{CLOUD_CONFIG_APP_NAME}/prod"
        assert result == {
            "server": {"mode": "prod", "port": 8080},
            "database": {"host": "prod-db"},
        }


class TestSettingsPrecedence:
    def _clear_env(self, monkeypatch: pytest.MonkeyPatch) -> None:
        for var in ("APP_CONFIG_SERVER", "APP_CONFIG_PROFILE", "APP_CONFIG_LABEL"):
            monkeypatch.delenv(var, raising=False)
        monkeypatch.delenv(f"{ENV_PREFIX}__DATABASE__HOST", raising=False)
        monkeypatch.delenv(f"{ENV_PREFIX}__SERVER__MODE", raising=False)

    def test_cloud_values_applied_and_env_overrides(self, monkeypatch: pytest.MonkeyPatch) -> None:
        self._clear_env(monkeypatch)

        monkeypatch.setattr(
            cc,
            "fetch_cloud_config",
            lambda _app: {
                "server": {"mode": "prod"},
                "database": {"host": "cloudhost", "port": 9999},
            },
        )
        monkeypatch.setenv(f"{ENV_PREFIX}__DATABASE__HOST", "envhost")

        settings = Settings()

        assert settings.database.host == "envhost"
        assert settings.database.port == 9999
        assert settings.server.mode == "prod"
        assert settings.database.db_name == DatabaseConfig().db_name

    def test_defaults_when_cloud_empty(self, monkeypatch: pytest.MonkeyPatch) -> None:
        self._clear_env(monkeypatch)
        monkeypatch.setattr(cc, "fetch_cloud_config", lambda _app: {})

        settings = Settings()

        assert settings.server.mode == ServerConfig().mode
        assert settings.database.host == DatabaseConfig().host
