import os
from pathlib import Path

import pytest

from internal.config.config import DatabaseConfig
from internal.infrastructure.support import database

ENV_PREFIX = "NOTIFICATION"


def _write_env_file(path: Path, username: str, password: str) -> None:
    path.write_text(
        f"{ENV_PREFIX}__DATABASE__USERNAME={username}\n{ENV_PREFIX}__DATABASE__PASSWORD={password}\n",
        encoding="utf-8",
    )


class TestDoConnect:
    def test_do_connect_file_present_overrides_credentials(
        self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        env_file = tmp_path / "db.env"
        _write_env_file(env_file, "vault-user", "vault-pass")
        monkeypatch.setattr(database, "_VAULT_ENV_FILE", str(env_file))

        do_connect = database._make_do_connect(ENV_PREFIX, database._CredsCache())
        cparams = {"user": "static-user", "password": "static-pass"}

        assert do_connect(None, None, [], cparams) is None
        assert cparams == {"user": "vault-user", "password": "vault-pass"}

    def test_do_connect_file_rewritten_picks_up_new_password(
        self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        env_file = tmp_path / "db.env"
        _write_env_file(env_file, "vault-user", "old-pass")
        monkeypatch.setattr(database, "_VAULT_ENV_FILE", str(env_file))

        do_connect = database._make_do_connect(ENV_PREFIX, database._CredsCache())
        first = {}
        do_connect(None, None, [], first)

        _write_env_file(env_file, "vault-user", "new-pass")
        stat = os.stat(env_file)
        os.utime(env_file, (stat.st_atime + 10, stat.st_mtime + 10))

        second = {}
        do_connect(None, None, [], second)

        assert first == {"user": "vault-user", "password": "old-pass"}
        assert second == {"user": "vault-user", "password": "new-pass"}

    def test_do_connect_file_missing_leaves_params(self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch) -> None:
        monkeypatch.setattr(database, "_VAULT_ENV_FILE", str(tmp_path / "absent.env"))

        do_connect = database._make_do_connect(ENV_PREFIX, database._CredsCache())
        cparams = {"user": "static-user", "password": "static-pass"}

        assert do_connect(None, None, [], cparams) is None
        assert cparams == {"user": "static-user", "password": "static-pass"}

    def test_do_connect_prefix_empty_is_noop(self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch) -> None:
        env_file = tmp_path / "db.env"
        _write_env_file(env_file, "vault-user", "vault-pass")
        monkeypatch.setattr(database, "_VAULT_ENV_FILE", str(env_file))

        do_connect = database._make_do_connect("", database._CredsCache())
        cparams = {"user": "static-user", "password": "static-pass"}

        assert do_connect(None, None, [], cparams) is None
        assert cparams == {"user": "static-user", "password": "static-pass"}


class TestCreateSessionFactory:
    def test_create_session_factory_empty_prefix_raises(self) -> None:
        with pytest.raises(ValueError):
            database.create_session_factory(DatabaseConfig(), "")


class TestCreateReadinessEngine:
    def test_create_readiness_engine_empty_prefix_raises(self) -> None:
        with pytest.raises(ValueError):
            database.create_readiness_engine(DatabaseConfig(), "")

    async def test_create_readiness_engine_pool_smaller_than_application_pool(self) -> None:
        app_engine = database.create_session_factory(DatabaseConfig(), ENV_PREFIX).kw["bind"]
        readiness_engine = database.create_readiness_engine(DatabaseConfig(), ENV_PREFIX)
        try:
            assert database._READINESS_POOL_SIZE < database._APP_POOL_SIZE
            assert database._READINESS_MAX_OVERFLOW < database._APP_MAX_OVERFLOW
            assert readiness_engine.pool.size() == database._READINESS_POOL_SIZE
            assert readiness_engine.pool.size() < app_engine.pool.size()
        finally:
            await readiness_engine.dispose()
            await app_engine.dispose()

    async def test_create_readiness_engine_vault_file_present_overrides_credentials(
        self, tmp_path: Path, monkeypatch: pytest.MonkeyPatch
    ) -> None:
        env_file = tmp_path / "db.env"
        _write_env_file(env_file, "vault-user", "vault-pass")
        monkeypatch.setattr(database, "_VAULT_ENV_FILE", str(env_file))

        engine = database.create_readiness_engine(DatabaseConfig(), ENV_PREFIX)
        try:
            listeners = list(engine.sync_engine.dialect.dispatch.do_connect)
            assert len(listeners) == 1

            cparams = {"user": "static-user", "password": "static-pass"}
            listeners[0](None, None, [], cparams)

            assert cparams == {"user": "vault-user", "password": "vault-pass"}
        finally:
            await engine.dispose()
