class ShipmentStatusCount:
    def __init__(self, status: str, count: int):
        self.status = status
        self.count = count


class CarrierPerformance:
    def __init__(
        self,
        carrier_uid: str,
        carrier_name: str,
        total_shipments: int,
        delivered_count: int,
        failed_count: int,
        avg_delivery_hours: float | None,
        total_freight_cost: float,
    ):
        self.carrier_uid = carrier_uid
        self.carrier_name = carrier_name
        self.total_shipments = total_shipments
        self.delivered_count = delivered_count
        self.failed_count = failed_count
        self.avg_delivery_hours = avg_delivery_hours
        self.total_freight_cost = total_freight_cost


class DeliveryDashboard:
    def __init__(
        self,
        total_shipments: int,
        status_counts: list[ShipmentStatusCount],
        carrier_performances: list[CarrierPerformance],
        avg_delivery_hours: float | None,
        total_freight_cost: float,
        on_time_delivery_rate: float | None,
    ):
        self.total_shipments = total_shipments
        self.status_counts = status_counts
        self.carrier_performances = carrier_performances
        self.avg_delivery_hours = avg_delivery_hours
        self.total_freight_cost = total_freight_cost
        self.on_time_delivery_rate = on_time_delivery_rate
