namespace DeliveryService;

// ─── Carrier enums ──────────────────────────────────────────────────────────

enum CarrierStatus {
  ACTIVE: "ACTIVE",
  INACTIVE: "INACTIVE",
}

// ─── Carrier models ─────────────────────────────────────────────────────────

model CarrierCreateRequest {
  name: string;
  code: string;
  contact_name?: string;
  contact_phone?: string;
  contact_email?: string;
  base_rate: float64;
  rate_per_kg: float64;
}

model CarrierUpdateRequest {
  name?: string;
  contact_name?: string;
  contact_phone?: string;
  contact_email?: string;
  base_rate?: float64;
  rate_per_kg?: float64;
  status?: CarrierStatus;
}

model CarrierItem {
  uid: uuid;
  name: string;
  code: string;
  contact_name?: string;
  contact_phone?: string;
  contact_email?: string;
  base_rate: float64;
  rate_per_kg: float64;
  status: CarrierStatus;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model CarrierFindResponse is CarrierItem;

model CarrierUpdateResponse is CarrierItem;

model CarrierFindAllResponse is PaginatedResponse<CarrierItem>;
