namespace DeliveryService;

// ─── Shipment enums ─────────────────────────────────────────────────────────

enum ShipmentStatus {
  PLANNED: "PLANNED",
  DISPATCHED: "DISPATCHED",
  PICKED_UP: "PICKED_UP",
  IN_TRANSIT: "IN_TRANSIT",
  OUT_FOR_DELIVERY: "OUT_FOR_DELIVERY",
  DELIVERED: "DELIVERED",
  FAILED: "FAILED",
  CANCELLED: "CANCELLED",
}

// ─── Shipment models ────────────────────────────────────────────────────────

model ShipmentCreateRequest {
  order_uid: uuid;
  carrier_uid?: uuid;
  origin_address: string;
  destination_address: string;
  destination_city: string;
  destination_state: string;
  destination_country_code: string;
  destination_postal_code: string;
  weight_kg?: float64;
  planned_pickup_at?: utcDateTime;
  planned_delivery_at?: utcDateTime;
}

model ShipmentAssignCarrierRequest {
  carrier_uid: uuid;
  tracking_number?: string;
}

model ShipmentItem {
  uid: uuid;
  order_uid: uuid;
  carrier_uid?: uuid;
  origin_address: string;
  destination_address: string;
  destination_city: string;
  destination_state: string;
  destination_country_code: string;
  destination_postal_code: string;
  status: ShipmentStatus;
  tracking_number?: string;
  weight_kg?: float64;
  planned_pickup_at?: utcDateTime;
  planned_delivery_at?: utcDateTime;
  actual_pickup_at?: utcDateTime;
  actual_delivery_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model ShipmentFindResponse is ShipmentItem;

model ShipmentUpdateResponse is ShipmentItem;

model ShipmentFindAllResponse is PaginatedResponse<ShipmentItem>;

// ─── Tracking models ────────────────────────────────────────────────────────

model TrackingCreateRequest {
  status: ShipmentStatus;
  location: string;
  description: string;
}

model TrackingItem {
  uid: uuid;
  shipment_uid: uuid;
  status: ShipmentStatus;
  location: string;
  description: string;
  occurred_at: utcDateTime;
  created_at: utcDateTime;
}

model TrackingFindAllResponse {
  @minItems(0)
  items: TrackingItem[];
}
