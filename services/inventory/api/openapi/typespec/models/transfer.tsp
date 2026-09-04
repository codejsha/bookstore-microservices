namespace InventoryService;

// ─── Transfer enums ────────────────────────────────────────────────────────

enum TransferStatus {
  PENDING,
  IN_TRANSIT,
  COMPLETED,
  CANCELLED,
}

// ─── Transfer models ───────────────────────────────────────────────────────

model TransferCreateRequest {
  edition_uid: uuid;
  source_warehouse_uid: uuid;
  target_warehouse_uid: uuid;

  @minValue(1)
  quantity: int32;

  @maxLength(255)
  reason?: string;
}

model TransferFindResponse {
  uid: uuid;
  edition_uid: uuid;
  source_warehouse_uid: uuid;
  target_warehouse_uid: uuid;
  quantity: int32;
  status: TransferStatus;
  reason?: string;
  completed_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model TransferFindAllResponse is PaginatedResponse<TransferFindResponse>;
