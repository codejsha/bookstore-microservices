namespace InventoryService;

// ─── Audit enums ───────────────────────────────────────────────────────────

enum AuditStatus {
  IN_PROGRESS,
  COMPLETED,
}

// ─── Audit models ──────────────────────────────────────────────────────────

model AuditItemRequest {
  edition_uid: uuid;
  actual_quantity: int32;
}

model AuditCreateRequest {
  warehouse_uid: uuid;
  items: AuditItemRequest[];
  notes?: string;
}

model AuditItemResponse {
  edition_uid: uuid;
  system_quantity: int32;
  actual_quantity: int32;
  difference: int32;
}

model AuditFindResponse {
  uid: uuid;
  warehouse_uid: uuid;
  status: AuditStatus;
  items: AuditItemResponse[];
  notes?: string;
  completed_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model AuditFindAllResponse is PaginatedResponse<AuditFindResponse>;
