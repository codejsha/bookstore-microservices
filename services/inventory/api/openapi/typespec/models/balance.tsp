namespace InventoryService;

// ─── Balance models ────────────────────────────────────────────────────────

model BalanceFindResponse {
  edition_uid: uuid;
  warehouse_uid: uuid;
  warehouse_name: string;
  inbound_quantity: int32;
  outbound_quantity: int32;
  adjust_quantity: int32;
  current_quantity: int32;
}

model BalanceFindAllResponse is PaginatedResponse<BalanceFindResponse>;
