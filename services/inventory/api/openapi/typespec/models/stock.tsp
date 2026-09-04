namespace InventoryService;

// ─── Stock change type ──────────────────────────────────────────────────────

enum StockChangeType {
  INBOUND,
  OUTBOUND,
  ADJUSTMENT,
  RESERVATION,
  RELEASE,
}

// ─── Stock models ───────────────────────────────────────────────────────────

model StockWarehouseItem {
  warehouse_uid: uuid;
  warehouse_name: string;
  quantity: int32;
}

model StockFindResponse {
  uid: uuid;
  edition_uid: uuid;
  total_quantity: int32;
  warehouses: StockWarehouseItem[];
}

model StockFindAllResponse is PaginatedResponse<StockFindResponse>;

// ─── Stock operation requests ───────────────────────────────────────────────

@doc("Receive stock into a warehouse (inbound)")
model StockReceiveRequest {
  edition_uid: uuid;
  warehouse_uid: uuid;

  @minValue(1)
  quantity: int32;

  @maxLength(255)
  reason?: string;
}

@doc("Release stock from a warehouse (outbound)")
model StockReleaseRequest {
  edition_uid: uuid;
  warehouse_uid: uuid;

  @minValue(1)
  quantity: int32;

  @maxLength(255)
  reason?: string;
}

@doc("Adjust stock quantity (correction)")
model StockAdjustRequest {
  edition_uid: uuid;
  warehouse_uid: uuid;

  @doc("Positive or negative quantity adjustment")
  quantity: int32;

  @maxLength(255)
  reason: string;
}

@doc("Reserve stock for an order")
model StockReserveRequest {
  edition_uid: uuid;
  warehouse_uid: uuid;

  @minValue(1)
  quantity: int32;

  @doc("Order reference")
  @maxLength(255)
  reason?: string;
}

// ─── Stock history models ───────────────────────────────────────────────────

model StockHistoryFindResponse {
  uid: uuid;
  change_type: StockChangeType;
  reason?: string;
  change_qty: int32;
  created_at: utcDateTime;
}

model StockHistoryFindAllResponse is PaginatedResponse<StockHistoryFindResponse>;
