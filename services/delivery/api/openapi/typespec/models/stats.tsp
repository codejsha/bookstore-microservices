namespace DeliveryService;

// ─── Stats models ───────────────────────────────────────────────────────────

model ShipmentStatusCount {
  status: string;
  count: int64;
}

model CarrierPerformance {
  carrier_uid: uuid;
  carrier_name: string;
  total_shipments: int64;
  delivered_count: int64;
  failed_count: int64;
  avg_delivery_hours?: float64;
  total_freight_cost: float64;
}

model DashboardResponse {
  total_shipments: int64;

  @minItems(0)
  status_counts: ShipmentStatusCount[];

  @minItems(0)
  carrier_performances: CarrierPerformance[];

  avg_delivery_hours?: float64;
  total_freight_cost: float64;
  on_time_delivery_rate?: float64;
}
