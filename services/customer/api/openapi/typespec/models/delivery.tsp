namespace CustomerService;

// ─── Delivery / shipment models (proxied from delivery service) ──────────────

enum ShipmentStatus {
  UNKNOWN,
  PLANNED,
  DISPATCHED,
  PICKED_UP,
  IN_TRANSIT,
  OUT_FOR_DELIVERY,
  DELIVERED,
  FAILED,
  CANCELLED,
}

model ShipmentResponse {
  uid: uuid;
  order_uid?: uuid;
  carrier_uid?: uuid;
  status?: ShipmentStatus;
  tracking_number?: string;
  destination_city?: string;
  destination_state?: string;
  destination_country_code?: string;
  destination_postal_code?: string;
  planned_delivery_at?: string;
  actual_delivery_at?: string;
  created_at?: string;
  updated_at?: string;
}

model ShipmentFindAllResponse is PaginatedResponse<ShipmentResponse>;
