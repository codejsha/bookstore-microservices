from __future__ import annotations

import asyncio
import os
from collections.abc import Iterator
from pathlib import Path
from urllib.parse import urlsplit

import pytest
from alembic import command
from alembic.autogenerate import compare_metadata
from alembic.config import Config
from alembic.migration import MigrationContext
from sqlalchemy import text
from sqlalchemy.ext.asyncio import create_async_engine

from internal.infrastructure.adapter.pgsql.models import Base
from tests.postgres_container import container_runtime_available, postgres_container_url

os.environ.setdefault("TESTCONTAINERS_RYUK_DISABLED", "true")

pytestmark = pytest.mark.skipif(
    not container_runtime_available(),
    reason="no reachable container runtime (Docker/Podman) for testcontainers",
)

SERVICE_ROOT = Path(__file__).resolve().parent.parent
UID_COLUMNS = {"uid", "order_uid", "carrier_uid", "shipment_uid", "edition_uid", "warehouse_uid", "stock_uid"}
AMOUNT_COLUMNS = {"base_rate", "rate_per_kg", "weight_kg", "base_cost", "total_cost", "discount"}


@pytest.fixture(scope="module")
def postgres_url() -> Iterator[str]:
    yield from postgres_container_url()


@pytest.fixture
def migrated_url(postgres_url: str, monkeypatch: pytest.MonkeyPatch) -> Iterator[str]:
    parts = urlsplit(postgres_url)
    monkeypatch.setenv("FULFILLMENT__DATABASE__HOST", parts.hostname or "")
    monkeypatch.setenv("FULFILLMENT__DATABASE__PORT", str(parts.port))
    monkeypatch.setenv("FULFILLMENT__DATABASE__USERNAME", parts.username or "")
    monkeypatch.setenv("FULFILLMENT__DATABASE__PASSWORD", parts.password or "")
    monkeypatch.setenv("FULFILLMENT__DATABASE__DB_NAME", parts.path.lstrip("/"))
    config = Config(str(SERVICE_ROOT / "alembic.ini"))
    config.set_main_option("script_location", str(SERVICE_ROOT / "db" / "migrations"))
    command.upgrade(config, "head")
    try:
        yield postgres_url
    finally:
        command.downgrade(config, "base")


async def _columns(url: str) -> list[tuple[str, str, str, int | None]]:
    engine = create_async_engine(url)
    try:
        async with engine.connect() as conn:
            rows = await conn.execute(
                text(
                    "SELECT table_name, column_name, data_type, datetime_precision "
                    "FROM information_schema.columns WHERE table_schema = 'public'"
                )
            )
            return [tuple(row) for row in rows.all()]
    finally:
        await engine.dispose()


async def _metadata_diff(url: str) -> list:
    engine = create_async_engine(url)
    try:
        async with engine.connect() as conn:
            return await conn.run_sync(
                lambda sync_conn: compare_metadata(MigrationContext.configure(sync_conn), Base.metadata)
            )
    finally:
        await engine.dispose()


def test_upgrade_head_creates_every_mapped_table(migrated_url: str) -> None:
    tables = {row[0] for row in asyncio.run(_columns(migrated_url))}

    assert tables == set(Base.metadata.tables) | {"alembic_version"}


def test_upgrade_head_uid_columns_native_uuid(migrated_url: str) -> None:
    types = {row[2] for row in asyncio.run(_columns(migrated_url)) if row[1] in UID_COLUMNS}

    assert types == {"uuid"}


def test_upgrade_head_time_columns_timestamp_6_without_time_zone(migrated_url: str) -> None:
    time_columns = [row for row in asyncio.run(_columns(migrated_url)) if row[1].endswith("_at")]

    assert time_columns
    assert {(row[2], row[3]) for row in time_columns} == {("timestamp without time zone", 6)}


def test_upgrade_head_amount_columns_double_precision(migrated_url: str) -> None:
    types = {row[2] for row in asyncio.run(_columns(migrated_url)) if row[1] in AMOUNT_COLUMNS}

    assert types == {"double precision"}


def test_upgrade_head_id_columns_bigint(migrated_url: str) -> None:
    types = {row[2] for row in asyncio.run(_columns(migrated_url)) if row[1] == "id"}

    assert types == {"bigint"}


def test_upgrade_head_schema_matches_orm_metadata(migrated_url: str) -> None:
    assert asyncio.run(_metadata_diff(migrated_url)) == []
