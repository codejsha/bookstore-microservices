namespace DeliveryService;

// ─── Freight enums ──────────────────────────────────────────────────────────

enum FreightStatus {
  ESTIMATED: "ESTIMATED",
  CONFIRMED: "CONFIRMED",
  INVOICED: "INVOICED",
  PAID: "PAID",
}

// ─── Freight models ─────────────────────────────────────────────────────────

model FreightCreateRequest {
  shipment_uid: uuid;
  carrier_uid: uuid;
  @minValue(0)
  base_cost: float64;
  @minValue(0)
  weight_surcharge?: float64;
  @minValue(0)
  distance_surcharge?: float64;
  @minValue(0)
  discount?: float64;
  @pattern("^[A-Z]{3}$")
  currency?: string;
}

model FreightUpdateStatusRequest {
  status: FreightStatus;
}

model FreightItem {
  uid: uuid;
  shipment_uid: uuid;
  carrier_uid: uuid;
  base_cost: float64;
  weight_surcharge: float64;
  distance_surcharge: float64;
  discount: float64;
  total_cost: float64;
  currency: string;
  status: FreightStatus;
  invoiced_at?: utcDateTime;
  paid_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model FreightFindResponse is FreightItem;

model FreightUpdateResponse is FreightItem;

model FreightFindAllResponse is PaginatedResponse<FreightItem>;
