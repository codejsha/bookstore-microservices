import "@typespec/http";

using Http;

namespace AdminService;

// ─── Warehouses ─────────────────────────────────────────────────────────────

model AdminWarehouseResponse {
  uid: uuid;
  name: string;
  address?: string;

  @doc("Units the warehouse can hold, as recorded by inventory")
  capacity: int32;

  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model AdminWarehouseFindAllResponse is PaginatedResponse<AdminWarehouseResponse>;

// ─── Stock ──────────────────────────────────────────────────────────────────

@doc("How much of an edition one warehouse holds")
model AdminStockWarehouseItem {
  warehouse_uid: uuid;
  warehouse_name: string;
  quantity: int32;
}

@doc("Stock for a single edition, broken down by warehouse")
model AdminStockResponse {
  uid: uuid;
  edition_uid: uuid;

  @doc("The sum of the per-warehouse quantities")
  total_quantity: int32;

  warehouses: AdminStockWarehouseItem[];
}

model AdminStockFindAllResponse is PaginatedResponse<AdminStockResponse>;
