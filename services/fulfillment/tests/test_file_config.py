from pathlib import Path

import pytest
import yaml

from internal.config import file_config as fc
from internal.config.config import DatabaseConfig, ServerConfig, Settings, vault_env_prefix

ENV_PREFIX = "FULFILLMENT"

CONFIG_YAML = """\
app:
  logging_level: debug
server:
  port: 8080
  mode: prod
database:
  host: cfghost
  port: 9999
  db_name: fulfillment_db
temporal:
  host: temporal.local:7233
  auth:
    enabled: true
    token_url: ${TEST_CFG_TOKEN_URL}
    client_secret: ${TEST_CFG_SECRET}
"""


def _write(tmp_path: Path, content: str) -> Path:
    path = tmp_path / "application.yaml"
    path.write_text(content, encoding="utf-8")
    return path


@pytest.fixture(autouse=True)
def _isolated_env(monkeypatch: pytest.MonkeyPatch) -> None:
    for var in (
        fc.CONFIG_FILE_ENV,
        f"{ENV_PREFIX}__DATABASE__HOST",
        f"{ENV_PREFIX}__SERVER__MODE",
        "TEST_CFG_TOKEN_URL",
        "TEST_CFG_SECRET",
    ):
        monkeypatch.delenv(var, raising=False)


class TestConfigFilePath:
    def test_config_file_path_env_unset_uses_mounted_default(self) -> None:
        assert fc.config_file_path() == Path("/app/config/application.yaml")

    def test_config_file_path_env_blank_uses_mounted_default(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv(fc.CONFIG_FILE_ENV, "   ")
        assert fc.config_file_path() == Path("/app/config/application.yaml")

    def test_config_file_path_env_set_uses_env_path(self, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv(fc.CONFIG_FILE_ENV, "/tmp/custom.yaml")
        assert fc.config_file_path() == Path("/tmp/custom.yaml")


class TestLoadConfigFile:
    def test_load_config_file_absent_file_is_empty(self, tmp_path: Path) -> None:
        assert fc.load_config_file(tmp_path / "missing.yaml") == {}

    def test_load_config_file_empty_file_is_empty(self, tmp_path: Path) -> None:
        assert fc.load_config_file(_write(tmp_path, "")) == {}

    def test_load_config_file_nested_yaml_keeps_structure(self, tmp_path: Path) -> None:
        data = fc.load_config_file(_write(tmp_path, "server:\n  port: 8080\nitems:\n  - a\n  - b\n"))
        assert data == {"server": {"port": 8080}, "items": ["a", "b"]}

    def test_load_config_file_env_placeholders_expands_from_environment(
        self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        secret = "p@ss: #w$rd 'x\""
        monkeypatch.setenv("TEST_CFG_TOKEN_URL", "http://idp/token")
        monkeypatch.setenv("TEST_CFG_SECRET", secret)

        data = fc.load_config_file(_write(tmp_path, CONFIG_YAML))

        assert data["temporal"]["auth"]["token_url"] == "http://idp/token"
        assert data["temporal"]["auth"]["client_secret"] == secret

    def test_load_config_file_placeholder_in_list_expands_each_item(
        self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        monkeypatch.setenv("TEST_CFG_SECRET", "broker:9093")
        data = fc.load_config_file(_write(tmp_path, "brokers:\n  - ${TEST_CFG_SECRET}\n  - plain\n"))
        assert data["brokers"] == ["broker:9093", "plain"]

    def test_load_config_file_scalar_top_level_raises(self, tmp_path: Path) -> None:
        with pytest.raises(fc.ConfigFileInvalidError):
            fc.load_config_file(_write(tmp_path, "- a\n- b\n"))

    def test_load_config_file_malformed_yaml_raises(self, tmp_path: Path) -> None:
        with pytest.raises(yaml.YAMLError):
            fc.load_config_file(_write(tmp_path, "server: [unclosed\n"))


class TestSettingsPrecedence:
    def test_settings_env_and_file_both_set_prefers_env(self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv(fc.CONFIG_FILE_ENV, str(_write(tmp_path, CONFIG_YAML)))
        monkeypatch.setenv(f"{ENV_PREFIX}__DATABASE__HOST", "envhost")

        settings = Settings(_env_file=None)

        assert settings.database.host == "envhost"
        assert settings.database.port == 9999
        assert settings.server.mode == "prod"
        assert settings.app.logging_level == "debug"
        assert settings.temporal.host == "temporal.local:7233"
        assert settings.temporal.auth.enabled is True
        assert settings.database.username == DatabaseConfig().username

    def test_settings_vault_env_file_overrides_file_values(
        self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        monkeypatch.setenv(fc.CONFIG_FILE_ENV, str(_write(tmp_path, CONFIG_YAML)))
        db_env = tmp_path / "db.env"
        db_env.write_text(
            f"{ENV_PREFIX}__DATABASE__USERNAME=vault-user\n{ENV_PREFIX}__DATABASE__PASSWORD=vault-pass\n",
            encoding="utf-8",
        )

        settings = Settings(_env_file=(str(db_env), str(tmp_path / "temporal.env")))

        assert settings.database.username == "vault-user"
        assert settings.database.password == "vault-pass"
        assert settings.database.host == "cfghost"

    def test_settings_absent_file_uses_defaults(self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setenv(fc.CONFIG_FILE_ENV, str(tmp_path / "absent.yaml"))

        settings = Settings(_env_file=None)

        assert settings.server.mode == ServerConfig().mode
        assert settings.database.host == DatabaseConfig().host

    def test_settings_database_from_init_keeps_vault_prefix_available(
        self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        monkeypatch.setenv(fc.CONFIG_FILE_ENV, str(tmp_path / "absent.yaml"))

        settings = Settings(_env_file=None, database={"host": "db", "port": 5432, "db_name": "x"})

        assert settings.database.host == "db"
        assert vault_env_prefix() == ENV_PREFIX
