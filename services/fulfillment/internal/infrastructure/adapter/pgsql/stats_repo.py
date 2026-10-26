from sqlalchemy import case, func, select, text
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from internal.application.port.repo.repos import StatsRepository
from internal.domain.aggregate.stats_aggregate import (
    CarrierPerformance,
    DeliveryDashboard,
    ShipmentStatusCount,
)
from internal.domain.constant.shipment_status import ShipmentStatus
from internal.domain.model.option.delivery_option import StatsFilterOption
from internal.infrastructure.adapter.mysql.models import (
    CarrierEntity,
    FreightEntity,
    ShipmentEntity,
)
from internal.infrastructure.adapter.mysql.uuid_helper import bytes_to_uuid, uuid_to_bytes


class MySQLStatsRepository(StatsRepository):
    def __init__(self, session_factory: async_sessionmaker[AsyncSession]):
        self._session_factory = session_factory

    async def get_dashboard(self, option: StatsFilterOption) -> DeliveryDashboard:
        async with self._session_factory() as session:
            base_filter = ShipmentEntity.deleted_at.is_(None)
            filters = [base_filter]
            if option.carrier_uid is not None:
                filters.append(ShipmentEntity.carrier_uid == uuid_to_bytes(option.carrier_uid))
            if option.date_from is not None:
                filters.append(ShipmentEntity.created_at >= option.date_from)
            if option.date_to is not None:
                filters.append(ShipmentEntity.created_at <= option.date_to)

            total = (
                await session.execute(select(func.count()).select_from(ShipmentEntity).where(*filters))
            ).scalar() or 0

            status_rows = (
                await session.execute(
                    select(ShipmentEntity.status, func.count()).where(*filters).group_by(ShipmentEntity.status)
                )
            ).all()
            status_counts = [ShipmentStatusCount(status=ShipmentStatus(r[0]), count=r[1]) for r in status_rows]

            avg_hours = (
                await session.execute(
                    select(
                        func.avg(
                            func.timestampdiff(
                                text("HOUR"),
                                ShipmentEntity.actual_pickup_at,
                                ShipmentEntity.actual_delivery_at,
                            )
                        )
                    ).where(
                        *filters,
                        ShipmentEntity.status == ShipmentStatus.DELIVERED,
                        ShipmentEntity.actual_pickup_at.isnot(None),
                        ShipmentEntity.actual_delivery_at.isnot(None),
                    )
                )
            ).scalar()

            delivered_count = (
                await session.execute(
                    select(func.count())
                    .select_from(ShipmentEntity)
                    .where(*filters, ShipmentEntity.status == ShipmentStatus.DELIVERED)
                )
            ).scalar() or 0

            on_time_count = (
                await session.execute(
                    select(func.count())
                    .select_from(ShipmentEntity)
                    .where(
                        *filters,
                        ShipmentEntity.status == ShipmentStatus.DELIVERED,
                        ShipmentEntity.planned_delivery_at.isnot(None),
                        ShipmentEntity.actual_delivery_at <= ShipmentEntity.planned_delivery_at,
                    )
                )
            ).scalar() or 0

            on_time_rate = (on_time_count / delivered_count * 100) if delivered_count > 0 else None

            total_freight = (
                await session.execute(
                    select(func.coalesce(func.sum(FreightEntity.total_cost), 0))
                    .select_from(FreightEntity)
                    .join(ShipmentEntity, ShipmentEntity.uid == FreightEntity.shipment_uid)
                    .where(FreightEntity.deleted_at.is_(None), *filters)
                )
            ).scalar() or 0.0

            freight_by_shipment = (
                select(
                    FreightEntity.shipment_uid.label("shipment_uid"),
                    func.sum(FreightEntity.total_cost).label("freight_cost"),
                )
                .where(FreightEntity.deleted_at.is_(None))
                .group_by(FreightEntity.shipment_uid)
                .subquery()
            )

            carrier_rows = (
                await session.execute(
                    select(
                        CarrierEntity.uid,
                        CarrierEntity.name,
                        func.count(ShipmentEntity.id),
                        func.count(case((ShipmentEntity.status == ShipmentStatus.DELIVERED, ShipmentEntity.id))),
                        func.count(case((ShipmentEntity.status == ShipmentStatus.FAILED, ShipmentEntity.id))),
                        func.avg(
                            case(
                                (
                                    ShipmentEntity.status == ShipmentStatus.DELIVERED,
                                    func.timestampdiff(
                                        text("HOUR"),
                                        ShipmentEntity.actual_pickup_at,
                                        ShipmentEntity.actual_delivery_at,
                                    ),
                                ),
                                else_=None,
                            )
                        ),
                        func.coalesce(func.sum(freight_by_shipment.c.freight_cost), 0),
                    )
                    .select_from(ShipmentEntity)
                    .join(CarrierEntity, CarrierEntity.uid == ShipmentEntity.carrier_uid)
                    .outerjoin(freight_by_shipment, freight_by_shipment.c.shipment_uid == ShipmentEntity.uid)
                    .where(*filters, ShipmentEntity.carrier_uid.isnot(None))
                    .group_by(CarrierEntity.uid, CarrierEntity.name)
                )
            ).all()

            carrier_performances = [
                CarrierPerformance(
                    carrier_uid=str(bytes_to_uuid(r[0])),
                    carrier_name=r[1],
                    total_shipments=r[2],
                    delivered_count=r[3] or 0,
                    failed_count=r[4] or 0,
                    avg_delivery_hours=float(r[5]) if r[5] is not None else None,
                    total_freight_cost=float(r[6]),
                )
                for r in carrier_rows
            ]

            return DeliveryDashboard(
                total_shipments=total,
                status_counts=status_counts,
                carrier_performances=carrier_performances,
                avg_delivery_hours=float(avg_hours) if avg_hours is not None else None,
                total_freight_cost=float(total_freight),
                on_time_delivery_rate=on_time_rate,
            )
