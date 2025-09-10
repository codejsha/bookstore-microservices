namespace InventoryService;

// ─── Closing enums ─────────────────────────────────────────────────────────

enum ClosingStatus {
  OPEN,
  CLOSED,
}

// ─── Closing models ────────────────────────────────────────────────────────

model ClosingCreateRequest {
  warehouse_uid: uuid;
  year: int32;
  month: int32;
}

model ClosingItemResponse {
  edition_uid: uuid;
  opening_quantity: int32;
  inbound_quantity: int32;
  outbound_quantity: int32;
  adjust_quantity: int32;
  closing_quantity: int32;
}

model ClosingFindResponse {
  uid: uuid;
  warehouse_uid: uuid;
  year: int32;
  month: int32;
  status: ClosingStatus;
  items: ClosingItemResponse[];
  closed_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model ClosingFindAllResponse is PaginatedResponse<ClosingFindResponse>;
