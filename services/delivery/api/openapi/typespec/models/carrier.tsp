namespace DeliveryService;

// ─── Carrier enums ──────────────────────────────────────────────────────────

enum CarrierStatus {
  ACTIVE: "ACTIVE",
  INACTIVE: "INACTIVE",
}

// ─── Carrier models ─────────────────────────────────────────────────────────

model CarrierCreateRequest {
  @minLength(1)
  @maxLength(255)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  name: string;
  @minLength(1)
  @maxLength(50)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  code: string;
  @minLength(1)
  @maxLength(255)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  contact_name?: string;
  @minLength(1)
  @maxLength(50)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  contact_phone?: string;
  @maxLength(255)
  @pattern("^[^@\\s]+@[^@\\s]+\\.[^@\\s]+$")
  contact_email?: string;
  @minValue(0)
  base_rate: float64;
  @minValue(0)
  rate_per_kg: float64;
}

model CarrierUpdateRequest {
  @minLength(1)
  @maxLength(255)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  name?: string;
  @minLength(1)
  @maxLength(255)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  contact_name?: string;
  @minLength(1)
  @maxLength(50)
  @pattern("[\\s\\S]*\\S[\\s\\S]*")
  contact_phone?: string;
  @maxLength(255)
  @pattern("^[^@\\s]+@[^@\\s]+\\.[^@\\s]+$")
  contact_email?: string;
  @minValue(0)
  base_rate?: float64;
  @minValue(0)
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
