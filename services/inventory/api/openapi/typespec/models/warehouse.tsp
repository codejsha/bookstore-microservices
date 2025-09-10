namespace InventoryService;

// ─── Warehouse models ───────────────────────────────────────────────────────

model WarehouseCreateRequest {
  name: string;
  address?: string;

  @minValue(0)
  capacity: int32;
}

model WarehouseUpdateRequest {
  name?: string;
  address?: string;

  @minValue(0)
  capacity?: int32;
}

model WarehouseFindResponse {
  uid: uuid;
  name: string;
  address?: string;
  capacity: int32;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model WarehouseUpdateResponse is WarehouseFindResponse;

model WarehouseFindAllResponse is PaginatedResponse<WarehouseFindResponse>;
