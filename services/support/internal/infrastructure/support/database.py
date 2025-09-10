import os
from dataclasses import dataclass
from typing import Any

from sqlalchemy import event
from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, async_sessionmaker, create_async_engine

from internal.config.config import DatabaseConfig

_VAULT_ENV_FILE = "/vault/secrets/db.env"


@dataclass
class _CredsCache:
    mtime: float = 0.0
    username: str = ""
    password: str = ""


def _make_do_connect(prefix: str, cache: _CredsCache):
    user_key = f"{prefix}__DATABASE__USERNAME"
    pass_key = f"{prefix}__DATABASE__PASSWORD"

    def do_connect(dialect, conn_rec, cargs: list[Any], cparams: dict[str, Any]):
        if not prefix:
            return None

        try:
            stat = os.stat(_VAULT_ENV_FILE)
        except OSError:
            return None

        if stat.st_mtime != cache.mtime:
            username = ""
            password = ""
            try:
                with open(_VAULT_ENV_FILE, encoding="utf-8") as fh:
                    for raw_line in fh:
                        line = raw_line.strip()
                        if "=" not in line or line.startswith("#"):
                            continue
                        key, _, value = line.partition("=")
                        if key == user_key:
                            username = value
                        elif key == pass_key:
                            password = value
            except OSError:
                return None

            if username and password:
                cache.mtime = stat.st_mtime
                cache.username = username
                cache.password = password
            else:
                return None

        if cache.username and cache.password:
            cparams["user"] = cache.username
            cparams["password"] = cache.password

        return None

    return do_connect


def create_engine_and_session_factory(
    config: DatabaseConfig,
) -> tuple[AsyncEngine, async_sessionmaker[AsyncSession]]:
    engine = create_async_engine(
        config.url,
        pool_size=10,
        max_overflow=20,
        pool_timeout=30,
        pool_recycle=1800,
        pool_pre_ping=True,
        connect_args={"connect_timeout": 10},
        echo=False,
    )

    _cache: _CredsCache = _CredsCache()
    event.listen(
        engine.sync_engine,
        "do_connect",
        _make_do_connect(config.vault_env_prefix, _cache),
    )

    return engine, async_sessionmaker(bind=engine, expire_on_commit=False)


def create_session_factory(config: DatabaseConfig) -> async_sessionmaker[AsyncSession]:
    _, factory = create_engine_and_session_factory(config)
    return factory
