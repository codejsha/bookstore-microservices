import argparse
import asyncio
import json
import sys
from dataclasses import dataclass
from pathlib import Path
from uuid import UUID

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from sqlalchemy import select  # noqa: E402
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker  # noqa: E402

from internal.config.config import Settings, vault_env_prefix  # noqa: E402
from internal.domain.constant.stock_change_type import StockChangeType  # noqa: E402
from internal.infrastructure.adapter.pgsql.models import WarehouseEntity  # noqa: E402
from internal.infrastructure.adapter.pgsql.stock_movement import StockDelta, apply_stock_delta, lock_stock  # noqa: E402
from internal.infrastructure.support.database import create_engine_and_session_factory  # noqa: E402

SEED_REASON = "seed"


@dataclass(frozen=True)
class SeedItem:
    edition_id: int
    edition_uid: UUID
    warehouse_uid: UUID
    quantity: int


@dataclass
class SeedSummary:
    created: int = 0
    skipped: int = 0


def load_items_from_file(path: Path) -> list[SeedItem]:
    raw = json.loads(path.read_text(encoding="utf-8"))
    return [
        SeedItem(
            edition_id=int(entry["edition_id"]),
            edition_uid=UUID(str(entry["edition_uid"])),
            warehouse_uid=UUID(str(entry["warehouse_uid"])),
            quantity=int(entry["quantity"]),
        )
        for entry in raw
    ]


def load_items_from_catalog(dsn: str, warehouse_uid: UUID, quantity: int) -> list[SeedItem]:
    try:
        import psycopg
    except ImportError as exc:
        raise SystemExit("psycopg is required for --catalog-dsn: run `uv run --with psycopg[binary] ...`") from exc
    with psycopg.connect(dsn) as conn, conn.cursor() as cursor:
        cursor.execute("SELECT id, uid FROM edition WHERE deleted_at IS NULL ORDER BY id")
        rows = cursor.fetchall()
    return [
        SeedItem(edition_id=int(row[0]), edition_uid=UUID(str(row[1])), warehouse_uid=warehouse_uid, quantity=quantity)
        for row in rows
    ]


async def seed_stock(session_factory: async_sessionmaker[AsyncSession], items: list[SeedItem]) -> SeedSummary:
    summary = SeedSummary()
    for item in items:
        async with session_factory() as session, session.begin():
            warehouse = (
                await session.execute(
                    select(WarehouseEntity).where(
                        WarehouseEntity.uid == item.warehouse_uid,
                        WarehouseEntity.deleted_at.is_(None),
                    )
                )
            ).scalar_one_or_none()
            if warehouse is None:
                raise SystemExit(f"warehouse {item.warehouse_uid} not found; register it first")
            if await lock_stock(session, item.edition_uid, warehouse.uid) is not None:
                summary.skipped += 1
                continue
            await apply_stock_delta(
                session,
                StockDelta(
                    edition_id=item.edition_id,
                    edition_uid=item.edition_uid,
                    warehouse_id=warehouse.id,
                    warehouse_uid=warehouse.uid,
                    delta=item.quantity,
                    change_type=StockChangeType.INBOUND,
                    reason=SEED_REASON,
                    allow_create=True,
                ),
            )
            summary.created += 1
    return summary


def parse_args(argv: list[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Seed initial stock rows (idempotent: existing rows are skipped)")
    source = parser.add_mutually_exclusive_group(required=True)
    source.add_argument("--input", type=Path, help="JSON list of {edition_id, edition_uid, warehouse_uid, quantity}")
    source.add_argument("--catalog-dsn", help="PostgreSQL DSN of the catalog database to read editions from")
    parser.add_argument("--warehouse-uid", type=UUID, help="target warehouse when reading from the catalog")
    parser.add_argument("--quantity", type=int, default=100, help="quantity per edition when reading from the catalog")
    args = parser.parse_args(argv)
    if args.catalog_dsn and args.warehouse_uid is None:
        parser.error("--warehouse-uid is required with --catalog-dsn")
    if args.quantity <= 0:
        parser.error("--quantity must be positive")
    return args


async def run(args: argparse.Namespace) -> SeedSummary:
    items = (
        load_items_from_file(args.input)
        if args.input
        else load_items_from_catalog(args.catalog_dsn, args.warehouse_uid, args.quantity)
    )
    settings = Settings()
    engine, session_factory = create_engine_and_session_factory(settings.database, vault_env_prefix())
    try:
        return await seed_stock(session_factory, items)
    finally:
        await engine.dispose()


def main(argv: list[str] | None = None) -> None:
    args = parse_args(sys.argv[1:] if argv is None else argv)
    summary = asyncio.run(run(args))
    print(f"seeded stock rows: created={summary.created} skipped={summary.skipped}")


if __name__ == "__main__":
    main()
